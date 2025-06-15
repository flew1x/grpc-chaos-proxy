BINARY=proxy
BUILD_DIR=bin
CMD_PATH=./cmd/proxy

.PHONY: all build run clean docker docker-build docker-run lint test

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)

run: build
	./$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR) $(BINARY)

lint:
	golangci-lint run ./...

test:
	go test ./...

mod:
	go mod tidy

