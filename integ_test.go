package snail

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vodkaslime/snail/utils"
)

const TestWordCount int = 1000000
const TestTableSize int = 1000
const TestFilePath string = "./test.db"

type Verify struct {
	words []string
	ptr   int
}

func readAndVerify(
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
				assert.Nil(t, err.Error())
			}
		}

		l := int64(binary.LittleEndian.Uint64(buf))

		if l == 0 {
			if isMemTable {
				break
			} else {
				assert.Nil(t, "reading 0 in word len in file")
			}
		}

		buf = make([]byte, l)
		_, err = r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				assert.Nil(t, err.Error())
			}
		}

		w := string(buf)
		assert.Equal(t, v.words[v.ptr], w)
		v.ptr += 1
	}
}

func TestFS(t *testing.T) {

	err := utils.ClearFile(TestFilePath)
	assert.Nil(t, err)

	defer func() {
		err = utils.ClearFile(TestFilePath)
		assert.Nil(t, err)
	}()

	cfg := Config{
		TableSize: TestTableSize,
		FilePath:  TestFilePath,
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
	for i := 0; i < TestWordCount; i++ {
		w := wordsPool[rand.Intn(len(wordsPool))]

		v.words = append(v.words, w)

		b := []byte(w)
		_, err := fs.Write(&b)
		if err != nil {
			assert.Nil(t, err.Error())
		}
	}

	// Read and verify
	fstat, err := os.Stat(TestFilePath)
	if err != nil && !errors.Is(err, os.ErrExist) {
		assert.Nil(t, err.Error())
	}

	// Verify words in file
	f, err := os.Open(TestFilePath)
	if err != nil {
		assert.Nil(t, err.Error())
	}
	fBuf := make([]byte, fstat.Size())
	_, err = f.Read(fBuf)
	if err != nil {
		assert.Nil(t, "error reading from file to fbuf")
	}
	fReader := bytes.NewReader(fBuf)
	readAndVerify(fReader, &v, false, t)

	// Verify words in memory
	bReader := bytes.NewReader(fs.table.buf)
	readAndVerify(bReader, &v, true, t)

	assert.Equal(t, len(v.words), v.ptr)
}
