package logger

import (
	"log"
	"os"
)

type goLogger struct {
	logger *log.Logger
}

func newGoLogger() *goLogger {
	return &goLogger{
		logger: log.New(os.Stdout, "golog", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
}

func (l *goLogger) Log(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *goLogger) Logf(format string, args ...interface{}) {
	l.logger.Printf(format+"\n", args)
}
