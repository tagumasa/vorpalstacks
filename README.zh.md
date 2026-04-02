# Vorpalstacks

[English](README.md) | [日本語](README.ja.md)

> **注意：目前为测试版。** Vorpalstacks 正在积极开发中。已实现 30 个 AWS 服务，Go SDK 测试 1544 项通过（另有 Python 631 项、TypeScript 629 项、C# 606 项），但并未完全覆盖所有边界情况和 AWS 行为。可能存在破坏性变更。欢迎提交 Bug 报告和贡献代码。

提供 AWS 兼容服务的轻量边缘及本地云平台。

## 概述

Vorpalstacks 的目标是在无法使用完整 AWS 连接的环境中运行 AWS 兼容服务：

- 边缘计算场景
- 本地部署
- 开发和测试环境
- 物理隔离网络

## 特性

> **关于本项目的定位**：这是 AWS 兼容 API 的真实实现，而非 Mock 框架。各服务将数据存储在 PebbleDB 中，并支持多区域隔离。

> **关于本项目的局限**：并非对所有 AWS 行为的完整复刻。某些边界情况、未公开文档的行为和高级功能可能与 AWS 存在差异。各服务的具体支持范围请参阅 [docs/services.md](docs/services.md)。

- **AWS API 兼容**：可直接使用现有 AWS SDK 和 CLI
- **30 个 AWS 服务**：S3、SQS、SNS、Lambda、DynamoDB、API Gateway、Step Functions、WAF、WAFv2、Kinesis、KMS、Neptune 等
- **IAM 授权**：完整的 IAM 策略评估，支持基于用户/组/角色的访问控制
- **DynamoDB PartiQL**：支持 WHERE 函数（attribute_exists、begins_with、contains、size）的类 SQL 查询
- **S3 SelectObjectContent**：通过事件流对 CSV/JSON 对象执行 SQL 查询
- **多区域支持**：每个区域使用独立的 PebbleDB 实现区域隔离存储
- **全局服务**：IAM、Route53、CloudFront、STS 共享全局存储
- **gRPC-Web 管理 API**：在独立端口上为所有服务提供 Connect-RPC 管理界面
- **轻量**：单一二进制文件，依赖极少
- **持久存储**：基于 Pebble 的键值存储
- **Docker 集成**：Lambda 函数在容器中运行
- **服务间联动**：基于事件驱动的服务间通信
- **TLS 支持**：可选 HTTPS，支持自动生成或自定义证书
- **与 LocalStack 的对比**：技术比较请参阅 [docs/localstack_vs_vorpalstacks_report.md](docs/localstack_vs_vorpalstacks_report.md)

## 已实现服务

| 服务 | 覆盖范围 | 备注 |
|------|---------|------|
| ACM | 完整 | |
| API Gateway | 较全面 | 不支持客户端证书、文档或 SDK 生成 |
| Athena | 较全面 | 不支持容量预留或 Notebook 会话 |
| CloudFront | 有限 | 不提供实际边缘流量分发 |
| CloudTrail | 较全面 | 不支持事件数据存储或 SQL 查询 |
| CloudWatch Logs | 有限 | 不支持 Logs Insights 查询或导出 |
| CloudWatch Metrics | 较全面 | 不支持指标流或异常检测 |
| Cognito IDP | 有限 | 不支持外部 IdP、托管 UI |
| Cognito Identity | 有限 | 仅支持基本身份池 |
| DynamoDB | 完整 | |
| EventBridge | 较全面 | 不支持全局终端节点或合作伙伴事件源 |
| IAM | 较全面 | 不支持策略模拟器或 Organizations 集成 |
| Kinesis | 完整 | |
| KMS | 完整 | |
| Lambda | 较全面 | 不支持 Durable Functions 或代码签名 |
| Neptune | 完整 | Property graph + RDF、openCypher/Gremlin、批量加载器 |
| Route53 | 有限 | 仅 DNS 记录管理 |
| S3 | 较全面 | 不支持分析、清单或 S3 Express |
| Scheduler | 完整 | |
| Secrets Manager | 完整 | |
| SESv2 | 较全面 | 不支持送达测试或多租户 |
| SNS | 较全面 | 不支持 SMS 发送 |
| SQS | 完整 | |
| SSM | 有限 | 仅 Parameter Store |
| STS | 完整 | |
| Step Functions | 完整 | |
| Timestream | 完整 | |
| WAF | 有限 | 不支持托管规则组或日志配置 |
| WAFv2 | 较全面 | |

覆盖范围的详细分级和服务联动方式请参阅 [docs/services.md](docs/services.md)。

## 快速开始

### 构建

```bash
go build -o vorpalstacks .
```

### 运行（开发模式）

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### 使用 Docker 运行（Lambda 用）

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### 使用 AWS CLI

```bash
export AWS_ENDPOINT_URL=http://localhost:8080

aws --endpoint-url=http://localhost:8080 --no-sign-request sns list-topics
aws --endpoint-url=http://localhost:8080 --no-sign-request sqs list-queues
aws --endpoint-url=http://localhost:8080 --no-sign-request lambda list-functions
```

