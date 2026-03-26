#!/bin/bash
# Scheduler Service Tests for VorpalStacks
# Usage: bash scripts/services/test_scheduler.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Scheduler Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

SCHEDULE_GROUP_NAME="test-schedule-group-$$"
SCHEDULE_NAME="test-schedule-$$"
TARGET_ARN="arn:aws:sqs:$REGION:$ACCOUNT_ID:test-queue"

cleanup() {
    aws_noauth scheduler delete-schedule --name "$SCHEDULE_NAME" --group-name "$SCHEDULE_GROUP_NAME" 2>/dev/null || true
    aws_noauth scheduler delete-schedule-group --name "$SCHEDULE_GROUP_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. CreateScheduleGroup" \
    "aws_noauth scheduler create-schedule-group --name '$SCHEDULE_GROUP_NAME' --query 'ScheduleGroupArn' --output text" \
    "$SCHEDULE_GROUP_NAME"

run_test "2. GetScheduleGroup" \
    "aws_noauth scheduler get-schedule-group --name '$SCHEDULE_GROUP_NAME' --query 'Name' --output text" \
    "$SCHEDULE_GROUP_NAME"

run_test "3. ListScheduleGroups" \
    "aws_noauth scheduler list-schedule-groups --query \"ScheduleGroups[?@.Name=='$SCHEDULE_GROUP_NAME'] | [0].Name\" --output text" \
    "$SCHEDULE_GROUP_NAME"

run_test "4. CreateSchedule" \
    "aws_noauth scheduler create-schedule --name '$SCHEDULE_NAME' --group-name '$SCHEDULE_GROUP_NAME' --schedule-expression 'rate(1 day)' --target '{\"Arn\":\"$TARGET_ARN\",\"RoleArn\":\"arn:aws:iam::$ACCOUNT_ID:role/scheduler-role\"}' --flexible-time-window '{\"Mode\":\"OFF\"}' --query 'ScheduleArn' --output text" \
    "$SCHEDULE_NAME"

run_test "5. GetSchedule" \
    "aws_noauth scheduler get-schedule --name '$SCHEDULE_NAME' --group-name '$SCHEDULE_GROUP_NAME' --query 'Name' --output text" \
    "$SCHEDULE_NAME"

run_test "6. ListSchedules" \
    "aws_noauth scheduler list-schedules --group-name '$SCHEDULE_GROUP_NAME' --query \"Schedules[?@.Name=='$SCHEDULE_NAME'] | [0].Name\" --output text" \
    "$SCHEDULE_NAME"

run_test "7. UpdateSchedule" \
    "aws_noauth scheduler update-schedule --name '$SCHEDULE_NAME' --group-name '$SCHEDULE_GROUP_NAME' --schedule-expression 'rate(2 days)' --target '{\"Arn\":\"$TARGET_ARN\",\"RoleArn\":\"arn:aws:iam::$ACCOUNT_ID:role/scheduler-role\"}' --flexible-time-window '{\"Mode\":\"OFF\"}' --query 'ScheduleArn' --output text" \
    "$SCHEDULE_NAME"

run_test "8. TagResource" \
    "aws_noauth scheduler tag-resource --resource-arn 'arn:aws:scheduler:$REGION:$ACCOUNT_ID:schedule/$SCHEDULE_GROUP_NAME/$SCHEDULE_NAME' --tags '[{\"Key\":\"Environment\",\"Value\":\"Test\"}]'" \
    ""

run_test "9. ListTagsForResource" \
    "aws_noauth scheduler list-tags-for-resource --resource-arn 'arn:aws:scheduler:$REGION:$ACCOUNT_ID:schedule/$SCHEDULE_GROUP_NAME/$SCHEDULE_NAME' --query 'Tags[0].Value' --output text" \
    "Test"

run_test "10. DeleteSchedule" \
    "aws_noauth scheduler delete-schedule --name '$SCHEDULE_NAME' --group-name '$SCHEDULE_GROUP_NAME'" \
    ""

run_test "11. DeleteScheduleGroup" \
    "aws_noauth scheduler delete-schedule-group --name '$SCHEDULE_GROUP_NAME'" \
    ""

print_summary
