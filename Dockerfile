# build stage
FROM golang:1.23.1-alpine AS build

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY api api
COPY cmd cmd
COPY internal internal

RUN go build -ldflags="-s -w" cmd/domaincrawler/main.go

# app stage
FROM alpine:3.18.4 AS app

WORKDIR /app

COPY --from=build /app/main main

CMD [ "./main" ]
