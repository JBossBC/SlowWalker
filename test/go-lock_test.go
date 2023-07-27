package test

import (
	"fmt"
	"replite_web/internal/app/utils"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestLock(t *testing.T) {
	utils.AssemblyMutex(utils.WithStorageClient(redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 8})))
	mutex := utils.NewMutex("lock")
	// blockSingal := make(chan struct{})
	// go func() {
	mutex.Lock()
	fmt.Println("lock success")
	// relay
	time.Sleep(10 * time.Second)
	mutex.UnLock()
	fmt.Println("unlock sueccess")
	// err := mutex.UnLock()
	// if err != nil {
	// 	panic(err)
	// }
	// }()
	// <-blockSingal
}

// func TestUnLock(t *testing.T) {
// 	utils.AssemblyMutex(utils.WithStorageClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
// 	mutex := utils.NewMutex("unlock")
// 	blockSingal := make(chan struct{})
// 	// go func() {
// 	mutex.UnLock()
// 	<-blockSingal
// }

// func TestDistributedEnvDelay(t *testing.T) {
// 	utils.AssemblyMutex(utils.WithStorageClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
// 	wg := sync.WaitGroup{}
// 	wg.Add(10)
// 	ids := getUniqueIDS(10)
// 	for i := 0; i < 10; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			mutex := utils.NewMutex("distributed")
// 			utils.ChangeValue(ids[i])
// 			mutex.Lock()
// 			fmt.Println(ids[i], "Lock")
// 			time.Sleep(5 * time.Second)
// 			mutex.UnLock()
// 			fmt.Println(ids[i], "UnLock")
// 		}(i)
// 	}
// 	wg.Wait()
// }

// imitate the distribute environment
// func TestPerformanceForLock(t *testing.T) {
// 	utils.AssemblyMutex(utils.WithStorageClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
// 	wg := sync.WaitGroup{}
// 	wg.Add(1000)
// 	ids := getUniqueIDS(1000)
// 	for i := 0; i < 1000; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			mutex := utils.NewMutex("PDistributed")
// 			utils.ChangeValue(ids[i])
// 			mutex.Lock()
// 			fmt.Println(ids[i], "Lock")
// 			time.Sleep(2 * time.Millisecond)
// 			mutex.UnLock()
// 			fmt.Println(ids[i], "UnLock")
// 		}(i)
// 	}
// 	wg.Wait()
// }

func getUniqueIDS(number int) []string {
	var result = make([]string, number)
	for i := 0; i < number; i++ {
		result[i] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		time.Sleep(1 * time.Millisecond)
	}
	return result
}
