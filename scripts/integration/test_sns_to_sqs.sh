#!/bin/bash
# SNS to SQS Integration Test
# Tests SNS message delivery to SQS subscription

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "SNS → SQS Integration Test"
echo ""

QUEUE_NAME="integration-queue-$$"
TOPIC_NAME="integration-topic-$$"
QUEUE_URL=""
TOPIC_ARN=""
SUB_ARN=""

cleanup() {
    [ -n "$SUB_ARN" ] && aws_cmd sns unsubscribe --subscription-arn "$SUB_ARN" 2>/dev/null || true
    [ -n "$TOPIC_ARN" ] && aws_cmd sns delete-topic --topic-arn "$TOPIC_ARN" 2>/dev/null || true
    [ -n "$QUEUE_URL" ] && aws_noauth sqs delete-queue --queue-url "$QUEUE_URL" 2>/dev/null || true
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

run_test "2. Get Queue ARN" \
    "aws_noauth sqs get-queue-attributes --queue-url '$QUEUE_URL' --attribute-names QueueArn --query 'Attributes.QueueArn' --output text" \
    "$QUEUE_NAME"

# Test 3: Create SNS Topic (use aws_cmd for SNS)
echo -n "Test 3: Create SNS Topic... "
TOPIC_ARN=$(aws_cmd sns create-topic --name "$TOPIC_NAME" --query 'TopicArn' --output text 2>&1)
if echo "$TOPIC_ARN" | grep -q "$TOPIC_NAME"; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $TOPIC_ARN"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "4. Subscribe SQS to SNS" \
    "SUB_ARN=\$(aws_cmd sns subscribe --topic-arn '$TOPIC_ARN' --protocol sqs --notification-endpoint '$QUEUE_URL' --query 'SubscriptionArn' --output text) && echo \$SUB_ARN" \
    "$TOPIC_NAME"

sleep 1

run_test "5. Publish Message to SNS" \
    "aws_cmd sns publish --topic-arn '$TOPIC_ARN' --message 'Hello from SNS to SQS!' --subject 'Integration Test'" \
    "MessageId"

sleep 2

run_test "6. Receive Message from SQS" \
    "aws_noauth sqs receive-message --queue-url '$QUEUE_URL' --max-number-of-messages 1 --wait-time-seconds 5" \
    "Hello from SNS"

print_summary
