# Use Debian-based VS Code dev container image
FROM mcr.microsoft.com/vscode/devcontainers/base:bookworm

# Install dependencies in a single layer to reduce image size
RUN apt-get update && apt-get install -y \
    curl \
    git \
    make \
    lib32gcc-s1 \
    dbus \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.24.2
RUN wget -q https://go.dev/dl/go1.24.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz && \
    rm go1.24.2.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
RUN mkdir -p /go/bin && chown -R vscode:vscode /go


# Set up non-root user and workspace
USER vscode
WORKDIR /workspaces