package dao

import (
	"context"
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFUALT_QUERYS_USER_NUMBER = 10

const DEFAULT_REDIS_USER_PREFIX = "user-"

const DEFAULT_REDIS_USERS_PREFIX = "users-"

var INVALID_REDIS_USER_VALUE = User{}

const DEFAULT_USER_EXPIRE_TIME = 3 * time.Minute

var INVALID_REDIS_USERS_VALUE = struct{}{}
var emptyUser = User{}

type User struct {
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	Authority   string `json:"athority" bson:"authority"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	Code        string `json:"-" bson:"-"`
}

const DEFAULT_USER_COLLECTION = "user"

func getUserCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(DEFAULT_USER_COLLECTION).(string))
}
func CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := getUserCollection().InsertOne(ctx, user)
	if err != nil {
		log.Printf("创建document出错:%s", err.Error())
		return err
	}
	//delete the cache invalid key
	Create(user.Username, user, DEFAULT_USER_EXPIRE_TIME)
	return err
}
func UpdateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//to keep the operation atomic
	result := getUserCollection().FindOneAndUpdate(ctx, bson.M{"username": user.Username}, user)
	if result.Err() != nil {
		log.Printf("修改document异常:%s", result.Err().Error())
		return result.Err()
	}
	updateUser := new(User)
	err := result.Decode(updateUser)
	if err != nil {
		log.Printf("解析mongoDB修改后的document异常:%s", err.Error())
		return err
	}
	err = Create(user.Username, updateUser, DEFAULT_USER_EXPIRE_TIME)
	if err != nil {
		log.Printf("创建redis缓存(%v)失败:%s", updateUser, err.Error())
	}
	return err
}

func DeleteUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//del the mongo data
	result, err := getUserCollection().DeleteOne(ctx, bson.M{"username": user.Username})
	if err != nil && result.DeletedCount <= 0 {
		log.Printf("删除User document出错,影响document数量%d:%s", result.DeletedCount, err.Error())
		return err
	}
	// del redis key
	err = Del(getUserKey(user.Username))
	if err != nil {
		log.Printf("删除redis缓存(%s)出错:%s", user.Username, err.Error())
	}
	return err
}
func QueryUser(user *User) (User, error) {
	//query for redis cache
	var model = User{}
	redisKey := getUserKey(user.Username)
	err := Get(redisKey, &model)
	// defend the invalid key to access the mongoDB
	if !model.IsEmpty() || (err != nil && err != redis.Nil) {
		log.Printf("查询(%s)缓存失败：%s", redisKey, err.Error())
		return model, err
	}
	// if err != redis.Nil {
	// 	log.Printf("查询缓存失败：%s", err.Error())
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// query for mongo database
	result := getUserCollection().FindOne(ctx, bson.M{"username": user.Username})
	if result.Err() != nil {
		// defend the cache breakdown
		if result.Err() == mongo.ErrNoDocuments {
			// invalid key expire time be set
			//TODO  the query user will appear in the register stage,this situation will take much redis appear
			Create(redisKey, INVALID_REDIS_USER_VALUE, 1*time.Minute)
			return User{}, nil
		}
		return User{}, result.Err()
	}
	err = result.Decode(&model)
	//keep cache
	if err == nil {
		err = Create(redisKey, &model, DEFAULT_USER_EXPIRE_TIME)
		if err != nil {
			log.Printf("redis cache the user info %v error:%s", model, err.Error())
		}
	}
	return model, err
}

func QueryUsers(page int, pageNumber int) ([]*User, error) {
	//redis cache
	users := make([]*User, DEFUALT_QUERYS_USER_NUMBER)
	redisKey := getUsersKey(page, pageNumber)
	err := GetList(redisKey, users, 0, -1)
	if err == nil {
		return users, nil
	}
	if err != redis.Nil {
		log.Printf("redis查询缓存(%s)失败:%s", redisKey, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//get users in mongo database
	result, err := getUserCollection().Find(ctx, bson.D{}, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page)-1))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Create(redisKey, INVALID_REDIS_USERS_VALUE, 1*time.Minute)
		}
		return nil, err
	}
	defer result.Close(context.Background())
	result.All(context.Background(), users)
	err = Create(redisKey, users, DEFAULT_USER_EXPIRE_TIME)
	if err != nil {
		log.Printf("创建redis缓存(key:%s,value:%v)失败%s", redisKey, users, err.Error())
	}
	return users, nil
}

func (user User) IsEmpty() bool {
	return user == emptyUser
}

/*
create the unique id for user in redis
*/
func getUserKey(username string) string {
	return utils.MergeStr(DEFAULT_REDIS_USER_PREFIX, username)
}

/* create the unique id for user list in redis*/
func getUsersKey(page int, pageNumber int) string {
	return utils.MergeStr(DEFAULT_REDIS_USERS_PREFIX, strconv.FormatInt(int64(page), 10), "-", strconv.FormatInt(int64(pageNumber), 10))
}
