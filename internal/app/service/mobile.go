package service

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/infrastructure"
	"replite_web/internal/app/utils"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type mobileService struct {
}

var (
	mobileSvc  *mobileService
	mobileOnce sync.Once
)

func GetMobileService() *mobileService {
	mobileOnce.Do(func() {
		mobileSvc = new(mobileService)
	})
	return mobileSvc
}

const DEFAULT_CODE_VALID_TIME = 1 * time.Minute

const DEFAULT_REDIS_PHONE_CODE_PREFIX = "phoneCode-"

func (mobile *mobileService) SendMessage(phone string, ip string) (response utils.Response) {
	redisKey := getPhoneCodeKey(phone)
	// add the repeat send message failed
	_, err := dao.GetStr(redisKey)
	if err != nil {
		if err != redis.Nil {
			log.Printf("访问redis缓存失败(%s):%s", redisKey, err.Error())
			return utils.NewFailedResponse("系统错误")
		}
	} else {
		if !Incre_Warning_IP(ip) {
			return utils.NewFailedResponse("系统错误")
		}
	}
	code := utils.NewRandomCode()
	err = infrastructure.Send(phone, code)
	if err != nil {
		response = utils.NewFailedResponse("发送验证码失败")
		return
	}
	//发送成功
	err = dao.Create(redisKey, utils.IgnoreQuotationMarks(code), DEFAULT_CODE_VALID_TIME)
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
	if strings.Compare(code, redisCode) == 0 {
		return true
	} else {
		return false
	}
}

func getPhoneCodeKey(phoneNumber string) string {
	return utils.MergeStr(DEFAULT_REDIS_PHONE_CODE_PREFIX, phoneNumber)
}
