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
	"github.com/dsnet/compress/bzip2"
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
	flag.Parse()

	var b datasize.ByteSize
	err := b.UnmarshalText([]byte(*sampleSizeStr))
	check(err)
	sampleSize := int64(b.Bytes())

	if len(flag.Args()) < 1 {
		usage()
	}

	// open file for reading
	fmt.Printf("opening %s\n", flag.Args()[0])
	f, err := os.Open(flag.Args()[0])
	defer f.Close()
	check(err)

	// get file size
	fi, err := f.Stat()
	check(err)

	size := fi.Size()
	fmt.Printf("file size: %d bytes (%s)\n", size, datasize.ByteSize(size).HR())

	// Create a wrapper to the file reader and writer
	cf := NewCountReader(f)

	// gzip
	gzcount := NewCountWriter()
	gz := NewTimedWriter(gzip.NewWriter(gzcount))

	// lzw
	lzcount := NewCountWriter()
	lz := NewTimedWriter(lzw.NewWriter(lzcount, lzw.LSB, 8))

	// zlib
	zlcount := NewCountWriter()
	zl := NewTimedWriter(zlib.NewWriter(zlcount))

	// bzip2
	bzcount := NewCountWriter()
	bzwriter, err := bzip2.NewWriter(bzcount, &bzip2.WriterConfig{ Level: 9 })
	check(err)
	bz := NewTimedWriter(bzwriter)

	w := io.MultiWriter(gz, lz, zl, bz)

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
		_, err = w.Write(buf)
		check(err)
	}
	// close and flush
	check(gz.Close())
	check(lz.Close())
	check(zl.Close())
	check(bz.Close())

	// print stats
	report("gzip", cf, gzcount, gz)
	report("lzw", cf, lzcount, lz)
	report("zlib", cf, zlcount, zl)
	report("bzip2", cf, bzcount, bz)
}

func report(algorithm string, cr *CountReader, cw *CountWriter, tw *TimedWriter) {
	fmt.Printf("%s:\tread: %d bytes (%s)\twritten: %d bytes (%s)\tcompress rate: %f%%\trate:%s/s\n",
		algorithm,
		cr.BytesRead,
		datasize.ByteSize(cr.BytesRead).HR(),
		cw.BytesWritten,
		datasize.ByteSize(cw.BytesWritten).HR(),
		(float64(cr.BytesRead)/float64(cw.BytesWritten)  * 100),
		datasize.ByteSize(float64(cw.BytesWritten)/tw.Elapsed.Seconds()).HR())
}
