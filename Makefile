VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
CONFIG := ./volume/

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint build test-unit

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/forbole/juno/v3/cmd.Version=$(VERSION) \
	-X github.com/forbole/juno/v3/cmd.Commit=$(COMMIT)
BUILD_FLAGS :=  -ldflags '$(LD_FLAGS)'


###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building bdjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bdjuno.exe ./cmd/bdjuno
else
	@echo "building bdjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bdjuno ./cmd/bdjuno
endif
.PHONY: build

docker-build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building bdjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bdjuno.exe ./cmd/bdjuno
else
	@echo "building bdjuno binary..."
	@LEDGER_ENABLED=false CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --tags muslc -mod=vendor $(BUILD_FLAGS) -o build/bdjuno ./cmd/bdjuno
endif
.PHONY: docker-build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing bdjuno binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/bdjuno
.PHONY: install

###############################################################################
###                                   Run                                   ###
###############################################################################

run:
ifeq ($(OS),Windows_NT)
	@echo "running bdjuno for Windows..."
	@go run cmd/bdjuno/main.go parse --home $(CONFIG)
else
	@echo "running bdjuno for Linux..."
	@go run cmd/bdjuno/main.go parse --home $(CONFIG)
endif
.PHONY: run

run-build: build
ifeq ($(OS),Windows_NT)
	@echo "running bdjuno for Windows..."
	@./build/bdjuno.exe parse --home $(CONFIG)
else
	@echo "running bdjuno for Linux..."
	@./build/bdjuno parse --home $(CONFIG)
endif
.PHONY: run-build

start:
ifeq ($(OS),Windows_NT)
	@echo "running bdjuno for Windows..."
	@./build/bdjuno.exe parse --home $(CONFIG)
else
	@echo "running bdjuno for Linux..."
	@./build/bdjuno parse --home $(CONFIG)
endif
.PHONY: start

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop bdjuno-test-db || true && docker rm bdjuno-test-db || true
.PHONY: stop-docker-test

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name bdjuno-test-db -e POSTGRES_USER=bdjuno -e POSTGRES_PASSWORD=password -e POSTGRES_DB=bdjuno -d -p 6433:5432 postgres
.PHONY: start-docker-test

test-unit: start-docker-test
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -local github.com/forbole/bdjuno
.PHONY: format

clean:
	rm -f tools-stamp ./build/**
.PHONY: clean
