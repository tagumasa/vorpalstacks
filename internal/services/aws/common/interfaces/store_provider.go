package interfaces

import (
	acmstore "vorpalstacks/internal/store/aws/acm"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	athenastore "vorpalstacks/internal/store/aws/athena"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	cloudwatchstore "vorpalstacks/internal/store/aws/cloudwatch"
	cloudwatchlogsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	eventbridgestore "vorpalstacks/internal/store/aws/eventbridge"
	iamstore "vorpalstacks/internal/store/aws/iam"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	route53store "vorpalstacks/internal/store/aws/route53"
	s3store "vorpalstacks/internal/store/aws/s3"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
	stsstore "vorpalstacks/internal/store/aws/sts"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// StoreProvider defines the interface for accessing various service stores.
type StoreProvider interface {
	GetACMStore() acmstore.CertificateStoreInterface
	GetAPIGatewayStore() apigatewaystore.APIGatewayStoresInterface
	GetAthenaStore() athenastore.AthenaStoresInterface
	GetCloudFrontStores() cloudfrontstore.CloudFrontStoresInterface
	GetCloudTrailStore() cloudtrailstore.CloudTrailStoreInterface
	GetCloudWatchAlarmStore() cloudwatchstore.AlarmStoreInterface
	GetCloudWatchMetricStore() cloudwatchstore.MetricChunkStoreInterface
	GetCloudWatchLogsStore() cloudwatchlogsstore.CloudWatchLogsStoreInterface
	GetCognitoIdentityStore() cognitoidentitystore.CognitoIdentityStoreInterface
	GetCognitoStore() cognitostore.CognitoStoreInterface
	GetDynamoDBStore() dynamodbstore.DynamoDBStoreInterface
	GetEventBridgeStore() eventbridgestore.EventsStoreInterface
	GetIAMStore() iamstore.IAMStoreInterface
	GetKMSStores() kmsstore.KMSStoresInterface
	GetLambdaStore() lambdastore.LambdaStoreInterface
	GetRoute53Stores() route53store.Route53StoresInterface
	GetS3Store() s3store.S3StoreInterface
	GetSchedulerStore() schedulerstore.SchedulerStoreInterface
	GetSecretsManagerStore() secretsmanagerstore.SecretStoreInterface
	GetSESv2Store() sesv2store.SESv2StoreInterface
	GetSFNStore() sfnstore.StepFunctionStoreInterface
	GetSNSStore() snsstore.SNSStoreInterface
	GetSQSStore() sqsstore.SQSStoreInterface
	GetSSMStore() ssmstore.SSMStoreInterface
	GetSTSStore() stsstore.SessionStoreInterface
	GetTimestreamStores() timestreamstore.TimestreamStoresInterface
	GetWAFStores() wafstore.WAFStoresInterface
}

// StoreProviderImpl provides access to all service stores.
type StoreProviderImpl struct {
	iamStore              iamstore.IAMStoreInterface
	acmStore              acmstore.CertificateStoreInterface
	s3Store               s3store.S3StoreInterface
	ssmStore              ssmstore.SSMStoreInterface
	stsStore              stsstore.SessionStoreInterface
	dynamoDBStore         dynamodbstore.DynamoDBStoreInterface
	lambdaStore           lambdastore.LambdaStoreInterface
	eventBridgeStore      eventbridgestore.EventsStoreInterface
	cloudWatchAlarmStore  cloudwatchstore.AlarmStoreInterface
	cloudWatchMetricStore cloudwatchstore.MetricChunkStoreInterface
	cloudWatchLogsStore   cloudwatchlogsstore.CloudWatchLogsStoreInterface
	apiGatewayStore       apigatewaystore.APIGatewayStoresInterface
	sfnStore              sfnstore.StepFunctionStoreInterface
	athenaStore           athenastore.AthenaStoresInterface
	cognitoStore          cognitostore.CognitoStoreInterface
	timestreamStores      timestreamstore.TimestreamStoresInterface
	wafStores             wafstore.WAFStoresInterface
	kmsStores             kmsstore.KMSStoresInterface
	sesv2Store            sesv2store.SESv2StoreInterface
	schedulerStore        schedulerstore.SchedulerStoreInterface
	route53Stores         route53store.Route53StoresInterface
	cognitoIdentityStore  cognitoidentitystore.CognitoIdentityStoreInterface
	cloudFrontStores      cloudfrontstore.CloudFrontStoresInterface
	cloudTrailStore       cloudtrailstore.CloudTrailStoreInterface
	secretsManagerStore   secretsmanagerstore.SecretStoreInterface
	snsStore              snsstore.SNSStoreInterface
	sqsStore              sqsstore.SQSStoreInterface
}

// NewStoreProvider creates a new StoreProviderImpl with the provided store dependencies.
func NewStoreProvider(
	iamStore iamstore.IAMStoreInterface,
	acmStore acmstore.CertificateStoreInterface,
	s3Store s3store.S3StoreInterface,
	ssmStore ssmstore.SSMStoreInterface,
	stsStore stsstore.SessionStoreInterface,
	dynamoDBStore dynamodbstore.DynamoDBStoreInterface,
	lambdaStore lambdastore.LambdaStoreInterface,
	eventBridgeStore eventbridgestore.EventsStoreInterface,
	cloudWatchAlarmStore cloudwatchstore.AlarmStoreInterface,
	cloudWatchMetricStore cloudwatchstore.MetricChunkStoreInterface,
	cloudWatchLogsStore cloudwatchlogsstore.CloudWatchLogsStoreInterface,
	apiGatewayStore apigatewaystore.APIGatewayStoresInterface,
	sfnStore sfnstore.StepFunctionStoreInterface,
	athenaStore athenastore.AthenaStoresInterface,
	cognitoStore cognitostore.CognitoStoreInterface,
	timestreamStores timestreamstore.TimestreamStoresInterface,
	wafStores wafstore.WAFStoresInterface,
	kmsStores kmsstore.KMSStoresInterface,
	sesv2Store sesv2store.SESv2StoreInterface,
	schedulerStore schedulerstore.SchedulerStoreInterface,
	route53Stores route53store.Route53StoresInterface,
	cognitoIdentityStore cognitoidentitystore.CognitoIdentityStoreInterface,
	cloudFrontStores cloudfrontstore.CloudFrontStoresInterface,
	cloudTrailStore cloudtrailstore.CloudTrailStoreInterface,
	secretsManagerStore secretsmanagerstore.SecretStoreInterface,
	snsStore snsstore.SNSStoreInterface,
	sqsStore sqsstore.SQSStoreInterface,
) *StoreProviderImpl {
	return &StoreProviderImpl{
		iamStore:              iamStore,
		acmStore:              acmStore,
		s3Store:               s3Store,
		ssmStore:              ssmStore,
		stsStore:              stsStore,
		dynamoDBStore:         dynamoDBStore,
		lambdaStore:           lambdaStore,
		eventBridgeStore:      eventBridgeStore,
		cloudWatchAlarmStore:  cloudWatchAlarmStore,
		cloudWatchMetricStore: cloudWatchMetricStore,
		cloudWatchLogsStore:   cloudWatchLogsStore,
		apiGatewayStore:       apiGatewayStore,
		sfnStore:              sfnStore,
		athenaStore:           athenaStore,
		cognitoStore:          cognitoStore,
		timestreamStores:      timestreamStores,
		wafStores:             wafStores,
		kmsStores:             kmsStores,
		sesv2Store:            sesv2Store,
		schedulerStore:        schedulerStore,
		route53Stores:         route53Stores,
		cognitoIdentityStore:  cognitoIdentityStore,
		cloudFrontStores:      cloudFrontStores,
		cloudTrailStore:       cloudTrailStore,
		secretsManagerStore:   secretsManagerStore,
		snsStore:              snsStore,
		sqsStore:              sqsStore,
	}
}

// GetIAMStore returns the IAM store.
func (s *StoreProviderImpl) GetIAMStore() iamstore.IAMStoreInterface { return s.iamStore }

// GetACMStore returns the ACM store.
func (s *StoreProviderImpl) GetACMStore() acmstore.CertificateStoreInterface { return s.acmStore }

// GetS3Store returns the S3 store.
func (s *StoreProviderImpl) GetS3Store() s3store.S3StoreInterface { return s.s3Store }

// GetSSMStore returns the SSM store.
func (s *StoreProviderImpl) GetSSMStore() ssmstore.SSMStoreInterface { return s.ssmStore }

// GetSTSStore returns the STS store.
func (s *StoreProviderImpl) GetSTSStore() stsstore.SessionStoreInterface { return s.stsStore }

// GetDynamoDBStore returns the DynamoDB store.
func (s *StoreProviderImpl) GetDynamoDBStore() dynamodbstore.DynamoDBStoreInterface {
	return s.dynamoDBStore
}

// GetLambdaStore returns the Lambda store.
func (s *StoreProviderImpl) GetLambdaStore() lambdastore.LambdaStoreInterface { return s.lambdaStore }

// GetEventBridgeStore returns the EventBridge store.
func (s *StoreProviderImpl) GetEventBridgeStore() eventbridgestore.EventsStoreInterface {
	return s.eventBridgeStore
}

// GetCloudWatchAlarmStore returns the CloudWatch Alarm store.
func (s *StoreProviderImpl) GetCloudWatchAlarmStore() cloudwatchstore.AlarmStoreInterface {
	return s.cloudWatchAlarmStore
}

// GetCloudWatchMetricStore returns the CloudWatch Metric store.
func (s *StoreProviderImpl) GetCloudWatchMetricStore() cloudwatchstore.MetricChunkStoreInterface {
	return s.cloudWatchMetricStore
}

// GetCloudWatchLogsStore returns the CloudWatch Logs store.
func (s *StoreProviderImpl) GetCloudWatchLogsStore() cloudwatchlogsstore.CloudWatchLogsStoreInterface {
	return s.cloudWatchLogsStore
}

// GetAPIGatewayStore returns the API Gateway store.
func (s *StoreProviderImpl) GetAPIGatewayStore() apigatewaystore.APIGatewayStoresInterface {
	return s.apiGatewayStore
}

// GetSFNStore returns the Step Functions store.
func (s *StoreProviderImpl) GetSFNStore() sfnstore.StepFunctionStoreInterface { return s.sfnStore }

// GetAthenaStore returns the Athena store.
func (s *StoreProviderImpl) GetAthenaStore() athenastore.AthenaStoresInterface { return s.athenaStore }

// GetCognitoStore returns the Cognito User Pool store.
func (s *StoreProviderImpl) GetCognitoStore() cognitostore.CognitoStoreInterface {
	return s.cognitoStore
}

// GetTimestreamStores returns the Timestream stores.
func (s *StoreProviderImpl) GetTimestreamStores() timestreamstore.TimestreamStoresInterface {
	return s.timestreamStores
}

// GetWAFStores returns the WAF stores.
func (s *StoreProviderImpl) GetWAFStores() wafstore.WAFStoresInterface { return s.wafStores }

// GetKMSStores returns the KMS stores.
func (s *StoreProviderImpl) GetKMSStores() kmsstore.KMSStoresInterface { return s.kmsStores }

// GetSESv2Store returns the SESv2 store.
func (s *StoreProviderImpl) GetSESv2Store() sesv2store.SESv2StoreInterface { return s.sesv2Store }

// GetSchedulerStore returns the EventBridge Scheduler store.
func (s *StoreProviderImpl) GetSchedulerStore() schedulerstore.SchedulerStoreInterface {
	return s.schedulerStore
}

// GetRoute53Stores returns the Route 53 stores.
func (s *StoreProviderImpl) GetRoute53Stores() route53store.Route53StoresInterface {
	return s.route53Stores
}

// GetCognitoIdentityStore returns the Cognito Identity store.
func (s *StoreProviderImpl) GetCognitoIdentityStore() cognitoidentitystore.CognitoIdentityStoreInterface {
	return s.cognitoIdentityStore
}

// GetCloudFrontStores returns the CloudFront stores.
func (s *StoreProviderImpl) GetCloudFrontStores() cloudfrontstore.CloudFrontStoresInterface {
	return s.cloudFrontStores
}

// GetCloudTrailStore returns the CloudTrail store.
func (s *StoreProviderImpl) GetCloudTrailStore() cloudtrailstore.CloudTrailStoreInterface {
	return s.cloudTrailStore
}

// GetSecretsManagerStore returns the Secrets Manager store.
func (s *StoreProviderImpl) GetSecretsManagerStore() secretsmanagerstore.SecretStoreInterface {
	return s.secretsManagerStore
}

// GetSNSStore returns the SNS store.
func (s *StoreProviderImpl) GetSNSStore() snsstore.SNSStoreInterface { return s.snsStore }

// GetSQSStore returns the SQS store.
func (s *StoreProviderImpl) GetSQSStore() sqsstore.SQSStoreInterface { return s.sqsStore }

var _ StoreProvider = (*StoreProviderImpl)(nil)
