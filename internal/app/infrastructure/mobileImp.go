package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type MobileProvider struct {
}

var (
	mobileProvider *MobileProvider
	mobileOnce     sync.Once
)

func getMobileProvider() *MobileProvider {
	mobileOnce.Do(func() {
		mobileProvider = new(MobileProvider)
	})
	return mobileProvider
}

func init() {
	client, err := createClient(&config.GetServerConfig().SMSConfig.Key, &config.GetServerConfig().SMSConfig.Secret)
	if err != nil {
		panic(fmt.Sprintf("无法连接到阿里云SMS服务(key:%s,secret:%s):%s", config.GetServerConfig().SMSConfig.Key, config.GetServerConfig().SMSConfig.Secret, err.Error()))
	}
	smsClient = client
}

var smsClient *dysmsapi20170525.Client

// func NewSMSClient() *dysmsapi20170525.Client {
// 	return smsClient
// }

func createClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func (mobileProvider *MobileProvider) Send(phone string, code string) error {
	response, err := smsClient.SendSms(&dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("repliteWeb"),
		PhoneNumbers:  tea.String(phone),
		TemplateCode:  tea.String(config.GetServerConfig().SMSConfig.TemplateCode),
		TemplateParam: tea.String(utils.NewSMSRequest(code).MarshalJSON()),
	})
	if err != nil || *response.Body.Code != "OK" {
		if response.Body != nil {
			log.Printf("发送短信(%s)失败:%v", response.Body.String(), err)
		}
		return errors.New("发送短信失败")
	}
	return nil

}
