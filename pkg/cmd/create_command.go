package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

type CreateCommand struct {
	Annotations []string

	writer *change.Writer

	CommandBase
}

func (c *CreateCommand) RunE(_ *cobra.Command, args []string) error {
	entry := change.NewEntry()
	prompt := change.NewEntryPrompt(change.TargetUiApi())

	entry.Title = strings.Join(args, " ")
	entry.Annotations = change.ParseAnnotations(c.Annotations)

	if ok, err := prompt.Run(entry); err != nil {
		return err
	} else if !ok {
		fmt.Println("Operation aborted; not writing anything")
		return nil
	}

	if path, err := c.writer.Write(entry); err != nil {
		return err
	} else {
		fmt.Println("Change entry written to", path)
	}

	return nil
}

func NewCreateCommand(output io.EntityWriter) *CreateCommand {
	writer := change.NewWriter(output)
	command := &CreateCommand{
		writer: writer,
	}
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Interactively create a new changelog entry",
		Long:              "Interactive prompt for creating a new changelog entry",
		RunE:              command.RunE,
		Aliases:           []string{"new"},
		Args:              cobra.ArbitraryArgs,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	decorateCreateFlags(cmd, writer, command)

	command.command = cmd

	return command
}

func decorateCreateFlags(cmd *cobra.Command,
	writer *change.Writer,
	store *CreateCommand) {
	if value := os.Getenv(EnvUnreleasedDir); len(value) != 0 {
		writer.UnreleasedDir = value
	}

	cmd.Flags().StringArrayVarP(&store.Annotations, "annotation", "a", store.Annotations, "Arbitrary annotations to attach to the entry")
	cmd.Flags().StringVarP(&writer.UnreleasedDir, "directory", "d", writer.UnreleasedDir, "Directory to write the changelog to")
	_ = cmd.MarkFlagDirname("directory")
}
