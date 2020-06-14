package cmd

import (
	"github.com/spf13/cobra"
)

type Command interface {
	Command() *cobra.Command
}

type CommandBase struct {
	command *cobra.Command
}

func (b *CommandBase) Command() *cobra.Command {
	return b.command
}
