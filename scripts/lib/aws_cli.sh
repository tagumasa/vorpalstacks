#!/bin/bash
# AWS CLI wrapper functions for tests

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

# AWS CLI with authentication (for tests that need auth context)
# Usage: aws_cmd <service> <command> [args...]
aws_cmd() {
    local service="$1"
    shift
    
    AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-$DEFAULT_AWS_ACCESS_KEY_ID}" \
    AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-$DEFAULT_AWS_SECRET_ACCESS_KEY}" \
    aws "$service" "$@" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION"
}

# AWS CLI without authentication (for basic operations)
# Usage: aws_noauth <service> <command> [args...]
aws_noauth() {
    local service="$1"
    shift
    
    aws "$service" "$@" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        --no-sign-request
}

# SQS helpers
sqs_create_queue() {
    local queue_name="$1"
    aws_noauth sqs create-queue --queue-name "$queue_name" --query 'QueueUrl' --output text
}

sqs_get_queue_arn() {
    local queue_url="$1"
    aws_noauth sqs get-queue-attributes --queue-url "$queue_url" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text
}

# SNS helpers
sns_create_topic() {
    local topic_name="$1"
    aws_noauth sns create-topic --name "$topic_name" --query 'TopicArn' --output text
}

# DynamoDB helpers
ddb_create_table() {
    local table_name="$1"
    local pk="$2"
    local pk_type="${3:-S}"
    aws_noauth dynamodb create-table \
        --table-name "$table_name" \
        --attribute-definitions AttributeName="$pk",AttributeType="$pk_type" \
        --key-schema AttributeName="$pk",KeyType=HASH \
        --billing-mode PAY_PER_REQUEST
}

# S3 helpers
s3_create_bucket() {
    local bucket_name="$1"
    aws_noauth s3 mb "s3://$bucket_name" 2>/dev/null || true
}

# Lambda helpers
lambda_create_function() {
    local fn_name="$1"
    local zip_file="$2"
    local runtime="${3:-nodejs20.x}"
    local handler="${4:-index.handler}"
    
    aws_noauth lambda create-function \
        --function-name "$fn_name" \
        --runtime "$runtime" \
        --role "arn:aws:iam::$ACCOUNT_ID:role/lambda-role" \
        --handler "$handler" \
        --zip-file "fileb://$zip_file"
}

# Step Functions helpers
sfn_create_state_machine() {
    local name="$1"
    local definition="$2"
    
    aws_noauth stepfunctions create-state-machine \
        --name "$name" \
        --definition "$definition" \
        --role-arn "arn:aws:iam::$ACCOUNT_ID:role/test-role"
}

sfn_start_execution() {
    local sm_arn="$1"
    local input="${2:-{}}"
    
    aws_noauth stepfunctions start-execution \
        --state-machine-arn "$sm_arn" \
        --input "$input"
}

# Events helpers
events_create_event_bus() {
    local bus_name="$1"
    aws_noauth events create-event-bus --name "$bus_name"
}

events_put_rule() {
    local rule_name="$1"
    local bus_name="${2:-default}"
    local pattern="$3"
    
    aws_noauth events put-rule \
        --name "$rule_name" \
        --event-bus-name "$bus_name" \
        --event-pattern "$pattern"
}

# API Gateway helpers
apigw_create_rest_api() {
    local api_name="$1"
    aws_noauth apigateway create-rest-api --name "$api_name" --query 'id' --output text
}

# Generate ARN helper
make_arn() {
    local service="$1"
    local resource="$2"
    echo "arn:aws:$service:$REGION:$ACCOUNT_ID:$resource"
}

# Wait for execution to complete (Step Functions)
wait_for_execution() {
    local exec_arn="$1"
    local timeout="${2:-30}"
    local start_time=$(date +%s)
    
    while true; do
        local status=$(aws_noauth stepfunctions describe-execution --execution-arn "$exec_arn" --query 'status' --output text 2>/dev/null || echo "RUNNING")
        
        if [ "$status" != "RUNNING" ]; then
            echo "$status"
            return 0
        fi
        
        local current_time=$(date +%s)
        if [ $((current_time - start_time)) -ge $timeout ]; then
            echo "TIMEOUT"
            return 1
        fi
        
        sleep 1
    done
}
