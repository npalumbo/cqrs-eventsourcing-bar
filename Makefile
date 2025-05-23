#@ Helpers
# from https://www.thapaliya.com/en/writings/well-documented-makefiles/
help:  ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tools
tools: ## Installs required binaries locally
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@v2.49.0
	go install fyne.io/fyne/v2/cmd/fyne@v2.6.0-alpha1

build: check## Builds cqrs-eventsourcing-bar go binaries for local arch. Outputs to `bin/app, bin/readservice, bin/writeservice`
	@echo "== build"
	CGO_ENABLED=1 go build -o bin/ ./...

##@ Cleanup
clean: ## Deletes binaries from the bin folder
	@echo "== clean"
	rm -rfv ./bin

##@ Tests
test: ## Run unit tests
	@echo "== unit test"
	go test ./...

##@ Run static checks
check: ## Runs lint, fmt and vet checks against the codebase
	golangci-lint --timeout 280s run
	go fmt ./...
	go vet ./...

##@ Golang Generate
generate: ## Calls golang generate
	go mod tidy
	go generate ./...
