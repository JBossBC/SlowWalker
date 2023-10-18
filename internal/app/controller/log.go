package controller

type Log interface {
}

func GetLogController() Log {
	return getLogController()
}
