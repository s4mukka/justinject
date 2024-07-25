package domain

import "github.com/spf13/cobra"

type ICli interface {
	AddCommand(cmds ...*cobra.Command)
	Execute() error
}

type ICliFactory interface {
	MakeCli() ICli
}
