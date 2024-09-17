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

The service uses an in-memory cache to store the HTML document response. The cache is not persisted and is cleared on every restart.

## CI

The CI pipeline is configured using GitHub Actions.

This runs the tests and linting on every push to the repository.

## Tests

```
make test
```

## Linting

```
brew install golangci-lint
brew upgrade golangci-lint

make lint
```

## Assumptions

1. Cache key - I'm caching the entire HTML document returned by the URL. I didn't cache the `ExtractResult` because different keywords can be used for the same URL. The only caveat here is that the cache footprint can go big since we're caching the entire HTML document.
2. Rate limiting - I'm using a simple rate limiter which resets every minute.
