#!/bin/bash
# Cognito Service Tests for VorpalStacks
# Usage: bash scripts/services/test_cognito.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Cognito Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

USER_POOL_NAME="test-user-pool-$$"
IDENTITY_POOL_NAME="test-identity-pool-$$"
USER_POOL_ID=""
IDENTITY_POOL_ID=""

cleanup() {
    [ -n "$IDENTITY_POOL_ID" ] && aws_noauth cognito-identity delete-identity-pool --identity-pool-id "$IDENTITY_POOL_ID" 2>/dev/null || true
    [ -n "$USER_POOL_ID" ] && aws_noauth cognito-idp delete-user-pool --user-pool-id "$USER_POOL_ID" 2>/dev/null || true
}
trap cleanup EXIT

RESULT=$(aws_noauth cognito-idp create-user-pool --pool-name "$USER_POOL_NAME" --query 'UserPool.Id' --output text 2>&1)
USER_POOL_ID="$RESULT"
run_test "1. CreateUserPool" \
    "echo '$USER_POOL_ID' | grep -v '^$' | grep -v '^None$'" \
    ""

run_test "2. DescribeUserPool" \
    "aws_noauth cognito-idp describe-user-pool --user-pool-id '$USER_POOL_ID' --query 'UserPool.Name' --output text" \
    "$USER_POOL_NAME"

run_test "3. ListUserPools" \
    "aws_noauth cognito-idp list-user-pools --max-results 10 --query 'UserPools[?Name==\`$USER_POOL_NAME\`].Name' --output text" \
    "$USER_POOL_NAME"

run_test "4. CreateUserPoolClient" \
    "aws_noauth cognito-idp create-user-pool-client --user-pool-id '$USER_POOL_ID' --client-name 'test-client' --query 'UserPoolClient.ClientName' --output text" \
    "test-client"

run_test "5. ListUserPoolClients" \
    "aws_noauth cognito-idp list-user-pool-clients --user-pool-id '$USER_POOL_ID' --query 'UserPoolClients[0].ClientName' --output text" \
    "test-client"

run_test "6. AdminCreateUser" \
    "aws_noauth cognito-idp admin-create-user --user-pool-id '$USER_POOL_ID' --username 'testuser' --message-action SUPPRESS --query 'User.Username' --output text" \
    "testuser"

run_test "7. AdminGetUser" \
    "aws_noauth cognito-idp admin-get-user --user-pool-id '$USER_POOL_ID' --username 'testuser' --query 'Username' --output text" \
    "testuser"

run_test "8. ListUsers" \
    "aws_noauth cognito-idp list-users --user-pool-id '$USER_POOL_ID' --query 'Users[0].Username' --output text" \
    "testuser"

run_test "9. AdminDeleteUser" \
    "aws_noauth cognito-idp admin-delete-user --user-pool-id '$USER_POOL_ID' --username 'testuser'" \
    ""

RESULT=$(aws_noauth cognito-identity create-identity-pool --identity-pool-name "$IDENTITY_POOL_NAME" --allow-unauthenticated-identities --query 'IdentityPoolId' --output text 2>&1)
IDENTITY_POOL_ID="$RESULT"
run_test "10. CreateIdentityPool" \
    "echo '$IDENTITY_POOL_ID' | grep -v '^$' | grep -v '^None$'" \
    ""

run_test "11. DescribeIdentityPool" \
    "aws_noauth cognito-identity describe-identity-pool --identity-pool-id '$IDENTITY_POOL_ID' --query 'IdentityPoolName' --output text" \
    "$IDENTITY_POOL_NAME"

run_test "12. ListIdentityPools" \
    "aws_noauth cognito-identity list-identity-pools --max-results 10 --query 'IdentityPools[?IdentityPoolName==\`$IDENTITY_POOL_NAME\`].IdentityPoolName' --output text" \
    "$IDENTITY_POOL_NAME"

run_test "13. DeleteUserPool" \
    "aws_noauth cognito-idp delete-user-pool --user-pool-id '$USER_POOL_ID'" \
    ""

run_test "14. DeleteIdentityPool" \
    "aws_noauth cognito-identity delete-identity-pool --identity-pool-id '$IDENTITY_POOL_ID'" \
    ""

print_summary
