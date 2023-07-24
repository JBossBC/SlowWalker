package dao

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFAULT_LOG_DOCUMENT = "log"

const ERROR_LOG_STORAGE = "/var/repliteLog.json"

// const DEFALT_LOG_NUMBER = 10

const DEFAULT_REDIS_LOGS_PREFIX = "logs-"

const DEFUALT_REDIS_LOGS_EXPIRE = 30 * time.Second

const DEFAULT_REDIS_LOG_PREFIX = "log-"

var Empty_Log = Log{}

var (
	error_log_info *bufio.Writer
	mu             sync.Mutex
)

func init() {
	file, err := os.Open(ERROR_LOG_STORAGE)
	if _, ok := err.(*os.PathError); ok {
		return
	}
	fileInfo, _ := file.Stat()
	var error_log = make([]Log, fileInfo.Size()/int64(unsafe.Sizeof(Log{})))
	err = json.NewDecoder(bufio.NewReader(file)).Decode(error_log)
	if err != nil {
		panic(fmt.Sprintf("log(%s) recover is error: %v", ERROR_LOG_STORAGE, err))
	}
	file, err = os.OpenFile(ERROR_LOG_STORAGE, os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		panic(fmt.Sprintf("打开日志文件出错%s:%v", ERROR_LOG_STORAGE, err))
	}
	error_log_info = bufio.NewWriter(file)
	mu = sync.Mutex{}
}

type Log struct {
	Level    LogLevel `json:"level" bson:"level"`
	IP       string   `json:"ip" bson:"ip"`
	Message  string   `json:"message" bson:"message"`
	Operator string   `json:"operator" bson:"operator"`
	//unix time
	Date int64 `json:"date" bson:"date"`
}

type LogLevel string

var (
	PRINT LogLevel = "print"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	INFO  LogLevel = "info"
	PANIC LogLevel = "panic"
)

func Error(operator string, ip string, message string) {
	Errorf(operator, ip, message)
}

