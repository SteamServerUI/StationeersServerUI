FROM debian:12-slim AS runner

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application source code
COPY ./build/StationeersServerControl*.x86_64 /app/StationeersServerControl
COPY ./LICENSE /app/LICENSE
RUN chmod +x /app/StationeersServerControl

RUN dpkg --add-architecture i386 \
 && apt-get update -y \
 && apt-get install -y --no-install-recommends ca-certificates locales lib32gcc-s1 file

# Verify that the executable was copied and renamed successfully
RUN echo "Verifying the copied and renamed StationeersServerControl executable:" && \
    if ls -l /app/StationeersServerControl; then \
        echo "StationeersServerControl copy and rename successful."; \
    else \
        echo "Error: StationeersServerControl executable not found after copy."; \
        exit 1; \
    fi

# Expose the ports
EXPOSE 8443 27016 27015

# Set the entrypoint to the application
ENTRYPOINT ["/app/StationeersServerControl"]

# Provide default arguments to the entrypoint
CMD []