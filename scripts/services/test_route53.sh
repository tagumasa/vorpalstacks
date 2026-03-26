#!/bin/bash
# Route53 Service Tests for VorpalStacks
# Usage: bash scripts/services/test_route53.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Route53 Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

ZONE_NAME="test-zone-$$.example.com."
ZONE_ID=""

cleanup() {
    if [ -n "$ZONE_ID" ]; then
        aws_noauth route53 delete-hosted-zone --id "$ZONE_ID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

ZONE_ID=$(aws_noauth route53 create-hosted-zone --name "$ZONE_NAME" --caller-reference "ref-$$" --query 'HostedZone.Id' --output text 2>/dev/null || echo "")

if [ -n "$ZONE_ID" ] && [ "$ZONE_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateHostedZone (ID: $ZONE_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateHostedZone"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "2. GetHostedZone" \
    "aws_noauth route53 get-hosted-zone --id '$ZONE_ID' --query 'HostedZone.Name' --output text" \
    "$ZONE_NAME"

run_test "3. ListHostedZones" \
    "aws_noauth route53 list-hosted-zones --query 'HostedZones[?Name==\`$ZONE_NAME\`].Name' --output text" \
    "$ZONE_NAME"

run_test "4. ChangeResourceRecordSets (Create)" \
    "aws_noauth route53 change-resource-record-sets --hosted-zone-id '$ZONE_ID' --change-batch '{\"Changes\":[{\"Action\":\"CREATE\",\"ResourceRecordSet\":{\"Name\":\"test.$ZONE_NAME\",\"Type\":\"A\",\"TTL\":300,\"ResourceRecords\":[{\"Value\":\"1.2.3.4\"}]}}]}' --query 'ChangeInfo.Status' --output text" \
    "INSYNC"

run_test "5. ListResourceRecordSets" \
    "aws_noauth route53 list-resource-record-sets --hosted-zone-id '$ZONE_ID' --query 'ResourceRecordSets[?Type==\`A\`].Type' --output text" \
    "A"

run_test "6. ChangeResourceRecordSets (Delete)" \
    "aws_noauth route53 change-resource-record-sets --hosted-zone-id '$ZONE_ID' --change-batch '{\"Changes\":[{\"Action\":\"DELETE\",\"ResourceRecordSet\":{\"Name\":\"test.$ZONE_NAME\",\"Type\":\"A\",\"TTL\":300,\"ResourceRecords\":[{\"Value\":\"1.2.3.4\"}]}}]}' --query 'ChangeInfo.Status' --output text" \
    "INSYNC"

run_test "7. UpdateHostedZoneComment" \
    "aws_noauth route53 update-hosted-zone-comment --id '$ZONE_ID' --comment 'Updated comment' --query 'HostedZone.Name' --output text" \
    "$ZONE_NAME"

run_test "8. DeleteHostedZone" \
    "aws_noauth route53 delete-hosted-zone --id '$ZONE_ID' --query 'ChangeInfo.Status' --output text" \
    ""

ZONE_ID=""

print_summary
