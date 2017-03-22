Estimate compression rate of large files

## Usage

```
usage: cmprs [flags] <file>

estimate compression rate of large files

  -nsamples int
        number of samples (default 100)
  -ssize string
        sample size (default "1kb")
```

## Example

```shell
$./cmprs -nsamples 30 -ssize 2kb ./cmprs 
opening ./cmprs
file size: 2443751 bytes (2.330543.1 MB)
gzip:   read: 61440 bytes (60.000000.1 KB)              written: 27271 bytes (26.631836.1 KB)           compress rate: 225.294269%
lzw:    read: 61440 bytes (60.000000.1 KB)              written: 38289 bytes (37.391602.1 KB)           compress rate: 160.463841%
zlib:   read: 61440 bytes (60.000000.1 KB)              written: 27259 bytes (26.620117.1 KB)           compress rate: 225.393448%
bzip2:  read: 61440 bytes (60.000000.1 KB)              written: 28826 bytes (28.150391.1 KB)           compress rate: 213.140914%
```

