#!/bin/bash
# Kinesis Service Tests for VorpalStacks
# Usage: bash scripts/services/test_kinesis.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Kinesis Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

STREAM_NAME="test-stream-$$"
STREAM_ARN=""

cleanup() {
    aws_noauth kinesis delete-stream --stream-name "$STREAM_NAME" 2>/dev/null || true
}
trap cleanup EXIT

RESULT=$(aws_noauth kinesis create-stream --stream-name "$STREAM_NAME" --shard-count 1 2>&1)
sleep 2
run_test "1. CreateStream" \
    "echo '$RESULT' | grep -v 'error\|Error\|ERROR'" \
    ""

run_test "2. DescribeStream" \
    "aws_noauth kinesis describe-stream --stream-name '$STREAM_NAME' --query 'StreamDescription.StreamName' --output text" \
    "$STREAM_NAME"

run_test "3. ListStreams" \
    "aws_noauth kinesis list-streams --query 'StreamNames[0]' --output text" \
    "$STREAM_NAME"

run_test "4. DescribeStreamSummary" \
    "aws_noauth kinesis describe-stream-summary --stream-name '$STREAM_NAME' --query 'StreamDescriptionSummary.StreamName' --output text" \
    "$STREAM_NAME"

run_test "5. PutRecord" \
    "aws_noauth kinesis put-record --stream-name '$STREAM_NAME' --data 'test-data' --partition-key 'test-key' --query 'SequenceNumber' --output text" \
    ""

run_test "6. PutRecords" \
    "aws_noauth kinesis put-records --stream-name '$STREAM_NAME' --records '[{\"Data\":\"dGVzdC1kYXRh\",\"PartitionKey\":\"test-key\"}]' --query 'Records[0].SequenceNumber' --output text" \
    ""

# Get shard iterator for reading
SHARD_ITERATOR=$(aws_noauth kinesis get-shard-iterator --stream-name "$STREAM_NAME" --shard-id "shardId-000000000000" --shard-iterator-type TRIM_HORIZON --query 'ShardIterator' --output text 2>/dev/null || echo "")

if [ -n "$SHARD_ITERATOR" ] && [ "$SHARD_ITERATOR" != "None" ]; then
    run_test "7. GetShardIterator" \
        "echo '$SHARD_ITERATOR'" \
        ""
    
    run_test "8. GetRecords" \
        "aws_noauth kinesis get-records --shard-iterator '$SHARD_ITERATOR' --query 'Records' --output text" \
        ""
else
    echo -e "${YELLOW}SKIP${NC}: GetShardIterator/GetRecords"
fi

run_test "9. AddTagsToStream" \
    "aws_noauth kinesis add-tags-to-stream --stream-name '$STREAM_NAME' --tags Environment=Test" \
    ""

run_test "10. ListTagsForStream" \
    "aws_noauth kinesis list-tags-for-stream --stream-name '$STREAM_NAME' --query 'Tags[?Key==\`Environment\`].Key' --output text" \
    "Environment"

run_test "11. IncreaseStreamRetentionPeriod" \
    "aws_noauth kinesis increase-stream-retention-period --stream-name '$STREAM_NAME' --retention-period-hours 48" \
    ""

run_test "12. DecreaseStreamRetentionPeriod" \
    "aws_noauth kinesis decrease-stream-retention-period --stream-name '$STREAM_NAME' --retention-period-hours 24" \
    ""

run_test "13. RemoveTagsFromStream" \
    "aws_noauth kinesis remove-tags-from-stream --stream-name '$STREAM_NAME' --tag-keys Environment" \
    ""

run_test "14. DeleteStream" \
    "aws_noauth kinesis delete-stream --stream-name '$STREAM_NAME'" \
    ""

print_summary
