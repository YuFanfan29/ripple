# Common development tasks. Run `make` or `make help` to list them.

GO ?= go
MVNW := cd services/content && ./mvnw -B -ntp -q

.DEFAULT_GOAL := help
.PHONY: help test test-go test-java lint lint-go lint-java fmt fmt-go fmt-java ci

help: ## List available targets
	@grep -hE '^[a-z][a-z-]*:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-10s %s\n", $$1, $$2}'

test: test-go test-java ## Run all tests

test-go: ## Run Go tests with the race detector
	$(GO) test -race ./...

test-java: ## Run Java tests
	$(MVNW) test

lint: lint-go lint-java ## Check formatting and run static analysis

lint-go: ## Check Go formatting, then run go vet and golangci-lint
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "These files need formatting (run make fmt):"; echo "$$unformatted"; exit 1; fi
	$(GO) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; \
	else echo "golangci-lint is not installed, so it was skipped. CI still runs it."; fi

lint-java: ## Check Java formatting
	$(MVNW) spotless:check

fmt: fmt-go fmt-java ## Format all code

fmt-go: ## Format Go code
	gofmt -w .

fmt-java: ## Format Java code
	$(MVNW) spotless:apply

ci: lint test ## Run the same checks as CI
