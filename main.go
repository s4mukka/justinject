package main

import (
	"github.com/s4mukka/justinject/broker"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(broker.Init())

	rootCmd.Execute()
}
