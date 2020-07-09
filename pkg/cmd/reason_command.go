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
	Annotations  []string
	Reason       change.Reason
	Author       string
	Component    string
	MergeRequest int

	writer *change.Writer

	CommandBase
}

func (c *ReasonCommand) RunE(_ *cobra.Command, args []string) error {
	entry := change.NewEntry()

	entry.Title = strings.Join(args, " ")
	entry.Reason = c.Reason
	entry.Author = c.Author
	entry.Component = c.Component
	entry.MergeRequest = c.MergeRequest
	entry.Annotations = change.ParseAnnotations(c.Annotations)

	if path, err := c.writer.Write(entry); err != nil {
		return err
	} else {
		fmt.Println("Change entry written to", path)
	}

	return nil
}

func NewReasonCommand(reason change.Reason, output io.EntityWriter) *ReasonCommand {
	writer := change.NewWriter(output)
	author:= change.DefaultAuthor()
	command := &ReasonCommand{
		Reason: reason,
		Author: author,
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
	if value := os.Getenv(EnvGroupAuthor); len(value) != 0 {
		store.Author = value
	}

	cmd.Flags().StringArrayVarP(&store.Annotations, "annotation", "a", store.Annotations, "Arbitrary annotations to attach to the entry")
	cmd.Flags().IntVarP(&store.MergeRequest, "merge-request", "r", store.MergeRequest, "Merge request this change is related to")
	cmd.Flags().StringVarP(&store.Component, "component", "c", store.Component, "Application component this change is related to")
	cmd.Flags().StringVarP(&store.Author, "author", "A", store.Author, "Name to associate with this change")
	cmd.Flags().StringVarP(&writer.UnreleasedDir, "directory", "d", writer.UnreleasedDir, "Directory to write the changelog to")
	_ = cmd.MarkFlagDirname("directory")
}
