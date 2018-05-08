# FilePiece

Concurrently read and write from files

## Installation
```BASH
go get github.com/aleitner/FilePiece
```

## Usage
```Golang
import "github.com/aleitner/FilePiece"
```

### Chunk struct
```Golang
type Chunk struct {
	File       *os.File
	Offset     int64
	Length     int64
	CurrentPos int64
}
```
* Chunk.File - os.File being read from
* Chunk.Offset - starting position for reading/writing data
* Chunk.Length - length of data to be read/written
* Chunk.CurrentPos - Keeps track to know where to write to or read from next

### NewChunk
Create a chunk from a file
```Golang
func NewChunk(file *os.File, offset int64, length int64) Chunk
```

### Read
Concurrently read from a file
```Golang
func (f Chunk) Read(b []byte) (n int, err error)
```

### Write
Concurrently write to a file
```Golang
func (f Chunk) Write(b []byte) (n int, err error)
```
