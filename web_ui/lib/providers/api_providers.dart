import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc_web.dart';

import '../src/generated/acm.pbgrpc.dart';
import '../src/generated/admin_auth.pbgrpc.dart';
import '../src/generated/admin_config.pbgrpc.dart';
import '../src/generated/apigateway.pbgrpc.dart';
import '../src/generated/athena.pbgrpc.dart';
import '../src/generated/cloudfront.pbgrpc.dart';
import '../src/generated/cloudtrail.pbgrpc.dart';
import '../src/generated/cognito_identity.pbgrpc.dart';
import '../src/generated/cognito_idp.pbgrpc.dart';
import '../src/generated/dynamodb.pbgrpc.dart';
import '../src/generated/email.pbgrpc.dart';
import '../src/generated/events.pbgrpc.dart';
import '../src/generated/iam.pbgrpc.dart';
import '../src/generated/kinesis.pbgrpc.dart';
import '../src/generated/kms.pbgrpc.dart';
import '../src/generated/lambda.pbgrpc.dart';
import '../src/generated/logs.pbgrpc.dart';
import '../src/generated/monitoring.pbgrpc.dart';
import '../src/generated/route53.pbgrpc.dart';
import '../src/generated/s3.pbgrpc.dart';
import '../src/generated/scheduler.pbgrpc.dart';
import '../src/generated/secretsmanager.pbgrpc.dart';
import '../src/generated/sns.pbgrpc.dart';
import '../src/generated/sqs.pbgrpc.dart';
import '../src/generated/ssm.pbgrpc.dart';
import '../src/generated/states.pbgrpc.dart';
import '../src/generated/sts.pbgrpc.dart';
import '../src/generated/waf.pbgrpc.dart';
import '../src/generated/wafv2.pbgrpc.dart';
import '../src/generated/query.timestream.pbgrpc.dart';
import '../src/generated/ingest.timestream.pbgrpc.dart';

const kDefaultBaseUrl = 'http://localhost:9090';

final channelProvider = Provider<GrpcWebClientChannel>((ref) {
  final channel = GrpcWebClientChannel.xhr(Uri.parse(kDefaultBaseUrl));
  ref.onDispose(() => channel.shutdown());
  return channel;
});

final adminAuthProvider = Provider<AdminAuthServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return AdminAuthServiceClient(channel);
});

final adminConfigProvider = Provider<AdminConfigServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return AdminConfigServiceClient(channel);
});

final acmProvider = Provider<ACMServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return ACMServiceClient(channel);
});

final apiGatewayProvider = Provider<APIGatewayServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return APIGatewayServiceClient(channel);
});

final athenaProvider = Provider<AthenaServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return AthenaServiceClient(channel);
});

final cloudFrontProvider = Provider<CloudFrontServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CloudFrontServiceClient(channel);
});

final cloudTrailProvider = Provider<CloudTrailServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CloudTrailServiceClient(channel);
});

final cloudWatchProvider = Provider<CloudWatchServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CloudWatchServiceClient(channel);
});

final cloudWatchLogsProvider = Provider<CloudWatchLogsServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CloudWatchLogsServiceClient(channel);
});

final cloudWatchEventsProvider = Provider<CloudWatchEventsServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CloudWatchEventsServiceClient(channel);
});

final cognitoIdentityProvider = Provider<CognitoIdentityServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return CognitoIdentityServiceClient(channel);
});

final cognitoIdpProvider = Provider<CognitoIdentityProviderServiceClient>((
  ref,
) {
  final channel = ref.watch(channelProvider);
  return CognitoIdentityProviderServiceClient(channel);
});

final dynamoDbProvider = Provider<DynamoDBServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return DynamoDBServiceClient(channel);
});

final iamProvider = Provider<IAMServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return IAMServiceClient(channel);
});

final kinesisProvider = Provider<KinesisServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return KinesisServiceClient(channel);
});

final kmsProvider = Provider<KMSServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return KMSServiceClient(channel);
});

final lambdaProvider = Provider<LambdaServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return LambdaServiceClient(channel);
});

final route53Provider = Provider<Route53ServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return Route53ServiceClient(channel);
});

final s3Provider = Provider<S3ServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return S3ServiceClient(channel);
});

final schedulerProvider = Provider<SchedulerServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SchedulerServiceClient(channel);
});

final secretsManagerProvider = Provider<SecretsManagerServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SecretsManagerServiceClient(channel);
});

final sesProvider = Provider<SESv2ServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SESv2ServiceClient(channel);
});

final snsProvider = Provider<SNSServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SNSServiceClient(channel);
});

final sqsProvider = Provider<SQSServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SQSServiceClient(channel);
});

final ssmProvider = Provider<SSMServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SSMServiceClient(channel);
});

final stepFunctionsProvider = Provider<SFNServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return SFNServiceClient(channel);
});

final stsProvider = Provider<STSServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return STSServiceClient(channel);
});

final wafProvider = Provider<WAFServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return WAFServiceClient(channel);
});

final wafv2Provider = Provider<WAFV2ServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return WAFV2ServiceClient(channel);
});

final timestreamQueryProvider = Provider<TimestreamQueryServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return TimestreamQueryServiceClient(channel);
});

final timestreamWriteProvider = Provider<TimestreamWriteServiceClient>((ref) {
  final channel = ref.watch(channelProvider);
  return TimestreamWriteServiceClient(channel);
});
