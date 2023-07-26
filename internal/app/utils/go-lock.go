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

func init() {
	singleConfig.cancelTime = defaultCancelTime
	singleConfig.expiresTime = defaultExpiresTime
	singleConfig.maxOffsetTime = defaultMaxOffsetTime
	singleConfig.reties = defaultReties
	nodeID, err := getMachineID()
	if err != nil {
		panic(fmt.Sprintf("init the mutex unique nodeId false:%s", err.Error()))
	}
	singleConfig.nodeID = nodeID
}

const defaultCancelTime = 1 * time.Second

const defaultExpiresTime = 3 * time.Second

const defaultMaxOffsetTime = 10 * time.Millisecond
const defaultReties = 2

type Mutex struct {
	// auto delay
	delayDone chan struct{}
	name      string
	config    *config
	ending    chan error
}

type config struct {
	cancelTime    time.Duration
	maxOffsetTime time.Duration
	expiresTime   time.Duration
	reties        int
	delegate      *redis.Client
	nodeID        string
}

type ConfigOption func(*config)

// WithCancelTime a time for redis conn need spend the max time
func WithCancelTime(cancelTime time.Duration) ConfigOption {
	return func(c *config) {
		c.cancelTime = cancelTime
	}
}

// WithExpiresTime   Duration of lock
func WithExpiresTime(expireTime time.Duration) ConfigOption {
	return func(c *config) {
		c.expiresTime = expireTime
	}
}

// WithMaxOffsetTime be set to priority about the max gradient times for the interval of requesting lock
func WithMaxOffsetTime(maxOffsetTime time.Duration) ConfigOption {
	return func(c *config) {
		c.maxOffsetTime = maxOffsetTime
	}
}

// WithReties be set to priority about  the gradient decreases for the interval of requesting lock,After how many repetitions
func WithReties(reties int) ConfigOption {
	return func(c *config) {
		c.reties = reties
	}
}

// WithStorageClient the must be init
func WithStorageClient(client *redis.Client) ConfigOption {
	return func(c *config) {
		c.delegate = client
	}
}

// AssemblyMutex the mutex config init
func AssemblyMutex(options ...ConfigOption) {
	once.Do(func() {
		for _, value := range options {
			value(singleConfig)
		}
	})
}

// must init before use
var (
	singleConfig *config
	once         sync.Once
)

func (mutex *Mutex) Lock() {
	timeOffset := mutex.config.maxOffsetTime
	retryTimes := 0
	for {
		ok := mutex.TryLock()
		if ok {
			//delay
			go func() {
				for {
					select {
					case <-mutex.delayDone:
						return
					default:
						err := mutex.delay()
						if err != nil {
							mutex.ending <- err
							return
						}
						time.Sleep(mutex.config.expiresTime / 4)
					}
				}
			}()
			return
		}
		time.Sleep(timeOffset)
		retryTimes++
		if mutex.config.reties <= retryTimes {
			timeOffset /= 2
			retryTimes = 0
		}
	}
}

func (mutex *Mutex) TryLock() bool {
	ok, err := mutex.acquire()
	if err != nil || !ok {
		return false
	}
	return true
}

func (mutex *Mutex) acquire() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mutex.config.cancelTime)
	defer cancel()
	cmd := mutex.config.delegate.SetNX(ctx, mutex.name, mutex.config.nodeID, mutex.config.expiresTime)
	return cmd.Val(), cmd.Err()
}

// -1 resprent the delay error
const delayScript = `
   if (redis.call('GET',KEYS[1])==ARGS[1])then
          return redis.call('TTL',ARGS[2])
	else
	      return -1
	done	  	   
`
const delayCacheKey = "delay"

func (mutex *Mutex) delay() error {
	ctx, cancel := context.WithTimeout(context.TODO(), mutex.config.cancelTime)
	defer cancel()
	if _, ok := cacheHash[delayCacheKey]; !ok {
		cmd := mutex.config.delegate.Eval(ctx, delayScript, []string{mutex.name}, mutex.config.nodeID)
		hash := sha256.Sum256([]byte(delayScript))
		cacheHash[delayCacheKey] = string(hash[:])
		return cmd.Err()
	}
	return mutex.config.delegate.EvalSha(ctx, cacheHash[delayCacheKey], []string{mutex.name}, mutex.config.nodeID, mutex.config.expiresTime).Err()
}

func (mutex *Mutex) UnLock() error {
	return mutex.release()
}

const releaseScript = `
   if(redis.call('GET',KEYS[1])==ARGS[1])then
           return redis.call('DEL',KEYS[1])
	done	   
`

var cacheHash = make(map[string]string)

const releaseCacheKey = "release"

// const deleteCacheKey = "delete"

func (mutex *Mutex) release() error {
	ctx, cancel := context.WithTimeout(context.TODO(), mutex.config.cancelTime)
	defer cancel()
	if _, ok := cacheHash[releaseCacheKey]; !ok {
		cmd := mutex.config.delegate.Eval(ctx, releaseScript, []string{mutex.name}, mutex.config.nodeID)
		hash := sha256.Sum256([]byte(releaseScript))
		cacheHash[releaseCacheKey] = string(hash[:])
		return cmd.Err()
	}
	return mutex.config.delegate.EvalSha(ctx, cacheHash[releaseCacheKey], []string{mutex.name}, mutex.config.nodeID).Err()
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
