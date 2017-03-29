package logger

import (
	"abacus/Godeps/_workspace/src/github.com/op/go-logging"
	"os"
)

var log *logging.Logger


func InitLogger(){
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset}%{message}`,
	)
	log = logging.MustGetLogger("")

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func Debug(format string, args ...interface{}){
	log.Debug(format, args...)
}

func Info(format string, args ...interface{}){
	log.Info(format, args...)
}

func Notice(format string, args ...interface{}){
	log.Notice(format, args...)
}

func Warning(format string, args ...interface{}){
	log.Warning(format, args...)
}

func Error(format string, args ...interface{}){
	log.Error(format, args...)
}

func Critical(format string, args ...interface{}){
	log.Critical(format, args...)
}