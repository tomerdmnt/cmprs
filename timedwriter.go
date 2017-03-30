package main

import (
	"io"
	"time"
)

type TimedWriter struct {
	writer  io.WriteCloser
	Elapsed time.Duration
}

func NewTimedWriter(w io.WriteCloser) *TimedWriter {
	return &TimedWriter{w, 0}
}

func (tw *TimedWriter) Write(p []byte) (int, error) {
	t := time.Now()
	n, err := tw.writer.Write(p)
	tw.Elapsed += time.Since(t)
	return n, err
}

func (tw *TimedWriter) Close() error {
	t := time.Now()
	err :=  tw.writer.Close()
	tw.Elapsed += time.Since(t)
	return err
}
