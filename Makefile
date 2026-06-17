.PHONY: build test lint run docker-build

build:
	go build -o bin/shipping-service ./cmd/server

test:
	go test ./... -v

run:
	go run ./cmd/server

docker-build:
	docker build -t ecom-shipping-service .
