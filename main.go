package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/s4mukka/justinject/broker"
	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/cli"
	internalCtx "github.com/s4mukka/justinject/internal/context"
)

var (
	brokerInit                    = broker.Init
	cliFactory domain.ICliFactory = &cli.CliFactory{}
)

func main() {
	rootCmd := cliFactory.MakeCli()
	rootCmd.AddCommand(BuildCmd("broker", brokerInit))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting CLI: %v\n", err)
	}
}

type Cmd func(ctx context.Context) error

func BuildCmd(command string, fn Cmd) *cobra.Command {
	return &cobra.Command{
		Use:   command,
		Short: fmt.Sprintf("Starting %s", command),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := internalCtx.InitializeContext(command)
			defer internalCtx.ShutdownComponents(ctx)
			environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)
			logger := environment.Logger
			logger.Infof("Starting %s...", command)
			if err := fn(ctx); err != nil {
				logger.Errorf("Error starting %s: %v\n", command, err)
			}
		},
	}
}
