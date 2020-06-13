package io

import (
	"bufio"
	sysio "io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Newline = []byte{'\n'}

type FilePrepender struct {
	File string
	Perm os.FileMode

	temp *os.File
}

func (p *FilePrepender) Write(data []byte) (int, error) {
	return p.temp.Write(data)
}

func (p *FilePrepender) Close() (err error) {
	var swap string

	defer os.Remove(p.temp.Name())

	err = p.temp.Close()
	if err != nil {
		return
	}

	swap, err = p.closeLogic()
	if err != nil {
		return
	}

	defer os.Remove(swap)

	err = os.Chmod(swap, p.Perm)
	if err != nil {
		return
	}

	err = os.Rename(swap, p.File)
	if err != nil {
		return
	}

	return
}

func (p *FilePrepender) closeLogic() (path string, err error) {
	var prepend, append, temp *os.File
	name := filepath.Base(p.File)
	dir := filepath.Dir(p.File)

	temp, err = ioutil.TempFile(dir, name+".concat")
	if err != nil {
		return
	}

	defer func() {
		if cErr := temp.Close(); err == nil {
			err = cErr
		}
		if err != nil {
			os.Remove(temp.Name())
		}
	}()

	prepend, err = os.Open(p.temp.Name())
	if err != nil {
		return
	}

	defer func() {
		if cErr := prepend.Close(); err == nil {
			err = cErr
		}
	}()

	append, err = os.Open(p.File)
	if os.IsNotExist(err) {
		path, err = p.concat(temp, prepend)
		return
	} else if err != nil {
		return
	}

	defer func() {
		if cErr := append.Close(); err == nil {
			err = cErr
		}
	}()

	path, err = p.concat(temp, prepend, append)
	return
}

func (p *FilePrepender) concat(out *os.File, in ...sysio.Reader) (path string, err error) {
	for _, reader := range in {
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			line := scanner.Bytes()
			if _, err = WriteLine(out, line); err != nil {
				return
			}
		}

		if err = scanner.Err(); err != nil {
			return
		}
	}

	return out.Name(), out.Sync()
}

// WriteLine forwards the input date to the write and will write
// and additional newline. The number of writte bites is the sum
// of those operations (unless it exits early due to an error).
// The error (if any), is passed on from the underlying call to io.Writer.Write
func WriteLine(writer sysio.Writer, line []byte) (len int, err error) {
	var n int

	if len, err = writer.Write(line); err != nil {
		return
	}

	if n, err = writer.Write(Newline); err != nil {
		return
	}

	len += n
	return
}

func OpenFilePrepender(file string, perm os.FileMode) (prepender *FilePrepender, err error) {
	name := filepath.Base(file)
	dir := filepath.Dir(file)
	prepender = &FilePrepender{
		File: file,
		Perm: perm,
	}

	prepender.temp, err = ioutil.TempFile(dir, name+".buffer")
	if err != nil {
		return
	}

	return
}
