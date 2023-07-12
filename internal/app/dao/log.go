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
const DEFALT_LOG_NUMBER = 10

const DEFAULT_REDIS_LOGS_PREFIX = "logs-"

const DEFUALT_REDIS_LOGS_EXPIRE = 3 * time.Minute

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

func Error(operator string, message string) {
	Errorf(operator, message)
}

func Errorf(operator string, format string, v ...any) {
	log := newLog(ERROR, operator, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Info(operator string, message string) {
	Infof(operator, message)
}

func Infof(operator string, format string, v ...any) {
	log := newLog(INFO, operator, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Panic(operator string, message string) {
	Panicf(operator, message)
}
func Panicf(operator string, format string, v ...any) {
	log := newLog(PANIC, operator, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Print(operator string, message string) {
	Printf(operator, message)
}

func Printf(operator string, format string, v ...any) {
	log := newLog(PRINT, operator, fmt.Sprintf(format, v...))
	insertLog(&log)
}

func Warn(operator string, message string) {
	Warnf(operator, message)
}

func Warnf(operator string, format string, v ...any) {
	log := newLog(WARN, operator, fmt.Sprintf(format, v...))
	insertLog(&log)
}
func newLog(level LogLevel, operator string, message string) Log {
	return Log{
		Level:    level,
		Operator: operator,
		Message:  message,
		Date:     time.Now().Unix(),
	}
}

func insertLog(l *Log) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := getLogCollection().InsertOne(ctx, l)
	if err != nil {
		bLog, err := json.Marshal(l)
		if err != nil {
			log.Printf("序列化日志信息出错:%v \r\n", bLog)
			return
		}
		go func() {
			mu.Lock()
			defer mu.Unlock()
			error_log_info.Write(bLog)
		}()
	}
}

func getLogCollection() *mongo.Collection {
	return getMongoConn().Collection(config.CollectionConfig.Get(DEFAULT_LOG_DOCUMENT).(string))
}

func QueryLogs(page int, pageNumber int) ([]*Log, error) {
	var logs []*Log = make([]*Log, 0, DEFALT_LOG_NUMBER)
	/* the cache start and end is set the list length by default*/
	redisKey := getLogsKey(page, pageNumber)
	err := GetList(redisKey, logs, 0, -1)
	if err != nil {
		log.Printf("query logs info (page: %d, pageNumber: %d) error: %s \r\n", page, pageNumber, err.Error())
		return nil, err
	}
	if len(logs) <= 0 {
		// invalid the cache
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := getLogCollection().Find(ctx, bson.D{}, options.Find().SetLimit(int64(pageNumber)), options.Find().SetSkip(int64(page)-1))
	if err != nil {
		log.Printf("query the logs(page: %d,pageNumber: %d) error:%s", page, pageNumber, err.Error())
		if err == mongo.ErrNoDocuments {
			err = CreateList(redisKey, nil, 1*time.Minute)
			if err != nil {
				log.Printf("create the logs invalid key error:%s", err.Error())
			}
			return nil, nil
		}
		return nil, err
	}
	defer result.Close(context.Background())
	err = result.All(context.Background(), result)
	if err != nil {
		return nil, err
	}
	err = CreateList(redisKey, result, DEFUALT_REDIS_LOGS_EXPIRE)
	if err != nil {
		log.Printf("create the logs %v cache error: %s", result, err.Error())
	}
	return logs, nil
}

func getLogsKey(page int, pageNumber int) string {
	return utils.MergeStr(DEFAULT_REDIS_LOGS_PREFIX, strconv.FormatInt(int64(page), 10), "-", strconv.FormatInt(int64(pageNumber), 10))
}
