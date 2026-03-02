.PHONY: all build run test cover lint integration-test load-test migrate docker-up docker-down tidy

BINARY   := bin/api
CMD_PATH := ./cmd/api
COVER_MIN := 70

# Unit-testable packages (excludes DB-dependent infrastructure adapters,
# which are covered exclusively by integration tests).
UNIT_PKGS := $(shell go list ./internal/... | grep -v '/persistence/')

all: lint test build

# --- Build ---
build:
	@echo "==> Building..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BINARY) $(CMD_PATH)

run:
	@go run $(CMD_PATH)

# --- Dependencies ---
tidy:
	go mod tidy

# --- Tests ---
test:
	@echo "==> Running unit tests..."
	go test $(UNIT_PKGS) -count=1 -timeout=60s

cover:
	@echo "==> Checking unit-test coverage (min $(COVER_MIN)%, excludes DB-dependent persistence)..."
	go test $(UNIT_PKGS) -count=1 \
		-coverprofile=coverage.out \
		-covermode=atomic \
		-coverpkg=$(shell go list ./internal/... | grep -v '/persistence/' | tr '\n' ',')
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $${COVERAGE}%"; \
	if [ "$$(echo "$${COVERAGE} < $(COVER_MIN)" | bc -l)" = "1" ]; then \
		echo "Coverage $${COVERAGE}% is below $(COVER_MIN)% threshold"; exit 1; \
	fi

cover-html:
	go test $(UNIT_PKGS) -coverprofile=coverage.out -covermode=atomic \
		-coverpkg=$(shell go list ./internal/... | grep -v '/persistence/' | tr '\n' ',')
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> Coverage report: coverage.html"

integration-test:
	@echo "==> Running integration tests..."
	go test ./tests/integration/... -tags integration -count=1 -timeout=120s

load-test:
	@echo "==> Running load tests (requires k6)..."
	k6 run tests/load/transaction_load_test.js

# --- Lint ---
lint:
	@echo "==> Linting..."
	golangci-lint run --timeout=5m

lint-fix:
	golangci-lint run --fix

# --- Database ---
migrate:
	@echo "==> Applying migrations..."
	PGPASSWORD=$${POSTGRES_PASSWORD:-postgres} psql \
		-h $${POSTGRES_HOST:-localhost} \
		-U $${POSTGRES_USER:-postgres} \
		-d $${POSTGRES_DB:-antifraude} \
		-f scripts/migrate.sql

# --- Docker ---
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

docker-clean:
	docker compose down -v --remove-orphans
