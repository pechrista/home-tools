#!/bin/bash

# Create a script that adds housekeeping tasks to crontab.

# Create a temporary file for the crontab
CRON_FILE=$(mktemp)

# Get the current crontab
crontab -l > "$CRON_FILE"

# Add housekeeping tasks to the crontab
# Example: run housekeeping.sh every day at 2am
echo "0 2 * * * /path/to/housekeeping.sh" >> "$CRON_FILE"

# Install the new crontab
crontab "$CRON_FILE"

# Cleanup
rm "$CRON_FILE"

# Inform the user
echo "Housekeeping tasks added to crontab."