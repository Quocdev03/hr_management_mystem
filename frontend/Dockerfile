# Stage 1: Build static assets
FROM node:20-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Set Vite build environment variables
ENV VITE_API_URL=/api/v1

# Build application
RUN npm run build

# Stage 2: Serve using Nginx
FROM nginx:stable-alpine

# Create directories for SSL certificates and logs
RUN mkdir -p /etc/nginx/certs /var/log/nginx

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build static assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose HTTP and HTTPS ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

# Set proper permissions
RUN chmod -R 755 /usr/share/nginx/html /var/log/nginx

CMD ["nginx", "-g", "daemon off;"]
