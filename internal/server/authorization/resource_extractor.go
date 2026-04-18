package authorization

import (
	"net/url"
	"strings"

	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// ResourceExtractorFunc is a function that extracts the resource ARN from request parameters.
// It takes the request parameters, account ID, and region, and returns the resource ARN.
type ResourceExtractorFunc func(params map[string]interface{}, accountID, region string) string

// ResourceExtractor extracts resource ARNs from AWS API request parameters.
// Different AWS services require different extraction strategies.
type ResourceExtractor struct {
	extractors map[string]map[string]ResourceExtractorFunc
}

// NewResourceExtractor creates a new ResourceExtractor instance.
// It registers default extractors for all supported AWS services.
func NewResourceExtractor() *ResourceExtractor {
	e := &ResourceExtractor{
		extractors: make(map[string]map[string]ResourceExtractorFunc),
	}

	e.registerDefaults()
	return e
}

// Extract extracts the resource ARN from the given request parameters.
// It looks for a service-specific extractor, falling back to a wildcard extractor if available.
// Returns "*" if no specific extractor is found.
func (e *ResourceExtractor) Extract(service, operation string, params map[string]interface{}, accountID, region string) string {
	if serviceExtractors, ok := e.extractors[service]; ok {
		if extractor, ok := serviceExtractors[operation]; ok {
			return extractor(params, accountID, region)
		}
		if extractor, ok := serviceExtractors["*"]; ok {
			return extractor(params, accountID, region)
		}
	}

	return "*"
}

// Register registers a custom resource extractor for a specific service and operation.
// This allows custom ARN extraction logic for operations that don't follow standard patterns.
func (e *ResourceExtractor) Register(service, operation string, fn ResourceExtractorFunc) {
	if e.extractors[service] == nil {
		e.extractors[service] = make(map[string]ResourceExtractorFunc)
	}
	e.extractors[service][operation] = fn
}

func (e *ResourceExtractor) registerDefaults() {
	e.registerS3Extractors()
	e.registerSQSExtractors()
	e.registerDynamoDBExtractors()
	e.registerKMSExtractors()
	e.registerLambdaExtractors()
	e.registerSNExtractors()
	e.registerEventBridgeExtractors()
	e.registerIAMExtractors()
	e.registerKinesisExtractors()
	e.registerCloudWatchExtractors()
	e.registerStepFunctionsExtractors()
	e.registerAPIGatewayExtractors()
	e.registerAthenaExtractors()
	e.registerSecretsManagerExtractors()
	e.registerAppSyncExtractors()
	e.registerNeptuneExtractors()
	e.registerNeptuneGraphExtractors()
}

func s3BucketArn(accountID, region string, bucket string) string {
	return arnutil.NewARNBuilder(accountID, region).S3().Bucket(bucket)
}

func s3ObjectArn(accountID, region string, bucket, key string) string {
	return arnutil.NewARNBuilder(accountID, region).S3().Object(bucket, key)
}

func sqsQueueArn(accountID, region string, queueName string) string {
	return arnutil.NewARNBuilder(accountID, region).SQS().Queue(queueName)
}

func dynamodbTableArn(accountID, region string, tableName string) string {
	return arnutil.NewARNBuilder(accountID, region).DynamoDB().Table(tableName)
}

func kmsKeyArn(accountID, region string, keyID string) string {
	return arnutil.NewARNBuilder(accountID, region).KMS().Key(keyID)
}

func kmsAliasArn(accountID, region string, aliasName string) string {
	return arnutil.NewARNBuilder(accountID, region).KMS().Alias(aliasName)
}

func lambdaFunctionArn(accountID, region string, functionName string) string {
	return arnutil.NewARNBuilder(accountID, region).Lambda().Function(functionName)
}

func snsTopicArn(accountID, region string, topicName string) string {
	return arnutil.NewARNBuilder(accountID, region).SNS().Topic(topicName)
}

func eventsEventBusArn(accountID, region string, busName string) string {
	return arnutil.NewARNBuilder(accountID, region).Events().EventBus(busName)
}

func iamUserArn(accountID, region string, userName string) string {
	return arnutil.NewARNBuilder(accountID, region).IAM().User(userName)
}

func iamRoleArn(accountID, region string, roleName string) string {
	return arnutil.NewARNBuilder(accountID, region).IAM().Role(roleName)
}

func iamPolicyArn(accountID, region string, policyName string) string {
	return arnutil.NewARNBuilder(accountID, region).IAM().Policy(policyName)
}

func kinesisStreamArn(accountID, region string, streamName string) string {
	return arnutil.NewARNBuilder(accountID, region).Kinesis().Stream(streamName)
}

func logsLogGroupArn(accountID, region string, logGroupName string) string {
	return arnutil.NewARNBuilder(accountID, region).CloudWatch().LogGroup(logGroupName)
}

func statesStateMachineArn(accountID, region string, name string) string {
	return arnutil.NewARNBuilder(accountID, region).StepFunctions().StateMachine(name)
}

func secretsManagerSecretArn(accountID, region string, secretID string) string {
	return arnutil.NewARNBuilder(accountID, region).SecretsManager().Secret(secretID)
}

func athenaWorkGroupArn(accountID, region string, workGroup string) string {
	return arnutil.NewARNBuilder(accountID, region).Athena().WorkGroup(workGroup)
}

func appsyncApiArn(accountID, region string, apiId string) string {
	return arnutil.NewARNBuilder(accountID, region).AppSync().Api(apiId)
}

func neptuneClusterArn(accountID, region string, clusterId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "cluster/"+clusterId)
}

