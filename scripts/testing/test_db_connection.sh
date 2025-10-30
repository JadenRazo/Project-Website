#!/bin/bash

set -e

echo "========================================"
echo "Database Connection Test Script"
echo "========================================"
echo ""

# Load environment variables from .env file
if [ -f "/main/Project-Website/backend/.env" ]; then
    export $(cat /main/Project-Website/backend/.env | grep -v '^#' | xargs)
    echo "✓ Loaded environment variables from .env"
else
    echo "✗ .env file not found"
    exit 1
fi

echo ""
echo "Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  SSL Mode: $DB_SSL_MODE"
echo ""

# Test 1: Check if PostgreSQL is accessible
echo "Test 1: Checking PostgreSQL accessibility..."
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "SELECT version();" > /dev/null 2>&1; then
    echo "✓ PostgreSQL server is accessible"
else
    echo "✗ Cannot connect to PostgreSQL server"
    exit 1
fi

# Test 2: Check if database exists
echo ""
echo "Test 2: Checking if database '$DB_NAME' exists..."
DB_EXISTS=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")
if [ "$DB_EXISTS" = "1" ]; then
    echo "✓ Database '$DB_NAME' exists"
else
    echo "✗ Database '$DB_NAME' does not exist"
    echo ""
    echo "To create the database, run:"
    echo "  PGPASSWORD='$DB_PASSWORD' psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c 'CREATE DATABASE $DB_NAME;'"
    exit 1
fi

# Test 3: Check if required extensions are installed
echo ""
echo "Test 3: Checking required extensions..."
EXTENSIONS=("uuid-ossp" "pgcrypto" "citext")
for ext in "${EXTENSIONS[@]}"; do
    EXT_EXISTS=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 1 FROM pg_extension WHERE extname='$ext'")
    if [ "$EXT_EXISTS" = "1" ]; then
        echo "  ✓ Extension '$ext' is installed"
    else
        echo "  ✗ Extension '$ext' is NOT installed"
    fi
done

# Test 4: Count tables
echo ""
echo "Test 4: Checking database tables..."
TABLE_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
echo "  Found $TABLE_COUNT tables in the database"

# Test 5: Check for critical tables
echo ""
echo "Test 5: Checking critical tables..."
CRITICAL_TABLES=("users" "projects" "shortened_urls" "skills" "certifications")
for table in "${CRITICAL_TABLES[@]}"; do
    TABLE_EXISTS=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '$table'")
    if [ "$TABLE_EXISTS" = "1" ]; then
        ROW_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM $table")
        echo "  ✓ Table '$table' exists ($ROW_COUNT rows)"
    else
        echo "  ✗ Table '$table' does NOT exist"
    fi
done

# Test 6: Check for visitor analytics tables (should be in schema.sql)
echo ""
echo "Test 6: Checking visitor analytics tables..."
VISITOR_TABLES=("visitor_sessions" "page_views" "privacy_consents" "visitor_metrics" "visitor_daily_summary" "visitor_realtime" "visitor_locations")
MISSING_VISITOR_TABLES=()
for table in "${VISITOR_TABLES[@]}"; do
    TABLE_EXISTS=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '$table'")
    if [ "$TABLE_EXISTS" = "1" ]; then
        echo "  ✓ Table '$table' exists"
    else
        echo "  ✗ Table '$table' does NOT exist (defined in schema.sql)"
        MISSING_VISITOR_TABLES+=("$table")
    fi
done

# Test 7: Check for messaging tables with correct naming
echo ""
echo "Test 7: Checking messaging table naming consistency..."
OLD_MESSAGING_TABLES=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('channels', 'messages', 'channel_members', 'message_attachments', 'message_reactions', 'read_receipts')")
NEW_MESSAGING_TABLES=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('messaging_channels', 'messaging_messages', 'messaging_channel_members', 'messaging_attachments', 'messaging_reactions', 'messaging_read_receipts')")

if [ "$OLD_MESSAGING_TABLES" -gt 0 ]; then
    echo "  ⚠ Found $OLD_MESSAGING_TABLES old-style messaging tables (channels, messages, etc.)"
    echo "    Schema.sql expects prefixed tables (messaging_channels, messaging_messages, etc.)"
fi
if [ "$NEW_MESSAGING_TABLES" -gt 0 ]; then
    echo "  ✓ Found $NEW_MESSAGING_TABLES new-style messaging tables"
fi

# Test 8: Test actual query execution
echo ""
echo "Test 8: Testing query execution..."
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT COUNT(*) FROM users;" > /dev/null 2>&1; then
    USER_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM users")
    echo "✓ Successfully queried users table: $USER_COUNT users found"
else
    echo "✗ Failed to query users table"
fi

# Test 9: Check for admin user
echo ""
echo "Test 9: Checking for admin user..."
ADMIN_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT COUNT(*) FROM users WHERE role = 'admin'")
if [ "$ADMIN_COUNT" -gt 0 ]; then
    echo "✓ Found $ADMIN_COUNT admin user(s)"
    echo ""
    echo "  Admin users:"
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT username, email, created_at FROM users WHERE role = 'admin' ORDER BY created_at;"
else
    echo "⚠ No admin users found"
    echo "  To create an admin user, run:"
    echo "    go run /main/Project-Website/backend/cmd/init-admin/main.go"
fi

# Summary
echo ""
echo "========================================"
echo "Summary"
echo "========================================"
echo "Total tables: $TABLE_COUNT"

if [ ${#MISSING_VISITOR_TABLES[@]} -gt 0 ]; then
    echo ""
    echo "⚠ Schema discrepancies detected:"
    echo "  Missing visitor analytics tables: ${#MISSING_VISITOR_TABLES[@]}"
    echo ""
    echo "To fix schema discrepancies, apply the schema.sql file:"
    echo "  PGPASSWORD='$DB_PASSWORD' psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f /main/Project-Website/backend/schema.sql"
fi

if [ "$OLD_MESSAGING_TABLES" -gt 0 ] && [ "$NEW_MESSAGING_TABLES" -eq 0 ]; then
    echo ""
    echo "⚠ Messaging table naming mismatch detected"
    echo "  Current database uses old naming (channels, messages)"
    echo "  Schema.sql expects new naming (messaging_channels, messaging_messages)"
    echo "  This may require a migration or schema update"
fi

echo ""
echo "✓ Database connection test completed"
