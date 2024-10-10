APP_NAME:= coldfront-adserver
VERSION := $(shell git describe --tags --always)

all: build

build:
	@mkdir -p bin/
	@go build -o bin/$(APP_NAME) -ldflags="-X 'main.version=$(VERSION)'" cmd/$(APP_NAME)/main.go

clean: 
	@rm -rf bin/

run: build
	@bin/($APP_NAME)

test:
	@go test ./...

.PHONY: all build clean run test
