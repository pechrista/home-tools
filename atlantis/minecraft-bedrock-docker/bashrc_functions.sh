# Dynamic Minecraft Docker Functions
function mc_start() {
  if [ -f ./docker-compose.yaml ] || [ -f ./docker-compose.yml ]; then
    docker compose up --detach && docker compose logs -f
  else
    echo "Error: No docker-compose file found in the current directory."
  fi
}

function mc_stop() {
  if [ -f ./docker-compose.yaml ] || [ -f ./docker-compose.yml ]; then
    docker compose down
  else
    echo "Error: No docker-compose file found in the current directory."
  fi
}

function mc_restart() {
  if [ -f ./docker-compose.yaml ] || [ -f ./docker-compose.yml ]; then
    docker compose restart
  else
    echo "Error: No docker-compose file found in the current directory."
  fi
}

function mc_console() {
  if [ -f ./docker-compose.yaml ] || [ -f ./docker-compose.yml ]; then
    # Dynamically fetches the name of the first running service in this compose file
    local container_name=$(docker compose ps --format "{{.Name}}" | head -n 1)
    if [ -n "$container_name" ]; then
      # Do not forward Ctrl+C (SIGINT) to the container process.
      docker container attach --sig-proxy=false "$container_name"
    else
      echo "Error: The server is not currently running."
    fi
  else
    echo "Error: No docker-compose file found in the current directory."
  fi
}

function mc_logs() {
  if [ -f ./docker-compose.yaml ] || [ -f ./docker-compose.yml ]; then
    docker compose logs -f
  else
    echo "Error: No docker-compose file found in the current directory."
  fi
}

# ==============================================================================
# LOGIN DASHBOARD: MINECRAFT CONTEXT & SECURITY SNAPSHOT
# ==============================================================================

# 1. Minecraft Functions Context (Clean List)
echo -e "\n\033[1;32m🎮 Docker Minecraft Functions Loaded (Context-Aware)\033[0m"
echo -e "Navigate to your server directory and run:"
echo -e "  • \033[1;36mmc_start\033[0m   \033[1;30m-\033[0m Boot the instance"
echo -e "  • \033[1;36mmc_stop\033[0m    \033[1;30m-\033[0m Safely halt the server"
echo -e "  • \033[1;36mmc_restart\033[0m \033[1;30m-\033[0m Quick recycle"
echo -e "  • \033[1;36mmc_console\033[0m \033[1;30m-\033[0m Attach to interactive CLI"
echo -e "  • \033[1;36mmc_logs\033[0m    \033[1;30m-\033[0m Follow runtime stdout"