## 测试

### 单元测试

```bash
make test
```

### SDK 测试（AWS Go SDK v2）

```bash
SIGNATURE_VERIFICATION_ENABLED=false TEST_MODE=true DATA_PATH=./data-test ./vorpalstacks &

cd sdk-tests && ./sdk-tests-test -service all -v
```

### CLI 集成测试

```bash
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
cd scripts/services && bash test_s3.sh
cd scripts/services && bash test_iam.sh
```

## 配置

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `8080` | HTTP 服务器端口 |
| `DATA_PATH` | `./data` | 持久化数据存储路径 |
| `AWS_REGION` | `us-east-1` | 默认区域 |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS 账号 ID |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | 启用 AWS Signature V4 验证 |
| `GRPC_WEB_PORT` | `9090` | gRPC-Web 管理服务器端口 |
| `TLS_ENABLED` | `false` | 启用 TLS |
| `AUTHORIZATION_ENABLED` | `false` | 启用 IAM 策略评估 |

完整列表请参阅 [配置](docs/configuration.md)。

### 生产环境配置示例

```bash
export AWS_ACCESS_KEY_ID=your-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-key
export SIGNATURE_VERIFICATION_ENABLED=true
export AUTHORIZATION_ENABLED=true
export DATA_PATH=/var/lib/vorpalstacks
./vorpalstacks
```

## 数据存储结构

```
DATA_PATH/
├── us-east-1/               # 区域专用存储
│   ├── pebble/              # us-east-1 的 PebbleDB
│   ├── code/                # Lambda 函数代码
│   └── logs/                # CloudWatch Logs 数据块
├── global/                  # 全局服务（IAM、Route53、CloudFront、STS）
│   └── pebble/
└── uploads/                 # S3 分段上传（临时文件）
```

每个区域拥有独立的 PebbleDB、Lambda 代码和日志存储。全局服务共享专用存储。

## Docker 要求

使用 Lambda 功能时需满足以下条件：

1. 已安装 Docker 且正在运行
2. Lambda 运行时镜像会自动拉取
3. Docker Socket 需在 `DOCKER_HOST`（默认: `unix:///var/run/docker.sock`）处可访问

## 文档

- [架构](docs/architecture.md) - 系统架构概览
- [服务](docs/services.md) - 已实现的 AWS 服务
- [配置](docs/configuration.md) - 环境变量和设置
- [联动](docs/integration.md) - 服务间通信

## 已知限制

- 并非所有服务的全部 AWS 操作都已实现 — 详情请参阅 [docs/services.md](docs/services.md)
- 某些边界情况和未公开文档的 AWS 行为可能存在差异
- 不支持跨账号访问控制（单账号模式）
- CloudFront 分发不提供实际边缘流量
- 不支持 Cognito 托管 UI 域名（需要 CloudFront 边缘）
- SQS FIFO 队列的支持有限
- DynamoDB Streams 和 Global Tables 尚未实现

## 路线图

- **短期**：Bug 修复、重构和稳定性改进
- **服务扩展**：AWS IoT Core, AppSync
- **Terraform**：28 个服务已通过基本适用性测试 — 详情及运行方式请参阅 [vorpalstacks-conformance-tests](https://github.com/tagumasa/vorpalstacks-conformance-tests)

发布历史请参阅 [CHANGELOG.md](CHANGELOG.md)。

## 系统要求

- Go 1.25+
- Docker（使用 Lambda 功能时）

## 性能

Vorpalstacks 将全部 29 个服务实现为基于 PebbleDB 的原生 Go 二进制文件，避免了解释型语言和外部进程依赖的开销。

这种架构使核心操作达到亚毫秒级延迟，可以在 CI/CD 流水线中直接运行大量 API 测试（Go SDK 1544 项、Python 631 项、TypeScript 629 项、C# 606 项），无需容器化开销。

### 基准测试结果（参考值）

平台: AMD Ryzen 7 5700U (16 核), Linux, Go 1.25.8, Pebble v2.1.4

> **注意**：以下数据受环境影响。在硬件、配置和工作负载完全相同的情况下，与其他系统的直接比较才有意义。

| 服务 | 操作 | 平均延迟 | ops/sec |
|------|------|---------|---------|
| DynamoDB | GetItem | 0.38ms | ~2,630 |
| S3 | GetObject (1KB) | 0.27ms | ~3,700 |
| SQS | SendMessage | 0.67ms | ~1,490 |

## 许可证

本项目基于 [Functional Source License, Version 1.1, MIT Future License (FSL-1.1-MIT)](LICENSE) 许可。

> **注意**：根许可证将在项目达到生产稳定性后更改为 MIT。目前暂定为 FSL-1.1-MIT，待充分验证运行稳定性后再行更改。

`pkg/` 目录下包含基于 Apache License 2.0 许可的代码 — 详情请参阅 `pkg/sqlparser/LICENSE.md` 和 `pkg/vsjwt/LICENSE`。
