run:
	go run cmd/main.go

test:
	go test ./... -v

build:
	go build -o tui_proxy_client cmd/main.go

clean:
	rm -f tui_proxy_client