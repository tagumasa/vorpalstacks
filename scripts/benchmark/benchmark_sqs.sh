#!/bin/bash
# SQS Performance Benchmark for VorpalStacks

QUEUE_NAME="bench-queue-$$"
ITERATIONS=${ITERATIONS:-100}

cleanup() {
    aws_noauth sqs delete-queue --queue-url "http://localhost:8080/123456789012/$QUEUE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

aws_noauth sqs create-queue --queue-name "$QUEUE_NAME" > /dev/null 2>&1
sleep 1

echo "=== SendMessage ($ITERATIONS messages) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    aws_noauth sqs send-message \
        --queue-url "http://localhost:8080/123456789012/$QUEUE_NAME" \
        --message-body "message-$i" \
        > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( ITERATIONS * 1000 / duration_ms ))
avg_latency=$(( duration_ms / ITERATIONS ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} messages/sec"
echo "Avg latency: ${avg_latency}ms"

echo ""
echo "=== ReceiveMessage ($ITERATIONS receives) ==="
start=$(date +%s%N)
for i in $(seq 1 $ITERATIONS); do
    aws_noauth sqs receive-message \
        --queue-url "http://localhost:8080/123456789012/$QUEUE_NAME" \
        > /dev/null 2>&1
done
end=$(date +%s%N)

duration_ms=$(( (end - start) / 1000000 ))
ops=$(( ITERATIONS * 1000 / duration_ms ))

echo "Total time: ${duration_ms}ms"
echo "Throughput: ${ops} receives/sec"

echo ""
echo "=== Benchmark Complete ==="
