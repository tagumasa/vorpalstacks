#!/bin/bash
# API Gateway Service Tests for VorpalStacks
# Usage: bash scripts/services/test_apigateway.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "API Gateway Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

API_NAME="test-api-$$"
REST_API_ID=""
RESOURCE_ID=""

cleanup() {
    [ -n "$REST_API_ID" ] && aws_noauth apigateway delete-rest-api --rest-api-id "$REST_API_ID" 2>/dev/null || true
}
trap cleanup EXIT

REST_API_ID=$(aws_noauth apigateway create-rest-api --name "$API_NAME" --query 'id' --output text 2>/dev/null || echo "")

if [ -n "$REST_API_ID" ] && [ "$REST_API_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateRestApi (ID: $REST_API_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateRestApi"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "2. GetRestApi" \
    "aws_noauth apigateway get-rest-api --rest-api-id '$REST_API_ID' --query 'name' --output text" \
    "$API_NAME"

run_test "3. ListRestApis" \
    "aws_noauth apigateway get-rest-apis --query 'items[?name==\`$API_NAME\`].name' --output text" \
    "$API_NAME"

RESOURCE_ID=$(aws_noauth apigateway get-resources --rest-api-id "$REST_API_ID" --query 'items[0].id' --output text 2>/dev/null || echo "")

if [ -n "$RESOURCE_ID" ] && [ "$RESOURCE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 4. GetResources (root: $RESOURCE_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 4. GetResources"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "5. CreateResource" \
    "aws_noauth apigateway create-resource --rest-api-id '$REST_API_ID' --parent-id '$RESOURCE_ID' --path-part 'test' --query 'id' --output text" \
    ""

NEW_RESOURCE_ID=$(aws_noauth apigateway create-resource --rest-api-id "$REST_API_ID" --parent-id "$RESOURCE_ID" --path-part "items" --query 'id' --output text 2>/dev/null || echo "")

if [ -n "$NEW_RESOURCE_ID" ] && [ "$NEW_RESOURCE_ID" != "None" ]; then
    run_test "6. PutMethod" \
        "aws_noauth apigateway put-method --rest-api-id '$REST_API_ID' --resource-id '$NEW_RESOURCE_ID' --http-method GET --authorization-type NONE --query 'httpMethod' --output text" \
        "GET"

    run_test "7. GetMethod" \
        "aws_noauth apigateway get-method --rest-api-id '$REST_API_ID' --resource-id '$NEW_RESOURCE_ID' --http-method GET --query 'httpMethod' --output text" \
        "GET"

    run_test "8. PutIntegration" \
        "aws_noauth apigateway put-integration --rest-api-id '$REST_API_ID' --resource-id '$NEW_RESOURCE_ID' --http-method GET --type MOCK --request-templates '{\"application/json\":\"{\\\"statusCode\\\": 200}\"}'" \
        ""

    DEPLOY_ID=$(aws_noauth apigateway create-deployment --rest-api-id "$REST_API_ID" --description "Test deployment" --query 'id' --output text 2>/dev/null || echo "")

    if [ -n "$DEPLOY_ID" ] && [ "$DEPLOY_ID" != "None" ]; then
        echo -e "${GREEN}PASS${NC}: 9. CreateDeployment (ID: $DEPLOY_ID)"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}FAIL${NC}: 9. CreateDeployment"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    TEST_COUNT=$((TEST_COUNT + 1))

    run_test "10. GetDeployments" \
        "aws_noauth apigateway get-deployments --rest-api-id '$REST_API_ID' --query 'items | length(@)' --output text" \
        "1"

    run_test "11. DeleteMethod" \
        "aws_noauth apigateway delete-method --rest-api-id '$REST_API_ID' --resource-id '$NEW_RESOURCE_ID' --http-method GET" \
        ""

    run_test "12. DeleteResource" \
        "aws_noauth apigateway delete-resource --rest-api-id '$REST_API_ID' --resource-id '$NEW_RESOURCE_ID'" \
        ""
fi

run_test "13. TagResource" \
    "aws_noauth apigateway tag-resource --resource-arn 'arn:aws:apigateway:$REGION::$ACCOUNT_ID:/restapis/$REST_API_ID' --tags Environment=Test" \
    ""

run_test "14. DeleteRestApi" \
    "aws_noauth apigateway delete-rest-api --rest-api-id '$REST_API_ID'" \
    ""

REST_API_ID=""

print_summary
