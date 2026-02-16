#!/bin/bash
#
# Setup script to add housekeeping tasks to crontab
# Run with: curl -fsSL https://raw.githubusercontent.com/pechrista/home-tools/main/setup-housekeeping.sh | bash
# Or download and run: bash setup-housekeeping.sh
#
# This script will:
# 1. Create the housekeeping.sh script
# 2. Add it to sudo crontab to run weekly (Sundays at 3 AM)

# Check if running with sudo
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (sudo)"
   exit 1
fi

# Define the housekeeping script path
HOUSEKEEPING_SCRIPT="/usr/local/bin/housekeeping.sh"

# Create the housekeeping script
cat > "$HOUSEKEEPING_SCRIPT" << 'EOF'
#!/bin/bash
# 
# Weekly housekeeping script for system maintenance
# Scheduled to run Sundays at 3 AM via crontab

# Update package lists
/usr/bin/apt update -y

# Upgrade all installed packages
/usr/bin/apt upgrade -y

# Remove obsolete packages
/usr/bin/apt autoremove -y

# Clean the local repository of retrieved package files
/usr/bin/apt autoclean -y

# Log the update time (optional, for debugging/monitoring)
echo "Weekly housekeeping completed on $(date)" >> /var/log/housekeeping.log
EOF

# Make the housekeeping script executable
chmod +x "$HOUSEKEEPING_SCRIPT"

# Add to sudo crontab
echo "Adding housekeeping to sudo crontab..."
(sudo crontab -l 2>/dev/null | grep -v "$HOUSEKEEPING_SCRIPT"; echo "0 3 * * 0 $HOUSEKEEPING_SCRIPT") | sudo crontab -

echo "✓ Housekeeping script installed at $HOUSEKEEPING_SCRIPT"
echo "✓ Cron job scheduled to run weekly on Sundays at 3 AM"
echo ""
echo "To view your sudo crontab, run: sudo crontab -l"
echo "To remove this job, run: sudo crontab -e and delete the housekeeping line"