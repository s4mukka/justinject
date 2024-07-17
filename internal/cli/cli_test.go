package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddCommand(t *testing.T) {
	cli := &Cli{
		handler: &cobra.Command{},
	}

	cmd1 := &cobra.Command{
		Use:   "cmd1",
		Short: "This is command 1",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cmd2 := &cobra.Command{
		Use:   "cmd2",
		Short: "This is command 2",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cli.AddCommand(cmd1, cmd2)

	assert.Contains(t, cli.handler.Commands(), cmd1)
	assert.Contains(t, cli.handler.Commands(), cmd2)
}

func TestExecute(t *testing.T) {
	executed := 0
	cli := &Cli{
		handler: &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				executed++
			},
		},
	}

	err := cli.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 1, executed)
}

func TestCliFactory(t *testing.T) {
	factory := &CliFactory{}
	cli := factory.MakeCli()

	assert.NotNil(t, cli)
	assert.IsType(t, &Cli{}, cli)
}
