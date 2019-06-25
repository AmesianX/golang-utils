package log

var globalLogger *Logger

type LogFunc func(string, ...Field)

var (
	Debug LogFunc
	Info  LogFunc
	Warn  LogFunc
	Error LogFunc
	Panic LogFunc
)

// InitGlobalLogger globalLogger 초기화
func InitGlobalLogger(logger *Logger) {
	globalLogger = logger

	Debug = globalLogger.Debug
	Info = globalLogger.Info
	Warn = globalLogger.Warn
	Error = globalLogger.Error
	Panic = globalLogger.Panic
}
