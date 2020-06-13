package cmd

import (
	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/version"
)

type RootCommand struct {
	command *cobra.Command
}

func (c *RootCommand) Register(cmd *cobra.Command) {
	cmd.AddCommand(c.command)
}

func (c *RootCommand) RegisterCommand(cmd Command) {
	cmd.Register(c.command)
}

func (c *RootCommand) Execute() error {
	return c.command.Execute()
}

func NewRootCommand(application string) *RootCommand {
	command := &RootCommand{}
	cmd := &cobra.Command{
		Use:     application,
		Version: version.Version(),
		Short:   "Changelog generator using per-change fragments",
		Long: `Compile fragments of changelog entries
into a single file. This avoids merge conflicts
as only a single person works on the main changelog
file, while contributors work on dedicated files.`,
		Annotations: map[string]string{
			"commit": version.Commit(),
		},
	}

	cmd.SetVersionTemplate("{{.Version}}-{{.Annotations.commit}}\n")

	command.command = cmd

	return command
}
