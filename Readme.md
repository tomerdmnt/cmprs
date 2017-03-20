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
$ ./cmprs -nsamples 30 -ssize 2kb cmprs 
opening cmprs
file size: 2443751 bytes (2.330543.1 MB)
gzip:   read: 30720 bytes (30.000000.1 KB)              written: 12102 bytes (11.818359.1 KB)           compress rate: 253.842340%
lzw:    read: 30720 bytes (30.000000.1 KB)              written: 17743 bytes (17.327148.1 KB)           compress rate: 173.138703%
zlib:   read: 30720 bytes (30.000000.1 KB)              written: 12090 bytes (11.806641.1 KB)           compress rate: 254.094293%
bzip2:  read: 30720 bytes (30.000000.1 KB)              written: 13083 bytes (12.776367.1 KB)           compress rate: 234.808530%
```

