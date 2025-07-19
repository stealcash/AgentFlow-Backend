# Makefile for AgentFlow Backend

# The binary name you want
APP_NAME=AgentFlow

# Where to put built binaries
BUILD_DIR=build

# Phony targets (not real files)
.PHONY: build run clean

# Build the Go binary
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

# Run the built binary
run: build
	./$(BUILD_DIR)/$(APP_NAME)

# Remove the build output
clean:
	rm -rf $(BUILD_DIR)
