package cache

import (
	"sync"
	"time"
)

type CachePool struct {
	rwMutex         sync.RWMutex
	ExpireDevice    map[string]time.Duration
	valueContainers sync.Map
}

var (
	singlePool     *CachePool
	singlePoolOnce sync.Once
)

func newCachePool() *CachePool {
	cachePool := new(CachePool)
	cachePool.rwMutex = sync.RWMutex{}
	cachePool.valueContainers = sync.Map{}
	cachePool.ExpireDevice = make(map[string]time.Duration)
	return cachePool
}

func GetCachePool() *CachePool {
	singlePoolOnce.Do(func() {
		singlePool = newCachePool()
		singlePool.circleCleanTask()
	})
	return singlePool
}
func (cachePool *CachePool) Get(key string) any {
	value, _ := cachePool.valueContainers.Load(key)
	return value

}
func (cachePool *CachePool) TryGet(key string) (any, bool) {
	return cachePool.valueContainers.Load(key)
}
func (cachePool *CachePool) Delete(key string) {
	cachePool.valueContainers.Delete(key)
}

func (cachePool *CachePool) Store(key string, value any, Duration time.Duration) {
	cachePool.valueContainers.Store(key, value)
	cachePool.rwMutex.Lock()
	defer cachePool.rwMutex.Unlock()
	cachePool.ExpireDevice[key] = Duration
}

const DEFAULT_CLEARN_CACHE_TIME = 15 * time.Second

func (cachePool *CachePool) circleCleanTask() {
	timer := time.NewTimer(DEFAULT_CLEARN_CACHE_TIME)
	for {
		select {
		case <-timer.C:
			if len(cachePool.ExpireDevice) <= 0 {
				continue
			}
			cachePool.rwMutex.Lock()
			deleteKeys := make([]string, 0, 10)
			for key, value := range cachePool.ExpireDevice {
				newExpires := value - DEFAULT_CLEARN_CACHE_TIME
				if newExpires <= 0 {
					deleteKeys = append(deleteKeys, key)
				} else {
					cachePool.ExpireDevice[key] = newExpires
				}
			}
			for i := 0; i < len(deleteKeys); i++ {
				cachePool.valueContainers.Delete(deleteKeys[i])
			}
			cachePool.rwMutex.Unlock()
		default:
			time.Sleep(DEFAULT_CLEARN_CACHE_TIME / 3)
		}
	}
}
