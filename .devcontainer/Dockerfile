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

# Install Node.js 22.15.0
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest

# Install Wine for Windows .exe builds
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y wine wine32 wine64 && \
    rm -rf /var/lib/apt/lists/*

# Install NSIS for Windows NSIS installer
RUN wget -q https://downloads.sourceforge.net/project/nsis/NSIS%203/3.09/nsis-3.09-setup.exe && \
    wine nsis-3.09-setup.exe /S && \
    rm nsis-3.09-setup.exe

# Start D-Bus to avoid Electron errors
RUN service dbus start

# Set up non-root user and workspace
USER vscode
WORKDIR /workspaces