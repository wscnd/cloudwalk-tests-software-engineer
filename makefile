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

run-tool-trace:
	go tool trace trace.out

trace: run-trace run-tool-trace

# ==============================================================================
# Cleanup

.PHONY: clean

clean:
	rm -f main match_data.json
