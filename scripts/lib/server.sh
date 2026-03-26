#!/bin/bash
# Server management functions for tests

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/colors.sh"

SERVER_PID=""

# Check if server is running on specified port
is_server_running() {
    local port="${1:-$PORT}"
    lsof -i :$port | grep -q LISTEN 2>/dev/null
}

# Wait for server to be ready
wait_for_server() {
    local endpoint="${1:-$ENDPOINT_URL}"
    local timeout="${2:-$HEALTH_CHECK_TIMEOUT}"
    local start_time=$(date +%s)
    
    echo -n "Waiting for server at $endpoint..."
    
    while true; do
        if curl -s -o /dev/null -w "%{http_code}" "$endpoint" 2>/dev/null | grep -q "200\|404\|403"; then
            echo -e " ${GREEN}OK${NC}"
            return 0
        fi
        
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        if [ $elapsed -ge $timeout ]; then
            echo -e " ${RED}TIMEOUT${NC}"
            return 1
        fi
        
        echo -n "."
        sleep 1
    done
}

# Start test server
start_test_server() {
    local port="${1:-$PORT}"
    local data_path="${2:-$TEST_DATA_PATH}"
    
    if is_server_running "$port"; then
        print_warning "Server already running on port $port"
        return 0
    fi
    
    print_info "Starting server on port $port..."
    
    # Set environment variables
    export SIGNATURE_VERIFICATION_ENABLED=false
    export PORT=$port
    export DATA_PATH="$data_path"
    export TEST_MODE=true
    export LOG_LEVEL=ERROR
    
    # Clean up old data
    rm -rf "$data_path"
    mkdir -p "$data_path"
    mkdir -p "$(dirname "$SERVER_LOG")"
    
    # Start server
    cd "$PROJECT_ROOT"
    go run main.go > "$SERVER_LOG" 2>&1 &
    SERVER_PID=$!
    
    # Wait for server to be ready
    if wait_for_server "http://localhost:$port"; then
        print_info "Server started (PID: $SERVER_PID)"
        return 0
    else
        print_error "Failed to start server"
        [ -f "$SERVER_LOG" ] && cat "$SERVER_LOG"
        return 1
    fi
}

# Stop test server
stop_test_server() {
    local port="${1:-$PORT}"
    
    if [ -n "$SERVER_PID" ] && kill -0 $SERVER_PID 2>/dev/null; then
        print_info "Stopping server (PID: $SERVER_PID)..."
        kill -TERM $SERVER_PID 2>/dev/null
        
        # Wait for graceful shutdown
        local count=0
        while kill -0 $SERVER_PID 2>/dev/null && [ $count -lt 10 ]; do
            sleep 1
            count=$((count + 1))
        done
        
        # Force kill if still running
        if kill -0 $SERVER_PID 2>/dev/null; then
            print_warning "Force killing server..."
            kill -KILL $SERVER_PID 2>/dev/null
        fi
        
        SERVER_PID=""
    fi
    
    # Also check for any process on the port
    local pid=$(lsof -t -i :$port 2>/dev/null || true)
    if [ -n "$pid" ]; then
        print_info "Killing process on port $port (PID: $pid)..."
        kill -TERM $pid 2>/dev/null || true
    fi
}

# Cleanup function for trap
cleanup_server() {
    stop_test_server
}

# Setup test environment
setup_test_env() {
    local port="${1:-$PORT}"
    
    start_test_server "$port"
    trap cleanup_server EXIT
}

# Teardown test environment
teardown_test_env() {
    stop_test_server
    rm -rf "$TEST_DATA_PATH"
}
