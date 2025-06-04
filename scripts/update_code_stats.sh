#!/bin/bash

# --- Configuration ---
# IMPORTANT: List FULL SYSTEM ABSOLUTE paths for tokei to scan.
# Example: PROJECT_DIRECTORIES=("/home/user/projects/quiz_bot" "/var/www/main/Project-Website")
# Make sure the user running this script has read access to these directories.
PROJECT_DIRECTORIES=( \
  "/quiz_bot" \
  "/main/Project-Website" \
  # Add more FULL ABSOLUTE project paths here, each on a new line ending with \
)

# Directories to exclude
EXCLUDE_DIRS=("build" "node_modules" "logs" "bin")

# Path to the output JSON file, relative to the determined workspace root.
OUTPUT_JSON_FILE_RELATIVE="main/Project-Website/frontend/public/code_stats.json"

# --- Helper Functions ---

# Get the absolute path of the script itself
get_script_abs_path() {
  local script_path
  script_path=$(readlink -f "$0" 2>/dev/null)
  if [ -z "$script_path" ]; then # Fallback
    local script_dir
    script_dir=$(cd "$(dirname "$0")" && pwd)
    script_path="$script_dir/$(basename "$0")"
  fi
  echo "$script_path"
}

# --- Global Variables (derived after functions) ---
SCRIPT_ABS_PATH=$(get_script_abs_path)
SCRIPT_DIR=$(dirname "$SCRIPT_ABS_PATH")
# Assuming the script is always at <workspace_root>/main/Project-Website/scripts/update_code_stats.sh
# Therefore, workspace_root is 3 levels above SCRIPT_DIR
# This WORKSPACE_ROOT is primarily for locating the OUTPUT_JSON_FILE_RELATIVE correctly.
WORKSPACE_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)

# Absolute path for the output JSON file
OUTPUT_JSON_FILE="$WORKSPACE_ROOT/$OUTPUT_JSON_FILE_RELATIVE"

# Function to check and guide cron job setup
setup_cron_job() {
  CRON_LOG_FILE="/tmp/update_code_stats.log"
  CRON_JOB_COMMAND="0 * * * * \"$SCRIPT_ABS_PATH\" > \"$CRON_LOG_FILE\" 2>&1"

  echo "-----------------------------------------------------"
  echo "Cron Job Setup for Automatic Code Statistics Update"
  echo "-----------------------------------------------------"
  echo "This script should run hourly to keep stats fresh."
  echo "Script absolute path: $SCRIPT_ABS_PATH"
  echo "Workspace root (for output file path): $WORKSPACE_ROOT"
  echo "Checking your crontab for an existing job for this script..."

  if crontab -l 2>/dev/null | grep -Fq "$SCRIPT_ABS_PATH"; then
    echo "Found an existing cron job for this script. No action needed."
    crontab -l 2>/dev/null | grep -F "$SCRIPT_ABS_PATH"
  else
    echo "No existing cron job found."
    echo "To set it up, please run 'crontab -e' and add the following line:"
    echo ""
    echo "    $CRON_JOB_COMMAND"
    echo ""
    echo "Output/errors will be logged to: $CRON_LOG_FILE"
  fi
  echo "-----------------------------------------------------"
}

# --- Main Script Logic ---

echo "Starting tokei code statistics generation..."
echo "Script location: $SCRIPT_ABS_PATH"
echo "Workspace root (for output file): $WORKSPACE_ROOT"
echo "Output JSON file will be: $OUTPUT_JSON_FILE"

# Ensure the target directory for the JSON file exists
echo "Ensuring directory exists: $(dirname "$OUTPUT_JSON_FILE")"
mkdir -p "$(dirname "$OUTPUT_JSON_FILE")"

