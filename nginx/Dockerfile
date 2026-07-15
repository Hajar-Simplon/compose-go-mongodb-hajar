# ======================================================
# Nginx Reverse Proxy
# ======================================================

# Use a specific version of the official Nginx image.
FROM nginx:1.28-alpine

# Copy our custom Nginx configuration
# into the default configuration directory.
COPY default.conf /etc/nginx/conf.d/default.conf# ======================================================
# Nginx Reverse Proxy
# ======================================================

FROM nginx:1.28-alpine

# Copy Nginx configuration
COPY default.conf /etc/nginx/conf.d/default.conf

# Copy TLS certificate and private key
COPY certs/server.crt /etc/nginx/certs/server.crt
COPY certs/server.key /etc/nginx/certs/server.key