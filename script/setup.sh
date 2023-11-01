export LC_ALL=en_US.UTF-8

# setup timezone and locale
timedatectl set-timezone Africa/Lagos
apt --yes install locales-all

add-apt-repository --yes universe

# add new user
useradd --create-home --shell "/bin/bash" --groups sudo hnews

# copy root ssh keys to hnews user
rsync --archive --chown=hnews:hnews /root/.ssh /home/hnews

# enable firewall
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# install PostgreSQL
apt --yes install postgresql
sudo -i -u postgres psql -c "CREATE DATABASE hnews"
sudo -i -u postgres psql -d hnews -c "CREATE ROLE hnews WITH LOGIN PASSWORD 'KJS93sds91A9129S2as'"

echo "HNEWS_DB_DSN='postgres://hnews:KJS93sds91A9129S2as@localhost/hnews?sslmode=disable'" >> /etc/environment

# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy


echo "Script complete! Rebooting..."
reboot


