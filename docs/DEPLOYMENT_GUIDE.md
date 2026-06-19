# 🚀 Production Deployment Guide

**Last Updated:** 2026-06-14  
**Status:** ✅ Critical Security Fixes Implemented

---

## 📋 Table of Contents

1. [Pre-Deployment Checklist](#pre-deployment-checklist)
2. [Environment Setup](#environment-setup)
3. [SSL/TLS Configuration](#ssltls-configuration)
4. [Deployment Process](#deployment-process)
5. [Post-Deployment Validation](#post-deployment-validation)
6. [Monitoring & Maintenance](#monitoring--maintenance)
7. [Troubleshooting](#troubleshooting)

---

## ✅ Pre-Deployment Checklist

### Security Requirements (MUST DO)

- [ ] **Secrets Management**
   - [ ] Create `.env` file from `.env.example`
   - [ ] Generate strong passwords (min 16 characters, mix of uppercase, lowercase, numbers, special chars)
   - [ ] Generate JWT secret (min 32 characters, cryptographically random)
   - [ ] **NEVER commit `.env` to Git** (already in `.gitignore`)
   - [ ] Store `.env` in secure vault (AWS Secrets Manager, HashiCorp Vault, etc.)

- [ ] **SSL/TLS Configuration**
   - [ ] For development: Run `./scripts/generate-certs.sh`
   - [ ] For production: Obtain CA-signed certificates (Let's Encrypt recommended)
   - [ ] Place certificates in `certs/` directory
   - [ ] Verify certificate validity: `openssl x509 -in certs/cert.pem -text -noout`

- [ ] **Database**
   - [ ] Change default database user from `root` to `hrm_app_user`
   - [ ] Set strong password for database user
   - [ ] Verify database backup strategy in place
   - [ ] Test database restoration process

- [ ] **Nginx Configuration**
   - [ ] Verify security headers in `frontend/nginx.prod.conf`
   - [ ] Test HTTP → HTTPS redirect
   - [ ] Verify SSL protocols (TLSv1.2, TLSv1.3)
   - [ ] Check rate limiting configuration

- [ ] **Application Configuration**
   - [ ] Set `APP_ENV=production`
   - [ ] Set `APP_SEED=false`
   - [ ] Verify all environment variables are set
   - [ ] Run validation script: `./scripts/validate-production.sh`

---

## 🔧 Environment Setup

### 1. Create Environment Configuration

```bash
# Copy template to actual config
cp .env.example .env

# Edit .env with production values
nano .env
```

### 2. .env File Template (Update these values)

```env
# Application
APP_PORT=8080
APP_ENV=production

# Database (Use strong passwords!)
DB_HOST=mysql
DB_PORT=3306
DB_USER=hrm_app_user
DB_PASSWORD=<generate-strong-password>  # Min 16 chars
DB_NAME=hrm_db
MYSQL_ROOT_PASSWORD=<generate-strong-root-password>  # Min 16 chars

# JWT (Use cryptographically random key!)
JWT_SECRET=<generate-32-char-random-key>  # Min 32 chars, use: openssl rand -base64 32
JWT_EXPIRE_HOUR=1
JWT_REFRESH_EXPIRE_DAY=7

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<generate-strong-password>
REDIS_DB=0
```

### 3. Generate Strong Passwords

```bash
# Generate 32-character random password
openssl rand -base64 32

# Generate 16-character random password
openssl rand -base64 16

# Store these values in .env file
```

---

## 🔐 SSL/TLS Configuration

### Development Setup (Self-Signed Certificates)

```bash
# Generate self-signed certificates (valid for 365 days)
./scripts/generate-certs.sh

# This creates:
# - certs/cert.pem (public certificate)
# - certs/key.pem (private key)
```

### Production Setup (CA-Signed Certificates)

#### Option 1: Let's Encrypt (Recommended, FREE)

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Generate certificate for your domain
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Certificates will be at:
# /etc/letsencrypt/live/yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/yourdomain.com/privkey.pem

# Copy to certs directory
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem certs/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem certs/key.pem
sudo chown $(whoami):$(whoami) certs/*
sudo chmod 600 certs/key.pem
sudo chmod 644 certs/cert.pem
```

#### Option 2: Commercial Certificate Authority

1. Generate CSR (Certificate Signing Request)
2. Submit to your CA (Comodo, Sectigo, etc.)
3. Receive certificate files
4. Place in `certs/` directory with names:
   - `certs/cert.pem` (certificate)
   - `certs/key.pem` (private key)

### Certificate Verification

```bash
# Check certificate expiration
openssl x509 -in certs/cert.pem -noout -dates

# Check certificate details
openssl x509 -in certs/cert.pem -text -noout

# Verify certificate chain
openssl verify -CAfile certs/cert.pem certs/cert.pem
```

---

## 🚀 Deployment Process

> [!NOTE]
> On Windows, you can use the **`start-prod.bat`** script in CMD to manage the production stack easily:
> * Start & Setup: `start-prod.bat setup`
> * Start only: `start-prod.bat`
> * Stop: `start-prod.bat down`
> * Logs: `start-prod.bat logs`
> * Status: `start-prod.bat status`

### Step 1: Validation

```bash
# Run pre-deployment validation
./scripts/validate-production.sh

# Expected output:
# ✓ .env file exists
# ✓ All environment variables validated
# ✓ SSL certificates found
# ✓ Docker and Docker Compose installed
# ✓ All security checks passed
```

### Step 2: Build Images

```bash
# Build production Docker images
docker compose -f docker-compose.prod.yml build

# Verify images built successfully
docker images | grep hrm
```

### Step 3: Database Initialization

```bash
# Start only database services first
docker compose -f docker-compose.prod.yml up -d mysql redis

# Wait for MySQL to be ready (check healthcheck)
docker compose -f docker-compose.prod.yml ps

# Initialize database (if needed)
docker compose -f docker-compose.prod.yml exec mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD}" < database.sql
```

### Step 4: Start All Services

```bash
# Start all services
docker compose -f docker-compose.prod.yml up -d

# Verify all containers are running
docker compose -f docker-compose.prod.yml ps

# Check logs for errors
docker compose -f docker-compose.prod.yml logs -f
```

### Step 5: Run Migration (REQUIRED — first deploy only)

```bash
# Tạo bảng database — chạy 1 lần trước khi dùng
docker compose -f docker-compose.prod.yml run --rm migrate

# Seed KHÔNG chạy tự động trong production.
# Nếu cần dữ liệu khởi tạo, nhập thủ công qua SQL.
```

### Step 6: Verify Services

```bash
# Check backend health
curl https://localhost/api/v1/health

# Check frontend
curl https://localhost

# View logs
docker compose -f docker-compose.prod.yml logs backend
docker compose -f docker-compose.prod.yml logs frontend
```

---

## ✅ Post-Deployment Validation

### Application Health Checks

```bash
# Backend health endpoint
curl -k https://api.yourdomain.com/api/v1/health

# Expected response:
# {
#   "status": "ok",
#   "message": "HRM API đang chạy!"
# }

# Frontend check
curl -k https://yourdomain.com

# Should return HTML content
```

### SSL/TLS Verification

```bash
# Check SSL certificate from client perspective
openssl s_client -connect yourdomain.com:443

# Test SSL grade
# Visit: https://www.ssllabs.com/ssltest/

# Check security headers
curl -i -k https://yourdomain.com | grep -i "strict-transport-security\|x-content-type-options\|x-frame-options"
```

### Database Verification

```bash
# Connect to database container
docker compose -f docker-compose.prod.yml exec mysql mysql -u hrm_app_user -p${DB_PASSWORD} hrm_db

# List tables
SHOW TABLES;

# Check user count
SELECT COUNT(*) FROM users;

# Exit
EXIT;
```

### Redis Verification

```bash
# Connect to Redis
docker compose -f docker-compose.prod.yml exec redis redis-cli

# Test connectivity
PING

# Check memory usage
INFO memory

# Exit
EXIT
```

---

## 📊 Monitoring & Maintenance

### Daily Checks

```bash
# Check container status
docker compose -f docker-compose.prod.yml ps

# Check system resource usage
docker stats

# View recent logs
docker compose -f docker-compose.prod.yml logs --tail=100

# Check disk space
df -h
```

### Weekly Tasks

- [ ] Review application logs for errors
- [ ] Monitor database growth
- [ ] Verify backup procedures running
- [ ] Test SSL certificate renewal (if using Let's Encrypt)
- [ ] Review security audit logs

### Monthly Tasks

- [ ] Performance review
- [ ] Security updates check
- [ ] Database maintenance (optimization)
- [ ] Capacity planning review
- [ ] Disaster recovery testing

### SSL Certificate Renewal (Let's Encrypt)

```bash
# Automatic renewal (runs daily via cron)
sudo certbot renew

# Manual renewal
sudo certbot renew --force-renewal

# After renewal, copy to certs directory
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem certs/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem certs/key.pem

# Reload nginx
docker compose -f docker-compose.prod.yml exec frontend nginx -s reload
```

### Backup Strategy

```bash
# Daily database backup script
#!/bin/bash
docker compose -f docker-compose.prod.yml exec mysql mysqldump \
  -u hrm_app_user \
  -p"${DB_PASSWORD}" \
  hrm_db | gzip > backups/hrm_db_$(date +%Y%m%d_%H%M%S).sql.gz

# Add to crontab (daily at 2 AM)
# 0 2 * * * /path/to/backup-script.sh
```

---

## 🔧 Troubleshooting

### Container Issues

```bash
# Container won't start - check logs
docker compose -f docker-compose.prod.yml logs <service-name>

# Restart container
docker compose -f docker-compose.prod.yml restart <service-name>

# View container status details
docker inspect <container-name>
```

### SSL/TLS Issues

```bash
# Certificate mismatch error
# → Check that certificate CN matches your domain name

# Connection refused on port 443
# → Verify nginx is listening on 443
# → Check certificate files exist and are readable
# → Check firewall rules allow port 443

# Mixed content warning (HTTP content on HTTPS page)
# → Verify nginx redirect HTTP to HTTPS
# → Check proxy headers (X-Forwarded-Proto)
```

### Database Issues

```bash
# Database connection failed
docker compose -f docker-compose.prod.yml logs mysql

# Check database credentials in .env match
grep "DB_" .env

# MySQL port conflict
docker compose -f docker-compose.prod.yml ps | grep mysql
netstat -tlnp | grep 3306
```

### Performance Issues

```bash
# Check container resource usage
docker stats

# Check database query performance
docker compose -f docker-compose.prod.yml exec mysql mysql \
  -u hrm_app_user -p"${DB_PASSWORD}" hrm_db \
  -e "SHOW PROCESSLIST;"

# Check Redis memory
docker compose -f docker-compose.prod.yml exec redis redis-cli INFO memory
```

---

## 🚨 Emergency Rollback

If deployment fails or issues occur:

```bash
# Stop current deployment
docker compose -f docker-compose.prod.yml down

# Restore database from backup
gunzip < backups/hrm_db_<timestamp>.sql.gz | \
  docker compose -f docker-compose.prod.yml exec -T mysql mysql \
  -u root -p"${MYSQL_ROOT_PASSWORD}" hrm_db

# Deploy previous version
docker compose -f docker-compose.prod.yml up -d
```

---

## 📚 Security Checklist (OWASP)

- [ ] **A01:2021 – Broken Access Control** → JWT authentication enforced
- [ ] **A02:2021 – Cryptographic Failures** → HTTPS/TLS enabled, secrets in .env
- [ ] **A03:2021 – Injection** → Input validation in backend
- [ ] **A04:2021 – Insecure Design** → Architecture reviewed
- [ ] **A05:2021 – Security Misconfiguration** → Security headers set
- [ ] **A06:2021 – Vulnerable & Outdated Components** → Docker images updated regularly
- [ ] **A07:2021 – Authentication Failures** → JWT + Refresh token implemented
- [ ] **A08:2021 – Data Integrity Failures** → Database integrity checks
- [ ] **A09:2021 – Logging & Monitoring Failures** → Logging implemented
- [ ] **A10:2021 – SSRF** → Input validation in place

---

## 📞 Support & Contact

For deployment issues or questions:

- Review logs: `docker compose -f docker-compose.prod.yml logs`
- Check health endpoint: `curl https://api.yourdomain.com/api/v1/health`
- Contact development team for debugging

---

**Note:** Always test deployment process in staging environment first!
