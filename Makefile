.DEFAULT_GOAL := build

.PHONY: build
build: lint test

.PHONY: lint
lint:
	golint -set_exit_status ./...
	go vet -v ./...

.PHONY: test
test:
	go test -v -short ./...

.PHONY: coverage
coverage:
	mkdir -p coverage
	go test -covermode=count -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
