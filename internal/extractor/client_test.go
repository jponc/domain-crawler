package extractor_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jponc/domain-crawler/internal/extractor"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

type mockExtractorCache struct {
	getFn func(url string) (*extractor.ExtractResult, bool)
	setFn func(url string, result *extractor.ExtractResult)
}

func (m *mockExtractorCache) Get(url string) (*extractor.ExtractResult, bool) {
	if m != nil && m.getFn != nil {
		return m.getFn(url)
	}

	return nil, false
}

func (m *mockExtractorCache) Set(url string, result *extractor.ExtractResult) {
	if m != nil && m.setFn != nil {
		m.setFn(url, result)
	}
}

func TestClient_Extract(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		keywords           []string
		roundTripFunc      roundTripFunc
		mockExtractorCache *mockExtractorCache
		expectedError      string
		expectedResult     *extractor.ExtractResult
	}{
		{
			name:     "returns cached result when available",
			url:      "http://example.com",
			keywords: []string{"keyword1", "keyword2"},
			mockExtractorCache: &mockExtractorCache{
				getFn: func(url string) (*extractor.ExtractResult, bool) {
					return &extractor.ExtractResult{
						Title:            "Example Domain",
						MetaDescriptions: []string{"This is an example domain"},
						Links:            []string{"http://example.com/link1", "http://example.com/link2"},
						KeywordCounts: map[string]int{
							"keyword1": 1,
							"keyword2": 2,
						},
					}, true
				},
			},
			expectedResult: &extractor.ExtractResult{
				Title:            "Example Domain",
				MetaDescriptions: []string{"This is an example domain"},
				Links:            []string{"http://example.com/link1", "http://example.com/link2"},
				KeywordCounts: map[string]int{
					"keyword1": 1,
					"keyword2": 2,
				},
			},
		},
		{
			name:     "returns err when failed to get url",
			url:      "http://example.com",
			keywords: []string{"keyword1", "keyword2"},
			roundTripFunc: func(r *http.Request) (*http.Response, error) {
				return nil, http.ErrHandlerTimeout
			},
			expectedError: "failed to get url: Get \"http://example.com\": http: Handler timeout",
		},
		{
			name:     "returns err when unexpected status code",
			url:      "http://example.com",
			keywords: []string{"keyword1", "keyword2"},
			roundTripFunc: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     "Not Found",
				}, nil
			},
			expectedError: "unexpected status code: 404 Not Found",
		},
		{
			name:     "returns err when failed to parse html",
			url:      "http://example.com",
			keywords: []string{"keyword1", "keyword2"},
			roundTripFunc: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(errorReader{}),
				}, nil
			},
			expectedError: "failed to parse html: unexpected EOF",
		},
		{
			name: "returns result when successfully parsed the html body",
			url:  "http://example.com",
			keywords: []string{
				"keyword1",
				"keyword2",
			},
			roundTripFunc: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`
						<html>
							<head>
								<title>Example Domain</title>
								<meta name="description" content="This is the first meta description." />
								<meta name="description" content="This is the second meta description, which might be ignored by search engines." />
							</head>
							<body>
								<span>keyword1 keyword1 keyword2</span>

								<a href="http://example.com/link1">Link 1</a>
								<a href="http://example.com/link2">Link 2</a>
							</body>
						</html>
					`)),
				}, nil
			},
			expectedResult: &extractor.ExtractResult{
				URL:   "http://example.com",
				Title: "Example Domain",
				MetaDescriptions: []string{
					"This is the first meta description.",
					"This is the second meta description, which might be ignored by search engines.",
				},
				Links: []string{
					"http://example.com/link1",
					"http://example.com/link2",
				},
				KeywordCounts: map[string]int{
					"keyword1": 2,
					"keyword2": 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			httpClient := &http.Client{
				Transport: tt.roundTripFunc,
			}

			client := extractor.NewExtractorClient(httpClient, tt.mockExtractorCache)

			result, err := client.Extract(ctx, tt.url, tt.keywords)
			if tt.expectedError != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
