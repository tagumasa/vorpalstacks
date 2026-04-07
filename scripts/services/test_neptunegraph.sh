#!/bin/bash
# NeptuneGraph Service Tests for VorpalStacks
# Usage: bash scripts/services/test_neptunegraph.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "NeptuneGraph Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

GRAPH_ID=""
GRAPH_ID_2=""
SNAPSHOT_ID=""
TAG_KEY="testkey"
TAG_VALUE="testvalue"
GRAPH_ARN=""

cleanup() {
    [ -n "$GRAPH_ID_2" ] && aws_noauth neptune-graph delete-graph --graph-identifier "$GRAPH_ID_2" --skip-snapshot 2>/dev/null || true
    [ -n "$GRAPH_ID" ] && aws_noauth neptune-graph delete-graph --graph-identifier "$GRAPH_ID" --skip-snapshot 2>/dev/null || true
    [ -n "$SNAPSHOT_ID" ] && aws_noauth neptune-graph delete-graph-snapshot --snapshot-identifier "$SNAPSHOT_ID" 2>/dev/null || true
}
trap cleanup EXIT

GRAPH_ID=$(aws_noauth neptune-graph create-graph --graph-name "test-graph-$$" --provisioned-memory 128 --query 'id' --output text 2>/dev/null || echo "")

if [ -n "$GRAPH_ID" ] && [ "$GRAPH_ID" != "None" ] && [ "$GRAPH_ID" != "" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateGraph (ID: ${GRAPH_ID:0:8}...)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateGraph"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

GRAPH_ARN=$(aws_noauth neptune-graph get-graph --graph-identifier "$GRAPH_ID" --query 'arn' --output text 2>/dev/null || echo "")

run_test "2. GetGraph" \
    "aws_noauth neptune-graph get-graph --graph-identifier '$GRAPH_ID' --query 'id' --output text" \
    "$GRAPH_ID"

run_test "3. ListGraphs" \
    "aws_noauth neptune-graph list-graphs --query 'graphs[?id==\`$GRAPH_ID\`].id' --output text" \
    "$GRAPH_ID"

run_test "4. TagResource" \
    "aws_noauth neptune-graph tag-resource --resource-arn '$GRAPH_ARN' --tags Key=$TAG_KEY,Value=$TAG_VALUE" \
    ""

run_test "5. ListTagsForResource" \
    "aws_noauth neptune-graph list-tags-for-resource --resource-arn '$GRAPH_ARN' --query 'tags.$TAG_KEY' --output text" \
    "$TAG_VALUE"

run_test "6. UntagResource" \
    "aws_noauth neptune-graph untag-resource --resource-arn '$GRAPH_ARN' --tag-keys '$TAG_KEY'" \
    ""

run_test "7. UpdateGraph" \
    "aws_noauth neptune-graph update-graph --graph-identifier '$GRAPH_ID' --provisioned-memory 256" \
    ""

SNAPSHOT_ID=$(aws_noauth neptune-graph create-graph-snapshot --graph-identifier "$GRAPH_ID" --snapshot-name "test-snapshot-$$" --query 'id' --output text 2>/dev/null || echo "")

if [ -n "$SNAPSHOT_ID" ] && [ "$SNAPSHOT_ID" != "None" ] && [ "$SNAPSHOT_ID" != "" ]; then
    echo -e "${GREEN}PASS${NC}: 8. CreateGraphSnapshot (ID: ${SNAPSHOT_ID:0:8}...)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 8. CreateGraphSnapshot"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "9. GetGraphSnapshot" \
    "aws_noauth neptune-graph get-graph-snapshot --snapshot-identifier '$SNAPSHOT_ID' --query 'id' --output text" \
    "$SNAPSHOT_ID"

run_test "10. ListGraphSnapshots" \
    "aws_noauth neptune-graph list-graph-snapshots --query 'graphSnapshots[?id==\`$SNAPSHOT_ID\`].id' --output text" \
    "$SNAPSHOT_ID"

GRAPH_ID_2=$(aws_noauth neptune-graph restore-graph-from-snapshot --snapshot-identifier "$SNAPSHOT_ID" --graph-name "restored-graph-$$" --provisioned-memory 128 --query 'id' --output text 2>/dev/null || echo "")

if [ -n "$GRAPH_ID_2" ] && [ "$GRAPH_ID_2" != "None" ] && [ "$GRAPH_ID_2" != "" ]; then
    echo -e "${GREEN}PASS${NC}: 11. RestoreGraphFromSnapshot (ID: ${GRAPH_ID_2:0:8}...)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 11. RestoreGraphFromSnapshot"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "12. StopGraph (restored)" \
    "aws_noauth neptune-graph stop-graph --graph-identifier '$GRAPH_ID_2'" \
    ""

run_test "13. StartGraph (restored)" \
    "aws_noauth neptune-graph start-graph --graph-identifier '$GRAPH_ID_2'" \
    ""

run_test "14. DeleteGraphSnapshot" \
    "aws_noauth neptune-graph delete-graph-snapshot --snapshot-identifier '$SNAPSHOT_ID'" \
    ""
SNAPSHOT_ID=""

run_test "15. DeleteGraph (restored)" \
    "aws_noauth neptune-graph delete-graph --graph-identifier '$GRAPH_ID_2' --skip-snapshot" \
    ""
GRAPH_ID_2=""

run_test "16. DeleteGraph" \
    "aws_noauth neptune-graph delete-graph --graph-identifier '$GRAPH_ID' --skip-snapshot" \
    ""
GRAPH_ID=""

print_summary
