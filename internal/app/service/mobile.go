package service

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/infrastructure"
	"replite_web/internal/app/utils"
	"strings"
	"time"
)

const DEFAULT_CODE_VALID_TIME = 1 * time.Minute

const DEFAULT_REDIS_PHONE_CODE_PREFIX = "phoneCode-"

func SendMessage(phone string) (response utils.Response) {
	code := utils.NewRandomCode()
	err := infrastructure.Send(phone, code)
	if err != nil {
		response = utils.NewFailedResponse("发送验证码失败")
		return
	}
	//发送成功
	err = dao.Create(getPhoneCodeKey(phone), utils.IgnoreQuotationMarks(code), DEFAULT_CODE_VALID_TIME)
	if err != nil {
		response = utils.NewFailedResponse("发送验证码失败")
		return
	}
	return utils.NewSuccessResponse(nil)
}
func DeleteCode(key string) {
	realKey := getPhoneCodeKey(key)
	err := dao.Del(realKey)
	if err != nil {
		log.Printf("删除验证码(%s)出错:%s", realKey, err.Error())
	}
}

func IsMatching(key string, code string) bool {
	realKey := getPhoneCodeKey(key)
	redisCode, err := dao.GetStr(realKey)
	if err != nil {
		log.Printf("获取redis的值(key:%s)出错:%s", realKey, err.Error())
		return false
	}
	log.Printf("redis code:%s,user code:%s", redisCode, code)
	if strings.Compare(code, redisCode) == 0 {
		return true
	} else {
		return false
	}
}

func getPhoneCodeKey(phoneNumber string) string {
	return utils.MergeStr(DEFAULT_REDIS_PHONE_CODE_PREFIX, phoneNumber)
}