func neptuneClusterSnapshotArn(accountID, region string, snapshotId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "cluster-snapshot/"+snapshotId)
}

func neptuneParamGroupArn(accountID, region string, groupName string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "cluster-pg/"+groupName)
}

func neptuneSubnetGroupArn(accountID, region string, groupName string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "subnet-group/"+groupName)
}

func neptuneInstanceArn(accountID, region string, instanceId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "db/"+instanceId)
}

func neptuneGlobalClusterArn(accountID, region string, clusterId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "global-cluster/"+clusterId)
}

func neptuneEventSubscriptionArn(accountID, region string, subName string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune", "event-subscription/"+subName)
}

func neptuneGraphArn(accountID, region string, graphId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune-graph", "graph/"+graphId)
}

func neptuneGraphSnapshotArn(accountID, region string, snapshotId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune-graph", "snapshot/"+snapshotId)
}

func neptuneGraphImportTaskArn(accountID, region string, taskId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune-graph", "import-task/"+taskId)
}

func neptuneGraphExportTaskArn(accountID, region string, taskId string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("neptune-graph", "export-task/"+taskId)
}

func (e *ResourceExtractor) registerS3Extractors() {
	s3Extractor := func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		key, _ := params["Key"].(string)
		if key == "" {
			return s3BucketArn(accountID, region, bucket)
		}
		return s3ObjectArn(accountID, region, bucket, key)
	}

	services := []string{"s3"}
	operations := []string{
		"GetObject", "PutObject", "DeleteObject", "HeadObject",
		"CopyObject", "CreateMultipartUpload", "UploadPart", "CompleteMultipartUpload",
		"AbortMultipartUpload", "ListParts", "SelectObjectContent",
	}
	for _, svc := range services {
		for _, op := range operations {
			e.Register(svc, op, s3Extractor)
		}
	}

	s3BucketExtractor := func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return s3BucketArn(accountID, region, bucket)
	}

	for _, op := range []string{
		"ListBucket", "CreateBucket", "DeleteBucket", "GetBucketLocation",
		"PutBucketPolicy", "GetBucketPolicy", "DeleteBucketPolicy",
	} {
		e.Register("s3", op, s3BucketExtractor)
	}
}

