# domain-crawler

## Description

This is a domain crawler service which accepts a list of url's to crawl and keywords to search for.

## OpenAPI Specification

The OpenAPI specification can be found in the `api/openapi/api.yaml` file.

## Running Locally

### Local Go Installation

```
make start-local
```

### Docker & Docker Compose

```
make start-local-docker
```

You can update the environment variables injected in docker-compose.yml file.

## Environment Variables

```
PORT - The port the server will listen on
EXTRACTOR_CONCURRENT_LIMIT - The number of concurrent requests the extractor will make
RATE_LIMIT_RPM - The rate limit configured for the service
```

## Rate Limiting

The service is rate limited to 60 requests per minute (default).
This can be configured using the `RATE_LIMIT_RPM` environment variable.

## Cache

The service uses an in-memory cache to store the results of the data extraction.

## CI

The CI pipeline is configured using GitHub Actions.

This runs the tests and linters on every push to the repository.

## Tests

```
make tests
```

## Linting

```
brew install golangci-lint
brew upgrade golangci-lint

make lint
```
