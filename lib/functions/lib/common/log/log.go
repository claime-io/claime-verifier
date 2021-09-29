package log

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func Info(i interface{}) {
	logrus.Info(i)
}

func Error(i interface{}, err error) {
	logrus.WithFields(logrus.Fields{"error": err.Error()}).Error(i)
}
