.PHONY: build
build:
	go build -v ./cmd/coinmonserver

.DEFAULT_GOAL := build
