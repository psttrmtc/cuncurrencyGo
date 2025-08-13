package main

import (
	ncache "concurrency/05-non-blocking-cache/memo"
	"fmt"
	"reflect"
)

type Client interface {
	Get(address string) (string, error)
}

type Cache struct {
	client Client
	memo   *ncache.Memo
	// You can add new fields if needed
}

// Don't update signature of NewCache
func NewCache(client Client) *Cache {
	// TODO: Implement
	m := ncache.New(func(address string) (any, error) {
		return client.Get(address)
	})
	return &Cache{client: client, memo: m}
}

func (c *Cache) Close() {
	c.memo.Close()
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	// TODO: Implement. Right now it doesn't cache
	val, err := c.memo.Get(address)
	if err != nil {
		return "", err
	}
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("function signature don't return string %d", reflect.TypeOf(c.client.Get))
	}
	return s, nil
}
