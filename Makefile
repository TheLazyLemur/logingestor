gen:
	templ generate ./...

run: gen
	go run ./cmd/logingestor/...

build: gen
	go build -o bin/logingestor ./cmd/logingestor/...
