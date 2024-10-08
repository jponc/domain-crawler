package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/jponc/domain-crawler/api/openapi"
	"github.com/jponc/domain-crawler/internal/cache"
	"github.com/jponc/domain-crawler/internal/config"
	"github.com/jponc/domain-crawler/internal/crawl/handlers"
	"github.com/jponc/domain-crawler/internal/crawl/services"
	"github.com/jponc/domain-crawler/internal/extractor"
	"github.com/jponc/domain-crawler/internal/middlewares"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("starting domaincrawler..")

	// Load config
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Load OpenAPI spec
	openapiSpec, err := openapi.FS.ReadFile(openapi.OpenAPISpecFilename)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read openapi spec")
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(openapiSpec)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load openapi spec")
	}

	err = doc.Validate(loader.Context)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to validate openapi spec")
	}

	oapiValidatorMiddleware := middlewares.OpenAPIValidatorMiddleware(doc)

	// Setup chi router and middlewares
	r := chi.NewRouter()
	r.Use(oapiValidatorMiddleware)
	r.Use(httprate.LimitByIP(config.RateLimitRPM, time.Minute))

	// Setup dependencies
	httpClient := &http.Client{}
	inmemoryCache := cache.NewInMemoryCache()
	extractorClient := extractor.NewExtractorClient(httpClient, inmemoryCache)
	crawlService := services.NewCrawlService(extractorClient, config.ExtractorConcurrentLimit)

	// Setup handlers
	crawlHandler := handlers.NewCrawlHandler(crawlService)

	// Setup routes
	r.Post("/crawl", crawlHandler.Crawl)

	// Start server
	addr := fmt.Sprintf(":%s", config.Port)
	log.Info().Msgf("listening on %s", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
