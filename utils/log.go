package utils

import "github.com/sirupsen/logrus"

func LogInfo(s string, args ...interface{}) {
	logrus.Infof(s, args...)
}

func LogErrorMessage(s string, args ...interface{}) {
	logrus.Errorf(s, args...)
}

func LogError(err error) {
	logrus.Error(err.Error())
}
