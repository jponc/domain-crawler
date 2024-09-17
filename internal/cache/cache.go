package cache

// NOTE: This is a very basic cache implementation and is not thread safe.
// I also didn't introduce any TTL for the cache.

type cache struct {
	cache map[string]string
}

func NewInMemoryCache() *cache {
	return &cache{
		cache: make(map[string]string),
	}
}

func (c *cache) Get(k string) (string, bool) {
	v, ok := c.cache[k]
	return v, ok
}

func (c *cache) Set(k string, v string) {
	c.cache[k] = v
}
