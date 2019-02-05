package logger

var defaultLogger = (Logger)(func() *goLogger {
	ret := newGoLogger()
	ret.callDepth = 4
	return ret
}())

type Logger interface {
	Log(v ...interface{})
	Logf(format string, args ...interface{})
}

func SetDefaultLogger(l Logger) {
	defaultLogger = l
}

func Log(v ...interface{}) {
	defaultLogger.Log(v...)
}

func Logf(format string, args ...interface{}) {
	defaultLogger.Logf(format, args...)
}
