# Use to automate the common task while building chatter.
BINARY_NAME := chatter
BIN_DIR := release
VERSION ?= vlatest
PLATFORMS := windows
os := windows

build:
	go build -o $(BINARY_NAME) -race -v

clean:
	go clean -i -cache -testcache
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

.PHONY: vendor
vendor:
	dep ensure -update 

