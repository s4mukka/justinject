package domain

import "github.com/sirupsen/logrus"

type IHook interface {
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}
