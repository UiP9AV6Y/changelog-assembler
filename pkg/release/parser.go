package release

import (
	"bufio"
	"bytes"
	"fmt"
	sysio "io"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/io"
)

type Parser struct {
	InputFile              string
	OutputFile             string
	StartDelimiterTemplate string
	EndDelimiterTemplate   string

	ios io.IOFactory
}

func (p *Parser) Read(version string) (err error) {
	var reader sysio.ReadCloser
	var writer sysio.WriteCloser

	if reader, err = p.ios.Reader(p.InputFile); err != nil {
		return
	}

	defer func() {
		if cErr := reader.Close(); err == nil {
			err = cErr
		}
	}()

	if writer, err = p.ios.Writer(p.OutputFile); err != nil {
		return
	}

	defer func() {
		if cErr := writer.Close(); err == nil {
			err = cErr
		}
	}()

	return p.parseSection(version, reader, writer)
}

func (p *Parser) parseSection(version string, reader sysio.Reader, writer sysio.Writer) error {
	startDelimiter := []byte(fmt.Sprintf(p.StartDelimiterTemplate, version))
	endDelimiter := []byte(fmt.Sprintf(p.EndDelimiterTemplate, version))
	scanner := bufio.NewScanner(reader)
	found := false

	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.Equal(line, endDelimiter) {
			break
		} else if bytes.Equal(line, startDelimiter) {
			found = true
		} else if found {
			if _, err := io.WriteLine(writer, line); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Unable to find release %q in %q", version, p.InputFile)
	}

	return nil
}

func NewParser(ios io.IOFactory) *Parser {
	parser := &Parser{
		StartDelimiterTemplate: StartDelimiterTemplate,
		EndDelimiterTemplate:   EndDelimiterTemplate,
		InputFile:              DefaultOutputFile,
		OutputFile:             io.ProtoStdout,
		ios:                    ios,
	}

	return parser
}
