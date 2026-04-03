#!/bin/bash
# Batch SSL certificate setup for all domains
# Usage: bash batch_ssl.sh /path/to/brandsite.db [email]
set -e

DB="${1:-/opt/brandsite/data/brandsite.db}"
EMAIL="${2:-admin@yourdomain.com}"

if [ ! -f "$DB" ]; then
    echo "Database not found: $DB"
    exit 1
fi

DOMAINS=$(sqlite3 "$DB" "SELECT domain FROM domains WHERE status IN ('active','built')")

for domain in $DOMAINS; do
    echo "Setting up SSL for: $domain"
    # Check DNS resolves
    if ! host "$domain" > /dev/null 2>&1; then
        echo "  SKIP: DNS not resolving for $domain"
        continue
    fi
    certbot --nginx -d "$domain" --non-interactive --agree-tos -m "$EMAIL" || echo "  WARN: certbot failed for $domain"
done

echo "Done! Reloading nginx..."
systemctl reload nginx
