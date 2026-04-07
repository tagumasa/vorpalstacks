#!/bin/bash
set -e

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1

BASE_URL="http://localhost:8080"
ACCOUNT_ID="${ACCOUNT_ID:-000000000000}"

echo "=== IAM Authorization Test Setup ==="

echo "1. Creating admin user..."
aws --endpoint-url=$BASE_URL iam create-user --user-name admin-user 2>/dev/null || true

echo "2. Creating AdminPolicy..."
ADMIN_POLICY='{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}'
aws --endpoint-url=$BASE_URL iam create-policy \
    --policy-name AdminPolicy \
    --policy-document "$ADMIN_POLICY" 2>/dev/null || true

echo "3. Attaching AdminPolicy to admin-user..."
ADMIN_POLICY_ARN="arn:aws:iam::$ACCOUNT_ID:policy/AdminPolicy"
aws --endpoint-url=$BASE_URL iam attach-user-policy \
    --user-name admin-user \
    --policy-arn "$ADMIN_POLICY_ARN" 2>/dev/null || true

echo "4. Getting or creating access key for admin-user..."
ADMIN_KEYS=$(aws --endpoint-url=$BASE_URL iam list-access-keys --user-name admin-user --output json 2>/dev/null)
ADMIN_KEY_COUNT=$(echo "$ADMIN_KEYS" | jq -r '.AccessKeyMetadata | length')
if [ "$ADMIN_KEY_COUNT" -gt 0 ]; then
    ADMIN_ACCESS_KEY=$(echo "$ADMIN_KEYS" | jq -r '.AccessKeyMetadata[0].AccessKeyId')
    echo "Using existing admin key: $ADMIN_ACCESS_KEY"
else
    ADMIN_KEY=$(aws --endpoint-url=$BASE_URL iam create-access-key --user-name admin-user --output json)
    ADMIN_ACCESS_KEY=$(echo "$ADMIN_KEY" | jq -r '.AccessKey.AccessKeyId')
    echo "Admin Access Key: $ADMIN_ACCESS_KEY"
fi

echo "5. Creating no-policy-user..."
aws --endpoint-url=$BASE_URL iam create-user --user-name no-policy-user 2>/dev/null || true

echo "6. Getting or creating access key for no-policy-user..."
NO_POLICY_KEYS=$(aws --endpoint-url=$BASE_URL iam list-access-keys --user-name no-policy-user --output json 2>/dev/null)
NO_POLICY_KEY_COUNT=$(echo "$NO_POLICY_KEYS" | jq -r '.AccessKeyMetadata | length')
if [ "$NO_POLICY_KEY_COUNT" -gt 0 ]; then
    NO_POLICY_ACCESS_KEY=$(echo "$NO_POLICY_KEYS" | jq -r '.AccessKeyMetadata[0].AccessKeyId')
    echo "Using existing no-policy key: $NO_POLICY_ACCESS_KEY"
else
    NO_POLICY_KEY=$(aws --endpoint-url=$BASE_URL iam create-access-key --user-name no-policy-user --output json)
    NO_POLICY_ACCESS_KEY=$(echo "$NO_POLICY_KEY" | jq -r '.AccessKey.AccessKeyId')
    echo "No-Policy Access Key: $NO_POLICY_ACCESS_KEY"
fi

echo ""
echo "=== Test Results ==="
echo "ADMIN_ACCESS_KEY=$ADMIN_ACCESS_KEY"
echo "NO_POLICY_ACCESS_KEY=$NO_POLICY_ACCESS_KEY"
echo ""
echo "To test authorization, set:"
echo "  AUTHORIZATION_ENABLED=true"
echo "  AUTHORIZATION_DEFAULT_ACCESS_KEY_ID=$ADMIN_ACCESS_KEY"
echo ""
echo "Or use header:"
echo "  X-Test-Access-Key-ID: $ADMIN_ACCESS_KEY  (should be ALLOWED)"
echo "  X-Test-Access-Key-ID: $NO_POLICY_ACCESS_KEY  (should be DENIED)"
