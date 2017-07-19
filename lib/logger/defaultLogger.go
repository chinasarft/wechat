package logger

var (
	defaultLogger *JsonLogger
)

func Start(path, localFile string) (err error) {

	defaultLogger, err = NewLogger(path, localFile)
	return
}

func SetCallDepth(depth int) {
	if defaultLogger != nil {
		defaultLogger.SetCallDepth(depth)
	}
}

func SetLogLevel(logLevel int) {
	if defaultLogger != nil {
		defaultLogger.SetLogLevel(logLevel)
	}
}

func SetRotateSize(s int64) {
	if defaultLogger != nil {
		defaultLogger.SetRotateSize(s)
	}
}
func End() {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
}
func Flush() {
	if defaultLogger != nil {
		defaultLogger.Flush()
	}
}
func Close() {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
}

func Fatal(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(FATAL, flag, msg, v...)
	}
}

func Error(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(ERROR, flag, msg, v...)
	}
}

func Warn(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(WARN, flag, msg, v...)
	}
}

func Info(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(INFO, flag, msg, v...)
	}
}

func Debug(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(DEBUG, flag, msg, v...)
	}
}
func Trace(flag, msg string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Trace(TRACE, flag, msg, v...)
	}
}
