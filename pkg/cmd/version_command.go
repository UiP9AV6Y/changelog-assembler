package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/version"
)

type VersionCommand struct {
	VersionShort  bool
	VersionCommit bool

	application string
	command     *cobra.Command
}

func (c *VersionCommand) Register(cmd *cobra.Command) {
	cmd.AddCommand(c.command)
}

func (c *VersionCommand) Run(cmd *cobra.Command, args []string) {
	if c.VersionShort && c.VersionCommit {
		fmt.Println(version.Version() + "-" + version.Commit())
	} else if c.VersionShort {
		fmt.Println(version.Version())
	} else if c.VersionCommit {
		fmt.Println(version.Commit())
	} else {
		fmt.Println(version.Application(c.application))
	}
}

func NewVersionCommand(application string) *VersionCommand {
	command := &VersionCommand{
		application: application,
	}
	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Print the version number of " + application,
		Long:         `Emit the application version to stdout and exit.`,
		Run:          command.Run,
		SuggestFor:   []string{"info"},
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	decorateVersionFlags(cmd, command)

	command.command = cmd

	return command
}

func decorateVersionFlags(cmd *cobra.Command,
	store *VersionCommand) {
	cmd.Flags().BoolVarP(&store.VersionShort, "short", "s", store.VersionShort, "Emit only the version number")
	cmd.Flags().BoolVarP(&store.VersionCommit, "commit", "c", store.VersionCommit, "Emit the version control reference identifier")
}
