#!/bin/bash
# KMS Service Tests for VorpalStacks
# Usage: bash scripts/services/test_kms.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "KMS Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

KEY_ALIAS="alias/test-key-$$"
KEY_ID=""

cleanup() {
    aws_cmd kms delete-alias --alias-name "$KEY_ALIAS" 2>/dev/null || true
    [ -n "$KEY_ID" ] && aws_cmd kms schedule-key-deletion --key-id "$KEY_ID" --pending-window-in-days 7 2>/dev/null || true
}
trap cleanup EXIT

RESULT=$(aws_cmd kms create-key --description "Test key" --query 'KeyMetadata.KeyId' --output text 2>&1)
KEY_ID="$RESULT"
run_test "1. CreateKey" \
    "echo '$KEY_ID' | grep -v '^$' | grep -v '^None$'" \
    ""

run_test "2. DescribeKey" \
    "aws_cmd kms describe-key --key-id '$KEY_ID' --query 'KeyMetadata.KeyState' --output text" \
    "Enabled"

run_test "3. ListKeys" \
    "aws_cmd kms list-keys --query 'Keys[0].KeyId' --output text" \
    ""

run_test "4. CreateAlias" \
    "aws_cmd kms create-alias --alias-name '$KEY_ALIAS' --target-key-id '$KEY_ID'" \
    ""

run_test "5. ListAliases" \
    "aws_cmd kms list-aliases --query 'Aliases[?AliasName==\`$KEY_ALIAS\`].AliasName' --output text" \
    "$KEY_ALIAS"

run_test "6. EnableKeyRotation" \
    "aws_cmd kms enable-key-rotation --key-id '$KEY_ID'" \
    ""

run_test "7. GetKeyRotationStatus" \
    "aws_cmd kms get-key-rotation-status --key-id '$KEY_ID' --query 'KeyRotationEnabled' --output text" \
    "True"

run_test "8. PutKeyPolicy" \
    "aws_cmd kms put-key-policy --key-id '$KEY_ID' --policy-name default --policy '{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"*\"},\"Action\":\"kms:*\",\"Resource\":\"*\"}]}'" \
    ""

run_test "9. GetKeyPolicy" \
    "aws_cmd kms get-key-policy --key-id '$KEY_ID' --policy-name default" \
    "Version"

run_test "10. Encrypt" \
    "aws_cmd kms encrypt --key-id '$KEY_ID' --plaintext 'test-data' --query 'CiphertextBlob' --output text" \
    ""

run_test "11. DisableKey" \
    "aws_cmd kms disable-key --key-id '$KEY_ID'" \
    ""

run_test "12. EnableKey" \
    "aws_cmd kms enable-key --key-id '$KEY_ID'" \
    ""

run_test "13. DeleteAlias" \
    "aws_cmd kms delete-alias --alias-name '$KEY_ALIAS'" \
    ""

run_test "14. ScheduleKeyDeletion" \
    "aws_cmd kms schedule-key-deletion --key-id '$KEY_ID' --pending-window-in-days 7 --query 'KeyState' --output text" \
    "PendingDeletion"

print_summary
