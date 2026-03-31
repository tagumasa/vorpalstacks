package request

import (
	"context"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/services/aws/common/iam"
	"vorpalstacks/internal/services/aws/common/interfaces"
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
	neptunestore "vorpalstacks/internal/store/aws/neptune"
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
	"vorpalstacks/pkg/graphengine"
)

// AuditRecorder records audit events for CloudTrail.
type AuditRecorder interface {
	Record(event interface{}) error
}

// PrincipalType represents the type of principal making a request.
type PrincipalType string

const (
	// PrincipalTypeUser represents an IAM user.
	PrincipalTypeUser PrincipalType = "User"
	// PrincipalTypeRole represents an IAM role.
	PrincipalTypeRole PrincipalType = "Role"
	// PrincipalTypeAnonymous represents an unauthenticated request.
	PrincipalTypeAnonymous PrincipalType = "Anonymous"
)

// RequestContext holds the context information for an AWS API request.
type RequestContext struct {
	context.Context
	storageManager *storage.RegionStorageManager
	storeProvider  interfaces.StoreProvider
	iamValidator   *iam.IAMValidator
	AccountID      string
	Region         string
	Principal      string
	PrincipalID    string
	PrincipalType  PrincipalType
	SourceIP       string
	UserAgent      string
	auditRecorder  AuditRecorder
	graphDBManager *graphengine.DB
}

// NewRequestContext creates a new RequestContext with the given parameters.
func NewRequestContext(ctx context.Context, storageMgr *storage.RegionStorageManager, accountID, region string) *RequestContext {
	if ctx == nil {
		ctx = context.Background()
	}
	if region == "" {
		region = DefaultRegion
	}
	return &RequestContext{
		Context:        ctx,
		storageManager: storageMgr,
		AccountID:      accountID,
		Region:         region,
	}
}

// GetRegion returns the AWS region for the request.
func (c *RequestContext) GetRegion() string {
	if c.Region == "" {
		return DefaultRegion
	}
	return c.Region
}

// GetAccountID returns the AWS account ID for the request.
func (c *RequestContext) GetAccountID() string {
	return c.AccountID
}

// GetStorage returns the storage for the request's region.
func (c *RequestContext) GetStorage() (storage.BasicStorage, error) {
	return c.storageManager.GetStorage(c.GetRegion())
}

// GetGlobalStorage returns the global storage for region-independent services.
func (c *RequestContext) GetGlobalStorage() (storage.BasicStorage, error) {
	return c.storageManager.GetGlobalStorage()
}

// GetStorageForService returns the appropriate storage for the given service.
func (c *RequestContext) GetStorageForService(serviceName string) (storage.BasicStorage, error) {
	if c.IsGlobalService(serviceName) {
		return c.GetGlobalStorage()
	}
	return c.GetStorage()
}

// IsGlobalService returns true if the service is global (not region-specific).
func (c *RequestContext) IsGlobalService(serviceName string) bool {
	globalServices := map[string]bool{
		"iam":        true,
		"route53":    true,
		"cloudfront": true,
		"sts":        true,
	}
	return globalServices[serviceName]
}

// SetAuditRecorder sets the audit recorder for the request context.
func (c *RequestContext) SetAuditRecorder(recorder AuditRecorder) {
	c.auditRecorder = recorder
}

// GetAuditRecorder returns the audit recorder for the request context.
func (c *RequestContext) GetAuditRecorder() AuditRecorder {
	return c.auditRecorder
}

// HasAuditRecorder returns true if an audit recorder is set.
func (c *RequestContext) HasAuditRecorder() bool {
	return c.auditRecorder != nil
}

// SetStoreProvider sets the store provider for the request context.
func (c *RequestContext) SetStoreProvider(provider interfaces.StoreProvider) {
	c.storeProvider = provider
}

// GetStoreProvider returns the store provider for the request context.
func (c *RequestContext) GetStoreProvider() interfaces.StoreProvider {
	return c.storeProvider
}

// GetIAMStore returns the IAM store from the store provider.
func (c *RequestContext) GetIAMStore() iamstore.IAMStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetIAMStore()
}

// GetIAMValidator returns the IAM validator for the request context.
func (c *RequestContext) GetIAMValidator() *iam.IAMValidator {
	if c.iamValidator == nil && c.storeProvider != nil {
		iamStore := c.storeProvider.GetIAMStore()
		if iamStore != nil {
			c.iamValidator = iam.NewIAMValidator(iamStore.Roles(), c.AccountID)
		}
	}
	return c.iamValidator
}

// GetACMStore returns the ACM store from the store provider.
func (c *RequestContext) GetACMStore() acmstore.CertificateStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetACMStore()
}

// GetS3Store returns the S3 store from the store provider.
func (c *RequestContext) GetS3Store() s3store.S3StoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetS3Store()
}

// GetSSMStore returns the SSM store from the store provider.
func (c *RequestContext) GetSSMStore() ssmstore.SSMStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSSMStore()
}

