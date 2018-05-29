package fpiece // import "github.com/aleitner/FilePiece"

import (
	"errors"
	"io"
	"os"
)

// Section of data to be concurrently read
type Chunk struct {
	file       *os.File
	start      int64
	final      int64
	currentPos int64
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
	return f.final - f.start
}

func (f *Chunk) Close() error {
	return f.file.Close()
}

// Concurrently read from Chunk
func (f *Chunk) Read(b []byte) (n int, err error) {
	if f.currentPos >= f.final {
		return 0, io.EOF
	}

	var readLen int64 = 0
	if f.final-f.currentPos > int64(len(b)) {
		readLen = int64(len(b))
	} else {
		readLen = f.final - f.currentPos
	}

	n, err = f.file.ReadAt(b[:readLen], f.currentPos)
	f.currentPos += int64(n)
	return n, err
}

func (f *Chunk) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= f.final-f.start {
		return 0, io.EOF
	}

	off += f.start

	if max := f.final - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = f.file.ReadAt(p, off)
		if err == nil {
			err = io.EOF
		}
		return n, err
	}
	return f.file.ReadAt(p, off)
}

// Concurrently write to Chunk
func (f *Chunk) Write(b []byte) (n int, err error) {
	if f.currentPos >= f.final {
		return 0, io.EOF
	}

	var writeLen int64 = 0
	if f.final-f.currentPos > int64(len(b)) {
		writeLen = int64(len(b))
	} else {
		writeLen = f.final - f.currentPos
	}

	n, err = f.file.WriteAt(b[:writeLen], f.currentPos)
	f.currentPos += int64(n)
	return n, err
}

func (f *Chunk) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= f.final-f.start {
		return 0, io.EOF
	}

	off += f.start

	if max := f.final - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = f.file.WriteAt(p, off)
		if err == nil {
			err = io.EOF
		}
		return n, err
	}
	return f.file.WriteAt(p, off)
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (f *Chunk) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case SeekStart:
		offset += f.start
	case SeekCurrent:
		offset += f.currentPos
	case SeekEnd:
		offset += f.final
	}

	// Do not seek to where somewhere outside the chunk
	if offset < f.start || offset > f.final {
		return 0, errOffset
	}

	f.currentPos = offset
	return offset - f.start, nil
}
