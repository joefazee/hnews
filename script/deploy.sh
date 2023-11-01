

SERVER_IP=146.190.28.127
SERVER_PORT=22
SERVER_USER=hnews

@echo "deploying..."
rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -rP --delete ./migrations  ${SERVER_USER}@${SERVER_IP}:~
rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./bin/web  ${SERVER_USER}@${SERVER_IP}:~
rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./script/web.service  ${SERVER_USER}@${SERVER_IP}:~
rsync -e "ssh -p ${SERVER_PORT} -i ./keys/hnews" -P ./script/Caddyfile  ${SERVER_USER}@${SERVER_IP}:~

# services
sudo mv web.service /etc/systemd/system/
sudo systemctl enable web
sudo systemctl restart web
sudo mv Caddyfile /etc/caddy/
sudo systemctl reload caddy
