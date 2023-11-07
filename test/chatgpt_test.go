package test

import (
	"log"
	"replite_web/internal/app/infrastructure"
	"testing"
)

func TestChatgpt(t *testing.T) {
	result, err := infrastructure.GetAIModel().Guess("ip查询", "可以输入多个ip地址,每个IP地址占据一行")
	if err != nil {
		log.Panicln(err)
	}
	log.Println(result)
}
