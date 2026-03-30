# Vorpalstacks

[English](README.md) | [中文](README.zh.md)

> **注意：これはベータ版です。** Vorpalstacksは現在アクティブに開発中です。29のAWSサービスを実装し、595件のGo SDKテスト（他にPython 631件、TypeScript 629件、C# 606件）が通過していますが、すべてのエッジケースとAWSの動作が完全にカバーされているわけではありません。破壊的変更が発生する可能性があります。バグ報告とコントリビューションを歓迎します。

AWS互換サービスを提供する軽量エッジ・オンプレミスクラウドプラットフォーム。

## 概要

Vorpalstacksは、完全なAWS接続が利用できない環境でAWS互換サービスを実行可能にすることを目標にしてます：

- エッジコンピューティングシナリオ
- オンプレミスデプロイ
- 開発・テスト環境
- エアギャップネットワーク

## 特徴

> **本プロジェクトが提供するもの**: モックフレームワークではなく、AWS互換APIの本格的な実装です。各サービスはPebbleDBにデータを保存し、マルチリージョン分離をサポートします。

> **本プロジェクトが提供しないもの**: すべてのAWS動作の完全な再現ではありません。一部のエッジケース、公式ドキュメントにない動作、高度な機能はAWSと異なる場合があります。各サービスの対応範囲については[docs/services.md](docs/services.md)を参照してください。

- **AWS API互換**: 既存のAWS SDKおよびCLIで動作
- **29のAWSサービス**: S3、SQS、SNS、Lambda、DynamoDB、API Gateway、Step Functions、WAF、WAFv2、Kinesis、KMSなど
- **IAM認可**: ユーザー/グループ/ロールベースのアクセス制御による完全なIAMポリシー評価
- **DynamoDB PartiQL**: WHERE関数（attribute_exists、begins_with、contains、size）を含むSQLライククエリ
- **S3 SelectObjectContent**: イベントストリーミングによるCSV/JSONオブジェクトへのSQLクエリ
- **マルチリージョンサポート**: リージョンごとの専用PebbleDBによる分離ストレージ
- **グローバルサービス**: IAM、Route53、CloudFront、STSの共有グローバルストレージ
- **gRPC-Web管理API**: 全サービス向けの別ポートConnect-RPC管理インターフェース
- **軽量**: シングルバイナリ、最小限の依存関係
- **永続ストレージ**: Pebbleベースのキーバリューストア
- **Docker統合**: Lambda関数をコンテナで実行
- **サービス間連携**: イベント駆動型サービス間通信
- **TLSサポート**: 自動生成またはカスタム証明書によるオプションHTTPS

## 実装済みサービス

| サービス | カバレッジ | 備考 |
|---------|-----------|------|
| ACM | Full | |
| API Gateway | Broad | クライアント証明書、ドキュメント、SDK生成なし |
| Athena | Broad | キャパシティ予約、ノートブックセッションなし |
| CloudFront | Selective | 実際のエッジトラフィック分散なし |
| CloudTrail | Broad | イベントデータストア、SQLクエリなし |
| CloudWatch Logs | Selective | Logs Insightsクエリ、エクスポートなし |
| CloudWatch Metrics | Broad | メトリクスストリーム、異常検知なし |
| Cognito IDP | Selective | 外部IdP、ホストUIなし |
| Cognito Identity | Selective | 基本的なアイデンティティプールサポートのみ |
| DynamoDB | Full | |
| EventBridge | Broad | グローバルエンドポイント、パートナーイベントソースなし |
| IAM | Broad | ポリシーシモュレータ、Organizations統合なし |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | Durable Functions、コード署名なし |
| Route53 | Selective | DNSレコード管理のみ |
| S3 | Broad | アナリティクス、インベントリ、S3 Expressなし |
| Scheduler | Full | |
| Secrets Manager | Full | |
| SESv2 | Broad | 配信テスト、マルチテナンシーなし |
| SNS | Broad | SMS送信未対応 |
| SQS | Full | |
| SSM | Selective | Parameter Storeのみ |
| STS | Full | |
| Step Functions | Full | |
| Timestream | Full | |
| WAF | Selective | マネージドルールグループ、ログ設定なし |
| WAFv2 | Broad | |

詳細なカバレッジ階層とサービス統合パターンについては[docs/services.md](docs/services.md)を参照してください。

## クイックスタート

### ビルド

```bash
go build -o vorpalstacks .
```

### 実行（開発モード）

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### Dockerで実行（Lambda用）

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### AWS CLIでの使用

