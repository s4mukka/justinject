package domain

import log "github.com/sirupsen/logrus"

type Environment struct {
	Instance string

	Logger *log.Entry

	Teste string
}
