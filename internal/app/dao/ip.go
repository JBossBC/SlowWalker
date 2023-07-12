package dao

import (
	"context"
	"log"
	"replite_web/internal/app/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

/**

the ip model only exist the redis to defind more user to access the system in the same time

**/

// TODO this inspection is too bad to effect the system high performance effort
// to defind the register more account from ip
const Register_FAILED_TIMES_PREDIXX = "register-failed-"

const DEFAULT_IP_EXPIRE_TIME = 1 * time.Minute

var REDIS_INSERT_IP_FAILED = 1 << (utils.GetOperationBit() - 1)

// whether insert the key and value is unnecessary operation
func InsertIP(ip string) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_IP_EXPIRE_TIME)
	defer cancel()
	// sctx, scancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer scancel()
	// realKey := getRegisterKey(ip)
	tx := GetRedisClient().TxPipeline()
	tx.Incr(ctx, ip)
	tx.Expire(ctx, ip, 2*time.Minute)
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Printf("查询redis缓存(%s)失败:%s", ip, err.Error())
		// return REDIS_INSERT_IP_FAILED
	}
	// result, err := cmd.Uint64()
	// if err != nil {
	// 	return REDIS_INSERT_IP_FAILED
	// }
	// return int(result)

}

func QueryIP(ip string) int {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_IP_EXPIRE_TIME)
	defer cancel()
	// realKey := getRegisterKey(ip)
	cmd := GetRedisClient().Get(ctx, ip)
	err := cmd.Err()
	if err != nil && err != redis.Nil {
		return REDIS_INSERT_IP_FAILED
	}
	result, err := cmd.Int()
	if err != nil {
		return REDIS_INSERT_IP_FAILED
	}
	return result
}