func (e *ResourceExtractor) registerSQSExtractors() {
	sqsExtractor := func(params map[string]interface{}, accountID, region string) string {
		queueURL, _ := params["QueueUrl"].(string)
		if queueURL == "" {
			return "*"
		}
		queueName := extractQueueNameFromURL(queueURL)
		return sqsQueueArn(accountID, region, queueName)
	}

	operations := []string{
		"SendMessage", "ReceiveMessage", "DeleteMessage", "ChangeMessageVisibility",
		"SendMessageBatch", "DeleteMessageBatch", "ChangeMessageVisibilityBatch",
		"GetQueueAttributes", "SetQueueAttributes", "PurgeQueue",
	}
	for _, op := range operations {
		e.Register("sqs", op, sqsExtractor)
	}

	e.Register("sqs", "CreateQueue", func(params map[string]interface{}, accountID, region string) string {
		queueName, _ := params["QueueName"].(string)
		if queueName == "" {
			return "*"
		}
		return sqsQueueArn(accountID, region, queueName)
	})

	e.Register("sqs", "DeleteQueue", func(params map[string]interface{}, accountID, region string) string {
		queueURL, _ := params["QueueUrl"].(string)
		if queueURL == "" {
			return "*"
		}
		queueName := extractQueueNameFromURL(queueURL)
		return sqsQueueArn(accountID, region, queueName)
	})

	e.Register("sqs", "ListQueues", func(params map[string]interface{}, accountID, region string) string {
		return sqsQueueArn(accountID, region, "*")
	})
}

func (e *ResourceExtractor) registerDynamoDBExtractors() {
	dynamoExtractor := func(params map[string]interface{}, accountID, region string) string {
		tableName, _ := params["TableName"].(string)
		if tableName == "" {
			return "*"
		}
		return dynamodbTableArn(accountID, region, tableName)
	}

	operations := []string{
		"GetItem", "PutItem", "UpdateItem", "DeleteItem", "Query", "Scan",
		"BatchGetItem", "BatchWriteItem", "TransactGetItems", "TransactWriteItems",
		"UpdateTable", "DeleteTable", "DescribeTable", "GetRecords",
	}
	for _, op := range operations {
		e.Register("dynamodb", op, dynamoExtractor)
	}

	e.Register("dynamodb", "CreateTable", func(params map[string]interface{}, accountID, region string) string {
		tableName, _ := params["TableName"].(string)
		if tableName == "" {
			return "*"
		}
		return dynamodbTableArn(accountID, region, tableName)
	})

	e.Register("dynamodb", "ListTables", func(params map[string]interface{}, accountID, region string) string {
		return dynamodbTableArn(accountID, region, "*")
	})
}

func (e *ResourceExtractor) registerKMSExtractors() {
	kmsExtractor := func(params map[string]interface{}, accountID, region string) string {
		keyID := extractKMSKeyID(params)
		if keyID == "" {
			return "*"
		}
		if strings.HasPrefix(keyID, "arn:") {
			return keyID
		}
		return kmsKeyArn(accountID, region, keyID)
	}

	operations := []string{
		"Encrypt", "Decrypt", "ReEncrypt", "GenerateDataKey", "GenerateDataKeyWithoutPlaintext",
		"GenerateRandom", "Sign", "Verify", "GetPublicKey",
	}
	for _, op := range operations {
		e.Register("kms", op, kmsExtractor)
	}

	e.Register("kms", "CreateKey", func(params map[string]interface{}, accountID, region string) string {
		return kmsKeyArn(accountID, region, "*")
	})

	for _, op := range []string{
		"DescribeKey", "ScheduleKeyDeletion", "CancelKeyDeletion",
		"EnableKey", "DisableKey", "GetKeyPolicy", "PutKeyPolicy",
	} {
		e.Register("kms", op, kmsExtractor)
	}

	e.Register("kms", "ListKeys", func(params map[string]interface{}, accountID, region string) string {
		return kmsKeyArn(accountID, region, "*")
	})

	kmsAliasExtractor := func(params map[string]interface{}, accountID, region string) string {
		aliasName, _ := params["AliasName"].(string)
		if aliasName == "" {
			return "*"
		}
		return kmsAliasArn(accountID, region, aliasName)
	}

	for _, op := range []string{"CreateAlias", "DeleteAlias"} {
		e.Register("kms", op, kmsAliasExtractor)
	}
}

