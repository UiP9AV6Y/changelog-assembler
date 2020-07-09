package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/release"
)

type ReleaseCommand struct {
	changelog *release.Changelog

	CommandBase
}

func (c *ReleaseCommand) RunE(_ *cobra.Command, args []string) error {
	return c.changelog.Write(args[0])
}

func NewReleaseCommand(input io.EntityReader, ioFactory io.IOFactory) *ReleaseCommand {
	renderer := release.NewRenderer(ioFactory)
	changelog := release.NewChangelog(input, renderer)
	command := &ReleaseCommand{
		changelog: changelog,
	}
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Assemble any unreleased changelogs",
		Long: `Concatenates changelog fragments
in preparation of a new release
of your application`,
		RunE:              command.RunE,
		Aliases:           []string{"assemble", "build"},
		Args:              cobra.ExactArgs(1),
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	decorateReleaseFlags(cmd, renderer, changelog)

	command.command = cmd

	return command
}

func decorateReleaseFlags(cmd *cobra.Command,
	renderer *release.Renderer,
	changelog *release.Changelog) {
	if value := os.Getenv(EnvUnreleasedDir); len(value) != 0 {
		changelog.UnreleasedDir = value
	}
	if value := os.Getenv(EnvOutputFile); len(value) != 0 {
		renderer.OutputFile = value
	}
	if value := os.Getenv(EnvTemplateFile); len(value) != 0 {
		renderer.TemplateFile = value
	}
	if value := os.Getenv(EnvGroupComponents); len(value) != 0 {
		renderer.GroupComponents = true
	}
	if value := os.Getenv(EnvGroupAuthor); len(value) != 0 {
		changelog.Author = value
	}

	cmd.Flags().BoolVarP(&changelog.RetainInput, "keep", "k", changelog.RetainInput, "Do not delete changelog fragments after assembly")
	cmd.Flags().StringVarP(&changelog.Author, "author", "A", changelog.Author, "Name to associate with this release")
	cmd.Flags().StringVarP(&changelog.UnreleasedDir, "directory", "d", changelog.UnreleasedDir, "Directory to read the unreleased changelogs from")
	cmd.Flags().StringVarP(&renderer.OutputFile, "file", "f", renderer.OutputFile, "File to use for assembly")
	cmd.Flags().StringVarP(&renderer.TemplateFile, "template", "t", renderer.TemplateFile, "Custom template to render changelog fragment")
	cmd.Flags().BoolVarP(&renderer.GroupComponents, "group-components", "c", renderer.GroupComponents, "Group changes by their affected components")
	cmd.Flags().StringVarP(&renderer.StartDelimiterTemplate, "start-delimiter", "S", "", "Delimiter which introduces a version section")
	cmd.Flags().StringVarP(&renderer.EndDelimiterTemplate, "end-delimiter", "E", "", "Delimiter which terminates a version section")
	_ = cmd.MarkFlagDirname("directory")
	_ = cmd.MarkFlagFilename("file")
	_ = cmd.MarkFlagFilename("template")
}
