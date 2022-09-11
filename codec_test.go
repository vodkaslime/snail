package main

import (
	"fmt"
	"testing"
)

func TestUintConvert(t *testing.T) {
	p := []byte("yo here we go")
	s := encodePayload(&p)
	fmt.Println(s)
}
