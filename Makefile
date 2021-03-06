clean: ## clean build output
	rm -rf bin/*

go-clean: ## Invoke gofmt's "simplify" option to streamline the source code.
	gofmt -w -s ./pkg/
	gofmt -w -s ./pkg/api/
	gofmt -w -s ./pkg/api/client/
	gofmt -w -s ./pkg/api/util/
	gofmt -w -s ./pkg/api/v2/
	gofmt -w -s ./cmd/
	goimports -w $(shell git ls-files "**/*.go" "*.go" | grep -v -e "vendor" | xargs echo)

.PHONY: install-tools
install-tools: ## install tools needed by go-link-checks
	GOIMPORTS_CMD=$(shell command -v goimports 2> /dev/null)
ifndef GOIMPORTS_CMD
	go get golang.org/x/tools/cmd/goimports
endif

	GOLINT_CMD=$(shell command -v golint 2> /dev/null)
ifndef GOLINT_CMD
	go get github.com/golang/lint/golint
endif

	GOCYCLO_CMD=$(shell command -v gocyclo 2> /dev/null)
ifndef GOLINT_CMD
	go get github.com/fzipp/gocyclo
endif

test: ## run go test (must be maintained)
	go test -v ./pkg/api/client
	go test -v ./pkg/api/util
	go test -v ./pkg/api/v2

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
