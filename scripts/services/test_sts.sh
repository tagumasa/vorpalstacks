#!/bin/bash
# STS Service Tests for VorpalStacks
# Usage: bash scripts/services/test_sts.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "STS Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

echo "Creating test IAM role for AssumeRole tests..."
aws_cmd iam create-role --role-name test-sts-role --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*","Federated":"arn:aws:iam::'"$ACCOUNT_ID"':saml-provider/test-provider"},"Action":["sts:AssumeRole","sts:AssumeRoleWithWebIdentity","sts:AssumeRoleWithSAML"]}]}' 2>/dev/null || true

run_test "1. GetCallerIdentity" \
    "aws_cmd sts get-caller-identity --query 'Arn' --output text" \
    "arn:aws:iam"

run_test "2. GetCallerIdentity (Account)" \
    "aws_cmd sts get-caller-identity --query 'Account' --output text" \
    ""

run_test "3. GetCallerIdentity (UserId)" \
    "aws_cmd sts get-caller-identity --query 'UserId' --output text" \
    ""

run_test "4. AssumeRole" \
    "aws_cmd sts assume-role --role-arn 'arn:aws:iam::$ACCOUNT_ID:role/test-sts-role' --role-session-name 'test-session' --query 'Credentials.AccessKeyId' --output text" \
    ""

run_test "5. AssumeRoleWithWebIdentity" \
    "aws_cmd sts assume-role-with-web-identity --role-arn 'arn:aws:iam::$ACCOUNT_ID:role/test-sts-role' --role-session-name 'test-session' --web-identity-token 'test-token' --query 'Credentials.AccessKeyId' --output text" \
    ""

run_test "6. AssumeRoleWithSAML" \
    "aws_cmd sts assume-role-with-saml --role-arn 'arn:aws:iam::$ACCOUNT_ID:role/test-sts-role' --principal-arn 'arn:aws:iam::$ACCOUNT_ID:saml-provider/test-provider' --saml-assertion 'test-assertion' --query 'Credentials.AccessKeyId' --output text" \
    ""

run_test "7. GetSessionToken" \
    "aws_cmd sts get-session-token --query 'Credentials.AccessKeyId' --output text" \
    ""

run_test "8. GetFederationToken" \
    "aws_cmd sts get-federation-token --name 'test-user' --query 'Credentials.AccessKeyId' --output text" \
    ""

run_test "9. DecodeAuthorizationMessage" \
    "aws_cmd sts decode-authorization-message --encoded-message 'eyJtZXNzYWdlIjoidGVzdCJ9' --query 'DecodedMessage' --output text" \
    "test"

print_summary
