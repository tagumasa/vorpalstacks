package eventbus

// RegisterInvoker registers a ServiceInvoker for the given service type.
func (b *EventBus) RegisterInvoker(invoker ServiceInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.invokers[invoker.ServiceType()] = invoker
}

// GetInvoker retrieves the ServiceInvoker registered for the given service type.
func (b *EventBus) GetInvoker(serviceType string) (ServiceInvoker, bool) {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	inv, ok := b.invokers[serviceType]
	return inv, ok
}

// SetLambdaInjector sets the Lambda invoker used for dispatching Lambda function
// invocations from bus events.
func (b *EventBus) SetLambdaInvoker(invoker LambdaInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.lambdaInvoker = invoker
}

// LambdaInvoker returns the configured Lambda invoker.
func (b *EventBus) LambdaInvoker() LambdaInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.lambdaInvoker
}

// SetSQSInvoker sets the SQS invoker used for dispatching SQS SendMessage calls
// from bus events.
func (b *EventBus) SetSQSInvoker(invoker SQSInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.sqsInvoker = invoker
}

// SQSInvoker returns the configured SQS invoker.
func (b *EventBus) SQSInvoker() SQSInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.sqsInvoker
}

// SetSNSInvoker sets the SNS invoker used for dispatching SNS Publish calls
// from bus events.
func (b *EventBus) SetSNSInvoker(invoker SNSInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.snsInvoker = invoker
}

// SNSInvoker returns the configured SNS invoker.
func (b *EventBus) SNSInvoker() SNSInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.snsInvoker
}

// SetKinesisInvoker sets the Kinesis invoker used for dispatching Kinesis
// PutRecord calls from bus events.
func (b *EventBus) SetKinesisInvoker(invoker KinesisInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.kinesisInvoker = invoker
}

// KinesisInvoker returns the configured Kinesis invoker.
func (b *EventBus) KinesisInvoker() KinesisInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.kinesisInvoker
}

// SetEventsInvoker sets the EventBridge invoker used for dispatching
// EventBridge PutEvents calls from bus events.
func (b *EventBus) SetEventsInvoker(invoker EventsInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.eventsInvoker = invoker
}

// EventsInvoker returns the configured EventBridge invoker.
func (b *EventBus) EventsInvoker() EventsInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.eventsInvoker
}

// SetEC2Invoker sets the EC2 invoker used for dispatching EC2 API calls
// from bus events.
func (b *EventBus) SetEC2Invoker(invoker EC2Invoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.ec2Invoker = invoker
}

// EC2Invoker returns the configured EC2 invoker.
func (b *EventBus) EC2Invoker() EC2Invoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.ec2Invoker
}

// SetDynamoDBInvoker sets the DynamoDB invoker used for dispatching DynamoDB
// item operations from bus events (e.g. AppSync GraphQL resolvers).
func (b *EventBus) SetDynamoDBInvoker(invoker DynamoDBInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.dynamoDBInvoker = invoker
}

// DynamoDBInvoker returns the configured DynamoDB invoker.
func (b *EventBus) DynamoDBInvoker() DynamoDBInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.dynamoDBInvoker
}

// SetDynamoDBStreamsInvoker sets the DynamoDB Streams invoker used for
// polling stream records (e.g. by the Lambda ESM poller).
func (b *EventBus) SetDynamoDBStreamsInvoker(invoker DynamoDBStreamsInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.dynamoDBStreamsInvoker = invoker
}

// DynamoDBStreamsInvoker returns the configured DynamoDB Streams invoker.
func (b *EventBus) DynamoDBStreamsInvoker() DynamoDBStreamsInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.dynamoDBStreamsInvoker
}

// SetNeptuneGraphInvoker sets the NeptuneGraph invoker used for dispatching
// graph queries from bus consumers (e.g. AppSync GraphQL resolvers).
func (b *EventBus) SetNeptuneGraphInvoker(invoker NeptuneGraphInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.neptuneGraphInvoker = invoker
}

// NeptuneGraphInvoker returns the configured NeptuneGraph invoker.
func (b *EventBus) NeptuneGraphInvoker() NeptuneGraphInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.neptuneGraphInvoker
}

// SetKMSInvoker sets the KMS invoker used for encryption key operations.
func (b *EventBus) SetKMSInvoker(invoker KMSInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.kmsInvoker = invoker
}

// KMSInvoker returns the configured KMS invoker.
func (b *EventBus) KMSInvoker() KMSInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.kmsInvoker
}

// SetS3Invoker sets the S3 invoker used for cross-service S3 object operations.
func (b *EventBus) SetS3Invoker(invoker S3Invoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.s3Invoker = invoker
}

// S3Invoker returns the configured S3 invoker.
func (b *EventBus) S3Invoker() S3Invoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.s3Invoker
}

// SetWAFInvoker sets the WAF invoker used for WebACL association operations.
func (b *EventBus) SetWAFInvoker(invoker WAFInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.wafInvoker = invoker
}

// WAFInvoker returns the configured WAF invoker.
func (b *EventBus) WAFInvoker() WAFInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.wafInvoker
}

