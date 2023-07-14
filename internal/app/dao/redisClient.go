package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	var redisConfig = config.DBConfig.RedisConfig
	db, err := strconv.ParseInt(redisConfig.Database, 10, utils.GetOperationBit())
	if err != nil {
		panic(fmt.Sprintf("convert %s to int error:%v", redisConfig.Database, err))
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Username: redisConfig.Username,
		Password: redisConfig.Passowrd,
		DB:       int(db),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis %v ping error:%s", redisConfig, err.Error()))
	}
	log.Printf("Connected to Redis")
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetStr(key string) (string, error) {
	cmd := get(key)
	return cmd.Val(), cmd.Err()
}
func Get(key string, value any) error {
	cmd := get(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return json.Unmarshal([]byte(cmd.Val()), value)
}

func get(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return GetRedisClient().Get(ctx, key)
}

func GetStrList(key string, start int, end int) ([]string, error) {
	cmd := getList(key, start, end)
	return cmd.Result()
}

func GetList(key string, value any, start int, end int) error {
	cmd := getList(key, start, end)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	val, err := cmd.Result()
	if err != nil {
		return err
	}
	if len(val) <= 0 {
		return redis.Nil
	}
	rValue := reflect.ValueOf(value)
	if rValue.Type().Kind() != reflect.Slice {
		return errors.New("value 应该是一个slice")
	}
	if rValue.Len() <= len(val) {
		rValue = reflect.MakeSlice(rValue.Type(), len(val), len(val))
	}
	//************** TODO should inspect **********
	for i := 0; i < len(val); i++ {
		err = json.Unmarshal([]byte(val[i]), rValue.Index(i).Interface())
		if err != nil {
			return err
		}
	}
	return nil
}
func getList(key string, start int, end int) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return GetRedisClient().LRange(ctx, key, int64(start), int64(end))
}

func Del(key ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := GetRedisClient().Del(ctx, key...)
	return cmd.Err()
}

// func create(key string, value any, expire time.Duration) *redis.StatusCmd {
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	return GetRedisClient().SetEx(ctx, key, value, expire)
// }

func CreateList(key string, value any, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// sctx, scancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer scancel()
	// start transaction
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// valSlice, ok := value.([]any)
	// if !ok {
	// 	return errors.New("createList method: the value which is input  should be slice type")
	// }
	tx := GetRedisClient().TxPipeline()
	tx.LPush(ctx, key, string(str))
	tx.Expire(ctx, key, expire)
	_, err = tx.Exec(ctx)
	return err
}

// if value is str,should appear the ""
func Create(key string, value any, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var str string
	if valueStr, ok := value.(string); ok {
		str = valueStr
	} else {
		valueStr, err := json.Marshal(value)
		if err != nil {
			return err
		}
		str = string(valueStr)
	}
	return GetRedisClient().SetEx(ctx, key, str, expire).Err()
}
