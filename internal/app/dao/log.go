package dao

type Log interface {
	Error(operator string, ip string, message string)
	Errorf(operator string, ip string, format string, v ...any)
	Info(operator string, ip string, message string)
	Infof(operator string, ip string, format string, v ...any)
	Panic(operator string, ip string, message string)
	Panicf(operator string, ip string, format string, v ...any)
	Print(operator string, ip string, message string)
	Printf(operator string, ip string, format string, v ...any)
	Warn(operator string, ip string, message string)
	Warnf(operator string, ip string, format string, v ...any)
	insertLog(l *LogInfo)
	FilterLogs(l *LogInfo, page int, pageNumber int) (*[]*LogInfo, error)
	AggregateLogSum() (int32, error)
	RemoveLogs(filters []LogInfo) error
}

func GetLogDao() Log {
	return getLogDao()
}
