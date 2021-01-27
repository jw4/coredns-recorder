package main

import "io"

type nopCloser struct{ inner io.Writer }

func (n *nopCloser) Close() error                { return nil }
func (n *nopCloser) Write(p []byte) (int, error) { return n.inner.Write(p) }
