# syntax=docker/dockerfile:1.7
# ----------------------------------------------------------------------------
# reputation-be container image
# ----------------------------------------------------------------------------
#   builder  -> pinned golang alpine, runs `make build` (which `setup`s config
#               and compiles the cobra CLI as `main`)
#   runtime  -> minimal alpine, non-root user, only the binary + config +
#               db migrations shipped in.
#   entrypoint: `./main start` (cobra subcommand) listening on :4011
# ----------------------------------------------------------------------------

ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.25.9

# ---- builder ----
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Build-time tools: ca-certificates for outbound TLS, gcc/musl for any cgo,
# make for the `make build` target.
RUN apk add --no-cache make ca-certificates gcc g++ libc-dev git

WORKDIR /app

# Cache go.mod/go.sum first so dependency resolution doesn't bust the layer
# on every source change.
COPY go.mod go.sum ./
RUN go mod download

# Then the rest of the source.
COPY . .

# `make build` runs `setup` (copies config/example -> config/resources) and
# then `go build -o ./main main.go`. Reproducible by virtue of go build's
# default ldflags behaviour.
RUN make build

# ---- runtime ----
# Pin to the same alpine major/minor so `golang:` and `alpine:` stay in sync.
FROM alpine:${ALPINE_VERSION}

# Only what's needed at runtime: tzdata + ca-certificates (for outbound TLS
# to Groq / MiniMax / etc.). bash stays in for `start.sh`-like ops.
RUN apk add --no-cache ca-certificates bash tzdata

WORKDIR /app

# Non-root user (matches the FE runner pattern). uid 1001 / gid 1001.
RUN addgroup -S -g 1001 app \
 && adduser  -S -u 1001 -G app -h /app app

# Layout matches what `make build` produced.
RUN mkdir -p /app/config /app/db /app/logs \
 && chown -R app:app /app

COPY --from=builder --chown=app:app /app/main         /app/main
COPY --from=builder --chown=app:app /app/config        /app/config
COPY --from=builder --chown=app:app /app/db            /app/db
COPY --from=builder --chown=app:app /app/Makefile      /app/Makefile

USER app

EXPOSE 4011

ENV PATH=/app:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin \
    GIN_MODE=release \
    TZ=UTC

# The CLI binary accepts subcommands: `start` boots the REST server on :4011.
ENTRYPOINT ["./main"]
CMD ["start"]


# ----------------------------------------------------------------------------
# OCI image metadata
# ----------------------------------------------------------------------------
LABEL org.opencontainers.image.title="reputation-be" \
      org.opencontainers.image.description="CekReputasi — Go/Fiber REST API + admin/customer/auth modules" \
      org.opencontainers.image.source="https://github.com/raymondsugiarto/reputation-be" \
      org.opencontainers.image.licenses="Proprietary"
