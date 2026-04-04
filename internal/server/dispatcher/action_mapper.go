package dispatcher

import (
	"strings"
)

// ActionMapper maps AWS service operations to IAM action strings.
// It is used to translate gRPC-Web operation names to their corresponding IAM actions.
type ActionMapper struct {
	mappings map[string]map[string]string
}

// NewActionMapper creates a new ActionMapper instance.
// It initialises the mapper with default service-to-action mappings.
func NewActionMapper() *ActionMapper {
	m := &ActionMapper{
		mappings: make(map[string]map[string]string),
	}
	m.registerDefaults()
	return m
}

// Map translates a service name and operation to its corresponding IAM action.
// If a specific mapping exists, it returns that action.
// If a wildcard mapping exists, it replaces the wildcard with the operation name.
// Otherwise, it returns "service:operation" as the default action.
func (m *ActionMapper) Map(service, operation string) string {
	if serviceMappings, ok := m.mappings[service]; ok {
		if action, ok := serviceMappings[operation]; ok {
			return action
		}
		if action, ok := serviceMappings["*"]; ok {
			return strings.Replace(action, "*", operation, 1)
		}
	}
	return service + ":" + operation
}

// Register registers a specific action mapping for a service and operation.
// This allows custom action translations for individual operations.
func (m *ActionMapper) Register(service, operation, action string) {
	if m.mappings[service] == nil {
		m.mappings[service] = make(map[string]string)
	}
	m.mappings[service][operation] = action
}

// RegisterWildcard registers a wildcard mapping for a service.
// All operations for that service will use the prefix with the operation name substituted.
func (m *ActionMapper) RegisterWildcard(service, prefix string) {
	if m.mappings[service] == nil {
		m.mappings[service] = make(map[string]string)
	}
	m.mappings[service]["*"] = prefix + ":*"
}

func (m *ActionMapper) registerDefaults() {
	m.registerS3Mappings()
	m.registerSQSMappings()
	m.registerDynamoDBMappings()
	m.registerKMSMappings()
	m.registerLambdaMappings()
	m.registerSNSMappings()
	m.registerEventBridgeMappings()
	m.registerIAMMappings()
	m.registerKinesisMappings()
	m.registerCloudWatchMappings()
	m.registerStepFunctionsMappings()
	m.registerAPIGatewayMappings()
	m.registerAthenaMappings()
	m.registerSecretsManagerMappings()
	m.registerSTSMappings()
	m.registerCognitoMappings()
	m.registerSchedulerMappings()
	m.registerCloudFrontMappings()
	m.registerCloudTrailMappings()
	m.registerTimestreamMappings()
	m.registerWAFMappings()
	m.registerElastiCacheMappings()
	m.registerRedshiftMappings()
	m.registerRDSMappings()
	m.registerEC2Mappings()
	m.registerNeptuneMappings()
	m.registerAppSyncMappings()
}

func (m *ActionMapper) registerS3Mappings() {
	m.RegisterWildcard("s3", "s3")
}

func (m *ActionMapper) registerSQSMappings() {
	m.RegisterWildcard("sqs", "sqs")
}

func (m *ActionMapper) registerDynamoDBMappings() {
	m.RegisterWildcard("dynamodb", "dynamodb")
}

func (m *ActionMapper) registerKMSMappings() {
	m.RegisterWildcard("kms", "kms")
}

func (m *ActionMapper) registerLambdaMappings() {
	m.RegisterWildcard("lambda", "lambda")
}

func (m *ActionMapper) registerSNSMappings() {
	m.RegisterWildcard("sns", "sns")
}

func (m *ActionMapper) registerEventBridgeMappings() {
	m.RegisterWildcard("events", "events")
}

func (m *ActionMapper) registerIAMMappings() {
	m.RegisterWildcard("iam", "iam")
}

func (m *ActionMapper) registerKinesisMappings() {
	m.RegisterWildcard("kinesis", "kinesis")
}

func (m *ActionMapper) registerCloudWatchMappings() {
	m.RegisterWildcard("logs", "logs")
	m.RegisterWildcard("cloudwatch", "cloudwatch")
}

func (m *ActionMapper) registerStepFunctionsMappings() {
	m.RegisterWildcard("states", "states")
}

func (m *ActionMapper) registerAPIGatewayMappings() {
	m.RegisterWildcard("apigateway", "apigateway")
}

func (m *ActionMapper) registerAthenaMappings() {
	m.RegisterWildcard("athena", "athena")
}

func (m *ActionMapper) registerSecretsManagerMappings() {
	m.RegisterWildcard("secretsmanager", "secretsmanager")
}

func (m *ActionMapper) registerSTSMappings() {
	m.RegisterWildcard("sts", "sts")
}

func (m *ActionMapper) registerCognitoMappings() {
	m.RegisterWildcard("cognito-idp", "cognito-idp")
	m.RegisterWildcard("cognito-identity", "cognito-identity")
}

func (m *ActionMapper) registerSchedulerMappings() {
	m.RegisterWildcard("scheduler", "scheduler")
}

func (m *ActionMapper) registerCloudFrontMappings() {
	m.RegisterWildcard("cloudfront", "cloudfront")
}

func (m *ActionMapper) registerCloudTrailMappings() {
	m.RegisterWildcard("cloudtrail", "cloudtrail")
}

func (m *ActionMapper) registerTimestreamMappings() {
	m.RegisterWildcard("timestream", "timestream")
	m.Register("timestream-write", "WriteRecords", "timestream:WriteRecords")
	m.Register("timestream-write", "DescribeEndpoints", "timestream:DescribeEndpoints")
	m.Register("timestream-write", "*", "timestream:*")
	m.Register("timestream-query", "Query", "timestream:Query")
	m.Register("timestream-query", "DescribeEndpoints", "timestream:DescribeEndpoints")
	m.Register("timestream-query", "*", "timestream:*")
}

func (m *ActionMapper) registerWAFMappings() {
	m.RegisterWildcard("waf", "waf")
	m.RegisterWildcard("wafv2", "wafv2")
}

func (m *ActionMapper) registerElastiCacheMappings() {
	m.RegisterWildcard("elasticache", "elasticache")
}

func (m *ActionMapper) registerRedshiftMappings() {
	m.RegisterWildcard("redshift", "redshift")
}

func (m *ActionMapper) registerRDSMappings() {
	m.RegisterWildcard("rds", "rds")
}

func (m *ActionMapper) registerEC2Mappings() {
	m.RegisterWildcard("ec2", "ec2")
}

func (m *ActionMapper) registerNeptuneMappings() {
	m.RegisterWildcard("neptune", "neptune")
	m.RegisterWildcard("neptunedata", "neptune-db")
}

// registerAppSyncMappings registers IAM action mappings for the AWS AppSync service.
func (m *ActionMapper) registerAppSyncMappings() {
	m.RegisterWildcard("appsync", "appsync")
}
