#!/bin/bash
# Reproduction script for identified DynamoDB bugs

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "DynamoDB Bug Reproduction Tests"

TABLE_NAME="BugReproTable-$$"

cleanup() {
    echo "Cleaning up..."
    aws_noauth dynamodb delete-table --table-name "$TABLE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

echo "1. Creating table..."
aws_noauth dynamodb create-table --table-name "$TABLE_NAME" --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST > /dev/null

echo "2. Putting item with nested map..."
aws_noauth dynamodb put-item --table-name "$TABLE_NAME" --item '{"id":{"S":"test1"},"info":{"M":{"name":{"S":"John"},"age":{"N":"30"}}}}'

echo "3. Testing ProjectionExpression with nested path..."
# This is expected to FAIL if the bug exists
RESULT=$(aws_noauth dynamodb get-item --table-name "$TABLE_NAME" --key '{"id":{"S":"test1"}}' --projection-expression "info.name" --query 'Item.info.M.name.S' --output text)
if [ "$RESULT" == "John" ]; then
    echo -e "${GREEN}PASS: ProjectionExpression with nested path works${NC}"
else
    echo -e "${RED}FAIL: ProjectionExpression with nested path failed (Result: $RESULT)${NC}"
fi

echo "4. Testing ConditionExpression with nested path..."
# This is expected to FAIL if the bug exists
aws_noauth dynamodb put-item --table-name "$TABLE_NAME" --item '{"id":{"S":"test1"},"info":{"M":{"name":{"S":"Jane"},"age":{"N":"31"}}}}' --condition-expression "info.name = :val" --expression-attribute-values '{":val":{"S":"John"}}' 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS: ConditionExpression with nested path worked${NC}"
else
    echo -e "${RED}FAIL: ConditionExpression with nested path failed${NC}"
fi

echo "5. Testing ADD action on NS with normalization inconsistency..."
aws_noauth dynamodb put-item --table-name "$TABLE_NAME" --item '{"id":{"S":"test-ns"},"nums":{"NS":["1.0"]}}'
aws_noauth dynamodb update-item --table-name "$TABLE_NAME" --key '{"id":{"S":"test-ns"}}' --update-expression 'ADD nums :n' --expression-attribute-values '{":n":{"NS":["1"]}}'
RESULT=$(aws_noauth dynamodb get-item --table-name "$TABLE_NAME" --key '{"id":{"S":"test-ns"}}' --query 'length(Item.nums.NS)' --output text)
if [ "$RESULT" == "1" ]; then
    echo -e "${GREEN}PASS: ADD on NS handles normalization (Result: $RESULT)${NC}"
else
    echo -e "${RED}FAIL: ADD on NS might have duplicates (Result: $RESULT)${NC}"
fi

echo "6. Testing Binary Set (BS) with invalid base64..."
# AWS CLI might handle base64 encoding, so we might need to be careful how we send "invalid base64"
# Actually, if we send a string that is NOT base64 but labeled as BS, it should fail or be handled consistently.
# But let's skip this for now as it's harder to test with CLI.

print_summary