func Errorf(operator string, ip string, format string, v ...any) {
	log := newLog(ERROR, operator, ip, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Info(operator string, ip string, message string) {
	Infof(operator, ip, message)
}

func Infof(operator string, ip string, format string, v ...any) {
	log := newLog(INFO, operator, ip, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Panic(operator string, ip string, message string) {
	Panicf(operator, ip, message)
}
func Panicf(operator string, ip string, format string, v ...any) {
	log := newLog(PANIC, operator, ip, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Print(operator string, ip string, message string) {
	Printf(operator, ip, message)
}

func Printf(operator string, ip string, format string, v ...any) {
	log := newLog(PRINT, operator, ip, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Warn(operator string, ip string, message string) {
	Warnf(operator, ip, message)
}

func Warnf(operator string, ip string, format string, v ...any) {
	log := newLog(WARN, operator, ip, fmt.Sprintf(format, v...))
	insertLog(&log)
}
func newLog(level LogLevel, operator string, ip string, message string) Log {
	return Log{
		Level:    level,
		Operator: operator,
		IP:       ip,
		Message:  message,
		Date:     time.Now().Unix(),
	}
}

func insertLog(l *Log) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := GetLogCollection().InsertOne(ctx, l)
	if err != nil {
		bLog, err := json.Marshal(l)
		if err != nil {
			log.Printf("序列化日志信息出错:%v \r\n", bLog)
			return
		}
		go func(bLog []byte) {
			mu.Lock()
			defer mu.Unlock()
			error_log_info.Write(bLog)
		}(bLog)
	}
}

func GetLogCollection() *mongo.Collection { //返回一个mongoDB数据库链接
	return getMongoConn().Collection(config.CollectionConfig.Get(DEFAULT_LOG_DOCUMENT).(string))
}

// func QueryLogs(page int, pageNumber int) (*[]*Log, error) {
// 	var logs []*Log = make([]*Log, 0, pageNumber)
// 	/* the cache start and end is set the list length by default*/
// 	redisKey := getLogsKey(page, pageNumber)
// 	err := GetList(redisKey, logs, 0, -1)
// 	// if err != nil && err != redis.Nil {
// 	// 	log.Printf("query logs info (page: %d, pageNumber: %d) error: %s \r\n", page, pageNumber, err.Error())
// 	// 	return nil, err
// 	// }
// 	// 如果每次返回错误的时候返回logs,err的话会导致内存空间占用增大
// 	if err == nil {
// 		return &logs, nil
// 	}
// 	if err != redis.Nil {
// 		log.Printf("查询日志的redis缓存(%s)出错:%s", redisKey, err.Error())
// 	}
// 	// if len(logs) <= 0 {
// 	// 	// invalid the cache
// 	// 	return nil, nil
// 	// }
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	result, err := getLogCollection().Find(ctx, bson.D{}, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page)-1))
// 	if err != nil {
// 		log.Printf("query the logs(page: %d,pageNumber: %d) error:%s", page, pageNumber, err.Error())
// 		if err == mongo.ErrNoDocuments {
// 			err = CreateList(redisKey, nil, 30*time.Second)
// 			if err != nil {
// 				log.Printf("create the logs invalid key error:%s", err.Error())
// 			}
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	defer result.Close(context.Background())
// 	err = result.All(context.Background(), &logs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = CreateList(redisKey, result, DEFUALT_REDIS_LOGS_EXPIRE)
// 	if err != nil {
// 		log.Printf("create the logs %v cache error: %s", result, err.Error())
// 	}
// 	return &logs, nil
// }

// cant need to redis cache
func FilterLogs(l *Log, page int, pageNumber int) (*[]*Log, int, error) {
	var logs []*Log = make([]*Log, 0, pageNumber) //定义一个日志结构体类型的切片，长度为0，容量为每一页的最大显示条数
	//最后这个函数返回的是这个切片的指针，那么就代表着，从数据库中查询到日志，每pageNumber条日志作为

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 设置超时时间
	defer cancel()

	filter := bson.M{} //使用 bson.M{} 初始化的变量 filter   // 初始化一个空的 filter，用于存储查询条件
	if l.IP != "" {
		filter["ip"] = bson.M{ // 如果 IP 不为空，则添加 IP 条件
			"$regex":   fmt.Sprintf(".*%s.*", l.IP),
			"$options": "i",
		}
	}
	if l.Operator != "" {
		filter["operator"] = bson.M{ // 如果 Operator 不为空，则添加 Operator 条件
			"$regex":   fmt.Sprintf(".*%s.*", l.Operator),
			"$options": "i",
		}
	}
	if l.Level != "" {
		filter["level"] = l.Level // 如果 Level 不为空，则添加 Level 条件
	}
	//执行mongodb数据库查询，并且将查询结果存入result中
	result, err := GetLogCollection().Find(ctx, filter, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page)-1))
	//这是一个选项设置方法链，用于设置查询结果的最大返回条数。.SetLimit(int64(pageNumber)) 表示将查询结果限制为 pageNumber 条记录。
	//这也是一个选项设置方法链，用于设置查询结果的跳过记录数。.SetSkip(int64(page)-1) 表示跳过前面 (page-1) 条记录，从第 page 条记录开始获取。
	//比如，page=2.表示查询第二页的结果。pageNumber=20，表示该页最大展示20条数据

	//获取这个log结构体对应的主人的全部日志记录
	var count = 0
	//fliter2 := bson.M{}

	result2, err := GetLogCollection().Find(ctx, filter)
	for result2.Next(ctx) {
		count++
	}

	if err != nil {
		log.Printf("query the logs(filter:%v,page: %d,pageNumber: %d) error:%s", l, page, pageNumber, err.Error())
		return nil, 0, err
	}
	err = result.All(context.Background(), &logs)

	//result 是 MongoDB 查询的结果对象，它拥有一些方法用于获取查询结果。
	//.All(context.Background(), &logs) 是一个方法调用，使用 context.Background() 作为上下文对象，将查询结果解码到 logs 变量中。

	if err != nil {
		log.Println("日志数据处理出错%s", err.Error())
	}
	return &logs, count, err
	// var result = new(Log)
	// redisKey := getLogKey(l)
	// err := Get(redisKey, result)
	// if err == nil {
	// 	// to defend the Variable escape
	// 	return *result, nil
	// }
	// if err != redis.Nil {
	// 	log.Printf("查询日志redis缓存(%s)失效:%s", redisKey, err.Error())
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	// filter := bson.M{}
	// if l.IP != "" {
	// 	filter["ip"] = bson.M{
	// 		"$regex":   fmt.Sprintf(".*%s.*", l.IP),
	// 		"$options": "i",
	// 	}
	// }
	// if l.Operator != "" {
	// 	filter["operator"] = bson.M{
	// 		"$regex":   fmt.Sprintf(".*%s.*", l.Operator),
	// 		"$options": "i",
	// 	}
	// }
	// if l.Level != "" {
	// 	filter["level"] = l.Level
	// }
	// err = getLogCollection().FindOne(ctx, filter).Decode(result)
	// if err != nil {
	// 	log.Printf("query the log %v error:%s", filter, err.Error())
	// 	if err == mongo.ErrNoDocuments {
	// 		err = Create(redisKey, nil, 30*time.Second)
	// 		if err != nil {
	// 			log.Printf("create the log invalid key error:%s", err.Error())
	// 		}
	// 		return Log{}, nil
	// 	}
	// 	return Log{}, err
	// }
	// err = Create(redisKey, result, DEFUALT_REDIS_LOGS_EXPIRE)
	// if err != nil {
	// 	log.Printf("create the logs %v cache error: %s", result, err.Error())
	// }
	// return *result, nil
}

// func queryMaxPage(pageNumber int) int {

// }

func getLogsKey(page int, pageNumber int) string {
	return utils.MergeStr(DEFAULT_REDIS_LOGS_PREFIX, strconv.FormatInt(int64(page), 10), "-", strconv.FormatInt(int64(pageNumber), 10))
}
func getLogKey(log *Log) string {
	return utils.MergeStr(DEFAULT_REDIS_LOG_PREFIX, log.Operator, "-", log.IP)
}

func (l *Log) IsEmpty() bool {
	return *l == Empty_Log
}
