#!/bin/bash
# Lambda Service Tests for VorpalStacks
# Usage: bash scripts/services/test_lambda.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "Lambda Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

FN_NAME="test-lambda-$$"
LAYER_NAME="test-layer-$$"
FUNCTION_ZIP="$PROJECT_ROOT/tmp/lambda-test-$$/function.zip"
RESPONSE_FILE="$PROJECT_ROOT/tmp/lambda-response-$$"

cleanup() {
    aws_noauth lambda delete-function --function-name "$FN_NAME" 2>/dev/null || true
    aws_noauth lambda delete-layer-version --layer-name "$LAYER_NAME" --version-number 1 2>/dev/null || true
    rm -rf "$(dirname "$FUNCTION_ZIP")" "$RESPONSE_FILE" 2>/dev/null || true
}
trap cleanup EXIT

mkdir -p "$(dirname "$FUNCTION_ZIP")"
create_lambda_zip "$FUNCTION_ZIP"

run_test "1. GetAccountSettings" \
    "aws_noauth lambda get-account-settings --query 'AccountLimit.TotalCodeSize' --output text" \
    ""

run_test "2. ListFunctions (empty)" \
    "aws_noauth lambda list-functions --query 'Functions' --output text" \
    ""

run_test "3. CreateFunction" \
    "aws_noauth lambda create-function --function-name '$FN_NAME' --runtime nodejs20.x --role 'arn:aws:iam:$REGION:$ACCOUNT_ID:role/lambda-role' --handler index.handler --zip-file 'fileb://$FUNCTION_ZIP' --query 'FunctionName' --output text" \
    "$FN_NAME"

run_test "4. GetFunction" \
    "aws_noauth lambda get-function --function-name '$FN_NAME' --query 'Configuration.FunctionName' --output text" \
    "$FN_NAME"

run_test "5. GetFunctionConfiguration" \
    "aws_noauth lambda get-function-configuration --function-name '$FN_NAME' --query 'Runtime' --output text" \
    "nodejs20.x"

run_test "6. ListFunctions" \
    "aws_noauth lambda list-functions --query 'Functions[0].FunctionName' --output text" \
    "$FN_NAME"

run_test "7. UpdateFunctionConfiguration" \
    "aws_noauth lambda update-function-configuration --function-name '$FN_NAME' --timeout 60 --memory-size 512 --query 'Timeout' --output text" \
    "60"

run_test "8. UpdateFunctionCode" \
    "aws_noauth lambda update-function-code --function-name '$FN_NAME' --zip-file 'fileb://$FUNCTION_ZIP' --query 'CodeSize' --output text" \
    ""

run_test "9. PublishVersion" \
    "aws_noauth lambda publish-version --function-name '$FN_NAME' --query 'Version' --output text" \
    "1"

run_test "10. ListVersionsByFunction" \
    "aws_noauth lambda list-versions-by-function --function-name '$FN_NAME' --query 'Versions' --output text" \
    "1"

run_test "11. CreateAlias" \
    "aws_noauth lambda create-alias --function-name '$FN_NAME' --name test-alias --function-version '1' --query 'Name' --output text" \
    "test-alias"

run_test "12. GetAlias" \
    "aws_noauth lambda get-alias --function-name '$FN_NAME' --name test-alias --query 'Name' --output text" \
    "test-alias"

run_test "13. ListAliases" \
    "aws_noauth lambda list-aliases --function-name '$FN_NAME' --query 'Aliases[0].Name' --output text" \
    "test-alias"

run_test "14. DeleteAlias" \
    "aws_noauth lambda delete-alias --function-name '$FN_NAME' --name test-alias" \
    ""

run_test "15. TagResource" \
    "aws_noauth lambda tag-resource --resource 'arn:aws:lambda:$REGION:$ACCOUNT_ID:function:$FN_NAME' --tags Environment=Test" \
    ""

run_test "16. ListTags" \
    "aws_noauth lambda list-tags --resource 'arn:aws:lambda:$REGION:$ACCOUNT_ID:function:$FN_NAME' --query 'Tags.Environment' --output text" \
    "Test"

run_test "17. UntagResource" \
    "aws_noauth lambda untag-resource --resource 'arn:aws:lambda:$REGION:$ACCOUNT_ID:function:$FN_NAME' --tag-keys Environment" \
    ""

run_test "18. PublishLayerVersion" \
    "aws_noauth lambda publish-layer-version --layer-name '$LAYER_NAME' --compatible-runtimes nodejs20.x --zip-file 'fileb://$FUNCTION_ZIP' --query 'LayerVersionArn' --output text" \
    "$LAYER_NAME"

run_test "19. ListLayers" \
    "aws_noauth lambda list-layers --query 'Layers[?LayerName==\`$LAYER_NAME\`].LayerName' --output text" \
    "$LAYER_NAME"

run_test "20. DeleteLayerVersion" \
    "aws_noauth lambda delete-layer-version --layer-name '$LAYER_NAME' --version-number 1" \
    ""

run_test "21. DeleteFunction" \
    "aws_noauth lambda delete-function --function-name '$FN_NAME'" \
    ""

print_summary
