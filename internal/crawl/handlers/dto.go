package handlers

import "github.com/jponc/domain-crawler/internal/crawl/services"

// Requsts

type CrawlRequest struct {
	URLs     []string `json:"urls"`
	Keywords []string `json:"keywords"`
}

// Responses

type CrawlResponse struct {
	Results []SuccessResult `json:"results"`
	Errors  []ErrorResult   `json:"errors,omitempty"`
}

// Types

type SuccessResult struct {
	URL              string         `json:"url"`
	Title            string         `json:"title"`
	MetaDescriptions []string       `json:"meta_descriptions"`
	Links            []string       `json:"links"`
	KeywordCounts    map[string]int `json:"keyword_counts"`
}

type ErrorResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// Domain to DTO converters

func convertSuccessCrawlResultsToSuccessResults(crawlResults []services.SuccessCrawlResult) []SuccessResult {
	results := make([]SuccessResult, 0, len(crawlResults))
	for _, crawlResult := range crawlResults {
		result := SuccessResult{
			URL:              crawlResult.URL,
			Title:            crawlResult.Title,
			MetaDescriptions: crawlResult.MetaDescriptions,
			Links:            crawlResult.Links,
			KeywordCounts:    crawlResult.KeywordCounts,
		}
		results = append(results, result)
	}
	return results
}

func convertErrorCrawlResultsToErrorResults(crawlResults []services.ErrorCrawlResult) []ErrorResult {
	results := make([]ErrorResult, 0, len(crawlResults))
	for _, crawlResult := range crawlResults {
		result := ErrorResult{
			URL:   crawlResult.URL,
			Error: crawlResult.Error,
		}
		results = append(results, result)
	}
	return results
}
