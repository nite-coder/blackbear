.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

lint:
	golangci-lint run ./... -v

lint.docker:
	docker run --rm -v ${pwd}:/app -w /app golangci/golangci-lint:v1.41-alpine golangci-lint run 