package mock

import (
	"fmt"

	"github.com/s4mukka/justinject/internal/cli"
	"github.com/spf13/cobra"
)

type MockCli struct{}

func (c *MockCli) AddCommand(cmds ...*cobra.Command) {}

func (c *MockCli) Execute() error {
	return fmt.Errorf("simulated error")
}

type MockCliFactory struct {
}

func (c *MockCliFactory) MakeCli() cli.ICli {
	return &MockCli{}
}
