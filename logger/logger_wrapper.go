package logger_wrapper

import (
	"github.com/op/go-logging"
	"os"
)

func InitLogger(module string) *logging.Logger {

	var logLevel = logging.INFO

	logger := logging.MustGetLogger(module)
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.4s} ▶ [%{shortfunc}]: %{color:reset} %{message}`,
	)

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend1 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend1Formatter := logging.NewBackendFormatter(backend1, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logLevel, module)

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled)

	return logger
}
