package test

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"testing"
)

func TestLog(t *testing.T) {
	// 编写 input
	l := &dao.LogInfo{
		Level:    dao.PRINT,
		IP:       "127.0.0.1",
		Operator: "audit",
	}
	// 运行Dao
	logs, err := dao.GetLogDao().FilterLogs(l, 2, 5)
	if err != nil {
		log.Println(err)
		return
	}
	// 检查output 并且执行比对
	for _, log := range *logs {
		t.Fail()
		fmt.Println(log)
	}
}

func TestLogService(t *testing.T) {
	//
}

func TestLogController(t *testing.T) {

}