```bash
export AWS_ENDPOINT_URL=http://localhost:8080

aws --endpoint-url=http://localhost:8080 --no-sign-request sns list-topics
aws --endpoint-url=http://localhost:8080 --no-sign-request sqs list-queues
aws --endpoint-url=http://localhost:8080 --no-sign-request lambda list-functions
```

## テスト

### ユニットテスト

```bash
make test
```

### SDKテスト（AWS Go SDK v2）

```bash
SIGNATURE_VERIFICATION_ENABLED=false TEST_MODE=true DATA_PATH=./data-test ./vorpalstacks &

cd sdk-tests && ./sdk-tests-test -service all -v
```

### CLI統合テスト

```bash
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
cd scripts/services && bash test_s3.sh
cd scripts/services && bash test_iam.sh
```

## 設定

### 環境変数

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `PORT` | `8080` | HTTPサーバーポート |
| `DATA_PATH` | `./data` | 永続データストレージのパス |
| `AWS_REGION` | `us-east-1` | デフォルトリージョン |
| `AWS_ACCOUNT_ID` | `000000000000` | AWSアカウントID |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | AWS Signature V4検証を有効化 |
| `GRPC_WEB_PORT` | `9090` | gRPC-Web管理サーバーポート |
| `TLS_ENABLED` | `false` | TLSを有効化 |
| `AUTHORIZATION_ENABLED` | `false` | IAMポリシー評価を有効化 |

完全な一覧については[設定](docs/configuration.md)を参照してください。

### 本番構成例

```bash
export AWS_ACCESS_KEY_ID=your-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-key
export SIGNATURE_VERIFICATION_ENABLED=true
export AUTHORIZATION_ENABLED=true
export DATA_PATH=/var/lib/vorpalstacks
./vorpalstacks
```

## データストレージ構造

```
DATA_PATH/
├── us-east-1/               # リージョン固有のストレージ
│   ├── pebble/              # us-east-1用PebbleDB
│   ├── code/                # Lambda関数コード
│   └── logs/                # CloudWatch Logsチャンク
├── global/                  # グローバルサービス（IAM、Route53、CloudFront、STS）
│   └── pebble/
└── uploads/                 # S3マルチパートアップロード（一時的）
```

各リージョンにはPebbleDB、Lambdaコード、ログを含む分離されたストレージがあります。グローバルサービスは専用ストレージを共有します。

## Docker要件

Lambda機能を使用する場合：

1. Dockerがインストールされ実行中であること
2. Lambdaランタイムイメージが自動的にプルされます
3. Dockerソケットが`DOCKER_HOST`（デフォルト: `unix:///var/run/docker.sock`）でアクセス可能であること

## ドキュメント

- [アーキテクチャ](docs/architecture.md) - システムアーキテクチャ概要
- [サービス](docs/services.md) - 実装済みAWSサービス
- [設定](docs/configuration.md) - 環境変数と設定
- [統合](docs/integration.md) - サービス間通信

## 既知の制限事項

- すべてのサービスのすべてのAWSオペレーションが実装されているわけではありません — 詳細は[docs/services.md](docs/services.md)を参照
- 一部のエッジケースと未文書化のAWS動作は異なる場合があります
- クロスアカウントアクセス制御なし（シングルアカウントモード）
- CloudFrontディストリビューションは実際のエッジトラフィックを配信しません
- CognitoホストUIドメインは未対応（CloudFrontエッジが必要）
- SQS FIFOキューのサポートは限定されています
- DynamoDB StreamsとGlobal Tablesは未実装です

## ロードマップ

- **短期**: バグ修正、リファクタリング、安定性向上
- **サービス拡張**: AWS IoT Core
- **Terraform**: 既存サービスのプロバイダー互換性改善

リリース履歴については[CHANGELOG.md](CHANGELOG.md)を参照してください。

## 動作要件

- Go 1.25+
- Docker（Lambda機能を使用する場合）

## ライセンス

本プロジェクトは [Functional Source License, Version 1.1, MIT Future License (FSL-1.1-MIT)](LICENSE) の下でライセンスされています。

> **注意**: ルートライセンスは本プロジェクトがプロダクション安定性に達した後、MITに変更します。現在は暫定的にFSL-1.1-MITにしてますが、ある程度の十分な動作確認が取れたタイミングでライセンスを変更する予定です。

`pkg/`ディレクトリにはApache License 2.0でライセンスされたコードが含まれています — 詳細は`pkg/sqlparser/LICENSE.md`および`pkg/vsjwt/LICENSE`を参照してください。
