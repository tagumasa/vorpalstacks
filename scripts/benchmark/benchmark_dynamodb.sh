#!/bin/bash
# DynamoDB Performance Benchmark for VorpalStacks
# Usage: bash scripts/benchmark/benchmark_dynamodb.sh

source "$(dirname "$0")/../services/lib/config.sh"
source "$(dirname "$0")/../services/lib/colors.sh"
source "$(dirname "$0")/../services/lib/aws_cli.sh"

print_header "DynamoDB Performance Benchmark"
echo "Endpoint: $ENDPOINT_URL"
echo "Region: $REGION"
echo ""

TABLE_NAME="BenchmarkTable-$$"
ITERATIONS=${ITERATIONS:-1000}
BATCH_SIZE=${BATCH_SIZE:-25}
CONCURRENT=${CONCURRENT:-10}

cleanup() {
    aws_noauth dynamodb delete-table --table-name "$TABLE_NAME" 2>/dev/null || true
}
trap cleanup EXIT

setup_table() {
    echo -e "${CYAN}Setting up benchmark table...${NC}"
    aws_noauth dynamodb create-table \
        --table-name "$TABLE_NAME" \
        --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S \
        --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE \
        --billing-mode PAY_PER_REQUEST \
        > /dev/null 2>&1
    
    sleep 1
    echo -e "${GREEN}Table created${NC}"
}

benchmark_put_item() {
    echo -e "\n${CYAN}=== PutItem (Single Write) ===${NC}"
    echo "Iterations: $ITERATIONS"
    
    local start=$(date +%s%N)
    for i in $(seq 1 $ITERATIONS); do
        aws_noauth dynamodb put-item \
            --table-name "$TABLE_NAME" \
            --item "{\"pk\":{\"S\":\"bench\"},\"sk\":{\"S\":\"item-$i\"},\"data\":{\"S\":\"data-$i\"}}" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local ops=$(( ITERATIONS * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / ITERATIONS ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Throughput: ${ops} ops/sec"
    echo "Avg latency: ${avg_latency}ms"
}

benchmark_get_item() {
    echo -e "\n${CYAN}=== GetItem (Single Read) ===${NC}"
    echo "Iterations: $ITERATIONS"
    
    local start=$(date +%s%N)
    for i in $(seq 1 $ITERATIONS); do
        aws_noauth dynamodb get-item \
            --table-name "$TABLE_NAME" \
            --key "{\"pk\":{\"S\":\"bench\"},\"sk\":{\"S\":\"item-$i\"}}" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local ops=$(( ITERATIONS * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / ITERATIONS ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Throughput: ${ops} ops/sec"
    echo "Avg latency: ${avg_latency}ms"
}

benchmark_batch_write() {
    echo -e "\n${CYAN}=== BatchWriteItem (Batch Write, size=$BATCH_SIZE) ===${NC}"
    local batches=$(( ITERATIONS / BATCH_SIZE ))
    echo "Batches: $batches"
    
    local start=$(date +%s%N)
    for batch in $(seq 1 $batches); do
        local items=""
        for i in $(seq 1 $BATCH_SIZE); do
            local idx=$(( (batch - 1) * BATCH_SIZE + i + ITERATIONS ))
            items="${items}{\"PutRequest\":{\"Item\":{\"pk\":{\"S\":\"batch\"},\"sk\":{\"S\":\"item-$idx\"},\"data\":{\"S\":\"data-$idx\"}}}}"
            if [ $i -lt $BATCH_SIZE ]; then
                items="${items},"
            fi
        done
        aws_noauth dynamodb batch-write-item \
            --request-items "{\"$TABLE_NAME\":[$items]}" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local total_items=$(( batches * BATCH_SIZE ))
    local ops=$(( total_items * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / batches ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Items written: $total_items"
    echo "Throughput: ${ops} items/sec"
    echo "Avg batch latency: ${avg_latency}ms"
}

benchmark_query() {
    echo -e "\n${CYAN}=== Query (Read with filter) ===${NC}"
    echo "Iterations: $ITERATIONS"
    
    local start=$(date +%s%N)
    for i in $(seq 1 $ITERATIONS); do
        aws_noauth dynamodb query \
            --table-name "$TABLE_NAME" \
            --key-condition-expression "pk = :pk" \
            --expression-attribute-values "{\":pk\":{\"S\":\"bench\"}}" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local ops=$(( ITERATIONS * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / ITERATIONS ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Throughput: ${ops} queries/sec"
    echo "Avg latency: ${avg_latency}ms"
}

benchmark_scan() {
    echo -e "\n${CYAN}=== Scan (Full table scan) ===${NC}"
    echo "Iterations: 10"
    
    local start=$(date +%s%N)
    for i in $(seq 1 10); do
        aws_noauth dynamodb scan \
            --table-name "$TABLE_NAME" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local ops=$(( 10 * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / 10 ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Throughput: ${ops} scans/sec"
    echo "Avg latency: ${avg_latency}ms"
}

benchmark_update_item() {
    echo -e "\n${CYAN}=== UpdateItem (Write with expression) ===${NC}"
    echo "Iterations: $ITERATIONS"
    
    local start=$(date +%s%N)
    for i in $(seq 1 $ITERATIONS); do
        aws_noauth dynamodb update-item \
            --table-name "$TABLE_NAME" \
            --key "{\"pk\":{\"S\":\"bench\"},\"sk\":{\"S\":\"item-$i\"}}" \
            --update-expression "SET #d = :val" \
            --expression-attribute-names '{"#d":"data"}' \
            --expression-attribute-values "{\":val\":{\"S\":\"updated-$i\"}}" \
            > /dev/null 2>&1
    done
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local ops=$(( ITERATIONS * 1000 / duration_ms ))
    local avg_latency=$(( duration_ms / ITERATIONS ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Throughput: ${ops} ops/sec"
    echo "Avg latency: ${avg_latency}ms"
}

benchmark_concurrent_reads() {
    echo -e "\n${CYAN}=== Concurrent Reads ($CONCURRENT parallel) ===${NC}"
    
    local start=$(date +%s%N)
    for i in $(seq 1 $CONCURRENT); do
        (
            for j in $(seq 1 $(( ITERATIONS / CONCURRENT ))); do
                aws_noauth dynamodb get-item \
                    --table-name "$TABLE_NAME" \
                    --key "{\"pk\":{\"S\":\"bench\"},\"sk\":{\"S\":\"item-$j\"}}" \
                    > /dev/null 2>&1
            done
        ) &
    done
    wait
    local end=$(date +%s%N)
    
    local duration_ms=$(( (end - start) / 1000000 ))
    local total_ops=$ITERATIONS
    local ops=$(( total_ops * 1000 / duration_ms ))
    
    echo "Total time: ${duration_ms}ms"
    echo "Total ops: $total_ops"
    echo "Throughput: ${ops} ops/sec"
}

main() {
    setup_table
    
    echo -e "\n${YELLOW}Starting benchmarks (ITERATIONS=$ITERATIONS)...${NC}"
    
    benchmark_put_item
    benchmark_get_item
    benchmark_update_item
    benchmark_query
    benchmark_scan
    benchmark_batch_write
    benchmark_concurrent_reads
    
    echo -e "\n${GREEN}=== Benchmark Complete ===${NC}"
}

main
