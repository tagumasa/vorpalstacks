#!/bin/bash
# SSM (Systems Manager) Service Tests for VorpalStacks
# Usage: bash scripts/services/test_ssm.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "SSM Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

PARAM_NAME="/test/param-$$"
SECRET_PARAM_NAME="/test/secret-param-$$"

cleanup() {
    aws_noauth ssm delete-parameter --name "$PARAM_NAME" 2>/dev/null || true
    aws_noauth ssm delete-parameter --name "$SECRET_PARAM_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. PutParameter (String)" \
    "aws_noauth ssm put-parameter --name '$PARAM_NAME' --value 'test-value' --type String --query 'Version' --output text" \
    "1"

run_test "2. GetParameter" \
    "aws_noauth ssm get-parameter --name '$PARAM_NAME' --query 'Parameter.Value' --output text" \
    "test-value"

run_test "3. DescribeParameters" \
    "aws_noauth ssm describe-parameters --parameter-filters '[{\"Key\":\"Name\",\"Values\":[\"$PARAM_NAME\"]}]' --query 'Parameters[?Name==\`$PARAM_NAME\`].Name' --output text" \
    "$PARAM_NAME"

run_test "4. PutParameter (update)" \
    "aws_noauth ssm put-parameter --name '$PARAM_NAME' --value 'updated-value' --type String --overwrite --query 'Version' --output text" \
    "2"

run_test "5. GetParameterHistory" \
    "aws_noauth ssm get-parameter-history --name '$PARAM_NAME' --query 'length(Parameters)' --output text" \
    "2"

run_test "6. PutParameter (SecureString)" \
    "aws_noauth ssm put-parameter --name '$SECRET_PARAM_NAME' --value 'secret-value' --type SecureString --query 'Version' --output text" \
    "1"

run_test "7. GetParameter (with decryption)" \
    "aws_noauth ssm get-parameter --name '$SECRET_PARAM_NAME' --with-decryption --query 'Parameter.Value' --output text" \
    "secret-value"

run_test "8. AddTagsToResource" \
    "aws_noauth ssm add-tags-to-resource --resource-type Parameter --resource-id '$PARAM_NAME' --tags '[{\"Key\":\"Environment\",\"Value\":\"Test\"}]'" \
    ""

run_test "9. ListTagsForResource" \
    "aws_noauth ssm list-tags-for-resource --resource-type Parameter --resource-id '$PARAM_NAME' --query 'TagList[0].Value' --output text" \
    "Test"

run_test "10. RemoveTagsFromResource" \
    "aws_noauth ssm remove-tags-from-resource --resource-type Parameter --resource-id '$PARAM_NAME' --tag-keys Environment" \
    ""

run_test "11. DeleteParameter" \
    "aws_noauth ssm delete-parameter --name '$SECRET_PARAM_NAME'" \
    ""

run_test "12. DeleteParameters" \
    "aws_noauth ssm delete-parameters --names '$PARAM_NAME' --query 'DeletedParameters[0]' --output text" \
    "$PARAM_NAME"

print_summary
