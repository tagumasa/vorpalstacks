#!/bin/bash
# AppSync Service Tests for VorpalStacks
# Usage: bash scripts/services/test_appsync.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "AppSync Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

API_NAME="test-api-$$"
API_ID=""
API_ARN=""
TAG_KEY="testkey"
TAG_VALUE="testvalue"

cleanup() {
    [ -n "$API_ID" ] && aws_noauth appsync delete-graphql-api --api-id "$API_ID" 2>/dev/null || true
}
trap cleanup EXIT

API_ID=$(aws_noauth appsync create-graphql-api \
    --name "$API_NAME" \
    --authentication-type API_KEY \
    --query 'graphqlApi.apiId' --output text 2>/dev/null || echo "")

if [ -n "$API_ID" ] && [ "$API_ID" != "None" ] && [ "$API_ID" != "" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateGraphqlApi (ID: ${API_ID:0:8}...)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateGraphqlApi"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

API_ARN=$(aws_noauth appsync get-graphql-api --api-id "$API_ID" --query 'graphqlApi.arn' --output text 2>/dev/null || echo "")

run_test "2. GetGraphqlApi" \
    "aws_noauth appsync get-graphql-api --api-id '$API_ID' --query 'graphqlApi.name' --output text" \
    "$API_NAME"

run_test "3. ListGraphqlApis" \
    "aws_noauth appsync list-graphql-apis --query 'graphqlApis[?name==\`$API_NAME\`].apiId' --output text" \
    "$API_ID"

run_test "4. UpdateGraphqlApi" \
    "aws_noauth appsync update-graphql-api --api-id '$API_ID' --name '${API_NAME}-updated' --authentication-type API_KEY --query 'graphqlApi.name' --output text" \
    "${API_NAME}-updated"

run_test "5. CreateApiKey" \
    "aws_noauth appsync create-api-key --api-id '$API_ID' --query 'apiKey.id' --output text" \
    ""

run_test "6. ListApiKeys" \
    "aws_noauth appsync list-api-keys --api-id '$API_ID' --query 'apiKeys | length(@)' --output text" \
    ""

run_test "7. TagResource" \
    "aws_noauth appsync tag-resource --resource-arn '$API_ARN' --tags Key=$TAG_KEY,Value=$TAG_VALUE" \
    ""

run_test "8. ListTagsForResource" \
    "aws_noauth appsync list-tags-for-resource --resource-arn '$API_ARN' --query 'tags.$TAG_KEY' --output text" \
    "$TAG_VALUE"

run_test "9. UntagResource" \
    "aws_noauth appsync untag-resource --resource-arn '$API_ARN' --tag-keys '$TAG_KEY'" \
    ""

run_test "10. ListTypes" \
    "aws_noauth appsync list-types --api-id '$API_ID' --format SDL --query 'types | length(@)' --output text" \
    ""

run_test "11. DeleteGraphqlApi" \
    "aws_noauth appsync delete-graphql-api --api-id '$API_ID'" \
    ""
API_ID=""

print_summary
