// Package serviceconfig defines the service registry and types used by
// the admin console settings UI and the port allocation system.
// This package is a leaf — it imports only stdlib.
package serviceconfig

import (
	"sort"
)

// PortMode determines how a service exposes its resources.
type PortMode string

const (
	// PortModeFQDN routes all resources through a single shared listener using host-based routing.
	PortModeFQDN PortMode = "fqdn"
	// PortModeIndividual allocates a separate port for each resource.
	PortModeIndividual PortMode = "individual"
)

// ServiceDef describes a single AWS-compatible service managed by vorpalstacks.
type ServiceDef struct {
	Name         string   // unique key: "sns", "s3", "neptune_data"
	EnvVar       string   // e.g. "SNS_ENABLED"
	PortKey      string   // e.g. "ports.s3_website" (empty if no per-resource ports)
	Category     string   // display category
	HostSuffix   string   // FQDN suffix for FQDN mode (empty if not supported)
	SupportsFQDN bool     // true if this service can operate in FQDN mode
	Tags         []string // optional tags for grouping
}

// Services is the complete registry of all managed services (excludes EC2 —
// internal helper, not user-facing; excludes admin_auth/admin_config — gRPC-Web
// admin RPCs, always enabled).
var Services = []ServiceDef{
	{Name: "sns", EnvVar: "SNS_ENABLED", Category: "messaging"},
	{Name: "sqs", EnvVar: "SQS_ENABLED", Category: "messaging"},
	{Name: "lambda", EnvVar: "LAMBDA_ENABLED", Category: "compute"},
	{Name: "events", EnvVar: "EVENTS_ENABLED", Category: "integration"},
	{Name: "logs", EnvVar: "LOGS_ENABLED", Category: "management"},
	{Name: "kinesis", EnvVar: "KINESIS_ENABLED", Category: "analytics"},
	{Name: "stepfunctions", EnvVar: "STEPFUNCTIONS_ENABLED", Category: "integration"},
	{Name: "apigateway", EnvVar: "APIGATEWAY_ENABLED", PortKey: "ports.apigateway", Category: "networking", HostSuffix: "execute-api.{region}.amazonaws.com", SupportsFQDN: true},
	{Name: "cognito", EnvVar: "COGNITO_ENABLED", PortKey: "ports.cognito_hosted", Category: "security", HostSuffix: "auth.{region}.amazoncognito.com", SupportsFQDN: true},
	{Name: "cognito_identity", EnvVar: "COGNITO_IDENTITY_ENABLED", Category: "security"},
	{Name: "ssm", EnvVar: "SSM_ENABLED", Category: "management"},
	{Name: "sesv2", EnvVar: "SESV2_ENABLED", Category: "messaging"},
	{Name: "scheduler", EnvVar: "SCHEDULER_ENABLED", Category: "integration"},
	{Name: "cloudtrail", EnvVar: "CLOUDTRAIL_ENABLED", Category: "management"},
	{Name: "acm", EnvVar: "ACM_ENABLED", Category: "security"},
	{Name: "cloudwatch", EnvVar: "CLOUDWATCH_ENABLED", Category: "management"},
	{Name: "dynamodb", EnvVar: "DYNAMODB_ENABLED", Category: "database"},
	{Name: "kms", EnvVar: "KMS_ENABLED", Category: "security"},
	{Name: "s3", EnvVar: "S3_ENABLED", PortKey: "ports.s3_website", Category: "storage", HostSuffix: "s3-website.{region}.amazonaws.com", SupportsFQDN: true},
	{Name: "secretsmanager", EnvVar: "SECRETSMANAGER_ENABLED", Category: "security"},
	{Name: "sts", EnvVar: "STS_ENABLED", Category: "security"},
	{Name: "iam", EnvVar: "IAM_ENABLED", Category: "security"},
	{Name: "timestream_write", EnvVar: "TIMESTREAM_WRITE_ENABLED", Category: "analytics"},
	{Name: "timestream_query", EnvVar: "TIMESTREAM_QUERY_ENABLED", Category: "analytics"},
	{Name: "athena", EnvVar: "ATHENA_ENABLED", Category: "analytics"},
	{Name: "appsync", EnvVar: "APPSYNC_ENABLED", PortKey: "ports.appsync_events", Category: "integration", HostSuffix: "appsync-api.{region}.amazonaws.com", SupportsFQDN: true},
	{Name: "neptune", EnvVar: "NEPTUNE_ENABLED", Category: "database"},
	{Name: "neptune_data", EnvVar: "NEPTUNE_DATA_ENABLED", Category: "database"},
	{Name: "neptune_graph", EnvVar: "NEPTUNE_GRAPH_ENABLED", Category: "database"},
	{Name: "cloudfront", EnvVar: "CLOUDFRONT_ENABLED", PortKey: "ports.cloudfront", Category: "networking", HostSuffix: "cloudfront.net", SupportsFQDN: true},
	{Name: "wafv2", EnvVar: "WAFV2_ENABLED", Category: "security"},
	{Name: "route53", EnvVar: "ROUTE53_ENABLED", Category: "networking"},
	{Name: "lambda_url", EnvVar: "LAMBDA_ENABLED", PortKey: "ports.lambda_url", Category: "compute", HostSuffix: "lambda-url.{region}.on.aws", SupportsFQDN: true},
	{Name: "rdsdata", EnvVar: "RDS_DATA_ENABLED", Category: "database"},
	{Name: "iot", EnvVar: "IOT_ENABLED", Category: "iot"},
	{Name: "iotevents", EnvVar: "IOT_EVENTS_ENABLED", Category: "iot"},
}

// defaultPortMode returns the default port mode for this service.
func (s *ServiceDef) DefaultPortMode() string {
	if s.SupportsFQDN {
		return string(PortModeFQDN)
	}
	return string(PortModeIndividual)
}

// ByName returns the ServiceDef for the given name, or nil.
func ByName(name string) *ServiceDef {
	for i := range Services {
		if Services[i].Name == name {
			return &Services[i]
		}
	}
	return nil
}

// FQDNServices returns only services that support FQDN mode.
func FQDNServices() []ServiceDef {
	var result []ServiceDef
	for _, s := range Services {
		if s.SupportsFQDN {
			result = append(result, s)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// ByCategory groups services by their Category field.
func ByCategory() map[string][]ServiceDef {
	m := make(map[string][]ServiceDef)
	for _, s := range Services {
		m[s.Category] = append(m[s.Category], s)
	}
	return m
}
