# First stage: Build the Go application using the Golang image
FROM golang:1.22.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN echo "Downloading dependencies..." && go mod download

# Copy the rest of the application source code
COPY . .

# Build the initial executable
RUN echo "Building StationeersServerUI..." && go build -o StationeersServerUI ./build.go

#TODO: Add cross compilation during build and split of into two images, one for windows and one for linux

# Second stage: Bootstrap the server using a Debian slim image
FROM golang:1.22.1 AS bootstrapper

# Set the working directory inside the container
WORKDIR /app

# Install required libraries
RUN echo "Installing required libraries..." && apt-get update && apt-get install -y \
    lib32gcc-s1 \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

# Copy the initial executable from the builder stage
COPY --from=builder /app/StationeersServerUI /app/StationeersServerUI

# Copy the rest of the application source code
COPY --from=builder /app /app

# Download and install SteamCMD
#RUN wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz && \
#    tar -xvzf steamcmd_linux.tar.gz && \
#    rm steamcmd_linux.tar.gz && \
#    chmod +x steamcmd.sh

# Run the initial executable to build StationeersServerControl
RUN echo "Running StationeersServerUI to build StationeersServerControl..." && ./StationeersServerUI

# Verify that the resulting executable exists
RUN echo "Verifying the existence of StationeersServerControl executable:" && \
    if ls -l /app/StationeersServerControl*; then \
        echo "StationeersServerControl build successful."; \
    else \
        echo "Error: StationeersServerControl executable not found."; \
        exit 1; \
    fi

# Copy the UIMod directory
COPY --from=builder /app/UIMod /app/UIMod

# Third stage: Run the final application using the steamcmd/steamcmd image
FROM debian:12-slim AS runner

# Set the working directory inside the container
WORKDIR /app

RUN dpkg --add-architecture i386 \
 && apt-get update -y \
 && apt-get install -y --no-install-recommends ca-certificates locales lib32gcc-s1

# Install required libraries
#RUN echo "Installing required libraries..." && apt-get update && apt-get install -y \
#    lib32gcc-s1 \
#    libc6 \
#    && rm -rf /var/lib/apt/lists/*

# Install required libraries
#RUN set -x \
#    apt-get update \
#    && apt-get upgrade -yq \
#    && apt-get install -y \
#    lib32gcc-s1 \
#    libc6 \
#    wget \
#    tar \
#    gzip \
#    && rm -rf /var/lib/apt/lists/*

# Copy the resulting executable from the bootstrapper stage and rename it
COPY --from=bootstrapper /app/StationeersServerControl* /app/StationeersServerControl

# Verify that the executable was copied and renamed successfully
RUN echo "Verifying the copied and renamed StationeersServerControl executable:" && \
    if ls -l /app/StationeersServerControl; then \
        echo "StationeersServerControl copy and rename successful."; \
    else \
        echo "Error: StationeersServerControl executable not found after copy."; \
        exit 1; \
    fi

# Copy the UIMod directory
COPY --from=bootstrapper /app/UIMod /app/UIMod

# Expose the ports
EXPOSE 8080 27016

# Set the entrypoint to the application
ENTRYPOINT ["/app/StationeersServerControl"]

# Provide default arguments to the entrypoint
CMD []

#TODO: Fix image to launch rocketstation_DedicatedServer.exe/.x84_64 with the right mono dependencies