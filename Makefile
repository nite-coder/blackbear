.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

lint:
	golangci-lint run ./... -v

docker_lint:
	docker run -it --rm -v "${LOCAL_WORKSPACE_FOLDER}:/app" -w /app golangci/golangci-lint:v1.59.1-alpine golangci-lint run ./... -v

release: lint test

