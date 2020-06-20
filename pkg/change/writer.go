package change

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

const (
	DefaultUnreleasedDir = "changelogs/unreleased"
	FilenameExtension    = ".yml"
)

type Writer struct {
	UnreleasedDir string
	Force         bool

	output io.EntityWriter
}

func (w *Writer) Write(e *Entry) (string, error) {
	return w.WriteFile(e, e.Slug()+FilenameExtension)
}

func (w *Writer) WriteFile(e *Entry, filename string) (string, error) {
	path := filepath.Join(w.UnreleasedDir, filename)

	if _, err := os.Stat(path); err == nil {
		if !w.Force {
			return path, fmt.Errorf("change entry %q already exists", path)
		}
	} else if !os.IsNotExist(err) {
		return path, err
	}

	return path, w.output.Write(e, path)
}

func NewWriter(output io.EntityWriter) *Writer {
	entryWriter := &Writer{
		UnreleasedDir: DefaultUnreleasedDir,
		output:        output,
	}

	return entryWriter
}
