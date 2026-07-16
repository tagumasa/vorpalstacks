#!/bin/bash
# Common configuration for VorpalStacks test scripts

# Default settings
DEFAULT_PORT=8080
DEFAULT_REGION="us-east-1"
DEFAULT_ACCOUNT_ID="123456789012"

# Override with environment variables
PORT="${PORT:-$DEFAULT_PORT}"
ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:${PORT}}"
REGION="${REGION:-$DEFAULT_REGION}"
ACCOUNT_ID="${ACCOUNT_ID:-$DEFAULT_ACCOUNT_ID}"

# Export for AWS CLI (needs AWS_DEFAULT_REGION)
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-$REGION}"

# Test data path (use project's tmp directory, not OS tmp)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
TEST_DATA_PATH="${TEST_DATA_PATH:-$PROJECT_ROOT/tmp/test-data}"

# Server settings
SERVER_LOG="$PROJECT_ROOT/tmp/vorpalstacks-test.log"
HEALTH_CHECK_TIMEOUT=30

# Default credentials for tests that need authentication context
DEFAULT_AWS_ACCESS_KEY_ID="test"
DEFAULT_AWS_SECRET_ACCESS_KEY="test"
