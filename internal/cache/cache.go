package cache

import "errors"

var CacheEngine Cache

var cache map[string]Cache = map[string]Cache{}

type CacheFactory func() Cache

type Cache interface {
	Connect() error
	Disconnect() error
	GetString(string) string
	GetBytes(string) []byte
	GetInt(string) int

	Set(string, any)

	Delete(string)
}

func RegisterCache(name string, c Cache) {
	cache[name] = c
}

func UseCache(name string) (Cache, error) {
	c, ok := cache[name]
	if !ok {
		return nil, errors.New("cache engine not found")
	}
	CacheEngine = c
	return c, nil
}
