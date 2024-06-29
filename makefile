# ==============================================================================
# Running

dev-up:
	go run cmd/parser/main.go logs/qgames.log

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
