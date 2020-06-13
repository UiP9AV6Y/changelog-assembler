package io

import (
	sysio "io"
	"os"
	"strings"
)

const (
	ProtoStdout   = "-"
	ProtoTruncate = "file://"
	ProtoPrepend  = "prepend://"
)

type IOFactory interface {
	Reader(string) (sysio.ReadCloser, error)
	Writer(string) (sysio.WriteCloser, error)
}

type OsIOFactory struct {
	FileMode os.FileMode
}

func (f *OsIOFactory) Reader(source string) (sysio.ReadCloser, error) {
	truncateSource := strings.TrimPrefix(source, ProtoTruncate)

	if len(source) == 0 || source == ProtoStdout {
		return os.Stdin, nil
	}

	return os.Open(truncateSource)
}

func (f *OsIOFactory) Writer(target string) (sysio.WriteCloser, error) {
	prependTarget := strings.TrimPrefix(target, ProtoPrepend)
	truncateTarget := strings.TrimPrefix(target, ProtoTruncate)

	switch {
	case len(target) == 0, target == ProtoStdout:
		return os.Stdout, nil
	case target != prependTarget:
		return OpenFilePrepender(prependTarget, f.FileMode)
	}

	return os.OpenFile(truncateTarget, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileMode)
}

func NewOsIOFactory(fileMode os.FileMode) *OsIOFactory {
	ioFactory := &OsIOFactory{
		FileMode: fileMode,
	}

	return ioFactory
}
