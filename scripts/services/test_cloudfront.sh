#!/bin/bash
# CloudFront Service Tests for VorpalStacks
# Usage: bash scripts/services/test_cloudfront.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "CloudFront Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

DIST_ID=""
CALLER_REF="caller-ref-$$"
TMP_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$TMP_DIR" 2>/dev/null || true
    if [ -n "$DIST_ID" ]; then
        ETAG=$(aws_noauth cloudfront get-distribution --id "$DIST_ID" --query 'ETag' --output text 2>/dev/null || echo "")
        if [ -n "$ETAG" ] && [ "$ETAG" != "None" ]; then
            aws_noauth cloudfront delete-distribution --id "$DIST_ID" --if-match "$ETAG" 2>/dev/null || true
        fi
    fi
}
trap cleanup EXIT

cat > "$TMP_DIR/create.json" <<EOF
{
  "DistributionConfig": {
    "CallerReference": "$CALLER_REF",
    "Comment": "Test distribution",
    "Enabled": true,
    "Origins": {
      "Quantity": 1,
      "Items": [
        {
          "Id": "origin1",
          "DomainName": "example.com",
          "CustomOriginConfig": {
            "HTTPPort": 80,
            "HTTPSPort": 443,
            "OriginProtocolPolicy": "http-only"
          }
        }
      ]
    },
    "DefaultCacheBehavior": {
      "TargetOriginId": "origin1",
      "ViewerProtocolPolicy": "allow-all",
      "ForwardedValues": {
        "QueryString": false,
        "Cookies": {
          "Forward": "none"
        }
      },
      "MinTTL": 0
    }
  }
}
EOF

DIST_ID=$(aws_noauth cloudfront create-distribution --cli-input-json "file://$TMP_DIR/create.json" --query 'Distribution.Id' --output text 2>/dev/null || echo "")

