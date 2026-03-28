#!/bin/bash
# API Gateway → Lambda → S3 Integration Test
# Tests full request path: HTTP → API Gateway → Lambda invocation → S3 PutObject → S3 GetObject

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "API Gateway → Lambda → S3 Integration Test"
echo ""

API_NAME="integ-apigw-s3-$$"
FN_NAME="integ-apigw-s3-fn-$$"
BUCKET_NAME="integ-apigw-s3-$$"
API_ID=""
RESOURCE_ID=""
ROOT_RESOURCE_ID=""
STAGE_NAME="test"
DEPLOYMENT_ID=""
FN_ZIP="$PROJECT_ROOT/tmp/lambda-apigw-s3-$$/function.zip"
TEST_KEY="test-object-$(date +%s).txt"

cleanup() {
    [ -n "$API_ID" ] && aws_noauth apigateway delete-rest-api --rest-api-id "$API_ID" 2>/dev/null || true
    aws_noauth lambda delete-function --function-name "$FN_NAME" 2>/dev/null || true
    aws_noauth s3 rb "s3://$BUCKET_NAME" --force 2>/dev/null || true
    rm -rf "$(dirname "$FN_ZIP")" 2>/dev/null || true
}
trap cleanup EXIT

mkdir -p "$(dirname "$FN_ZIP")"

S3_HANDLER='const { S3Client, PutObjectCommand, GetObjectCommand } = require("@aws-sdk/client-s3"); const s3 = new S3Client({ endpoint: process.env.AWS_ENDPOINT_URL || "http://localhost:8080", region: "us-east-1", forcePathStyle: true }); exports.handler = async (event) => { const body = JSON.parse(event.body || "{}"); const key = body.key || "default.txt"; const content = body.content || "hello"; if (event.httpMethod === "GET") { const result = await s3.send(new GetObjectCommand({ Bucket: process.env.BUCKET_NAME, Key: key })); const str = await result.Body.transformToString("utf-8"); return { statusCode: 200, body: JSON.stringify({ key, content: str, operation: "read" }) }; } else { await s3.send(new PutObjectCommand({ Bucket: process.env.BUCKET_NAME, Key: key, Body: content })); return { statusCode: 200, body: JSON.stringify({ key, status: "uploaded", operation: "write" }) }; } };'

echo "$S3_HANDLER" > /tmp/apigw-s3-handler-$$.js
(cd /tmp && zip -q -r "$FN_ZIP" apigw-s3-handler-$$.js)

run_test "1. Create S3 Bucket" \
    "aws_noauth s3 mb 's3://$BUCKET_NAME' 2>&1" \
    "make_bucket"

run_test "2. Create Lambda Function" \
    "aws_noauth lambda create-function --function-name '$FN_NAME' --runtime nodejs20.x --role 'arn:aws:iam::$ACCOUNT_ID:role/lambda-role' --handler apigw-s3-handler-$$.handler --zip-file 'fileb://$FN_ZIP' --environment Variables='{BUCKET_NAME=$BUCKET_NAME,AWS_ENDPOINT_URL=$ENDPOINT_URL}' --query 'FunctionName' --output text" \
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

echo -n "5. Create /s3/{key} Resource... "
RESOURCE_ID=$(aws_noauth apigateway create-resource --rest-api-id "$API_ID" --parent-id "$ROOT_RESOURCE_ID" --path-part '{key}' --query 'id' --output text 2>&1)
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

run_test "6. PutMethod (PUT /s3/{key})" \
    "aws_noauth apigateway put-method --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method PUT --authorization-type NONE --query 'httpMethod' --output text" \
    "PUT"

run_test "7. PutIntegration (LAMBDA_PROXY, PUT)" \
    "aws_noauth apigateway put-integration --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method PUT --type AWS_PROXY --integration-http-method POST --uri 'arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/$FN_ARN/invocations' --query 'type' --output text" \
    "AWS_PROXY"

run_test "8. PutMethod (GET /s3/{key})" \
    "aws_noauth apigateway put-method --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --authorization-type NONE --query 'httpMethod' --output text" \
    "GET"

run_test "9. PutIntegration (LAMBDA_PROXY, GET)" \
    "aws_noauth apigateway put-integration --rest-api-id '$API_ID' --resource-id '$RESOURCE_ID' --http-method GET --type AWS_PROXY --integration-http-method POST --uri 'arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/$FN_ARN/invocations' --query 'type' --output text" \
    "AWS_PROXY"

echo -n "10. Create Deployment... "
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

run_test "11. Create Stage" \
    "aws_noauth apigateway create-stage --rest-api-id '$API_ID' --stage-name '$STAGE_NAME' --deployment-id '$DEPLOYMENT_ID' --query 'stageName' --output text" \
    "$STAGE_NAME"

INVOKE_URL="$ENDPOINT_URL/restapis/$API_ID/$STAGE_NAME/_user_request_/s3/$TEST_KEY"

echo -n "12. Invoke API PUT (write to S3)... "
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$INVOKE_URL" \
    -H "Content-Type: application/json" \
    -d "{\"content\":\"apigw-lambda-s3-integration\"}" 2>&1)
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ] && echo "$BODY" | grep -q "uploaded"; then
    echo -e "${GREEN}PASS${NC} (HTTP $HTTP_CODE, body: ${BODY:0:100})"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  HTTP: $HTTP_CODE"
    echo "  Body: ${BODY:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

echo -n "13. Invoke API GET (read from S3)... "
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$INVOKE_URL" 2>&1)
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ] && echo "$BODY" | grep -q "apigw-lambda-s3-integration"; then
    echo -e "${GREEN}PASS${NC} (HTTP $HTTP_CODE, body: ${BODY:0:100})"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  HTTP: $HTTP_CODE"
    echo "  Body: ${BODY:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "14. Verify S3 Object via AWS CLI" \
    "aws_noauth s3api get-object --bucket '$BUCKET_NAME' --key '$TEST_KEY' /dev/stdout 2>/dev/null" \
    "apigw-lambda-s3-integration"

print_summary
