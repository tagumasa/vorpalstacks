package authutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAWSCredential(t *testing.T) {
	t.Run("valid header", func(t *testing.T) {
		header := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;range;x-amz-date, Signature=..."
		cred := ParseAWSCredential(header)
		assert.NotNil(t, cred)
		assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", cred.AccessKeyID)
		assert.Equal(t, "20130524", cred.Date)
		assert.Equal(t, "us-east-1", cred.Region)
		assert.Equal(t, "s3", cred.SigningService)
		assert.Equal(t, "aws4_request", cred.RequestType)
	})

	t.Run("empty header", func(t *testing.T) {
		assert.Nil(t, ParseAWSCredential(""))
	})

	t.Run("no credential component", func(t *testing.T) {
		assert.Nil(t, ParseAWSCredential("Bearer some-token"))
	})

	t.Run("malformed credential", func(t *testing.T) {
		assert.Nil(t, ParseAWSCredential("Credential=incomplete"))
	})

	t.Run("credential not first part", func(t *testing.T) {
		header := "SignedHeaders=host, Credential=AKIA/20200101/eu-west-1/lambda/aws4_request"
		cred := ParseAWSCredential(header)
		assert.NotNil(t, cred)
		assert.Equal(t, "AKIA", cred.AccessKeyID)
		assert.Equal(t, "eu-west-1", cred.Region)
	})
}
