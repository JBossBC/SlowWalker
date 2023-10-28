package dao

import (
	"context"
	"log"
	"replite_web/internal/app/cache"
	"replite_web/internal/app/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DepartmentInfo struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	CreateTime  int64    `json:"createTime" bson:"createTime"`
	Leaders     []string `json:"leaders" bson:"leaders"`
}
type DepartmentDao struct {
}

var (
	departmentDao     *DepartmentDao
	departmentDaoOnce sync.Once
)

func getDepartmentDao() *DepartmentDao {
	departmentDaoOnce.Do(func() {
		departmentDao = new(DepartmentDao)
	})
	return departmentDao
}

var (
	departmentCollection     *mongo.Collection
	departmentCollectionOnce sync.Once
)

const DEFAULT_DEPARTMENT_COLLECTION = "department"

func getDepartmentCollection() *mongo.Collection {
	departmentCollectionOnce.Do(func() {
		departmentCollection = getMongoConn().Collection(config.GetCollectionConfig().Get(DEFAULT_DEPARTMENT_COLLECTION).(string))
	})
	return departmentCollection
}
func (department *DepartmentDao) CreateDepartment(departmentInfo DepartmentInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	_, err := getDepartmentCollection().InsertOne(ctx, departmentInfo)
	if err != nil {
		log.Printf("创建document出错:%s", err.Error())
		return err
	}
	cache.GetCachePool().Store(departmentInfo.Name, departmentInfo, 3*time.Minute)
	return err
}

// 保证一致性
func (departmentDao *DepartmentDao) UpdateDepartment(departmentInfo DepartmentInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// redisKey := getUserKey(user.Username)
	// err := Del(redisKey)
	// if err != nil && err != redis.Nil {
	// 	log.Printf("删除redis缓存(%s)失败:%s", redisKey, err.Error())
	// 	return errors.New("修改失败")
	// }
	//to keep the operation atomic
	result := getDepartmentCollection().FindOneAndUpdate(ctx, bson.M{"name": departmentInfo.Name}, departmentInfo)
	if result.Err() != nil {
		log.Printf("修改document异常:%s", result.Err().Error())
		return result.Err()
	}
	updateDepartment := new(DepartmentInfo)
	err := result.Decode(updateDepartment)
	if err != nil {
		log.Printf("解析mongoDB修改后的document异常:%s", err.Error())
		return err
	}
	cache.GetCachePool().Store(departmentInfo.Name, *updateDepartment, 3*time.Minute)
	// err = Create(redisKey, updateUser, DEFAULT_USER_EXPIRE_TIME)
	// if err != nil {
	// 	log.Printf("创建redis缓存(key:%s,value:%v)失败:%s", redisKey, updateUser, err.Error())
	// }
	return err
}

func (departmentDao *DepartmentDao) DeleteDepartment(departmentInfo *DepartmentInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//del the mongo data
	result, err := getDepartmentCollection().DeleteOne(ctx, bson.M{"name": departmentInfo})
	if err != nil && result.DeletedCount <= 0 {
		log.Printf("删除Department document出错,影响document数量%d:%s", result.DeletedCount, err.Error())
		return err
	}
	cache.GetCachePool().Delete(departmentInfo.Name)
	// del redis key,should defer delete if concurrency scene
	// err = Del(getUserKey(user.Username))
	// if err != nil {
	// 	log.Printf("删除redis缓存(%s)出错:%s", user.Username, err.Error())
	// }
	return err
}
func (departmentDao *DepartmentDao) QueryDepartment(department *DepartmentInfo) (DepartmentInfo, error) {
	//query for redis cache
	cacheModel, ok := cache.GetCachePool().TryGet(department.Name)
	if ok {
		return cacheModel.(DepartmentInfo), nil
	}
	var model = DepartmentInfo{}
	// redisKey := getUserKey(user.Username)
	// err := Get(redisKey, &model)
	// // defend the invalid key to access the mongoDB
	// if !model.IsEmpty() {
	// 	return model, nil
	// }
	// if err != nil && err != redis.Nil {
	// 	log.Printf("查询(%s)缓存失败：%v", redisKey, err)
	// }
	// if err != redis.Nil {
	// 	log.Printf("查询缓存失败：%s", err.Error())
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// query for mongo database
	result := getDepartmentCollection().FindOne(ctx, bson.M{"name": department.Name})
	if result.Err() != nil {
		// defend the cache breakdown
		// if result.Err() == mongo.ErrNoDocuments {
		// 	// invalid key expire time be set
		// 	//TODO  the query user will appear in the register stage,this situation will take much redis key appear
		// 	Create(redisKey, INVALID_REDIS_USER_VALUE, 1*time.Minute)
		// 	return UserInfo{}, nil
		// }
		return DepartmentInfo{}, result.Err()
	}
	err := result.Decode(&model)
	//keep cache
	if err == nil {
		cache.GetCachePool().Store(model.Name, model, 3*time.Minute)
	}
	// if err == nil {
	// 	err = Create(redisKey, &model, DEFAULT_USER_EXPIRE_TIME)
	// 	if err != nil {
	// 		log.Printf("redis cache the user info %v error:%s", model, err.Error())
	// 	}
	// }
	return model, err
}

const ALL_DEPARTMENT_KEYS = "all_departments"

const DEFAULT_QUERYS_DEPARTMENT_NUMBER = 10

// in the scene, actually the query operation is query all departments

func (departmentDao *DepartmentDao) QueryDepartments() ([]*DepartmentInfo, error) {
	//redis cache
	cacheDepartments, ok := cache.GetCachePool().TryGet(ALL_DEPARTMENT_KEYS)
	if ok {
		return cacheDepartments.([]*DepartmentInfo), nil
	}
	departments := make([]*DepartmentInfo, 0, DEFAULT_QUERYS_DEPARTMENT_NUMBER)
	// redisKey := getUsersKey(page, pageNumber)
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
	result, err := getDepartmentCollection().Find(ctx, bson.M{})
	if err != nil {
		// if err == mongo.ErrNoDocuments {
		// 	Create(redisKey, INVALID_REDIS_USERS_VALUE, 1*time.Minute)
		// }
		return nil, err
	}
	defer result.Close(context.Background())
	result.All(context.Background(), &departments)
	// err = CreateList(redisKey, users, DEFAULT_USER_EXPIRE_TIME)
	// if err != nil {
	// 	log.Printf("创建redis缓存(key:%s,value:%v)失败%s", redisKey, users, err.Error())
	// }
	cache.GetCachePool().Store(ALL_DEPARTMENT_KEYS, departments, 3*time.Minute)
	return departments, nil
}

// func (user UserInfo) IsEmpty() bool {
// 	return user == emptyUser
// }

// /*
// create the unique id for user in redis
// */
// func getDepartmentKey(username string) string {
// 	return utils.MergeStr(DEFAULT_REDIS_USER_PREFIX, username)
// }

// /* create the unique id for user list in redis*/
// func getUsersKey(page int, pageNumber int) string {
// 	return utils.MergeStr(DEFAULT_REDIS_USERS_PREFIX, strconv.FormatInt(int64(page), 10), "-", strconv.FormatInt(int64(pageNumber), 10))
// }
