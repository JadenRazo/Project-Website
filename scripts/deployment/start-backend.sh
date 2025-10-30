#!/bin/bash
# This script is a wrapper around the new start-prod.sh script.
# Please use start-prod.sh in the future.

echo "This script is a wrapper around the new start-prod.sh script."
echo "Please use start-prod.sh in the future."

"$(dirname "${BASH_SOURCE[0]}")/../start-prod.sh"