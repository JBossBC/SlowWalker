package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

// TODO will make sure that what times we should to merge the string to finish the complete redis key
const DEFUALT_ALLOW_MAX_REGISTER_NUMBER = 1

func IsValidRegisterInternal(ip string) bool {
	return dao.QueryIP(getRegisterKey(ip)) <= DEFUALT_ALLOW_MAX_REGISTER_NUMBER
}

func RegisterSuccessHook(ip string) bool {
	return dao.InsertIP(getRegisterKey(ip))
}

func getRegisterKey(ip string) string {
	return utils.MergeStr(Register_FAILED_TIMES_PREDIXX, ip)
}

const Warn_IP_MAX = 8

// be used to defend the dangerous ip to access system by global middleware
func IsWarnIP(ip string) bool {
	return dao.QueryIP(getWarnIPkey(ip)) <= Warn_IP_MAX
}
func Incre_Warning_IP(ip string) bool {
	return dao.InsertIP(getWarnIPkey(ip))
}

func getWarnIPkey(ip string) string {
	return utils.MergeStr(Warn_IP_Times_PREFIX, ip)
}

const Warn_IP_Times_PREFIX = "dangerous-ip-"

// to defind the register more account from ip
const Register_FAILED_TIMES_PREDIXX = "register-failed-"
