run:
	@go run ./cmd/web


run/migrate:
	@go run ./cmd/web -migrate=true

build:
	@go build -o=./bin/web ./cmd/web
