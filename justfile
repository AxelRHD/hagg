# HAGG justfile

# Default recipe (list all recipes)
default:
    @just --list

# CSS Build (production, minified)
css-build:
    tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --minify

# CSS Watch (development, auto-rebuild)
css-watch:
    tailwindcss -i ./static/css/base.css -o ./static/css/styles.css --watch

# Run the server
run:
    go run cmd/main.go

# Build the binary
build:
    go build -o hagg cmd/main.go

# Run tests
test:
    go test ./...

# Clean generated files
clean:
    rm -f hagg
    rm -f static/css/styles.css
