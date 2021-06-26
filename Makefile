.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

lint:
	