func (e *ResourceExtractor) registerLambdaExtractors() {
	lambdaExtractor := func(params map[string]interface{}, accountID, region string) string {
		functionName := extractLambdaFunctionName(params)
		if functionName == "" {
			return "*"
		}
		if strings.HasPrefix(functionName, "arn:") {
			return functionName
		}
		return lambdaFunctionArn(accountID, region, functionName)
	}

	operations := []string{
		"Invoke", "InvokeAsync", "GetFunction", "UpdateFunctionCode", "UpdateFunctionConfiguration",
		"DeleteFunction", "GetFunctionConfiguration", "PublishVersion", "CreateAlias", "UpdateAlias",
		"DeleteAlias", "GetPolicy", "AddPermission", "RemovePermission",
	}
	for _, op := range operations {
		e.Register("lambda", op, lambdaExtractor)
	}

	e.Register("lambda", "CreateFunction", func(params map[string]interface{}, accountID, region string) string {
		functionName, _ := params["FunctionName"].(string)
		if functionName == "" {
			return "*"
		}
		return lambdaFunctionArn(accountID, region, functionName)
	})

	e.Register("lambda", "ListFunctions", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, region).Build("lambda", "function:*")
	})

	e.Register("lambda", "CreateEventSourceMapping", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, region).Build("lambda", "event-source-mapping:*")
	})

	e.Register("lambda", "ListEventSourceMappings", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, region).Build("lambda", "event-source-mapping:*")
	})
}

func (e *ResourceExtractor) registerSNExtractors() {
	snsExtractor := func(params map[string]interface{}, accountID, region string) string {
		topicArn, _ := params["TopicArn"].(string)
		if topicArn != "" {
			return topicArn
		}
		topicArn, _ = params["TargetArn"].(string)
		if topicArn != "" {
			return topicArn
		}
		return "*"
	}

	operations := []string{
		"Publish", "PublishBatch", "GetTopicAttributes", "SetTopicAttributes",
		"AddPermission", "RemovePermission", "DeleteTopic",
	}
	for _, op := range operations {
		e.Register("sns", op, snsExtractor)
	}

	e.Register("sns", "CreateTopic", func(params map[string]interface{}, accountID, region string) string {
		name, _ := params["Name"].(string)
		if name == "" {
			return "*"
		}
		return snsTopicArn(accountID, region, name)
	})

	e.Register("sns", "ListTopics", func(params map[string]interface{}, accountID, region string) string {
		return snsTopicArn(accountID, region, "*")
	})

	e.Register("sns", "Subscribe", func(params map[string]interface{}, accountID, region string) string {
		topicArn, _ := params["TopicArn"].(string)
		if topicArn == "" {
			return "*"
		}
		return topicArn
	})
}

func (e *ResourceExtractor) registerEventBridgeExtractors() {
	eventsExtractor := func(params map[string]interface{}, accountID, region string) string {
		eventBusArn, _ := params["EventBusArn"].(string)
		if eventBusArn != "" {
			return eventBusArn
		}
		eventBusName, _ := params["EventBusName"].(string)
		if eventBusName == "" {
			eventBusName = "default"
		}
		return eventsEventBusArn(accountID, region, eventBusName)
	}

	operations := []string{
		"PutEvents", "PutRule", "DeleteRule", "EnableRule", "DisableRule",
		"DescribeRule", "ListRules", "PutTargets", "RemoveTargets",
	}
	for _, op := range operations {
		e.Register("events", op, eventsExtractor)
	}

	eventsBusExtractor := func(params map[string]interface{}, accountID, region string) string {
		name, _ := params["Name"].(string)
		if name == "" {
			return "*"
		}
		return eventsEventBusArn(accountID, region, name)
	}

	for _, op := range []string{"CreateEventBus", "DeleteEventBus"} {
		e.Register("events", op, eventsBusExtractor)
	}

	e.Register("events", "ListEventBuses", func(params map[string]interface{}, accountID, region string) string {
		return eventsEventBusArn(accountID, region, "*")
	})
}

