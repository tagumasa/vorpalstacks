package config

import (
	"os"

	"vorpalstacks/internal/common/defaults"
)

func loadDefaults() map[string]ConfigEntry {
	return map[string]ConfigEntry{
		// Server Configuration
		"server.port": {
			Key:         "server.port",
			Value:       getEnvInt("PORT", 8080),
			Type:        ConfigTypePort,
			Description: "Main HTTP server port",
			ReadOnly:    false,
			EnvVar:      "PORT",
			Category:    CategoryServer,
		},
		"server.grpc_web_port": {
			Key:         "server.grpc_web_port",
			Value:       getEnvInt("GRPC_WEB_PORT", 9090),
			Type:        ConfigTypePort,
			Description: "gRPC-Web admin port",
			ReadOnly:    false,
			EnvVar:      "GRPC_WEB_PORT",
			Category:    CategoryServer,
		},
		"server.bind_addr": {
			Key:         "server.bind_addr",
			Value:       getEnvString("BIND_ADDR", "127.0.0.1"),
			Type:        ConfigTypeString,
			Description: "Server bind address",
			ReadOnly:    false,
			EnvVar:      "BIND_ADDR",
			Category:    CategoryServer,
		},
		"server.tls_enabled": {
			Key:         "server.tls_enabled",
			Value:       getEnvBool("TLS_ENABLED", false),
			Type:        ConfigTypeBool,
			Description: "Enable TLS for HTTPS server",
			ReadOnly:    false,
			EnvVar:      "TLS_ENABLED",
			Category:    CategoryServer,
		},
		"server.tls_port": {
			Key:         "server.tls_port",
			Value:       getEnvInt("TLS_PORT", 8443),
			Type:        ConfigTypePort,
			Description: "HTTPS server port",
			ReadOnly:    false,
			EnvVar:      "TLS_PORT",
			Category:    CategoryServer,
		},
		"server.tls_cert_path": {
			Key:         "server.tls_cert_path",
			Value:       getEnvString("TLS_CERT_PATH", "auto"),
			Type:        ConfigTypeString,
			Description: "TLS certificate path ('auto' generates and saves to pebble)",
			ReadOnly:    false,
			EnvVar:      "TLS_CERT_PATH",
			Category:    CategoryServer,
		},
		"server.tls_key_path": {
			Key:         "server.tls_key_path",
			Value:       getEnvString("TLS_KEY_PATH", "auto"),
			Type:        ConfigTypeString,
			Description: "TLS private key path ('auto' generates and saves to pebble)",
			ReadOnly:    false,
			EnvVar:      "TLS_KEY_PATH",
			Category:    CategoryServer,
		},

		// AWS Identity (Read-only at runtime)
		"aws.account_id": {
			Key:         "aws.account_id",
			Value:       getEnvString("AWS_ACCOUNT_ID", ""),
			Type:        ConfigTypeString,
			Description: "AWS account ID",
			ReadOnly:    true,
			EnvVar:      "AWS_ACCOUNT_ID",
			Category:    CategoryAWS,
		},
		"aws.region": {
			Key:         "aws.region",
			Value:       getEnvString("AWS_REGION", defaults.DefaultRegion),
			Type:        ConfigTypeString,
			Description: "Default AWS region",
			ReadOnly:    true,
			EnvVar:      "AWS_REGION",
			Category:    CategoryAWS,
		},

		// Storage Configuration
		"storage.data_path": {
			Key:         "storage.data_path",
			Value:       getEnvString("DATA_PATH", "./data"),
			Type:        ConfigTypeString,
			Description: "Data storage path",
			ReadOnly:    false,
			EnvVar:      "DATA_PATH",
			Category:    CategoryStorage,
		},
		"storage.metadata_path": {
			Key:         "storage.metadata_path",
			Value:       getEnvString("METADATA_PATH", ""),
			Type:        ConfigTypeString,
			Description: "Metadata path (defaults to data_path if empty)",
			ReadOnly:    false,
			EnvVar:      "METADATA_PATH",
			Category:    CategoryStorage,
		},

		// Feature Flags
		"features.audit_enabled": {
			Key:         "features.audit_enabled",
			Value:       getEnvBool("VS_AUDIT_ENABLED", false),
			Type:        ConfigTypeBool,
			Description: "CloudTrail audit logging",
			ReadOnly:    false,
			EnvVar:      "VS_AUDIT_ENABLED",
			Category:    CategoryFeatures,
		},
		"features.signature_verification": {
			Key:         "features.signature_verification",
			Value:       getEnvBool("SIGNATURE_VERIFICATION_ENABLED", true),
			Type:        ConfigTypeBool,
			Description: "AWS signature verification",
			ReadOnly:    false,
			EnvVar:      "SIGNATURE_VERIFICATION_ENABLED",
			Category:    CategoryFeatures,
		},
		"features.route53_dns": {
			Key:         "features.route53_dns",
			Value:       getEnvBool("ROUTE53_DNS_ENABLED", false),
			Type:        ConfigTypeBool,
			Description: "Route53 DNS server",
			ReadOnly:    false,
			EnvVar:      "ROUTE53_DNS_ENABLED",
			Category:    CategoryFeatures,
		},

		// External Service Endpoints
		"endpoints.base_url": {
			Key:         "endpoints.base_url",
			Value:       getEnvString("VS_BASE_URL", "http://localhost:8080"),
			Type:        ConfigTypeURL,
			Description: "Base URL for generated endpoints",
			ReadOnly:    false,
			EnvVar:      "VS_BASE_URL",
			Category:    CategoryEndpoints,
		},
		"endpoints.s3_website_suffix": {
			Key:         "endpoints.s3_website_suffix",
			Value:       getEnvString("VS_S3_WEBSITE_SUFFIX", "s3-website.{region}.amazonaws.com"),
			Type:        ConfigTypeString,
			Description: "S3 website domain suffix",
			ReadOnly:    false,
			EnvVar:      "VS_S3_WEBSITE_SUFFIX",
			Category:    CategoryEndpoints,
		},
		"endpoints.cognito_suffix": {
			Key:         "endpoints.cognito_suffix",
			Value:       getEnvString("VS_COGNITO_SUFFIX", "auth.{region}.amazoncognito.com"),
			Type:        ConfigTypeString,
			Description: "Cognito hosted UI suffix",
			ReadOnly:    false,
			EnvVar:      "VS_COGNITO_SUFFIX",
			Category:    CategoryEndpoints,
		},
		"endpoints.apigateway_suffix": {
			Key:         "endpoints.apigateway_suffix",
			Value:       getEnvString("VS_APIGW_SUFFIX", "execute-api.{region}.amazonaws.com"),
			Type:        ConfigTypeString,
			Description: "API Gateway suffix",
			ReadOnly:    false,
			EnvVar:      "VS_APIGW_SUFFIX",
			Category:    CategoryEndpoints,
		},

		// Default Public Endpoint Ports
		"ports.s3_website": {
			Key:         "ports.s3_website",
			Value:       getEnvInt("VS_PORT_S3_WEBSITE", 8081),
			Type:        ConfigTypePort,
			Description: "S3 website default port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_S3_WEBSITE",
			Category:    CategoryPorts,
		},
		"ports.apigateway": {
			Key:         "ports.apigateway",
			Value:       getEnvInt("VS_PORT_APIGW", 8082),
			Type:        ConfigTypePort,
			Description: "API Gateway invoke URL port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_APIGW",
			Category:    CategoryPorts,
		},
		"ports.cognito_hosted": {
			Key:         "ports.cognito_hosted",
			Value:       getEnvInt("VS_PORT_COGNITO_HOSTED", 8083),
			Type:        ConfigTypePort,
			Description: "Cognito Hosted UI port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_COGNITO_HOSTED",
			Category:    CategoryPorts,
		},
		"ports.cloudfront": {
			Key:         "ports.cloudfront",
			Value:       getEnvInt("VS_PORT_CLOUDFRONT", 8084),
			Type:        ConfigTypePort,
			Description: "CloudFront distribution port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_CLOUDFRONT",
			Category:    CategoryPorts,
		},
		"ports.lambda_url": {
			Key:         "ports.lambda_url",
			Value:       getEnvInt("VS_PORT_LAMBDA_URL", 8085),
			Type:        ConfigTypePort,
			Description: "Lambda Function URL port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_LAMBDA_URL",
			Category:    CategoryPorts,
		},
		"ports.appsync_events": {
			Key:         "ports.appsync_events",
			Value:       getEnvInt("VS_PORT_APPSYNC_EVENTS", 8086),
			Type:        ConfigTypePort,
			Description: "AppSync Events (WebSocket + HTTP publish) port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_APPSYNC_EVENTS",
			Category:    CategoryPorts,
		},
		"ports.neptune": {
			Key:         "ports.neptune",
			Value:       getEnvInt("VS_PORT_NEPTUNE", 8087),
			Type:        ConfigTypePort,
			Description: "Neptune DB cluster default port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_NEPTUNE",
			Category:    CategoryPorts,
		},
		"ports.route53_dns": {
			Key:         "ports.route53_dns",
			Value:       getEnvInt("VS_PORT_ROUTE53_DNS", 8088),
			Type:        ConfigTypePort,
			Description: "Route53 DNS server port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_ROUTE53_DNS",
			Category:    CategoryPorts,
		},
		"ports.route53_healthcheck": {
			Key:         "ports.route53_healthcheck",
			Value:       getEnvInt("VS_PORT_ROUTE53_HEALTHCHECK", 8089),
			Type:        ConfigTypePort,
			Description: "Route53 health check default endpoint port",
			ReadOnly:    false,
			EnvVar:      "VS_PORT_ROUTE53_HEALTHCHECK",
			Category:    CategoryPorts,
		},

		// HTTP Configuration (CORS)
		"http.cors_allowed_origins": {
			Key:         "http.cors_allowed_origins",
			Value:       getEnvString("VS_CORS_ALLOWED_ORIGINS", "*"),
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS origins ('*' allows all)",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_ORIGINS",
			Category:    CategoryHTTP,
		},
		"http.cors_allowed_methods": {
			Key:         "http.cors_allowed_methods",
			Value:       getEnvString("VS_CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS, HEAD"),
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS methods",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_METHODS",
			Category:    CategoryHTTP,
		},
		"http.cors_allowed_headers": {
			Key:         "http.cors_allowed_headers",
			Value:       getEnvString("VS_CORS_ALLOWED_HEADERS", "Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256"),
			Type:        ConfigTypeString,
			Description: "Comma-separated list of allowed CORS request headers",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_ALLOWED_HEADERS",
			Category:    CategoryHTTP,
		},
		"http.cors_expose_headers": {
			Key:         "http.cors_expose_headers",
			Value:       getEnvString("VS_CORS_EXPOSE_HEADERS", "x-amzn-RequestId, x-amzn-ErrorType, x-amzn-ErrorMessage"),
			Type:        ConfigTypeString,
			Description: "Comma-separated list of CORS headers exposed to the client",
			ReadOnly:    false,
			EnvVar:      "VS_CORS_EXPOSE_HEADERS",
			Category:    CategoryHTTP,
		},
	}
}

func getEnvString(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		if _, err := parseInt(val, &i); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1" || val == "yes"
	}
	return defaultVal
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
//
// Returns:
//   - []ConfigSchema: The configuration schema
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
//
// Parameters:
//   - category: The category to filter by
//
// Returns:
//   - []string: Configuration keys in the category
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
