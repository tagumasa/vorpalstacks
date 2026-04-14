package authorization

import (
	"fmt"
	"net/url"
	"strings"
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

func (e *ResourceExtractor) registerS3Extractors() {
	s3Extractor := func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		key, _ := params["Key"].(string)
		if key == "" {
			return fmt.Sprintf("arn:aws:s3:::%s", bucket)
		}
		return fmt.Sprintf("arn:aws:s3:::%s/%s", bucket, key)
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

	e.Register("s3", "ListBucket", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "CreateBucket", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "DeleteBucket", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "GetBucketLocation", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "PutBucketPolicy", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "GetBucketPolicy", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})

	e.Register("s3", "DeleteBucketPolicy", func(params map[string]interface{}, accountID, region string) string {
		bucket, _ := params["Bucket"].(string)
		if bucket == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:s3:::%s", bucket)
	})
}

func (e *ResourceExtractor) registerSQSExtractors() {
	sqsExtractor := func(params map[string]interface{}, accountID, region string) string {
		queueURL, _ := params["QueueUrl"].(string)
		if queueURL == "" {
			return "*"
		}
		queueName := extractQueueNameFromURL(queueURL)
		return fmt.Sprintf("arn:aws:sqs:%s:%s:%s", region, accountID, queueName)
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
		return fmt.Sprintf("arn:aws:sqs:%s:%s:%s", region, accountID, queueName)
	})

	e.Register("sqs", "DeleteQueue", func(params map[string]interface{}, accountID, region string) string {
		queueURL, _ := params["QueueUrl"].(string)
		if queueURL == "" {
			return "*"
		}
		queueName := extractQueueNameFromURL(queueURL)
		return fmt.Sprintf("arn:aws:sqs:%s:%s:%s", region, accountID, queueName)
	})

	e.Register("sqs", "ListQueues", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:sqs:%s:%s:*", region, accountID)
	})
}

