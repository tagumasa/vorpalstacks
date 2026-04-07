#!/bin/bash
# Timestream Performance Benchmark for VorpalStacks

DATABASE_NAME="bench_db_$$"
TABLE_NAME="bench_table_$$"
ITERATIONS=${ITERATIONS:-1000}

cleanup() {
    aws_noauth timestream-write delete-table --database-name "$DATABASE_NAME" --table-name "$TABLE_NAME" 2>/dev/null || true
    aws_noauth timestream-write delete-database --database-name "$DATABASE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

aws_noauth timestream-write create-database --database-name "$DATABASE_NAME" > /dev/null 2>&1
sleep 1
aws_noauth timestream-write create-table --database-name "$DATABASE_NAME" --table-name "$TABLE_NAME" --retention-properties 'MemoryStoreRetentionPeriodInHours=1,MagneticStoreRetentionPeriodInDays=1' > /dev/null 2>&1
sleep 1

echo "=== WriteRecords ($ITERATIONS records) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    aws_noauth timestream-write write-records \
        --database-name "$DATABASE_NAME" \
        --table-name "$TABLE_NAME" \
        --records "[{\"MeasureName\":\"cpu\",\"MeasureValue\":\"$i\",\"Time\":\"$(date +%s%s)000\",\"Dimensions\":[{\"Name\":\"host\",\"Value\":\"server$i\"}]}]" \
        > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( ITERATIONS * 1000 / duration_ms ))
avg_latency=$(( duration_ms / ITERATIONS ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} records/sec"
echo "Avg latency: ${avg_latency}ms"

echo ""
echo "=== Query (10 queries) ==="
start=$(date +%s%N)
for i in $(seq 1 10); do
    aws_noauth timestream-query query \
        --query "SELECT * FROM \"$DATABASE_NAME\".\"$TABLE_NAME\" LIMIT 10" \
        > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( 10 * 1000 / duration_ms ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} queries/sec"

echo ""
echo "=== Benchmark Complete ==="
