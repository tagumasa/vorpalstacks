#!/usr/bin/env python3
"""S3 SelectObjectContent test using boto3"""

import boto3
import io

endpoint_url = "http://localhost:8080"
bucket_name = "test-select-py"
key_name = "test.csv"

s3 = boto3.client(
    "s3",
    endpoint_url=endpoint_url,
    aws_access_key_id="test",
    aws_secret_access_key="test",
    region_name="us-east-1",
)

# Create bucket
try:
    s3.create_bucket(Bucket=bucket_name)
    print(f"Created bucket: {bucket_name}")
except Exception as e:
    print(f"Bucket exists or error: {e}")

# Upload CSV
csv_data = b"name,age,city\nAlice,30,Tokyo\nBob,25,Osaka"
s3.put_object(Bucket=bucket_name, Key=key_name, Body=csv_data)
print(f"Uploaded {key_name}")

# Select content
print("\nTest 1: SELECT * FROM s3object")
try:
    response = s3.select_object_content(
        Bucket=bucket_name,
        Key=key_name,
        ExpressionType="SQL",
        Expression="SELECT * FROM s3object",
        InputSerialization={"CSV": {"FileHeaderInfo": "USE"}},
        OutputSerialization={"CSV": {}},
    )

    result = b""
    for event in response["Payload"]:
        if "Records" in event:
            result += event["Records"]["Payload"]
        elif "Stats" in event:
            print(f"Stats: {event['Stats']}")
        elif "End" in event:
            print("End event received")

    print(f"Result:\n{result.decode()}")
except Exception as e:
    print(f"Error: {e}")

# Cleanup
s3.delete_object(Bucket=bucket_name, Key=key_name)
s3.delete_bucket(Bucket=bucket_name)
print("\nCleanup complete")
