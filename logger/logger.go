package logger

var defaultLogger = (Logger)(func() *goLogger {
	ret := newGoLogger()
	ret.callDepth = 4
	return ret
}())

type LogLevel int

const (
	noLevel LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, args ...interface{})

	SetLogLevel(l LogLevel)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func SetDefaultLogger(l Logger) {
	defaultLogger = l
}

func SetLogLevel(l LogLevel) {
	defaultLogger.SetLogLevel(l)
}

func Log(v ...interface{}) {
	defaultLogger.Log(v...)
}

func Logf(format string, args ...interface{}) {
	defaultLogger.Logf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}
