# Go FileCache

`filecache` is a package to write data from reader to cache files 
with blackjack and metadata.

## Installation

```sh
go get -u gitlab.com/kukymbrgo/filecache
```

## Usage

```go
package main

import (
	"gitlab.com/kukymbrgo/filecache"
	"io"
	"io/ioutil"
)

func main()  {
	
    // Set defaults for all instances:
    // items namespace
    filecache.NamespaceDefault = "dft"
    
    // default cache files extension
    filecache.ExtDefault = ".cache"
    
    // default time-to-live in seconds; set -1 to eternal
    filecache.TTLDefault = -1
	
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
    
    var reader io.Reader
    
    item, err := fc.Read(pageUrl, "")
    
    if err != nil {
        // Some slow function call
        reader := getPageDownloaderReader()
        _, err = fc.Write(&filecache.Meta{Key: pageUrl}, reader)
        if err != nil {
            // If failed to cache, handle the error       
            panic(err)
        }
    } else {
        reader = item.Reader
    }
    
    // Do some stuff
    _, _ = ioutil.ReadAll(reader)
}
```