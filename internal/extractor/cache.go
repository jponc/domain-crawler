package extractor

// NOTE: This is a very basic cache implementation and is not thread safe.
// I also didn't introduce any TTL for the cache.

type extractorCache struct {
	cache map[string]*ExtractResult
}

func NewExtractorCache() *extractorCache {
	return &extractorCache{
		cache: make(map[string]*ExtractResult),
	}
}

func (c *extractorCache) Get(url string) (*ExtractResult, bool) {
	result, exists := c.cache[url]
	return result, exists
}

func (c *extractorCache) Set(url string, result *ExtractResult) {
	c.cache[url] = result
}
