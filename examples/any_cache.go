package main

import (
	"context"
	"fmt"
	"time"

	"github.com/streamdp/microcache"
)

const cacheReadTimeout = time.Second

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

type AnyCache struct {
	ctx context.Context
	c   Cache
}

func NewAnyCache(ctx context.Context, c Cache) *AnyCache {
	return &AnyCache{
		ctx: ctx,
		c:   c,
	}
}

func (a *AnyCache) Set(key string, value any) (err error) {
	ctx, cancel := context.WithTimeout(a.ctx, cacheReadTimeout)
	defer cancel()

	if err = a.c.Set(ctx, key, value, time.Hour); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func (a *AnyCache) Get(key string) (result any, err error) {
	ctx, cancel := context.WithTimeout(a.ctx, cacheReadTimeout)
	defer cancel()

	if result, err = a.c.Get(ctx, key); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return
}

func main() {
	ctx := context.Background()

	c := NewAnyCache(ctx, microcache.New(ctx, nil))
	_ = c.Set("key1", "val1")

	fmt.Println(c.Get("key1"))
}
