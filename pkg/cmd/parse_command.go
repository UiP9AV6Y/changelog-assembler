package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/release"
)

type ParseCommand struct {
	parser *release.Parser

	CommandBase
}

func (c *ParseCommand) RunE(_ *cobra.Command, args []string) error {
	return c.parser.Read(args[0])
}

func NewParseCommand(ios io.IOFactory) *ParseCommand {
	parser := release.NewParser(ios)
	command := &ParseCommand{
		parser: parser,
	}
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the changes for a specific version",
		Long: `Parses a previously generated changelog
to extract the changes of the given version`,
		RunE:              command.RunE,
		Aliases:           []string{"extract", "cut"},
		SuggestFor:        []string{"grep", "find", "show"},
		Args:              cobra.ExactArgs(1),
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	decorateParseFlags(cmd, parser)

	command.command = cmd

	return command
}

func decorateParseFlags(cmd *cobra.Command,
	parser *release.Parser) {
	if value := os.Getenv(EnvOutputFile); len(value) != 0 {
		parser.InputFile = value
	}

	cmd.Flags().StringVarP(&parser.StartDelimiterTemplate, "start-delimiter", "S", "", "Delimiter which introduces a version section")
	cmd.Flags().StringVarP(&parser.EndDelimiterTemplate, "end-delimiter", "E", "", "Delimiter which terminates a version section")
	cmd.Flags().StringVarP(&parser.InputFile, "file", "f", parser.InputFile, "File to process")
	_ = cmd.MarkFlagFilename("file")
}
