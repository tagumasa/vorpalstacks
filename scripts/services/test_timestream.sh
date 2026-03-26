#!/bin/bash
# Timestream Service Tests for VorpalStacks
# Usage: bash scripts/services/test_timestream.sh

source "$(dirname "$0")/../lib/config.sh"
source "$(dirname "$0")/../lib/colors.sh"
source "$(dirname "$0")/../lib/assertions.sh"
source "$(dirname "$0")/../lib/aws_cli.sh"

print_header "Timestream Service Tests"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

DATABASE_NAME="test_db_$$"
TABLE_NAME="test_table_$$"

cleanup() {
    aws_noauth timestream-write delete-table --database-name "$DATABASE_NAME" --table-name "$TABLE_NAME" 2>/dev/null || true
    aws_noauth timestream-write delete-database --database-name "$DATABASE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

run_test "1. DescribeEndpoints (write)" \
    "aws_noauth timestream-write describe-endpoints --query 'Endpoints[0].Address' --output text" \
    ""

run_test "2. CreateDatabase" \
    "aws_noauth timestream-write create-database --database-name '$DATABASE_NAME' --query 'Database.DatabaseName' --output text" \
    "$DATABASE_NAME"

run_test "3. DescribeDatabase" \
    "aws_noauth timestream-write describe-database --database-name '$DATABASE_NAME' --query 'Database.DatabaseName' --output text" \
    "$DATABASE_NAME"

run_test "4. ListDatabases" \
    "aws_noauth timestream-write list-databases --query 'Databases[?DatabaseName==\`$DATABASE_NAME\`].DatabaseName' --output text" \
    "$DATABASE_NAME"

run_test "5. UpdateDatabase" \
    "aws_noauth timestream-write update-database --database-name '$DATABASE_NAME' --kms-key-id 'alias/aws/timestream'" \
    ""

run_test "6. CreateTable" \
    "aws_noauth timestream-write create-table --database-name '$DATABASE_NAME' --table-name '$TABLE_NAME' --query 'Table.TableName' --output text" \
    "$TABLE_NAME"

run_test "7. DescribeTable" \
    "aws_noauth timestream-write describe-table --database-name '$DATABASE_NAME' --table-name '$TABLE_NAME' --query 'Table.TableName' --output text" \
    "$TABLE_NAME"

run_test "8. ListTables" \
    "aws_noauth timestream-write list-tables --database-name '$DATABASE_NAME' --query 'Tables[?TableName==\`$TABLE_NAME\`].TableName' --output text" \
    "$TABLE_NAME"

run_test "9. WriteRecords" \
    "aws_noauth timestream-write write-records --database-name '$DATABASE_NAME' --table-name '$TABLE_NAME' --records '[{\"MeasureName\":\"cpu\",\"MeasureValue\":\"50\",\"Time\":\"$(date +%s)000\",\"Dimensions\":[{\"Name\":\"host\",\"Value\":\"server1\"}]}]'" \
    ""

run_test "10. UpdateTable" \
    "aws_noauth timestream-write update-table --database-name '$DATABASE_NAME' --table-name '$TABLE_NAME' --retention-properties 'MemoryStoreRetentionPeriodInHours=1,MagneticStoreRetentionPeriodInDays=1'" \
    ""

run_test "11. DeleteTable" \
    "aws_noauth timestream-write delete-table --database-name '$DATABASE_NAME' --table-name '$TABLE_NAME'" \
    ""

run_test "12. DeleteDatabase" \
    "aws_noauth timestream-write delete-database --database-name '$DATABASE_NAME'" \
    ""

run_test "13. DescribeEndpoints (query)" \
    "aws_noauth timestream-query describe-endpoints --query 'Endpoints[0].Address' --output text" \
    ""

print_summary
