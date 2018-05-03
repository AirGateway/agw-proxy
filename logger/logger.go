package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

// Logger is a server logger
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Critical(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// LoggerID ia a logger ID
const LoggerID = "NDC Requester Logger"

var (
	// Logger settings
	logger           = logging.MustGetLogger(LoggerID)
	logConsoleFormat = logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05} (%{shortfile}) >> %{message} %{color:reset}`,
	)
)

func init() {
	// Prepare logger
	logConsoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logConsolePrettyBackend := logging.NewBackendFormatter(logConsoleBackend, logConsoleFormat)
	logging.SetBackend(logConsolePrettyBackend)

	// Set log level based on env
	switch os.Getenv("NDC_REQUESTER_LOGLEVEL") {
	case "debug":
		logging.SetLevel(logging.DEBUG, LoggerID)
	case "info":
		logging.SetLevel(logging.INFO, LoggerID)
	case "warn":
		logging.SetLevel(logging.WARNING, LoggerID)
	case "err":
		logging.SetLevel(logging.ERROR, LoggerID)
	default:
		logging.SetLevel(logging.DEBUG, LoggerID) // log everything by default
	}
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Fatal(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Critical(format string, args ...interface{}) {
	logger.Critical(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}

func Warning(format string, args ...interface{}) {
	logger.Warning(format, args...)
}

func Notice(format string, args ...interface{}) {
	logger.Notice(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}
