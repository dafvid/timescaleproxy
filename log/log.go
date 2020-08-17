package log

import (
	l "log"

	"github.com/dafvid/timescaleproxy/util"
)

type LevelType int

var Off LevelType = 0
var ErrorLevel LevelType = 1
var InfoLevel LevelType = 2
var DebugLevel LevelType = 3

var levelString []string = []string{
	"Off",
	"Error",
	"Info",
	"Debug",
}

var LogLevel LevelType = ErrorLevel

func PrintLevel() {
	Infof("Log level is: %v", levelString[LogLevel])
}

func Print(v ...interface{}) {
	l.Print(v...)
}

func Printf(fmt string, v ...interface{}) {
	l.Printf(fmt, v...)
}

func Info(v ...interface{}) {
	if LogLevel >= InfoLevel {
		v := util.Prepend("INFO: ", v)
		l.Print(v...)
	}
}

func Infof(fmt string, v ...interface{}) {
	if LogLevel >= InfoLevel {
		fmt = "INFO: " + fmt
		l.Printf(fmt, v...)
	}
}

func Debug(v ...interface{}) {
	if LogLevel >= DebugLevel {
		v := util.Prepend("DEBUG: ", v)
		l.Print(v...)
	}
}

func Debugf(fmt string, v ...interface{}) {
	if LogLevel >= DebugLevel {
		fmt = "DEBUG: " + fmt
		l.Printf(fmt, v...)
	}
}

func Error(v ...interface{}) {
	if LogLevel >= ErrorLevel {
		v := util.Prepend("DEBUG: ", v)
		l.Print(v...)
	}
}

func Errorf(fmt string, v ...interface{}) {
	if LogLevel >= ErrorLevel {
		fmt = "ERROR: " + fmt
		l.Printf(fmt, v...)
	}
}
