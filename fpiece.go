package fpiece // import "github.com/aleitner/FilePiece"

import (
	"errors"
	"io"
	"os"
)

// Section of data to be concurrently read
type Chunk struct {
	File       *os.File
	Start      int64
	Final      int64
	CurrentPos int64
}

const (
	SeekStart   = 0 // seek relative to the origin of the file
	SeekCurrent = 1 // seek relative to the current offset
	SeekEnd     = 2 // seek relative to the end
)

// Create Chunk
func NewChunk(file *os.File, offset int64, length int64) *Chunk {
	return &Chunk{file, offset, length + offset, offset}
}

func (f *Chunk) Size() int64 {
	return f.Final - f.Start
}

// Concurrently read from Chunk
func (f *Chunk) Read(b []byte) (n int, err error) {
	if f.CurrentPos >= f.Final {
		return 0, io.EOF
	}

	var readLen int64 = 0
	if f.Final-f.CurrentPos > int64(len(b)) {
		readLen = int64(len(b))
	} else {
		readLen = f.Final - f.CurrentPos
	}

	n, err = f.File.ReadAt(b[:readLen], f.CurrentPos)
	f.CurrentPos += int64(n)
	return n, err
}

func (f *Chunk) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= f.Final-f.Start {
		return 0, io.EOF
	}

	off += f.Start

	if max := f.Final - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = f.File.ReadAt(p, off)
		if err == nil {
			err = io.EOF
		}
		return n, err
	}
	return f.File.ReadAt(p, off)
}

// Concurrently write to Chunk
func (f *Chunk) Write(b []byte) (n int, err error) {
	if f.CurrentPos >= f.Final {
		return 0, io.EOF
	}

	var writeLen int64 = 0
	if f.Final-f.CurrentPos > int64(len(b)) {
		writeLen = int64(len(b))
	} else {
		writeLen = f.Final - f.CurrentPos
	}

	n, err = f.File.WriteAt(b[:writeLen], f.CurrentPos)
	f.CurrentPos += int64(n)
	return n, err
}

func (f *Chunk) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= f.Final-f.Start {
		return 0, io.EOF
	}

	off += f.Start

	if max := f.Final - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = f.File.WriteAt(p, off)
		if err == nil {
			err = io.EOF
		}
		return n, err
	}
	return f.File.WriteAt(p, off)
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (f *Chunk) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case SeekStart:
		offset += f.Start
	case SeekCurrent:
		offset += f.CurrentPos
	case SeekEnd:
		offset += f.Final
	}

	// Do not seek to where somewhere outside the chunk
	if offset < f.Start || offset > f.Final {
		return 0, errOffset
	}

	f.CurrentPos = offset
	return offset - f.Start, nil
}
