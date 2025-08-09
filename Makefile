run:
	go run cmd/main.go

test:
	go test ./... -v

build:
	go build -o go_v2ray_client cmd/main.go

clean:
	rm -f go_v2ray_client