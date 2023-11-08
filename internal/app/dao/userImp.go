package dao

import (
	"context"
	"errors"
	"fmt"
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFUALT_QUERYS_USER_NUMBER = 10

const DEFAULT_REDIS_USER_PREFIX = "user-"

const DEFAULT_REDIS_USERS_PREFIX = "users-"

// TODO this field has exist for emptyUser
var INVALID_REDIS_USER_VALUE = UserInfo{}

const DEFAULT_USER_EXPIRE_TIME = 5 * time.Minute

var INVALID_REDIS_USERS_VALUE = struct{}{}
var emptyUser = UserInfo{}

// TODO increase  the create time field
type UserInfo struct {
	Username    string `json:"username" bson:"username"`
	Password    string `json:"-" bson:"password"`
	Authority   string `json:"authority" bson:"authority"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	RealName    string `json:"realName" bson:"realName"`
	//only user operation has value for department,the department field includes two means. first: what the users belong to. second: what the user display role in the department
	Department string `json:"department" bson:"department"`
	CreateTime int64  `json:"createTime" bson:"createTime"`
	Code       string `json:"-" bson:"-"`
	IP         string `json:"-" bson:"-"`
	Page       int    `json:"-" bson:"-"`
	PageNumber int    `json:"-" bson:"-"`
	Start      int64  `json:"-" bson:"-"`
	End        int64  `json:"-" bson:"-"`
}

type UserDao struct {
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

func getUserDao() *UserDao {
	userOnce.Do(func() {
		userDao = new(UserDao)
	})
	return userDao
}

const DEFAULT_USER_COLLECTION = "user"

var (
	userCollection     *mongo.Collection
	userCollectionOnce sync.Once
)

func getUserCollection() *mongo.Collection {
	userCollectionOnce.Do(func() {
		userCollection = getMongoConn().Collection(config.GetCollectionConfig().Get(DEFAULT_USER_COLLECTION).(string))
	})
	return userCollection
}
func (userDao *UserDao) CreateUser(user *UserInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	_, err := getUserCollection().InsertOne(ctx, user)
	if err != nil {
		log.Printf("创建document出错:%s", err.Error())
		return err
	}
	//delete the cache invalid key
	Create(getUserKey(user.Username), user, DEFAULT_USER_EXPIRE_TIME)
	return err
}

// 保证一致性
func (userDao *UserDao) UpdateUser(user *UserInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisKey := getUserKey(user.Username)
	err := Del(redisKey)
	if err != nil && err != redis.Nil {
		log.Printf("删除redis缓存(%s)失败:%s", redisKey, err.Error())
		return errors.New("修改失败")
	}
	//to keep the operation atomic
	result := getUserCollection().FindOneAndUpdate(ctx, bson.M{"username": user.Username}, user)
	if result.Err() != nil {
		log.Printf("修改document异常:%s", result.Err().Error())
		return result.Err()
	}
	updateUser := new(UserInfo)
	err = result.Decode(updateUser)
	if err != nil {
		log.Printf("解析mongoDB修改后的document异常:%s", err.Error())
		return err
	}
	err = Create(redisKey, updateUser, DEFAULT_USER_EXPIRE_TIME)
	if err != nil {
		log.Printf("创建redis缓存(key:%s,value:%v)失败:%s", redisKey, updateUser, err.Error())
	}
	return err
}

func (userDao *UserDao) DeleteUser(user *UserInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//del the mongo data
	result, err := getUserCollection().DeleteOne(ctx, bson.M{"username": user.Username})
	if err != nil && result.DeletedCount <= 0 {
		log.Printf("删除User document出错,影响document数量%d:%s", result.DeletedCount, err.Error())
		return err
	}
	// del redis key,should defer delete if concurrency scene
	err = Del(getUserKey(user.Username))
	if err != nil {
		log.Printf("删除redis缓存(%s)出错:%s", user.Username, err.Error())
	}
	return err
}

// func (userDao *UserDao) QueryUserByUsername(user *UserInfo) (UserInfo, error) {
// 	//query for redis cache
// 	var model = UserInfo{}
// 	redisKey := getUserKey(user.Username)
// 	err := Get(redisKey, &model)
// 	// defend the invalid key to access the mongoDB
// 	if !model.IsEmpty() {
// 		return model, nil
// 	}
// 	if err != nil && err != redis.Nil {
// 		log.Printf("查询(%s)缓存失败：%v", redisKey, err)
// 	}
// 	// if err != redis.Nil {
// 	// 	log.Printf("查询缓存失败：%s", err.Error())
// 	// }
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	// query for mongo database
// 	result := getUserCollection().FindOne(ctx, bson.M{"username": user.Username})
// 	if result.Err() != nil {
// 		// defend the cache breakdown
// 		if result.Err() == mongo.ErrNoDocuments {
// 			// invalid key expire time be set
// 			//TODO  the query user will appear in the register stage,this situation will take much redis key appear
// 			Create(redisKey, INVALID_REDIS_USER_VALUE, 1*time.Minute)
// 			return UserInfo{}, nil
// 		}
// 		return UserInfo{}, result.Err()
// 	}
// 	err = result.Decode(&model)
// 	//keep cache
// 	if err == nil {
// 		err = Create(redisKey, &model, DEFAULT_USER_EXPIRE_TIME)
// 		if err != nil {
// 			log.Printf("redis cache the user info %v error:%s", model, err.Error())
// 		}
// 	}
// 	return model, err
// }

// func (userDao *UserDao) QueryUsersByDepartment(user *UserInfo) ([]*UserInfo, error) {
// 	//query for redis cache
// 	var model = make([]*UserInfo, 0, DEFUALT_QUERYS_USER_NUMBER)
// 	redisKey := getUserKey(user.Department)
// 	err := GetList(redisKey, &model, 0, int(user.PageNumber))
// 	// defend the invalid key to access the mongoDB
// 	// if len(model) <= 0 {
// 	// 	return model, nil
// 	// }
// 	if err != nil && err != redis.Nil {
// 		log.Printf("查询(%s)缓存失败：%v", redisKey, err)
// 	}
// 	// if err != redis.Nil {
// 	// 	log.Printf("查询缓存失败：%s", err.Error())
// 	// }
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	// query for mongo database
// 	result, err := getUserCollection().Find(ctx, bson.M{"department": user.Department})
// 	if err != nil {
// 		// defend the cache breakdown
// 		if result.Err() == mongo.ErrNoDocuments {
// 			// invalid key expire time be set
// 			//TODO  the query user will appear in the register stage,this situation will take much redis key appear
// 			CreateList(redisKey, INVALID_REDIS_USER_VALUE, 1*time.Minute)
// 			return nil, nil
// 		}
// 		return nil, result.Err()
// 	}
// 	defer result.Close(context.Background())
// 	err = result.All(context.Background(), &model)
// 	//keep cache
// 	if err == nil {
// 		err = CreateList(redisKey, &model, DEFAULT_USER_EXPIRE_TIME)
// 		if err != nil {
// 			log.Printf("redis cache the user info %v error:%s", model, err.Error())
// 		}
// 	}
// 	return model, err
// }

// TODO should be updated
func (userDao *UserDao) FilterUsers(userInfo *UserInfo) ([]*UserInfo, error) {
	filter := bson.M{}
	if userInfo.Username != "" {
		filter["username"] = userInfo.Username
	}
	if userInfo.RealName != "" {
		filter["realName"] = userInfo.RealName
	}
	if userInfo.Authority != "" {
		filter["authority"] = userInfo.Authority
	}
	if userInfo.Department != "" {
		filter["department"] = userInfo.Department
	}
	if userInfo.PhoneNumber != "" {
		filter["phoneNumber"] = userInfo.PhoneNumber
	}
	//redis cache
	users := make([]*UserInfo, 0, DEFUALT_QUERYS_USER_NUMBER)
	// redisKey := getUsersKey(userInfo.Page, userInfo.PageNumber)
	// err := GetList(redisKey, users, 0, -1)
	// if err == nil {
	// 	return users, nil
	// }
	// if err != redis.Nil {
	// 	log.Printf("redis查询缓存(%s)失败:%s", redisKey, err.Error())
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//get users in mongo database
	result, err := getUserCollection().Find(ctx, filter, options.Find().SetLimit(int64(userInfo.PageNumber)), options.Find().SetSkip(int64(userInfo.Page)-1))
	if err != nil {
		// if err == mongo.ErrNoDocuments {
		// 	Create(redisKey, INVALID_REDIS_USERS_VALUE, 1*time.Minute)
		// }
		return nil, err
	}
	defer result.Close(context.Background())
	result.All(context.Background(), &users)
	//hide secret
	// err = CreateList(redisKey, users, DEFAULT_USER_EXPIRE_TIME)
	// if err != nil {
	// 	log.Printf("创建redis缓存(key:%s,value:%v)失败%s", redisKey, users, err.Error())
	// }
	return users, nil
}

func (userDao *UserDao) QueryUsers(queryTempalte *UserInfo) ([]*UserInfo, error) {
	query := bson.M{}
	if queryTempalte.Username != "" {
		query["username"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", queryTempalte.Username),
			"$options": "i",
		}
	}
	if queryTempalte.RealName != "" {
		query["realName"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", queryTempalte.RealName),
			"$options": "i",
		}
	}
	if queryTempalte.Authority != "" {
		query["authority"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", queryTempalte.Authority),
			"$options": "i",
		}
	}
	if queryTempalte.Department != "" {
		query["department"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", queryTempalte.Department),
			"$options": "i",
		}
	}
	if queryTempalte.PhoneNumber != "" {
		query["phoneNumber"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", queryTempalte.PhoneNumber),
			"$options": "i",
		}
	}
	if queryTempalte.Start != 0 || queryTempalte.End != 0 {
		query["createTime"] = bson.M{
			"$gt": queryTempalte.Start,
			"$lt": queryTempalte.End,
		}
	}
	users := make([]*UserInfo, 0, DEFUALT_QUERYS_USER_NUMBER)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	result, err := getUserCollection().Find(ctx, query, options.Find().SetLimit(int64(queryTempalte.PageNumber)), options.Find().SetSkip(int64(queryTempalte.Page)-1))
	if err != nil {
		log.Printf("query the user(filter:%v) error:%s", queryTempalte, err.Error())
		return nil, err
	}
	defer result.Close(context.Background())
	err = result.All(context.Background(), &users)
	if err != nil {
		log.Printf("analysis the user collections by filterUser function error:%s ", err.Error())
		return nil, err
	}
	return users, nil
}

func (user UserInfo) IsEmpty() bool {
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
