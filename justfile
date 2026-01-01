# Justfile for HAGG development

# Run with live reload (air)
dev:
    air

# Build the project
build:
    go build -o hagg cmd/main.go

# Run without live reload
run:
    go run cmd/main.go serve

# Run database migrations
migrate-up:
    goose -dir migrations sqlite db.sqlite3 up

# Create a new migration
migrate-create name:
    goose -dir migrations create {{name}} sql

# Run tests
test:
    go test ./...

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run

# Clean build artifacts
clean:
    rm -f hagg hagg_test
    rm -f db.sqlite3

# Show this help
help:
    @just --list
