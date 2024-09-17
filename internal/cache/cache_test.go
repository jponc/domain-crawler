package cache_test

import (
	"testing"

	"github.com/jponc/domain-crawler/internal/cache"
	"github.com/stretchr/testify/require"
)

func TestExtractorCache_Get(t *testing.T) {
	// Setup cache
	cache := cache.NewInMemoryCache()
	cache.Set("http://example.com", "<html>Test</html>")

	tests := []struct {
		name           string
		key            string
		expectedValue  string
		expectedExists bool
	}{
		{
			name:           "returns no result when key is not in cache",
			key:            "http://not-found-in-cache.com",
			expectedValue:  "",
			expectedExists: false,
		},
		{
			name:           "returns value when key is in cache",
			key:            "http://example.com",
			expectedValue:  "<html>Test</html>",
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := cache.Get(tt.key)
			require.Equal(t, tt.expectedValue, value)
			require.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestExtractorCache_Set(t *testing.T) {
	// Setup cache
	cache := cache.NewInMemoryCache()
	cache.Set("http://example.com", "<html>Test</html>")

	tests := []struct {
		name          string
		key           string
		value         string
		expectedValue string
	}{
		{
			name:          "sets result when key is not in cache",
			key:           "http://not-found-in-cache.com",
			value:         "<html>Test 2</html>",
			expectedValue: "<html>Test 2</html>",
		},
		{
			name:          "overwrites result when key is in cache",
			key:           "http://example.com",
			value:         "<html>Test 3</html>",
			expectedValue: "<html>Test 3</html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set first
			cache.Set(tt.key, tt.value)

			// Then assert on get
			result, exists := cache.Get(tt.key)
			require.Equal(t, tt.expectedValue, result)
			require.True(t, exists)
		})
	}
}
