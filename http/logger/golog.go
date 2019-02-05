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
}

func newGoLogger() *goLogger {
	return &goLogger{
		logger:    log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		callDepth: defaultCallDepth,
	}
}

func (l *goLogger) Log(v ...interface{}) {
	l.output(fmt.Sprintln(v...))
}

func (l *goLogger) Logf(format string, args ...interface{}) {
	l.output(fmt.Sprintf(format+"\n", args...))
}

func (l *goLogger) output(msg string) {
	l.logger.Output(l.callDepth, msg)
}
