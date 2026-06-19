# ═══════════════════════════════════════════════════════════════════════════
# Frontend Dockerfile — Production Build
# ═══════════════════════════════════════════════════════════════════════════
# Stage 1: Build static assets với Vite
# Stage 2: Serve với Nginx (production config)
# ═══════════════════════════════════════════════════════════════════════════

# Stage 1: Build static assets
FROM node:20-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# VITE_API_URL cho production: dùng relative path qua nginx reverse proxy
# Nginx sẽ forward /api/v1/* → backend:8080
ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=${VITE_API_URL}

# Build production bundle
RUN npm run build

# Stage 2: Serve with Nginx (production config with SSL support)
FROM nginx:stable-alpine

# Create directories for SSL certificates and logs
RUN mkdir -p /etc/nginx/certs /var/log/nginx

# Copy production nginx config (HTTPS + security headers)
COPY nginx.prod.conf /etc/nginx/conf.d/default.conf

# Copy build static assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Set proper permissions
RUN chmod -R 755 /usr/share/nginx/html /var/log/nginx

# Expose HTTP and HTTPS ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

CMD ["nginx", "-g", "daemon off;"]
