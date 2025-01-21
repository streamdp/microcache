package main

import (
	"context"
	"fmt"
	"time"

	"github.com/streamdp/microcache"
)

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
}

type AnyCache struct {
	c Cache
}

func NewAnyCache(c Cache) *AnyCache {
	return &AnyCache{c: c}
}

func (a *AnyCache) Set(key string, value any) (err error) {
	if err = a.c.Set(key, value, time.Hour); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func (a *AnyCache) Get(key string) (result any, err error) {
	if result, err = a.c.Get(key); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return
}

func main() {
	c := NewAnyCache(microcache.New(context.Background(), -1))
	_ = c.Set("key1", "val1")

	v, err := c.Get("key1")
	if err != nil {
		fmt.Println("key not found")
	}
	fmt.Println(v)
}
