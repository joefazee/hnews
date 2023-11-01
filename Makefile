run:
	@go run ./cmd/web


run/migrate:
	@go run ./cmd/web -migrate=true

build:
	@echo "building binary..."
	@go build -o=./bin/web ./cmd/web
	@GOOS=linux GOARCH=amd64 go build -o=./bin/linux/web ./cmd/web

ssh-root:
	@ssh -i ./keys/hnews root@146.190.28.127

ssh-hnews:
	@ssh -i ./keys/hnews hnews@146.190.28.127

SERVER_IP=167.71.76.148
SERVER_PORT=22
SERVER_USER=hnews

deploy: build
	@echo "deploying..."
	rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -rP --delete ./migrations  ${SERVER_USER}@${SERVER_IP}:~
	rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./bin/linux/web  ${SERVER_USER}@${SERVER_IP}:~
	rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./script/web.service  ${SERVER_USER}@${SERVER_IP}:~
	rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./script/Caddyfile  ${SERVER_USER}@${SERVER_IP}:~

