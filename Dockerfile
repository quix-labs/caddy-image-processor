FROM caddy:builder-alpine AS builder

LABEL org.opencontainers.image.source="https://github.com/quix-labs/caddy-image-processor"


RUN apk add --update --no-cache make vips-dev gcc musl-dev

ADD . .
RUN make build

# FINAL IMAGE
FROM caddy:alpine
RUN apk add --update --no-cache vips

COPY --from=builder /usr/bin/out/caddy /usr/bin/caddy