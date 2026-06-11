SHELL:=/usr/bin/env bash

dev: build
	./pac

build:
	CGO_ENABLED=0 go build .

.PHONY: proto
proto:
	protoc --go_out=core proto/*.proto
