#!/bin/bash
# Test fixtures and common test data

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

# Generate unique test resource names (to avoid conflicts in parallel tests)
make_test_name() {
    local prefix="$1"
    local suffix="${2:-$$}"
    echo "${prefix}-${suffix}"
}

# Common state machine definitions
STATE_MACHINE_PASS='{
  "Comment": "Simple Pass state",
  "StartAt": "Hello",
  "States": {
    "Hello": {
      "Type": "Pass",
      "Result": "Hello World!",
      "End": true
    }
  }
}'

STATE_MACHINE_CHOICE='{
  "Comment": "Choice state example",
  "StartAt": "CheckValue",
  "States": {
    "CheckValue": {
      "Type": "Choice",
      "Choices": [
        {"Variable": "$.value", "NumericEquals": 1, "Next": "ValueIsOne"},
        {"Variable": "$.value", "NumericEquals": 2, "Next": "ValueIsTwo"}
      ],
      "Default": "DefaultValue"
    },
    "ValueIsOne": {"Type": "Pass", "Result": "One", "End": true},
    "ValueIsTwo": {"Type": "Pass", "Result": "Two", "End": true},
    "DefaultValue": {"Type": "Pass", "Result": "Default", "End": true}
  }
}'

STATE_MACHINE_WAIT='{
  "Comment": "Wait state example",
  "StartAt": "Wait",
  "States": {
    "Wait": {
      "Type": "Wait",
      "Seconds": 1,
      "End": true
    }
  }
}'

STATE_MACHINE_PARALLEL='{
  "Comment": "Parallel state example",
  "StartAt": "Parallel",
  "States": {
    "Parallel": {
      "Type": "Parallel",
      "Branches": [
        {"StartAt": "Branch1", "States": {"Branch1": {"Type": "Pass", "Result": "A", "End": true}}},
        {"StartAt": "Branch2", "States": {"Branch2": {"Type": "Pass", "Result": "B", "End": true}}}
      ],
      "End": true
    }
  }
}'

STATE_MACHINE_MAP='{
  "Comment": "Map state example",
  "StartAt": "Map",
  "States": {
    "Map": {
      "Type": "Map",
      "ItemsPath": "$.items",
      "Iterator": {
        "StartAt": "Process",
        "States": {
          "Process": {"Type": "Pass", "ResultPath": "$.processed", "End": true}
        }
      },
      "End": true
    }
  }
}'

STATE_MACHINE_SQS='{
  "Comment": "SQS integration",
  "StartAt": "SendMessage",
  "States": {
    "SendMessage": {
      "Type": "Task",
      "Resource": "arn:aws:states:::sqs:sendMessage",
      "Parameters": {"QueueName": "TEST_QUEUE", "MessageBody": "Hello from Step Functions!"},
      "End": true
    }
  }
}'

STATE_MACHINE_SNS='{
  "Comment": "SNS integration",
  "StartAt": "Publish",
  "States": {
    "Publish": {
      "Type": "Task",
      "Resource": "arn:aws:states:::sns:publish",
      "Parameters": {"TopicName": "TEST_TOPIC", "Message": "Hello from Step Functions!"},
      "End": true
    }
  }
}'

STATE_MACHINE_INPUTPATH='{
  "Comment": "InputPath test",
  "StartAt": "FilterInput",
  "States": {
    "FilterInput": {
      "Type": "Pass",
      "InputPath": "$.detail",
      "Result": {"filtered": true},
      "End": true
    }
  }
}'

STATE_MACHINE_OUTPUTPATH='{
  "Comment": "OutputPath test",
  "StartAt": "FilterOutput",
  "States": {
    "FilterOutput": {
      "Type": "Pass",
      "Result": {"full": {"nested": "value"}, "extra": "data"},
      "OutputPath": "$.full",
      "End": true
    }
  }
}'

STATE_MACHINE_RESULTPATH='{
  "Comment": "ResultPath test",
  "StartAt": "AddResult",
  "States": {
    "AddResult": {
      "Type": "Pass",
      "Result": {"status": "processed"},
      "ResultPath": "$.result",
      "End": true
    }
  }
}'

STATE_MACHINE_PARAMETERS='{
  "Comment": "Parameters test",
  "StartAt": "WithParams",
  "States": {
    "WithParams": {
      "Type": "Pass",
      "Parameters": {
        "staticValue": "hello",
        "dynamicValue.$": "$.input",
        "nested": {"key.$": "$.input"}
      },
      "End": true
    }
  }
}'

