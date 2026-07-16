#!/bin/bash
# Athena Service Tests for VorpalStacks
# Usage: bash scripts/services/test_athena.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Athena Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

WORKGROUP_NAME="test-workgroup-$$"
NAMED_QUERY_NAME="test-query-$$"
NAMED_QUERY_ID=""

cleanup() {
    [ -n "$NAMED_QUERY_ID" ] && aws_noauth athena delete-named-query --named-query-id "$NAMED_QUERY_ID" 2>/dev/null || true
    aws_noauth athena delete-work-group --work-group "$WORKGROUP_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. CreateWorkGroup" \
    "aws_noauth athena create-work-group --name '$WORKGROUP_NAME' --description 'Test workgroup'" \
    ""

run_test "2. GetWorkGroup" \
    "aws_noauth athena get-work-group --work-group '$WORKGROUP_NAME' --query 'WorkGroup.Name' --output text" \
    "$WORKGROUP_NAME"

run_test "3. ListWorkGroups" \
    "aws_noauth athena list-work-groups --query 'WorkGroups[?Name==\`$WORKGROUP_NAME\`].Name' --output text" \
    "$WORKGROUP_NAME"

run_test "4. UpdateWorkGroup" \
    "aws_noauth athena update-work-group --work-group '$WORKGROUP_NAME' --description 'Updated description'" \
    ""

NAMED_QUERY_ID=$(aws_noauth athena create-named-query --name "$NAMED_QUERY_NAME" --database default --query-string "SELECT 1" --work-group "$WORKGROUP_NAME" --query 'NamedQueryId' --output text 2>/dev/null || echo "")

if [ -n "$NAMED_QUERY_ID" ] && [ "$NAMED_QUERY_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 5. CreateNamedQuery (ID: ${NAMED_QUERY_ID:0:8}...)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 5. CreateNamedQuery"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "6. GetNamedQuery" \
    "aws_noauth athena get-named-query --named-query-id '$NAMED_QUERY_ID' --query 'NamedQuery.Name' --output text" \
    "$NAMED_QUERY_NAME"

run_test "7. ListNamedQueries" \
    "aws_noauth athena list-named-queries --work-group '$WORKGROUP_NAME' --query 'NamedQueryIds | length(@)' --output text" \
    "1"

run_test "8. DeleteNamedQuery" \
    "aws_noauth athena delete-named-query --named-query-id '$NAMED_QUERY_ID'" \
    ""

NAMED_QUERY_ID=""

run_test "9. DeleteWorkGroup" \
    "aws_noauth athena delete-work-group --work-group '$WORKGROUP_NAME'" \
    ""

print_summary
