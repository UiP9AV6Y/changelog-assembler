package release

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

type Changelog struct {
	Author        string
	UnreleasedDir string
	RetainInput   bool

	input  io.EntityReader
	output *Renderer
}

func (w *Changelog) Write(version string) (err error) {
	var entries change.Entries
	var files []string

	if files, entries, err = w.parseEntries(); err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("no unreleased changes found in %q", w.UnreleasedDir)
	}

	if err = w.writeRelease(version, entries); err != nil {
		return err
	}

	if w.RetainInput {
		return nil
	}

	return w.removeFiles(files)
}

func (w *Changelog) parseEntries() ([]string, change.Entries, error) {
	pattern := filepath.Join(w.UnreleasedDir, "*.yml")
	files, err := filepath.Glob(pattern)

	if err != nil {
		return files, nil, err
	}

	entries := make(change.Entries, 0, len(files))

	for _, file := range files {
		entity, err := w.input.Read(file)

		if err != nil {
			return files, nil, err
		}

		if entry, ok := entity.(*change.Entry); ok {
			entries = append(entries, entry)
		}
	}

	return files, entries, nil
}

func (w *Changelog) removeFiles(files []string) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (w *Changelog) writeRelease(version string, entries change.Entries) error {
	data := NewRenderContext(version, entries)

	data.Author = w.Author

	return w.output.Write(data)
}

func NewChangelog(input io.EntityReader, output *Renderer) *Changelog {
	author := change.DefaultAuthor()
	writer := &Changelog{
		UnreleasedDir: change.DefaultUnreleasedDir,
		Author:        author,
		input:         input,
		output:        output,
	}

	return writer
}
