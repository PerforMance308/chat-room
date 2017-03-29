package logger

import (
	"os"
	"github.com/op/go-logging"
)

var logger *logging.Logger


func InitLogger(){
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset}%{message}`,
	)
	logger = logging.MustGetLogger("")

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func Logger() *logging.Logger {
	return logger
}