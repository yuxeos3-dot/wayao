#!/bin/bash
# Generate Nginx configs for all active domains
# Usage: ./gen_configs.sh /path/to/brandsite.db

DB="${1:-/opt/brandsite/data/brandsite.db}"
TEMPLATE_DIR="$(dirname "$0")"
OUTPUT_DIR="/etc/nginx/sites-available"
ENABLED_DIR="/etc/nginx/sites-enabled"

if [ ! -f "$DB" ]; then
    echo "Database not found: $DB"
    exit 1
fi

echo "Generating Nginx configs from $DB..."

# Get all active domains
DOMAINS=$(sqlite3 "$DB" "SELECT domain FROM domains WHERE status IN ('active','built')")

for domain in $DOMAINS; do
    CONF="$OUTPUT_DIR/$domain.conf"
    echo "  Generating: $CONF"
    sed "s/DOMAIN_NAME/$domain/g" "$TEMPLATE_DIR/site.conf.tmpl" > "$CONF"

    # Enable if not already
    if [ ! -L "$ENABLED_DIR/$domain.conf" ]; then
        ln -sf "$CONF" "$ENABLED_DIR/$domain.conf"
    fi
done

echo "Testing nginx config..."
nginx -t && echo "OK - run 'systemctl reload nginx' to apply"
