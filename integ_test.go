package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TestWordCount int = 10000000
const TestTableSize int = 1000

type Verify struct {
	words []string
	ptr   int
}

func ReadAndVerify(
	r io.Reader,
	v *Verify,
	isMemTable bool,
	t *testing.T,
) {
	for {
		buf := make([]byte, INT_LEN)
		_, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				assert.Empty(t, err.Error())
			}
		}

		l := int64(binary.LittleEndian.Uint64(buf))

		if l == 0 {
			if isMemTable {
				break
			} else {
				assert.Empty(t, "reading 0 in word len in file")
			}
		}

		buf = make([]byte, l)
		_, err = r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				assert.Empty(t, err.Error())
			}
		}

		w := string(buf)
		assert.Equal(t, v.words[v.ptr], w)
		v.ptr += 1
	}
}
func TestFS(t *testing.T) {
	cfg := Config{
		TableSize: TestTableSize,
		FilePath:  "./test.db",
	}
	fs := NewLocalFS(&cfg)

	wordsPool := []string{
		"yo this is really good",
		"the quick brown fox jumps over the lazy dog",
		"how about this",
	}

	v := Verify{
		words: []string{},
		ptr:   0,
	}

	// Write to the FS
	println("writing to fs")
	start := time.Now()
	for i := 0; i < TestWordCount; i++ {
		w := wordsPool[rand.Intn(len(wordsPool))]

		v.words = append(v.words, w)

		b := []byte(w)
		_, err := fs.Write(&b)
		if err != nil {
			assert.Empty(t, err.Error())
		}
	}

	println("write took time " + time.Since(start).String())

	// Read and verify
	println("trying to read from disk")
	start = time.Now()
	fstat, err := os.Stat(cfg.FilePath)
	if err != nil && !errors.Is(err, os.ErrExist) {
		assert.Empty(t, err.Error())
	}

	// Verify words in file
	f, err := os.Open(cfg.FilePath)
	if err != nil {
		assert.Empty(t, err.Error())
	}
	fBuf := make([]byte, fstat.Size())
	_, err = f.Read(fBuf)
	if err != nil {
		assert.Empty(t, "error reading from file to fbuf")
	}
	println("read from disk took time " + time.Since(start).String())
	fReader := bytes.NewReader(fBuf)
	ReadAndVerify(fReader, &v, false, t)

	// Verify words in memory
	println("trying to read from mem")
	bReader := bytes.NewReader(fs.table.buf)
	ReadAndVerify(bReader, &v, true, t)
	println("read and verify took time " + time.Since(start).String())

	assert.Equal(t, len(v.words), v.ptr)
}