// GetSTSStore returns the STS store from the store provider.
func (c *RequestContext) GetSTSStore() stsstore.SessionStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSTSStore()
}

// GetDynamoDBStore returns the DynamoDB store from the store provider.
func (c *RequestContext) GetDynamoDBStore() dynamodbstore.DynamoDBStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetDynamoDBStore()
}

// GetLambdaStore returns the Lambda store from the store provider.
func (c *RequestContext) GetLambdaStore() lambdastore.LambdaStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetLambdaStore()
}

// GetEventBridgeStore returns the EventBridge store from the store provider.
func (c *RequestContext) GetEventBridgeStore() eventbridgestore.EventsStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetEventBridgeStore()
}

// GetCloudWatchAlarmStore returns the CloudWatch Alarm store from the store provider.
func (c *RequestContext) GetCloudWatchAlarmStore() cloudwatchstore.AlarmStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCloudWatchAlarmStore()
}

// GetCloudWatchMetricStore returns the CloudWatch Metric store from the store provider.
func (c *RequestContext) GetCloudWatchMetricStore() cloudwatchstore.MetricChunkStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCloudWatchMetricStore()
}

// GetAPIGatewayStore returns the API Gateway store from the store provider.
func (c *RequestContext) GetAPIGatewayStore() apigatewaystore.APIGatewayStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetAPIGatewayStore()
}

// GetSFNStore returns the Step Functions store from the store provider.
func (c *RequestContext) GetSFNStore() sfnstore.StepFunctionStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSFNStore()
}

// GetAthenaStore returns the Athena store from the store provider.
func (c *RequestContext) GetAthenaStore() athenastore.AthenaStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetAthenaStore()
}

// GetCognitoStore returns the Cognito User Pool store from the store provider.
func (c *RequestContext) GetCognitoStore() cognitostore.CognitoStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCognitoStore()
}

// GetTimestreamStores returns the Timestream stores from the store provider.
func (c *RequestContext) GetTimestreamStores() timestreamstore.TimestreamStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetTimestreamStores()
}

// GetWAFStores returns the WAF stores from the store provider.
func (c *RequestContext) GetWAFStores() wafstore.WAFStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetWAFStores()
}

// GetKMSStores returns the KMS stores from the store provider.
func (c *RequestContext) GetKMSStores() kmsstore.KMSStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetKMSStores()
}

// GetSESv2Store returns the SESv2 store from the store provider.
func (c *RequestContext) GetSESv2Store() sesv2store.SESv2StoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSESv2Store()
}

// GetCloudWatchLogsStore returns the CloudWatch Logs store from the store provider.
func (c *RequestContext) GetCloudWatchLogsStore() cloudwatchlogsstore.CloudWatchLogsStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCloudWatchLogsStore()
}

// GetSchedulerStore returns the EventBridge Scheduler store from the store provider.
func (c *RequestContext) GetSchedulerStore() schedulerstore.SchedulerStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSchedulerStore()
}

// GetRoute53Stores returns the Route 53 stores from the store provider.
func (c *RequestContext) GetRoute53Stores() route53store.Route53StoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetRoute53Stores()
}

// GetCognitoIdentityStore returns the Cognito Identity store from the store provider.
func (c *RequestContext) GetCognitoIdentityStore() cognitoidentitystore.CognitoIdentityStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCognitoIdentityStore()
}

// GetCloudFrontStores returns the CloudFront stores from the store provider.
func (c *RequestContext) GetCloudFrontStores() cloudfrontstore.CloudFrontStoresInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCloudFrontStores()
}

// GetCloudTrailStore returns the CloudTrail store from the store provider.
func (c *RequestContext) GetCloudTrailStore() cloudtrailstore.CloudTrailStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetCloudTrailStore()
}

// GetSecretsManagerStore returns the Secrets Manager store from the store provider.
func (c *RequestContext) GetSecretsManagerStore() secretsmanagerstore.SecretStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSecretsManagerStore()
}

// GetSNSStore returns the SNS store from the store provider.
func (c *RequestContext) GetSNSStore() snsstore.SNSStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSNSStore()
}

// GetSQSStore returns the SQS store from the store provider.
func (c *RequestContext) GetSQSStore() sqsstore.SQSStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetSQSStore()
}

func (c *RequestContext) GetNeptuneStore() neptunestore.NeptuneStoreInterface {
	if c.storeProvider == nil {
		return nil
	}
	return c.storeProvider.GetNeptuneStore()
}

func (c *RequestContext) SetGraphDBManager(db *graphengine.DB) {
	c.graphDBManager = db
}

func (c *RequestContext) GraphReader() graphengine.GraphReader {
	return c.graphDBManager
}

func (c *RequestContext) GraphWriter() graphengine.GraphWriter {
	return c.graphDBManager
}
