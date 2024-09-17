package services

import (
	"context"

	"github.com/jponc/domain-crawler/internal/extractor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type extractorClient interface {
	Extract(ctx context.Context, url string, keywords []string) (*extractor.ExtractResult, error)
}

type crawlService struct {
	extractorClient extractorClient
	concurrentLimit int
	logger          zerolog.Logger
}

func NewCrawlService(extractorClient extractorClient, concurrentLimit int) *crawlService {
	return &crawlService{
		extractorClient: extractorClient,
		concurrentLimit: concurrentLimit,
		logger:          log.With().Str("package", "services").Str("service", "CrawlService").Logger(),
	}
}

func (s *crawlService) Crawl(ctx context.Context, urls []string, keywords []string) ([]SuccessCrawlResult, []ErrorCrawlResult, error) {
	successCrawlResults := []SuccessCrawlResult{}
	errorCrawlResults := []ErrorCrawlResult{}

	// Define errgroup
	eg, egCtx := errgroup.WithContext(ctx)

	// Set limit for concurrent requests
	eg.SetLimit(s.concurrentLimit)

	// Iterate over urls and extract data
	for _, url := range urls {
		url := url
		eg.Go(func() error {
			s.logger.Info().Str("url", url).Msg("Extracting data from URL")
			result, err := s.extractorClient.Extract(egCtx, url, keywords)
			// Handle error
			if err != nil {
				s.logger.Error().Str("url", url).Msg("Failed to extract data from URL")
				errorCrawlResults = append(errorCrawlResults, ErrorCrawlResult{
					URL:   url,
					Error: err.Error(),
				})
				return nil
			}

			// Handle success result
			s.logger.Info().Str("url", url).Msg("Successfully extracted data from URL")
			successCrawlResult := SuccessCrawlResult{
				URL:              result.URL,
				Title:            result.Title,
				MetaDescriptions: result.MetaDescriptions,
				Links:            result.Links,
				KeywordCounts:    result.KeywordCounts,
			}
			successCrawlResults = append(successCrawlResults, successCrawlResult)
			return nil
		})
	}

	// Wait for all requests to finish
	err := eg.Wait()
	if err != nil {
		// This is actually not gonna happen as we are not returning any error from the goroutines
		// But handling the error just in case
		return nil, nil, err
	}

	// Return both success and error results
	return successCrawlResults, errorCrawlResults, nil
}
