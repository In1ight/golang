.PHONY: build test run lint


build:
	go build -o bin/control-plane ./cmd/controlplane


test:
	go test ./...

run:
	go run ./cmd/controlplane

lint:
	golangci-lint run ./...