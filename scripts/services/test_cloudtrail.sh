#!/bin/bash
# CloudTrail Service Tests for VorpalStacks
# Usage: bash scripts/services/test_cloudtrail.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "CloudTrail Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

TRAIL_NAME="test-trail-$$"
BUCKET_NAME="cloudtrail-logs-$$"

cleanup() {
    aws_noauth cloudtrail delete-trail --name "$TRAIL_NAME" 2>/dev/null || true
    aws_noauth s3 rm "s3://$BUCKET_NAME" --recursive 2>/dev/null || true
    aws_noauth s3 rb "s3://$BUCKET_NAME" 2>/dev/null || true
}
trap cleanup EXIT

aws_noauth s3 mb "s3://$BUCKET_NAME" 2>/dev/null || true

# Create trail
aws_noauth cloudtrail create-trail --name "$TRAIL_NAME" --s3-bucket-name "$BUCKET_NAME" > /dev/null 2>&1

run_test "1. CreateTrail" \
    "aws_noauth cloudtrail describe-trails --query \"trailList[?@.Name=='$TRAIL_NAME'] | [0].Name\" --output text" \
    "$TRAIL_NAME"

run_test "2. DescribeTrails" \
    "aws_noauth cloudtrail describe-trails --query \"trailList[?@.Name=='$TRAIL_NAME'] | [0].Name\" --output text" \
    "$TRAIL_NAME"

run_test "3. GetTrail" \
    "aws_noauth cloudtrail get-trail --name '$TRAIL_NAME' --query 'Trail.Name' --output text" \
    "$TRAIL_NAME"

run_test "4. GetTrailStatus" \
    "aws_noauth cloudtrail get-trail-status --name '$TRAIL_NAME' --query 'IsLogging' --output text" \
    "False"

run_test "5. StartLogging" \
    "aws_noauth cloudtrail start-logging --name '$TRAIL_NAME'" \
    ""

run_test "6. GetTrailStatus (logging)" \
    "aws_noauth cloudtrail get-trail-status --name '$TRAIL_NAME' --query 'IsLogging' --output text" \
    "True"

run_test "7. StopLogging" \
    "aws_noauth cloudtrail stop-logging --name '$TRAIL_NAME'" \
    ""

run_test "8. GetTrailStatus (stopped)" \
    "aws_noauth cloudtrail get-trail-status --name '$TRAIL_NAME' --query 'IsLogging' --output text" \
    "False"

run_test "9. ListTrails" \
    "aws_noauth cloudtrail list-trails --query \"Trails[?@.Name=='$TRAIL_NAME'] | [0].Name\" --output text" \
    "$TRAIL_NAME"

# Get the actual TrailARN from describe-trails (account ID comes from auth header)
TRAIL_ARN=$(aws_noauth cloudtrail describe-trails --query "trailList[?Name=='$TRAIL_NAME'].TrailARN | [0]" --output text)

run_test "10. AddTags" \
    "aws_noauth cloudtrail add-tags --resource-id '$TRAIL_ARN' --tags-list Key=Environment,Value=Test" \
    ""

run_test "11. ListTags" \
    "aws_noauth cloudtrail list-tags --resource-id-list '$TRAIL_ARN' --query 'ResourceTagList[0].TagsList[0].Value' --output text" \
    "Test"

run_test "12. RemoveTags" \
    "aws_noauth cloudtrail remove-tags --resource-id '$TRAIL_ARN' --tags-list Key=Environment,Value=Test" \
    ""

run_test "13. UpdateTrail" \
    "aws_noauth cloudtrail update-trail --name '$TRAIL_NAME' --include-global-service-events --query 'IncludeGlobalServiceEvents' --output text" \
    "True"

run_test "14. PutEventSelectors" \
    "aws_noauth cloudtrail put-event-selectors --trail-name '$TRAIL_NAME' --event-selectors '[{\"ReadWriteType\":\"All\",\"IncludeManagementEvents\":true}]'" \
    ""

run_test "15. GetEventSelectors" \
    "aws_noauth cloudtrail get-event-selectors --trail-name '$TRAIL_NAME' --query 'EventSelectors[0].IncludeManagementEvents' --output text" \
    "True"

run_test "16. LookupEvents" \
    "aws_noauth cloudtrail lookup-events --max-results 1 --query 'Events' --output text" \
    ""

run_test "17. DeleteTrail" \
    "aws_noauth cloudtrail delete-trail --name '$TRAIL_NAME'" \
    ""

run_test "18. Verify deletion" \
    "aws_noauth cloudtrail get-trail --name '$TRAIL_NAME' 2>&1 | head -5" \
    "TrailNotFoundException\|InternalFailure"

print_summary
