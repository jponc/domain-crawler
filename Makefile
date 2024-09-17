.PHONY: start-local start-local-docker run-tests

install:
	go mod download


start-local:
	PORT=8080 go run cmd/domaincrawler/main.go

start-local-docker:
	# This runs docker compose with the local docker-compose.yml file
	# This always rebuilds the docker image and runs the container
	docker compose up --build

run-tests:
	go test ./... -v
