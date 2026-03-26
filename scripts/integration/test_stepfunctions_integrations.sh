#!/bin/bash
# Step Functions Integrations Test
# Tests SFN integration with SQS, SNS

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "Step Functions Integrations Test"
echo ""

QUEUE_NAME="sfn-queue-$$"
TOPIC_NAME="sfn-topic-$$"
SM_SQS_NAME="sfn-sqs-sm-$$"
SM_SNS_NAME="sfn-sns-sm-$$"

QUEUE_URL=""
TOPIC_ARN=""
SM_SQS_ARN=""
SM_SNS_ARN=""

cleanup() {
    [ -n "$SM_SQS_ARN" ] && aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_SQS_ARN" 2>/dev/null || true
    [ -n "$SM_SNS_ARN" ] && aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_SNS_ARN" 2>/dev/null || true
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

# Test 2: Create SNS Topic
echo -n "Test 2: Create SNS Topic... "
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

SQS_SM_DEF=$(echo "$STATE_MACHINE_SQS" | sed "s/TEST_QUEUE/$QUEUE_NAME/g")
SNS_SM_DEF=$(echo "$STATE_MACHINE_SNS" | sed "s/TEST_TOPIC/$TOPIC_NAME/g")

# Test 3: Create SFM with SQS Task
echo -n "Test 3: Create SFM with SQS Task... "
SM_SQS_ARN=$(aws_noauth stepfunctions create-state-machine --name "$SM_SQS_NAME" --definition "$SQS_SM_DEF" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
if echo "$SM_SQS_ARN" | grep -q "stateMachine"; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $SM_SQS_ARN"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "4. Start SQS Execution" \
    "aws_noauth stepfunctions start-execution --state-machine-arn '$SM_SQS_ARN' --input '{}' --query 'executionArn' --output text" \
    "execution"

sleep 2

run_test "5. Verify SQS Message" \
    "aws_noauth sqs receive-message --queue-url '$QUEUE_URL' --max-number-of-messages 1" \
    "Step Functions"

# Test 6: Create SFM with SNS Task
echo -n "Test 6: Create SFM with SNS Task... "
SM_SNS_ARN=$(aws_noauth stepfunctions create-state-machine --name "$SM_SNS_NAME" --definition "$SNS_SM_DEF" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
if echo "$SM_SNS_ARN" | grep -q "stateMachine"; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $SM_SNS_ARN"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "7. Start SNS Execution" \
    "aws_noauth stepfunctions start-execution --state-machine-arn '$SM_SNS_ARN' --input '{}' --query 'executionArn' --output text" \
    "execution"

print_summary
