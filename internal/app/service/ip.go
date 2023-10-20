package service

type IP interface {
	IsValidRegisterInternal(ip string) bool
	RegisterSuccessHook(ip string) bool
	getRegisterKey(ip string) string
	IsWarnIP(ip string) bool
	Incre_Warning_IP(ip string) bool
	getWarnIPkey(ip string) string
}

func GetIPService() IP {
	return getIPService()
}
