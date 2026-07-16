#!/bin/bash
set -e

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:8081}"
REGION="${REGION:-us-east-1}"
TEST_USER="test-cli-user"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CREDENTIALS_FILE="$SCRIPT_DIR/../tmp/vorpalstacks_test_credentials.sh"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=========================================="
echo "Setup Test Credentials"
echo "==========================================${NC}"
echo "Endpoint URL: $ENDPOINT_URL"
echo "Region: $REGION"
echo "Test User: $TEST_USER"
echo ""

check_server() {
    echo -n "Checking server availability... "
    if curl -s -o /dev/null -w "%{http_code}" "$ENDPOINT_URL/" 2>/dev/null | grep -q "200\|403\|404\|500"; then
        echo -e "${GREEN}OK${NC}"
        return 0
    else
        echo -e "${RED}FAILED${NC}"
        echo "Server is not available at $ENDPOINT_URL"
        exit 1
    fi
}

cleanup_existing_user() {
    echo -n "Cleaning up existing test user... "
    aws iam delete-user --user-name "$TEST_USER" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        --no-sign-request 2>/dev/null || true
    echo -e "${GREEN}Done${NC}"
}

create_test_user() {
    echo -n "Creating test user '$TEST_USER'... "
    RESULT=$(aws iam create-user \
        --user-name "$TEST_USER" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        --no-sign-request 2>&1)
    
    if echo "$RESULT" | grep -q "UserName"; then
        USER_ARN=$(echo "$RESULT" | jq -r '.User.Arn' 2>/dev/null || echo "N/A")
        echo -e "${GREEN}OK${NC} (ARN: $USER_ARN)"
    else
        echo -e "${RED}FAILED${NC}"
        echo "$RESULT"
        exit 1
    fi
}

create_access_key() {
    echo -n "Creating access key for '$TEST_USER'... "
    RESULT=$(aws iam create-access-key \
        --user-name "$TEST_USER" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        --no-sign-request 2>&1)
    
    if echo "$RESULT" | grep -q "AccessKeyId"; then
        ACCESS_KEY_ID=$(echo "$RESULT" | jq -r '.AccessKey.AccessKeyId')
        SECRET_ACCESS_KEY=$(echo "$RESULT" | jq -r '.AccessKey.SecretAccessKey')
        echo -e "${GREEN}OK${NC}"
        echo "  Access Key ID: $ACCESS_KEY_ID"
        echo "  Secret Access Key: ${SECRET_ACCESS_KEY:0:10}..."
        
        cat > "$CREDENTIALS_FILE" << EOF
#!/bin/bash
export AWS_ACCESS_KEY_ID="$ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="$SECRET_ACCESS_KEY"
export AWS_DEFAULT_REGION="$REGION"
export ENDPOINT_URL="$ENDPOINT_URL"
export TEST_USER="$TEST_USER"
EOF
        chmod +x "$CREDENTIALS_FILE"
        echo ""
        echo -e "${GREEN}Credentials saved to: $CREDENTIALS_FILE${NC}"
        
        export AWS_ACCESS_KEY_ID
        export AWS_SECRET_ACCESS_KEY
        export AWS_DEFAULT_REGION="$REGION"
    else
        echo -e "${RED}FAILED${NC}"
        echo "$RESULT"
        exit 1
    fi
}

verify_credentials() {
    echo ""
    echo -n "Verifying credentials with GetCallerIdentity... "
    RESULT=$(aws sts get-caller-identity \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" 2>&1)
    
    if echo "$RESULT" | grep -q "Account"; then
        ACCOUNT=$(echo "$RESULT" | jq -r '.Account' 2>/dev/null || echo "N/A")
        echo -e "${GREEN}OK${NC} (Account: $ACCOUNT)"
    else
        echo -e "${YELLOW}WARNING${NC}"
        echo "  STS GetCallerIdentity may not be fully functional"
        echo "  Response: $RESULT"
    fi
}

check_server
cleanup_existing_user
create_test_user
create_access_key
verify_credentials

echo ""
echo -e "${GREEN}=========================================="
echo "Setup Complete!"
echo "==========================================${NC}"
echo ""
echo "To use these credentials in other scripts:"
echo "  source $CREDENTIALS_FILE"
echo ""
echo "Or run tests directly:"
echo "  ./scripts/test_iam_cli.sh"
echo "  ./scripts/test_sts_cli.sh"
echo "  ./scripts/test_kms_cli.sh"
