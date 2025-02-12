package contracts

type ILoggerProvider interface {
	Error(message string, err error)
	Panic(message string, err error)
	ErrorMsg(message string)
	Info(message string)
	Warning(message string)
	Debug(message string, data any)
}
