#!/bin/bash
# ACM (Certificate Manager) Service Tests for VorpalStacks
# Usage: bash scripts/services/test_acm.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "ACM Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

DOMAIN_NAME="test-example-$$.com"
CERT_ARN=""

cleanup() {
    if [ -n "$CERT_ARN" ]; then
        aws_noauth acm delete-certificate --certificate-arn "$CERT_ARN" 2>/dev/null || true
    fi
}
trap cleanup EXIT

CERT_ARN=$(aws_noauth acm request-certificate --domain-name "$DOMAIN_NAME" --validation-method DNS --query 'CertificateArn' --output text 2>/dev/null || echo "")

if [ -n "$CERT_ARN" ] && [ "$CERT_ARN" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 1. RequestCertificate"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. RequestCertificate"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "2. DescribeCertificate" \
    "aws_noauth acm describe-certificate --certificate-arn '$CERT_ARN' --query 'Certificate.DomainName' --output text" \
    "$DOMAIN_NAME"

run_test "3. ListCertificates" \
    "aws_noauth acm list-certificates --query 'CertificateSummaryList[?DomainName==\`$DOMAIN_NAME\`].DomainName' --output text" \
    "$DOMAIN_NAME"

run_test "4. AddTagsToCertificate" \
    "aws_noauth acm add-tags-to-certificate --certificate-arn '$CERT_ARN' --tags Key=Environment,Value=Test" \
    ""

run_test "5. ListTagsForCertificate" \
    "aws_noauth acm list-tags-for-certificate --certificate-arn '$CERT_ARN' --query 'Tags[?Key==\`Environment\`].Value' --output text" \
    "Test"

run_test "6. RemoveTagsFromCertificate" \
    "aws_noauth acm remove-tags-from-certificate --certificate-arn '$CERT_ARN' --tags Key=Environment" \
    ""

run_test "7. ResendValidationEmail" \
    "aws_noauth acm resend-validation-email --certificate-arn '$CERT_ARN' --domain '$DOMAIN_NAME' --validation-domain 'example.com'" \
    ""

run_test "8. UpdateCertificateOptions" \
    "aws_noauth acm update-certificate-options --certificate-arn '$CERT_ARN' --options CertificateTransparencyLoggingPreference=DISABLED" \
    ""

run_test "9. DeleteCertificate" \
    "aws_noauth acm delete-certificate --certificate-arn '$CERT_ARN'" \
    ""

CERT_ARN=""

print_summary
