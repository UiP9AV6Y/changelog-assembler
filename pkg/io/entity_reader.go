package io

import (
	sysio "io"
	"io/ioutil"
)

type EntityDeserializer interface {
	Deserialize(data []byte) (interface{}, error)
}

type EntityReader interface {
	Read(string) (interface{}, error)
}

type EntityIOReader struct {
	readers      IOFactory
	deserializer EntityDeserializer
}

func (w *EntityIOReader) Read(path string) (entity interface{}, err error) {
	var reader sysio.ReadCloser
	var data []byte

	if reader, err = w.readers.Reader(path); err != nil {
		return
	}

	defer func() {
		if cErr := reader.Close(); err == nil {
			err = cErr
		}
	}()

	if data, err = ioutil.ReadAll(reader); err != nil {
		return
	}

	entity, err = w.deserializer.Deserialize(data)
	return
}

func NewEntityIOReader(readers IOFactory, deserializer EntityDeserializer) *EntityIOReader {
	entityReader := &EntityIOReader{
		readers:      readers,
		deserializer: deserializer,
	}

	return entityReader
}
