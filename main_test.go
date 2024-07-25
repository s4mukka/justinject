package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/s4mukka/justinject/mock"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	originalBrokerInit := brokerInit
	defer func() { brokerInit = originalBrokerInit }()
	brokerInit = func(ctx context.Context) error {
		return nil
	}

	os.Args = []string{"cmd", "broker"}

	output := mock.CaptureError(main)
	assert.Equal(t, "", output)
}

func TestMainFunctionWithError(t *testing.T) {
	originalBrokerInit := brokerInit
	defer func() { brokerInit = originalBrokerInit }()
	brokerInit = func(ctx context.Context) error {
		return nil
	}

	cliFactory = &mock.MockCliFactory{}

	output := mock.CaptureError(main)

	assert.Equal(t, "Error starting CLI: simulated error\n", output)
}

func TestMainIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmd := func(ctx context.Context) error {
		return nil
	}

	rootCmd := BuildCmd("anything", mockCmd)

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Error starting CLI: %v\n", err)
	}
}

func TestMainIntegrationWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmd := func(ctx context.Context) error {
		return fmt.Errorf("simulated error") // Simulate an error
	}

	rootCmd := BuildCmd("anything", mockCmd)

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Error starting CLI: %v\n", err)
	}
}
