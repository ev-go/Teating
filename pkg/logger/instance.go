package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	I Logger
	l *logrus.Logger
)

func Init(name, environment, level string) {
	l = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter),
	}

	SetLogLevel(level)

	logger := logrus.NewEntry(l).
		WithField("name", name).
		WithField("environment", environment)

	I = Logger{logger}
}

func SetLogLevel(level string) {
	var lvl logrus.Level

	switch strings.ToLower(level) {
	case "error":
		lvl = logrus.ErrorLevel
	case "warn":
		lvl = logrus.WarnLevel
	case "info":
		lvl = logrus.InfoLevel
	case "debug":
		lvl = logrus.DebugLevel
	default:
		lvl = logrus.InfoLevel
	}

	l.SetLevel(lvl)
}
