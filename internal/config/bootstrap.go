// Package config provides bootstrap and runtime configuration for vorpalstacks.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// BootstrapConfig holds all configuration values read from environment variables
// at server startup, before the storage layer is available.
type BootstrapConfig struct {
	Port                  string
	GRPCWebPort           string
	GRPCWebBindAddr       string
	DataPath              string
	MetadataPath          string
	AccountID             string
	Region                string
	AccessKeyID           string
	SecretAccessKey       string
	SignatureVerification bool
	UseChainGateway       bool
	TLSEnabled            bool
	TLSPort               string
	TLSCertPath           string
	TLSKeyPath            string
	TLSHostname           string
	Route53DNSEnabled     bool
	Route53DNSBindAddr    string
	DockerHost            string

	SNS             bool
	SQS             bool
	Lambda          bool
	Events          bool
	Logs            bool
	Kinesis         bool
	StepFunctions   bool
	APIGateway      bool
	Cognito         bool
	CognitoIdentity bool
	SSM             bool
	SESv2           bool
	Scheduler       bool
	CloudTrail      bool
	TimestreamWrite bool
	TimestreamQuery bool
	Athena          bool
}

// ServerHost returns the server hostname suitable for self-referencing URLs.
func (c *BootstrapConfig) ServerHost() string {
	return "localhost:" + c.Port
}

// LoadBootstrapConfig reads all bootstrap configuration from environment variables
// and returns a populated BootstrapConfig. Defaults are applied when variables
// are unset.
func LoadBootstrapConfig() *BootstrapConfig {
	accountId := envOr("AWS_ACCOUNT_ID", "")
	if accountId == "" {
		accountId = "000000000000"
	}

	return &BootstrapConfig{
		Port:                  envOr("PORT", "8080"),
		GRPCWebPort:           envOr("GRPC_WEB_PORT", "9090"),
		GRPCWebBindAddr:       envOr("GRPC_WEB_BIND_ADDR", "127.0.0.1"),
		DataPath:              envOr("DATA_PATH", "./data"),
		MetadataPath:          envOr("METADATA_PATH", ""),
		AccountID:             accountId,
		Region:                envOr("AWS_REGION", "us-east-1"),
		AccessKeyID:           os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SignatureVerification: envBool("SIGNATURE_VERIFICATION_ENABLED", true),
		UseChainGateway:       envBool("USE_CHAIN_GATEWAY", false),
		TLSEnabled:            envBool("TLS_ENABLED", false),
		TLSPort:               envOr("TLS_PORT", "8443"),
		TLSCertPath:           envOr("TLS_CERT_PATH", "auto"),
		TLSKeyPath:            envOr("TLS_KEY_PATH", "auto"),
		TLSHostname:           envOr("TLS_HOSTNAME", ""),
		Route53DNSEnabled:     envBool("ROUTE53_DNS_ENABLED", false),
		Route53DNSBindAddr:    envOr("ROUTE53_DNS_BIND_ADDR", "127.0.0.1"),
		DockerHost:            envOr("DOCKER_HOST", "unix:///var/run/docker.sock"),

		SNS:             envBool("SNS_ENABLED", true),
		SQS:             envBool("SQS_ENABLED", true),
		Lambda:          envBool("LAMBDA_ENABLED", true),
		Events:          envBool("EVENTS_ENABLED", true),
		Logs:            envBool("LOGS_ENABLED", true),
		Kinesis:         envBool("KINESIS_ENABLED", true),
		StepFunctions:   envBool("STEPFUNCTIONS_ENABLED", true),
		APIGateway:      envBool("APIGATEWAY_ENABLED", true),
		Cognito:         envBool("COGNITO_ENABLED", true),
		CognitoIdentity: envBool("COGNITO_IDENTITY_ENABLED", true),
		SSM:             envBool("SSM_ENABLED", true),
		SESv2:           envBool("SESV2_ENABLED", true),
		Scheduler:       envBool("SCHEDULER_ENABLED", true),
		CloudTrail:      envBool("CLOUDTRAIL_ENABLED", true),
		TimestreamWrite: envBool("TIMESTREAM_WRITE_ENABLED", true),
		TimestreamQuery: envBool("TIMESTREAM_QUERY_ENABLED", true),
		Athena:          envBool("ATHENA_ENABLED", true),
	}
}

func envOr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, _ := strconv.ParseBool(v)
	return b
}

// PrintStartupBanner writes the server startup information to stdout.
func (c *BootstrapConfig) PrintStartupBanner() {
	fmt.Printf("Starting AWS-compatible server on :%s\n", c.Port)
	fmt.Printf("Data path: %s\n", c.DataPath)
	fmt.Printf("Account ID: %s\n", c.AccountID)
	if c.UseChainGateway {
		fmt.Printf("Gateway mode: Chain (metadata-based routing)\n")
	} else {
		fmt.Printf("Gateway mode: Standard (hardcoded routing)\n")
	}
	fmt.Printf("Services: IAM, STS, KMS, S3, Route53, CloudWatch, DynamoDB, ACM, CloudFront, WAF, SecretsManager, SNS, SQS, Events, Lambda, APIGateway, StepFunctions, Cognito, SSM, Logs, SESv2, Scheduler, Kinesis, CloudTrail, TimestreamWrite, TimestreamQuery, Athena\n")
	if c.Route53DNSEnabled {
		fmt.Printf("Route53 DNS server: enabled on %s:53\n", c.Route53DNSBindAddr)
	}
	if c.SignatureVerification {
		fmt.Printf("Signature verification: enabled\n")
	} else {
		fmt.Printf("Signature verification: disabled\n")
	}
	if c.AccessKeyID == "" {
		fmt.Println("Warning: AWS_ACCESS_KEY_ID not set")
	}
	if c.SecretAccessKey == "" {
		fmt.Println("Warning: AWS_SECRET_ACCESS_KEY not set")
	}
}
