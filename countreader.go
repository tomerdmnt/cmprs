package main

import "io"

type CountReader struct {
	BytesRead    int
	reader       io.Reader
}

func NewCountReader(r io.Reader) *CountReader {
	c := &CountReader{
		BytesRead:    0,
	}
	c.reader = r
	return c
}

func (c *CountReader) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	c.BytesRead += n
	return n, err
}

