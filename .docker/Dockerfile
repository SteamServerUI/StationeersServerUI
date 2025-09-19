###############
# Frontend build
###############
FROM node:24-bookworm-slim AS frontend-builder

WORKDIR /src/frontend

# Only copy frontend to leverage Docker layer caching for npm install
COPY ./frontend/package.json ./frontend/package-lock.json* ./
RUN npm install

COPY ./frontend/ ./
# Build Svelte UI into ../UIMod/onboard_bundled/v2 (as configured in vite.config.js)
RUN npm run build

############
# Go build
############
FROM golang:1.25-bookworm AS go-builder

WORKDIR /src

# Pre-copy go.mod and go.sum for better caching
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . ./

COPY --from=frontend-builder /src/UIMod /src/UIMod

# Build the server (embed will include UIMod/* at build time)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/StationeersServerControl ./server.go

############
# Runtime
############
FROM debian:13-slim AS runner

WORKDIR /app

ENV HOME=/app

# Minimal runtime dependencies. SteamCMD needs 32-bit runtime.
RUN dpkg --add-architecture i386 \
 && apt-get update -y \
 && apt-get install -y --no-install-recommends \
    ca-certificates \
    locales \
    lib32gcc-s1 \
    file \
    gosu \
 && rm -rf /var/lib/apt/lists/*

RUN useradd -m -d /app -s /bin/sh app

# Copy compiled binary and license
COPY --from=go-builder /out/StationeersServerControl /usr/local/bin/StationeersServerControl
COPY ./LICENSE /app/LICENSE
COPY ./.docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/StationeersServerControl /usr/local/bin/docker-entrypoint.sh

RUN chown -R app:app /app

EXPOSE 8443/tcp 27016/udp 27015/udp

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
