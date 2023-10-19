package service

import "replite_web/internal/app/utils"

type Mobile interface {
	SendMessage(phone string, ip string) (response utils.Response)
}

func GetMobileService() Mobile {
	return getMobileService()
}
