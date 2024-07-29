package broker

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/s4mukka/justinject/broker/mocks"
	"github.com/s4mukka/justinject/domain"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	environment := &domain.Environment{
		Logger: logrus.NewEntry(logrus.New()),
	}
	ctx = context.WithValue(ctx, domain.EnvironmentKey, environment)

	serverFactory = &mocks.MockServerFactory{}

	err := Init(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}
