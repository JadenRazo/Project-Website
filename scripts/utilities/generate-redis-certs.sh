#!/bin/bash

# Generate self-signed TLS certificates for Redis
# This script creates certificates for development and testing
# For production, use proper CA-signed certificates

set -e

CERT_DIR="/main/Project-Website/deploy/redis/certs"
DAYS_VALID=365

echo "Generating Redis TLS certificates..."

# Create certificate directory
mkdir -p "$CERT_DIR"

# Generate CA private key
openssl genrsa -out "$CERT_DIR/ca.key" 4096

# Generate CA certificate
openssl req -new -x509 -days $DAYS_VALID -key "$CERT_DIR/ca.key" \
    -out "$CERT_DIR/ca.crt" \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=Redis CA"

# Generate Redis server private key
openssl genrsa -out "$CERT_DIR/redis.key" 4096

# Generate Redis server certificate signing request
openssl req -new -key "$CERT_DIR/redis.key" \
    -out "$CERT_DIR/redis.csr" \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=redis"

# Sign the Redis certificate with CA
openssl x509 -req -days $DAYS_VALID \
    -in "$CERT_DIR/redis.csr" \
    -CA "$CERT_DIR/ca.crt" \
    -CAkey "$CERT_DIR/ca.key" \
    -CAcreateserial \
    -out "$CERT_DIR/redis.crt"

# Set proper permissions
chmod 600 "$CERT_DIR"/*.key
chmod 644 "$CERT_DIR"/*.crt

# Clean up CSR
rm -f "$CERT_DIR/redis.csr"

echo "Redis TLS certificates generated successfully in $CERT_DIR"
echo "Files created:"
echo "  - ca.crt (CA certificate)"
echo "  - ca.key (CA private key)"
echo "  - redis.crt (Redis server certificate)"
echo "  - redis.key (Redis server private key)"
echo ""
echo "Note: These are self-signed certificates for development/testing."
echo "For production, use certificates from a trusted Certificate Authority."
