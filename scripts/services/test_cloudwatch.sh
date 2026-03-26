#!/bin/bash
# CloudWatch Service Tests for VorpalStacks
# Usage: bash scripts/services/test_cloudwatch.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "CloudWatch Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

ALARM_NAME="test-alarm-$$"
NAMESPACE="Test/Namespace"
METRIC_NAME="TestMetric"

START_TIME=$(date -u -d "1 hour ago" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-1H +%Y-%m-%dT%H:%M:%SZ)
END_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

cleanup() {
    aws_noauth cloudwatch delete-alarms --alarm-names "$ALARM_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. PutMetricData" \
    "aws_noauth cloudwatch put-metric-data --namespace '$NAMESPACE' --metric-data '[{\"MetricName\":\"$METRIC_NAME\",\"Value\":1,\"Unit\":\"Count\"}]'" \
    ""

run_test "2. ListMetrics" \
    "aws_noauth cloudwatch list-metrics --namespace '$NAMESPACE' --query 'Metrics[?MetricName==\`$METRIC_NAME\`].MetricName' --output text" \
    "$METRIC_NAME"

run_test "3. GetMetricStatistics" \
    "aws_noauth cloudwatch get-metric-statistics --namespace '$NAMESPACE' --metric-name '$METRIC_NAME' --start-time '$START_TIME' --end-time '$END_TIME' --period 3600 --statistics Average" \
    ""

run_test "4. PutMetricAlarm" \
    "aws_noauth cloudwatch put-metric-alarm --alarm-name '$ALARM_NAME' --metric-name '$METRIC_NAME' --namespace '$NAMESPACE' --statistic Average --period 60 --evaluation-periods 1 --threshold 10 --comparison-operator GreaterThanThreshold" \
    ""

run_test "5. DescribeAlarms" \
    "aws_noauth cloudwatch describe-alarms --alarm-names '$ALARM_NAME' --query 'MetricAlarms[?AlarmName==\`$ALARM_NAME\`].AlarmName' --output text" \
    "$ALARM_NAME"

run_test "6. DescribeAlarmsForMetric" \
    "aws_noauth cloudwatch describe-alarms-for-metric --metric-name '$METRIC_NAME' --namespace '$NAMESPACE' --query 'MetricAlarms[?AlarmName==\`$ALARM_NAME\`].AlarmName' --output text" \
    "$ALARM_NAME"

run_test "7. GetMetricData" \
    "aws_noauth cloudwatch get-metric-data --metric-data-queries '[{\"Id\":\"m1\",\"MetricStat\":{\"Metric\":{\"Namespace\":\"$NAMESPACE\",\"MetricName\":\"$METRIC_NAME\"},\"Period\":3600,\"Stat\":\"Average\"}}]' --start-time '$START_TIME' --end-time '$END_TIME'" \
    ""

run_test "8. SetAlarmState" \
    "aws_noauth cloudwatch set-alarm-state --alarm-name '$ALARM_NAME' --state-value ALARM --state-reason 'Test'" \
    ""

run_test "9. EnableAlarmActions" \
    "aws_noauth cloudwatch enable-alarm-actions --alarm-names '$ALARM_NAME'" \
    ""

run_test "10. DisableAlarmActions" \
    "aws_noauth cloudwatch disable-alarm-actions --alarm-names '$ALARM_NAME'" \
    ""

run_test "11. DeleteAlarms" \
    "aws_noauth cloudwatch delete-alarms --alarm-names '$ALARM_NAME'" \
    ""

print_summary
