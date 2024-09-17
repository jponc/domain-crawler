package extractor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cache interface {
	Get(k string) (string, bool)
	Set(k, v string)
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
	html, err := c.fetchHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch html: %w", err)
	}

	// Parse HTML doc
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
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

	// Return result
	return &result, nil
}

func (c *client) fetchHTML(ctx context.Context, url string) (string, error) {
	// Check cache if available
	if cachedHTML, exists := c.resultCache.Get(url); exists {
		c.logger.Info().Str("url", url).Msg("Returning cached HTML")
		return cachedHTML, nil
	}

	c.logger.Info().Str("url", url).Msg("Fetching HTML from origin")
	res, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", res.Status)
	}

	// Read all the data from the ReadCloser
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Store result to cache
	html := string(data)
	c.resultCache.Set(url, html)

	return html, nil
}

func getKeywordCounts(doc *goquery.Document, keywords []string) KeywordCounts {
	keywordCounts := KeywordCounts{}

	for _, keyword := range keywords {
		count := strings.Count(doc.Text(), keyword)
		keywordCounts[keyword] = count
	}

	return keywordCounts
}
