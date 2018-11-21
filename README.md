# Go FileCache

Simple file-based cache.

## Usage

```go
package main

import "gitlab.com/kukymbrgo/filecache"

func main()  {
	
    // Initialize cache
    c, err := filecache.New("/path/to/cache/dir", "")
    if err != nil {
    	panic(err)
    } 
    
    // ... make some reader
}
```