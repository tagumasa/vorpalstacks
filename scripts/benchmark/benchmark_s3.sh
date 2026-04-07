#!/bin/bash
# S3 Performance Benchmark for VorpalStacks

BUCKET_NAME="bench-bucket-$$"
ITERATIONS=${ITERATIONS:-100}

cleanup() {
    aws_noauth s3 rb s3://$BUCKET_NAME --force 2>/dev/null || true
}
trap cleanup EXIT

aws_noauth s3 mb s3://$BUCKET_NAME > /dev/null 2>&1
sleep 1

echo "=== PutObject ($ITERATIONS objects) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    echo "data-$i" | aws_noauth s3 cp - s3://$BUCKET_NAME/object-$i.txt > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( ITERATIONS * 1000 / duration_ms ))
avg_latency=$(( duration_ms / ITERATIONS ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} objects/sec"
echo "Avg latency: ${avg_latency}ms"

echo ""
echo "=== GetObject ($ITERATIONS reads) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    aws_noauth s3 cp s3://$BUCKET_NAME/object-$i.txt /dev/null > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( ITERATIONS * 1000 / duration_ms ))
avg_latency=$(( duration_ms / ITERATIONS ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} objects/sec"
echo "Avg latency: ${avg_latency}ms"

echo ""
echo "=== ListObjects (10 lists) ==="
start=$(date +%s%N)
for i in $(seq 1 10); do
    aws_noauth s3 ls s3://$BUCKET_NAME/ > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( 10 * 1000 / duration_ms ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} lists/sec"

echo ""
echo "=== Benchmark Complete ==="
