#!/bin/bash

# ═══════════════════════════════════════════════════════════════════════════
# SSL Certificate Generation Script
# ═══════════════════════════════════════════════════════════════════════════
# This script generates self-signed SSL certificates for development/testing.
# For production, use certificates from a Certificate Authority (e.g., Let's Encrypt).
#
# Usage: ./scripts/generate-certs.sh
# ═══════════════════════════════════════════════════════════════════════════

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

CERTS_DIR="certs"
CERT_FILE="$CERTS_DIR/cert.pem"
KEY_FILE="$CERTS_DIR/key.pem"

echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}SSL Certificate Generation for Development/Testing${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"

# Create certs directory if it doesn't exist
if [ ! -d "$CERTS_DIR" ]; then
    echo -e "${YELLOW}[*] Creating certificates directory: $CERTS_DIR${NC}"
    mkdir -p "$CERTS_DIR"
fi

# Check if certificates already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo -e "${YELLOW}[!] Certificates already exist:${NC}"
    echo "    - $CERT_FILE"
    echo "    - $KEY_FILE"
    
    # Show certificate info
    echo -e "${YELLOW}[*] Certificate Information:${NC}"
    openssl x509 -in "$CERT_FILE" -text -noout | grep -E "Subject:|Issuer:|Not Before|Not After"
    
    read -p "Do you want to regenerate? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${GREEN}[✓] Using existing certificates.${NC}"
        exit 0
    fi
fi

# Generate self-signed certificate valid for 365 days
echo -e "${YELLOW}[*] Generating self-signed certificate...${NC}"

openssl req -x509 \
    -newkey rsa:2048 \
    -keyout "$KEY_FILE" \
    -out "$CERT_FILE" \
    -days 365 \
    -nodes \
    -subj "/C=VN/ST=HCM/L=HCM/O=HR Management/CN=localhost" \
    2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}[✓] Self-signed certificate generated successfully!${NC}"
    echo ""
    echo -e "${YELLOW}Certificate Details:${NC}"
    echo "  - Location: $CERT_FILE"
    echo "  - Private Key: $KEY_FILE"
    echo ""
    openssl x509 -in "$CERT_FILE" -text -noout | grep -E "Subject:|Issuer:|Not Before|Not After"
    echo ""
    echo -e "${YELLOW}[!] Important Notes:${NC}"
    echo "  - This certificate is self-signed and suitable for development only"
    echo "  - For production, use certificates from Let's Encrypt or your CA"
    echo "  - Browser will show security warnings - this is normal for self-signed certs"
    echo "  - Certificate expires in 365 days"
else
    echo -e "${RED}[✗] Failed to generate certificate${NC}"
    exit 1
fi

# Set proper permissions
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo -e "${GREEN}[✓] Done!${NC}"
