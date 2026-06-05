# Minecraft Bedrock Docker Stack

This setup runs a Bedrock server with an automated backup companion.
Both services use `restart: unless-stopped`.

## What This Includes

1. `bedrock-vanilla`: the main Minecraft Bedrock server container.
	- Pulls the `latest` server image on startup (`pull_policy: always`).
	- Mounts persistent world/server data at `./data`.
1. `bedrock-backup`: a backup service container for scheduled world backups.
	- Pulls the `latest` backup image on startup (`pull_policy: always`).
	- Uses `BACKUP_CRON` for backup cadence (for example, daily at 4 AM).
	- Reads world data from `./data` and writes backup archives to `./backups`.

## How to Run

## Prerequisites

1. Make sure your computer has the requirements to run this server. See [this page](https://www.minecraft.net/en-us/download/server/bedrock) for more information. 
1. You need Docker and Docker Compose installed on your machine.

## Steps to Run

1. Create a folder called `minecraft-{SERVER_NAME}` and go into it.
1. Copy `docker-compose.yaml` from this repo into that folder.
1. Tweak values as needed. (ex. Port mappings, backup frequency / timezone)
1. Run `docker compose up -d`