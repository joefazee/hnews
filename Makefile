run:
	@go run ./cmd/web


run/migrate:
	@go run ./cmd/web -migrate=true
