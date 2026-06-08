SHELL:=/usr/bin/env bash

dev: build
	./pac

build:
	go build .

.PHONY: proto
proto:
	protoc --go_out=core proto/*.proto
