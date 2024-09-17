package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jponc/domain-crawler/internal/crawl/services"
	"github.com/jponc/domain-crawler/internal/errs"
	"github.com/jponc/domain-crawler/internal/utils"
)

type crawlService interface {
	Crawl(ctx context.Context, urls []string, keywords []string) ([]services.SuccessCrawlResult, []services.ErrorCrawlResult, error)
}

type crawlHandler struct {
	crawlService crawlService
}

func NewCrawlHandler(crawlService crawlService) *crawlHandler {
	h := &crawlHandler{
		crawlService: crawlService,
	}

	return h
}

func (h *crawlHandler) Crawl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var reqBody CrawlRequest

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := errs.ErrorResponse{Error: "failed to decode request body"}
		json.NewEncoder(w).Encode(errResp)
		return
	}

	// Remove duplicate urls if any
	uniqueURLs := utils.RemoveDuplicates(reqBody.URLs)

	// Crawl the URLs
	successCrawlResults, errorCrawlResults, err := h.crawlService.Crawl(ctx, uniqueURLs, reqBody.Keywords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := errs.ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(errResp)
		return
	}

	// Convert success crawl results to success results
	successResults := convertSuccessCrawlResultsToSuccessResults(successCrawlResults)

	// Convert error crawl results to error Results
	errorResults := convertErrorCrawlResultsToErrorResults(errorCrawlResults)

	// Create response Body
	respBody := CrawlResponse{
		Results: successResults,
		Errors:  errorResults,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}
