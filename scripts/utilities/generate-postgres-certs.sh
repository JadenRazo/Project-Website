#!/bin/bash
set -e

echo "=========================================="
echo "PostgreSQL SSL/TLS Certificate Generator"
echo "=========================================="
echo

CERT_DIR="/main/Project-Website/deploy/postgres/certs"

if [ ! -d "$CERT_DIR" ]; then
    echo "Creating certificate directory: $CERT_DIR"
    mkdir -p "$CERT_DIR"
fi

echo "Generating PostgreSQL SSL certificates..."
echo

echo "[1/4] Generating Certificate Authority (CA)..."
openssl req -new -x509 -days 3650 -nodes \
    -out "$CERT_DIR/ca.crt" \
    -keyout "$CERT_DIR/ca.key" \
    -subj "/C=US/ST=State/L=City/O=Portfolio/CN=PostgreSQL-CA" \
    2>/dev/null

echo "[2/4] Generating PostgreSQL server private key..."
openssl genrsa -out "$CERT_DIR/server.key" 2048 2>/dev/null

echo "[3/4] Generating server certificate signing request..."
openssl req -new \
    -key "$CERT_DIR/server.key" \
    -out "$CERT_DIR/server.csr" \
    -subj "/C=US/ST=State/L=City/O=Portfolio/CN=postgres" \
    2>/dev/null

echo "[4/4] Signing server certificate with CA..."
openssl x509 -req -days 3650 \
    -in "$CERT_DIR/server.csr" \
    -CA "$CERT_DIR/ca.crt" \
    -CAkey "$CERT_DIR/ca.key" \
    -CAcreateserial \
    -out "$CERT_DIR/server.crt" \
    2>/dev/null

echo
echo "Setting proper file permissions..."
chmod 600 "$CERT_DIR/server.key"
chmod 600 "$CERT_DIR/ca.key"
chmod 644 "$CERT_DIR/server.crt"
chmod 644 "$CERT_DIR/ca.crt"
chmod 644 "$CERT_DIR/server.csr"

echo
echo "Cleaning up temporary files..."
rm -f "$CERT_DIR/server.csr"

echo
echo "✅ PostgreSQL SSL certificates generated successfully!"
echo
echo "Generated files:"
echo "  - CA Certificate:     $CERT_DIR/ca.crt"
echo "  - CA Private Key:     $CERT_DIR/ca.key"
echo "  - Server Certificate: $CERT_DIR/server.crt"
echo "  - Server Private Key: $CERT_DIR/server.key"
echo
echo "Certificate validity: 10 years (3650 days)"
echo
echo "Note: PostgreSQL requires strict permissions on server.key (600)"
echo "      All certificates are ready for Docker volume mounting"
echo
