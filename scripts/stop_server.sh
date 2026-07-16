#!/bin/bash
set -e

SERVER_PORT=${1:-8080}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="$SCRIPT_DIR/../tmp/vorpalstacks-server-${SERVER_PORT}.pid"

echo "Stopping VorpalStacks server on port $SERVER_PORT..."

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "Stopping server (PID: $PID)..."
        kill -TERM "$PID"
        
        for i in {1..10}; do
            if ! kill -0 "$PID" 2>/dev/null; then
                echo "Server stopped gracefully"
                rm -f "$PID_FILE"
                exit 0
            fi
            sleep 1
        done
        
        echo "Server did not stop gracefully, forcing..."
        kill -KILL "$PID"
        rm -f "$PID_FILE"
    else
        echo "PID file exists but process is not running, cleaning up..."
        rm -f "$PID_FILE"
    fi
fi

PID=$(lsof -t -i :$SERVER_PORT 2>/dev/null || true)
if [ -n "$PID" ]; then
    echo "Found process on port $SERVER_PORT (PID: $PID), stopping..."
    kill -TERM "$PID"
    sleep 2
    if kill -0 "$PID" 2>/dev/null; then
        echo "Force killing process..."
        kill -KILL "$PID"
    fi
fi

echo "Server stopped"
