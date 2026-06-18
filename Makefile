SHELL:=/usr/bin/env bash

dev: build
	./pac -config-file=data/config.json

build:
	CGO_ENABLED=0 go build .

.PHONY: proto
proto:
	protoc --go_out=core proto/*.proto
