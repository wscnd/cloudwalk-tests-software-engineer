# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/stretchr/testify
	go install honnef.co/go/tools/cmd/staticcheck@latest

# ==============================================================================
# Running

dev-init: dev-gotooling

dev-run: build dev-run-build

dev-up:
	go run cmd/parser/main.go logs/qgames.log

dev-run-build:
	./main logs/qgames.log

view-output:
	cat match_data.json | jq .

# ==============================================================================
# Linting and tests

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

test:
	CGO_ENABLED=0 go test -count=1 ./...

# ==============================================================================
# Build

build:
	go build cmd/parser/main.go

# ==============================================================================
# Tracing

run-trace:
	./main logs/qgames.log > trace.out

tool-trace:
	go tool trace trace.out

trace: run-trace tool-trace

# ==============================================================================
# Cleanup

.PHONY: clean

clean:
	rm -f main match_data.json trace.out

# ==============================================================================
# Container

BASE_IMAGE_NAME := localhost
APP := quake-parser
APP_VERSION := 0.1
IMAGE_NAME := $(BASE_IMAGE_NAME)/$(APP):$(APP_VERSION)
CONTAINER_NAME := $(APP)-container

docker-build:
	docker build \
		-f zarf/docker/dockerfile \
		-t $(IMAGE_NAME) \
		--build-arg BUILD_REF=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

docker-run:
	docker run \
	--name $(CONTAINER_NAME) \
	-d $(IMAGE_NAME)

docker-copy:
	docker cp \
	 $(CONTAINER_NAME):/parser/match_data.json \
	 match_data.json

docker-rm:
	docker rm $(CONTAINER_NAME)

docker-all: docker-build docker-run docker-copy docker-rm