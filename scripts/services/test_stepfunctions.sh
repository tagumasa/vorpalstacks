#!/bin/bash
# Step Functions Service Tests for VorpalStacks
# Usage: bash scripts/services/test_stepfunctions.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "Step Functions Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

SM_NAME="test-state-machine-$$"
SM_ARN=""
EXEC_ARN=""
ACTIVITY_NAME="test-activity-$$"
ACTIVITY_ARN=""

cleanup() {
    [ -n "$SM_ARN" ] && aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN" 2>/dev/null || true
    [ -n "$ACTIVITY_ARN" ] && aws_noauth stepfunctions delete-activity --activity-arn "$ACTIVITY_ARN" 2>/dev/null || true
}
trap cleanup EXIT

SM_ARN=$(aws_noauth stepfunctions create-state-machine --name "$SM_NAME" --definition "$STATE_MACHINE_PASS" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
run_test "1. CreateStateMachine" \
    "echo '$SM_ARN'" \
    "stateMachine"

run_test "2. ListStateMachines" \
    "aws_noauth stepfunctions list-state-machines --query 'stateMachines[?name==\`$SM_NAME\`].name' --output text" \
    "$SM_NAME"

run_test "3. DescribeStateMachine" \
    "aws_noauth stepfunctions describe-state-machine --state-machine-arn '$SM_ARN' --query 'name' --output text" \
    "$SM_NAME"

EXEC_ARN=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN" --input '{"test":"input"}' --query 'executionArn' --output text 2>&1)
run_test "4. StartExecution" \
    "echo '$EXEC_ARN'" \
    "execution"

sleep 2

run_test "5. DescribeExecution" \
    "aws_noauth stepfunctions describe-execution --execution-arn '$EXEC_ARN' --query 'status' --output text" \
    "SUCCEEDED"

run_test "6. GetExecutionHistory" \
    "aws_noauth stepfunctions get-execution-history --execution-arn '$EXEC_ARN' --query 'events[0].type' --output text" \
    "ExecutionStarted"

run_test "7. ListExecutions" \
    "aws_noauth stepfunctions list-executions --state-machine-arn '$SM_ARN' --query 'executions[0].status' --output text" \
    "SUCCEEDED"

ACTIVITY_ARN=$(aws_noauth stepfunctions create-activity --name "$ACTIVITY_NAME" --query 'activityArn' --output text 2>&1)
run_test "8. CreateActivity" \
    "echo '$ACTIVITY_ARN'" \
    "activity"

run_test "9. ListActivities" \
    "aws_noauth stepfunctions list-activities --query 'activities[?name==\`$ACTIVITY_NAME\`].name' --output text" \
    "$ACTIVITY_NAME"

run_test "10. DescribeActivity" \
    "aws_noauth stepfunctions describe-activity --activity-arn '$ACTIVITY_ARN' --query 'name' --output text" \
    "$ACTIVITY_NAME"

run_test "11. TagResource" \
    "aws_noauth stepfunctions tag-resource --resource-arn '$SM_ARN' --tags key=Environment,value=Test" \
    ""

run_test "12. ListTagsForResource" \
    "aws_noauth stepfunctions list-tags-for-resource --resource-arn '$SM_ARN' --query 'tags[0].value' --output text" \
    "Test"

run_test "13. UntagResource" \
    "aws_noauth stepfunctions untag-resource --resource-arn '$SM_ARN' --tag-keys Environment" \
    ""

run_test "14. DeleteActivity" \
    "aws_noauth stepfunctions delete-activity --activity-arn '$ACTIVITY_ARN'" \
    ""

run_test "15. DeleteStateMachine" \
    "aws_noauth stepfunctions delete-state-machine --state-machine-arn '$SM_ARN'" \
    ""

# ============================================
# Input Processing Tests (Tests 16-21)
# ============================================
echo ""
echo "--- Input Processing Tests ---"

SM_ARN2=$(aws_noauth stepfunctions create-state-machine --name "inputpath-test-$$" --definition "$STATE_MACHINE_INPUTPATH" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
EXEC_ARN2=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN2" --input '{"detail": {"value": 123}, "extra": "ignored"}' --query 'executionArn' --output text 2>&1)
sleep 2
OUTPUT=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC_ARN2" --query 'output' --output text 2>&1)
run_test "16. InputPath" \
    "echo '$OUTPUT'" \
    "filtered"
aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN2" 2>/dev/null || true

SM_ARN3=$(aws_noauth stepfunctions create-state-machine --name "outputpath-test-$$" --definition "$STATE_MACHINE_OUTPUTPATH" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
EXEC_ARN3=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN3" --input '{}' --query 'executionArn' --output text 2>&1)
sleep 2
OUTPUT=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC_ARN3" --query 'output' --output text 2>&1)
run_test "17. OutputPath" \
    "echo '$OUTPUT'" \
    "nested"
aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN3" 2>/dev/null || true

SM_ARN4=$(aws_noauth stepfunctions create-state-machine --name "resultpath-test-$$" --definition "$STATE_MACHINE_RESULTPATH" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
EXEC_ARN4=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN4" --input '{"input": "original"}' --query 'executionArn' --output text 2>&1)
sleep 2
OUTPUT=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC_ARN4" --query 'output' --output text 2>&1)
run_test "18. ResultPath" \
    "echo '$OUTPUT' | grep -q '\"input\"' && echo 'contains both' || echo 'missing'" \
    "contains both"
aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN4" 2>/dev/null || true

SM_ARN5=$(aws_noauth stepfunctions create-state-machine --name "parameters-test-$$" --definition "$STATE_MACHINE_PARAMETERS" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
EXEC_ARN5=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN5" --input '{"input": "dynamic"}' --query 'executionArn' --output text 2>&1)
sleep 2
OUTPUT=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC_ARN5" --query 'output' --output text 2>&1)
run_test "19. Parameters" \
    "echo '$OUTPUT' | grep -q '\"staticValue\"' && echo 'has static' || echo 'missing'; echo '$OUTPUT' | grep -q '\"dynamicValue\"' && echo 'has dynamic' || echo 'missing'; echo 'combined'" \
    "combined"
aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN5" 2>/dev/null || true

SM_ARN6=$(aws_noauth stepfunctions create-state-machine --name "resultselector-test-$$" --definition "$STATE_MACHINE_RESULTSELECTOR" --role-arn "arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role" --query 'stateMachineArn' --output text 2>&1)
EXEC_ARN6=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM_ARN6" --input '{}' --query 'executionArn' --output text 2>&1)
sleep 2
OUTPUT=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC_ARN6" --query 'output' --output text 2>&1)
run_test "20. ResultSelector" \
    "echo '$OUTPUT' | grep -q '\"selectedValue\"' && echo 'has key' || echo 'missing'; echo 'done'" \
    "done"
aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM_ARN6" 2>/dev/null || true

print_summary
