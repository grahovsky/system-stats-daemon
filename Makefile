BIN := "./bin/smdaemon"
BIN_CLI := "./bin/sm-client"
DOCKER_IMG := "smd:develop"

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

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img-service: build-img
	docker run -p 8086:8086 $(DOCKER_IMG) 

run-img-client: build-img
	docker run -it --network host smd:develop sh -c "\${BIN_FILE_CLIENT}"

version: build
	$(BIN) --version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run --timeout=90s ./...

test-grpc:
	go test -race ./tests/grpc/... -count 3

test: test-grpc
	go test -race ./internal/... -count 100

test-integration:
	go test -tags integration ./tests/integration/... -count 1 -v

integration-test:
	docker compose -f deployments/integration-test-compose.yaml up --build --attach tester --exit-code-from tester
	docker compose -f deployments/integration-test-compose.yaml down

generate: 
	go generate ./...

up:
	docker compose -f deployments/compose.yaml --project-directory deployments up -d --build

rebuild: docker compose -f deployments/compose.yaml --project-directory deployments build --no-cache

logs:
	docker compose -f deployments/compose.yaml logs -f

down:
	docker compose -f deployments/compose.yaml down

clean:
	docker compose -f deployments/compose.yaml down --rmi all -v
	docker compose -f deployments/integration-test-compose.yaml down --rmi all -v


.PHONY: build run build-img-service run-img-service version test lint
