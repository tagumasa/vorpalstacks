#!/bin/bash
# Kinesis Performance Benchmark for VorpalStacks

STREAM_NAME="bench_stream_$$"
ITERATIONS=${ITERATIONS:-100}

cleanup() {
    aws_noauth kinesis delete-stream --stream-name "$STREAM_NAME" 2>/dev/null || true
}
trap cleanup EXIT

aws_noauth kinesis create-stream --stream-name "$STREAM_NAME" --shard-count 1 > /dev/null 2>&1
sleep 2

echo "=== PutRecord ($ITERATIONS records) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    aws_noauth kinesis put-record \
        --stream-name "$STREAM_NAME" \
        --data "data-$i" \
        --partition-key "key-$i" \
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
echo "=== GetRecords (10 reads) ==="
ITERATOR=$(aws_noauth kinesis get-shard-iterator --stream-name "$STREAM_NAME" --shard-id "shardId-000000000000" --shard-iterator-type TRIM_HORIZON --query 'ShardIterator' --output text 2>/dev/null)

start=$(date +%s%N)
for i in $(seq 1 10); do
    aws_noauth kinesis get-records --shard-iterator "$ITERATOR" > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( 10 * 1000 / duration_ms ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} reads/sec"

echo ""
echo "=== Benchmark Complete ==="
