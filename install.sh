echo "Getting latest asset"

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
After=network.target

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

echo "Edit the config file at /etc/vm-manager-client/config.json"
echo "And add your server cert to /etc/vm-manager-client/keys/server-cert.pem"