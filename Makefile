APP_NAME := goval
AIR_BIN := $(shell go env GOPATH)/bin/air
# Default target
all: test build

# Run tests
test:
	go test -v ./...

# Build the project
build:
	mkdir -p $(APP_NAME)
	go build -o ./$(APP_NAME)/$(APP_NAME)

# Clean the build
clean:
	go clean
	rm -f $(APP_NAME)
		
# Run the application
run: build
	./$(APP_NAME)/$(APP_NAME)

dev:
	CONFIG_DIR=./goval/config DATA_DIR=./goval/data $(AIR_BIN)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format the code
fmt:
	go fmt ./...

.PHONY: