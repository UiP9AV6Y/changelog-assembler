package cmd

import (
	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

func New(application string) Command {
	stream := io.NewOsIOFactory(0644)
	serializer := change.NewEntrySerializer()
	input := io.NewEntityIOReader(stream, serializer)
	output := io.NewEntityIOWriter(stream, serializer)
	rootCmd := NewRootCommand(application)
	cmds := []Command{
		NewVersionCommand(),
		NewReleaseCommand(input, stream),
		NewCreateCommand(output),
		NewParseCommand(stream),
	}

	for _, reason := range change.ReasonInstances[1:] {
		cmds = append(cmds, NewReasonCommand(reason, output))
	}

	for _, c := range cmds {
		rootCmd.Command().AddCommand(c.Command())
	}

	return rootCmd
}
