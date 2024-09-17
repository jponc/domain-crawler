package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jponc/domain-crawler/internal/errs"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

// OpenAPIValidatorMiddleware validates requests against the OpenAPI spec
func OpenAPIValidatorMiddleware(doc *openapi3.T) func(http.Handler) http.Handler {
	return nethttpmiddleware.OapiRequestValidatorWithOptions(
		doc,
		&nethttpmiddleware.Options{
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)

				response, _ := json.Marshal(&errs.ErrorResponse{Error: message})

				// Write the Response
				w.Write(response)
			},
		},
	)
}
