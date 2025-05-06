.PHONY: build unit e2e test

build:
	@mkdir -p bin
	go mod tidy
	go build -o bin/biathlon ./cmd/biathlon

unit:
	go test ./internal/...

cov:
	touch cover.out
	go tool cover -func=cover.out

e2e:
	cd scripts && pytest

test: build unit e2e
