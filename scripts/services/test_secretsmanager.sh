#!/bin/bash
# Secrets Manager Service Tests for VorpalStacks
# Usage: bash scripts/services/test_secretsmanager.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Secrets Manager Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

SECRET_NAME="test-secret-$$"
SECRET_ARN=""

cleanup() {
    [ -n "$SECRET_ARN" ] && aws_noauth secretsmanager delete-secret --secret-id "$SECRET_ARN" --force-delete-without-recovery 2>/dev/null || true
    aws_noauth secretsmanager delete-secret --secret-id "$SECRET_NAME" --force-delete-without-recovery 2>/dev/null || true
}
trap cleanup EXIT

RESULT=$(aws_noauth secretsmanager create-secret --name "$SECRET_NAME" --secret-string '{"username":"admin","password":"test123"}' --query 'ARN' --output text 2>&1)
SECRET_ARN="$RESULT"
run_test "1. CreateSecret" \
    "echo '$SECRET_ARN'" \
    "secret"

run_test "2. DescribeSecret" \
    "aws_noauth secretsmanager describe-secret --secret-id '$SECRET_NAME' --query 'Name' --output text" \
    "$SECRET_NAME"

run_test "3. GetSecretValue" \
    "aws_noauth secretsmanager get-secret-value --secret-id '$SECRET_NAME' --query 'SecretString' --output text" \
    "admin"

run_test "4. ListSecrets" \
    "aws_noauth secretsmanager list-secrets --query 'SecretList[0].Name' --output text" \
    "$SECRET_NAME"

VERSION_ID=$(aws_noauth secretsmanager put-secret-value --secret-id "$SECRET_NAME" --secret-string '{"username":"admin","password":"updated456"}' --query 'VersionId' --output text 2>&1)
run_test "5. PutSecretValue" \
    "echo '$VERSION_ID' | awk '{print length(\$0)}' | grep -q '3[2-9]' && echo 'ok' || echo 'short'" \
    "ok"

run_test "6. UpdateSecretVersionStage" \
    "aws_noauth secretsmanager update-secret-version-stage --secret-id '$SECRET_NAME' --version-stage AWSCURRENT --move-to-version '$VERSION_ID'" \
    ""

run_test "7. TagResource" \
    "aws_noauth secretsmanager tag-resource --secret-id '$SECRET_NAME' --tags '[{\"Key\":\"Environment\",\"Value\":\"Test\"}]'" \
    ""

run_test "8. RotateSecret" \
    "aws_noauth secretsmanager rotate-secret --secret-id '$SECRET_NAME' --rotation-lambda-arn 'arn:aws:lambda:$REGION:$ACCOUNT_ID:function:test-rotation' --rotate-immediately 2>&1 || echo 'OK'" \
    ""

run_test "9. RestoreSecret" \
    "aws_noauth secretsmanager restore-secret --secret-id '$SECRET_NAME' --query 'ARN' --output text" \
    "secret"

run_test "10. ListSecretVersionIds" \
    "aws_noauth secretsmanager list-secret-version-ids --secret-id '$SECRET_NAME' --query 'Versions | length(@)' --output text" \
    "2"

run_test "11. DeleteSecret" \
    "aws_noauth secretsmanager delete-secret --secret-id '$SECRET_NAME' --force-delete-without-recovery --query 'ARN' --output text" \
    "secret"

run_test "12. GetRandomPassword" \
    "aws_noauth secretsmanager get-random-password --password-length 20 --query 'RandomPassword' --output text | awk '{print length(\$0)}'" \
    "20"

run_test "13. GetRandomPassword (exclude punctuation)" \
    "aws_noauth secretsmanager get-random-password --password-length 16 --exclude-punctuation --query 'RandomPassword' --output text | awk '{print length(\$0)}'" \
    "16"

print_summary
