GOLINT := golangci-lint

PACKAGES_FOR_TEST := $(shell go list ./... | grep -v model | grep -v "mock")


all: dep gen-mock lint vet test

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

test:
	@go test -cover -race -tags=unit -parallel 10 -count=1 -v $(PACKAGES_FOR_TEST)

run:
	go run -race main.go --node.host=$(NODE_HOST) --node.port=$(NODE_PORT)

vet:
	go vet ./...

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.57.2

lint: dep check-lint ## Lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml

check-mockgen:
	@which mockgen || go install go.uber.org/mock/mockgen@latest

gen-mock:
	mockgen -package mock -source core/interface.go -destination core/mock/interface.go
	mockgen -package mock -source client/interface.go -destination client/mock/interface.go