func (e *ResourceExtractor) registerIAMExtractors() {
	iamUserExtractor := func(params map[string]interface{}, accountID, region string) string {
		userName, _ := params["UserName"].(string)
		if userName == "" {
			return "*"
		}
		return iamUserArn(accountID, region, userName)
	}

	for _, op := range []string{"CreateUser", "GetUser", "DeleteUser", "CreateAccessKey"} {
		e.Register("iam", op, iamUserExtractor)
	}

	iamRoleExtractor := func(params map[string]interface{}, accountID, region string) string {
		roleName, _ := params["RoleName"].(string)
		if roleName == "" {
			return "*"
		}
		return iamRoleArn(accountID, region, roleName)
	}

	for _, op := range []string{"CreateRole", "GetRole", "DeleteRole"} {
		e.Register("iam", op, iamRoleExtractor)
	}

	iamPolicyExtractor := func(params map[string]interface{}, accountID, region string) string {
		policyName, _ := params["PolicyName"].(string)
		if policyName == "" {
			return "*"
		}
		return iamPolicyArn(accountID, region, policyName)
	}

	e.Register("iam", "CreatePolicy", iamPolicyExtractor)

	e.Register("iam", "GetPolicy", func(params map[string]interface{}, accountID, region string) string {
		policyArn, _ := params["PolicyArn"].(string)
		if policyArn == "" {
			return "*"
		}
		return policyArn
	})

	e.Register("iam", "ListUsers", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, "").Build("iam", "user/*")
	})
	e.Register("iam", "ListRoles", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, "").Build("iam", "role/*")
	})
	e.Register("iam", "ListPolicies", func(params map[string]interface{}, accountID, region string) string {
		return arnutil.NewARNBuilder(accountID, "").Build("iam", "policy/*")
	})
}

func (e *ResourceExtractor) registerKinesisExtractors() {
	kinesisExtractor := func(params map[string]interface{}, accountID, region string) string {
		streamName := extractKinesisStreamName(params)
		if streamName == "" {
			return "*"
		}
		return kinesisStreamArn(accountID, region, streamName)
	}

	operations := []string{
		"PutRecord", "PutRecords", "GetRecords", "GetShardIterator",
		"DescribeStream", "DescribeStreamSummary", "ListShards",
		"EnableEnhancedMonitoring", "DisableEnhancedMonitoring",
	}
	for _, op := range operations {
		e.Register("kinesis", op, kinesisExtractor)
	}

	e.Register("kinesis", "CreateStream", func(params map[string]interface{}, accountID, region string) string {
		streamName, _ := params["StreamName"].(string)
		if streamName == "" {
			return "*"
		}
		return kinesisStreamArn(accountID, region, streamName)
	})

	e.Register("kinesis", "DeleteStream", func(params map[string]interface{}, accountID, region string) string {
		streamName, _ := params["StreamName"].(string)
		if streamName == "" {
			return "*"
		}
		return kinesisStreamArn(accountID, region, streamName)
	})

	e.Register("kinesis", "ListStreams", func(params map[string]interface{}, accountID, region string) string {
		return kinesisStreamArn(accountID, region, "*")
	})
}

func (e *ResourceExtractor) registerCloudWatchExtractors() {
	e.Register("logs", "*", func(params map[string]interface{}, accountID, region string) string {
		logGroupName, _ := params["LogGroupName"].(string)
		if logGroupName == "" {
			return "*"
		}
		return logsLogGroupArn(accountID, region, logGroupName)
	})

	e.Register("cloudwatch", "*", func(params map[string]interface{}, accountID, region string) string {
		namespace, _ := params["Namespace"].(string)
		if namespace == "" {
			return "*"
		}
		return arnutil.NewARNBuilder(accountID, region).Build("cloudwatch", namespace)
	})
}

