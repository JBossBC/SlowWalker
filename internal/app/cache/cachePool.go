package cache

import "sync"

var cachePool sync.Map

func init() {
	cachePool = sync.Map{}
	
}
