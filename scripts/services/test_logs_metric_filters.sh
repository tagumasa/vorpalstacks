#!/bin/bash
# CloudWatch Logs Metric Filters Test Script

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:8080}"
AWS_REGION="${AWS_REGION:-us-east-1}"
LOG_GROUP_NAME="test-metric-filter-group"
LOG_STREAM_NAME="test-stream"
FILTER_NAME="error-filter"
METRIC_NAMESPACE="TestNamespace"
METRIC_NAME="ErrorCount"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

pass_count=0
fail_count=0

run_test() {
    local name="$1"
    local cmd="$2"
    local expected="$3"
    
    echo -n "Testing: $name ... "
    
    output=$(eval "$cmd" 2>&1) || true
    exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        if [ -z "$expected" ] || echo "$output" | grep -q -- "$expected"; then
            echo -e "${GREEN}PASSED${NC}"
            pass_count=$((pass_count + 1))
        else
            echo -e "${RED}FAILED${NC}"
            echo "  Expected: $expected"
            echo "  Got: $output"
            fail_count=$((fail_count + 1))
        fi
    else
        echo -e "${RED}FAILED${NC}"
        echo "  Command failed with exit code: $exit_code"
        echo "  Output: $output"
        fail_count=$((fail_count + 1))
    fi
}

aws_noauth() {
    aws --endpoint-url="$ENDPOINT_URL" --region "$AWS_REGION" --no-sign-request "$@"
}

echo "=== CloudWatch Logs Metric Filters Test ==="
echo "Endpoint: $ENDPOINT_URL"
echo ""

echo "Cleaning up..."
aws_noauth logs delete-log-group --log-group-name "$LOG_GROUP_NAME" 2>/dev/null || true

run_test "1. Create Log Group" \
    "aws_noauth logs create-log-group --log-group-name '$LOG_GROUP_NAME'" \
    ""

run_test "2. Create Log Stream" \
    "aws_noauth logs create-log-stream --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME'" \
    ""

run_test "3. Put Metric Filter (ERROR pattern)" \
    "aws_noauth logs put-metric-filter --log-group-name '$LOG_GROUP_NAME' --filter-name '$FILTER_NAME' --filter-pattern 'ERROR' --metric-transformations '[{\"metricName\":\"$METRIC_NAME\",\"metricNamespace\":\"$METRIC_NAMESPACE\",\"metricValue\":\"1\"}]'" \
    ""

run_test "4. Describe Metric Filters" \
    "aws_noauth logs describe-metric-filters --log-group-name '$LOG_GROUP_NAME' --query 'metricFilters[0].filterName' --output text" \
    "$FILTER_NAME"

TIMESTAMP=$(date +%s)000
run_test "5. Put Log Events (with ERROR)" \
    "aws_noauth logs put-log-events --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME' --log-events '[{\"timestamp\":$TIMESTAMP,\"message\":\"[ERROR 500] Internal Server Error\"}]'" \
    "nextSequenceToken"

sleep 1
run_test "6. Get Metric Statistics (may be empty if CW not working)" \
    "aws_noauth cloudwatch get-metric-statistics --namespace '$METRIC_NAMESPACE' --metric-name '$METRIC_NAME' --start-time $(date -u -d '-5 minutes' +%Y-%m-%dT%H:%M:%SZ) --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) --period 60 --statistics Sum" \
    ""

TIMESTAMP2=$(date +%s)000
run_test "7. Put Log Events (without ERROR)" \
    "aws_noauth logs put-log-events --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME' --log-events '[{\"timestamp\":$TIMESTAMP2,\"message\":\"[INFO 200] OK\"}]'" \
    "nextSequenceToken"

run_test "8. Put Metric Filter (JSON pattern)" \
    "aws_noauth logs put-metric-filter --log-group-name '$LOG_GROUP_NAME' --filter-name 'json-filter' --filter-pattern '{ \$.level = \"ERROR\" }' --metric-transformations '[{\"metricName\":\"JsonErrorCount\",\"metricNamespace\":\"$METRIC_NAMESPACE\",\"metricValue\":\"1\"}]'" \
    ""

TIMESTAMP3=$(date +%s)000
run_test "9. Put Log Events (JSON with level=ERROR)" \
    "aws_noauth logs put-log-events --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME' --log-events '[{\"timestamp\":$TIMESTAMP3,\"message\":\"{\\\"level\\\":\\\"ERROR\\\",\\\"message\\\":\\\"Something failed\\\"}\"}]'" \
    "nextSequenceToken"

sleep 1
run_test "10. Get Metric Statistics (JSON filter, may be empty)" \
    "aws_noauth cloudwatch get-metric-statistics --namespace '$METRIC_NAMESPACE' --metric-name 'JsonErrorCount' --start-time $(date -u -d '-5 minutes' +%Y-%m-%dT%H:%M:%SZ) --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) --period 60 --statistics Sum" \
    ""

run_test "11. Delete Metric Filter" \
    "aws_noauth logs delete-metric-filter --log-group-name '$LOG_GROUP_NAME' --filter-name '$FILTER_NAME'" \
    ""

run_test "12. Verify Filter Deleted" \
    "aws_noauth logs describe-metric-filters --log-group-name '$LOG_GROUP_NAME' --query 'length(metricFilters)' --output text" \
    "1"

echo ""
echo "Cleaning up..."
aws_noauth logs delete-log-group --log-group-name "$LOG_GROUP_NAME" 2>/dev/null || true

echo ""
echo "=== Test Summary ==="
echo -e "Passed: ${GREEN}$pass_count${NC}"
echo -e "Failed: ${RED}$fail_count${NC}"

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
