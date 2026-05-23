.PHONY: build run test lint docker

build:
	go build -o bin/monitor ./cmd/monitor

run:
	go run ./cmd/monitor

test:
	go test -v ./...

lint:
	golangci-lint run

docker:
	docker build -t go-uptime-monitor .

docker-run:
	docker compose up -d
