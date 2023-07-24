package infrastructure

type PlatForm interface {
	GetExecLocation() string
	PushTask(Operate)
}

type basePlatForm struct {
	// the use Language
	Language string
	// the exec location
	Exec string
}

func PushTask(operate Operate) {
  
}

func (base *basePlatForm) GetExecLocation() string {
	return base.Exec
}

type PythonPlatForm struct {
	basePlatForm
}
