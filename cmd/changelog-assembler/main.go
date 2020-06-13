package main

import (
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/cmd"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

func main() {
	app := filepath.Base(os.Args[0])
	stream := io.NewOsIOFactory(0644)
	serializer := change.NewEntrySerializer()
	input := io.NewEntityIOReader(stream, serializer)
	output := io.NewEntityIOWriter(stream, serializer)
	rootCmd := cmd.NewRootCommand(app)
	cmds := []cmd.Command{
		cmd.NewVersionCommand(app),
		cmd.NewReleaseCommand(input, stream),
		cmd.NewCreateCommand(output),
		cmd.NewParseCommand(stream),
	}

	for _, reason := range change.ReasonInstances[1:] {
		cmds = append(cmds, cmd.NewReasonCommand(reason, output))
	}

	for _, cmd := range cmds {
		rootCmd.RegisterCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
