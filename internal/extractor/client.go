package extractor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cache interface {
	Get(url string) (*ExtractResult, bool)
	Set(url string, result *ExtractResult)
}

type client struct {
	httpClient  *http.Client
	resultCache cache
	logger      zerolog.Logger
}

func NewExtractorClient(httpClient *http.Client, resultCache cache) *client {
	return &client{
		httpClient:  httpClient,
		resultCache: resultCache,
		logger:      log.With().Str("package", "extractor").Str("client", "ExtractorClient").Logger(),
	}
}

func (c *client) Extract(ctx context.Context, url string, keywords []string) (*ExtractResult, error) {
	// Check cache if available
	if cachedResult, exists := c.resultCache.Get(url); exists {
		c.logger.Info().Str("url", url).Msg("Returning cached result")
		return cachedResult, nil
	}

	c.logger.Info().Str("url", url).Msg("Fetching data from origin")
	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", res.StatusCode, res.Status)
	}

	// Parse HTML doc
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	// Extract Title
	title := doc.Find("title").Text()

	// Extract Meta descriptions
	metaDescriptions := []string{}

	doc.Find("meta[name=description]").Each(func(i int, s *goquery.Selection) {
		if _, exists := s.Attr("content"); exists {
			metaDescriptions = append(metaDescriptions, s.AttrOr("content", ""))
		}
	})

	// Extract links
	links := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, href)
		}
	})

	// Extract keyword counts
	keywordCounts := getKeywordCounts(doc, keywords)

	result := ExtractResult{
		URL:              url,
		Title:            title,
		MetaDescriptions: metaDescriptions,
		Links:            links,
		KeywordCounts:    keywordCounts,
	}

	// Store result to cache
	c.resultCache.Set(url, &result)

	// Return result
	return &result, nil
}

func getKeywordCounts(doc *goquery.Document, keywords []string) KeywordCounts {
	keywordCounts := KeywordCounts{}

	for _, keyword := range keywords {
		count := strings.Count(doc.Text(), keyword)
		keywordCounts[keyword] = count
	}

	return keywordCounts
}