func (e *ResourceExtractor) registerStepFunctionsExtractors() {
	sfnExtractor := func(params map[string]interface{}, accountID, region string) string {
		stateMachineArn, _ := params["StateMachineArn"].(string)
		if stateMachineArn != "" {
			return stateMachineArn
		}
		return "*"
	}

	operations := []string{
		"StartExecution", "StopExecution", "DescribeExecution", "GetExecutionHistory",
		"DescribeStateMachine", "UpdateStateMachine", "DeleteStateMachine",
		"ListExecutions", "ListStateMachines",
	}
	for _, op := range operations {
		e.Register("states", op, sfnExtractor)
	}

	e.Register("states", "CreateStateMachine", func(params map[string]interface{}, accountID, region string) string {
		name, _ := params["Name"].(string)
		if name == "" {
			return "*"
		}
		return statesStateMachineArn(accountID, region, name)
	})
}

func (e *ResourceExtractor) registerAPIGatewayExtractors() {
	e.Register("apigateway", "*", func(params map[string]interface{}, accountID, region string) string {
		restApiId, _ := params["RestApiId"].(string)
		if restApiId == "" {
			return "*"
		}
		return arnutil.NewARNBuilder(accountID, region).APIGateway().RestApi(restApiId)
	})
}

func (e *ResourceExtractor) registerAthenaExtractors() {
	e.Register("athena", "*", func(params map[string]interface{}, accountID, region string) string {
		workGroup, _ := params["WorkGroup"].(string)
		if workGroup == "" {
			workGroup = "primary"
		}
		return athenaWorkGroupArn(accountID, region, workGroup)
	})
}

func (e *ResourceExtractor) registerSecretsManagerExtractors() {
	e.Register("secretsmanager", "*", func(params map[string]interface{}, accountID, region string) string {
		secretId, _ := params["SecretId"].(string)
		if secretId == "" {
			return "*"
		}
		if strings.HasPrefix(secretId, "arn:") {
			return secretId
		}
		return secretsManagerSecretArn(accountID, region, secretId)
	})
}

