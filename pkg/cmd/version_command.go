package cmd

import (
	"os"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/version"
)

const DefaultVersionFormat = `Version: {{ .Version }}
Commit: {{ .Commit }}
Build Date: {{ .BuildDate }}
Go Version: {{ .GoVersion }}
Compiler: {{ .Compiler }}
Platform: {{ .Platform }}
`

type VersionCommand struct {
	Format string

	CommandBase
}

func (c *VersionCommand) RunE(_ *cobra.Command, args []string) error {
	format := c.Format

	if len(format) == 0 {
		format = DefaultVersionFormat
	}

	tmpl, err := template.New("version").Parse(format)
	if err != nil {
		return err
	}

	info := version.NewInfo()
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		return err
	}

	return nil
}

func NewVersionCommand() *VersionCommand {
	command := &VersionCommand{
		// do not initialize Format with DefaultVersionFormat here
		// as it would be used in the help message of this command
		// (which would bloat the help message due to the length of
		// the default value)
	}
	cmd := &cobra.Command{
		Use:               "version",
		Short:             "Show version information",
		Long:              `Render the application version to stdout and exit.`,
		RunE:              command.RunE,
		SuggestFor:        []string{"info"},
		Args:              cobra.NoArgs,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	decorateVersionFlags(cmd, command)

	command.command = cmd

	return command
}

func decorateVersionFlags(cmd *cobra.Command,
	store *VersionCommand) {
	cmd.Flags().StringVarP(&store.Format, "format", "f", store.Format, "Render the output using the provided Go template")
}