if [ -n "$DIST_ID" ] && [ "$DIST_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 1. CreateDistribution (ID: $DIST_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 1. CreateDistribution"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "2. GetDistribution" \
    "aws_noauth cloudfront get-distribution --id '$DIST_ID' --query 'Distribution.DistributionConfig.Comment' --output text" \
    "Test distribution"

run_test "3. ListDistributions" \
    "aws_noauth cloudfront list-distributions --query 'DistributionList.Items[?Comment==\`Test distribution\`].Comment' --output text" \
    "Test distribution"

cat > "$TMP_DIR/update.json" <<EOF
{
  "DistributionConfig": {
    "CallerReference": "$CALLER_REF",
    "Comment": "Updated comment",
    "Enabled": true,
    "Origins": {
      "Quantity": 1,
      "Items": [
        {
          "Id": "origin1",
          "DomainName": "example.com",
          "CustomOriginConfig": {
            "HTTPPort": 80,
            "HTTPSPort": 443,
            "OriginProtocolPolicy": "http-only"
          }
        }
      ]
    },
    "DefaultCacheBehavior": {
      "TargetOriginId": "origin1",
      "ViewerProtocolPolicy": "allow-all",
      "ForwardedValues": {
        "QueryString": false,
        "Cookies": {
          "Forward": "none"
        }
      },
      "MinTTL": 0
    }
  },
  "Id": "$DIST_ID",
  "IfMatch": "*"
}
EOF

run_test "4. UpdateDistribution" \
    "aws_noauth cloudfront update-distribution --cli-input-json 'file://$TMP_DIR/update.json' --query 'Distribution.DistributionConfig.Comment' --output text" \
    "Updated comment"

run_test "5. GetDistributionConfig" \
    "aws_noauth cloudfront get-distribution-config --id '$DIST_ID' --query 'DistributionConfig.Comment' --output text" \
    "Updated comment"

run_test "6. TagResource" \
    "aws_noauth cloudfront tag-resource --resource 'arn:aws:cloudfront::$ACCOUNT_ID:distribution/$DIST_ID' --tags 'Items=[{Key=Environment,Value=Test}]'" \
    ""

run_test "7. ListTagsForResource" \
    "aws_noauth cloudfront list-tags-for-resource --resource 'arn:aws:cloudfront::$ACCOUNT_ID:distribution/$DIST_ID' --query 'Tags.Items[?Key==\`Environment\`].Value' --output text" \
    "Test"

cat > "$TMP_DIR/disable.json" <<EOF
{
  "DistributionConfig": {
    "CallerReference": "$CALLER_REF",
    "Comment": "Updated comment",
    "Enabled": false,
    "Origins": {
      "Quantity": 1,
      "Items": [
        {
          "Id": "origin1",
          "DomainName": "example.com",
          "CustomOriginConfig": {
            "HTTPPort": 80,
            "HTTPSPort": 443,
            "OriginProtocolPolicy": "http-only"
          }
        }
      ]
    },
    "DefaultCacheBehavior": {
      "TargetOriginId": "origin1",
      "ViewerProtocolPolicy": "allow-all",
      "ForwardedValues": {
        "QueryString": false,
        "Cookies": {
          "Forward": "none"
        }
      },
      "MinTTL": 0
    }
  },
  "Id": "$DIST_ID",
  "IfMatch": "*"
}
EOF

run_test "8. DisableDistribution" \
    "aws_noauth cloudfront update-distribution --cli-input-json 'file://$TMP_DIR/disable.json' --query 'Distribution.DistributionConfig.Enabled' --output text" \
    "False"

ETAG=$(aws_noauth cloudfront get-distribution --id "$DIST_ID" --query 'ETag' --output text 2>/dev/null || echo "")

run_test "9. DeleteDistribution" \
    "aws_noauth cloudfront delete-distribution --id '$DIST_ID' --if-match '$ETAG'" \
    ""

DIST_ID=""

# Phase 2: Origin Access Control Tests
echo ""
echo "=== Origin Access Control Tests ==="

OAC_ID=""

cleanup_oac() {
    if [ -n "$OAC_ID" ]; then
        aws_noauth cloudfront delete-origin-access-control --id "$OAC_ID" 2>/dev/null || true
    fi
}

cat > "$TMP_DIR/oac_create.json" <<EOF
{
  "OriginAccessControlConfig": {
    "Name": "test-oac-$$",
    "OriginAccessControlOriginType": "s3",
    "SigningBehavior": "always",
    "SigningProtocol": "sigv4"
  }
}
EOF

OAC_ID=$(aws_noauth cloudfront create-origin-access-control --cli-input-json "file://$TMP_DIR/oac_create.json" --query 'OriginAccessControl.Id' --output text 2>/dev/null || echo "")

if [ -n "$OAC_ID" ] && [ "$OAC_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 10. CreateOriginAccessControl (ID: $OAC_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 10. CreateOriginAccessControl"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "11. GetOriginAccessControl" \
    "aws_noauth cloudfront get-origin-access-control --id '$OAC_ID' --query 'OriginAccessControl.OriginAccessControlConfig.Name' --output text" \
    "test-oac-$$"

run_test "12. ListOriginAccessControls" \
    "aws_noauth cloudfront list-origin-access-controls --query 'OriginAccessControlList.Items[?Name==\`test-oac-$$\`].Name' --output text" \
    "test-oac-$$"

cat > "$TMP_DIR/oac_update.json" <<EOF
{
  "OriginAccessControlConfig": {
    "Name": "test-oac-$$",
    "Description": "Updated description",
    "OriginAccessControlOriginType": "s3",
    "SigningBehavior": "no-override",
    "SigningProtocol": "sigv4"
  },
  "Id": "$OAC_ID",
  "IfMatch": "*"
}
EOF

run_test "13. UpdateOriginAccessControl" \
    "aws_noauth cloudfront update-origin-access-control --cli-input-json 'file://$TMP_DIR/oac_update.json' --query 'OriginAccessControl.OriginAccessControlConfig.Description' --output text" \
    "Updated description"

run_test "14. DeleteOriginAccessControl" \
    "aws_noauth cloudfront delete-origin-access-control --id '$OAC_ID'" \
    ""

OAC_ID=""

# Phase 4: Response Headers Policy Tests
echo ""
echo "=== Response Headers Policy Tests ==="

RHP_ID=""

cat > "$TMP_DIR/rhp_create.json" <<EOF
{
  "ResponseHeadersPolicyConfig": {
    "Name": "test-rhp-$$",
    "Comment": "Test response headers policy",
    "SecurityHeadersConfig": {
      "ContentTypeOptions": {"Override": true},
      "FrameOptions": {"FrameOption": "SAMEORIGIN", "Override": true},
      "XSSProtection": {"Protection": true, "Override": true},
      "ReferrerPolicy": {"ReferrerPolicy": "strict-origin-when-cross-origin", "Override": true},
      "StrictTransportSecurity": {"AccessControlMaxAgeSec": 31536000, "Override": true}
    }
  }
}
EOF

RHP_ID=$(aws_noauth cloudfront create-response-headers-policy --cli-input-json "file://$TMP_DIR/rhp_create.json" --query 'ResponseHeadersPolicy.Id' --output text 2>/dev/null || echo "")

if [ -n "$RHP_ID" ] && [ "$RHP_ID" != "None" ]; then
    echo -e "${GREEN}PASS${NC}: 15. CreateResponseHeadersPolicy (ID: $RHP_ID)"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}: 15. CreateResponseHeadersPolicy"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
TEST_COUNT=$((TEST_COUNT + 1))

run_test "16. GetResponseHeadersPolicy" \
    "aws_noauth cloudfront get-response-headers-policy --id '$RHP_ID' --query 'ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Comment' --output text" \
    "Test response headers policy"

run_test "17. GetResponseHeadersPolicyConfig" \
    "aws_noauth cloudfront get-response-headers-policy-config --id '$RHP_ID' --query 'ResponseHeadersPolicyConfig.SecurityHeadersConfig.FrameOptions.FrameOption' --output text" \
    "SAMEORIGIN"

run_test "18. ListResponseHeadersPolicies" \
    "aws_noauth cloudfront list-response-headers-policies --query 'ResponseHeadersPolicyList.Items[?ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Name==\`test-rhp-$$\`].ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Name' --output text" \
    "test-rhp-$$"

cat > "$TMP_DIR/rhp_update.json" <<EOF
{
  "ResponseHeadersPolicyConfig": {
    "Name": "test-rhp-$$",
    "Comment": "Updated policy comment",
    "SecurityHeadersConfig": {
      "ContentTypeOptions": {"Override": true},
      "FrameOptions": {"FrameOption": "DENY", "Override": true}
    }
  },
  "Id": "$RHP_ID",
  "IfMatch": "*"
}
EOF

run_test "19. UpdateResponseHeadersPolicy" \
    "aws_noauth cloudfront update-response-headers-policy --cli-input-json 'file://$TMP_DIR/rhp_update.json' --query 'ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Comment' --output text" \
    "Updated policy comment"

run_test "20. DeleteResponseHeadersPolicy" \
    "aws_noauth cloudfront delete-response-headers-policy --id '$RHP_ID'" \
    ""

RHP_ID=""

print_summary
