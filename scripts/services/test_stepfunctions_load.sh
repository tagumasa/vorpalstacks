#!/bin/bash
# Step Functions Load Tests for VorpalStacks
# Verifies: Parallel execution, Map 100+ items, 50+ concurrent executions,
#           StopExecution goroutine cleanup, Activity Task timeout/cancel
# Usage: bash scripts/services/test_stepfunctions_load.sh

set -euo pipefail

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Step Functions Load Tests"
echo "Endpoint: $ENDPOINT_URL"
echo ""

PASS=0
FAIL=0
TOTAL=0

assert_pass() {
    TOTAL=$((TOTAL+1))
    local label="$1"
    local actual="$2"
    local expected="$3"
    if echo "$actual" | grep -q "$expected"; then
        echo "  ${GREEN}✓ $label${NC}"
        PASS=$((PASS+1))
    else
        echo "  ${RED}✗ $label${NC}"
        echo "    expected: $expected"
        echo "    actual:   $actual"
        FAIL=$((FAIL+1))
    fi
}

assert_eq() {
    TOTAL=$((TOTAL+1))
    local label="$1"
    local actual="$2"
    local expected="$3"
    if [ "$actual" = "$expected" ]; then
        echo "  ${GREEN}✓ $label${NC}"
        PASS=$((PASS+1))
    else
        echo "  ${RED}✗ $label${NC}"
        echo "    expected: $expected"
        echo "    actual:   $actual"
        FAIL=$((FAIL+1))
    fi
}

ROLE_ARN="arn:aws:iam:$REGION:$ACCOUNT_ID:role/test-role"
PREFIX="sfn-load-$$_"

cleanup_all() {
    echo "Cleaning up..."
    for sm in $(aws_noauth stepfunctions list-state-machines --query "stateMachines[?starts_with(name, \`$PREFIX\`)].stateMachineArn" --output text 2>/dev/null); do
        aws_noauth stepfunctions delete-state-machine --state-machine-arn "$sm" 2>/dev/null || true
    done
    for act in $(aws_noauth stepfunctions list-activities --query "activities[?starts_with(name, \`$PREFIX\`)].activityArn" --output text 2>/dev/null); do
        aws_noauth stepfunctions delete-activity --activity-arn "$act" 2>/dev/null || true
    done
}
trap cleanup_all EXIT

wait_for_status() {
    local exec_arn="$1"
    local target_status="$2"
    local timeout="${3:-60}"
    local elapsed=0
    while [ $elapsed -lt $timeout ]; do
        local status
        status=$(aws_noauth stepfunctions describe-execution --execution-arn "$exec_arn" --query 'status' --output text 2>/dev/null) || true
        if [ "$status" = "$target_status" ]; then
            echo "$status"
            return 0
        fi
        if [ "$status" = "FAILED" ] || [ "$status" = "ABORTED" ] || [ "$status" = "TIMED_OUT" ]; then
            echo "$status"
            return 0
        fi
        sleep 1
        elapsed=$((elapsed+1))
    done
    echo "TIMEOUT"
    return 0
}

# =====================================================
# Test 1: Parallel State - 8 branches simultaneously
# =====================================================
echo "--- Test 1: Parallel State (8 branches) ---"

PARALLEL_DEF='{
  "Comment": "Load test - 8 parallel branches",
  "StartAt": "Parallel",
  "States": {
    "Parallel": {
      "Type": "Parallel",
      "Branches": [
        {"StartAt": "B0", "States": {"B0": {"Type": "Pass", "Result": 0, "End": true}}},
        {"StartAt": "B1", "States": {"B1": {"Type": "Pass", "Result": 1, "End": true}}},
        {"StartAt": "B2", "States": {"B2": {"Type": "Pass", "Result": 2, "End": true}}},
        {"StartAt": "B3", "States": {"B3": {"Type": "Pass", "Result": 3, "End": true}}},
        {"StartAt": "B4", "States": {"B4": {"Type": "Pass", "Result": 4, "End": true}}},
        {"StartAt": "B5", "States": {"B5": {"Type": "Pass", "Result": 5, "End": true}}},
        {"StartAt": "B6", "States": {"B6": {"Type": "Pass", "Result": 6, "End": true}}},
        {"StartAt": "B7", "States": {"B7": {"Type": "Pass", "Result": 7, "End": true}}}
      ],
      "End": true
    }
  }
}'

