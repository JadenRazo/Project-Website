#!/bin/bash

[ -f .env ] && export $(grep -v '^#' .env | xargs)

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Database configuration
DB_NAME="${DB_NAME:-portfolio_db}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
BACKUP_DIR="backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo -e "${GREEN}Starting database backup...${NC}"

# Perform the backup
pg_dump -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_DIR/backup_$TIMESTAMP.sql"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Backup completed successfully!${NC}"
    echo -e "Backup file: $BACKUP_DIR/backup_$TIMESTAMP.sql"
else
    echo -e "${RED}Backup failed!${NC}"
    exit 1
fi
