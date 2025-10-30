#!/bin/bash

# Script to set up hourly cron job for code stats

SCRIPT_PATH="/main/Project-Website/scripts/monitoring/update_code_stats.sh"
CRON_LOG_FILE="/tmp/update_code_stats.log"

echo "Setting up hourly cron job for code statistics..."

# Check if the update script exists
if [ ! -f "$SCRIPT_PATH" ]; then
    echo "Error: Code stats script not found at $SCRIPT_PATH"
    exit 1
fi

# Make sure the script is executable
chmod +x "$SCRIPT_PATH"

# Check if cron job already exists
if crontab -l 2>/dev/null | grep -q "$SCRIPT_PATH"; then
    echo "Cron job already exists for code stats update"
    echo "Current cron entry:"
    crontab -l | grep "$SCRIPT_PATH"
else
    # Add the cron job
    (crontab -l 2>/dev/null; echo "0 * * * * $SCRIPT_PATH > $CRON_LOG_FILE 2>&1") | crontab -
    echo "✅ Cron job added successfully!"
    echo "Code stats will update every hour at :00"
    echo "Logs will be written to: $CRON_LOG_FILE"
fi

echo ""
echo "To view current cron jobs: crontab -l"
echo "To remove the cron job: crontab -e (then delete the line)"
echo "To check the last run: cat $CRON_LOG_FILE"