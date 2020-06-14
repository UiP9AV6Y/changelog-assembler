package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

type ReasonCommand struct {
	Reason       change.Reason
	Component    string
	MergeRequest int

	writer *change.Writer

	CommandBase
}

func (c *ReasonCommand) RunE(_ *cobra.Command, args []string) error {
	entry := change.NewEntry()

	entry.Title = strings.Join(args, " ")
	entry.Author = change.DefaultAuthor()
	entry.Reason = c.Reason
	entry.Component = c.Component
	entry.MergeRequest = c.MergeRequest

	if path, err := c.writer.Write(entry); err != nil {
		return err
	} else {
		fmt.Println("Change entry written to", path)
	}

	return nil
}

func NewReasonCommand(reason change.Reason, output io.EntityWriter) *ReasonCommand {
	writer := change.NewWriter(output)
	command := &ReasonCommand{
		Reason: reason,
		writer: writer,
	}
	cmd := &cobra.Command{
		Use:               reason.String(),
		Short:             fmt.Sprintf("Create a new %q changelog entry", reason),
		RunE:              command.RunE,
		Aliases:           []string{reason.Alias()},
		Args:              cobra.MinimumNArgs(1),
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	decorateReasonFlags(cmd, writer, command)

	command.command = cmd

	return command
}

func decorateReasonFlags(cmd *cobra.Command,
	writer *change.Writer,
	store *ReasonCommand) {
	if value := os.Getenv(EnvUnreleasedDir); len(value) != 0 {
		writer.UnreleasedDir = value
	}

	cmd.Flags().IntVarP(&store.MergeRequest, "merge-request", "r", store.MergeRequest, "Merge request this change is related to")
	cmd.Flags().StringVarP(&store.Component, "component", "c", store.Component, "Application component this change is related to")
	cmd.Flags().StringVarP(&writer.UnreleasedDir, "directory", "d", writer.UnreleasedDir, "Directory to write the changelog to")
	_ = cmd.MarkFlagDirname("directory")
}
