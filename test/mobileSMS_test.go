package test

import (
	"replite_web/internal/app/infrastructure"
	"replite_web/internal/app/utils"
	"testing"
)

func TestSendMessage(t *testing.T) {
	err := infrastructure.GetMobileProvider().Send("18080705675", utils.NewRandomCode())
	if err != nil {
		t.Fatalf("测试发送验证码错误:%s", err.Error())
	}
}
