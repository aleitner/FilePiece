package fpiece // import "github.com/aleitner/FilePiece"

import (
	"io"
	"os"
)

// Section of data to be concurrently read
type Chunk struct {
	File       *os.File
	Offset     int64
	Length     int64
	CurrentPos int64
}

// Create Chunk
func NewChunk(file *os.File, offset int64, length int64) Chunk {
	chunk := Chunk{}
	chunk.File = file
	chunk.Offset = offset
	chunk.Length = length
	chunk.CurrentPos = 0
	return chunk
}

// Concurrently read from Chunk
func (f Chunk) Read(b []byte) (n int, err error) {
	if f.CurrentPos >= f.Length {
		return 0, io.EOF
	}

	n, err = f.File.ReadAt(b[:f.Length], f.Offset+f.CurrentPos)
	f.CurrentPos += int64(n)
	return n, err
}

// Concurrently write to Chunk
func (f Chunk) Write(b []byte) (n int, err error) {
	if f.CurrentPos >= f.Length {
		return 0, io.EOF
	}

	n, err = f.File.WriteAt(b[:f.Length], f.Offset+f.CurrentPos)
	f.CurrentPos += int64(n)
	return n, err
}
