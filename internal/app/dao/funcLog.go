package dao

type FuncLog interface {
}

func GetFuncLogDao() FuncLog {
	return getFuncLogDao()
}
