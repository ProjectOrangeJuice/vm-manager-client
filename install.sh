# !/bin/bash
echo "Getting latest asset"
apt install -y curl jq
# require two args
if [ $# -lt 1 ]; then
    echo "Usage: $0 <client name> <server address> <allow insecure>"
    exit 1
fi
echo "Given the name $1 and server address $2"

ASSET_URL=$(curl -s "https://api.github.com/repos/ProjectOrangeJuice/vm-manager-client/releases/latest" | jq -r ".assets[0].browser_download_url")
echo "Downloading asset from $ASSET_URL"
curl -Lo /usr/local/bin/vm-manager-client -k $ASSET_URL

echo "Setting permissions"
chmod +x /usr/local/bin/vm-manager-client

echo "Creating config directory"
mkdir -p /etc/vm-manager-client

echo "Creating service"
echo "[Unit]
Description=VM Manager Client
Wants=network-online.target
After=network.target network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/etc/vm-manager-client
ExecStart=/usr/local/bin/vm-manager-client
Restart=always

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/vm-manager-client.service

echo "Starting service"

systemctl daemon-reload
systemctl enable vm-manager-client
systemctl start vm-manager-client

sleep 5

echo "Edit the config file at /etc/vm-manager-client/config.json"
echo "And add your server cert to /etc/vm-manager-client/keys/server-cert.pem"

# replace the config name with arg 1 using jq
jq --arg name "$1" '.Name = $name' /etc/vm-manager-client/config.json > /etc/vm-manager-client/config.json.tmp && mv /etc/vm-manager-client/config.json.tmp /etc/vm-manager-client/config.json

# if arg 2 is set, replace the server address using jq
if [ -n "$2" ]; then
    apt install -y openssl
    jq --arg address "$2:8080" '.ServerAddress = $address' /etc/vm-manager-client/config.json > /etc/vm-manager-client/config.json.tmp && mv /etc/vm-manager-client/config.json.tmp /etc/vm-manager-client/config.json
    # download the server cert via curl
    echo | openssl s_client -servername $2 -connect $2:8080 | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > /etc/vm-manager-client/keys/server-cert.pem
fi

# if arg 3 is set, replace the allow insecure
if [ -n "$3" ]; then
    jq  '.AllowInsecure = true' /etc/vm-manager-client/config.json > /etc/vm-manager-client/config.json.tmp && mv /etc/vm-manager-client/config.json.tmp /etc/vm-manager-client/config.json
fi

systemctl restart vm-manager-client