func (e *ResourceExtractor) registerDynamoDBExtractors() {
	dynamoExtractor := func(params map[string]interface{}, accountID, region string) string {
		tableName, _ := params["TableName"].(string)
		if tableName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:dynamodb:%s:%s:table/%s", region, accountID, tableName)
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
		return fmt.Sprintf("arn:aws:dynamodb:%s:%s:table/%s", region, accountID, tableName)
	})

	e.Register("dynamodb", "ListTables", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:dynamodb:%s:%s:table/*", region, accountID)
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
		if strings.HasPrefix(keyID, "mrk-") {
			return fmt.Sprintf("arn:aws:kms:%s:%s:key/%s", region, accountID, keyID)
		}
		if len(keyID) == 36 && strings.Contains(keyID, "-") {
			return fmt.Sprintf("arn:aws:kms:%s:%s:key/%s", region, accountID, keyID)
		}
		return fmt.Sprintf("arn:aws:kms:%s:%s:key/%s", region, accountID, keyID)
	}

	operations := []string{
		"Encrypt", "Decrypt", "ReEncrypt", "GenerateDataKey", "GenerateDataKeyWithoutPlaintext",
		"GenerateRandom", "Sign", "Verify", "GetPublicKey",
	}
	for _, op := range operations {
		e.Register("kms", op, kmsExtractor)
	}

	e.Register("kms", "CreateKey", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:kms:%s:%s:key/*", region, accountID)
	})

	e.Register("kms", "DescribeKey", kmsExtractor)
	e.Register("kms", "ScheduleKeyDeletion", kmsExtractor)
	e.Register("kms", "CancelKeyDeletion", kmsExtractor)
	e.Register("kms", "EnableKey", kmsExtractor)
	e.Register("kms", "DisableKey", kmsExtractor)
	e.Register("kms", "GetKeyPolicy", kmsExtractor)
	e.Register("kms", "PutKeyPolicy", kmsExtractor)
	e.Register("kms", "ListKeys", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:kms:%s:%s:key/*", region, accountID)
	})

	e.Register("kms", "CreateAlias", func(params map[string]interface{}, accountID, region string) string {
		aliasName, _ := params["AliasName"].(string)
		if aliasName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:kms:%s:%s:%s", region, accountID, aliasName)
	})

	e.Register("kms", "DeleteAlias", func(params map[string]interface{}, accountID, region string) string {
		aliasName, _ := params["AliasName"].(string)
		if aliasName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:kms:%s:%s:%s", region, accountID, aliasName)
	})
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
		return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", region, accountID, functionName)
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
		return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", region, accountID, functionName)
	})

	e.Register("lambda", "ListFunctions", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:lambda:%s:%s:function:*", region, accountID)
	})

	e.Register("lambda", "CreateEventSourceMapping", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:lambda:%s:%s:event-source-mapping:*", region, accountID)
	})

	e.Register("lambda", "ListEventSourceMappings", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:lambda:%s:%s:event-source-mapping:*", region, accountID)
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
		return fmt.Sprintf("arn:aws:sns:%s:%s:%s", region, accountID, name)
	})

	e.Register("sns", "ListTopics", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:sns:%s:%s:*", region, accountID)
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
		return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, accountID, eventBusName)
	}

	operations := []string{
		"PutEvents", "PutRule", "DeleteRule", "EnableRule", "DisableRule",
		"DescribeRule", "ListRules", "PutTargets", "RemoveTargets",
	}
	for _, op := range operations {
		e.Register("events", op, eventsExtractor)
	}

	e.Register("events", "CreateEventBus", func(params map[string]interface{}, accountID, region string) string {
		name, _ := params["Name"].(string)
		if name == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, accountID, name)
	})

	e.Register("events", "DeleteEventBus", func(params map[string]interface{}, accountID, region string) string {
		name, _ := params["Name"].(string)
		if name == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, accountID, name)
	})

	e.Register("events", "ListEventBuses", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/*", region, accountID)
	})
}

func (e *ResourceExtractor) registerIAMExtractors() {
	e.Register("iam", "CreateUser", func(params map[string]interface{}, accountID, region string) string {
		userName, _ := params["UserName"].(string)
		if userName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, userName)
	})

	e.Register("iam", "GetUser", func(params map[string]interface{}, accountID, region string) string {
		userName, _ := params["UserName"].(string)
		if userName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, userName)
	})

	e.Register("iam", "DeleteUser", func(params map[string]interface{}, accountID, region string) string {
		userName, _ := params["UserName"].(string)
		if userName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, userName)
	})

	e.Register("iam", "CreateRole", func(params map[string]interface{}, accountID, region string) string {
		roleName, _ := params["RoleName"].(string)
		if roleName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
	})

	e.Register("iam", "GetRole", func(params map[string]interface{}, accountID, region string) string {
		roleName, _ := params["RoleName"].(string)
		if roleName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
	})

	e.Register("iam", "DeleteRole", func(params map[string]interface{}, accountID, region string) string {
		roleName, _ := params["RoleName"].(string)
		if roleName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
	})

	e.Register("iam", "CreatePolicy", func(params map[string]interface{}, accountID, region string) string {
		policyName, _ := params["PolicyName"].(string)
		if policyName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountID, policyName)
	})

	e.Register("iam", "GetPolicy", func(params map[string]interface{}, accountID, region string) string {
		policyArn, _ := params["PolicyArn"].(string)
		if policyArn == "" {
			return "*"
		}
		return policyArn
	})

	e.Register("iam", "CreateAccessKey", func(params map[string]interface{}, accountID, region string) string {
		userName, _ := params["UserName"].(string)
		if userName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, userName)
	})

	e.Register("iam", "ListUsers", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:iam::%s:user/*", accountID)
	})

	e.Register("iam", "ListRoles", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:iam::%s:role/*", accountID)
	})

	e.Register("iam", "ListPolicies", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:iam::%s:policy/*", accountID)
	})
}

func (e *ResourceExtractor) registerKinesisExtractors() {
	kinesisExtractor := func(params map[string]interface{}, accountID, region string) string {
		streamName := extractKinesisStreamName(params)
		if streamName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/%s", region, accountID, streamName)
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
		return fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/%s", region, accountID, streamName)
	})

	e.Register("kinesis", "DeleteStream", func(params map[string]interface{}, accountID, region string) string {
		streamName, _ := params["StreamName"].(string)
		if streamName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/%s", region, accountID, streamName)
	})

	e.Register("kinesis", "ListStreams", func(params map[string]interface{}, accountID, region string) string {
		return fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/*", region, accountID)
	})
}

func (e *ResourceExtractor) registerCloudWatchExtractors() {
	e.Register("logs", "*", func(params map[string]interface{}, accountID, region string) string {
		logGroupName, _ := params["LogGroupName"].(string)
		if logGroupName == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", region, accountID, logGroupName)
	})

	e.Register("cloudwatch", "*", func(params map[string]interface{}, accountID, region string) string {
		namespace, _ := params["Namespace"].(string)
		if namespace == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:cloudwatch:%s:%s:%s", region, accountID, namespace)
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
		return fmt.Sprintf("arn:aws:states:%s:%s:stateMachine:%s", region, accountID, name)
	})
}

func (e *ResourceExtractor) registerAPIGatewayExtractors() {
	e.Register("apigateway", "*", func(params map[string]interface{}, accountID, region string) string {
		restApiId, _ := params["RestApiId"].(string)
		if restApiId == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s", region, restApiId)
	})
}

func (e *ResourceExtractor) registerAthenaExtractors() {
	e.Register("athena", "*", func(params map[string]interface{}, accountID, region string) string {
		workGroup, _ := params["WorkGroup"].(string)
		if workGroup == "" {
			workGroup = "primary"
		}
		return fmt.Sprintf("arn:aws:athena:%s:%s:workgroup/%s", region, accountID, workGroup)
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
		return fmt.Sprintf("arn:aws:secretsmanager:%s:%s:secret:%s", region, accountID, secretId)
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
		return fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", region, accountID, apiId)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:cluster/%s", region, accountID, clusterId)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:cluster-snapshot/%s", region, accountID, snapshotId)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:cluster-pg/%s", region, accountID, groupName)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:subnet-group/%s", region, accountID, groupName)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:db/%s", region, accountID, instanceId)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:global-cluster/%s", region, accountID, clusterId)
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
		return fmt.Sprintf("arn:aws:neptune:%s:%s:event-subscription/%s", region, accountID, subName)
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
		return fmt.Sprintf("arn:aws:neptune-graph:%s:%s:graph/%s", region, accountID, graphId)
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
		return fmt.Sprintf("arn:aws:neptune-graph:%s:%s:snapshot/%s", region, accountID, snapshotId)
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
		return fmt.Sprintf("arn:aws:neptune-graph:%s:%s:import-task/%s", region, accountID, taskId)
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
		return fmt.Sprintf("arn:aws:neptune-graph:%s:%s:export-task/%s", region, accountID, taskId)
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

	queryExtractor := func(params map[string]interface{}, accountID, region string) string {
		graphId, _ := params["graphIdentifier"].(string)
		if graphId == "" {
			return "*"
		}
		return fmt.Sprintf("arn:aws:neptune-graph:%s:%s:graph/%s", region, accountID, graphId)
	}

	for _, op := range []string{
		"GetQuery", "ListQueries", "CancelQuery",
	} {
		e.Register("neptunegraph", op, queryExtractor)
	}
}
