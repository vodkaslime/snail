package main

import (
	"errors"
	"os"
)

var (
	errorTableFull    = errors.New("table is full")
	errorTableFlushed = errors.New("table is already flushed")
)

type Table struct {
	ptr       int
	cap       int
	isFlushed bool

	buf []byte
}

func newTable(cap int) Table {
	return Table{
		ptr:       0,
		cap:       cap,
		isFlushed: false,

		buf: make([]byte, cap),
	}
}

func (t *Table) flush(filePath string) error {
	if t.isFlushed {
		return errorTableFlushed
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(t.buf[:t.ptr]); err != nil {
		return err
	}

	t.isFlushed = true
	return nil
}

func (t *Table) space(cfg *Config) int {
	return cfg.TableSize - t.ptr
}

func (t *Table) write(cfg *Config, src *[]byte) (int, error) {

	if t.isFlushed {
		return 0, errorTableFlushed
	}

	l := len(*src)

	if int(l) > cfg.TableSize-t.ptr {
		return 0, errorTableFull
	}
	n := copy(t.buf[t.ptr:], *src)

	t.ptr += n

	return int(n), nil
}
