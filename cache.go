package microcache

import (
	"context"
	"errors"
	"maps"
	"sync"
	"time"
)

const defaultCheckInterval = 10 * time.Second

var errKeyNotFound = errors.New("key is missing in the cache")

type record struct {
	value     any
	expiredAt time.Time
}

type MicroCache struct {
	ctx context.Context

	c  map[string]*record
	mu *sync.RWMutex

	checkInterval time.Duration
}

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

func (m *MicroCache) Get(_ context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if r, ok := m.c[key]; ok {
		return r.value, nil
	}

	return nil, errKeyNotFound
}

func (m *MicroCache) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.c[key] = &record{
		value:     value,
		expiredAt: time.Now().Add(expiration),
	}

	return nil
}
