package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jponc/domain-crawler/internal/crawl/services"
	"github.com/jponc/domain-crawler/internal/extractor"
	"github.com/stretchr/testify/require"
)

// Mocks
type mockExtractorClient struct {
	extractFn func(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error)
}

func (m *mockExtractorClient) Extract(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error) {
	if m != nil && m.extractFn != nil {
		return m.extractFn(ctx, url, keywords)
	}

	return &extractor.ExtractResult{
		URL:              url,
		Title:            "Title",
		MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
		Links:            []string{"https://link1.com", "https://link2.com"},
		KeywordCounts: map[string]int{
			"keyword1": 1,
			"keyword2": 2,
		},
	}, nil
}

// Tests

func TestCrawlService_Crawl(t *testing.T) {
	tests := []struct {
		name                        string
		urls                        []string
		keywords                    []string
		mockExtractorClient         *mockExtractorClient
		expectedSuccessCrawlResults []services.SuccessCrawlResult
		expectedErrorCrawlResults   []services.ErrorCrawlResult
	}{
		{
			name:     "returns error crawl results when failed to extract data from URL",
			urls:     []string{"http://example.com"},
			keywords: []string{"keyword1", "keyword2"},
			mockExtractorClient: &mockExtractorClient{
				extractFn: func(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error) {
					return nil, fmt.Errorf("failed to extract data")
				},
			},
			expectedSuccessCrawlResults: []services.SuccessCrawlResult{},
			expectedErrorCrawlResults: []services.ErrorCrawlResult{
				{
					URL:   "http://example.com",
					Error: "failed to extract data",
				},
			},
		},
		{
			name:     "returns success crawl results when successfully extracted data from URL",
			urls:     []string{"http://example.com"},
			keywords: []string{"keyword1", "keyword2"},
			mockExtractorClient: &mockExtractorClient{
				extractFn: func(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error) {
					return &extractor.ExtractResult{
						URL:              url,
						Title:            "Title",
						MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
						Links:            []string{"https://link1.com", "https://link2.com"},
						KeywordCounts: map[string]int{
							"keyword1": 1,
							"keyword2": 2,
						},
					}, nil
				},
			},
			expectedSuccessCrawlResults: []services.SuccessCrawlResult{
				{
					URL:              "http://example.com",
					Title:            "Title",
					MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
					Links:            []string{"https://link1.com", "https://link2.com"},
					KeywordCounts: map[string]int{
						"keyword1": 1,
						"keyword2": 2,
					},
				},
			},
			expectedErrorCrawlResults: []services.ErrorCrawlResult{},
		},
		{
			name:     "returns success and error crawl results when successfully extracted data from some URLs and failed to extract data from others",
			urls:     []string{"http://example.com", "http://example.com/404"},
			keywords: []string{"keyword1", "keyword2"},
			mockExtractorClient: &mockExtractorClient{
				extractFn: func(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error) {
					if url == "http://example.com" {
						return &extractor.ExtractResult{
							URL:              url,
							Title:            "Title",
							MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
							Links:            []string{"https://link1.com", "https://link2.com"},
							KeywordCounts: map[string]int{
								"keyword1": 1,
								"keyword2": 2,
							},
						}, nil
					}

					return nil, fmt.Errorf("failed to extract data")
				},
			},
			expectedSuccessCrawlResults: []services.SuccessCrawlResult{
				{
					URL:              "http://example.com",
					Title:            "Title",
					MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
					Links:            []string{"https://link1.com", "https://link2.com"},
					KeywordCounts: map[string]int{
						"keyword1": 1,
						"keyword2": 2,
					},
				},
			},
			expectedErrorCrawlResults: []services.ErrorCrawlResult{
				{
					URL:   "http://example.com/404",
					Error: "failed to extract data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawlService := services.NewCrawlService(tt.mockExtractorClient, 1)

			crawlSuccessResults, crawlErrorResults, err := crawlService.Crawl(context.Background(), tt.urls, tt.keywords)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tt.expectedSuccessCrawlResults, crawlSuccessResults)
			require.Equal(t, tt.expectedErrorCrawlResults, crawlErrorResults)
		})
	}
}
