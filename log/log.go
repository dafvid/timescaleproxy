package log

import (
	l "log"

	"github.com/dafvid/timescaleproxy/util"
)

var Off int = 0
var ErrorLevel int = 1
var InfoLevel int = 2
var DebugLevel int = 3

var Loglevel int = ErrorLevel

func Print(v ...interface{}) {
	Info(v...)
}

func Info(v ...interface{}) {
	if Loglevel >= InfoLevel {
		v := util.Prepend("INFO: ", v)
		l.Print(v...)
	}
}

func Debug(v ...interface{}) {
	if Loglevel >= DebugLevel {
		v := util.Prepend("DEBUG: ", v)
		l.Print(v...)
	}
}

func Error(v ...interface{}) {
	if Loglevel >= DebugLevel {
		v := util.Prepend("DEBUG: ", v)
		l.Print(v...)
	}
}