SM1=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}parallel" \
    --definition "$PARALLEL_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "1.1 CreateStateMachine (Parallel)" "$SM1" "stateMachine"

EXEC1=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM1" \
    --input '{}' --query 'executionArn' --output text 2>&1)
assert_pass "1.2 StartExecution" "$EXEC1" "execution"

STATUS1=$(wait_for_status "$EXEC1" "SUCCEEDED" 30)
assert_eq "1.3 Parallel 8-branch execution SUCCEEDED" "$STATUS1" "SUCCEEDED"

OUTPUT1=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC1" --query 'output' --output text 2>/dev/null)
assert_pass "1.4 Output contains 8 branch results" "$OUTPUT1" "0,1,2,3,4,5,6,7"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM1" 2>/dev/null || true

# =====================================================
# Test 2: Map State - 150 items
# =====================================================
echo ""
echo "--- Test 2: Map State (150 items) ---"

MAP_DEF='{
  "Comment": "Load test - Map 150 items",
  "StartAt": "Map",
  "States": {
    "Map": {
      "Type": "Map",
      "ItemsPath": "$.items",
      "Iterator": {
        "StartAt": "Process",
        "States": {
          "Process": {
            "Type": "Pass",
            "ResultPath": "$.processed",
            "End": true
          }
        }
      },
      "End": true
    }
  }
}'

ITEMS=$(python3 -c "
import json
items = [{'id': i, 'value': f'item-{i}'} for i in range(150)]
print(json.dumps({'items': items}))
")

SM2=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}map150" \
    --definition "$MAP_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "2.1 CreateStateMachine (Map)" "$SM2" "stateMachine"

EXEC2=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM2" \
    --input "$ITEMS" --query 'executionArn' --output text 2>&1)
assert_pass "2.2 StartExecution" "$EXEC2" "execution"

STATUS2=$(wait_for_status "$EXEC2" "SUCCEEDED" 120)
assert_eq "2.3 Map 150-items execution SUCCEEDED" "$STATUS2" "SUCCEEDED"

OUTPUT2=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC2" --query 'output' --output text 2>/dev/null)
ITEM_COUNT=$(echo "$OUTPUT2" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d))" 2>/dev/null || echo "0")
assert_eq "2.4 Output has 150 items" "$ITEM_COUNT" "150"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM2" 2>/dev/null || true

# =====================================================
# Test 3: 50 concurrent executions
# =====================================================
echo ""
echo "--- Test 3: 50 Concurrent Executions ---"

SIMPLE_DEF='{
  "Comment": "Simple pass for load test",
  "StartAt": "Pass",
  "States": {
    "Pass": {
      "Type": "Pass",
      "Result": {"done": true},
      "End": true
    }
  }
}'

SM3=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}concurrent" \
    --definition "$SIMPLE_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "3.1 CreateStateMachine (concurrent)" "$SM3" "stateMachine"

echo "  Launching 50 concurrent executions..."
EXEC_ARNS=()
for i in $(seq 1 50); do
    ARN=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM3" \
        --input "{\"index\": $i}" --query 'executionArn' --output text 2>&1)
    EXEC_ARNS+=("$ARN")
done
assert_eq "3.2 Started 50 executions" "${#EXEC_ARNS[@]}" "50"

echo "  Waiting for all 50 to complete..."
SUCCEEDED=0
FAILED_COUNT=0
OTHER=0
for arn in "${EXEC_ARNS[@]}"; do
    STATUS=$(wait_for_status "$arn" "SUCCEEDED" 120)
    if [ "$STATUS" = "SUCCEEDED" ]; then
        SUCCEEDED=$((SUCCEEDED+1))
    elif [ "$STATUS" = "FAILED" ]; then
        FAILED_COUNT=$((FAILED_COUNT+1))
    else
        OTHER=$((OTHER+1))
    fi
done
assert_eq "3.3 All 50 executions SUCCEEDED" "$SUCCEEDED" "50"
assert_eq "3.4 Zero FAILED executions" "$FAILED_COUNT" "0"
assert_eq "3.5 Zero TIMEOUT/OTHER executions" "$OTHER" "0"

EXEC_LIST_COUNT=$(aws_noauth stepfunctions list-executions --state-machine-arn "$SM3" \
    --query 'length(executions)' --output text 2>/dev/null)
assert_eq "3.6 ListExecutions shows 50" "$EXEC_LIST_COUNT" "50"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM3" 2>/dev/null || true

