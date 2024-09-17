package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/jponc/domain-crawler/api/openapi"
	"github.com/jponc/domain-crawler/internal/crawl/handlers"
	"github.com/jponc/domain-crawler/internal/crawl/services"
	"github.com/jponc/domain-crawler/internal/middlewares"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/require"
)

// Mocks
type mockCrawlService struct {
	crawlFn func(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error)
}

func (m *mockCrawlService) Crawl(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error) {
	if m != nil && m.crawlFn != nil {
		return m.crawlFn(ctx, urls, keywords)
	}

	return []services.SuccessCrawlResult{}, []services.ErrorCrawlResult{}, nil
}

func TestCrawlHandler_Crawl(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		mockCrawlService     *mockCrawlService
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:               "returns 400 when failed to decode request body",
			requestBody:        "invalid",
			mockCrawlService:   &mockCrawlService{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `
				{
					"error": "request body has an error: failed to decode request body: invalid character 'i' looking for beginning of value"
				}`,
		},
		{
			name:               "returns 400 when request body doesn't conform to openapi spec",
			requestBody:        "{}",
			mockCrawlService:   &mockCrawlService{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `
				{
					"error": "request body has an error: doesn't match schema #/components/schemas/CrawlRequest: Error at \"/urls\": property \"urls\" is missing"
				}`,
		},
		{
			name: "returns 500 when crawl service returns an error",
			requestBody: `
				{
					"urls": ["https://example.com"],
					"keywords": ["example"]
				}`,
			mockCrawlService: &mockCrawlService{
				crawlFn: func(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error) {
					return nil, nil, fmt.Errorf("error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: `
				{
					"error": "error"
				}`,
		},
		{
			name: "returns 200 when crawl service returns success results",
			requestBody: `
				{
					"urls": ["https://example.com"],
					"keywords": ["example"]
				}`,
			mockCrawlService: &mockCrawlService{
				crawlFn: func(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error) {
					return []services.SuccessCrawlResult{
						{
							URL:              "https://example.com",
							Title:            "Title",
							MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
							Links:            []string{"https://link1.com", "https://link2.com"},
							KeywordCounts: map[string]int{
								"keyword1": 1,
								"keyword2": 2,
							},
						},
					}, nil, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `
				{
					"results": [
						{
							"url": "https://example.com",
							"title": "Title",
							"meta_descriptions": ["Meta Description 1", "Meta Description 2"],
							"links": ["https://link1.com", "https://link2.com"],
							"keyword_counts": {
								"keyword1": 1,
								"keyword2": 2
							}
						}
					]
				}`,
		},
		{
			name: "removes duplicate urls before sending to crawl service",
			requestBody: `
				{
					"urls": ["https://example.com", "https://example.com", "https://example.com", "https://example.com"],
					"keywords": ["example"]
				}`,
			mockCrawlService: &mockCrawlService{
				crawlFn: func(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error) {
					require.Len(t, urls, 1)
					require.Equal(t, []string{"https://example.com"}, urls)

					return []services.SuccessCrawlResult{
						{
							URL:              "https://example.com",
							Title:            "Title",
							MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
							Links:            []string{"https://link1.com", "https://link2.com"},
							KeywordCounts: map[string]int{
								"keyword1": 1,
								"keyword2": 2,
							},
						},
					}, nil, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `
				{
					"results": [
						{
							"url": "https://example.com",
							"title": "Title",
							"meta_descriptions": ["Meta Description 1", "Meta Description 2"],
							"links": ["https://link1.com", "https://link2.com"],
							"keyword_counts": {
								"keyword1": 1,
								"keyword2": 2
							}
						}
					]
				}`,
		},
		{
			name: "returns 200 when crawl service returns success and error results",
			requestBody: `
				{
					"urls": ["https://example.com", "https://example.com/404"],
					"keywords": ["example"]
				}`,
			mockCrawlService: &mockCrawlService{
				crawlFn: func(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error) {
					successCrawlResults := []services.SuccessCrawlResult{
						{
							URL:              "https://example.com",
							Title:            "Title",
							MetaDescriptions: []string{"Meta Description 1", "Meta Description 2"},
							Links:            []string{"https://link1.com", "https://link2.com"},
							KeywordCounts: map[string]int{
								"keyword1": 1,
								"keyword2": 2,
							},
						},
					}

					errorCrawlResults := []services.ErrorCrawlResult{
						{
							URL:   "https://example.com/404",
							Error: "failed to extract data",
						},
					}

					return successCrawlResults, errorCrawlResults, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `
				{
					"results": [
						{
							"url": "https://example.com",
							"title": "Title",
							"meta_descriptions": ["Meta Description 1", "Meta Description 2"],
							"links": ["https://link1.com", "https://link2.com"],
							"keyword_counts": {
								"keyword1": 1,
								"keyword2": 2
							}
						}
					],
					"errors": [
						{
							"url": "https://example.com/404",
							"error": "failed to extract data"
						}
					]
				}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// initialise router with openapi spec
			openapiSpec, err := openapi.FS.ReadFile(openapi.OpenAPISpecFilename)
			require.NoError(t, err)

			loader := openapi3.NewLoader()
			doc, err := loader.LoadFromData(openapiSpec)
			require.NoError(t, err)

			oapiValidatorMiddleware := middlewares.OpenAPIValidatorMiddleware(doc)
			router := chi.NewRouter()
			router.Use(oapiValidatorMiddleware)

			// initialise handlers
			h := handlers.NewCrawlHandler(tt.mockCrawlService)

			// setup route
			router.Post("/crawl", h.Crawl)

			// create request
			r := httptest.NewRequest(http.MethodPost, "/crawl", strings.NewReader(tt.requestBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			require.Equal(t, tt.expectedStatusCode, w.Code)
			fmt.Println(w.Body.String())
			jsonassert.New(t).Assertf(w.Body.String(), tt.expectedResponseBody)
		})
	}
}
