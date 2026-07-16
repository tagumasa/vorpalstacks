package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	DefaultAccessKeyID     = "test"
	DefaultSecretAccessKey = "test"
	DefaultSessionToken    = ""
)

type AWSConfig struct {
	Endpoint   string
	Region     string
	HTTPClient *http.Client
}

func LoadDefaultAWSConfig(cfg AWSConfig) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("endpoint not configured")
	})

	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     DefaultAccessKeyID,
				SecretAccessKey: DefaultSecretAccessKey,
				SessionToken:    DefaultSessionToken,
			}, nil
		})),
		config.WithHTTPClient(client),
	)
}
