#!/bin/bash

# =====================================================
# RDP Platform Backup Script
# Version: 1.0
# Date: 2026-02-22
# =====================================================

set -e

# Configuration
RDP_DIR="/opt/rdp"
BACKUP_DIR="/opt/rdp/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="rdp"
DB_USER="rdp_user"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

log_info() {
    echo -e "[INFO] $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create backup directory
mkdir -p "$BACKUP_DIR"

log_info "Starting backup..."

# Backup PostgreSQL database
log_info "Backing up PostgreSQL database..."
pg_dump -U "$DB_USER" -d "$DB_NAME" -F c -b -v -f "$BACKUP_DIR/db_$DATE.dump"

# Backup configuration files
log_info "Backing up configuration files..."
tar -czf "$BACKUP_DIR/config_$DATE.tar.gz" -C "$RDP_DIR" config/

# Backup uploaded files
log_info "Backing up files..."
tar -czf "$BACKUP_DIR/files_$DATE.tar.gz" -C "$RDP_DIR" data/files/

# Backup logs
log_info "Backing up logs..."
tar -czf "$BACKUP_DIR/logs_$DATE.tar.gz" -C "$RDP_DIR" logs/

# Create backup manifest
cat > "$BACKUP_DIR/manifest_$DATE.txt" << EOF
RDP Platform Backup Manifest
===========================
Date: $DATE
Hostname: $(hostname)

Backup Contents:
- db_$DATE.dump (PostgreSQL database)
- config_$DATE.tar.gz (Configuration files)
- files_$DATE.tar.gz (Uploaded files)
- logs_$DATE.tar.gz (Log files)

Database: $DB_NAME
EOF

# Cleanup old backups (keep last 7 days)
log_info "Cleaning up old backups..."
find "$BACKUP_DIR" -name "db_*.dump" -mtime +7 -delete
find "$BACKUP_DIR" -name "config_*.tar.gz" -mtime +7 -delete
find "$BACKUP_DIR" -name "files_*.tar.gz" -mtime +7 -delete
find "$BACKUP_DIR" -name "logs_*.tar.gz" -mtime +7 -delete

# Create latest symlinks
ln -sf "$BACKUP_DIR/db_$DATE.dump" "$BACKUP_DIR/db_latest.dump"
ln -sf "$BACKUP_DIR/config_$DATE.tar.gz" "$BACKUP_DIR/config_latest.tar.gz"

echo ""
echo -e "${GREEN}Backup completed successfully!${NC}"
echo "Backup location: $BACKUP_DIR"
echo "Backup date: $DATE"
