package cli

import "github.com/spf13/cobra"

type ICli interface {
	AddCommand(cmds ...*cobra.Command)
	Execute() error
}

type Cli struct {
	handler *cobra.Command
}

func (c *Cli) AddCommand(cmds ...*cobra.Command) {
	c.handler.AddCommand(cmds...)
}

func (c *Cli) Execute() error {
	return c.handler.Execute()
}

type ICliFactory interface {
	MakeCli() ICli
}

type CliFactory struct {
}

func (c *CliFactory) MakeCli() ICli {
	return &Cli{
		handler: &cobra.Command{},
	}
}