STATE_MACHINE_RESULTSELECTOR='{
  "Comment": "ResultSelector test",
  "StartAt": "SelectResult",
  "States": {
    "SelectResult": {
      "Type": "Pass",
      "Result": {"data": {"value": 42, "name": "test"}},
      "ResultSelector": {
        "selectedValue.$": "$.data.value",
        "selectedName.$": "$.data.name"
      },
      "End": true
    }
  }
}'

# Event patterns
EVENT_PATTERN_SIMPLE='{"source": ["test.source"]}'
EVENT_PATTERN_DETAIL='{"source": ["test.source"], "detail-type": ["test-event"], "detail": {"status": ["success"]}}'

# DynamoDB test items
DDB_ITEM_SIMPLE='{"id": {"S": "test1"}, "name": {"S": "Test Item"}, "count": {"N": "42"}}'
DDB_ITEM_COMPLEX='{"id": {"S": "test2"}, "data": {"M": {"key1": {"S": "value1"}, "key2": {"N": "123"}}}, "tags": {"SS": ["tag1", "tag2"]}}'

# S3 test content
S3_TEST_CONTENT="Hello, S3 from VorpalStacks! $(date +%s)"

# Lambda test payload
LAMBDA_PAYLOAD='{"key": "value", "number": 42}'

# SQS message body
SQS_MESSAGE_BODY="Hello SQS! $(date +%s)"

# SNS message
SNS_MESSAGE="Hello SNS! $(date +%s)"

# Create test file
create_test_file() {
    local content="${1:-$S3_TEST_CONTENT}"
    local file_path="${2:-$PROJECT_ROOT/tmp/test-file-$$}"
    
    mkdir -p "$(dirname "$file_path")"
    echo "$content" > "$file_path"
    echo "$file_path"
}

# Create Lambda zip file (minimal handler)
create_lambda_zip() {
    local output_path="${1:-$PROJECT_ROOT/tmp/lambda/function.zip}"
    local handler_code="${2:-exports.handler = async (event) => ({ statusCode: 200, body: JSON.stringify(event) });}"
    
    mkdir -p "$(dirname "$output_path")"
    
    local tmp_dir=$(mktemp -d)
    echo "$handler_code" > "$tmp_dir/index.js"
    (cd "$tmp_dir" && zip -q -r "$output_path" index.js)
    rm -rf "$tmp_dir"
    
    echo "$output_path"
}

# Cleanup test resources helper
cleanup_aws_resources() {
    local service="$1"
    
    case "$service" in
        sqs)
            for queue in $(aws_noauth sqs list-queues --query 'QueueUrls' --output text 2>/dev/null); do
                aws_noauth sqs delete-queue --queue-url "$queue" 2>/dev/null || true
            done
            ;;
        sns)
            for topic in $(aws_noauth sns list-topics --query 'Topics[].TopicArn' --output text 2>/dev/null); do
                aws_noauth sns delete-topic --topic-arn "$topic" 2>/dev/null || true
            done
            ;;
        dynamodb)
            for table in $(aws_noauth dynamodb list-tables --query 'TableNames' --output text 2>/dev/null); do
                aws_noauth dynamodb delete-table --table-name "$table" 2>/dev/null || true
            done
            ;;
        s3)
            for bucket in $(aws_noauth s3 ls --query 'Buckets[].Name' --output text 2>/dev/null); do
                aws_noauth s3 rm "s3://$bucket" --recursive 2>/dev/null || true
                aws_noauth s3 rb "s3://$bucket" 2>/dev/null || true
            done
            ;;
        stepfunctions)
            for sm in $(aws_noauth stepfunctions list-state-machines --query 'stateMachines[].stateMachineArn' --output text 2>/dev/null); do
                aws_noauth stepfunctions delete-state-machine --state-machine-arn "$sm" 2>/dev/null || true
            done
            ;;
        events)
            for bus in $(aws_noauth events list-event-buses --query 'EventBuses[].Name' --output text 2>/dev/null); do
                [ "$bus" != "default" ] && aws_noauth events delete-event-bus --name "$bus" 2>/dev/null || true
            done
            ;;
        lambda)
            for fn in $(aws_noauth lambda list-functions --query 'Functions[].FunctionName' --output text 2>/dev/null); do
                aws_noauth lambda delete-function --function-name "$fn" 2>/dev/null || true
            done
            ;;
        apigateway)
            for api in $(aws_noauth apigateway get-rest-apis --query 'items[].id' --output text 2>/dev/null); do
                aws_noauth apigateway delete-rest-api --rest-api-id "$api" 2>/dev/null || true
            done
            ;;
    esac
}
