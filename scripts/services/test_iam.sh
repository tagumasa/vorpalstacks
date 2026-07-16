#!/bin/bash
# IAM Service Tests for VorpalStacks
# Usage: bash scripts/services/test_iam.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "IAM Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

USER_NAME="test-user-$$"
ROLE_NAME="test-role-$$"
POLICY_NAME="test-policy-$$"
GROUP_NAME="test-group-$$"

cleanup() {
    aws_cmd iam delete-user --user-name "$USER_NAME" 2>/dev/null || true
    aws_cmd iam delete-role --role-name "$ROLE_NAME" 2>/dev/null || true
    aws_cmd iam delete-group --group-name "$GROUP_NAME" 2>/dev/null || true
    local POLICY_ARN="arn:aws:iam::$ACCOUNT_ID:policy/$POLICY_NAME"
    aws_cmd iam delete-policy --policy-arn "$POLICY_ARN" 2>/dev/null || true
}
trap cleanup EXIT

RESULT=$(aws_cmd iam create-user --user-name "$USER_NAME" --query 'User.UserName' --output text 2>&1)
run_test "1. CreateUser" \
    "echo '$RESULT'" \
    "$USER_NAME"

run_test "2. GetUser" \
    "aws_cmd iam get-user --user-name '$USER_NAME' --query 'User.UserName' --output text" \
    "$USER_NAME"

run_test "3. ListUsers" \
    "aws_cmd iam list-users --query 'Users[?UserName==\`$USER_NAME\`].UserName' --output text" \
    "$USER_NAME"

run_test "4. CreateRole" \
    "aws_cmd iam create-role --role-name '$ROLE_NAME' --assume-role-policy-document '{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"lambda.amazonaws.com\"},\"Action\":\"sts:AssumeRole\"}]}' --query 'Role.RoleName' --output text" \
    "$ROLE_NAME"

run_test "5. GetRole" \
    "aws_cmd iam get-role --role-name '$ROLE_NAME' --query 'Role.RoleName' --output text" \
    "$ROLE_NAME"

run_test "6. ListRoles" \
    "aws_cmd iam list-roles --query 'Roles[?RoleName==\`$ROLE_NAME\`].RoleName' --output text" \
    "$ROLE_NAME"

run_test "7. CreateGroup" \
    "aws_cmd iam create-group --group-name '$GROUP_NAME' --query 'Group.GroupName' --output text" \
    "$GROUP_NAME"

run_test "8. ListGroups" \
    "aws_cmd iam list-groups --query 'Groups[?GroupName==\`$GROUP_NAME\`].GroupName' --output text" \
    "$GROUP_NAME"

run_test "9. AddUserToGroup" \
    "aws_cmd iam add-user-to-group --group-name '$GROUP_NAME' --user-name '$USER_NAME'" \
    ""

run_test "10. ListGroupsForUser" \
    "aws_cmd iam list-groups-for-user --user-name '$USER_NAME' --query 'Groups[0].GroupName' --output text" \
    "$GROUP_NAME"

POLICY_ARN=$(aws_cmd iam create-policy --policy-name "$POLICY_NAME" --policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}' --query 'Policy.Arn' --output text 2>&1)

run_test "11. CreatePolicy" \
    "echo '$POLICY_ARN'" \
    "arn:aws:iam::"

run_test "12. ListPolicies" \
    "aws_cmd iam list-policies --scope Local --query 'Policies[?PolicyName==\`$POLICY_NAME\`].PolicyName' --output text" \
    "$POLICY_NAME"

run_test "13. AttachUserPolicy" \
    "aws_cmd iam attach-user-policy --user-name '$USER_NAME' --policy-arn '$POLICY_ARN'" \
    ""

run_test "14. ListAttachedUserPolicies" \
    "aws_cmd iam list-attached-user-policies --user-name '$USER_NAME' --query 'AttachedPolicies[0].PolicyName' --output text" \
    "$POLICY_NAME"

run_test "15. TagUser" \
    "aws_cmd iam tag-user --user-name '$USER_NAME' --tags Key=Environment,Value=Test" \
    ""

run_test "16. ListUserTags" \
    "aws_cmd iam list-user-tags --user-name '$USER_NAME' --query 'Tags[0].Value' --output text" \
    "Test"

run_test "17. RemoveUserFromGroup" \
    "aws_cmd iam remove-user-from-group --group-name '$GROUP_NAME' --user-name '$USER_NAME'" \
    ""

run_test "18. DeleteUser" \
    "aws_cmd iam delete-user --user-name '$USER_NAME'" \
    ""

run_test "19. DeleteRole" \
    "aws_cmd iam delete-role --role-name '$ROLE_NAME'" \
    ""

run_test "20. DeleteGroup" \
    "aws_cmd iam delete-group --group-name '$GROUP_NAME'" \
    ""

print_summary
