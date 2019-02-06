package logger

import (
	"fmt"
	"log"
	"os"
)

const defaultCallDepth = 3

type goLogger struct {
	logger *log.Logger
	// default depath: 3
	callDepth int
	level     LogLevel
}

func newGoLogger() *goLogger {
	return &goLogger{
		logger:    log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		callDepth: defaultCallDepth,
		level:     DEBUG,
	}
}

func (l *goLogger) Log(v ...interface{}) {
	l.output(fmt.Sprint(v...))
}

func (l *goLogger) Logf(format string, args ...interface{}) {
	l.output(fmt.Sprintf(format, args...))
}

func (l *goLogger) output(msg string) {
	l.logger.Output(l.callDepth, msg)
}

func (l *goLogger) isLogEnable(lvl LogLevel) bool {
	return l.level <= lvl
}

func (l *goLogger) SetLogLevel(lvl LogLevel) {
	l.level = lvl
}

func (l *goLogger) Debugf(format string, args ...interface{}) {
	if l.isLogEnable(DEBUG) {
		l.output(fmt.Sprintf("[DEBUG] "+format, args...))
	}
}

func (l *goLogger) Infof(format string, args ...interface{}) {
	if l.isLogEnable(INFO) {
		l.output(fmt.Sprintf("[INFO] "+format, args...))
	}
}

func (l *goLogger) Warnf(format string, args ...interface{}) {
	if l.isLogEnable(WARN) {
		l.output(fmt.Sprintf("[WARN] "+format, args...))
	}
}

func (l *goLogger) Errorf(format string, args ...interface{}) {
	if l.isLogEnable(ERROR) {
		l.output(fmt.Sprintf("[ERROR] "+format, args...))
	}
}
