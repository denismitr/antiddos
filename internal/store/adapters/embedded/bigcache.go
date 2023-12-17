package embedded

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

type Store struct {
	bc *bigcache.BigCache
}

func New(ctx context.Context, eviction uint64) (*Store, error) {
	bc, err := bigcache.New(ctx, bigcache.DefaultConfig(time.Duration(eviction)*time.Second))
	if err != nil {
		return nil, err
	}

	return &Store{
		bc: bc,
	}, nil
}

func (s *Store) Validate(key string) bool {
	_, err := s.bc.Get(key)
	if err != nil {
		return false
	}
	return true
}

func (s *Store) Remember(key string) {
	_ = s.bc.Set(key, nil)
}
