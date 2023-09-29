BIN := "./bin/smdaemon"
BIN_CLI := "./bin/sm-client"
DOCKER_IMG="sm-service:develop"
DOCKER_IMG_CLI="sm-client:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/service

run: build
	$(BIN) --config ./configs/config.yaml

build-client:
	go build -v -o $(BIN_CLI) -ldflags "$(LDFLAGS)" ./cmd/client

run-client: build-client
	$(BIN_CLI)

build-img-service:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/service.Dockerfile .

build-img-client:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_CLI) \
		-f build/client.Dockerfile .

run-img-service: build-img-service
	docker run -p 8086:8086 $(DOCKER_IMG)

run-img-client: build-img-client
	docker run --network host ${DOCKER_IMG_CLI}

version: build
	$(BIN) --version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run --timeout=90s ./...

test-integration:
	go test -race ./tests/... -count 3

test: test-integration
	go test -race ./internal/... -count 100

generate: 
	go generate ./...

.PHONY: build run build-img-service run-img-service version test lint
