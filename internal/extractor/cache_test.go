package extractor_test

import (
	"testing"

	"github.com/jponc/domain-crawler/internal/extractor"
	"github.com/stretchr/testify/require"
)

func TestExtractorCache_Get(t *testing.T) {
	// Setup cache
	cache := extractor.NewExtractorCache()
	cache.Set("http://example.com", &extractor.ExtractResult{
		Title:            "Example Domain",
		MetaDescriptions: []string{"This is an example domain"},
		Links:            []string{"http://example.com/link1", "http://example.com/link2"},
		KeywordCounts: map[string]int{
			"keyword1": 1,
			"keyword2": 2,
		},
	})

	tests := []struct {
		name           string
		url            string
		expectedResult *extractor.ExtractResult
		expectedExists bool
	}{
		{
			name:           "returns no result when url is not in cache",
			url:            "http://not-found-in-cache.com",
			expectedResult: nil,
			expectedExists: false,
		},
		{
			name: "returns result when url is in cache",
			url:  "http://example.com",
			expectedResult: &extractor.ExtractResult{
				Title:            "Example Domain",
				MetaDescriptions: []string{"This is an example domain"},
				Links:            []string{"http://example.com/link1", "http://example.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 1,
					"keyword2": 2,
				},
			},
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := cache.Get(tt.url)
			require.Equal(t, tt.expectedResult, result)
			require.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestExtractorCache_Set(t *testing.T) {
	// Setup cache
	cache := extractor.NewExtractorCache()
	cache.Set("http://example.com", &extractor.ExtractResult{
		Title:            "Example Domain",
		MetaDescriptions: []string{"This is an example domain"},
		Links:            []string{"http://example.com/link1", "http://example.com/link2"},
		KeywordCounts: map[string]int{
			"keyword1": 1,
			"keyword2": 2,
		},
	})

	tests := []struct {
		name           string
		url            string
		result         *extractor.ExtractResult
		expectedResult *extractor.ExtractResult
	}{
		{
			name: "sets result when url is not in cache",
			url:  "http://not-found-in-cache.com",
			result: &extractor.ExtractResult{
				Title:            "Not Found",
				MetaDescriptions: []string{"This is not found"},
				Links:            []string{"http://not-found-in-cache.com/link1", "http://not-found-in-cache.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 1,
					"keyword2": 2,
				},
			},
			expectedResult: &extractor.ExtractResult{
				Title:            "Not Found",
				MetaDescriptions: []string{"This is not found"},
				Links:            []string{"http://not-found-in-cache.com/link1", "http://not-found-in-cache.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 1,
					"keyword2": 2,
				},
			},
		},
		{
			name: "overwrites result when url is in cache",
			url:  "http://example.com",
			result: &extractor.ExtractResult{
				Title:            "New Example Domain",
				MetaDescriptions: []string{"This is an example domain"},
				Links:            []string{"http://example.com/link1", "http://example.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 100,
					"keyword2": 200,
				},
			},
			expectedResult: &extractor.ExtractResult{
				Title:            "New Example Domain",
				MetaDescriptions: []string{"This is an example domain"},
				Links:            []string{"http://example.com/link1", "http://example.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 100,
					"keyword2": 200,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set first
			cache.Set(tt.url, tt.result)

			// Then assert on get
			result, exists := cache.Get(tt.url)
			require.Equal(t, tt.expectedResult, result)
			require.True(t, exists)
		})
	}
}
