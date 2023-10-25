package dao

type Function interface {
	GetFuncMap(function string) (fm *FuncMap)
	DeleteFuncMap(funcmap FuncMap) error
	CreateFuncMap(funcmap FuncMap) error
}

func GetFunctionDao() Function {
	return getFunctionDao()
}
