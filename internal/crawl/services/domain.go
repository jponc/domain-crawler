package services

import "github.com/jponc/domain-crawler/internal/extractor"

type SuccessCrawlResult struct {
	URL              string
	Title            string
	MetaDescriptions []string
	Links            []string
	KeywordCounts    map[string]int
}

type ErrorCrawlResult struct {
	URL   string
	Error string
}

func transformExtractResultToCrawlResult(result *extractor.ExtractResult) *SuccessCrawlResult {
	return &SuccessCrawlResult{
		URL:              result.URL,
		Title:            result.Title,
		MetaDescriptions: result.MetaDescriptions,
		Links:            result.Links,
		KeywordCounts:    result.KeywordCounts,
	}
}
