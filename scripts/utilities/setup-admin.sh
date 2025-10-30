#!/bin/bash

# Setup admin user script
# This script creates an admin user with the specified credentials

echo "=== Admin User Setup Script ==="
echo "This will create an admin user for the DevPanel"
echo ""

# Set environment variables if needed
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_NAME=${DB_NAME:-portfolio}
export DB_USER=${DB_USER:-postgres}
export DB_SSL_MODE=${DB_SSL_MODE:-disable}

# Admin credentials
ADMIN_EMAIL="admin@jadenrazo.dev"
ADMIN_PASSWORD="JadenRazoAdmin5795%"

echo "Creating admin user with email: $ADMIN_EMAIL"

# Create SQL script to insert admin user
cat > /tmp/create_admin.sql << EOF
-- Check if user exists
DO \$\$
BEGIN
    -- Check if the user already exists
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = '$ADMIN_EMAIL') THEN
        -- Create new admin user
        INSERT INTO users (
            id,
            username,
            email,
            password_hash,
            full_name,
            role,
            is_active,
            is_verified,
            created_at,
            updated_at
        ) VALUES (
            uuid_generate_v4(),
            'admin',
            '$ADMIN_EMAIL',
            '\$2a\$10\$DUMMY_HASH', -- This will be updated below
            'Administrator',
            'admin',
            true,
            true,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        );
        
        RAISE NOTICE 'Admin user created successfully';
    ELSE
        -- Update existing user to admin
        UPDATE users 
        SET 
            role = 'admin',
            is_active = true,
            is_verified = true,
            updated_at = CURRENT_TIMESTAMP
        WHERE email = '$ADMIN_EMAIL';
        
        RAISE NOTICE 'Existing user updated to admin';
    END IF;
END\$\$;
EOF

echo ""
echo "To complete the setup, you need to:"
echo "1. Run the database migration to create the certification tables"
echo "2. Generate a proper password hash for the admin user"
echo ""
echo "Since the backend uses bcrypt, you'll need to run the Go application to generate the proper password hash."
echo ""
echo "Alternatively, you can use the admin setup endpoint at: https://jadenrazo.dev/api/v1/auth/admin/setup/request"