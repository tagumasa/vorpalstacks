package invokers

// Registry is the typed invoker registry that the server application
// wires at startup and cross-service consumers read from. Each setter
// installs a service's invoker implementation; each getter returns it,
// or nil when the providing service is not initialised. The event bus
// implements this interface alongside its delivery surface, so a
// consumer holding the bus can reach every invoker through one object.
type Registry interface {
	LambdaInvoker() LambdaInvoker
	SetLambdaInvoker(invoker LambdaInvoker)
	SQSInvoker() SQSInvoker
	SetSQSInvoker(invoker SQSInvoker)
	SNSInvoker() SNSInvoker
	SetSNSInvoker(invoker SNSInvoker)
	KinesisInvoker() KinesisInvoker
	SetKinesisInvoker(invoker KinesisInvoker)
	EventsInvoker() EventsInvoker
	SetEventsInvoker(invoker EventsInvoker)
	EC2Invoker() EC2Invoker
	SetEC2Invoker(invoker EC2Invoker)
	DynamoDBInvoker() DynamoDBInvoker
	SetDynamoDBInvoker(invoker DynamoDBInvoker)
	DynamoDBStreamsInvoker() DynamoDBStreamsInvoker
	SetDynamoDBStreamsInvoker(invoker DynamoDBStreamsInvoker)
	NeptuneGraphInvoker() NeptuneGraphInvoker
	SetNeptuneGraphInvoker(invoker NeptuneGraphInvoker)
	KMSInvoker() KMSInvoker
	SetKMSInvoker(invoker KMSInvoker)
	S3Invoker() S3Invoker
	SetS3Invoker(invoker S3Invoker)
	WAFInvoker() WAFInvoker
	SetWAFInvoker(invoker WAFInvoker)
	CloudWatchMetricInvoker() CloudWatchMetricInvoker
	SetCloudWatchMetricInvoker(invoker CloudWatchMetricInvoker)
	CloudWatchAlarmInvoker() CloudWatchAlarmInvoker
	SetCloudWatchAlarmInvoker(invoker CloudWatchAlarmInvoker)
	TimestreamInvoker() TimestreamInvoker
	SetTimestreamInvoker(invoker TimestreamInvoker)
	CloudTrailInvoker() CloudTrailInvoker
	SetCloudTrailInvoker(invoker CloudTrailInvoker)
	LogsInvoker() LogsInvoker
	SetLogsInvoker(invoker LogsInvoker)
	RDSDataInvoker() RDSDataInvoker
	SetRDSDataInvoker(invoker RDSDataInvoker)
	CognitoTokenValidator() CognitoTokenValidator
	SetCognitoTokenValidator(validator CognitoTokenValidator)
	RegisterSubnetUsageChecker(checker SubnetUsageChecker)
	RegisterSecurityGroupUsageChecker(checker SecurityGroupUsageChecker)
	SubnetUsageCheckers() []SubnetUsageChecker
	SecurityGroupUsageCheckers() []SecurityGroupUsageChecker
}
