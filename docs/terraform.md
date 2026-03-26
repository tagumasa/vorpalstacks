# Terraform & OpenTofu Guide

Vorpalstacks works with HashiCorp Terraform and OpenTofu using the standard AWS provider.

## Prerequisites

- Vorpalstacks server running
- Terraform >= 1.6.0 or OpenTofu >= 1.6.0
- AWS provider ~> 6.0

## Quick Start

### 1. Start Vorpalstacks

```bash
SIGNATURE_VERIFICATION_ENABLED=false PORT=8080 DATA_PATH=./data TEST_MODE=true ./vorpalstacks &
```

### 2. Provider Configuration

```hcl
terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    # Regional services
    dynamodb     = "http://localhost:8080"
    s3           = "http://localhost:8080"
    sqs          = "http://localhost:8080"
    sns          = "http://localhost:8080"
    lambda       = "http://localhost:8080"
    kinesis      = "http://localhost:8080"
    kms          = "http://localhost:8080"
    secretsmanager = "http://localhost:8080"
    ssm          = "http://localhost:8080"
    events       = "http://localhost:8080"
    scheduler    = "http://localhost:8080"
    athena       = "http://localhost:8080"
    cloudwatch   = "http://localhost:8080"
    cloudwatchlogs = "http://localhost:8080"
    cloudtrail   = "http://localhost:8080"
    sesv2        = "http://localhost:8080"
    apigateway   = "http://localhost:8080"
    waf          = "http://localhost:8080"
    wafv2        = "http://localhost:8080"
    cognitoidp   = "http://localhost:8080"
    cognitoidentity = "http://localhost:8080"
    timestreamwrite = "http://localhost:8080"
    timestreamquery = "http://localhost:8080"
    stepfunctions = "http://localhost:8080"
    acm          = "http://localhost:8080"

    # Global services
    iam          = "http://localhost:8080"
    sts          = "http://localhost:8080"
    route53      = "http://localhost:8080"
    cloudfront   = "http://localhost:8080"
  }

  default_tags {
    tags = {
      ManagedBy = "vorpalstacks"
    }
  }
}
```

### 3. Example: SQS Queue

```hcl
resource "aws_sqs_queue" "test" {
  name = "test-queue"
}

output "queue_url" {
  value = aws_sqs_queue.test.id
}
```

### 4. Run

```bash
terraform init
terraform apply
terraform plan  # Should show "No changes"
terraform destroy
```

## Using OpenTofu

OpenTofu is a drop-in replacement for Terraform. Replace all `terraform` commands with `tofu`:

```bash
tofu init
tofu apply
tofu plan
tofu destroy
```

No changes to provider configuration are needed — OpenTofu uses the same `hashicorp/aws` provider.

## Verified Services

The following services have been verified to work with the full Terraform cycle (apply → plan → apply → destroy):

| Service | Resource(s) |
|---------|-------------|
| SQS | `aws_sqs_queue`, `aws_sqs_queue_policy` |
| DynamoDB | `aws_dynamodb_table` |
| S3 | `aws_s3_bucket` |
| SNS | `aws_sns_topic`, `aws_sns_topic_subscription` |
| Lambda | `aws_lambda_function` |
| IAM | `aws_iam_user`, `aws_iam_policy`, `aws_iam_role` |
| KMS | `aws_kms_key` |
| Kinesis | `aws_kinesis_stream` |
| EventBridge | `aws_cloudwatch_event_rule`, `aws_cloudwatch_event_target` |
| Step Functions | `aws_sfn_state_machine` |
| Athena | `aws_athena_workgroup`, `aws_athena_database` |
| CloudWatch Logs | `aws_cloudwatch_log_group`, `aws_cloudwatch_log_metric_filter` |
| Cognito | `aws_cognito_user_pool`, `aws_cognito_user_pool_client` |
| STS | Data source only |
| Timestream Write | `aws_timestreamwrite_database`, `aws_timestreamwrite_table` |
| WAF | `aws_waf_ipset`, `aws_waf_rule`, `aws_waf_web_acl`, `aws_waf_regex_pattern_set` |
| WAFv2 | `aws_wafv2_ip_set`, `aws_wafv2_regex_pattern_set`, `aws_wafv2_rule_group`, `aws_wafv2_web_acl` |

## Tips

- Always run `terraform plan` after `terraform apply` to verify no drift
- Use `TEST_MODE=true` when starting the server for consistent test behaviour
- Set `AWS_ACCOUNT_ID=123456789012` if ARN-based resources need a specific account ID
- Clean state between test runs: delete `terraform.tfstate*` and clear `data/` directory
