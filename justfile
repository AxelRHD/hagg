app_name := "hagg"
main_file := "cmd/main.go"
bin_dir := "bin"
bin_file := bin_dir / app_name

# Git version: tag, tag-hash, or just hash (with -dirty suffix if uncommitted changes)
git_version := `v=""; \
    tag=$(git describe --tags --exact-match 2>/dev/null); \
    if [ -n "$tag" ]; then \
        v="$tag"; \
    else \
        latest=$(git describe --tags --abbrev=0 2>/dev/null); \
        hash=$(git rev-parse --short=7 HEAD); \
        if [ -n "$latest" ]; then \
            v="$latest-$hash"; \
        else \
            v="$hash"; \
        fi; \
    fi; \
    git diff --quiet 2>/dev/null || v="$v-dirty"; \
    echo "$v"`

[private]
default:
    @just --list --unsorted

# --- Development ---

# Run with live reload (air)
[group('dev')]
dev:
    air

# Run without live reload
[group('dev')]
run:
    go run {{main_file}} serve

# Show active configuration
[group('dev')]
config:
    go run {{main_file}} config

# --- Build ---

# Private build helper with version embedding
[private]
_build env="":
    mkdir -p {{bin_dir}}
    {{env}} go build \
        -ldflags "-X 'github.com/axelrhd/hagg/internal/version.Version={{git_version}}'" \
        -o {{bin_file}} {{main_file}}

# Build binary (local)
[group('build')]
build: _build

# Build Linux binary (amd64, static)
[group('build')]
build-linux: (_build "CGO_ENABLED=0 GOOS=linux GOARCH=amd64")

# --- Database ---

# Run database migrations
[group('db')]
migrate-up:
    goose -dir migrations sqlite db.sqlite3 up

# Rollback last migration
[group('db')]
migrate-down:
    goose -dir migrations sqlite db.sqlite3 down

# Create a new migration
[group('db')]
migrate-create name:
    goose -dir migrations create {{name}} sql

# --- Quality ---

# Run tests
[group('quality')]
test:
    go test ./...

# Format code
[group('quality')]
fmt:
    go fmt ./...

# Lint code
[group('quality')]
lint:
    golangci-lint run

# --- Cleanup ---

# Clean build artifacts
[group('cleanup')]
clean:
    rm -rf {{bin_dir}}

# Tidy go modules
[group('cleanup')]
tidy:
    go mod tidy
