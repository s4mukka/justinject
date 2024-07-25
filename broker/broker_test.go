package broker

import (
	"context"
	"testing"

	"github.com/s4mukka/justinject/broker/mocks"
	"github.com/s4mukka/justinject/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	environment := &domain.Environment{
		Logger: logrus.NewEntry(logrus.New()),
	}
	ctx = context.WithValue(ctx, "environment", environment)

	serverFactory = &mocks.MockServerFactory{}

	err := Init(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}
