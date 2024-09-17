package openapi

import "embed"

//go:embed *.yml
var FS embed.FS

const OpenAPISpecFilename = "api.yml"