func extractQueueNameFromURL(queueURL string) string {
	if queueURL == "" {
		return ""
	}
	parsed, err := url.Parse(queueURL)
	if err != nil {
		parts := strings.Split(queueURL, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return ""
	}
	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	if len(parts) == 1 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

func extractKMSKeyID(params map[string]interface{}) string {
	if keyID, ok := params["KeyId"].(string); ok && keyID != "" {
		return keyID
	}
	if keyID, ok := params["KeyID"].(string); ok && keyID != "" {
		return keyID
	}
	if grantToken, ok := params["GrantTokens"].([]interface{}); ok && len(grantToken) > 0 {
		if token, ok := grantToken[0].(string); ok && strings.HasPrefix(token, "arn:") {
			return token
		}
	}
	return ""
}

func extractLambdaFunctionName(params map[string]interface{}) string {
	if fn, ok := params["FunctionName"].(string); ok && fn != "" {
		return fn
	}
	return ""
}

func extractKinesisStreamName(params map[string]interface{}) string {
	if streamName, ok := params["StreamName"].(string); ok && streamName != "" {
		return streamName
	}
	if streamARN, ok := params["StreamARN"].(string); ok && streamARN != "" {
		if strings.HasPrefix(streamARN, "arn:") {
			parts := strings.Split(streamARN, "/")
			if len(parts) > 1 {
				return parts[len(parts)-1]
			}
		}
	}
	return ""
}

// registerAppSyncExtractors registers resource ARN extractors for the AWS AppSync
// service. All AppSync resources share the ARN pattern
// arn:aws:appsync:{region}:{account}:apis/{apiId}.
func (e *ResourceExtractor) registerAppSyncExtractors() {
	e.Register("appsync", "*", func(params map[string]interface{}, accountID, region string) string {
		apiId, _ := params["apiId"].(string)
		if apiId == "" {
			return "*"
		}
		return appsyncApiArn(accountID, region, apiId)
	})
}

// registerNeptuneExtractors registers resource ARN extractors for the AWS Neptune
// service. Neptune resources use ARN patterns such as
// arn:aws:neptune:{region}:{account}:cluster/{clusterId}.
func (e *ResourceExtractor) registerNeptuneExtractors() {
	clusterExtractor := func(params map[string]interface{}, accountID, region string) string {
		clusterId, _ := params["DBClusterIdentifier"].(string)
		if clusterId == "" {
			return "*"
		}
		return neptuneClusterArn(accountID, region, clusterId)
	}

	for _, op := range []string{
		"CreateDBCluster", "DeleteDBCluster", "ModifyDBCluster", "DescribeDBClusters",
		"StartDBCluster", "StopDBCluster", "FailoverDBCluster",
		"AddRoleToDBCluster", "RemoveRoleFromDBCluster",
		"RestoreDBClusterFromSnapshot", "RestoreDBClusterToPointInTime",
		"PromoteReadReplicaDBCluster",
		"CreateDBClusterEndpoint", "DescribeDBClusterEndpoints",
		"ModifyDBClusterEndpoint", "DeleteDBClusterEndpoint",
	} {
		e.Register("neptune", op, clusterExtractor)
	}

	snapshotExtractor := func(params map[string]interface{}, accountID, region string) string {
		snapshotId, _ := params["DBClusterSnapshotIdentifier"].(string)
		if snapshotId == "" {
			return "*"
		}
		return neptuneClusterSnapshotArn(accountID, region, snapshotId)
	}

	for _, op := range []string{
		"CreateDBClusterSnapshot", "DeleteDBClusterSnapshot", "DescribeDBClusterSnapshots",
		"CopyDBClusterSnapshot",
		"DescribeDBClusterSnapshotAttributes", "ModifyDBClusterSnapshotAttribute",
	} {
		e.Register("neptune", op, snapshotExtractor)
	}

	paramGroupExtractor := func(params map[string]interface{}, accountID, region string) string {
		groupName, _ := params["DBClusterParameterGroupName"].(string)
		if groupName == "" {
			groupName, _ = params["DBParameterGroupName"].(string)
		}
		if groupName == "" {
			return "*"
		}
		return neptuneParamGroupArn(accountID, region, groupName)
	}

	for _, op := range []string{
		"CreateDBClusterParameterGroup", "DeleteDBClusterParameterGroup",
		"DescribeDBClusterParameterGroups", "DescribeDBClusterParameters",
		"ModifyDBClusterParameterGroup", "ResetDBClusterParameterGroup",
		"CopyDBClusterParameterGroup",
		"CreateDBParameterGroup", "DeleteDBParameterGroup",
		"DescribeDBParameterGroups", "DescribeDBParameters",
		"ModifyDBParameterGroup", "ResetDBParameterGroup",
		"CopyDBParameterGroup",
	} {
		e.Register("neptune", op, paramGroupExtractor)
	}

	subnetGroupExtractor := func(params map[string]interface{}, accountID, region string) string {
		groupName, _ := params["DBSubnetGroupName"].(string)
		if groupName == "" {
			return "*"
		}
		return neptuneSubnetGroupArn(accountID, region, groupName)
	}

	for _, op := range []string{
		"CreateDBSubnetGroup", "DeleteDBSubnetGroup",
		"DescribeDBSubnetGroups", "ModifyDBSubnetGroup",
	} {
		e.Register("neptune", op, subnetGroupExtractor)
	}

	instanceExtractor := func(params map[string]interface{}, accountID, region string) string {
		instanceId, _ := params["DBInstanceIdentifier"].(string)
		if instanceId == "" {
			return "*"
		}
		return neptuneInstanceArn(accountID, region, instanceId)
	}

	for _, op := range []string{
		"CreateDBInstance", "DeleteDBInstance", "ModifyDBInstance",
		"DescribeDBInstances", "RebootDBInstance",
	} {
		e.Register("neptune", op, instanceExtractor)
	}

	globalClusterExtractor := func(params map[string]interface{}, accountID, region string) string {
		clusterId, _ := params["GlobalClusterIdentifier"].(string)
		if clusterId == "" {
			return "*"
		}
		return neptuneGlobalClusterArn(accountID, region, clusterId)
	}

	for _, op := range []string{
		"CreateGlobalCluster", "DeleteGlobalCluster", "DescribeGlobalClusters",
		"ModifyGlobalCluster", "RemoveFromGlobalCluster",
	} {
		e.Register("neptune", op, globalClusterExtractor)
	}

	eventSubExtractor := func(params map[string]interface{}, accountID, region string) string {
		subName, _ := params["SubscriptionName"].(string)
		if subName == "" {
			return "*"
		}
		return neptuneEventSubscriptionArn(accountID, region, subName)
	}

	for _, op := range []string{
		"CreateEventSubscription", "DeleteEventSubscription",
		"DescribeEventSubscriptions", "ModifyEventSubscription",
		"AddSourceIdentifierToSubscription", "RemoveSourceIdentifierFromSubscription",
	} {
		e.Register("neptune", op, eventSubExtractor)
	}

	tagExtractor := func(params map[string]interface{}, accountID, region string) string {
		resourceName, _ := params["ResourceName"].(string)
		if resourceName == "" {
			return "*"
		}
		return resourceName
	}

	for _, op := range []string{
		"AddTagsToResource", "ListTagsForResource", "RemoveTagsFromResource",
	} {
		e.Register("neptune", op, tagExtractor)
	}
}

func (e *ResourceExtractor) registerNeptuneGraphExtractors() {
	graphExtractor := func(params map[string]interface{}, accountID, region string) string {
		graphId, _ := params["graphIdentifier"].(string)
		if graphId == "" {
			return "*"
		}
		return neptuneGraphArn(accountID, region, graphId)
	}

	for _, op := range []string{
		"CreateGraph", "DeleteGraph", "GetGraph", "ListGraphs", "UpdateGraph",
		"StartGraph", "StopGraph", "ResetGraph",
		"RestoreGraphFromSnapshot",
		"CreatePrivateGraphEndpoint", "GetPrivateGraphEndpoint",
		"ListPrivateGraphEndpoints", "DeletePrivateGraphEndpoint",
		"ExecuteQuery", "GetGraphSummary",
	} {
		e.Register("neptunegraph", op, graphExtractor)
	}

	snapshotExtractor := func(params map[string]interface{}, accountID, region string) string {
		snapshotId, _ := params["snapshotIdentifier"].(string)
		if snapshotId == "" {
			return "*"
		}
		return neptuneGraphSnapshotArn(accountID, region, snapshotId)
	}

	for _, op := range []string{
		"CreateGraphSnapshot", "DeleteGraphSnapshot", "GetGraphSnapshot", "ListGraphSnapshots",
	} {
		e.Register("neptunegraph", op, snapshotExtractor)
	}

	importTaskExtractor := func(params map[string]interface{}, accountID, region string) string {
		taskId, _ := params["taskIdentifier"].(string)
		if taskId == "" {
			return "*"
		}
		return neptuneGraphImportTaskArn(accountID, region, taskId)
	}

	for _, op := range []string{
		"CreateGraphUsingImportTask", "GetImportTask", "ListImportTasks",
		"CancelImportTask", "StartImportTask",
	} {
		e.Register("neptunegraph", op, importTaskExtractor)
	}

	exportTaskExtractor := func(params map[string]interface{}, accountID, region string) string {
		taskId, _ := params["taskIdentifier"].(string)
		if taskId == "" {
			return "*"
		}
		return neptuneGraphExportTaskArn(accountID, region, taskId)
	}

	for _, op := range []string{
		"StartExportTask", "GetExportTask", "ListExportTasks", "CancelExportTask",
	} {
		e.Register("neptunegraph", op, exportTaskExtractor)
	}

	ngTagExtractor := func(params map[string]interface{}, accountID, region string) string {
		resourceArn, _ := params["resourceArn"].(string)
		if resourceArn == "" {
			return "*"
		}
		return resourceArn
	}

	for _, op := range []string{
		"ListTagsForResource", "TagResource", "UntagResource",
	} {
		e.Register("neptunegraph", op, ngTagExtractor)
	}

	for _, op := range []string{
		"GetQuery", "ListQueries", "CancelQuery",
	} {
		e.Register("neptunegraph", op, graphExtractor)
	}
}
