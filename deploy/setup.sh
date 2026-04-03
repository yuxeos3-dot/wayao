#!/bin/bash
# BrandSite Pro - Server Setup Script
# Usage: bash setup.sh
set -e

echo "=== BrandSite Pro Server Setup ==="

# 1. System deps
apt update && apt install -y nginx sqlite3 certbot python3-certbot-nginx ufw git dnsutils

# 2. Hugo
if ! command -v hugo &> /dev/null; then
    echo "Installing Hugo..."
    curl -L https://github.com/gohugoio/hugo/releases/download/v0.147.4/hugo_extended_0.147.4_linux-amd64.tar.gz -o /tmp/hugo.tar.gz
    tar -xzf /tmp/hugo.tar.gz -C /usr/local/bin hugo
    hugo version
fi

# 3. Create directories
mkdir -p /opt/brandsite/{data,templates,frontend/dist}
mkdir -p /var/www/sites
mkdir -p /var/log/brandsite

# 4. Firewall
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# 5. Nginx rate limit zone
grep -q "limit_req_zone" /etc/nginx/nginx.conf || \
    sed -i '/http {/a\    limit_req_zone $binary_remote_addr zone=tracker:10m rate=10r/s;' /etc/nginx/nginx.conf

# 6. Systemd service
cat > /etc/systemd/system/brandsite.service << 'UNIT'
[Unit]
Description=BrandSite Pro Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/brandsite
ExecStart=/opt/brandsite/brandsite-server
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=DATA_DIR=/opt/brandsite/data
Environment=TEMPLATE_DIR=/opt/brandsite/templates
Environment=ADMIN_DIR=/opt/brandsite/frontend/dist
Environment=HUGO_PATH=/usr/local/bin/hugo
Environment=BUILD_OUTPUT=/var/www/sites
StandardOutput=append:/var/log/brandsite/server.log
StandardError=append:/var/log/brandsite/error.log
ProtectSystem=strict
ReadWritePaths=/opt/brandsite/data /var/www/sites /var/log/brandsite /tmp

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
echo "=== Setup Complete ==="
echo "Next steps:"
echo "  1. Build: cd /path/to/wayao && CGO_ENABLED=1 go build -o brandsite-server ./cmd/server/"
echo "  2. Copy binary + templates + frontend/dist to /opt/brandsite/"
echo "  3. Start: systemctl enable --now brandsite"
echo "  4. Access: http://YOUR_IP:8080/admin/"
