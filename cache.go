package microcache

import (
	"context"
	"errors"
	"maps"
	"sync"
	"time"
)

// defaultCheckInterval sets the default value of memory checks for expired entries.
const defaultCheckInterval = 10 * time.Second

var ErrKeyNotFound = errors.New("key is missing in the cache")

type record struct {
	value     any
	expiredAt time.Time
}

// MicroCache provides cache interface implementation.
//
//	type Cache interface {
//	    Get(ctx context.Context, key string) (any, error)
//	    Set(ctx context.Context, key string, value any, expiration time.Duration) error
//	}
type MicroCache struct {
	ctx context.Context

	c  map[string]*record
	mu *sync.RWMutex

	checkInterval time.Duration
}

// New create a new one micro cache instance, parameter "checkInterval" sets how often check memory map for the
// expired entries.
func New(ctx context.Context, checkInterval *time.Duration) *MicroCache {
	c := &MicroCache{
		ctx: ctx,

		c:  map[string]*record{},
		mu: &sync.RWMutex{},
	}

	if checkInterval != nil {
		c.checkInterval = *checkInterval
	} else {
		c.checkInterval = defaultCheckInterval
	}

	go c.processExpiration()

	return c
}

func (m *MicroCache) processExpiration() {
	t := time.NewTimer(m.checkInterval)
	defer t.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-t.C:
			t.Reset(m.checkInterval)

			if len(m.c) == 0 {
				continue
			}

			now := time.Now()
			m.mu.Lock()
			maps.DeleteFunc(m.c, func(_ string, v *record) bool { return now.After(v.expiredAt) })
			m.mu.Unlock()
		}
	}
}

// Get entry from the cache by "key" if present, otherwise it returns ErrKeyNotFound error.
func (m *MicroCache) Get(_ context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if r, ok := m.c[key]; ok {
		return r.value, nil
	}

	return nil, ErrKeyNotFound
}

// Set the entry to cache, "expiration" interval determines how long the entry will remain in the cache.
func (m *MicroCache) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.c[key] = &record{
		value:     value,
		expiredAt: time.Now().Add(expiration),
	}

	return nil
}
