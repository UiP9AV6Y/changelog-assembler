package io

import (
	sysio "io"
)

type EntitySerializer interface {
	Serialize(entity interface{}) ([]byte, error)
}

type EntityWriter interface {
	Write(interface{}, string) error
}

type EntityIOWriter struct {
	writers    IOFactory
	serializer EntitySerializer
}

func (w *EntityIOWriter) Write(entity interface{}, path string) (err error) {
	var writer sysio.WriteCloser
	var data []byte

	if writer, err = w.writers.Writer(path); err != nil {
		return
	}

	defer func() {
		if cErr := writer.Close(); err == nil {
			err = cErr
		}
	}()

	if data, err = w.serializer.Serialize(entity); err != nil {
		return
	}

	_, err = writer.Write(data)
	return
}

func NewEntityIOWriter(writers IOFactory, serializer EntitySerializer) *EntityIOWriter {
	entityWriter := &EntityIOWriter{
		writers:    writers,
		serializer: serializer,
	}

	return entityWriter
}
