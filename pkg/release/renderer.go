package release

import (
	"fmt"
	sysio "io"
	"text/template"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

const (
	DefaultOutputFile      = "CHANGELOG.md"
	StartDelimiterTemplate = `<a name="%s"></a>`
	EndDelimiterTemplate   = `<!-- changelog-assembler: %s -->`
)

type Renderer struct {
	OutputFile             string
	TemplateFile           string
	GroupComponents        bool
	StartDelimiterTemplate string
	EndDelimiterTemplate   string

	writers io.IOFactory
}

func (r *Renderer) Write(data *RenderContext) (err error) {
	var tpl *template.Template
	var writer sysio.WriteCloser

	if writer, err = r.writers.Writer(io.ProtoPrepend + r.OutputFile); err != nil {
		return
	}

	defer func() {
		if cErr := writer.Close(); err == nil {
			err = cErr
		}
	}()

	if tpl, err = r.newTemplate(); err != nil {
		return
	}

	err = r.writeData(data, tpl, writer)
	return
}

func (r *Renderer) newTemplate() (*template.Template, error) {
	tpl := template.New("changelog_entry")

	if len(r.TemplateFile) > 0 {
		return tpl.ParseFiles(r.TemplateFile)
	} else if r.GroupComponents {
		return tpl.Parse(GroupedTemplate)
	}

	return tpl.Parse(DefaultTemplate)
}

func (r *Renderer) writeData(data *RenderContext, tpl *template.Template, writer sysio.WriteCloser) (err error) {
	startDelimiter := []byte(fmt.Sprintf(r.StartDelimiterTemplate, data.Version))
	endDelimiter := []byte(fmt.Sprintf(r.EndDelimiterTemplate, data.Version))

	if _, err := io.WriteLine(writer, startDelimiter); err != nil {
		return err
	}

	if err := tpl.Execute(writer, data); err != nil {
		return err
	}

	if _, err := io.WriteLine(writer, endDelimiter); err != nil {
		return err
	}

	return nil
}

func NewRenderer(writers io.IOFactory) *Renderer {
	renderer := &Renderer{
		StartDelimiterTemplate: StartDelimiterTemplate,
		EndDelimiterTemplate:   EndDelimiterTemplate,
		OutputFile:             DefaultOutputFile,
		writers:                writers,
	}

	return renderer
}
