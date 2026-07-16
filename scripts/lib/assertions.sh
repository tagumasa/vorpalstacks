#!/bin/bash
# Test assertion functions

source "$(dirname "${BASH_SOURCE[0]}")/colors.sh"

# Test counters
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Run a single test
# Usage: run_test "test_name" "command" "expected_pattern"
run_test() {
    local name="$1"
    local cmd="$2"
    local expected="${3:-}"
    
    TEST_COUNT=$((TEST_COUNT + 1))
    echo -n "$name... "
    
    local result
    local exit_code
    result=$(eval "$cmd" 2>&1)
    exit_code=$?
    
    if [ -n "$expected" ]; then
        if echo "$result" | grep -q "$expected"; then
            echo -e "${GREEN}PASS${NC}"
            PASS_COUNT=$((PASS_COUNT + 1))
            return 0
        else
            echo -e "${RED}FAIL${NC}"
            echo "  Expected: $expected"
            echo "  Got: ${result:0:200}"
            FAIL_COUNT=$((FAIL_COUNT + 1))
            return 1
        fi
    else
        if [ $exit_code -eq 0 ]; then
            echo -e "${GREEN}PASS${NC}"
            PASS_COUNT=$((PASS_COUNT + 1))
            return 0
        else
            echo -e "${RED}FAIL${NC}"
            echo "  Exit code: $exit_code"
            echo "  Output: ${result:0:200}"
            FAIL_COUNT=$((FAIL_COUNT + 1))
            return 1
        fi
    fi
}

# Assert command succeeds (exit code 0)
assert_success() {
    local cmd="$1"
    local message="${2:-Command should succeed}"
    
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        return 1
    fi
}

# Assert command fails
assert_failure() {
    local cmd="$1"
    local message="${2:-Command should fail}"
    
    if ! eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        return 1
    fi
}

# Assert output contains pattern
assert_output_contains() {
    local output="$1"
    local pattern="$2"
    local message="${3:-Output should contain '$pattern'}"
    
    if echo "$output" | grep -q "$pattern"; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        echo "  Output: ${output:0:200}"
        return 1
    fi
}

# Assert output does not contain pattern
assert_output_not_contains() {
    local output="$1"
    local pattern="$2"
    local message="${3:-Output should not contain '$pattern'}"
    
    if ! echo "$output" | grep -q "$pattern"; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        return 1
    fi
}

# Assert two values are equal
assert_equals() {
    local expected="$1"
    local actual="$2"
    local message="${3:-Values should be equal}"
    
    if [ "$expected" = "$actual" ]; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        echo "  Expected: $expected"
        echo "  Actual: $actual"
        return 1
    fi
}

# Assert value is not empty
assert_not_empty() {
    local value="$1"
    local message="${2:-Value should not be empty}"
    
    if [ -n "$value" ]; then
        echo -e "${GREEN}✓${NC} $message"
        return 0
    else
        echo -e "${RED}✗${NC} $message"
        return 1
    fi
}

# Pass/fail tracking (for custom test implementations)
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    PASS_COUNT=$((PASS_COUNT + 1))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    [ -n "$2" ] && echo "  Error: $2"
    FAIL_COUNT=$((FAIL_COUNT + 1))
}

# Print test summary
print_summary() {
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}Test Summary${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo "Total: $TEST_COUNT"
    echo -e "${GREEN}Passed: $PASS_COUNT${NC}"
    echo -e "${RED}Failed: $FAIL_COUNT${NC}"
    echo ""
    
    if [ $FAIL_COUNT -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}Some tests failed.${NC}"
        return 1
    fi
}

# Reset counters (for multiple test suites)
reset_counters() {
    TEST_COUNT=0
    PASS_COUNT=0
    FAIL_COUNT=0
}
