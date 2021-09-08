#!/bin/bash
.PHONY: run
run:
	@go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: run-db
run-db:
	@docker run \
	--name=riak \
	-d \
	--rm \
	-p 8087:8087 \
	basho/riak-kv

.PHONY: stop-db
stop-db:
	@docker stop riak