package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry
var onceLog sync.Once

var defaultLogger = "INFO"

type utcFormatter struct {
	logrus.Formatter
}

func (u utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

func initLogger(logLevel string, logPrettyPrint bool) {

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Cannot parse Log level: %s", logLevel)
	}

	logrus.SetFormatter(utcFormatter{&logrus.JSONFormatter{
		PrettyPrint:     logPrettyPrint,
		TimestampFormat: "2006-01-02 15:04:05 -0700",
	}})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	Log = logrus.WithFields(logrus.Fields{"service": "filecoin-rpc-proxy"})
}

// InitLogger initializes logger instance
func InitLogger(logLevel string, logPrettyPrint bool) *logrus.Entry {
	onceLog.Do(func() {
		initLogger(logLevel, logPrettyPrint)
	})
	return Log
}

// InitDefaultLogger initializes logger instance
func InitDefaultLogger() *logrus.Entry {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		defaultLogger = logLevel
	}
	return InitLogger(defaultLogger, true)
}
