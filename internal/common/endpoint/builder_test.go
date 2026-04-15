package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilderWithConfig(t *testing.T) {
	b := NewBuilderWithConfig("http://localhost:9000", "eu-west-1")
	assert.Equal(t, "http://localhost:9000", b.BaseURL())
}

func TestBuilder_SQSQueueURL(t *testing.T) {
	b := NewBuilderWithConfig("http://localhost:8080", "us-east-1")
	url := b.SQSQueueURL("123456789012", "my-queue")
	assert.Equal(t, "http://localhost:8080/123456789012/my-queue", url)
}

func TestBuilder_CloudFrontURL(t *testing.T) {
	b := NewBuilderWithConfig("http://localhost:8080", "us-east-1")
	url := b.CloudFrontURL("E1234567890")
	assert.Equal(t, "https://E1234567890.cloudfront.net", url)
}

func TestBuilder_LambdaFunctionURL(t *testing.T) {
	b := NewBuilderWithConfig("http://localhost:8080", "us-east-1")
	url := b.LambdaFunctionURL("my-func", "eu-west-1")
	assert.Equal(t, "https://my-func.lambda-url.eu-west-1.on.aws", url)
}

func TestGetBuilder_Initialises(t *testing.T) {
	b := GetBuilder()
	assert.NotNil(t, b)
	same := GetBuilder()
	assert.Equal(t, b, same)
}
