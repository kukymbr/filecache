## Go FileCache

[![Build Status](https://gitlab.com/kukymbrgo/filecache/badges/master/build.svg)](https://gitlab.com/kukymbrgo/filecache/pipelines)
[![Coverage](https://gitlab.com/kukymbrgo/filecache/badges/master/coverage.svg)](https://gitlab.com/kukymbrgo/filecache)
[![GoDoc](https://godoc.org/gitlab.com/kukymbrgo/filecache?status.svg)](https://godoc.org/gitlab.com/kukymbrgo/filecache)

Store data from any io.Reader to cache files with TTL and metadata.

### Installation

```sh
go get -u gitlab.com/kukymbrgo/filecache
```

### Usage

```go
package main

import (
    "gitlab.com/kukymbrgo/filecache"
    "io/ioutil"
)

func main()  {
	
    // Set defaults for all instances (optional of course):
    // items namespace
    filecache.NamespaceDefault = "dft"
    
    // default cache files extension
    filecache.ExtDefault = ".cache"
    
    // default time-to-live in seconds; set -1 to eternal
    filecache.TTLDefault = -1
    
    // Set garbage collector run probability divisor
    // (e.g. 10 is 1/10 probability), optional
    filecache.GCDivisor = 10
	
    // Initialize cache instance
    fc, err := filecache.New("/path/to/cache/dir")
    if err != nil {
    	panic(err)
    } 
    
    // Set instance defaults:
    fc.NamespaceDefault = "wiki"
    fc.TTLDefault = 3600
    fc.Ext = ".html"
    
    // Read and write some data 
    
    pageUrl := "https://en.wikipedia.org/wiki/Main_Page"
    
    item, err := fc.Read(pageUrl, "")
    
    if err != nil {
        // Some slow function call
        downloader := getPageDownloaderReader()
        item, _, err = fc.WriteOpen(&filecache.Meta{Key: pageUrl}, downloader)
        if err != nil {
            // If failed to cache, handle the error       
            panic(err)
        }
    }
    
    // Do some stuff
    _, _ = ioutil.ReadAll(item.File)
}
```

### License

MIT. See the [LICENSE](/LICENSE) file.