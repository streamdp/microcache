## Microcache
Just a simple implementation of an in-memory cache with entry expiration

### Usage
You can use it directly:
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/streamdp/microcache"
)

func main() {
	cache := microcache.New(context.Background(), nil)
	_ = cache.Set("key1", "val1", time.Hour)

	fmt.Println(cache.Get("key1"))
}
```
Or create your own cache based on this solution, microCache implements the following interface:
```go
type Cache interface {
    Get(key string) (any, error)
    Set(key string, value any, expiration time.Duration) error
}
```
Look at the examples folder for explanations.