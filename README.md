# 📦 FileCache `v2`

[![Build](https://github.com/kukymbr/filecache/actions/workflows/validate.yml/badge.svg)](https://github.com/kukymbr/filecache/actions/workflows/validate.yml)
[![GoDoc](https://godoc.org/github.com/kukymbr/filecache/v2?status.svg)](https://godoc.org/github.com/kukymbr/filecache/v2)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/filecache/v2)](https://goreportcard.com/report/github.com/kukymbr/filecache/v2)

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
fc, err := filecache.NewInTemp()
```

```go
// With options
fc, err := filecache.New(
    "/path/to/cache/dir",
    filecache.InstanceOptions{
        PathGenerator: filecache.FilteredKeyPath,
        DefaultTTL:    time.Hour,
        GC:            filecache.NewIntervalGarbageCollector("/path/to/cache/dir", time.Hour),
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
if err != nil { 
    // Handle the error...
}

if res.Hit() {
    reader := res.Reader()
    // Read the data...
}
```

```go
// Read all the data
res, err := fc.Read(context.Background(), "key2")
if err == nil && res.Hit() {
    data := res.Data()
}
```

```go
// Read options
res, err := fc.Read(context.Background(), "key3")
if err == nil && res.Hit() {
    name := res.Options().Name
}
```

The `Open()` and `Read()` functions return an error only if context is canceled
or if the file open operation has failed. 
If there is no error, this doesn't mean the result is found, the `res.Hit()` function should be called. 

### Iterate through the cached items

To iterate through the cached items, use the `Scanner` tool:

```go
scanner := filecache.NewScanner("/path/to/cache/dir")

err := scanner.Scan(func(entry filecache.ScanEntry) error {
    // Do something with the found entry...
    return nil
})
if err := nil {
    // Handle the error...
}
```

### Removing the expired items

The expired cache items are removed by the `GarbageCollector`, assigned to the `FileCache` instance.

There are three predefined realizations of the `GarbageCollector`:

* `filecache.NewNopGarbageCollector()` — the `GarbageCollector` doing nothing, all the files are kept;
* `filecache.NewProbabilityGarbageCollector()` — the `GarbageCollector` running with the defined probability, used by default;
* `filecache.NewIntervalGarbageCollector()` — the `GarbageCollector` running by the time interval.

See the [gc.go's](gc.go) godocs for more info.

## License

[MIT](/LICENSE).