# Construct the tokei command arguments from configured ABSOLUTE paths
TOKEI_ARGS=()
for dir_path_input in "${PROJECT_DIRECTORIES[@]}"; do
  dir_to_check="$dir_path_input"
  if [[ "$dir_to_check" != /* ]]; then
    echo "Error: Path '$dir_to_check' in PROJECT_DIRECTORIES is not an absolute path. Please fix." >&2
    echo '{"totalLines": null, "error": "Invalid configuration: Not an absolute path in PROJECT_DIRECTORIES."}' > "$OUTPUT_JSON_FILE"
    exit 1
  fi

  if [ -d "$dir_to_check" ]; then
    TOKEI_ARGS+=("$dir_to_check")
  else
    echo "Warning: Directory '$dir_to_check' (from input '$dir_path_input') not found. Skipping."
  fi
done

if [ ${#TOKEI_ARGS[@]} -eq 0 ]; then
  echo "Error: No valid project directories found to scan. Exiting."
  echo '{"totalLines": null, "error": "No valid project directories specified or found."}' > "$OUTPUT_JSON_FILE"
  exit 1
fi

echo "Running tokei with the following arguments: ${TOKEI_ARGS[*]}"

# Use the full path to tokei
TOKEI_CMD="/root/.cargo/bin/tokei"
if [ ! -x "$TOKEI_CMD" ]; then
  echo "Error: Tokei not found at $TOKEI_CMD" >&2
  echo '{"totalLines": null, "error": "Tokei not found at expected location"}' > "$OUTPUT_JSON_FILE"
  exit 1
fi

# Build exclude arguments
EXCLUDE_ARGS=""
for exclude_dir in "${EXCLUDE_DIRS[@]}"; do
  EXCLUDE_ARGS="$EXCLUDE_ARGS --exclude $exclude_dir"
done

RAW_TOKEI_OUTPUT=$($TOKEI_CMD $EXCLUDE_ARGS "${TOKEI_ARGS[@]}" 2>&1)
TOKEI_EXIT_CODE=$?

echo "---- Raw Tokei Output (Exit Code: $TOKEI_EXIT_CODE) ----"
echo "$RAW_TOKEI_OUTPUT"
echo "---------------------------------------------------"

if [ $TOKEI_EXIT_CODE -ne 0 ]; then
    echo "Error: Tokei command failed with exit code $TOKEI_EXIT_CODE." >&2
    echo "Please check the Raw Tokei Output above for details."
    echo "{'totalLines': null, 'error': 'Tokei command failed. Exit code: $TOKEI_EXIT_CODE'}" > "$OUTPUT_JSON_FILE"
    exit 1
fi

# Parse the raw output for the FINAL "Total" line and extract the 3rd field (Lines).
# Grep for lines starting with optional spaces then "Total ", take the last one, then get 3rd field.
PARSED_TOTAL_LINES=$(echo "$RAW_TOKEI_OUTPUT" | grep -E '^[[:space:]]*Total[[:space:]]+' | tail -n 1 | awk '{print $3}')

if [[ "$PARSED_TOTAL_LINES" =~ ^[0-9]+$ ]]; then
  echo "Total lines (as per 'Lines' column from Tokei summary) extracted: $PARSED_TOTAL_LINES"
  # The frontend expects a field named "totalLines"
  JSON_CONTENT="{\"totalLines\": $PARSED_TOTAL_LINES}"
  echo "$JSON_CONTENT" > "$OUTPUT_JSON_FILE"
  echo "Successfully wrote statistics to $OUTPUT_JSON_FILE"
else
  echo "Error: Failed to parse total lines (Lines column) from tokei output."
  echo "Check the Raw Tokei Output. The 'Total' line or its 'Lines' column might be missing or malformed."
  echo '{"totalLines": null, "error": "Failed to parse tokei output (Lines column) or no code found."}' > "$OUTPUT_JSON_FILE"
  exit 1
fi

if [ "$1" == "--check-cron" ] || [ "$1" == "--setup-cron" ]; then
    setup_cron_job
elif [ -z "$1" ] && [ -t 1 ]; then
    setup_cron_job
fi

exit 0 