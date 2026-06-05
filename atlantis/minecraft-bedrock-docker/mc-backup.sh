#!/bin/sh

# 1. Use the built-in HOSTNAME env variable to get the current container's ID
MY_CONTAINER_ID=$HOSTNAME

# 2. Extract the unique Compose Project name for THIS specific stack
PROJECT_NAME=$(docker inspect --format "{{ index .Config.Labels \"com.docker.compose.project\" }}" "$MY_CONTAINER_ID")

# 3. Find the sibling game server running inside the exact same project
TARGET_CONTAINER=$(docker ps --filter "label=com.docker.compose.project=$PROJECT_NAME" --filter "label=com.docker.compose.service=bedrock-vanilla" --format "{{.ID}}")

# Setup runtime variables
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
echo "=== Starting Backup: $(date) ==="

if [ -z "$TARGET_CONTAINER" ]; then
    echo "Error: Could not find sibling bedrock-vanilla container in project $PROJECT_NAME!"
    exit 1
fi

echo "[Project: $PROJECT_NAME] Found target container: $TARGET_CONTAINER. Stopping safely..."
docker stop "$TARGET_CONTAINER"

if [ -d "/bedrock_data/worlds" ]; then
    tar -czf "/backups/worlds_backup_${TIMESTAMP}.tar.gz" -C /bedrock_data worlds
    echo "Backup saved: worlds_backup_${TIMESTAMP}.tar.gz"
fi

# Strict rotational cleanup based on total file count
echo "Checking retention limit (Max: $MAX_BACKUPS)..."
CURRENT_COUNT=$(ls -1 /backups/worlds_backup_*.tar.gz 2>/dev/null | wc -l)

if [ "$CURRENT_COUNT" -gt "$MAX_BACKUPS" ]; then
    EXCESS=$(($CURRENT_COUNT - $MAX_BACKUPS))
    echo "Found $CURRENT_COUNT backups. Removing $EXCESS oldest backup(s)...\""
    ls -1tr /backups/worlds_backup_*.tar.gz | head -n "$EXCESS" | xargs rm -f
fi

echo "Restarting target container..."
docker start "$TARGET_CONTAINER"
echo "=== Backup Complete ==="

