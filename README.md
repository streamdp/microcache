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
	ctx := context.Background()
	cache := microcache.New(ctx, nil)
	_ = cache.Set(ctx, "key1", "val1", time.Hour)

	fmt.Println(cache.Get(ctx, "key1"))
}
```
Or create your own cache based on this solution, microCache implements the following interface:
```go
type Cache interface {
    Get(ctx context.Context, key string) (any, error)
    Set(ctx context.Context, key string, value any, expiration time.Duration) error
}
```
Look at the examples folder for explanations.