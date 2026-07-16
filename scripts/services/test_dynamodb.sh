#!/bin/bash
# DynamoDB Service Tests for VorpalStacks
# Usage: bash scripts/services/test_dynamodb.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"
source "$(dirname "$0")/../lib/fixtures.sh"

print_header "DynamoDB Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

TABLE_NAME="TestTable-$$"
TABLE_ARN="arn:aws:dynamodb:$REGION:$ACCOUNT_ID:table/$TABLE_NAME"

cleanup() {
    aws_noauth dynamodb delete-table --table-name "$TABLE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. CreateTable" \
    "aws_noauth dynamodb create-table --table-name '$TABLE_NAME' --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --query 'TableDescription.TableName' --output text" \
    "$TABLE_NAME"

run_test "2. DescribeTable" \
    "aws_noauth dynamodb describe-table --table-name '$TABLE_NAME' --query 'Table.TableName' --output text" \
    "$TABLE_NAME"

run_test "3. ListTables" \
    "aws_noauth dynamodb list-tables --query \"TableNames[?@=='$TABLE_NAME'] | [0]\" --output text" \
    "$TABLE_NAME"

run_test "4. PutItem" \
    "aws_noauth dynamodb put-item --table-name '$TABLE_NAME' --item '{\"id\":{\"S\":\"test1\"},\"name\":{\"S\":\"Test Item\"},\"count\":{\"N\":\"42\"}}'" \
    ""

run_test "5. GetItem" \
    "aws_noauth dynamodb get-item --table-name '$TABLE_NAME' --key '{\"id\":{\"S\":\"test1\"}}' --query 'Item.name.S' --output text" \
    "Test Item"

run_test "6. UpdateItem" \
    "aws_noauth dynamodb update-item --table-name '$TABLE_NAME' --key '{\"id\":{\"S\":\"test1\"}}' --update-expression 'SET #n = :name' --expression-attribute-names '{\"#n\":\"name\"}' --expression-attribute-values '{\":name\":{\"S\":\"Updated\"}}' --return-values ALL_NEW --query 'Attributes.name.S' --output text" \
    "Updated"

run_test "7. Query" \
    "aws_noauth dynamodb query --table-name '$TABLE_NAME' --key-condition-expression 'id = :id' --expression-attribute-values '{\":id\":{\"S\":\"test1\"}}' --query 'Count' --output text" \
    "1"

run_test "8. Scan" \
    "aws_noauth dynamodb scan --table-name '$TABLE_NAME' --query 'Count' --output text" \
    "1"

run_test "9. PutItem (batch)" \
    "aws_noauth dynamodb batch-write-item --request-items '{\"$TABLE_NAME\":[{\"PutRequest\":{\"Item\":{\"id\":{\"S\":\"batch1\"},\"data\":{\"S\":\"batch item 1\"}}}},{\"PutRequest\":{\"Item\":{\"id\":{\"S\":\"batch2\"},\"data\":{\"S\":\"batch item 2\"}}}}]}'" \
    ""

run_test "10. BatchGetItem" \
    "aws_noauth dynamodb batch-get-item --request-items '{\"$TABLE_NAME\":{\"Keys\":[{\"id\":{\"S\":\"batch1\"}},{\"id\":{\"S\":\"batch2\"}}]}}' --query 'length(Responses) > \`0\`' --output text" \
    "True"

run_test "11. DeleteItem" \
    "aws_noauth dynamodb delete-item --table-name '$TABLE_NAME' --key '{\"id\":{\"S\":\"test1\"}}'" \
    ""

run_test "12. TagResource" \
    "aws_noauth dynamodb tag-resource --resource-arn '$TABLE_ARN' --tags Key=Environment,Value=Test" \
    ""

run_test "13. ListTagsForResource" \
    "aws_noauth dynamodb list-tags-of-resource --resource-arn '$TABLE_ARN' --query 'Tags[0].Key' --output text" \
    "Environment"

run_test "14. UntagResource" \
    "aws_noauth dynamodb untag-resource --resource-arn '$TABLE_ARN' --tag-keys Environment" \
    ""

run_test "15. UpdateTimeToLive" \
    "aws_noauth dynamodb update-time-to-live --table-name '$TABLE_NAME' --time-to-live-specification Enabled=true,AttributeName=ttl --query 'TimeToLiveSpecification.Enabled' --output text" \
    "True"

run_test "16. DescribeTimeToLive" \
    "aws_noauth dynamodb describe-time-to-live --table-name '$TABLE_NAME' --query 'TimeToLiveDescription.TimeToLiveStatus' --output text" \
    "ENABLED"

run_test "17. TransactWriteItems" \
    "aws_noauth dynamodb transact-write-items --transact-items '[{\"Put\":{\"TableName\":\"$TABLE_NAME\",\"Item\":{\"id\":{\"S\":\"tx1\"},\"data\":{\"S\":\"transaction test\"}}}}]'" \
    ""

run_test "18. TransactGetItems" \
    "aws_noauth dynamodb transact-get-items --transact-items '[{\"Get\":{\"TableName\":\"$TABLE_NAME\",\"Key\":{\"id\":{\"S\":\"tx1\"}}}}]' --query 'Responses[0].Item.data.S' --output text" \
    "transaction test"

run_test "19. CreateBackup" \
    "aws_noauth dynamodb create-backup --table-name '$TABLE_NAME' --backup-name 'test-backup-$$' --query 'BackupDetails.BackupName' --output text" \
    "test-backup"

run_test "20. ListBackups" \
    "aws_noauth dynamodb list-backups --table-name '$TABLE_NAME' --query 'BackupSummaries[0].BackupName' --output text" \
    "test-backup"

run_test "21. DescribeContinuousBackups" \
    "aws_noauth dynamodb describe-continuous-backups --table-name '$TABLE_NAME' --query 'ContinuousBackupsDescription.ContinuousBackupsStatus' --output text" \
    ""

run_test "22. UpdateContinuousBackups" \
    "aws_noauth dynamodb update-continuous-backups --table-name '$TABLE_NAME' --point-in-time-recovery-specification PointInTimeRecoveryEnabled=true --query 'ContinuousBackupsDescription.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus' --output text" \
    ""

run_test "23. ExecuteStatement (PartiQL)" \
    "aws_noauth dynamodb execute-statement --statement \"SELECT * FROM \\\"$TABLE_NAME\\\" WHERE id='tx1'\" --query 'length(Items)' --output text" \
    "1"

run_test "24. BatchExecuteStatement (PartiQL)" \
    "aws_noauth dynamodb batch-execute-statement --statements '[{\"Statement\":\"SELECT * FROM \\\"$TABLE_NAME\\\" WHERE id=''tx1''\"}]'" \
    ""

run_test "25. DescribeEndpoints" \
    "aws_noauth dynamodb describe-endpoints --query 'Endpoints[0].Address' --output text" \
    ""

run_test "26. DescribeLimits" \
    "aws_noauth dynamodb describe-limits --query 'AccountMaxReadCapacityUnits' --output text" \
    ""

run_test "27. DeleteTable" \
    "aws_noauth dynamodb delete-table --table-name '$TABLE_NAME' --query 'TableDescription.TableName' --output text" \
    "$TABLE_NAME"

print_summary
