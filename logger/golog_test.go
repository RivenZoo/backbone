package logger

import "testing"

func TestGoLog(t *testing.T)  {
	l := newGoLogger()
	l.Debugf("this is a %s", "debug")
	l.Infof("info")

	l.SetLogLevel(INFO)
	l.Debugf("no output")
	l.Infof("still info")
}