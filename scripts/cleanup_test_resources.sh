#!/bin/bash
set -e

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:8081}"
REGION="${REGION:-us-east-1}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CREDENTIALS_FILE="$SCRIPT_DIR/../tmp/vorpalstacks_test_credentials.sh"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=========================================="
echo "Cleanup Test Resources"
echo "==========================================${NC}"
echo "Endpoint URL: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

if [ -f "$CREDENTIALS_FILE" ]; then
    source "$CREDENTIALS_FILE"
fi

cleanup_iam() {
    echo "Cleaning up IAM resources..."
    
    for user in test-cli-user test-user-2 test-user-renamed; do
        echo -n "  Deleting user $user... "
        aws iam delete-user --user-name "$user" \
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
            --no-sign-request 2>/dev/null && echo -e "${GREEN}OK${NC}" || echo -e "${YELLOW}NOT FOUND${NC}"
    done
    
    for role in test-role test-assume-role; do
        echo -n "  Deleting role $role... "
        aws iam delete-role --role-name "$role" \
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
            --no-sign-request 2>/dev/null && echo -e "${GREEN}OK${NC}" || echo -e "${YELLOW}NOT FOUND${NC}"
    done
    
    for group in test-group; do
        echo -n "  Deleting group $group... "
        aws iam delete-group --group-name "$group" \
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
            --no-sign-request 2>/dev/null && echo -e "${GREEN}OK${NC}" || echo -e "${YELLOW}NOT FOUND${NC}"
    done
    
    POLICIES=$(aws iam list-policies --scope Local \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        --no-sign-request 2>/dev/null | jq -r '.Policies[].Arn' 2>/dev/null || true)
    
    for policy_arn in $POLICIES; do
        if echo "$policy_arn" | grep -q "test-policy"; then
            echo -n "  Deleting policy $policy_arn... "
            aws iam delete-policy --policy-arn "$policy_arn" \
                --endpoint-url "$ENDPOINT_URL" \
                --region "$REGION" \
                --no-sign-request 2>/dev/null && echo -e "${GREEN}OK${NC}" || echo -e "${YELLOW}FAILED${NC}"
        fi
    done
}

cleanup_kms() {
    echo "Cleaning up KMS resources..."
    
    echo -n "  Deleting alias alias/test-cli-key... "
    aws kms delete-alias --alias-name alias/test-cli-key \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" 2>/dev/null && echo -e "${GREEN}OK${NC}" || echo -e "${YELLOW}NOT FOUND${NC}"
}

cleanup_files() {
    echo "Cleaning up credential files..."
    
    if [ -f "$CREDENTIALS_FILE" ]; then
        echo -n "  Removing $CREDENTIALS_FILE... "
        rm -f "$CREDENTIALS_FILE"
        echo -e "${GREEN}OK${NC}"
    fi
}

cleanup_iam
cleanup_kms
cleanup_files

echo ""
echo -e "${GREEN}Cleanup complete!${NC}"
