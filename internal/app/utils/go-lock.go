package utils

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

//取锁,加时间,还锁

type Mutex struct {
	delegate *redis.Client
	// auto delay
	delayDone   chan struct{}
	value       string
	cancelTimes time.Duration
	name        string
}

var (
	mutex *Mutex
	once  sync.Once
)

func NewMutex(storage redis.Client) (mutex *Mutex, err error) {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		mutex = new(Mutex)
		err = storage.Ping(ctx).Err()
		if err != nil {
			return
		}
		mutex.delegate = &storage
		mutex.value, err = getMachineID()
	})
	return mutex, err
}

func (mutex *Mutex) Lock() {

}

func (mutex *Mutex) TryLock() {

}

func (mutex *Mutex) acquire() (bool, error) {

}

// -1 resprent the delay error
const delayScript = `
   if (redis.call('GET',KEYS[1])==ARGS[1])then
          return redis.call('TTL',ARGS[2])
	else
	      return -1
	done	  	   
`

func (mutex *Mutex) delay() {

}

const releaseScript = `
   if(redis.call('GET',KEYS[1])==ARGS[1])then
           return redis.call('DEL',KEYS[1])
	done	   
`

var cacheHash = make(map[string]string)

const releaseCacheKey = "release"
const deleteCacheKey = "delete"

func (mutex *Mutex) release() error {
	ctx, cancel := context.WithTimeout(context.TODO(), mutex.cancelTimes)
	defer cancel()
	if _, ok := cacheHash[releaseCacheKey]; !ok {
		cmd := mutex.delegate.Eval(ctx, releaseScript, []string{mutex.name}, mutex.value)
		hash := sha256.Sum256([]byte(releaseScript))
		cacheHash[releaseCacheKey] = string(hash[:])
		return cmd.Err()
	}
	return mutex.delegate.EvalSha(ctx, cacheHash[releaseCacheKey], []string{mutex.name}, mutex.value).Err()
}
func getMachineID() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {

		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			// 查找第一个有效的MAC地址
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					mac := iface.HardwareAddr.String()
					if mac != "" {
						return mac, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("无法获取机器的唯一标识")
}
