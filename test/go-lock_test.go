package test

import (

	// "net/http"
	// _ "net/http/pprof"

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
	blockSingal := make(chan struct{})
	// go func() {
	// for i := 0; i < 10; i++ {
	// now := time.Now()
	mutex.Lock()
	fmt.Println("lock success")
	//..... 15s
	mutex.UnLock()
	// relay
	// time.Sleep(10 * time.Second)
	// mutex.UnLock()
	// fmt.Println(time.Now().Sub(now).Microseconds())
	// fmt.Println("unlock sueccess")
	// }
	// err := mutex.UnLock()
	// if err != nil {
	// 	panic(err)
	// }
	// }()
	<-blockSingal
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
		time.Sleep(1 * time.Microsecond)
	}
	return result
}

// func BenchmarkLock(b *testing.B) {
// 	// go func() {
// 	// 	http.ListenAndServe("localhost:6060", nil)
// 	// }()
// 	f, err := os.Create("cpu.pprof")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pprof.StartCPUProfile(f)
// 	defer pprof.StopCPUProfile()
// 	// runtime.SetBlockProfileRate(1)
// 	// ch := make(chan struct{})
// 	utils.AssemblyMutex(utils.WithStorageClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"})))
// 	wg := sync.WaitGroup{}
// 	wg.Add(1000)
// 	ids := getUniqueIDS(1000)
// 	result := make([]string, 0, 1000)
// 	for i := 0; i < 1000; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			mutex := utils.NewMutex("PDistributed")
// 			utils.ChangeValue(ids[i])
// 			mutex.Lock()
// 			result = append(result, ids[i]+"Lock")
// 			prefix := len(result) - 1
// 			// fmt.Println(ids[i], "Lock")
// 			time.Sleep(2 * time.Millisecond)
// 			mutex.UnLock()
// 			if result[prefix] != ids[i]+"Lock" {
// 				log.Println(result)
// 				panic("concurrent error")
// 			}
// 			// fmt.Println(ids[i], "UnLock")
// 		}(i)
// 	}
// 	wg.Wait()
// }
