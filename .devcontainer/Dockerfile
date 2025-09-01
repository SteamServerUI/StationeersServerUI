# Use Debian-based VS Code dev container image
FROM mcr.microsoft.com/vscode/devcontainers/base:bookworm

ARG TZ
ENV TZ="$TZ"

# Install dependencies in a single layer to reduce image size
RUN apt-get update && apt-get install -y \
    curl \
    git \
    make \
    lib32gcc-s1 \
    dbus \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.25.0
RUN wget -q https://go.dev/dl/go1.25.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz && \
    rm go1.25.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
RUN mkdir -p /go/bin && chown -R vscode:vscode /go

# Install Node.js 22.15.0
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest

# Set up non-root user and workspace
USER vscode
WORKDIR /workspaces