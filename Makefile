# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/auth/main.go

# Run the application
run:
	@echo "Starting app..."
	@go run cmd/auth/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run test clean
