.PHONY: build unit e2e test

build:
	@mkdir -p bin
	go build -o bin/biathlon ./cmd/biathlon

unit:
	go test ./internal/...

e2e:
	cd scripts && pytest

test: build unit e2e
