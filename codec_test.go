package snail

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUintConvert(t *testing.T) {
	p := []byte("yo here we go")
	s := encodePayload(&p)

	toVerify := make([]byte, INT_LEN)
	binary.LittleEndian.PutUint64(toVerify, uint64(len(p)))
	toVerify = append(toVerify, p...)

	assert.Equal(t, toVerify, s)
}
