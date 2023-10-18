package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"sync"
)

type IPService struct {
}

var (
	ipService *IPService
	ipOnce    sync.Once
)

func getIPService() *IPService {
	ipOnce.Do(func() {
		ipService = &IPService{}
	})
	return ipService
}

// TODO will make sure that what times we should to merge the string to finish the complete redis key
const DEFUALT_ALLOW_MAX_REGISTER_NUMBER = 1

func (IP *IPService) IsValidRegisterInternal(ip string) bool {
	return dao.GetIPDao().QueryIP(IP.getRegisterKey(ip)) <= DEFUALT_ALLOW_MAX_REGISTER_NUMBER
}

func (IP *IPService) RegisterSuccessHook(ip string) bool {
	return dao.GetIPDao().InsertIP(IP.getRegisterKey(ip))
}

func (IP *IPService) getRegisterKey(ip string) string {
	return utils.MergeStr(Register_FAILED_TIMES_PREDIXX, ip)
}

const Warn_IP_MAX = 8

// be used to defend the dangerous ip to access system by global middleware
func (IP *IPService) IsWarnIP(ip string) bool {
	return dao.GetIPDao().QueryIP(IP.getWarnIPkey(ip)) >= Warn_IP_MAX
}
func (IP *IPService) Incre_Warning_IP(ip string) bool {
	return dao.GetIPDao().InsertIP(IP.getWarnIPkey(ip))
}

func (IP *IPService) getWarnIPkey(ip string) string {
	return utils.MergeStr(Warn_IP_Times_PREFIX, ip)
}

const Warn_IP_Times_PREFIX = "dangerous-ip-"

// to defind the register more account from ip
const Register_FAILED_TIMES_PREDIXX = "register-failed-"
