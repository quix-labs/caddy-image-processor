FROM caddy:2.8-builder-alpine AS builder

RUN apk add --update --no-cache make vips-dev gcc musl-dev

ADD . .
RUN make build

# FINAL IMAGE
FROM caddy:2.8-alpine
RUN apk add --update --no-cache vips

COPY --from=builder /usr/bin/out/caddy /usr/bin/caddy