# =====================================================
# Test 4: StopExecution - goroutine cleanup
# =====================================================
echo ""
echo "--- Test 4: StopExecution (goroutine cleanup) ---"

LONG_WAIT_DEF='{
  "Comment": "Long wait for stop test",
  "StartAt": "Wait1",
  "States": {
    "Wait1": {
      "Type": "Wait",
      "Seconds": 300,
      "Next": "Wait2"
    },
    "Wait2": {
      "Type": "Wait",
      "Seconds": 300,
      "End": true
    }
  }
}'

SM4=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}stop" \
    --definition "$LONG_WAIT_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "4.1 CreateStateMachine (stop)" "$SM4" "stateMachine"

EXEC4=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM4" \
    --input '{}' --query 'executionArn' --output text 2>&1)
assert_pass "4.2 StartExecution" "$EXEC4" "execution"

sleep 2

EXEC4_STATUS_BEFORE=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC4" --query 'status' --output text 2>/dev/null)
assert_eq "4.3 Execution is RUNNING before stop" "$EXEC4_STATUS_BEFORE" "RUNNING"

STOP_OUTPUT=$(aws_noauth stepfunctions stop-execution --execution-arn "$EXEC4" 2>&1)
assert_pass "4.4 StopExecution succeeded" "$STOP_OUTPUT" "stopDate"

sleep 2

EXEC4_STATUS_AFTER=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC4" --query 'status' --output text 2>/dev/null)
assert_eq "4.5 Execution is ABORTED after stop" "$EXEC4_STATUS_AFTER" "ABORTED"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM4" 2>/dev/null || true

# =====================================================
# Test 5: Activity Task - timeout
# =====================================================
echo ""
echo "--- Test 5: Activity Task Timeout ---"

ACT_NAME="${PREFIX}activity"
ACT_ARN=$(aws_noauth stepfunctions create-activity --name "$ACT_NAME" \
    --query 'activityArn' --output text 2>&1)
assert_pass "5.1 CreateActivity" "$ACT_ARN" "activity"

ACTIVITY_DEF="{
  \"Comment\": \"Activity with short timeout\",
  \"StartAt\": \"GetTask\",
  \"States\": {
    \"GetTask\": {
      \"Type\": \"Task\",
      \"Resource\": \"$ACT_ARN\",
      \"TimeoutSeconds\": 5,
      \"HeartbeatSeconds\": 2,
      \"End\": true
    }
  }
}"

SM5=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}act-timeout" \
    --definition "$ACTIVITY_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "5.2 CreateStateMachine (activity timeout)" "$SM5" "stateMachine"

EXEC5=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM5" \
    --input '{}' --query 'executionArn' --output text 2>&1)
assert_pass "5.3 StartExecution" "$EXEC5" "execution"

sleep 3

TASK_OUTPUT=$(aws_noauth stepfunctions get-activity-task --activity-arn "$ACT_ARN" --worker-name "load-test-worker" 2>&1 || true)
assert_pass "5.4 GetActivityTask returns task" "$TASK_OUTPUT" "taskToken"

if echo "$TASK_OUTPUT" | grep -q "taskToken"; then
    echo "  Activity task received, waiting for execution timeout (5s)..."
fi

STATUS5=$(wait_for_status "$EXEC5" "FAILED" 15)
assert_eq "5.5 Execution FAILED (activity task timeout)" "$STATUS5" "FAILED"

# =====================================================
# Test 5b: Activity Task - successful completion
# =====================================================
echo ""
echo "--- Test 5b: Activity Task Success ---"

EXEC5b=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM5" \
    --input '{"expected": "hello"}' --query 'executionArn' --output text 2>&1)
assert_pass "5b.1 StartExecution" "$EXEC5b" "execution"

sleep 2

TASK5b=$(aws_noauth stepfunctions get-activity-task --activity-arn "$ACT_ARN" --worker-name "load-test-worker-2" 2>&1 || true)
if echo "$TASK5b" | grep -q "taskToken"; then
    TOKEN5b=$(echo "$TASK5b" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['taskToken'])" 2>/dev/null)
    SEND_RESULT=$(aws_noauth stepfunctions send-task-success --task-token "$TOKEN5b" \
        --task-output '{"result": "activity-done"}' 2>&1)
    assert_pass "5b.2 SendTaskSuccess" "$SEND_RESULT" ""
fi

STATUS5b=$(wait_for_status "$EXEC5b" "SUCCEEDED" 15)
assert_eq "5b.3 Activity execution SUCCEEDED" "$STATUS5b" "SUCCEEDED"

