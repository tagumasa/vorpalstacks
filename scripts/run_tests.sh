#!/bin/bash
# Main test runner for VorpalStacks
# Usage: ./run_tests.sh [options]
#   --service X        Run only service X tests
#   --integration      Run integration tests only
#   --all              Run all tests (services + integration)
#   --no-server        Don't start server (assume running)
#   --quick            Run quick tests only (skip integration)
#   --verbose          Show verbose output
#   --list             List available tests

set +e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/lib/config.sh"
source "$SCRIPT_DIR/lib/colors.sh"
source "$SCRIPT_DIR/lib/server.sh"
source "$SCRIPT_DIR/lib/assertions.sh"

PARALLEL=false
SERVICE_FILTER=""
QUICK_MODE=false
VERBOSE=false
NO_SERVER=false
INTEGRATION_ONLY=false
RUN_ALL=false
LIST_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --parallel) PARALLEL=true; shift ;;
        --service) SERVICE_FILTER="$2"; shift 2 ;;
        --quick) QUICK_MODE=true; shift ;;
        --verbose) VERBOSE=true; shift ;;
        --no-server) NO_SERVER=true; shift ;;
        --integration) INTEGRATION_ONLY=true; shift ;;
        --all) RUN_ALL=true; shift ;;
        --list) LIST_ONLY=true; shift ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --service X        Run only service X tests"
            echo "  --integration      Run integration tests only"
            echo "  --all              Run all tests (default)"
            echo "  --no-server        Don't start server (assume running)"
            echo "  --quick            Run quick tests only (skip integration)"
            echo "  --verbose          Show verbose output"
            echo "  --list             List available tests"
            echo ""
            echo "Examples:"
            echo "  $0                          # Run all tests"
            echo "  $0 --service sqs            # Run SQS tests only"
            echo "  $0 --integration            # Run integration tests only"
            echo "  $0 --no-server --service s3 # Run S3 tests (server already running)"
            exit 0
            ;;
        *) shift ;;
    esac
done

# Test categories
declare -A SERVICE_TESTS=(
    ["sqs"]="services/test_sqs.sh"
    ["sns"]="services/test_sns.sh"
    ["stepfunctions"]="services/test_stepfunctions.sh"
    ["events"]="services/test_events.sh"
    ["dynamodb"]="services/test_dynamodb.sh"
    ["s3"]="services/test_s3.sh"
    ["lambda"]="services/test_lambda.sh"
    ["apigateway"]="services/test_apigateway.sh"
    ["kinesis"]="services/test_kinesis.sh"
    ["iam"]="services/test_iam.sh"
    ["kms"]="services/test_kms.sh"
    ["sts"]="services/test_sts.sh"
    ["cognito"]="services/test_cognito.sh"
    ["secretsmanager"]="services/test_secretsmanager.sh"
    ["cloudwatch"]="services/test_cloudwatch.sh"
    ["logs"]="services/test_logs.sh"
    ["ssm"]="services/test_ssm.sh"
    ["scheduler"]="services/test_scheduler.sh"
)

declare -A INTEGRATION_TESTS=(
    ["sns_to_sqs"]="integration/test_sns_to_sqs.sh"
    ["events_to_targets"]="integration/test_events_to_targets.sh"
    ["stepfunctions_integrations"]="integration/test_stepfunctions_integrations.sh"
    ["apigateway_flow"]="integration/test_apigateway_flow.sh"
)

# List tests
if [ "$LIST_ONLY" = true ]; then
    echo "Available Service Tests:"
    for service in $(echo "${!SERVICE_TESTS[@]}" | tr ' ' '\n' | sort); do
        echo "  $service"
    done
    echo ""
    echo "Available Integration Tests:"
    for test in $(echo "${!INTEGRATION_TESTS[@]}" | tr ' ' '\n' | sort); do
        echo "  $test"
    done
    exit 0
fi

# Track results
declare -A RESULTS
TOTAL_PASSED=0
TOTAL_FAILED=0

# Run a single test script
run_test_script() {
    local name="$1"
    local script="$2"
    local script_path="$SCRIPT_DIR/$script"
    
    if [ ! -f "$script_path" ]; then
        echo -e "${YELLOW}SKIP${NC}: $name (script not found: $script)"
        return 0
    fi
    
    print_section "Testing: $name"
    
    local start_time=$(date +%s)
    local output
    local exit_code
    
    if [ "$VERBOSE" = true ]; then
        bash "$script_path"
        exit_code=$?
    else
        output=$(bash "$script_path" 2>&1)
        exit_code=$?
    fi
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $name (${duration}s)"
        RESULTS[$name]="PASS"
        TOTAL_PASSED=$((TOTAL_PASSED + 1))
        return 0
    else
        echo -e "${RED}✗${NC} $name (${duration}s)"
        [ "$VERBOSE" = false ] && [ -n "$output" ] && echo "$output" | tail -20
        RESULTS[$name]="FAIL"
        TOTAL_FAILED=$((TOTAL_FAILED + 1))
        return 1
    fi
}

# Main
print_header "VorpalStacks Test Suite"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

# Start server (unless --no-server)
if [ "$NO_SERVER" = false ]; then
    print_info "Starting test server..."
    if ! start_test_server; then
        print_error "Failed to start server"
        exit 1
    fi
else
    print_info "Using existing server at $ENDPOINT_URL"
fi

# Run service tests (unless --integration)
if [ "$INTEGRATION_ONLY" = false ]; then
    print_header "Service Tests"
    for service in $(echo "${!SERVICE_TESTS[@]}" | tr ' ' '\n' | sort); do
        if [ -n "$SERVICE_FILTER" ] && [ "$service" != "$SERVICE_FILTER" ]; then
            continue
        fi
        run_test_script "$service" "${SERVICE_TESTS[$service]}"
    done
fi

# Run integration tests (unless --quick or --service)
if [ "$QUICK_MODE" = false ] && [ -z "$SERVICE_FILTER" ]; then
    print_header "Integration Tests"
    for test_name in $(echo "${!INTEGRATION_TESTS[@]}" | tr ' ' '\n' | sort); do
        run_test_script "$test_name" "${INTEGRATION_TESTS[$test_name]}"
    done
fi

# Stop server (if we started it)
if [ "$NO_SERVER" = false ]; then
    stop_test_server
fi

# Print summary
print_header "Test Results Summary"
echo ""
echo -e "${GREEN}Passed: $TOTAL_PASSED${NC}"
echo -e "${RED}Failed: $TOTAL_FAILED${NC}"
echo ""

if [ $TOTAL_FAILED -gt 0 ]; then
    echo "Failed tests:"
    for name in "${!RESULTS[@]}"; do
        if [ "${RESULTS[$name]}" = "FAIL" ]; then
            echo -e "  ${RED}✗${NC} $name"
        fi
    done
    exit 1
else
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi
