package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	t.Run("StaticCredentialsProvider", func(t *testing.T) {
		provider := NewStaticCredentialsProvider(
			"access-key",
			"secret-key",
			"us-east-1",
			"s3",
		)

		creds, err := provider.GetCredentials()
		assert.NoError(t, err)
		assert.Equal(t, "access-key", creds.AccessKeyID)
		assert.Equal(t, "secret-key", creds.SecretAccessKey)
		assert.Equal(t, "us-east-1", creds.Region)
		assert.Equal(t, "s3", creds.Service)
	})

	t.Run("Credentials struct", func(t *testing.T) {
		creds := &Credentials{
			AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			SessionToken:    "token123",
			Region:          "us-west-2",
			Service:         "dynamodb",
		}

		assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", creds.AccessKeyID)
		assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", creds.SecretAccessKey)
		assert.Equal(t, "token123", creds.SessionToken)
		assert.Equal(t, "us-west-2", creds.Region)
		assert.Equal(t, "dynamodb", creds.Service)
	})
}
