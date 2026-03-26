#!/bin/bash
# SQS Service Tests for VorpalStacks
# Usage: bash scripts/services/test_sqs.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "SQS Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

QUEUE_NAME="test-queue-$$"
FIFO_QUEUE_NAME="test-fifo-queue-$$.fifo"
QUEUE_URL=""
FIFO_QUEUE_URL=""

cleanup() {
    [ -n "$QUEUE_URL" ] && aws_noauth sqs delete-queue --queue-url "$QUEUE_URL" 2>/dev/null || true
    [ -n "$FIFO_QUEUE_URL" ] && aws_noauth sqs delete-queue --queue-url "$FIFO_QUEUE_URL" 2>/dev/null || true
}
trap cleanup EXIT

QUEUE_URL=$(aws_noauth sqs create-queue --queue-name "$QUEUE_NAME" --query 'QueueUrl' --output text 2>&1)
run_test "1. CreateQueue" \
    "echo '$QUEUE_URL'" \
    "$QUEUE_NAME"

run_test "2. ListQueues" \
    "aws_noauth sqs list-queues --query 'QueueUrls[?contains(@,\`$QUEUE_NAME\`)] | [0]' --output text" \
    "$QUEUE_NAME"

run_test "3. GetQueueUrl" \
    "aws_noauth sqs get-queue-url --queue-name $QUEUE_NAME --query 'QueueUrl' --output text" \
    "$QUEUE_NAME"

run_test "4. GetQueueAttributes" \
    "aws_noauth sqs get-queue-attributes --queue-url '$QUEUE_URL' --attribute-names All" \
    "QueueArn"

run_test "5. SetQueueAttributes" \
    "aws_noauth sqs set-queue-attributes --queue-url '$QUEUE_URL' --attributes VisibilityTimeout=30" \
    ""

run_test "6. SendMessage" \
    "aws_noauth sqs send-message --queue-url '$QUEUE_URL' --message-body '$SQS_MESSAGE_BODY'" \
    "MessageId"

run_test "7. SendMessageBatch" \
    "aws_noauth sqs send-message-batch --queue-url '$QUEUE_URL' --entries '[{\"Id\":\"msg1\",\"MessageBody\":\"Batch 1\"},{\"Id\":\"msg2\",\"MessageBody\":\"Batch 2\"}]'" \
    "Successful"

run_test "8. ReceiveMessage" \
    "aws_noauth sqs receive-message --queue-url '$QUEUE_URL' --max-number-of-messages 1" \
    "Body"

RECEIPT_HANDLE=$(aws_noauth sqs receive-message --queue-url "$QUEUE_URL" --max-number-of-messages 1 --query 'Messages[0].ReceiptHandle' --output text 2>/dev/null || echo "")

if [ -n "$RECEIPT_HANDLE" ] && [ "$RECEIPT_HANDLE" != "None" ]; then
    run_test "9. DeleteMessage" \
        "aws_noauth sqs delete-message --queue-url '$QUEUE_URL' --receipt-handle '$RECEIPT_HANDLE'" \
        ""
else
    echo -e "${YELLOW}SKIP${NC}: DeleteMessage (no message received)"
fi

run_test "10. TagQueue" \
    "aws_noauth sqs tag-queue --queue-url '$QUEUE_URL' --tags Environment=Test" \
    ""

run_test "11. ListQueueTags" \
    "aws_noauth sqs list-queue-tags --queue-url '$QUEUE_URL'" \
    "Environment"

run_test "12. UntagQueue" \
    "aws_noauth sqs untag-queue --queue-url '$QUEUE_URL' --tag-keys Environment" \
    ""

run_test "13. PurgeQueue" \
    "aws_noauth sqs purge-queue --queue-url '$QUEUE_URL'" \
    ""

run_test "14. DeleteQueue" \
    "aws_noauth sqs delete-queue --queue-url '$QUEUE_URL'" \
    ""

# FIFO Queue Tests
echo ""
echo "=== FIFO Queue Tests ==="

FIFO_QUEUE_URL=$(aws_noauth sqs create-queue --queue-name "$FIFO_QUEUE_NAME" --attributes '{"FifoQueue":"true","ContentBasedDeduplication":"true"}' --query 'QueueUrl' --output text 2>&1)
run_test "15. CreateFifoQueue" \
    "echo '$FIFO_QUEUE_URL'" \
    ".fifo"

run_test "16. FifoQueueAttributes" \
    "aws_noauth sqs get-queue-attributes --queue-url '$FIFO_QUEUE_URL' --attribute-names FifoQueue ContentBasedDeduplication --query 'Attributes.FifoQueue' --output text" \
    "true"

FIFO_MSG_ID1=$(aws_noauth sqs send-message --queue-url "$FIFO_QUEUE_URL" --message-body "FIFO test message" --message-group-id "group1" --message-deduplication-id "dup1" --query 'MessageId' --output text 2>&1)
run_test "17. SendFifoMessage" \
    "echo '$FIFO_MSG_ID1' | grep -v '^$' | grep -v '^None$'" \
    ""

FIFO_MSG_ID2=$(aws_noauth sqs send-message --queue-url "$FIFO_QUEUE_URL" --message-body "FIFO test message" --message-group-id "group1" --message-deduplication-id "dup1" --query 'MessageId' --output text 2>&1)
run_test "18. FifoDeduplication" \
    "echo '$FIFO_MSG_ID2'" \
    "$FIFO_MSG_ID1"

run_test "19. ReceiveFifoMessageGroupId" \
    "aws_noauth sqs receive-message --queue-url '$FIFO_QUEUE_URL' --attribute-names 'MessageGroupId,MessageDeduplicationId' --query 'Messages[0].Attributes.MessageGroupId' --output text" \
    "group1"

run_test "20. RejectNoMessageGroupId" \
    "aws_noauth sqs send-message --queue-url '$FIFO_QUEUE_URL' --message-body 'no group' 2>&1" \
    "MessageGroupId"

FIFO_MSG_ID3=$(aws_noauth sqs send-message --queue-url "$FIFO_QUEUE_URL" --message-body "no dedup id" --message-group-id "group2" --query 'MessageId' --output text 2>&1)
run_test "21. AcceptNoDedupIdWithContentBased" \
    "echo '$FIFO_MSG_ID3' | grep -v '^$' | grep -v '^None$'" \
    ""

print_summary
