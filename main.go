package main

import (
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/c2h5oh/datasize"
)

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cmprs [flags] <file>\n")
	fmt.Fprintln(os.Stderr, "estimate compression rate of large files\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	samples := flag.Int("nsamples", 100, "number of samples")
	sampleSizeStr := flag.String("ssize", "1kb", "sample size")

	var b datasize.ByteSize
	err := b.UnmarshalText([]byte(*sampleSizeStr))
	check(err)
	sampleSize := int64(b.Bytes())

	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
	}

	// open the file
	fmt.Printf("opening file %s\n", flag.Args()[0])
	f, err := os.Open(flag.Args()[0])
	check(err)

	// get file size
	fi, err := f.Stat()
	check(err)

	size := fi.Size()
	fmt.Printf("file size: %d (%s)\n", size, datasize.ByteSize(size).String())

	// Create a wrapper to the file reader
	cf := NewCountReader(f)

	// gzip
	gzcount := NewCountWriter()
	gz := gzip.NewWriter(gzcount)

	// lzw
	lzcount := NewCountWriter()
	lz := lzw.NewWriter(lzcount, lzw.LSB, 8)

	// zlib
	zlcount := NewCountWriter()
	zl := zlib.NewWriter(zlcount)

	buf := make([]byte, sampleSize)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < *samples; i++ {
		// find random place to seek to:
		loc := rand.Int63n(size - sampleSize)

		//fmt.Printf("sample %d: seeking to %d\n", i, loc)
		_, err := f.Seek(loc, 0)
		check(err)

		// read data from file
		_, err = io.ReadFull(cf, buf)
		check(err)

		// compress the data read
		_, err = gz.Write(buf)
		check(err)
		_, err = lz.Write(buf)
		check(err)
		_, err = zl.Write(buf)
		check(err)
	}
	// close and flush
	err = gz.Close()
	check(err)
	err = lz.Close()
	check(err)
	err = zl.Close()
	check(err)

	// print stats
	report("gzip", cf, gzcount)
	report("lzw", cf, lzcount)
	report("zlib", cf, zlcount)
}

func report(algorithm string, cr *CountReader, cw *CountWriter) {
	fmt.Printf("%s:\tread: %d (%s)\t\twritten: %d (%s)\t\tcompress rate: %f%%\n",
		algorithm,
		cr.BytesRead,
		datasize.ByteSize(cr.BytesRead).String(),
		cw.BytesWritten,
		datasize.ByteSize(cw.BytesWritten).String(),
		(float64(cw.BytesWritten) / float64(cr.BytesRead) * 100))
}
