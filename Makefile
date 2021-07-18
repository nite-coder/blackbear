.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

lint:
	docker run --rm -v ${LOCAL_WORKSPACE_FOLDER}:/app -w /app golangci/golangci-lint:v1.41-alpine golangci-lint run ./... -v