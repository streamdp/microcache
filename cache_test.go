package microcache

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMicroCache_getExpired(t *testing.T) {
	tests := []struct {
		name        string
		cache       map[string]*record
		wantExpired map[string]struct{}
		wantOk      bool
	}{
		{
			name: "get expired",
			cache: map[string]*record{
				"key1": {
					value:     "val1",
					expiredAt: int(time.Now().Add(-time.Minute).UnixMilli()),
				},
				"key2": {
					value:     "val2",
					expiredAt: int(time.Now().Add(-time.Minute).UnixMilli()),
				},
				"key3": {
					value:     "val3",
					expiredAt: int(time.Now().Add(+time.Minute).UnixMilli()),
				},
			},
			wantExpired: map[string]struct{}{
				"key1": {},
				"key2": {},
			},
			wantOk: true,
		},
		{
			name: "nothing expired yet",
			cache: map[string]*record{
				"key1": {
					value:     "val1",
					expiredAt: int(time.Now().Add(+time.Minute).UnixMilli()),
				},
				"key2": {
					value:     "val2",
					expiredAt: int(time.Now().Add(+time.Minute).UnixMilli()),
				},
				"key3": {
					value:     "val3",
					expiredAt: int(time.Now().Add(+time.Minute).UnixMilli()),
				},
			},
			wantExpired: map[string]struct{}{},
			wantOk:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MicroCache{
				ctx:           context.Background(),
				c:             tt.cache,
				mu:            &sync.RWMutex{},
				checkInterval: time.Minute,
			}
			gotExpired, gotOk := m.getExpired()
			if !reflect.DeepEqual(gotExpired, tt.wantExpired) {
				t.Errorf("getExpired() gotExpired = %v, want %v", gotExpired, tt.wantExpired)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getExpired() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMicroCache_Get(t *testing.T) {
	tests := []struct {
		name    string
		cache   map[string]*record
		key     string
		want    any
		wantErr bool
	}{
		{
			name: "key is present",
			cache: map[string]*record{
				"key1": {
					value:     "val1",
					expiredAt: int(time.Now().Add(time.Minute).UnixMilli()),
				},
			},
			key:     "key1",
			want:    "val1",
			wantErr: false,
		},
		{
			name:    "get error",
			cache:   map[string]*record{},
			key:     "key1",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MicroCache{
				ctx:           context.Background(),
				c:             tt.cache,
				mu:            &sync.RWMutex{},
				checkInterval: time.Minute,
			}
			got, err := m.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
