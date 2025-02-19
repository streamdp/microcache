package microcache

import (
	"context"
	"fmt"
	"time"
)

func ExampleMicroCache_Get() {
	cache := New(context.Background(), 0)

	if err := cache.Set("key1", "val1", time.Hour); err != nil {
		fmt.Println("couldn't add value to the cache:", err)
		return
	}

	v, err := cache.Get("key1")
	if err != nil {
		fmt.Println("key not found")
	}

	fmt.Println("v=", v)

	// Output:
	// v= val1
}

func ExampleMicroCache_Set() {
	var (
		cache = New(context.Background(), 0)
		val1  = map[string]int{"temp": 101}
	)

	if err := cache.Set("key1", val1, time.Minute); err != nil {
		fmt.Println("couldn't add value to the cache:", err)
		return
	}

	v, err := cache.Get("key1")
	if err != nil {
		fmt.Println("key not found")
	}

	fmt.Println("v=", v)

	// Output:
	// v= map[temp:101]
}
