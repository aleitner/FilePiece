package fpiece

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	// "github.com/stretchr/testify/assert"
)

var readTests = []struct {
	in     string
	offset int64
	len    int64
	out    string
}{
	{"butts", 0, 5, "butts"},
	{"butts", 0, 2, "bu"},
	{"butts", 3, 2, "ts"},
	{"butts", 0, 10, "butts"},
	{"butts", 1, 1100, "utts"},
}

var writeTests = []struct {
	in     string
	offset int64
	len    int64
	out    string
}{
	{"butts", 0, 5, "butts"},
	{"butts", 0, 2, "bu"},
	{"butts", 3, 2, "\x00\x00\x00bu"},
	{"butts", 0, 10, "butts"},
	{"butts", 1, 1100, "\x00butts"},
}

func TestRead(t *testing.T) {

	for _, tt := range readTests {
		t.Run("Reads data properly", func(t *testing.T) {

			tmpfilePtr, err := ioutil.TempFile("", "read_test")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpfilePtr.Name()) // clean up

			if _, err := tmpfilePtr.Write([]byte(tt.in)); err != nil {
				log.Fatal(err)
			}

			chunk := NewChunk(tmpfilePtr, tt.offset, tt.len)

			buffer := make([]byte, 100)
			n, _ := chunk.Read(buffer)

			if err := tmpfilePtr.Close(); err != nil {
				log.Fatal(err)
			}

			if string(buffer[:n]) != tt.out {
				t.Errorf("got %q, want %q", string(buffer[:n]), tt.out)
			}

		})
	}

}

func TestWrite(t *testing.T) {

	for _, tt := range writeTests {
		t.Run("Writes data properly", func(t *testing.T) {

			tmpfilePtr, err := ioutil.TempFile("", "write_test")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpfilePtr.Name()) // clean up

			chunk := NewChunk(tmpfilePtr, tt.offset, tt.len)
			chunk.Write([]byte(tt.in))

			buffer := make([]byte, 100)
			n, err := tmpfilePtr.Read(buffer)

			if err := tmpfilePtr.Close(); err != nil {
				log.Fatal(err)
			}

			if string(buffer[:n]) != tt.out {
				t.Errorf("got %q, want %q", string(buffer[:n]), tt.out)
			}

		})
	}

}

func TestMain(m *testing.M) {
	m.Run()
}
