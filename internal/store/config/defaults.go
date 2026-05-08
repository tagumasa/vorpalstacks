package config

import (
	"strconv"

	"vorpalstacks/internal/common/serviceports"
)

func loadDefaults() map[string]ConfigEntry {
	return map[string]ConfigEntry{
		// Server Configuration
		"server.port": {
			Key:         "server.port",
			Value:       serviceports.HTTP,
			Type:        ConfigTypePort,
			Description: "Main HTTP server port",
			ReadOnly:    false,
			EnvVar:      "PORT",
			Category:    CategoryServer,
		},
		"server.grpc_web_port": {
			Key:         "server.grpc_web_port",
			Value:       serviceports.GRPCWeb,
			Type:        ConfigTypePort,
			Description: "gRPC-Web admin port",
			ReadOnly:    false,
			EnvVar:      "GRPC_WEB_PORT",
			Category:    CategoryServer,
		},
		"server.tls_enabled": {
			Key:         "server.tls_enabled",
			Value:       false,
			Type:        ConfigTypeBool,
			Description: "Enable TLS for HTTPS server",
			ReadOnly:    false,
			EnvVar:      "TLS_ENABLED",
			Category:    CategoryServer,
		},
		"server.tls_port": {
			Key:         "server.tls_port",
			Value:       serviceports.TLS,
			Type:        ConfigTypePort,
			Description: "HTTPS server port",
			ReadOnly:    false,
			EnvVar:      "TLS_PORT",
			Category:    CategoryServer,
		},
		"server.tls_cert_path": {
			Key:         "server.tls_cert_path",
			Value:       "auto",
			Type:        ConfigTypeString,
			Description: "TLS certificate path ('auto' generates and saves to pebble)",
			ReadOnly:    false,
			EnvVar:      "TLS_CERT_PATH",
			Category:    CategoryServer,
		},
		"server.tls_key_path": {
			Key:         "server.tls_key_path",
			Value:       "auto",
			Type:        ConfigTypeString,
			Description: "TLS private key path ('auto' generates and saves to pebble)",
			ReadOnly:    false,
			EnvVar:      "TLS_KEY_PATH",
			Category:    CategoryServer,
		},
		"server.bind_mode": {
			Key:         "server.bind_mode",
			Value:       "all",
			Type:        ConfigTypeString,
			Description: "Console bind mode: all (0.0.0.0, default), localhost (127.0.0.1), or interface (custom IP)",
			ReadOnly:    false,
			EnvVar:      "BIND_MODE",
			Category:    CategoryServer,
		},
		"server.bind_interface": {
			Key:         "server.bind_interface",
			Value:       "",
			Type:        ConfigTypeString,
			Description: "Custom bind IP address (used only when bind_mode=interface)",
			ReadOnly:    false,
			EnvVar:      "BIND_INTERFACE",
			Category:    CategoryServer,
		},

		// AWS Identity (Read-only at runtime)
		"aws.account_id": {
			Key:         "aws.account_id",
			Value:       "",
			Type:        ConfigTypeString,
			Description: "AWS account ID",
			ReadOnly:    true,
			EnvVar:      "AWS_ACCOUNT_ID",
			Category:    CategoryAWS,
		},
		"aws.region": {
			Key:         "aws.region",
			Value:       "us-east-1",
			Type:        ConfigTypeString,
			Description: "Default AWS region",
			ReadOnly:    true,
			EnvVar:      "AWS_REGION",
			Category:    CategoryAWS,
		},

		// Storage Configuration
		"storage.data_path": {
			Key:         "storage.data_path",
			Value:       "./data",
			Type:        ConfigTypeString,
			Description: "Data storage path",
			ReadOnly:    true,
			EnvVar:      "DATA_PATH",
			Category:    CategoryStorage,
		},
		"storage.metadata_path": {
			Key:         "storage.metadata_path",
			Value:       "",
			Type:        ConfigTypeString,
			Description: "Metadata path (defaults to data_path if empty)",
			ReadOnly:    true,
			EnvVar:      "METADATA_PATH",
			Category:    CategoryStorage,
		},

		// Feature Flags
		"features.audit_enabled": {
			Key:         "features.audit_enabled",
			Value:       false,
			Type:        ConfigTypeBool,
			Description: "CloudTrail audit logging",
			ReadOnly:    false,
			EnvVar:      "VS_AUDIT_ENABLED",
			Category:    CategoryFeatures,
		},
		"features.signature_verification": {
			Key:         "features.signature_verification",
			Value:       true,
			Type:        ConfigTypeBool,
			Description: "AWS signature verification",
			ReadOnly:    false,
			EnvVar:      "SIGNATURE_VERIFICATION_ENABLED",
			Category:    CategoryFeatures,
		},
		"features.route53_dns": {
			Key:         "features.route53_dns",
			Value:       false,
			Type:        ConfigTypeBool,
			Description: "Route53 DNS server",
			ReadOnly:    false,
			EnvVar:      "ROUTE53_DNS_ENABLED",
			Category:    CategoryFeatures,
		},

		// External Service Endpoints
		"endpoints.base_url": {
			Key:         "endpoints.base_url",
			Value:       "http://localhost:" + strconv.Itoa(serviceports.HTTP),
			Type:        ConfigTypeURL,
			Description: "Base URL for generated endpoints",
			ReadOnly:    false,
			EnvVar:      "VS_BASE_URL",
			Category:    CategoryEndpoints,
		},
		"endpoints.s3_website_suffix": {
			Key:         "endpoints.s3_website_suffix",
			Value:       "s3-website.{region}.amazonaws.com",
			Type:        ConfigTypeString,
			Description: "S3 website domain suffix",
			ReadOnly:    false,
			EnvVar:      "VS_S3_WEBSITE_SUFFIX",
			Category:    CategoryEndpoints,
		},
		"endpoints.cognito_suffix": {
			Key:         "endpoints.cognito_suffix",
			Value:       "auth.{region}.amazoncognito.com",
			Type:        ConfigTypeString,
			Description: "Cognito hosted UI suffix",
			ReadOnly:    false,
			EnvVar:      "VS_COGNITO_SUFFIX",
			Category:    CategoryEndpoints,
		},
		"endpoints.apigateway_suffix": {
			Key:         "endpoints.apigateway_suffix",
			Value:       "execute-api.{region}.amazonaws.com",
			Type:        ConfigTypeString,
			Description: "API Gateway suffix",
			ReadOnly:    false,
			EnvVar:      "VS_APIGW_SUFFIX",
			Category:    CategoryEndpoints,
		},

		// Default Public Endpoint Ports
		"ports.s3_website": {
			Key:         "ports.s3_website",
			Value:       serviceports.S3Website,
			Type:        ConfigTypePort,
			Description: "S3 website default port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_S3_WEBSITE",
			Category:    CategoryPorts,
		},
		"ports.apigateway": {
			Key:         "ports.apigateway",
			Value:       serviceports.APIGateway,
			Type:        ConfigTypePort,
			Description: "API Gateway invoke URL port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_APIGW",
			Category:    CategoryPorts,
		},
		"ports.cognito_hosted": {
			Key:         "ports.cognito_hosted",
			Value:       serviceports.Cognito,
			Type:        ConfigTypePort,
			Description: "Cognito Hosted UI port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_COGNITO_HOSTED",
			Category:    CategoryPorts,
		},
		"ports.cloudfront": {
			Key:         "ports.cloudfront",
			Value:       serviceports.CloudFront,
			Type:        ConfigTypePort,
			Description: "CloudFront distribution port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_CLOUDFRONT",
			Category:    CategoryPorts,
		},
		"ports.lambda_url": {
			Key:         "ports.lambda_url",
			Value:       serviceports.LambdaURL,
			Type:        ConfigTypePort,
			Description: "Lambda Function URL port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_LAMBDA_URL",
			Category:    CategoryPorts,
		},
		"ports.appsync_events": {
			Key:         "ports.appsync_events",
			Value:       serviceports.AppSync,
			Type:        ConfigTypePort,
			Description: "AppSync Events (WebSocket + HTTP publish) port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_APPSYNC_EVENTS",
			Category:    CategoryPorts,
		},
		"ports.route53_dns": {
			Key:         "ports.route53_dns",
			Value:       serviceports.Route53DNS,
			Type:        ConfigTypePort,
			Description: "Route53 DNS server port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_ROUTE53_DNS",
			Category:    CategoryPorts,
		},
		"ports.route53_healthcheck": {
			Key:         "ports.route53_healthcheck",
			Value:       serviceports.Route53HC,
			Type:        ConfigTypePort,
			Description: "Route53 health check default endpoint port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_ROUTE53_HEALTHCHECK",
			Category:    CategoryPorts,
		},
		"ports.dynamic_range_start": {
			Key:         "ports.dynamic_range_start",
			Value:       serviceports.DynamicRangeStart,
			Type:        ConfigTypeInt,
			Description: "Start of dynamic port allocation range",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_DYNAMIC_START",
			Category:    CategoryPorts,
		},
		"ports.dynamic_range_end": {
			Key:         "ports.dynamic_range_end",
			Value:       serviceports.DynamicRangeEnd,
			Type:        ConfigTypeInt,
			Description: "End of dynamic port allocation range",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_DYNAMIC_END",
			Category:    CategoryPorts,
		},

		// Port mode entries (FQDN or Individual)
		"ports.s3_website.mode":     {Key: "ports.s3_website.mode", Value: "fqdn", Type: ConfigTypeString, Description: "S3 Website port mode (fqdn/individual)", ReadOnly: false, Category: CategoryPorts},
		"ports.apigateway.mode":     {Key: "ports.apigateway.mode", Value: "fqdn", Type: ConfigTypeString, Description: "API Gateway port mode", ReadOnly: false, Category: CategoryPorts},
		"ports.cognito_hosted.mode": {Key: "ports.cognito_hosted.mode", Value: "fqdn", Type: ConfigTypeString, Description: "Cognito Hosted UI port mode", ReadOnly: false, Category: CategoryPorts},
		"ports.cloudfront.mode":     {Key: "ports.cloudfront.mode", Value: "fqdn", Type: ConfigTypeString, Description: "CloudFront port mode", ReadOnly: false, Category: CategoryPorts},
		"ports.lambda_url.mode":     {Key: "ports.lambda_url.mode", Value: "fqdn", Type: ConfigTypeString, Description: "Lambda URL port mode", ReadOnly: false, Category: CategoryPorts},
		"ports.appsync_events.mode": {Key: "ports.appsync_events.mode", Value: "fqdn", Type: ConfigTypeString, Description: "AppSync Events port mode", ReadOnly: false, Category: CategoryPorts},

		// HTTP Configuration (CORS)
		"http.cors_allowed_origins": {
			Key:         "http.cors_allowed_origins",
			Value:       "*",
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS origins ('*' allows all)",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_ORIGINS",
			Category:    CategoryHTTP,
		},
		"http.cors_allowed_methods": {
			Key:         "http.cors_allowed_methods",
			Value:       "GET, POST, PUT, DELETE, OPTIONS, HEAD",
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS methods",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_METHODS",
			Category:    CategoryHTTP,
		},
		"http.cors_allowed_headers": {
			Key:         "http.cors_allowed_headers",
			Value:       "Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256",
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS request headers",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_HEADERS",
			Category:    CategoryHTTP,
		},
		"http.cors_expose_headers": {
			Key:         "http.cors_expose_headers",
			Value:       "x-amzn-RequestId, x-amzn-ErrorType, x-amzn-ErrorMessage",
			Type:        ConfigTypeString,
			Description: "Comma-separated list of CORS headers exposed to the client",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_EXPOSE_HEADERS",
			Category:    CategoryHTTP,
		},

		// Service enable/disable entries
		"services.sns.enabled":              {Key: "services.sns.enabled", Value: true, Type: ConfigTypeBool, Description: "SNS service", ReadOnly: false, EnvVar: "SNS_ENABLED", Category: CategoryServices},
		"services.sqs.enabled":              {Key: "services.sqs.enabled", Value: true, Type: ConfigTypeBool, Description: "SQS service", ReadOnly: false, EnvVar: "SQS_ENABLED", Category: CategoryServices},
		"services.lambda.enabled":           {Key: "services.lambda.enabled", Value: true, Type: ConfigTypeBool, Description: "Lambda service", ReadOnly: false, EnvVar: "LAMBDA_ENABLED", Category: CategoryServices},
		"services.events.enabled":           {Key: "services.events.enabled", Value: true, Type: ConfigTypeBool, Description: "EventBridge service", ReadOnly: false, EnvVar: "EVENTS_ENABLED", Category: CategoryServices},
		"services.logs.enabled":             {Key: "services.logs.enabled", Value: true, Type: ConfigTypeBool, Description: "CloudWatch Logs service", ReadOnly: false, EnvVar: "LOGS_ENABLED", Category: CategoryServices},
		"services.kinesis.enabled":          {Key: "services.kinesis.enabled", Value: true, Type: ConfigTypeBool, Description: "Kinesis service", ReadOnly: false, EnvVar: "KINESIS_ENABLED", Category: CategoryServices},
		"services.stepfunctions.enabled":    {Key: "services.stepfunctions.enabled", Value: true, Type: ConfigTypeBool, Description: "Step Functions service", ReadOnly: false, EnvVar: "STEPFUNCTIONS_ENABLED", Category: CategoryServices},
		"services.apigateway.enabled":       {Key: "services.apigateway.enabled", Value: true, Type: ConfigTypeBool, Description: "API Gateway service", ReadOnly: false, EnvVar: "APIGATEWAY_ENABLED", Category: CategoryServices},
		"services.cognito.enabled":          {Key: "services.cognito.enabled", Value: true, Type: ConfigTypeBool, Description: "Cognito IDP service", ReadOnly: false, EnvVar: "COGNITO_ENABLED", Category: CategoryServices},
		"services.cognito_identity.enabled": {Key: "services.cognito_identity.enabled", Value: true, Type: ConfigTypeBool, Description: "Cognito Identity service", ReadOnly: false, EnvVar: "COGNITO_IDENTITY_ENABLED", Category: CategoryServices},
		"services.ssm.enabled":              {Key: "services.ssm.enabled", Value: true, Type: ConfigTypeBool, Description: "SSM service", ReadOnly: false, EnvVar: "SSM_ENABLED", Category: CategoryServices},
		"services.sesv2.enabled":            {Key: "services.sesv2.enabled", Value: true, Type: ConfigTypeBool, Description: "SESv2 service", ReadOnly: false, EnvVar: "SESV2_ENABLED", Category: CategoryServices},
		"services.scheduler.enabled":        {Key: "services.scheduler.enabled", Value: true, Type: ConfigTypeBool, Description: "Scheduler service", ReadOnly: false, EnvVar: "SCHEDULER_ENABLED", Category: CategoryServices},
		"services.cloudtrail.enabled":       {Key: "services.cloudtrail.enabled", Value: true, Type: ConfigTypeBool, Description: "CloudTrail service", ReadOnly: false, EnvVar: "CLOUDTRAIL_ENABLED", Category: CategoryServices},
		"services.acm.enabled":              {Key: "services.acm.enabled", Value: true, Type: ConfigTypeBool, Description: "ACM service", ReadOnly: false, EnvVar: "ACM_ENABLED", Category: CategoryServices},
		"services.cloudwatch.enabled":       {Key: "services.cloudwatch.enabled", Value: true, Type: ConfigTypeBool, Description: "CloudWatch service", ReadOnly: false, EnvVar: "CLOUDWATCH_ENABLED", Category: CategoryServices},
		"services.dynamodb.enabled":         {Key: "services.dynamodb.enabled", Value: true, Type: ConfigTypeBool, Description: "DynamoDB service", ReadOnly: false, EnvVar: "DYNAMODB_ENABLED", Category: CategoryServices},
		"services.kms.enabled":              {Key: "services.kms.enabled", Value: true, Type: ConfigTypeBool, Description: "KMS service", ReadOnly: false, EnvVar: "KMS_ENABLED", Category: CategoryServices},
		"services.s3.enabled":               {Key: "services.s3.enabled", Value: true, Type: ConfigTypeBool, Description: "S3 service", ReadOnly: false, EnvVar: "S3_ENABLED", Category: CategoryServices},
		"services.secretsmanager.enabled":   {Key: "services.secretsmanager.enabled", Value: true, Type: ConfigTypeBool, Description: "Secrets Manager service", ReadOnly: false, EnvVar: "SECRETSMANAGER_ENABLED", Category: CategoryServices},
		"services.sts.enabled":              {Key: "services.sts.enabled", Value: true, Type: ConfigTypeBool, Description: "STS service", ReadOnly: false, EnvVar: "STS_ENABLED", Category: CategoryServices},
		"services.iam.enabled":              {Key: "services.iam.enabled", Value: true, Type: ConfigTypeBool, Description: "IAM service", ReadOnly: false, EnvVar: "IAM_ENABLED", Category: CategoryServices},
		"services.timestream_write.enabled": {Key: "services.timestream_write.enabled", Value: true, Type: ConfigTypeBool, Description: "Timestream Write service", ReadOnly: false, EnvVar: "TIMESTREAM_WRITE_ENABLED", Category: CategoryServices},
		"services.timestream_query.enabled": {Key: "services.timestream_query.enabled", Value: true, Type: ConfigTypeBool, Description: "Timestream Query service", ReadOnly: false, EnvVar: "TIMESTREAM_QUERY_ENABLED", Category: CategoryServices},
		"services.athena.enabled":           {Key: "services.athena.enabled", Value: true, Type: ConfigTypeBool, Description: "Athena service", ReadOnly: false, EnvVar: "ATHENA_ENABLED", Category: CategoryServices},
		"services.appsync.enabled":          {Key: "services.appsync.enabled", Value: true, Type: ConfigTypeBool, Description: "AppSync service", ReadOnly: false, EnvVar: "APPSYNC_ENABLED", Category: CategoryServices},
		"services.neptune.enabled":          {Key: "services.neptune.enabled", Value: true, Type: ConfigTypeBool, Description: "Neptune service", ReadOnly: false, EnvVar: "NEPTUNE_ENABLED", Category: CategoryServices},
		"services.neptune_data.enabled":     {Key: "services.neptune_data.enabled", Value: true, Type: ConfigTypeBool, Description: "Neptune Data service", ReadOnly: false, EnvVar: "NEPTUNE_DATA_ENABLED", Category: CategoryServices},
		"services.neptune_graph.enabled":    {Key: "services.neptune_graph.enabled", Value: true, Type: ConfigTypeBool, Description: "Neptune Graph service", ReadOnly: false, EnvVar: "NEPTUNE_GRAPH_ENABLED", Category: CategoryServices},
		"services.cloudfront.enabled":       {Key: "services.cloudfront.enabled", Value: true, Type: ConfigTypeBool, Description: "CloudFront service", ReadOnly: false, EnvVar: "CLOUDFRONT_ENABLED", Category: CategoryServices},
		"services.wafv2.enabled":            {Key: "services.wafv2.enabled", Value: true, Type: ConfigTypeBool, Description: "WAFv2 service", ReadOnly: false, EnvVar: "WAFV2_ENABLED", Category: CategoryServices},
		"services.route53.enabled":          {Key: "services.route53.enabled", Value: true, Type: ConfigTypeBool, Description: "Route53 service", ReadOnly: false, EnvVar: "ROUTE53_ENABLED", Category: CategoryServices},
	}
}

func parseInt(s string, dest *int) (int, error) {
	if s == "" {
		return 0, ErrConfigInvalid
	}
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, ErrConfigInvalid
		}
		result = result*10 + int(c-'0')
	}
	*dest = result
	return result, nil
}

// GetSchema returns the configuration schema.
func GetSchema() []ConfigSchema {
	defaults := loadDefaults()
	schema := make([]ConfigSchema, 0, len(defaults))
	for _, entry := range defaults {
		schema = append(schema, ConfigSchema{
			Key:         entry.Key,
			Type:        entry.Type,
			Default:     entry.Value,
			Description: entry.Description,
			ReadOnly:    entry.ReadOnly,
			EnvVar:      entry.EnvVar,
			Category:    entry.Category,
		})
	}
	return schema
}

// GetKeysByCategory returns all configuration keys in a category.
func GetKeysByCategory(category ConfigCategory) []string {
	defaults := loadDefaults()
	keys := make([]string, 0)
	for key, entry := range defaults {
		if entry.Category == category {
			keys = append(keys, key)
		}
	}
	return keys
}
