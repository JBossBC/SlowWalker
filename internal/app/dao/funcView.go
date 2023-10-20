package dao

type FuncView interface {
	GetFuncViews(function string) ([]*FuncViewInfo, error)
	CreateFuncViews(funcs ...FuncViewInfo) error
}

func GetFuncViewDao() FuncView {
	return getFuncViewDao()
}
