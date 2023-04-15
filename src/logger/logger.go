package logger

// LogLevel type
type LogLevel int32

const (
	//TRACE logs everything
	TRACE LogLevel = 1 << iota
	//DEBUG logs Debug, Info, Warnings and Errors
	DEBUG
	//INFO logs Info, Warnings and Errors
	INFO
	//WARN logs Warning and Errors
	WARN
	//ERROR logs just Errors
	ERROR
)
