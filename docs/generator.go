// +build docs

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/cmd"
)

var Application = "changelog-assembler"
var ManSection = "1"

type DocGenerator struct {
	Markdown   bool
	Output     string
	ManSection string

	command *cobra.Command
}

func (g *DocGenerator) Command() *cobra.Command {
	return g.command
}

func (g *DocGenerator) RunE(_ *cobra.Command, args []string) (err error) {
	cmd := cmd.New(g.command.Use)
	header := &doc.GenManHeader{
		Title:   strings.ToUpper(g.command.Use),
		Section: g.ManSection,
	}

	if err = os.MkdirAll(g.Output, 0755); err != nil {
		return err
	}

	if g.Markdown {
		err = doc.GenMarkdownTree(cmd.Command(), g.Output)
	} else {
		err = doc.GenManTree(cmd.Command(), header, g.Output)
	}

	return err
}

func NewDocGenerator(application, manSection string) cmd.Command {
	output, _ := os.Getwd()
	command := &DocGenerator{
		Output:     output,
		ManSection: manSection,
	}
	cmd := &cobra.Command{
		Use:          application,
		Short:        "Generate documentation",
		Long:         "Create command documentation automatically",
		RunE:         command.RunE,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&command.Output, "output", "o", command.Output, "Output directory")
	cmd.Flags().BoolVarP(&command.Markdown, "markdown", "m", command.Markdown, "Generate markdown documentation instead of man pages")

	command.command = cmd

	return command
}

func main() {
	gen := NewDocGenerator(Application, ManSection)

	if err := gen.Command().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
