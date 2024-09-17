package extractor

type KeywordCounts map[string]int

type ExtractResult struct {
	URL              string
	Title            string
	MetaDescriptions []string
	Links            []string
	KeywordCounts    KeywordCounts
}
