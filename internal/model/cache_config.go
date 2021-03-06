package model

import (
	"errors"
)

// CacheConfig defines cache configuration params for a service
type CacheConfig struct {
	TTL      int64    `json:"ttl"`
	Statuses []int    `json:"statuses"`
	Tags     []string `json:"tags"`
}

// DefaultCacheConfig is the value by default
var DefaultCacheConfig = CacheConfig{
	TTL:      0,
	Statuses: []int{},
	Tags:     []string{},
}

// IsEmpty checks if a cache is empty
func (c CacheConfig) IsEmpty() bool {
	return c.TTL == 0 && len(c.Statuses) == 0 && len(c.Tags) == 0
}

// Validate checks whether a cache is valid
func (c CacheConfig) Validate() error {
	if c.IsEmpty() {
		return nil
	}
	if c.TTL > 0 && len(c.Statuses) > 0 && len(c.Tags) > 0 {
		return nil
	}
	return ErrInvalidCacheConfig
}

// ErrCacheConfigNotFound error for not found service
var ErrCacheConfigNotFound = errors.New("Cache config not found")

// ErrInvalidCacheConfig error for invalid cache
var ErrInvalidCacheConfig = errors.New("Invalid cache config")
