#!/bin/bash
# WAF (Web Application Firewall) Service Tests for VorpalStacks
# Usage: bash scripts/services/test_waf.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "WAF Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

WEBACL_NAME="test-webacl-$$"
WEBACL_ID=""
WEBACL_ARN=""
WEBACL_LOCK_TOKEN=""

IPSET_NAME="test-ipset-$$"
IPSET_ID=""
IPSET_LOCK_TOKEN=""

RULEGROUP_NAME="test-rulegroup-$$"
RULEGROUP_ID=""
RULEGROUP_LOCK_TOKEN=""

REGEXPATTERNSET_NAME="test-regexpatternset-$$"
REGEXPATTERNSET_ID=""
REGEXPATTERNSET_LOCK_TOKEN=""

RESOURCE_ARN="arn:aws:apigateway:$REGION::/restapis/abc123/stages/prod"

cleanup() {
    if [ -n "$WEBACL_ID" ]; then
        aws_noauth wafv2 delete-web-acl --name "$WEBACL_NAME" --id "$WEBACL_ID" --lock-token "$WEBACL_LOCK_TOKEN" --scope REGIONAL 2>/dev/null || true
    fi
    if [ -n "$IPSET_ID" ]; then
        aws_noauth wafv2 delete-ip-set --name "$IPSET_NAME" --id "$IPSET_ID" --lock-token "$IPSET_LOCK_TOKEN" --scope REGIONAL 2>/dev/null || true
    fi
    if [ -n "$RULEGROUP_ID" ]; then
        aws_noauth wafv2 delete-rule-group --name "$RULEGROUP_NAME" --id "$RULEGROUP_ID" --lock-token "$RULEGROUP_LOCK_TOKEN" --scope REGIONAL 2>/dev/null || true
    fi
    if [ -n "$REGEXPATTERNSET_ID" ]; then
        aws_noauth wafv2 delete-regex-pattern-set --name "$REGEXPATTERNSET_NAME" --id "$REGEXPATTERNSET_ID" --lock-token "$REGEXPATTERNSET_LOCK_TOKEN" --scope REGIONAL 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Test 1: CreateWebACL
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: CreateWebACL... "
result=$(aws_noauth wafv2 create-web-acl \
    --name "$WEBACL_NAME" \
    --scope REGIONAL \
    --default-action '{"Allow":{}}' \
    --visibility-config '{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"test-metric"}' \
    --query 'Summary.Id' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    WEBACL_ID="$result"
    WEBACL_ARN=$(aws_noauth wafv2 get-web-acl --name "$WEBACL_NAME" --id "$WEBACL_ID" --scope REGIONAL --query 'WebACL.ARN' --output text 2>/dev/null)
    WEBACL_LOCK_TOKEN=$(aws_noauth wafv2 get-web-acl --name "$WEBACL_NAME" --id "$WEBACL_ID" --scope REGIONAL --query 'LockToken' --output text 2>/dev/null)
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 2: GetWebACL
run_test "2. GetWebACL" \
    "aws_noauth wafv2 get-web-acl --name '$WEBACL_NAME' --id '$WEBACL_ID' --scope REGIONAL --query 'WebACL.Name' --output text" \
    "$WEBACL_NAME"

# Test 3: ListWebACLs
run_test "3. ListWebACLs" \
    "aws_noauth wafv2 list-web-acls --scope REGIONAL --query 'WebACLs[?Name==\`$WEBACL_NAME\`].Name' --output text" \
    "$WEBACL_NAME"

# Test 4: UpdateWebACL
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: UpdateWebACL... "
result=$(aws_noauth wafv2 update-web-acl \
    --name "$WEBACL_NAME" \
    --id "$WEBACL_ID" \
    --scope REGIONAL \
    --lock-token "$WEBACL_LOCK_TOKEN" \
    --default-action '{"Block":{}}' \
    --visibility-config '{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"updated-metric"}' \
    --query 'NextLockToken' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    WEBACL_LOCK_TOKEN="$result"
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 5: CreateIPSet with Addresses (Bug #27)
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: CreateIPSet (with Addresses)... "
result=$(aws_noauth wafv2 create-ip-set \
    --name "$IPSET_NAME" \
    --scope REGIONAL \
    --ip-address-version IPV4 \
    --addresses '["192.0.2.0/24","203.0.113.0/24"]' \
    --query 'Summary.Id' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    IPSET_ID="$result"
    IPSET_LOCK_TOKEN=$(aws_noauth wafv2 get-ip-set --name "$IPSET_NAME" --id "$IPSET_ID" --scope REGIONAL --query 'LockToken' --output text 2>/dev/null)
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 6: GetIPSet - verify Addresses (Bug #27)
run_test "6. GetIPSet (verify Addresses)" \
    "aws_noauth wafv2 get-ip-set --name '$IPSET_NAME' --id '$IPSET_ID' --scope REGIONAL --query 'IPSet.Addresses' --output text" \
    "192.0.2.0/24"

# Test 7: ListIPSets
run_test "7. ListIPSets" \
    "aws_noauth wafv2 list-ip-sets --scope REGIONAL --query 'IPSets[?Name==\`$IPSET_NAME\`].Name' --output text" \
    "$IPSET_NAME"

# Test 8: UpdateIPSet
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: UpdateIPSet... "
result=$(aws_noauth wafv2 update-ip-set \
    --name "$IPSET_NAME" \
    --id "$IPSET_ID" \
    --scope REGIONAL \
    --lock-token "$IPSET_LOCK_TOKEN" \
    --addresses '["192.0.2.0/24","203.0.113.0/24","198.51.100.0/24"]' \
    --query 'NextLockToken' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    IPSET_LOCK_TOKEN="$result"
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 9: DeleteIPSet
run_test "9. DeleteIPSet" \
    "aws_noauth wafv2 delete-ip-set --name '$IPSET_NAME' --id '$IPSET_ID' --scope REGIONAL --lock-token '$IPSET_LOCK_TOKEN'" \
    ""
IPSET_ID=""

# Test 10: CreateRuleGroup with Rules (Bug #28)
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: CreateRuleGroup (with Rules)... "
result=$(aws_noauth wafv2 create-rule-group \
    --name "$RULEGROUP_NAME" \
    --scope REGIONAL \
    --capacity 100 \
    --visibility-config '{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"rule-group-metric"}' \
    --rules '[{"Name":"Rule1","Priority":1,"Statement":{"ByteMatchStatement":{"FieldToMatch":{"UriPath":{}},"PositionalConstraint":"STARTS_WITH","SearchString":"L2FkbWlu","TextTransformations":[{"Priority":0,"Type":"NONE"}]}},"Action":{"Allow":{}},"VisibilityConfig":{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"rule1-metric"}}]' \
    --query 'Summary.Id' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    RULEGROUP_ID="$result"
    RULEGROUP_LOCK_TOKEN=$(aws_noauth wafv2 get-rule-group --name "$RULEGROUP_NAME" --id "$RULEGROUP_ID" --scope REGIONAL --query 'LockToken' --output text 2>/dev/null)
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 11: GetRuleGroup - verify Rules (Bug #28)
run_test "11. GetRuleGroup (verify Rules)" \
    "aws_noauth wafv2 get-rule-group --name '$RULEGROUP_NAME' --id '$RULEGROUP_ID' --scope REGIONAL --query 'RuleGroup.Rules[0].Name' --output text" \
    "Rule1"

# Test 12: ListRuleGroups
run_test "12. ListRuleGroups" \
    "aws_noauth wafv2 list-rule-groups --scope REGIONAL --query 'RuleGroups[?Name==\`$RULEGROUP_NAME\`].Name' --output text" \
    "$RULEGROUP_NAME"

# Test 13: UpdateRuleGroup
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: UpdateRuleGroup... "
result=$(aws_noauth wafv2 update-rule-group \
    --name "$RULEGROUP_NAME" \
    --id "$RULEGROUP_ID" \
    --scope REGIONAL \
    --lock-token "$RULEGROUP_LOCK_TOKEN" \
    --visibility-config '{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"updated-rule-group-metric"}' \
    --rules '[{"Name":"Rule2","Priority":1,"Statement":{"ByteMatchStatement":{"FieldToMatch":{"UriPath":{}},"PositionalConstraint":"STARTS_WITH","SearchString":"L2FwaQ==","TextTransformations":[{"Priority":0,"Type":"NONE"}]}},"Action":{"Block":{}},"VisibilityConfig":{"SampledRequestsEnabled":true,"CloudWatchMetricsEnabled":true,"MetricName":"rule2-metric"}}]' \
    --query 'NextLockToken' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    RULEGROUP_LOCK_TOKEN="$result"
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 14: DeleteRuleGroup
run_test "14. DeleteRuleGroup" \
    "aws_noauth wafv2 delete-rule-group --name '$RULEGROUP_NAME' --id '$RULEGROUP_ID' --scope REGIONAL --lock-token '$RULEGROUP_LOCK_TOKEN'" \
    ""
RULEGROUP_ID=""

# Test 15: CreateRegexPatternSet
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: CreateRegexPatternSet... "
result=$(aws_noauth wafv2 create-regex-pattern-set \
    --name "$REGEXPATTERNSET_NAME" \
    --scope REGIONAL \
    --regular-expression-list '[{"RegexString":"^/admin/.*"}]' \
    --query 'Summary.Id' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    REGEXPATTERNSET_ID="$result"
    REGEXPATTERNSET_LOCK_TOKEN=$(aws_noauth wafv2 get-regex-pattern-set --name "$REGEXPATTERNSET_NAME" --id "$REGEXPATTERNSET_ID" --scope REGIONAL --query 'LockToken' --output text 2>/dev/null)
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 16: GetRegexPatternSet
run_test "16. GetRegexPatternSet" \
    "aws_noauth wafv2 get-regex-pattern-set --name '$REGEXPATTERNSET_NAME' --id '$REGEXPATTERNSET_ID' --scope REGIONAL --query 'RegexPatternSet.Name' --output text" \
    "$REGEXPATTERNSET_NAME"

# Test 17: ListRegexPatternSets
run_test "17. ListRegexPatternSets" \
    "aws_noauth wafv2 list-regex-pattern-sets --scope REGIONAL --query 'RegexPatternSets[?Name==\`$REGEXPATTERNSET_NAME\`].Name' --output text" \
    "$REGEXPATTERNSET_NAME"

# Test 18: UpdateRegexPatternSet
TEST_COUNT=$((TEST_COUNT + 1))
echo -n "Test $TEST_COUNT: UpdateRegexPatternSet... "
result=$(aws_noauth wafv2 update-regex-pattern-set \
    --name "$REGEXPATTERNSET_NAME" \
    --id "$REGEXPATTERNSET_ID" \
    --scope REGIONAL \
    --lock-token "$REGEXPATTERNSET_LOCK_TOKEN" \
    --regular-expression-list '[{"RegexString":"^/api/.*"}]' \
    --query 'NextLockToken' --output text 2>&1)
if [ $? -eq 0 ] && [ -n "$result" ] && [ "$result" != "None" ]; then
    REGEXPATTERNSET_LOCK_TOKEN="$result"
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}FAIL${NC}"
    echo "  Output: ${result:0:200}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Test 19: DeleteRegexPatternSet
run_test "19. DeleteRegexPatternSet" \
    "aws_noauth wafv2 delete-regex-pattern-set --name '$REGEXPATTERNSET_NAME' --id '$REGEXPATTERNSET_ID' --scope REGIONAL --lock-token '$REGEXPATTERNSET_LOCK_TOKEN'" \
    ""
REGEXPATTERNSET_ID=""

# Test 20: AssociateWebACL (Bug #25)
run_test "20. AssociateWebACL" \
    "aws_noauth wafv2 associate-web-acl --web-acl-arn '$WEBACL_ARN' --resource-arn '$RESOURCE_ARN'" \
    ""

# Test 21: ListResourcesAssociatedToWebACL
run_test "21. ListResourcesAssociatedToWebACL" \
    "aws_noauth wafv2 list-resources-for-web-acl --web-acl-arn '$WEBACL_ARN' --query 'ResourceArns[0]' --output text" \
    "$RESOURCE_ARN"

# Test 22: DisassociateWebACL
run_test "22. DisassociateWebACL" \
    "aws_noauth wafv2 disassociate-web-acl --resource-arn '$RESOURCE_ARN'" \
    ""

# Test 23: DeleteWebACL
run_test "23. DeleteWebACL" \
    "aws_noauth wafv2 delete-web-acl --name '$WEBACL_NAME' --id '$WEBACL_ID' --scope REGIONAL --lock-token '$WEBACL_LOCK_TOKEN'" \
    ""
WEBACL_ID=""

print_summary
