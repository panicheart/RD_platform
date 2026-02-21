#!/bin/bash

RDP_DIR="/opt/rdp"
BACKUP_DIR="/opt/rdp/backup"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/rdp_backup_${TIMESTAMP}.tar.gz"

mkdir -p $BACKUP_DIR

echo "Creating backup..."

pg_dump -U rdp -h localhost rdp > "${BACKUP_DIR}/db_${TIMESTAMP}.sql"

tar -czf $BACKUP_FILE \
    -C $RDP_DIR \
    --exclude='backup' \
    --exclude='logs' \
    --exclude='*.log' \
    .

echo "Backup created: $BACKUP_FILE"

find $BACKUP_DIR -name "rdp_backup_*.tar.gz" -mtime +7 -delete
find $BACKUP_DIR -name "db_*.sql" -mtime +7 -delete

echo "Old backups cleaned (keeping last 7 days)"
