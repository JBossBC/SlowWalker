package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

//TODO will make sure that what times we should to merge the string to finish the complete redis key
const DEFUALT_ALLOW_MAX_NUMBER = 1

func IsValidRegisterInternal(ip string) bool {
	return dao.QueryIP(getRegisterKey(ip)) <= DEFUALT_ALLOW_MAX_NUMBER
}

func RegisterSuccessHook(ip string) {
	dao.InsertIP(getRegisterKey(ip))
}

func getRegisterKey(ip string) string {
	return utils.MergeStr(Register_FAILED_TIMES_PREDIXX, ip)
}

// to defind the register more account from ip
const Register_FAILED_TIMES_PREDIXX = "register-failed-"
