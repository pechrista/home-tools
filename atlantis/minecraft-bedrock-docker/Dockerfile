FROM ubuntu:22.04

# Install standard dependencies and jq
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    libcurl4 \
    libssl3 \
    ca-certificates \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Set up active runtime directory
WORKDIR /bedrock
ENV LD_LIBRARY_PATH=.
EXPOSE 19132/udp

# Copy the startup script into the image
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Run it
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]