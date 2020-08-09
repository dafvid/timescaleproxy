package log

import (
	l "log"
)

var Off int = 0
var InfoLevel int = 1
var DebugLevel int = 2

var Loglevel int = InfoLevel

func Print(v ...interface{}) {
	Info(v...)
}

func Info(v ...interface{}) {
	if Loglevel >= InfoLevel {
		l.Print("INFO: ", v)
	}
}

func Debug(v ...interface{}) {
	if Loglevel >= DebugLevel {
		l.Print("DEBUG: ", v)
	}
}
