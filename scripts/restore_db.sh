#!/bin/bash

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Database configuration
DB_NAME="portfolio_db"
DB_USER="postgres"
DB_HOST="localhost"
BACKUP_DIR="backups"

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}Backup directory not found!${NC}"
    exit 1
fi

# List available backups
echo -e "${YELLOW}Available backups:${NC}"
ls -1 "$BACKUP_DIR"/*.sql 2>/dev/null | sed 's/.*\///'

# Prompt for backup file
echo -e "\n${YELLOW}Enter the backup file name to restore (e.g., backup_20240315_123456.sql):${NC}"
read -r backup_file

# Check if backup file exists
if [ ! -f "$BACKUP_DIR/$backup_file" ]; then
    echo -e "${RED}Backup file not found!${NC}"
    exit 1
fi

# Confirm restoration
echo -e "\n${YELLOW}Are you sure you want to restore from $backup_file? This will overwrite the current database. (y/N)${NC}"
read -r confirm

if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo -e "${RED}Restoration cancelled.${NC}"
    exit 0
fi

echo -e "${GREEN}Starting database restoration...${NC}"

# Drop and recreate the database
psql -h "$DB_HOST" -U "$DB_USER" -c "DROP DATABASE IF EXISTS ${DB_NAME}_temp;"
psql -h "$DB_HOST" -U "$DB_USER" -c "CREATE DATABASE ${DB_NAME}_temp;"

# Restore the backup to the temporary database
psql -h "$DB_HOST" -U "$DB_USER" -d "${DB_NAME}_temp" < "$BACKUP_DIR/$backup_file"

if [ $? -eq 0 ]; then
    # Drop the current database and rename the temporary one
    psql -h "$DB_HOST" -U "$DB_USER" -c "DROP DATABASE IF EXISTS ${DB_NAME};"
    psql -h "$DB_HOST" -U "$DB_USER" -c "ALTER DATABASE ${DB_NAME}_temp RENAME TO ${DB_NAME};"
    
    echo -e "${GREEN}Database restored successfully!${NC}"
else
    echo -e "${RED}Restoration failed!${NC}"
    # Clean up temporary database
    psql -h "$DB_HOST" -U "$DB_USER" -c "DROP DATABASE IF EXISTS ${DB_NAME}_temp;"
    exit 1
fi
