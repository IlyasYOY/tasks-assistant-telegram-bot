.DEFAULT_GOAL = verify

.PHONY: generate
generate: 
	go generate ./...

.PHONY: format
format: 
	go tool goimports -w .

.PHONY: lint 
lint:
	go tool golangci-lint run --config=.golangci.yml ./...

.PHONY: test 
test: 
	go tool gotestsum ./... -race

.PHONY: build 
build: 
	go build ./...

.PHONY: test-watch 
test-watch:
	go tool gotestsum --watch -- ./...

.PHONY: verify 
verify: lint test

