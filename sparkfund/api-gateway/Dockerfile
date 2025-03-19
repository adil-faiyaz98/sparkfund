FROM nginx:alpine

# Install required modules
RUN apk add --no-cache \
    nginx-mod-http-lua \
    nginx-mod-http-echo \
    nginx-mod-http-upstream-fair \
    nginx-mod-http-cache-purge

# Copy custom configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create cache directory
RUN mkdir -p /tmp/nginx_cache && \
    chown -R nginx:nginx /tmp/nginx_cache

# Expose ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1
