package main

import (
	"fmt"
	"io"
	"os"
)

type writer struct {
	out      io.WriteCloser
	filename string
}

func (w *writer) Close() error {
	if w == nil || w.out == nil {
		return nil
	}

	return w.out.Close()
}

func (w *writer) Write(p []byte) (int, error) {
	if w == nil {
		return 0, fmt.Errorf("nil: %w", os.ErrInvalid)
	}

	if w.out == nil {
		if len(w.filename) == 0 {
			w.out = &nopCloser{inner: os.Stdout}
		} else {
			o, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
			if err != nil {
				return 0, fmt.Errorf("opening file %s: %w", w.filename, err)
			}
			w.out = o
		}
	}

	return w.out.Write(p)
}
