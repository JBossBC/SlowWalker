package test

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"testing"
)

func TestLog(t *testing.T) {

	l := &dao.Log{
		Level:    dao.PRINT,
		IP:       "127.0.0.1",
		Operator: "audit",
	}
	logs, err := dao.FilterLogs(l, 2, 5)
	if err != nil {
		log.Println(err)
		return
	}
	for _, log := range *logs {
		fmt.Println(log)

	}

}
