# 📦 FileCache `v2`

[![Make](https://github.com/kukymbr/filecache/actions/workflows/make.yml/badge.svg)](https://github.com/kukymbr/filecache/actions/workflows/make.yml)
[![GoDoc](https://godoc.org/github.com/kukymbr/filecache?status.svg)](https://godoc.org/github.com/kukymbr/filecache)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/filecache)](https://goreportcard.com/report/github.com/kukymbr/filecache)

Store data from io.Reader or bytes to cache files with TTL and metadata.

## Installation

```sh
go get github.com/kukymbr/filecache/v2 
```

## Usage

### Initializing the cache instance

```go
// With target dir specified
fc, err := filecache.New("/path/to/cache/dir")
```

```go
// With temp dir as a target
fc, err := filecache.New("")
```

```go
// With options
fc, err := filecache.New(
	"/path/to/cache/dir",
    filecache.InstanceOptions{
        PathGenerator: filecache.FilteredKeyPath,
        DefaultTTL:    time.Hour,
        GCDivisor:     10,
    },
)
```

See the [`InstanceOptions` godoc](options.go) for the instance configuration values.

### Saving data to the cache

```go
// From the io.Reader
_, err := fc.Write(context.Background(), "key1", strings.NewReader("value1"))
```

```go
// From the byte array
_, err := fc.WriteData(context.Background(), "key2", []byte("value2"))
```

```go
// With the item options
_, err := fc.Write(
    context.Background(), 
    "key3", 
    strings.NewReader("value3"),
    filecache.ItemOptions{
        Name:   "Key 3",
        TTL:    time.Hour * 24,
        Fields: filecache.NewValues("field1", "val1", "field2", "val2"),
    },
)
```

See the [`ItemOptions` godoc](options.go) for the instance configuration values.

### Reading from cache

```go
// Opening the cache file reader
res, err := fc.Open(context.Background(), "key1")
reader := res.Reader()
```

```go
// Read all the data
res, err := fc.Read(context.Background(), "key2")
data := res.Data()
```

```go
// Read options
res, err := fc.Read(context.Background(), "key3")
name := res.Options().Name
```

### Iterate through the cached items

To iterate through the cached items, use the `Scanner` tool:

```go
// Initialize the scanner
scanner := filecache.NewScanner(fc.GetPath())

// Iterate
err = scanner.Scan(func(entry filecache.ScanEntry) error {
    // Do some nice things  
    return nil
})
```

## License

[MIT](/LICENSE).