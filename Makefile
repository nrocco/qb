.DEFAULT_GOAL := build

.PHONY: build
build: lint test

.PHONY: lint
lint:
	git ls-files | xargs misspell -error
	golint -set_exit_status ./...
	go vet -v ./...
	errcheck -blank -asserts ./...

.PHONY: test
test:
	go test -v -short ./...

.PHONY: coverage
coverage:
	mkdir -p coverage
	go test -covermode=count -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