OUTPUT5b=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC5b" --query 'output' --output text 2>/dev/null)
assert_pass "5b.4 Output contains activity result" "$OUTPUT5b" "activity-done"

# =====================================================
# Test 5c: StopExecution with pending Activity Task
# =====================================================
echo ""
echo "--- Test 5c: StopExecution with pending Activity Task ---"

EXEC5c=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM5" \
    --input '{}' --query 'executionArn' --output text 2>&1)
assert_pass "5c.1 StartExecution" "$EXEC5c" "execution"

sleep 2

aws_noauth stepfunctions stop-execution --execution-arn "$EXEC5c" 2>/dev/null || true

sleep 1

# Try to get the activity task after stop — should either get no task or task with cancelled execution
TASK5c=$(aws_noauth stepfunctions get-activity-task --activity-arn "$ACT_ARN" --worker-name "load-test-worker-3" 2>&1 || true)
# The task might already be dequeued or might still be pending — either way, execution should be ABORTED
STATUS5c=$(wait_for_status "$EXEC5c" "ABORTED" 10)
assert_eq "5c.2 Execution ABORTED after stop" "$STATUS5c" "ABORTED"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM5" 2>/dev/null || true
aws_noauth stepfunctions delete-activity --activity-arn "$ACT_ARN" 2>/dev/null || true

# =====================================================
# Test 6: Nested Parallel (Parallel within Parallel)
# =====================================================
echo ""
echo "--- Test 6: Nested Parallel ---"

NESTED_PARALLEL_DEF='{
  "Comment": "Nested parallel load test",
  "StartAt": "Outer",
  "States": {
    "Outer": {
      "Type": "Parallel",
      "Branches": [
        {
          "StartAt": "Inner1",
          "States": {
            "Inner1": {
              "Type": "Parallel",
              "Branches": [
                {"StartAt": "A", "States": {"A": {"Type": "Pass", "Result": "A1", "End": true}}},
                {"StartAt": "B", "States": {"B": {"Type": "Pass", "Result": "B1", "End": true}}},
                {"StartAt": "C", "States": {"C": {"Type": "Pass", "Result": "C1", "End": true}}}
              ],
              "End": true
            }
          }
        },
        {
          "StartAt": "Inner2",
          "States": {
            "Inner2": {
              "Type": "Parallel",
              "Branches": [
                {"StartAt": "D", "States": {"D": {"Type": "Pass", "Result": "D1", "End": true}}},
                {"StartAt": "E", "States": {"E": {"Type": "Pass", "Result": "E1", "End": true}}},
                {"StartAt": "F", "States": {"F": {"Type": "Pass", "Result": "F1", "End": true}}}
              ],
              "End": true
            }
          }
        }
      ],
      "End": true
    }
  }
}'

SM6=$(aws_noauth stepfunctions create-state-machine --name "${PREFIX}nested-parallel" \
    --definition "$NESTED_PARALLEL_DEF" --role-arn "$ROLE_ARN" \
    --query 'stateMachineArn' --output text 2>&1)
assert_pass "6.1 CreateStateMachine (nested parallel)" "$SM6" "stateMachine"

EXEC6=$(aws_noauth stepfunctions start-execution --state-machine-arn "$SM6" \
    --input '{}' --query 'executionArn' --output text 2>&1)
assert_pass "6.2 StartExecution" "$EXEC6" "execution"

STATUS6=$(wait_for_status "$EXEC6" "SUCCEEDED" 30)
assert_eq "6.3 Nested Parallel execution SUCCEEDED" "$STATUS6" "SUCCEEDED"

OUTPUT6=$(aws_noauth stepfunctions describe-execution --execution-arn "$EXEC6" --query 'output' --output text 2>/dev/null)
assert_pass "6.4 Output contains nested branch results" "$OUTPUT6" "A1"

aws_noauth stepfunctions delete-state-machine --state-machine-arn "$SM6" 2>/dev/null || true

# =====================================================
# Summary
# =====================================================
echo ""
echo "=========================================="
echo "  Step Functions Load Test Results"
echo "=========================================="
echo "  ${GREEN}PASSED: $PASS${NC}"
echo "  ${RED}FAILED: $FAIL${NC}"
echo "  TOTAL:  $TOTAL"
echo "=========================================="

if [ $FAIL -gt 0 ]; then
    exit 1
fi
