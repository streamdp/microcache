package microcache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// defaultCheckInterval sets the default value of memory checks for expired entries.
const defaultCheckInterval = 30000 * time.Millisecond

var ErrKeyNotFound = errors.New("key is missing in the cache")

type record struct {
	value     any
	expiredAt int
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

// New create a new one micro cache instance, "checkInterval" sets how often (in milliseconds) check memory
// map for the expired entries. Set "checkInterval" = 0 to use default value or -1 to disable expiration check.
//
//	// check for expired entries every minutes
//	cache := microcache.New(context.Background(), 60000);
//
//	// set check interval to the default value (30 seconds)
//	cache := microcache.New(context.Background(), 0);
//
//	// disable expiration check
//	cache := microcache.New(context.Background(), -1);
func New(ctx context.Context, checkInterval int) *MicroCache {
	c := &MicroCache{
		ctx: ctx,

		c:  map[string]*record{},
		mu: &sync.RWMutex{},

		checkInterval: defaultCheckInterval,
	}

	if checkInterval > 0 {
		c.checkInterval = time.Duration(checkInterval) * time.Millisecond
	}

	if checkInterval != -1 {
		go c.processExpiration()
	}

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

			if expired, ok := m.getExpired(); ok {
				m.mu.Lock()
				for k := range expired {
					delete(m.c, k)
				}
				m.mu.Unlock()
			}
		}
	}
}

func (m *MicroCache) getExpired() (expired map[string]struct{}, ok bool) {
	expired = map[string]struct{}{}

	now := int(time.Now().UnixMilli())

	m.mu.RLock()
	for k, v := range m.c {
		if now > v.expiredAt {
			expired[k] = struct{}{}
		}
	}
	m.mu.RUnlock()

	return expired, len(expired) > 0
}

// Get entry from the cache by "key" if present, otherwise it returns ErrKeyNotFound error.
func (m *MicroCache) Get(key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if r, ok := m.c[key]; ok {
		return r.value, nil
	}

	return nil, ErrKeyNotFound
}

// Set the entry to cache, "expiration" interval determines how long the entry will remain in the cache.
func (m *MicroCache) Set(key string, value any, expiration time.Duration) error {
	m.mu.Lock()
	m.c[key] = &record{
		value:     value,
		expiredAt: int(time.Now().UnixMilli() + expiration.Milliseconds()),
	}
	m.mu.Unlock()

	return nil
}

// Delete the entry from cache immediately.
func (m *MicroCache) Delete(key string) {
	m.mu.Lock()
	delete(m.c, key)
	m.mu.Unlock()
}
