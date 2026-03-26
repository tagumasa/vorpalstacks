#!/bin/bash
# API Gateway Full Flow Test
# Tests API Gateway -> Lambda/SQS/SNS integration

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "API Gateway Flow Test"
echo ""

API_NAME="test-api-$$"
API_ID=""
RESOURCE_ID=""
ROOT_RESOURCE_ID=""
STAGE_NAME="test"
DEPLOYMENT_ID=""

cleanup() {
    [ -n "$API_ID" ] && aws_noauth apigateway delete-rest-api --rest-api-id "$API_ID" 2>/dev/null || true
}
trap cleanup EXIT

# Test 1: CreateRestApi
echo -n "Test 1: CreateRestApi... "
API_ID=$(aws_noauth apigateway create-rest-api --name "$API_NAME" --query 'id' --output text 2>&1)
if [ -n "$API_ID" ] && [ "$API_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC} (API ID: $API_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $API_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

# Test 2: GetResources (root)
echo -n "Test 2: GetResources (root)... "
ROOT_RESOURCE_ID=$(aws_noauth apigateway get-resources --rest-api-id "$API_ID" --query 'items[0].id' --output text 2>&1)
if [ -n "$ROOT_RESOURCE_ID" ] && [ "$ROOT_RESOURCE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC} (Root: $ROOT_RESOURCE_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $ROOT_RESOURCE_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

# Test 3: CreateResource
echo -n "Test 3: CreateResource... "
RESOURCE_ID=$(aws_noauth apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ROOT_RESOURCE_ID" --path-part 'test' --query 'id' --output text 2>&1)
if [ -n "$RESOURCE_ID" ] && [ "$RESOURCE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC} (Resource: $RESOURCE_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $RESOURCE_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "4. PutMethod (GET)" \
    "aws_noauth apigateway put-method --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --authorization-type NONE --query 'httpMethod' --output text" \
    "GET"

run_test "5. PutIntegration (MOCK)" \
    "aws_noauth apigateway put-integration --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --type MOCK --query 'type' --output text" \
    "MOCK"

run_test "6. PutMethodResponse" \
    "aws_noauth apigateway put-method-response --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --status-code 200" \
    ""

run_test "7. PutIntegrationResponse" \
    "aws_noauth apigateway put-integration-response --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --status-code 200" \
    ""

# Test 8: CreateDeployment
echo -n "Test 8: CreateDeployment... "
DEPLOYMENT_ID=$(aws_noauth apigateway create-deployment --rest-api-id "$API_ID" --query 'id' --output text 2>&1)
if [ -n "$DEPLOYMENT_ID" ] && [ "$DEPLOYMENT_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC} (Deployment: $DEPLOYMENT_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $DEPLOYMENT_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "9. CreateStage" \
    "aws_noauth apigateway create-stage --rest-api-id '$API_ID' --stage-name '$STAGE_NAME' --deployment-id '$DEPLOYMENT_ID' --query 'stageName' --output text" \
    "$STAGE_NAME"

run_test "10. GetStage" \
    "aws_noauth apigateway get-stage --rest-api-id '$API_ID' --stage-name '$STAGE_NAME' --query 'stageName' --output text" \
    "$STAGE_NAME"

run_test "11. CreateModel" \
    "aws_noauth apigateway create-model --rest-api-id '$API_ID' --name 'TestModel' --content-type 'application/json' --schema '{\"type\":\"object\"}' --query 'name' --output text" \
    "TestModel"

run_test "12. GetModels" \
    "aws_noauth apigateway get-models --rest-api-id '$API_ID' --query 'items[0].name' --output text" \
    "TestModel"

run_test "13. CreateApiKey" \
    "aws_noauth apigateway create-api-key --name 'test-key' --enabled --query 'name' --output text" \
    "test-key"

run_test "14. GetApiKeys" \
    "aws_noauth apigateway get-api-keys --query 'items[0].name' --output text" \
    "test-key"

run_test "15. CreateUsagePlan" \
    "aws_noauth apigateway create-usage-plan --name 'test-plan' --query 'name' --output text" \
    "test-plan"

run_test "16. DeleteStage" \
    "aws_noauth apigateway delete-stage --rest-api-id '$API_ID' --stage-name '$STAGE_NAME'" \
    ""

run_test "17. DeleteRestApi" \
    "aws_noauth apigateway delete-rest-api --rest-api-id '$API_ID'" \
    ""

print_summary
