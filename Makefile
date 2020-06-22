VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

export GO111MODULE = on

all: ci-lint ci-test install

###############################################################################
# Build / Install
###############################################################################

LD_FLAGS = -X github.com/desmos-labs/juno/version.Version=$(VERSION) \
	-X github.com/desmos-labs/juno/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building bdjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bdjuno.exe ./cmd/bdjuno
else
	@echo "building bdjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bdjuno ./cmd/bdjuno
endif

install: go.sum
	@echo "installing bdjuno binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/bdjuno

###############################################################################
# Tests / CI
###############################################################################

coverage:
	@echo "Viewing test coverage..."
	@go tool cover --html=coverage.out

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop bdjuno-test-db || true && docker rm bdjuno-test-db || true

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name bdjuno-test-db -e POSTGRES_USER=bdjuno -e POSTGRES_PASSWORD=password -e POSTGRES_DB=bdjuno -d -p 5433:5432 postgres

ci-test: start-docker-test
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...

ci-lint:
	@echo "Running GolangCI-Lint..."
	@GO111MODULE=on golangci-lint run
	@echo "Formatting..."
	@find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	@echo "Verifying modules..."
	@go mod verify

clean:
	rm -f tools-stamp ./build/**

.PHONY: install build ci-test ci-lint coverage clean start-docker-test
