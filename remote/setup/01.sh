#!/bin/bash
set -eu
# ==================================================================================== #
# VARIABLES
# ==================================================================================== #
TIMEZONE=Asia/Singapore
USERNAME=bms
read -p "Enter password for bms DB user:" DB_PASSWORD
export LC_ALL=en_US.UTF-8

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #
add-apt-repository --yes universe
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 1883/tcp
ufw --force enable

apt --yes install fail2ban

curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

apt --yes install mysql-server

sudo mysql -e "CREATE DATABASE bms;"
sudo mysql -e "CREATE USER 'bms'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';"
sudo mysql -e "GRANT ALL PRIVILEGES ON bms.* TO 'bms'@'localhost';"
sudo mysql -e "FLUSH PRIVILEGES;"

echo "BMS_DB_DSN='mysql://root:${DB_PASSWORD}@localhost/bms'" >> /etc/environment

# Install and configure MQTT server
apt --yes install mosquitto mosquitto-clients

# Create password file for MQTT
sudo touch /etc/mosquitto/passwd
sudo mosquitto_passwd -b /etc/mosquitto/passwd admin qwe123456

# Configure MQTT
sudo cat > /etc/mosquitto/conf.d/default.conf << EOF
allow_anonymous false
password_file /etc/mosquitto/passwd
listener 1883
EOF

# Restart MQTT service
systemctl restart mosquitto
systemctl enable mosquitto

apt --yes install -y debian-keyring debian-archive-keyring apt-transport-https
curl -L https://dl.cloudsmith.io/public/caddy/stable/gpg.key | sudo apt-key add -
curl -L https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt | sudo tee -a /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy
echo "Script complete! Rebooting..."
reboot