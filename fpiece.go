package fpiece // import "github.com/aleitner/FilePiece"

import (
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
