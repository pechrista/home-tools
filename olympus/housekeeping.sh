#!/bin/bash
# 
# To add this script to sudo cron:
# 1. Run 'sudo crontab -e'
# 2. Add this entry: '0 3 * * 0 /path/to/housekeeping.sh'

# Update package lists
/usr/bin/apt update -y

# Upgrade all installed packages
/usr/bin/apt upgrade -y

# Remove obsolete packages
/usr/bin/apt autoremove -y

# Clean the local repository of retrieved package files
/usr/bin/apt autoclean -y

# Log the update time (optional, for debugging/monitoring)
echo "Weekly update and reboot completed on $(date)" >> /var/log/weekly-update-reboot.log

# Force a reboot
/sbin/shutdown -r now
