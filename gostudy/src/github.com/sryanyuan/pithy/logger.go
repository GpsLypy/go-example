package pithy

import (
	"log"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type Logger interface {
	Outputln(lv int, v ...interface{})
	Outputf(lv int, fmt string, v ...interface{})
	Flush()
}

var globalLogger Logger

func RegisterLogger(logger Logger) {
	globalLogger = logger
}

func Debugln(v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputln(LogLevelDebug, v...)
		return
	}
	log.Println(v...)
}

func Debugf(fmt string, v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputf(LogLevelDebug, fmt, v...)
		return
	}
	log.Printf(fmt, v...)
}

func Infoln(v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputln(LogLevelInfo, v...)
		return
	}
	log.Println(v...)
}

func Infof(fmt string, v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputf(LogLevelInfo, fmt, v...)
		return
	}
	log.Printf(fmt, v...)
}

func Warnln(v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputln(LogLevelWarn, v...)
		return
	}
	log.Println(v...)
}

func Warnf(fmt string, v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputf(LogLevelWarn, fmt, v...)
		return
	}
	log.Printf(fmt, v...)
}

func Errorln(v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputln(LogLevelError, v...)
		return
	}
	log.Println(v...)
}

func Errorf(fmt string, v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputf(LogLevelError, fmt, v...)
		return
	}
	log.Printf(fmt, v...)
}

func Fatalln(v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputln(LogLevelFatal, v...)
		return
	}
	log.Println(v...)
}

func Fatalf(fmt string, v ...interface{}) {
	if nil != globalLogger {
		globalLogger.Outputf(LogLevelFatal,fmt,  v...)
		return
	}
	log.Printf(fmt, v...)
}

func Flush() {
	if nil == globalLogger {
		return
	}
	globalLogger.Flush()
}
