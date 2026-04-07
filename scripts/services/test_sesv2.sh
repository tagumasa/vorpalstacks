#!/bin/bash
# SESv2 (Simple Email Service) Tests for VorpalStacks
# Usage: bash scripts/services/test_sesv2.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "SESv2 Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

IDENTITY="test-$$.example.com"
CONFIG_SET="test-config-set-$$"

cleanup() {
    aws_noauth sesv2 delete-email-identity --email-identity "$IDENTITY" 2>/dev/null || true
    aws_noauth sesv2 delete-configuration-set --configuration-set-name "$CONFIG_SET" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. GetAccount" \
    "aws_noauth sesv2 get-account --query 'DedicatedIpAutoWarmupMode' --output text" \
    ""

run_test "2. CreateConfigurationSet" \
    "aws_noauth sesv2 create-configuration-set --configuration-set-name '$CONFIG_SET'" \
    ""

run_test "3. GetConfigurationSet" \
    "aws_noauth sesv2 get-configuration-set --configuration-set-name '$CONFIG_SET' --query 'ConfigurationSetName' --output text" \
    "$CONFIG_SET"

run_test "4. ListConfigurationSets" \
    "aws_noauth sesv2 list-configuration-sets --query 'ConfigurationSets[?@==\`'$CONFIG_SET'\`]' --output text" \
    "$CONFIG_SET"

run_test "5. CreateEmailIdentity" \
    "aws_noauth sesv2 create-email-identity --email-identity '$IDENTITY' --query 'IdentityType' --output text" \
    "DOMAIN"

run_test "6. GetEmailIdentity" \
    "aws_noauth sesv2 get-email-identity --email-identity '$IDENTITY' --query 'IdentityType' --output text" \
    "DOMAIN"

run_test "7. ListEmailIdentities" \
    "aws_noauth sesv2 list-email-identities --query 'EmailIdentities[?IdentityName==\`$IDENTITY\`].IdentityName' --output text" \
    "$IDENTITY"

run_test "8. TagResource" \
    "aws_noauth sesv2 tag-resource --resource-arn 'arn:aws:ses:$REGION:$ACCOUNT_ID:identity/$IDENTITY' --tags Key=Environment,Value=Test" \
    ""

run_test "9. PutEmailIdentityDkimSigningAttributes" \
    "aws_noauth sesv2 put-email-identity-dkim-signing-attributes --email-identity '$IDENTITY' --signing-attributes-origin='AWS_SES'" \
    ""

run_test "10. SendEmail" \
    "aws_noauth sesv2 send-email --from-email-address 'sender@$IDENTITY' --destination 'ToAddresses=test@example.com' --content 'Simple={Subject={Data=Test},Body={Text={Data=Test message}}}' --query 'MessageId' --output text" \
    ""

run_test "11. DeleteEmailIdentity" \
    "aws_noauth sesv2 delete-email-identity --email-identity '$IDENTITY'" \
    ""

run_test "12. DeleteConfigurationSet" \
    "aws_noauth sesv2 delete-configuration-set --configuration-set-name '$CONFIG_SET'" \
    ""

print_summary
