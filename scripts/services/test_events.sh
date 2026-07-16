#!/bin/bash
# EventBridge Service Tests for VorpalStacks
# Usage: bash scripts/services/test_events.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "EventBridge Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

BUS_NAME="test-bus-$$"
REPLAY_BUS_NAME="replay-bus-$$"
RULE_NAME="test-rule-$$"
TARGET_ID="Target1"
ARCHIVE_NAME="test-archive-$$"
REPLAY_NAME="test-replay-$$"

cleanup() {
    aws_noauth events remove-targets --rule "$RULE_NAME" --event-bus-name "$BUS_NAME" --ids "$TARGET_ID" 2>/dev/null || true
    aws_noauth events delete-rule --name "$RULE_NAME" --event-bus-name "$BUS_NAME" 2>/dev/null || true
    aws_noauth events delete-event-bus --name "$BUS_NAME" 2>/dev/null || true
    aws_noauth events delete-event-bus --name "$REPLAY_BUS_NAME" 2>/dev/null || true
    aws_noauth events delete-archive --archive-name "$ARCHIVE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. CreateEventBus" \
    "aws_noauth events create-event-bus --name '$BUS_NAME'" \
    "EventBusArn"

run_test "2. DescribeEventBus" \
    "aws_noauth events describe-event-bus --name '$BUS_NAME' --query 'Name' --output text" \
    "$BUS_NAME"

run_test "3. ListEventBuses" \
    "aws_noauth events list-event-buses --query 'EventBuses[?Name==\`$BUS_NAME\`].Name' --output text" \
    "$BUS_NAME"

run_test "4. PutRule" \
    "aws_noauth events put-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME' --event-pattern '$EVENT_PATTERN_SIMPLE'" \
    "RuleArn"

run_test "5. DescribeRule" \
    "aws_noauth events describe-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME' --query 'Name' --output text" \
    "$RULE_NAME"

run_test "6. ListRules" \
    "aws_noauth events list-rules --event-bus-name '$BUS_NAME' --query 'Rules[0].Name' --output text" \
    "$RULE_NAME"

run_test "7. PutTargets" \
    "aws_noauth events put-targets --rule '$RULE_NAME' --event-bus-name '$BUS_NAME' --targets 'Id=$TARGET_ID,Arn=arn:aws:sqs:$REGION:$ACCOUNT_ID:test-queue'" \
    "FailedEntryCount"

run_test "8. ListTargetsByRule" \
    "aws_noauth events list-targets-by-rule --rule '$RULE_NAME' --event-bus-name '$BUS_NAME' --query 'Targets[0].Id' --output text" \
    "$TARGET_ID"

run_test "9. DisableRule" \
    "aws_noauth events disable-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME'" \
    ""

run_test "10. EnableRule" \
    "aws_noauth events enable-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME'" \
    ""

run_test "11. PutEvents" \
    "aws_noauth events put-events --entries '[{\"Source\":\"test.source\",\"DetailType\":\"test-event\",\"Detail\":\"{\\\"test\\\":\\\"data\\\"}\",\"EventBusName\":\"$BUS_NAME\"}]'" \
    "Entries"

run_test "12. RemoveTargets" \
    "aws_noauth events remove-targets --rule '$RULE_NAME' --event-bus-name '$BUS_NAME' --ids '$TARGET_ID'" \
    "FailedEntryCount"

run_test "13. DeleteRule" \
    "aws_noauth events delete-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME'" \
    ""

run_test "14. DeleteEventBus" \
    "aws_noauth events delete-event-bus --name '$BUS_NAME'" \
    ""

run_test "15. CreateEventBus for Replay" \
    "aws_noauth events create-event-bus --name '$REPLAY_BUS_NAME'" \
    "EventBusArn"

run_test "16. CreateArchive" \
    "aws_noauth events create-archive --archive-name '$ARCHIVE_NAME' --event-source-arn \"arn:aws:events:$REGION:$ACCOUNT_ID:event-bus/$REPLAY_BUS_NAME\"" \
    "ArchiveArn"

run_test "17. DescribeArchive" \
    "aws_noauth events describe-archive --archive-name '$ARCHIVE_NAME' --query 'ArchiveName' --output text" \
    "$ARCHIVE_NAME"

START_TIME=$(date -u -d '1 hour ago' +%s)
END_TIME=$(date -u +%s)

run_test "18. PutEvents to Archive" \
    "aws_noauth events put-events --entries '[{\"Source\":\"archive.test\",\"DetailType\":\"test-event\",\"Detail\":\"{\\\"test\\\":\\\"data\\\"}\",\"EventBusName\":\"$REPLAY_BUS_NAME\"}]'" \
    "Entries"

run_test "19. StartReplay" \
    "aws_noauth events start-replay --replay-name '$REPLAY_NAME' --event-source-arn \"arn:aws:events:$REGION:$ACCOUNT_ID:archive/$ARCHIVE_NAME\" --destination '{\"Arn\":\"arn:aws:events:$REGION:$ACCOUNT_ID:event-bus/$REPLAY_BUS_NAME\"}' --event-start-time $(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%SZ) --event-end-time $(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    "ReplayArn"

run_test "20. DescribeReplay" \
    "aws_noauth events describe-replay --replay-name '$REPLAY_NAME' --query 'State' --output text" \
    "COMPLETED"

run_test "21. ListReplays" \
    "aws_noauth events list-replays --query 'Replays[?ReplayName==\`$REPLAY_NAME\`].ReplayName' --output text" \
    "$REPLAY_NAME"

print_summary
