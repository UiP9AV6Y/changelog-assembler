// +build completions

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/cmd"
)

var Application = "changelog-assembler"

type CompletionGenerator struct {
	Output string

	command *cobra.Command
}

func (g *CompletionGenerator) Command() *cobra.Command {
	return g.command
}

func (g *CompletionGenerator) RunE(_ *cobra.Command, args []string) error {
	cmd := cmd.New(g.command.Use)

	for _, arg := range args {
		shell := strings.TrimPrefix(arg, ".")

		if err := g.generate(shell, cmd.Command()); err != nil {
			return err
		}
	}

	return nil
}

func (g *CompletionGenerator) generate(shell string, command *cobra.Command) error {
	filename := filepath.Join(g.Output, command.Use)

	switch shell {
	case "bash":
		return command.GenBashCompletionFile(filename + ".bash")
	case "fish":
		return command.GenFishCompletionFile(filename+".fish", true)
	case "zsh":
		return command.GenZshCompletionFile(filename + ".zsh")
	case "powershell", "ps":
		return command.GenPowerShellCompletionFile(filename + ".ps")
	}

	return fmt.Errorf("Target shell %q is not supported", shell)
}

func NewCompletionGenerator(application string) cmd.Command {
	output, _ := os.Getwd()
	command := &CompletionGenerator{
		Output: output,
	}
	cmd := &cobra.Command{
		Use:          application,
		Short:        "Generate shell completions",
		Long:         "Create command completion scripts for interactive shells",
		RunE:         command.RunE,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&command.Output, "output", "o", command.Output, "Output directory")

	command.command = cmd

	return command
}

func main() {
	gen := NewCompletionGenerator(Application)

	if err := gen.Command().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
