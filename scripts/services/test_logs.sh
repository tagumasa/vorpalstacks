#!/bin/bash
# CloudWatch Logs Service Tests for VorpalStacks
# Usage: bash scripts/services/test_logs.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "CloudWatch Logs Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

LOG_GROUP_NAME="/test/log-group-$$"
LOG_STREAM_NAME="test-stream-$$"

cleanup() {
    aws_noauth logs delete-log-group --log-group-name "$LOG_GROUP_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. CreateLogGroup" \
    "aws_noauth logs create-log-group --log-group-name '$LOG_GROUP_NAME'" \
    ""

run_test "2. DescribeLogGroups" \
    "aws_noauth logs describe-log-groups --log-group-name-prefix '$LOG_GROUP_NAME' --query 'logGroups[0].logGroupName' --output text" \
    "$LOG_GROUP_NAME"

run_test "3. PutRetentionPolicy" \
    "aws_noauth logs put-retention-policy --log-group-name '$LOG_GROUP_NAME' --retention-in-days 7" \
    ""

run_test "4. CreateLogStream" \
    "aws_noauth logs create-log-stream --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME'" \
    ""

run_test "5. DescribeLogStreams" \
    "aws_noauth logs describe-log-streams --log-group-name '$LOG_GROUP_NAME' --query 'logStreams[0].logStreamName' --output text" \
    "$LOG_STREAM_NAME"

TIMESTAMP=$(date +%s)000
run_test "6. PutLogEvents" \
    "aws_noauth logs put-log-events --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME' --log-events '[{\"timestamp\":$TIMESTAMP,\"message\":\"test log message\"}]'" \
    ""

run_test "7. GetLogEvents" \
    "aws_noauth logs get-log-events --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME' --query 'events[0].message' --output text" \
    "test log message"

run_test "8. TagLogGroup" \
    "aws_noauth logs tag-log-group --log-group-name '$LOG_GROUP_NAME' --tags Environment=Test" \
    ""

run_test "9. ListTagsLogGroup" \
    "aws_noauth logs list-tags-log-group --log-group-name '$LOG_GROUP_NAME' --query 'tags.Environment' --output text" \
    "Test"

run_test "10. FilterLogEvents" \
    "aws_noauth logs filter-log-events --log-group-name '$LOG_GROUP_NAME' --query 'events[0].message' --output text" \
    "test log message"

run_test "11. PutMetricFilter" \
    "aws_noauth logs put-metric-filter --log-group-name '$LOG_GROUP_NAME' --filter-name 'TestFilter' --filter-pattern '[timestamp, request_id, error_level=*...]' --metric-transformations '[{\"metricName\":\"TestMetric\",\"metricNamespace\":\"TestNamespace\",\"metricValue\":\"1\"}]'" \
    ""

run_test "12. DescribeMetricFilters" \
    "aws_noauth logs describe-metric-filters --log-group-name '$LOG_GROUP_NAME' --query 'metricFilters[0].filterName' --output text" \
    "TestFilter"

run_test "13. DeleteMetricFilter" \
    "aws_noauth logs delete-metric-filter --log-group-name '$LOG_GROUP_NAME' --filter-name 'TestFilter'" \
    ""

KINESIS_STREAM="test-logs-stream-$$"
aws_noauth kinesis create-stream --stream-name "$KINESIS_STREAM" --shard-count 1 2>/dev/null || true
sleep 2

KINESIS_ARN="arn:aws:kinesis:$REGION:000000000000:stream/$KINESIS_STREAM"
run_test "14. PutSubscriptionFilter" \
    "aws_noauth logs put-subscription-filter --log-group-name '$LOG_GROUP_NAME' --filter-name 'TestSubFilter' --filter-pattern '' --destination-arn '$KINESIS_ARN'" \
    ""

run_test "15. DescribeSubscriptionFilters" \
    "aws_noauth logs describe-subscription-filters --log-group-name '$LOG_GROUP_NAME' --query 'subscriptionFilters[0].filterName' --output text" \
    "TestSubFilter"

run_test "16. DeleteSubscriptionFilter" \
    "aws_noauth logs delete-subscription-filter --log-group-name '$LOG_GROUP_NAME' --filter-name 'TestSubFilter'" \
    ""

aws_noauth kinesis delete-stream --stream-name "$KINESIS_STREAM" 2>/dev/null || true

run_test "17. DeleteLogStream" \
    "aws_noauth logs delete-log-stream --log-group-name '$LOG_GROUP_NAME' --log-stream-name '$LOG_STREAM_NAME'" \
    ""

run_test "18. DeleteLogGroup" \
    "aws_noauth logs delete-log-group --log-group-name '$LOG_GROUP_NAME'" \
    ""

print_summary
