package dao

import (
	"context"
	"log"
	"math"
	"replite_web/internal/app/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

/**

the ip model only exist the redis to defind more user to access the system in the same time

**/

const DEFAULT_IP_EXPIRE_TIME = 1 * time.Minute

var REDIS_INSERT_IP_FAILED = 1 << (utils.GetOperationBit() - 1)

var REDIS_EMPTY_IP = math.MinInt

// whether insert the key and value is unnecessary operation
func InsertIP(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_IP_EXPIRE_TIME)
	defer cancel()
	// sctx, scancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer scancel()
	// realKey := getRegisterKey(ip)
	//每次增加的时候都会重新设置过期时间
	tx := GetRedisClient().TxPipeline()
	tx.Incr(ctx, ip)
	tx.Expire(ctx, ip, 2*time.Minute)
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Printf("插入redis缓存(%s)失败:%s", ip, err.Error())
		return false
		// return REDIS_INSERT_IP_FAILED
	}
	return true
	// result, err := cmd.Uint64()
	// if err != nil {
	// 	return REDIS_INSERT_IP_FAILED
	// }
	// return int(result)
}

// func IncreaseIPTimes(ip string) (int, bool) {
// 	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_IP_EXPIRE_TIME)
// 	defer cancel()
// 	cmd := GetRedisClient().Incr(ctx, ip)
// 	if cmd.Err() != nil {
// 		if cmd.Err() == redis.Nil {
// 			return REDIS_EMPTY_IP, true
// 		}
// 		log.Printf("redis缓存查询(%s)失败:%s", ip, cmd.Err().Error())
// 		return REDIS_INSERT_IP_FAILED, false
// 	}
// 	val, err := cmd.Uint64()
// 	if err != nil {
// 		log.Printf("redis的类型缓存出错(%s):%s", ip, err.Error())
// 		return REDIS_INSERT_IP_FAILED, false
// 	}
// 	return int(val), true
// }

func QueryIP(ip string) int {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_IP_EXPIRE_TIME)
	defer cancel()
	// realKey := getRegisterKey(ip)
	cmd := GetRedisClient().Get(ctx, ip)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return REDIS_EMPTY_IP
		}
		return REDIS_INSERT_IP_FAILED
	}
	result, err := cmd.Int()
	if err != nil {
		log.Printf("redis缓存的类型出错(%s):%s", ip, err.Error())
		return REDIS_INSERT_IP_FAILED
	}
	return result
}
