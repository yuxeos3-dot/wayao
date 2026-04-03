#!/bin/bash
# BrandSite Pro - Data Backup Script
# Usage: bash backup.sh [data_dir] [backup_dir]

DATA_DIR="${1:-/opt/brandsite/data}"
BACKUP_DIR="${2:-/backup/brandsite}"
DATE=$(date +%Y%m%d)
KEEP_DAYS=30

mkdir -p "$BACKUP_DIR"

echo "[$(date)] Starting backup..."

# 1. Backup SQLite database
if [ -f "$DATA_DIR/brandsite.db" ]; then
    sqlite3 "$DATA_DIR/brandsite.db" "PRAGMA wal_checkpoint(FULL);" 2>/dev/null
    cp "$DATA_DIR/brandsite.db" "$BACKUP_DIR/brandsite-$DATE.db"
    gzip -f "$BACKUP_DIR/brandsite-$DATE.db"
    echo "  DB backup: brandsite-$DATE.db.gz"
fi

# 2. Backup templates
if [ -d "/opt/brandsite/templates" ]; then
    tar -czf "$BACKUP_DIR/templates-$DATE.tar.gz" -C /opt/brandsite templates/ 2>/dev/null
    echo "  Templates backup: templates-$DATE.tar.gz"
fi

# 3. Backup env
if [ -f "/opt/brandsite/.env" ]; then
    cp /opt/brandsite/.env "$BACKUP_DIR/env-$DATE.bak"
    echo "  Config backup"
fi

# 4. Cleanup old backups
find "$BACKUP_DIR" -name "*.gz" -mtime +$KEEP_DAYS -delete
find "$BACKUP_DIR" -name "*.bak" -mtime +$KEEP_DAYS -delete
echo "  Cleaned backups older than $KEEP_DAYS days"

echo "[$(date)] Backup complete"
