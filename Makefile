BINARY_NAME := evcc-cli
OPENAPI_REF ?= master
GOBIN_PATH := $(shell go env GOPATH)/bin
OAPI_CODEGEN := $(GOBIN_PATH)/oapi-codegen

.PHONY: tools generate tidy test build run

tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

generate: tools
	mkdir -p internal/gen/evcc internal/gen/evccstate
	$(OAPI_CODEGEN) -config api/oapi-codegen.yaml api/upstream/openapi.yaml
	$(OAPI_CODEGEN) -config api/oapi-codegen.state.yaml api/upstream/openapi.state.yaml

tidy:
	go mod tidy

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) .

run:
	go run .
