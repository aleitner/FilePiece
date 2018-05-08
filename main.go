package piece // import "github.com/aleitner/FilePiece"

import (
	"io"
	"os"
)

// Section of data to be concurrently read
type FileChunk struct {
	file       *os.File
	offset     int64
	length 		 int64
	currentPos int64
}

func (f FileChunk) Read(b []byte) (n int, err error) {
	if f.currentPos >= f.length {
		return 0, io.EOF
	}

	n, err = f.file.ReadAt(b[:f.length], f.offset + f.currentPos)
	f.currentPos += int64(n)
	return n, err
}

func (f FileChunk) Write(b []byte) (n int, err error) {
	if f.currentPos >= f.length {
		return 0, io.EOF
	}

	n, err = f.file.WriteAt(b[:f.length], f.offset + f.currentPos)
	f.currentPos += int64(n)
	return n, err
}
