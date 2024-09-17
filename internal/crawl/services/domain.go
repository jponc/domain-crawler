package services

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
