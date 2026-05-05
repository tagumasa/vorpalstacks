// Package config provides bootstrap and runtime configuration for vorpalstacks.
package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/serviceports"
)

// BootstrapConfig holds all configuration values read from environment variables
// at server startup, before the storage layer is available.
type BootstrapConfig struct {
	Port                  int
	GRPCWebPort           int
	GRPCWebBindAddr       string
	DataPath              string
	AccountID             string
	Region                string
	AccessKeyID           string
	SecretAccessKey       string
	SignatureVerification bool
	UseChainGateway       bool
	TLSEnabled            bool
	TLSPort               int
	TLSCertPath           string
	TLSKeyPath            string
	TLSHostname           string
	Route53DNSEnabled     bool
	DockerHost            string
	BindMode              string
	BindInterface         string
	bindAddr              string

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
	ACM             bool
	CloudWatch      bool
	DynamoDB        bool
	KMS             bool
	S3              bool
	SecretsManager  bool
	STS             bool
	IAM             bool
	TimestreamWrite bool
	TimestreamQuery bool
	Athena          bool
	AppSync         bool
	Neptune         bool
	NeptuneData     bool
	NeptuneGraph    bool
	CloudFront      bool
	WAFv2           bool
	Route53         bool
	EC2             bool
}

// ServerHost returns the server hostname suitable for self-referencing URLs.
func (c *BootstrapConfig) ServerHost() string {
	addr := c.bindAddr
	if addr == "" {
		addr = "127.0.0.1"
	}
	return addr + ":" + strconv.Itoa(c.Port)
}

// ResolvedBindAddr returns the bind address for console ports (HTTP + gRPC-Web).
// Must be called after LoadBootstrapConfig. Returns an error when bind_mode is
// "interface" but the specified IP does not exist on any host network interface.
func (c *BootstrapConfig) ResolvedBindAddr() (string, error) {
	if c.bindAddr != "" {
		return c.bindAddr, nil
	}
	switch c.BindMode {
	case "localhost", "":
		c.bindAddr = "127.0.0.1"
	case "all":
		c.bindAddr = "0.0.0.0"
	case "interface":
		if c.BindInterface == "" {
			return "", fmt.Errorf("bind_mode=interface requires BIND_INTERFACE to be set")
		}
		ip := net.ParseIP(c.BindInterface)
		if ip == nil {
			return "", fmt.Errorf("BIND_INTERFACE %q is not a valid IP address", c.BindInterface)
		}
		ifaces, err := net.Interfaces()
		if err != nil {
			return "", fmt.Errorf("failed to list network interfaces: %w", err)
		}
		found := false
		for _, iface := range ifaces {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, a := range addrs {
				var ifaceIP net.IP
				switch v := a.(type) {
				case *net.IPNet:
					ifaceIP = v.IP
				case *net.IPAddr:
					ifaceIP = v.IP
				}
				if ifaceIP != nil && ifaceIP.Equal(ip) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return "", fmt.Errorf("BIND_INTERFACE %q not found on any host network interface", c.BindInterface)
		}
		c.bindAddr = c.BindInterface
	default:
		return "", fmt.Errorf("invalid BIND_MODE %q (must be localhost, all, or interface)", c.BindMode)
	}
	return c.bindAddr, nil
}

// LoadBootstrapConfig reads all bootstrap configuration from environment variables
// and returns a populated BootstrapConfig. Defaults come from serviceports constants.
func LoadBootstrapConfig() *BootstrapConfig {
	accountId := os.Getenv("AWS_ACCOUNT_ID")
	if accountId == "" {
		accountId = "000000000000"
	}

	return &BootstrapConfig{
		Port:                  envOrInt("PORT", serviceports.HTTP),
		GRPCWebPort:           envOrInt("GRPC_WEB_PORT", serviceports.GRPCWeb),
		GRPCWebBindAddr:       envOr("GRPC_WEB_BIND_ADDR", "127.0.0.1"),
		DataPath:              envOr("DATA_PATH", "./data"),
		AccountID:             accountId,
		Region:                envOr("AWS_REGION", defaults.DefaultRegion),
		AccessKeyID:           os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SignatureVerification: envBool("SIGNATURE_VERIFICATION_ENABLED", true),
		UseChainGateway:       envBool("USE_CHAIN_GATEWAY", false),
		TLSEnabled:            envBool("TLS_ENABLED", false),
		TLSPort:               envOrInt("TLS_PORT", serviceports.TLS),
		TLSCertPath:           envOr("TLS_CERT_PATH", "auto"),
		TLSKeyPath:            envOr("TLS_KEY_PATH", "auto"),
		TLSHostname:           envOr("TLS_HOSTNAME", ""),
		Route53DNSEnabled:     envBool("ROUTE53_DNS_ENABLED", false),
		DockerHost:            envOr("DOCKER_HOST", "unix:///var/run/docker.sock"),
		BindMode:              envOr("BIND_MODE", "localhost"),
		BindInterface:         envOr("BIND_INTERFACE", ""),

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
		ACM:             envBool("ACM_ENABLED", true),
		CloudWatch:      envBool("CLOUDWATCH_ENABLED", true),
		DynamoDB:        envBool("DYNAMODB_ENABLED", true),
		KMS:             envBool("KMS_ENABLED", true),
		S3:              envBool("S3_ENABLED", true),
		SecretsManager:  envBool("SECRETSMANAGER_ENABLED", true),
		STS:             envBool("STS_ENABLED", true),
		IAM:             envBool("IAM_ENABLED", true),
		TimestreamWrite: envBool("TIMESTREAM_WRITE_ENABLED", true),
		TimestreamQuery: envBool("TIMESTREAM_QUERY_ENABLED", true),
		Athena:          envBool("ATHENA_ENABLED", true),
		AppSync:         envBool("APPSYNC_ENABLED", true),
		Neptune:         envBool("NEPTUNE_ENABLED", true),
		NeptuneData:     envBool("NEPTUNE_DATA_ENABLED", true),
		NeptuneGraph:    envBool("NEPTUNE_GRAPH_ENABLED", true),
		CloudFront:      envBool("CLOUDFRONT_ENABLED", true),
		WAFv2:           envBool("WAFV2_ENABLED", true),
		Route53:         envBool("ROUTE53_ENABLED", true),
		EC2:             envBool("EC2_ENABLED", true),
	}
}

func envOr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envOrInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
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
