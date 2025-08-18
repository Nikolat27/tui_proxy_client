# Build the application
build:
	go build -o tui_proxy_client cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run all tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f tui_proxy_client
	rm -f coverage.out
	rm -f coverage.html
	rm -f test.json

# Install dependencies
deps:
	go mod download
	go mod tidy
