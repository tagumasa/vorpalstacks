#!/bin/bash
set -e

SERVER_PORT=${1:-8080}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="$SCRIPT_DIR/../tmp/vorpalstacks-server-${SERVER_PORT}.pid"

echo "Starting VorpalStacks server on port $SERVER_PORT..."
echo "Press Ctrl+C to stop the server"
echo ""

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "Server is already running on port $SERVER_PORT (PID: $PID)"
        echo "Use 'scripts/stop_server.sh $SERVER_PORT' to stop it first"
        exit 1
    fi
    rm -f "$PID_FILE"
fi

if lsof -i :$SERVER_PORT | grep -q LISTEN; then
    echo "Port $SERVER_PORT is already in use"
    echo "Use 'scripts/stop_server.sh $SERVER_PORT' to stop the existing server"
    exit 1
fi

export SIGNATURE_VERIFICATION_ENABLED=${SIGNATURE_VERIFICATION_ENABLED:-false}
export PORT=$SERVER_PORT
export LOG_LEVEL=${LOG_LEVEL:-INFO}
export LOG_FILE=${LOG_FILE:-}
export LOG_MAX_SIZE=${LOG_MAX_SIZE:-100}
export LOG_MAX_AGE=${LOG_MAX_AGE:-7}
export LOG_MAX_BACKUPS=${LOG_MAX_BACKUPS:-3}
export TEST_MODE=${TEST_MODE:-false}

cd "$(dirname "$0")/.."

echo "Log level: $LOG_LEVEL"
if [ -n "$LOG_FILE" ]; then
    echo "Log file: $LOG_FILE"
fi
echo "Test mode: $TEST_MODE"
echo "Signature verification: $SIGNATURE_VERIFICATION_ENABLED"
echo ""

go run main.go
