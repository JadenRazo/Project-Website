#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Database Restore Script${NC}"
echo -e "${BLUE}========================================${NC}"

# Navigate to backend directory
cd "$(dirname "${BASH_SOURCE[0]}")/.."
BACKEND_DIR="$(pwd)"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"

# Set default variables
DB_TYPE=${DB_TYPE:-"sqlite"}
BACKUP_DIR="${BACKEND_DIR}/backups"

# Check if backup directory exists
if [ ! -d "${BACKUP_DIR}" ]; then
    echo -e "${RED}Error: Backup directory not found: ${BACKUP_DIR}${NC}"
    exit 1
fi

# Load config values if possible
if [ -f "${BACKEND_DIR}/config/development.yaml" ]; then
    echo -e "${GREEN}Loading database configuration from development.yaml...${NC}"
    # This is a simple extraction - for production use a proper YAML parser
    if grep -q "database:" "${BACKEND_DIR}/config/development.yaml"; then
        DB_TYPE=$(grep -A 10 "database:" "${BACKEND_DIR}/config/development.yaml" | grep "driver:" | head -n 1 | awk -F': ' '{print $2}' | tr -d ' "')
        DB_DSN=$(grep -A 10 "database:" "${BACKEND_DIR}/config/development.yaml" | grep "dsn:" | head -n 1 | awk -F': ' '{print $2}' | tr -d ' "')
        
        # Extract database credentials if using PostgreSQL or MySQL
        if [[ "$DB_TYPE" == "postgres" ]]; then
            DB_HOST=$(echo $DB_DSN | grep -o "host=[^ ]*" | cut -d= -f2)
            DB_PORT=$(echo $DB_DSN | grep -o "port=[^ ]*" | cut -d= -f2)
            DB_NAME=$(echo $DB_DSN | grep -o "dbname=[^ ]*" | cut -d= -f2)
            DB_USER=$(echo $DB_DSN | grep -o "user=[^ ]*" | cut -d= -f2)
            DB_PASSWORD=$(echo $DB_DSN | grep -o "password=[^ ]*" | cut -d= -f2)
        elif [[ "$DB_TYPE" == "mysql" ]]; then
            DB_USER=$(echo $DB_DSN | cut -d: -f1)
            DB_PASSWORD=$(echo $DB_DSN | cut -d: -f2 | cut -d@ -f1)
            DB_HOST=$(echo $DB_DSN | cut -d@ -f2 | cut -d/ -f1 | cut -d: -f1)
            DB_PORT=$(echo $DB_DSN | cut -d@ -f2 | cut -d/ -f1 | cut -d: -f2)
            DB_NAME=$(echo $DB_DSN | cut -d/ -f2 | cut -d? -f1)
        elif [[ "$DB_TYPE" == "sqlite" ]]; then
            DB_FILE=$DB_DSN
        fi
    fi
fi

# List available backups
echo -e "\n${GREEN}Available backups:${NC}"
ls -lht "${BACKUP_DIR}" | grep -v "total"

# Get backup file from command line or prompt user
BACKUP_FILE=$1
if [ -z "$BACKUP_FILE" ]; then
    echo -e "\n${YELLOW}Please enter the backup file to restore (from the list above):${NC}"
    read -p "> " BACKUP_FILE
fi

# Check if the backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    # Try with backup directory prefix
    if [ ! -f "${BACKUP_DIR}/$BACKUP_FILE" ]; then
        echo -e "${RED}Error: Backup file not found: ${BACKUP_FILE}${NC}"
        exit 1
    else
        BACKUP_FILE="${BACKUP_DIR}/$BACKUP_FILE"
    fi
fi

echo -e "\n${GREEN}Database restore information:${NC}"
echo -e "  Database type: ${YELLOW}${DB_TYPE}${NC}"
echo -e "  Backup file: ${YELLOW}${BACKUP_FILE}${NC}"

# Confirm before proceeding
echo -e "\n${YELLOW}Warning: This will overwrite your current database. All current data will be lost.${NC}"
read -p "Do you want to continue? (y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Restore operation canceled.${NC}"
    exit 0
fi

