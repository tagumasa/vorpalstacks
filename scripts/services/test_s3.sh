#!/bin/bash
# S3 Service Tests for VorpalStacks
# Usage: bash scripts/services/test_s3.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "S3 Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

BUCKET_NAME="test-bucket-$$"
TEST_FILE="$PROJECT_ROOT/tmp/s3-test-$$"
DOWNLOAD_FILE="$PROJECT_ROOT/tmp/s3-download-$$"

cleanup() {
    aws_noauth s3 rm "s3://$BUCKET_NAME" --recursive 2>/dev/null || true
    aws_noauth s3 rb "s3://$BUCKET_NAME" 2>/dev/null || true
    rm -f "$TEST_FILE" "$DOWNLOAD_FILE"
}
trap cleanup EXIT

echo "$S3_TEST_CONTENT" > "$TEST_FILE"

run_test "1. CreateBucket (mb)" \
    "aws_noauth s3 mb 's3://$BUCKET_NAME'" \
    "make_bucket"

run_test "2. ListBuckets (ls)" \
    "aws_noauth s3 ls | grep '$BUCKET_NAME'" \
    "$BUCKET_NAME"

run_test "3. HeadBucket" \
    "aws_noauth s3api head-bucket --bucket '$BUCKET_NAME'" \
    ""

run_test "4. PutObject (cp)" \
    "aws_noauth s3 cp '$TEST_FILE' 's3://$BUCKET_NAME/test-object.txt'" \
    "upload"

run_test "5. ListObjects (ls)" \
    "aws_noauth s3 ls 's3://$BUCKET_NAME/'" \
    "test-object.txt"

run_test "6. ListObjectsV2" \
    "aws_noauth s3api list-objects-v2 --bucket '$BUCKET_NAME' --query \"Contents[?@.Key=='test-object.txt'] | [0].Key\" --output text" \
    "test-object.txt"

run_test "7. HeadObject" \
    "aws_noauth s3api head-object --bucket '$BUCKET_NAME' --key 'test-object.txt' --query 'ContentLength' --output text" \
    ""

run_test "8. GetObject (cp)" \
    "aws_noauth s3 cp 's3://$BUCKET_NAME/test-object.txt' '$DOWNLOAD_FILE'" \
    "download"

run_test "9. Verify content" \
    "cat '$DOWNLOAD_FILE'" \
    "VorpalStacks"

run_test "10. CopyObject" \
    "aws_noauth s3 cp 's3://$BUCKET_NAME/test-object.txt' 's3://$BUCKET_NAME/copied.txt'" \
    "copy"

run_test "11. DeleteObject (rm)" \
    "aws_noauth s3 rm 's3://$BUCKET_NAME/test-object.txt'" \
    "delete"

run_test "12. PutBucketTagging" \
    "aws_noauth s3api put-bucket-tagging --bucket '$BUCKET_NAME' --tagging 'TagSet=[{Key=Environment,Value=Test}]'" \
    ""

run_test "13. GetBucketTagging" \
    "aws_noauth s3api get-bucket-tagging --bucket '$BUCKET_NAME' --query \"TagSet[?@.Key=='Environment'] | [0].Value\" --output text" \
    "Test"

run_test "14. PutBucketVersioning" \
    "aws_noauth s3api put-bucket-versioning --bucket '$BUCKET_NAME' --versioning-configuration Status=Enabled" \
    ""

run_test "15. GetBucketVersioning" \
    "aws_noauth s3api get-bucket-versioning --bucket '$BUCKET_NAME' --query 'Status' --output text" \
    "Enabled"

run_test "16. PutObject (versioned)" \
    "aws_noauth s3 cp '$TEST_FILE' 's3://$BUCKET_NAME/versioned.txt'" \
    "upload"

run_test "17. ListObjectVersions" \
    "aws_noauth s3api list-object-versions --bucket '$BUCKET_NAME' --query \"Versions[?@.Key=='versioned.txt'] | [0].Key\" --output text" \
    "versioned.txt"

run_test "18. CreateMultipartUpload" \
    "aws_noauth s3api create-multipart-upload --bucket '$BUCKET_NAME' --key 'multipart.dat' --query 'UploadId' --output text" \
    ""

run_test "19. ListMultipartUploads" \
    "aws_noauth s3api list-multipart-uploads --bucket '$BUCKET_NAME' --query \"Uploads[?@.Key=='multipart.dat'] | [0].Key\" --output text" \
    "multipart.dat"

run_test "20. AbortMultipartUpload" \
    "UPLOAD_ID=\$(aws_noauth s3api list-multipart-uploads --bucket '$BUCKET_NAME' --query 'Uploads[0].UploadId' --output text 2>/dev/null) && aws_noauth s3api abort-multipart-upload --bucket '$BUCKET_NAME' --key 'multipart.dat' --upload-id \"\$UPLOAD_ID\"" \
    ""

run_test "21. DeleteObject (remaining)" \
    "aws_noauth s3 rm 's3://$BUCKET_NAME/copied.txt'" \
    "delete"

run_test "22. DeleteBucket (rb)" \
    "aws_noauth s3 rb 's3://$BUCKET_NAME'" \
    "remove_bucket"

print_summary
