#!/bin/sh
set -e

MARKER_FILE="/bedrock/.current_version"

echo 'Checking Mojang API for the latest Bedrock server version...'
API_RESPONSE=$(curl -s https://net-secondary.web.minecraft-services.net/api/v1.0/download/links)
DOWNLOAD_URL=$(echo "$API_RESPONSE" | jq -r '.result.links[] | select(.downloadType=="serverBedrockLinux") | .downloadUrl')

if [ -z "$DOWNLOAD_URL" ] || [ "$DOWNLOAD_URL" = "null" ]; then
    echo "Warning: Could not fetch download URL from API. Falling back to existing local installation..."
else
    # Extract a unique identifier from the URL (Mojang URLs usually contain the version number)
    # Example: https://.../bedrock-server-1.21.2.02.zip -> 1.21.2.02
    LATEST_VERSION=$(echo "$DOWNLOAD_URL" | grep -oE 'bedrock-server-[0-9.]+' | sed 's/bedrock-server-//' || echo "$DOWNLOAD_URL")

    # Read the previously installed version if it exists
    INSTALLED_VERSION=""
    if [ -f "$MARKER_FILE" ]; then
        INSTALLED_VERSION=$(cat "$MARKER_FILE")
    fi

    if [ "$LATEST_VERSION" = "$INSTALLED_VERSION" ] && [ -f "/bedrock/bedrock_server" ]; then
        echo "Server is up to date (Version: $LATEST_VERSION). Skipping download."
    else
        echo "New version detected or fresh install required!"
        echo "Installed: [${INSTALLED_VERSION:-None}] -> Latest: [$LATEST_VERSION]"
        echo "Downloading version from: $DOWNLOAD_URL"

        mkdir -p /tmp/bedrock_download
        cd /tmp/bedrock_download

        if curl -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" -fsSL "$DOWNLOAD_URL" -o bedrock-server.zip; then
            unzip -q bedrock-server.zip
            rm bedrock-server.zip

            echo 'Syncing core executable and system files to persistent storage...'
            cp -f bedrock_server /bedrock/
            cp -f *.so /bedrock/ 2>/dev/null || true
            cp -rf behavior_packs /bedrock/
            cp -rf resource_packs /bedrock/
            cp -f *_packs.json /bedrock/ 2>/dev/null || true

            if [ ! -f /bedrock/server.properties ]; then cp server.properties /bedrock/; fi
            if [ ! -f /bedrock/allowlist.json ]; then cp allowlist.json /bedrock/; fi
            if [ ! -f /bedrock/permissions.json ]; then cp permissions.json /bedrock/; fi

            # Save the version token to your persistent directory so it survives restarts
            echo "$LATEST_VERSION" > "$MARKER_FILE"
        else
            echo "Error: Download failed. Sticking with existing installation."
        fi

        cd /bedrock
        rm -rf /tmp/bedrock_download
    fi
fi

# Launch the server
cd /bedrock
chmod +x bedrock_server
exec ./bedrock_server