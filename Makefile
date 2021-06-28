.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

.PHONY: lint
lint:
	golangci-lint run ./... -v