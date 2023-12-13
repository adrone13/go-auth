# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/auth/main.go

# Run the application
run:
	@echo "Starting app..."
	@go run cmd/auth/main.go

# Run in watch mode
watch:
	@echo "Starting in watch mode..."
	@air

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run watch test clean
