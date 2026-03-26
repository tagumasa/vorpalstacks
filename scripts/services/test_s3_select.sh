#!/bin/bash
set -e

source "$(dirname "${BASH_SOURCE[0]}")/../lib/aws_cli.sh"

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:8080}"
BUCKET="test-select-bucket"
PREFIX="s3-select-test"

echo "=== S3 SelectObjectContent Integration Test ==="

# Cleanup
aws_noauth s3 rb "s3://${BUCKET}" --force 2>/dev/null || true

# Create bucket
echo "Creating bucket..."
aws_noauth s3 mb "s3://${BUCKET}"

# Test 1: CSV Select
echo ""
echo "=== Test 1: CSV Select ==="

CSV_KEY="${PREFIX}/test.csv"
CSV_DATA="name,age,city
Alice,30,Tokyo
Bob,25,Osaka
Charlie,35,Tokyo
David,40,Kyoto"

echo "Uploading CSV file..."
echo "$CSV_DATA" | aws_noauth s3 cp - "s3://${BUCKET}/${CSV_KEY}"

echo "Selecting all rows..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${CSV_KEY}" \
    --expression "SELECT * FROM s3object" \
    --expression-type SQL \
    --input-serialization '{"CSV": {"FileHeaderInfo": "USE"}}' \
    --output-serialization '{"CSV": {}}' \
    \
    "output-test1-all.txt" || echo "FAILED: Select all"

if [ -f "output-test1-all.txt" ]; then
    echo "Result:"
    cat output-test1-all.txt
    rm output-test1-all.txt
fi

echo "Selecting with WHERE clause..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${CSV_KEY}" \
    --expression "SELECT name, age FROM s3object WHERE city = 'Tokyo'" \
    --expression-type SQL \
    --input-serialization '{"CSV": {"FileHeaderInfo": "USE"}}' \
    --output-serialization '{"CSV": {}}' \
    \
    "output-test1-where.txt" || echo "FAILED: Select with WHERE"

if [ -f "output-test1-where.txt" ]; then
    echo "Result:"
    cat output-test1-where.txt
    rm output-test1-where.txt
fi

# Test 2: JSON Select
echo ""
echo "=== Test 2: JSON Select ==="

JSON_KEY="${PREFIX}/test.jsonl"
JSON_DATA='{"name":"Alice","age":30,"address":{"city":"Tokyo"}}
{"name":"Bob","age":25,"address":{"city":"Osaka"}}
{"name":"Charlie","age":35,"address":{"city":"Tokyo"}}
{"name":"David","age":40,"address":{"city":"Kyoto"}}'

echo "Uploading JSON file..."
echo "$JSON_DATA" | aws_noauth s3 cp - "s3://${BUCKET}/${JSON_KEY}"

echo "Selecting all rows..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${JSON_KEY}" \
    --expression "SELECT * FROM s3object" \
    --expression-type SQL \
    --input-serialization '{"JSON": {"Type": "LINES"}}' \
    --output-serialization '{"JSON": {}}' \
    \
    "output-test2-all.txt" || echo "FAILED: Select all JSON"

if [ -f "output-test2-all.txt" ]; then
    echo "Result:"
    cat output-test2-all.txt
    rm output-test2-all.txt
fi

echo "Selecting with nested field..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${JSON_KEY}" \
    --expression "SELECT name FROM s3object WHERE address.city = 'Tokyo'" \
    --expression-type SQL \
    --input-serialization '{"JSON": {"Type": "LINES"}}' \
    --output-serialization '{"JSON": {}}' \
    \
    "output-test2-nested.txt" || echo "FAILED: Select with nested field"

if [ -f "output-test2-nested.txt" ]; then
    echo "Result:"
    cat output-test2-nested.txt
    rm output-test2-nested.txt
fi

# Test 3: LIKE clause
echo ""
echo "=== Test 3: LIKE clause ==="

echo "Selecting with LIKE..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${CSV_KEY}" \
    --expression "SELECT name FROM s3object WHERE name LIKE 'A%'" \
    --expression-type SQL \
    --input-serialization '{"CSV": {"FileHeaderInfo": "USE"}}' \
    --output-serialization '{"CSV": {}}' \
    \
    "output-test3-like.txt" || echo "FAILED: Select with LIKE"

if [ -f "output-test3-like.txt" ]; then
    echo "Result:"
    cat output-test3-like.txt
    rm output-test3-like.txt
fi

# Test 4: IN clause
echo ""
echo "=== Test 4: IN clause ==="

echo "Selecting with IN..."
aws_noauth s3api select-object-content \
    --bucket "${BUCKET}" \
    --key "${CSV_KEY}" \
    --expression "SELECT name, city FROM s3object WHERE city IN ('Tokyo', 'Osaka')" \
    --expression-type SQL \
    --input-serialization '{"CSV": {"FileHeaderInfo": "USE"}}' \
    --output-serialization '{"CSV": {}}' \
    \
    "output-test4-in.txt" || echo "FAILED: Select with IN"

if [ -f "output-test4-in.txt" ]; then
    echo "Result:"
    cat output-test4-in.txt
    rm output-test4-in.txt
fi

# Cleanup
echo ""
echo "=== Cleanup ==="
aws_noauth s3 rm "s3://${BUCKET}" --recursive
aws_noauth s3 rb "s3://${BUCKET}"

echo ""
echo "=== S3 Select Tests Complete ==="
