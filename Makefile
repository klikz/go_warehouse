.PHONY: build
build:
	go build -v ./cmd/warehouse
.PHONY: run
run:
	go run -v ./cmd/warehouse
.DEFAULT_GOAL := run