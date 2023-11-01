# Yet another Hackernews clone

## How to run it locally
1. Run `docker compose up -d` to start Postgres on port `5455`
2. Run `go run ./cmd/web -migrate=true` then stop the application! This is just one time to create the tables
3. Run `go run ./cmd/web` to start the application on port 8009

For all command line options, type `go run ./cmd/web --help`

# Deployment Video
Want to deploy it to the web? Watch the video here.

