package dao

type IP interface {
	QueryIP(ip string) int
	InsertIP(ip string) bool
}

func GetIPDao() IP {
	return getIPDao()
}
