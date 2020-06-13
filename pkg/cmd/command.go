package cmd

import (
	"github.com/spf13/cobra"
)

type Command interface {
	Register(cmd *cobra.Command)
}
