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

db-run:
	@echo "Starting db..."
	@docker compose up

db-stop:
	@echo "Stopping db..."
	@docker compose down

migrate-create:
	@echo "Creating migration..."
	@go run cmd/migrate/migrate.go create $(name)

migrate-up:
	@echo "Migrating db up..."
	@go run cmd/migrate/migrate.go up

migrate-down:
	@echo "Migrating db down..."
	@go run cmd/migrate/migrate.go down

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run watch db-run db-stop test clean
