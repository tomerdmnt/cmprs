package main

type CountWriter struct {
	BytesWritten int
}

func NewCountWriter() *CountWriter {
	return &CountWriter{
		BytesWritten: 0,
	}
}

func (c *CountWriter) Write(p []byte) (int, error) {
	c.BytesWritten += len(p)
	return len(p), nil
}
