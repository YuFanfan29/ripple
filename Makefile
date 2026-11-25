# Common development tasks. Run `make` or `make help` to list them.

GO ?= go

.DEFAULT_GOAL := help
.PHONY: help test lint fmt ci

help: ## List available targets
	@grep -hE '^[a-z][a-z-]*:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-6s %s\n", $$1, $$2}'

test: ## Run all tests with the race detector
	$(GO) test -race ./...

lint: ## Check formatting and run static analysis
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "These files need formatting (run make fmt):"; echo "$$unformatted"; exit 1; fi
	$(GO) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; \
	else echo "golangci-lint is not installed, so it was skipped. CI still runs it."; fi

fmt: ## Format all Go code
	gofmt -w .

ci: lint test ## Run the same checks as CI
