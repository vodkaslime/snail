package main

import (
	"errors"
	"fmt"
)

type FS interface {
	Write(src []byte) error
}

type LocalFS struct {
	cfg   *Config
	table Table
}

func NewLocalFS(cfg *Config) LocalFS {
	return LocalFS{
		cfg:   cfg,
		table: newTable(cfg.TableSize),
	}
}

func (f *LocalFS) Write(src *[]byte) (int, error) {
	p := encodePayload(src)
	l := len(p)
	if l > f.cfg.TableSize {
		return 0, fmt.Errorf("oversized payload")
	}

	wl, err := f.table.write(f.cfg, &p)
	if err != nil {
		if errors.Is(err, errorTableFull) {
			// Handle case that table is full
			if err := f.table.flush(f.cfg.FilePath); err != nil {
				return 0, err
			}

			if err := f.rotateTable(); err != nil {
				return 0, err
			}

			return f.table.write(f.cfg, &p)
		} else {
			return 0, err
		}
	} else {
		return wl, nil
	}
}

func (f *LocalFS) rotateTable() error {

	err := f.table.flush(f.cfg.FilePath)
	if err != nil && !errors.Is(err, errorTableFlushed) {
		return err
	}

	t := newTable(f.cfg.TableSize)
	f.table = t

	return nil
}