// SetWebACLInspector sets the WAF request-inspection entry point used by
// protected-resource planes to enforce associated WebACLs.
func (b *EventBus) SetWebACLInspector(inspector WebACLInspector) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.webACLInspector = inspector
}

// WebACLInspector returns the configured WAF request-inspection entry
// point, or nil if none was set (WAFv2 not initialised).
func (b *EventBus) WebACLInspector() WebACLInspector {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.webACLInspector
}

// SetCloudWatchMetricInvoker sets the CloudWatch metric invoker used for
// cross-service metric data submission (e.g. CloudWatch Logs metric filters).
func (b *EventBus) SetCloudWatchMetricInvoker(invoker CloudWatchMetricInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.cloudWatchMetricInvoker = invoker
}

// CloudWatchMetricInvoker returns the configured CloudWatch metric invoker.
func (b *EventBus) CloudWatchMetricInvoker() CloudWatchMetricInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.cloudWatchMetricInvoker
}

// SetCloudWatchAlarmInvoker sets the CloudWatch alarm invoker.
func (b *EventBus) SetCloudWatchAlarmInvoker(invoker CloudWatchAlarmInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.cloudWatchAlarmInvoker = invoker
}

// CloudWatchAlarmInvoker returns the configured CloudWatch alarm invoker.
func (b *EventBus) CloudWatchAlarmInvoker() CloudWatchAlarmInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.cloudWatchAlarmInvoker
}

// SetTimestreamInvoker sets the Timestream invoker.
func (b *EventBus) SetTimestreamInvoker(invoker TimestreamInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.timestreamInvoker = invoker
}

// TimestreamInvoker returns the configured Timestream invoker.
func (b *EventBus) TimestreamInvoker() TimestreamInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.timestreamInvoker
}

// SetCloudTrailInvoker sets the CloudTrail invoker used for event lookup
// (e.g. IAM GenerateServiceLastAccessedDetails).
func (b *EventBus) SetCloudTrailInvoker(invoker CloudTrailInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.cloudTrailInvoker = invoker
}

// CloudTrailInvoker returns the configured CloudTrail invoker.
func (b *EventBus) CloudTrailInvoker() CloudTrailInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.cloudTrailInvoker
}

// SetLogsInvoker sets the CloudWatch Logs invoker used for cross-service log
// delivery (e.g. Lambda function log output).
func (b *EventBus) SetLogsInvoker(invoker LogsInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.logsInvoker = invoker
}

// LogsInvoker returns the configured CloudWatch Logs invoker.
func (b *EventBus) LogsInvoker() LogsInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.logsInvoker
}

// SetRDSDataInvoker sets the RDS Data API invoker for cross-service SQL execution.
func (b *EventBus) SetRDSDataInvoker(invoker RDSDataInvoker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.rdsDataInvoker = invoker
}

// RDSDataInvoker returns the configured RDS Data API invoker.
func (b *EventBus) RDSDataInvoker() RDSDataInvoker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.rdsDataInvoker
}

// SetCognitoTokenValidator sets the Cognito token validator for cross-service
// JWT validation (e.g. API Gateway COGNITO_USER_POOLS authorizer).
func (b *EventBus) SetCognitoTokenValidator(validator CognitoTokenValidator) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.cognitoTokenValidator = validator
}

// CognitoTokenValidator returns the configured Cognito token validator.
func (b *EventBus) CognitoTokenValidator() CognitoTokenValidator {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	return b.cognitoTokenValidator
}

// RegisterSubnetUsageChecker registers a service that can report whether a
// subnet is in use. Multiple services may register (e.g. Lambda, Neptune).
// EC2 calls all registered checkers before deleting a subnet.
func (b *EventBus) RegisterSubnetUsageChecker(checker SubnetUsageChecker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.subnetUsageCheckers = append(b.subnetUsageCheckers, checker)
}

// RegisterSecurityGroupUsageChecker registers a service that can report
// whether a security group is in use. Multiple services may register.
// EC2 calls all registered checkers before deleting a security group.
func (b *EventBus) RegisterSecurityGroupUsageChecker(checker SecurityGroupUsageChecker) {
	b.invokersMu.Lock()
	defer b.invokersMu.Unlock()
	b.securityGroupCheckers = append(b.securityGroupCheckers, checker)
}

// SubnetUsageCheckers returns all registered subnet usage checkers.
func (b *EventBus) SubnetUsageCheckers() []SubnetUsageChecker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	result := make([]SubnetUsageChecker, len(b.subnetUsageCheckers))
	copy(result, b.subnetUsageCheckers)
	return result
}

// SecurityGroupUsageCheckers returns all registered security group usage checkers.
func (b *EventBus) SecurityGroupUsageCheckers() []SecurityGroupUsageChecker {
	b.invokersMu.RLock()
	defer b.invokersMu.RUnlock()
	result := make([]SecurityGroupUsageChecker, len(b.securityGroupCheckers))
	copy(result, b.securityGroupCheckers)
	return result
}

// RoleResolver returns the configured RoleResolver, or nil if none was set.
func (b *EventBus) RoleResolver() RoleResolver {
	return b.roleResolver
}