# Perform restore based on database type and backup file extension
if [[ "$BACKUP_FILE" == *".db" ]]; then
    # It's a SQLite binary backup
    if [ -z "$DB_FILE" ]; then
        DB_FILE="${BACKEND_DIR}/data/database.db"
        echo -e "${YELLOW}Using default SQLite database location: ${DB_FILE}${NC}"
    fi
    
    echo -e "\n${GREEN}Restoring SQLite database from binary backup...${NC}"
    # Backup current database
    if [ -f "$DB_FILE" ]; then
        TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
        cp "$DB_FILE" "${DB_FILE}.bak.${TIMESTAMP}"
        echo -e "${YELLOW}Current database backed up to: ${DB_FILE}.bak.${TIMESTAMP}${NC}"
    fi
    
    # Copy the backup to the database location
    cp "$BACKUP_FILE" "$DB_FILE"
    echo -e "${GREEN}SQLite database restored successfully!${NC}"
    
elif [[ "$BACKUP_FILE" == *".pgdump" ]] && [[ "$DB_TYPE" == "postgres" ]]; then
    # It's a PostgreSQL custom format backup
    if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ]; then
        echo -e "${RED}Error: Missing PostgreSQL connection details.${NC}"
        exit 1
    fi
    
    echo -e "\n${GREEN}Restoring PostgreSQL database...${NC}"
    if command -v pg_restore &> /dev/null; then
        # Drop and recreate the database
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS ${DB_NAME};"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE ${DB_NAME};"
        
        # Restore the database
        PGPASSWORD=$DB_PASSWORD pg_restore -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -v "$BACKUP_FILE"
        echo -e "${GREEN}PostgreSQL database restored successfully!${NC}"
    else
        echo -e "${RED}Error: pg_restore command not found. Please install PostgreSQL client tools.${NC}"
        exit 1
    fi
    
elif [[ "$DB_TYPE" == "mysql" ]]; then
    # It's a MySQL backup
    if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ]; then
        echo -e "${RED}Error: Missing MySQL connection details.${NC}"
        exit 1
    fi
    
    echo -e "\n${GREEN}Restoring MySQL database...${NC}"
    if command -v mysql &> /dev/null; then
        # Drop and recreate the database
        mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD -e "DROP DATABASE IF EXISTS ${DB_NAME}; CREATE DATABASE ${DB_NAME};"
        
        # Restore the database
        mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME < "$BACKUP_FILE"
        echo -e "${GREEN}MySQL database restored successfully!${NC}"
    else
        echo -e "${RED}Error: mysql command not found. Please install MySQL client tools.${NC}"
        exit 1
    fi
    
elif [[ "$DB_TYPE" == "sqlite" ]]; then
    # It's a SQLite SQL dump
    if [ -z "$DB_FILE" ]; then
        DB_FILE="${BACKEND_DIR}/data/database.db"
        echo -e "${YELLOW}Using default SQLite database location: ${DB_FILE}${NC}"
    fi
    
    echo -e "\n${GREEN}Restoring SQLite database from SQL dump...${NC}"
    if command -v sqlite3 &> /dev/null; then
        # Backup current database
        if [ -f "$DB_FILE" ]; then
            TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
            cp "$DB_FILE" "${DB_FILE}.bak.${TIMESTAMP}"
            echo -e "${YELLOW}Current database backed up to: ${DB_FILE}.bak.${TIMESTAMP}${NC}"
        fi
        
        # Remove existing database and create new one from backup
        rm -f "$DB_FILE"
        sqlite3 "$DB_FILE" < "$BACKUP_FILE"
        echo -e "${GREEN}SQLite database restored successfully!${NC}"
    else
        echo -e "${RED}Error: sqlite3 command not found. Please install SQLite.${NC}"
        exit 1
    fi
    
else
    echo -e "${RED}Error: Unsupported backup file format or database type mismatch.${NC}"
    echo -e "  Database type: ${DB_TYPE}"
    echo -e "  Backup file: ${BACKUP_FILE}"
    exit 1
fi

echo -e "\n${GREEN}Restore operation completed successfully!${NC}"
