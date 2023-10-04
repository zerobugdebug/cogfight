package logging

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log  *logrus.Logger
	once sync.Once
)

func getLogger() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			ForceColors:     true,
			FullTimestamp:   true,
		})
	})
	return log
}

func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	getLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	getLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	getLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	getLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	getLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	getLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	getLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	getLogger().Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	getLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	getLogger().Panicf(format, args...)
}
