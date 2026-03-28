#!/bin/bash
# API Gateway → Lambda → DynamoDB Integration Test
# Tests full request path: HTTP → API Gateway → Lambda invocation → DynamoDB write → DynamoDB read

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "API Gateway → Lambda → DynamoDB Integration Test"
echo ""

API_NAME="integ-apigw-ddb-$$"
FN_NAME="integ-apigw-ddb-fn-$$"
TABLE_NAME="integ-apigw-ddb-$$"
API_ID=""
RESOURCE_ID=""
ROOT_RESOURCE_ID=""
STAGE_NAME="test"
DEPLOYMENT_ID=""
FN_ZIP="$PROJECT_ROOT/tmp/lambda-apigw-ddb-$$/function.zip"
TEST_ID=""

cleanup() {
    [ -n "$API_ID" ] && aws_noauth apigateway delete-rest-api --rest-api-id "$API_ID" 2>/dev/null || true
    aws_noauth lambda delete-function --function-name "$FN_NAME" 2>/dev/null || true
    aws_noauth dynamodb delete-table --table-name "$TABLE_NAME" 2>/dev/null || true
    rm -rf "$(dirname "$FN_ZIP")" 2>/dev/null || true
}
trap cleanup EXIT

mkdir -p "$(dirname "$FN_ZIP")"

DDB_HANDLER='const { DynamoDBClient, PutItemCommand } = require("@aws-sdk/client-dynamodb"); const ddb = new DynamoDBClient({ endpoint: process.env.AWS_ENDPOINT_URL || "http://localhost:8080", region: "us-east-1" }); exports.handler = async (event) => { const body = JSON.parse(event.body || "{}"); const id = body.id || Date.now().toString(); await ddb.send(new PutItemCommand({ TableName: process.env.TABLE_NAME, Item: { id: { S: id }, message: { S: body.message || "hello" }, timestamp: { S: new Date().toISOString() } } })); return { statusCode: 200, body: JSON.stringify({ id, status: "written" }) }; };'

echo "$DDB_HANDLER" > /tmp/apigw-ddb-handler-$$.js
(cd /tmp && zip -q -r "$FN_ZIP" apigw-ddb-handler-$$.js)

TEST_ID="test-$(date +%s)"

run_test "1. Create DynamoDB Table" \
    "aws_noauth dynamodb create-table --table-name '$TABLE_NAME' --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --query 'TableDescription.TableStatus' --output text" \
    "ACTIVE"

run_test "2. Create Lambda Function" \
    "aws_noauth lambda create-function --function-name '$FN_NAME' --runtime nodejs20.x --role 'arn:aws:iam::$ACCOUNT_ID:role/lambda-role' --handler apigw-ddb-handler-$$.handler --zip-file 'fileb://$FN_ZIP' --environment Variables='{TABLE_NAME=$TABLE_NAME,AWS_ENDPOINT_URL=$ENDPOINT_URL}' --query 'FunctionName' --output text" \
    "$FN_NAME"

echo -n "3. Create REST API... "
API_ID=$(aws_noauth apigateway create-rest-api --name "$API_NAME" --query 'id' --output text 2>&1)
if [ -n "$API_ID" ] && [ "$API_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC} (API: $API_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $API_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

echo -n "4. Get Root Resource... "
ROOT_RESOURCE_ID=$(aws_noauth apigateway get-resources --rest-api-id "$API_ID" --query 'items[0].id' --output text 2>&1)
if [ -n "$ROOT_RESOURCE_ID" ] && [ "$ROOT_RESOURCE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $ROOT_RESOURCE_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

echo -n "5. Create /write Resource... "
RESOURCE_ID=$(aws_noauth apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ROOT_RESOURCE_ID" --path-part 'write' --query 'id' --output text 2>&1)
if [ -n "$RESOURCE_ID" ] && [ "$RESOURCE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $RESOURCE_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

FN_ARN="arn:aws:lambda:$REGION:$ACCOUNT_ID:function:$FN_NAME"

run_test "6. PutMethod (POST /write)" \
    "aws_noauth apigateway put-method --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method POST --authorization-type NONE --query 'httpMethod' --output text" \
    "POST"

run_test "7. PutIntegration (LAMBDA_PROXY)" \
    "aws_noauth apigateway put-integration --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method POST --type AWS_PROXY --integration-http-method POST --uri 'arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/$FN_ARN/invocations' --query 'type' --output text" \
    "AWS_PROXY"

echo -n "8. Create Deployment... "
DEPLOYMENT_ID=$(aws_noauth apigateway create-deployment --rest-api-id "$API_ID" --query 'id' --output text 2>&1)
if [ -n "$DEPLOYMENT_ID" ] && [ "$DEPLOYMENT_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Got: $DEPLOYMENT_ID"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "9. Create Stage" \
    "aws_noauth apigateway create-stage --rest-api-id '$API_ID' --stage-name '$STAGE_NAME' --deployment-id '$DEPLOYMENT_ID' --query 'stageName' --output text" \
    "$STAGE_NAME"

INVOKE_URL="$ENDPOINT_URL/restapis/$API_ID/$STAGE_NAME/_user_request_/write"

echo -n "10. Invoke API (POST /write)... "
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$INVOKE_URL" \
    -H "Content-Type: application/json" \
    -d "{\"id\":\"$TEST_ID\",\"message\":\"apigw-lambda-ddb-integration\"}" 2>&1)
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ] && echo "$BODY" | grep -q "written"; then
    echo -e "${GREEN}PASS${NC} (HTTP $HTTP_CODE, body: ${BODY:0:100})"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  HTTP: $HTTP_CODE"
    echo "  Body: ${BODY:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

sleep 1

run_test "11. Verify DynamoDB Item Exists" \
    "aws_noauth dynamodb get-item --table-name '$TABLE_NAME' --key '{\"id\":{\"S\":\"$TEST_ID\"}}' --query 'Item.id.S' --output text" \
    "$TEST_ID"

run_test "12. Verify DynamoDB Item Content" \
    "aws_noauth dynamodb get-item --table-name '$TABLE_NAME' --key '{\"id\":{\"S\":\"$TEST_ID\"}}' --query 'Item.message.S' --output text" \
    "apigw-lambda-ddb-integration"

print_summary
