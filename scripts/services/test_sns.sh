#!/bin/bash
# SNS Service Tests for VorpalStacks
# Usage: bash scripts/services/test_sns.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "SNS Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

TOPIC_NAME="test-topic-$$"
TOPIC_ARN=""
SUB_ARN=""

cleanup() {
    [ -n "$SUB_ARN" ] && aws_noauth sns unsubscribe --subscription-arn "$SUB_ARN" 2>/dev/null || true
    [ -n "$TOPIC_ARN" ] && aws_noauth sns delete-topic --topic-arn "$TOPIC_ARN" 2>/dev/null || true
}
trap cleanup EXIT

TOPIC_ARN=$(aws_noauth sns create-topic --name "$TOPIC_NAME" --query 'TopicArn' --output text 2>&1)
run_test "1. CreateTopic" \
    "echo '$TOPIC_ARN'" \
    "$TOPIC_NAME"

run_test "2. ListTopics" \
    "aws_noauth sns list-topics --query 'Topics[?contains(TopicArn,\`$TOPIC_NAME\`)].TopicArn' --output text" \
    "$TOPIC_NAME"

run_test "3. GetTopicAttributes" \
    "aws_noauth sns get-topic-attributes --topic-arn '$TOPIC_ARN'" \
    "TopicArn"

run_test "4. SetTopicAttributes" \
    "aws_noauth sns set-topic-attributes --topic-arn '$TOPIC_ARN' --attribute-name DisplayName --attribute-value 'Test Display'" \
    ""

run_test "5. TagResource" \
    "aws_cmd sns tag-resource --resource-arn '$TOPIC_ARN' --tags Key=Environment,Value=Test" \
    ""

run_test "6. ListTagsForResource" \
    "aws_cmd sns list-tags-for-resource --resource-arn '$TOPIC_ARN'" \
    "Environment"

SUB_ARN=$(aws_noauth sns subscribe --topic-arn "$TOPIC_ARN" --protocol email --notification-endpoint 'test@example.com' --query 'SubscriptionArn' --output text 2>&1)
run_test "7. Subscribe (email)" \
    "echo '$SUB_ARN'" \
    "$TOPIC_NAME"

run_test "8. ListSubscriptions" \
    "aws_noauth sns list-subscriptions --query 'Subscriptions[?contains(TopicArn,\`$TOPIC_NAME\`)].TopicArn' --output text" \
    "$TOPIC_NAME"

run_test "9. ListSubscriptionsByTopic" \
    "aws_noauth sns list-subscriptions-by-topic --topic-arn '$TOPIC_ARN'" \
    "test@example.com"

run_test "10. GetSubscriptionAttributes" \
    "aws_noauth sns get-subscription-attributes --subscription-arn '$SUB_ARN'" \
    "TopicArn"

run_test "11. Publish" \
    "aws_noauth sns publish --topic-arn '$TOPIC_ARN' --message 'Hello from SNS test!'" \
    "MessageId"

run_test "12. Publish with Subject" \
    "aws_noauth sns publish --topic-arn '$TOPIC_ARN' --message 'Test message' --subject 'Test Subject'" \
    "MessageId"

run_test "13. Publish with MessageAttributes" \
    "aws_noauth sns publish --topic-arn '$TOPIC_ARN' --message 'Test' --message-attributes '{\"key\":{\"DataType\":\"String\",\"StringValue\":\"value\"}}'" \
    "MessageId"

run_test "14. Unsubscribe" \
    "aws_noauth sns unsubscribe --subscription-arn '$SUB_ARN'" \
    ""

run_test "15. UntagResource" \
    "aws_cmd sns untag-resource --resource-arn '$TOPIC_ARN' --tag-keys Environment" \
    ""

run_test "16. DeleteTopic" \
    "aws_noauth sns delete-topic --topic-arn '$TOPIC_ARN'" \
    ""

# FIFO Tests
FIFO_TOPIC_NAME="test-fifo-topic-$$"
FIFO_TOPIC_ARN=""

cleanup_fifo() {
    [ -n "$FIFO_TOPIC_ARN" ] && aws_noauth sns delete-topic --topic-arn "$FIFO_TOPIC_ARN" 2>/dev/null || true
}

echo ""
echo "=== SNS FIFO Tests ==="

FIFO_TOPIC_ARN=$(aws_noauth sns create-topic \
    --name "${FIFO_TOPIC_NAME}.fifo" \
    --attributes '{"FifoTopic":"true","ContentBasedDeduplication":"true"}' \
    --query 'TopicArn' --output text 2>&1)
run_test "17. CreateFifoTopic" \
    "echo '$FIFO_TOPIC_ARN'" \
    "fifo"

run_test "18. FifoTopicAttributes" \
    "aws_noauth sns get-topic-attributes --topic-arn '$FIFO_TOPIC_ARN' --query 'Attributes.FifoTopic' --output text" \
    "true"

FIFO_MSG_ID=$(aws_noauth sns publish \
    --topic-arn "$FIFO_TOPIC_ARN" \
    --message "Test FIFO message" \
    --message-group-id "group1" \
    --message-deduplication-id "dedup1" \
    --query 'MessageId' --output text 2>&1)
run_test "19. PublishFifoWithMessageGroupId" \
    "echo '$FIFO_MSG_ID' | grep -v '^$' | grep -v '^None$'" \
    ""

FIFO_MSG_ID2=$(aws_noauth sns publish \
    --topic-arn "$FIFO_TOPIC_ARN" \
    --message "Test FIFO message 2" \
    --message-group-id "group1" \
    --message-deduplication-id "dedup1" \
    --query 'MessageId' --output text 2>&1)
run_test "20. FifoDeduplication" \
    "echo '$FIFO_MSG_ID2'" \
    "$FIFO_MSG_ID"

run_test "21. RejectNoMessageGroupId" \
    "aws_noauth sns publish --topic-arn '$FIFO_TOPIC_ARN' --message 'Test without group ID' 2>&1" \
    "invalid\|required\|MessageGroupId"

CB_MSG_ID=$(aws_noauth sns publish \
    --topic-arn "$FIFO_TOPIC_ARN" \
    --message "Same content for dedup" \
    --message-group-id "group2" \
    --query 'MessageId' --output text 2>&1)
CB_MSG_ID2=$(aws_noauth sns publish \
    --topic-arn "$FIFO_TOPIC_ARN" \
    --message "Same content for dedup" \
    --message-group-id "group2" \
    --query 'MessageId' --output text 2>&1)
run_test "22. ContentBasedDeduplication" \
    "echo '$CB_MSG_ID2'" \
    "$CB_MSG_ID"

cleanup_fifo

print_summary
