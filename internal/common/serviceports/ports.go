package serviceports

const (
	// HTTP is the default HTTP listener port.
	HTTP = 50080
	// GRPCWeb is the gRPC-Web admin console port.
	GRPCWeb = 50090
	// TLS is the HTTPS listener port.
	TLS = 50443
	// Route53DNS is the Route53 DNS service port.
	Route53DNS = 50088
	// Route53HC is the Route53 health check port.
	Route53HC = 50089
)

const (
	// S3Website is the S3 static website hosting port.
	S3Website = 50101
	// APIGateway is the API Gateway REST API port.
	APIGateway = 50102
	// Cognito is the Cognito user pool port.
	Cognito = 50103
	// CloudFront is the CloudFront distribution port.
	CloudFront = 50104
	// LambdaURL is the Lambda function URL port.
	LambdaURL = 50105
	// AppSync is the AppSync GraphQL API port.
	AppSync = 50106
	// IotMQTT is the IoT Core MQTT broker port.
	IotMQTT = 50107
)

const (
	// DynamicRangeStart is the start of the dynamic port allocation range.
	DynamicRangeStart = 50200
	// DynamicRangeEnd is the end of the dynamic port allocation range.
	DynamicRangeEnd = 50400
)
