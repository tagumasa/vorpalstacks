#!/bin/bash
# Neptune Service Tests for VorpalStacks
# Usage: bash scripts/services/test_neptune.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Neptune Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

CLUSTER_NAME="test-cluster-$$"
CLUSTER_ARN="arn:aws:rds:$REGION:$ACCOUNT_ID:cluster:$CLUSTER_NAME"
TAG_KEY="testkey"
TAG_VALUE="testvalue"

cleanup() {
    aws_noauth neptune delete-db-cluster --db-cluster-identifier "$CLUSTER_NAME" --skip-final-snapshot 2>/dev/null || true
}
trap cleanup EXIT

CLUSTER_ID=$(aws_noauth neptune create-db-cluster \
    --db-cluster-identifier "$CLUSTER_NAME" \
    --engine neptune \
    --master-username admin \
    --master-user-password "TestPass123!" \
    --query 'DBCluster.DBClusterIdentifier' --output text 2>/dev/null || echo "")

if [ -n "$CLUSTER_ID" ] && [ "$CLUSTER_ID" != "None" ] && [ "$CLUSTER_ID" != "" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateDBCluster (ID: $CLUSTER_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateDBCluster"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "2. DescribeDBClusters" \
    "aws_noauth neptune describe-db-clusters --db-cluster-identifier '$CLUSTER_NAME' --query 'DBClusters[0].DBClusterIdentifier' --output text" \
    "$CLUSTER_NAME"

run_test "3. DescribeDBEngineVersions" \
    "aws_noauth neptune describe-db-engine-versions --engine neptune --query 'DBEngineVersions | length(@)' --output text" \
    ""

run_test "4. ListTagsForResource" \
    "aws_noauth neptune list-tags-for-resource --resource-name '$CLUSTER_ARN' --output text" \
    ""

run_test "5. AddTagsToResource" \
    "aws_noauth neptune add-tags-to-resource --resource-name '$CLUSTER_ARN' --tags Key=$TAG_KEY,Value=$TAG_VALUE" \
    ""

run_test "6. ListTagsForResource (with tag)" \
    "aws_noauth neptune list-tags-for-resource --resource-name '$CLUSTER_ARN' --query 'TagList[?Key==\`'$TAG_KEY'\`].Value' --output text" \
    "$TAG_VALUE"

run_test "7. RemoveTagsFromResource" \
    "aws_noauth neptune remove-tags-from-resource --resource-name '$CLUSTER_ARN' --tag-keys '$TAG_KEY'" \
    ""

run_test "8. DeleteDBCluster" \
    "aws_noauth neptune delete-db-cluster --db-cluster-identifier '$CLUSTER_NAME' --skip-final-snapshot" \
    ""

print_summary
