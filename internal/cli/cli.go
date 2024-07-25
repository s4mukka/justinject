package cli

import (
	"github.com/s4mukka/justinject/domain"
	"github.com/spf13/cobra"
)

type Cli struct {
	handler *cobra.Command
}

func (c *Cli) AddCommand(cmds ...*cobra.Command) {
	c.handler.AddCommand(cmds...)
}

func (c *Cli) Execute() error {
	return c.handler.Execute()
}

type CliFactory struct {
}

func (c *CliFactory) MakeCli() domain.ICli {
	return &Cli{
		handler: &cobra.Command{},
	}
}
