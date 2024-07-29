package mock

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/s4mukka/justinject/domain"
)

type MockCli struct{}

func (c *MockCli) AddCommand(cmds ...*cobra.Command) {}

func (c *MockCli) Execute() error {
	return fmt.Errorf("simulated error")
}

type MockCliFactory struct{}

func (c *MockCliFactory) MakeCli() domain.ICli {
	return &MockCli{}
}
