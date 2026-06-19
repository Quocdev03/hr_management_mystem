#!/bin/bash

# ═══════════════════════════════════════════════════════════════════════════
# Pre-Deployment Validation Script
# ═══════════════════════════════════════════════════════════════════════════
# This script validates that all security requirements are met before deployment
#
# Usage: ./scripts/validate-production.sh
# ═══════════════════════════════════════════════════════════════════════════

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Production Deployment Validation${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Check if .env file exists
echo -e "${YELLOW}[*] Checking environment configuration...${NC}"
if [ ! -f .env ]; then
    echo -e "${RED}[✗] ERROR: .env file not found${NC}"
    echo "    → Copy .env.example to .env and update with your values"
    ((ERRORS++))
else
    echo -e "${GREEN}[✓] .env file exists${NC}"
fi

# Check if .env.example exists
if [ ! -f .env.example ]; then
    echo -e "${RED}[✗] ERROR: .env.example file not found${NC}"
    ((ERRORS++))
else
    echo -e "${GREEN}[✓] .env.example file exists${NC}"
fi

echo ""
echo -e "${YELLOW}[*] Validating environment variables...${NC}"

# Function to check if variable is set and has strong value
check_env_var() {
    local var=$1
    local min_length=$2
    local description=$3
    local value=$(grep "^${var}=" .env | cut -d '=' -f 2- || echo "")
    
    if [ -z "$value" ]; then
        echo -e "${RED}[✗] ERROR: $var is not set${NC}"
        ((ERRORS++))
        return 1
    fi
    
    if [ "$min_length" -gt 0 ] && [ ${#value} -lt "$min_length" ]; then
        echo -e "${RED}[✗] ERROR: $var is too short (min $min_length chars, got ${#value})${NC}"
        ((ERRORS++))
        return 1
    fi
    
    # Check if value contains default placeholder
    if [[ "$value" == *"CHANGE_ME"* ]]; then
        echo -e "${RED}[✗] ERROR: $var still contains CHANGE_ME placeholder${NC}"
        ((ERRORS++))
        return 1
    fi
    
    echo -e "${GREEN}[✓] $description: OK (${#value} chars)${NC}"
    return 0
}

# Validate critical variables
check_env_var "JWT_SECRET" 32 "JWT_SECRET (strong encryption key)"
check_env_var "DB_PASSWORD" 16 "DB_PASSWORD (strong database password)"
check_env_var "MYSQL_ROOT_PASSWORD" 16 "MYSQL_ROOT_PASSWORD (strong root password)"

echo ""
echo -e "${YELLOW}[*] Checking SSL certificates...${NC}"
if [ -f certs/cert.pem ] && [ -f certs/key.pem ]; then
    echo -e "${GREEN}[✓] SSL certificate files found${NC}"
    
    # Check certificate expiration
    EXPIRATION=$(openssl x509 -in certs/cert.pem -noout -enddate | cut -d= -f 2)
    echo -e "${GREEN}    Expires: $EXPIRATION${NC}"
else
    echo -e "${YELLOW}[!] WARNING: SSL certificates not found${NC}"
    echo "    → Run: ./scripts/generate-certs.sh"
    ((WARNINGS++))
fi

echo ""
echo -e "${YELLOW}[*] Checking Docker and Docker Compose...${NC}"
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    echo -e "${GREEN}[✓] Docker found: $DOCKER_VERSION${NC}"
else
    echo -e "${RED}[✗] ERROR: Docker not installed${NC}"
    ((ERRORS++))
fi

if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version)
    echo -e "${GREEN}[✓] Docker Compose found: $COMPOSE_VERSION${NC}"
else
    echo -e "${RED}[✗] ERROR: Docker Compose not installed${NC}"
    ((ERRORS++))
fi

echo ""
echo -e "${YELLOW}[*] Checking required files...${NC}"
REQUIRED_FILES=(
    "docker-compose.prod.yml"
    "frontend/nginx.prod.conf"
    "backend/Dockerfile"
    "frontend/Dockerfile"
    "database.sql"
)

for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}[✓] $file exists${NC}"
    else
        echo -e "${RED}[✗] ERROR: $file not found${NC}"
        ((ERRORS++))
    fi
done

echo ""
echo -e "${YELLOW}[*] Security checks...${NC}"

# Check if secrets are in .env (not in docker-compose files)
if grep -q "password123" docker-compose.prod.yml 2>/dev/null; then
    echo -e "${RED}[✗] ERROR: Hardcoded password found in docker-compose.prod.yml${NC}"
    ((ERRORS++))
else
    echo -e "${GREEN}[✓] No hardcoded passwords in docker-compose.prod.yml${NC}"
fi

if grep -q "your_jwt_secret_key" docker-compose.prod.yml 2>/dev/null; then
    echo -e "${RED}[✗] ERROR: Default JWT secret found in docker-compose.prod.yml${NC}"
    ((ERRORS++))
else
    echo -e "${GREEN}[✓] No default JWT secrets in docker-compose.prod.yml${NC}"
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Validation Summary${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}[✓] All checks passed! Ready for production deployment.${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}[!] $WARNINGS warning(s) found. Review recommended before deployment.${NC}"
    exit 0
else
    echo -e "${RED}[✗] $ERRORS error(s) found. Fix before deployment.${NC}"
    exit 1
fi
