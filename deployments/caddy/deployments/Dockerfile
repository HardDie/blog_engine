FROM caddy:2.7.5-builder-alpine as builder

RUN xcaddy build \
    --with github.com/caddy-dns/cloudflare

# Build final image
FROM caddy:2.7.5-alpine
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
