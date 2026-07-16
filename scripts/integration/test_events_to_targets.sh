#!/bin/bash
# Events to Targets Integration Test
# Tests EventBridge rule delivery to SQS, SNS

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Events → Targets Integration Test"
echo ""

BUS_NAME="events-bus-$$"
RULE_NAME="events-rule-$$"
QUEUE_NAME="events-queue-$$"
TOPIC_NAME="events-topic-$$"
TARGET_ID_SQS="SQSTarget"
TARGET_ID_SNS="SNSTarget"

QUEUE_URL=""
TOPIC_ARN=""

cleanup() {
    aws_noauth events remove-targets --rule "$RULE_NAME" --event-bus-name "$BUS_NAME" --ids "$TARGET_ID_SQS" "$TARGET_ID_SNS" 2>/dev/null || true
    aws_noauth events delete-rule --name "$RULE_NAME" --event-bus-name "$BUS_NAME" 2>/dev/null || true
    aws_noauth events delete-event-bus --name "$BUS_NAME" 2>/dev/null || true
    [ -n "$QUEUE_URL" ] && aws_noauth sqs delete-queue --queue-url "$QUEUE_URL" 2>/dev/null || true
    [ -n "$TOPIC_ARN" ] && aws_noauth sns delete-topic --topic-arn "$TOPIC_ARN" 2>/dev/null || true
}
trap cleanup EXIT

# Test 1: Create SQS Queue
echo -n "Test 1: Create SQS Queue... "
QUEUE_URL=$(aws_noauth sqs create-queue --queue-name "$QUEUE_NAME" --query 'QueueUrl' --output text 2>&1)
if echo "$QUEUE_URL" | grep -q "$QUEUE_NAME"; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $QUEUE_URL"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

# Test 2: Create SNS Topic
echo -n "Test 2: Create SNS Topic... "
TOPIC_ARN=$(aws_noauth sns create-topic --name "$TOPIC_NAME" --query 'TopicArn' --output text 2>&1)
if echo "$TOPIC_ARN" | grep -q "$TOPIC_NAME"; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $TOPIC_ARN"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "3. Create Event Bus" \
    "aws_noauth events create-event-bus --name '$BUS_NAME' --query 'EventBusArn' --output text" \
    "$BUS_NAME"

run_test "4. PutRule" \
    "aws_noauth events put-rule --name '$RULE_NAME' --event-bus-name '$BUS_NAME' --event-pattern '{\"source\":[\"test.integration\"]}' --query 'RuleArn' --output text" \
    "$RULE_NAME"

run_test "5. PutTargets (SQS)" \
    "aws_noauth events put-targets --rule '$RULE_NAME' --event-bus-name '$BUS_NAME' --targets 'Id=$TARGET_ID_SQS,Arn=arn:aws:sqs:$REGION:$ACCOUNT_ID:$QUEUE_NAME' --query 'FailedEntryCount' --output text" \
    "0"

run_test "6. ListTargetsByRule" \
    "aws_noauth events list-targets-by-rule --rule '$RULE_NAME' --event-bus-name '$BUS_NAME' --query 'Targets | length' --output text" \
    "1"

run_test "7. PutEvents" \
    "aws_noauth events put-events --entries '[{\"Source\":\"test.integration\",\"DetailType\":\"test-event\",\"Detail\":\"{\\\"message\\\":\\\"hello\\\"}\",\"EventBusName\":\"$BUS_NAME\"}]' --query 'Entries[0].EventId' --output text" \
    ""

sleep 2

run_test "8. Verify SQS Received Event" \
    "aws_noauth sqs receive-message --queue-url '$QUEUE_URL' --max-number-of-messages 1" \
    "test.integration"

print_summary
