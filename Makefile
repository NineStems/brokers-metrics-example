BIN := "./bin/master"
RABBIT_READER := "./bin/rabbit-reader"
RABBIT_WRITER := "./bin/rabbit-writer"
KAFKA_READER := "./bin/kafka-reader"
KAFKA_WRITER := "./bin/kafka-writer"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/master
	go build -v -o $(RABBIT_READER) -ldflags "$(LDFLAGS)" ./cmd/rabbit-reader
	go build -v -o $(RABBIT_WRITER) -ldflags "$(LDFLAGS)" ./cmd/rabbit-writer
	go build -v -o $(KAFKA_READER) -ldflags "$(LDFLAGS)" ./cmd/kafka-reader
	go build -v -o $(KAFKA_WRITER) -ldflags "$(LDFLAGS)" ./cmd/kafka-writer

rabbit-r:
	$(RABBIT_READER)

rabbit-w:
	$(RABBIT_WRITER)

kafka-r:
	$(KAFKA_READER)

kafka-w:
	$(KAFKA_WRITER)

master:
	$(BIN)

up:
	docker-compose -f docker-compose.yaml up -d --force-recreate

down:
	docker-compose -f docker-compose.yaml down

deps:
	go mod download

.PHONY: build run