package release

import (
	"time"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/change"
)

const (
	VersionDateFormat = "2006-01-02"
)

type RenderContext struct {
	Version string
	Notes   string
	Entries change.Entries
	date    *time.Time
}

func (c *RenderContext) Date() string {
	return c.DateF(VersionDateFormat)
}

func (c *RenderContext) DateF(format string) string {
	return c.date.Format(format)
}

func NewRenderContext(version string, entries change.Entries) *RenderContext {
	date := time.Now()
	templateData := &RenderContext{
		Version: version,
		Entries: entries,
		date:    &date,
	}

	return templateData
}
