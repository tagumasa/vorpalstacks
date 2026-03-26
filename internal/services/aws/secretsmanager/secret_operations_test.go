//go:build ignore

package secretsmanager

import (
	"testing"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSecretValue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"SecretId": "my-secret",
			},
		}

		result, err := GetSecretValue(req)
		require.NoError(t, err)

		resp := result.(map[string]interface{})
		assert.Equal(t, "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret-abcdef", resp["ARN"])
		assert.Equal(t, "my-secret", resp["Name"])
		assert.NotEmpty(t, resp["VersionId"])
		assert.NotEmpty(t, resp["SecretString"])
		assert.NotEmpty(t, resp["SecretBinary"])
	})

	t.Run("Missing SecretId", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{},
		}

		result, err := GetSecretValue(req)
		assert.Nil(t, result)
		require.Error(t, err)
		assert.Equal(t, errors.ErrMissingParameter, err)
	})
}

func TestCreateSecret(t *testing.T) {
	t.Run("Success with SecretString", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"Name":         "new-secret",
				"SecretString": "my-secret-data",
			},
		}

		result, err := CreateSecret(req)
		require.NoError(t, err)

		resp := result.(map[string]interface{})
		assert.Equal(t, "arn:aws:secretsmanager:us-east-1:123456789012:secret:new-secret-abcdef", resp["ARN"])
		assert.Equal(t, "new-secret", resp["Name"])
		assert.Equal(t, "my-secret-data", resp["SecretString"])
	})

	t.Run("Missing Name", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"SecretString": "some-value",
			},
		}

		result, err := CreateSecret(req)
		assert.Nil(t, result)
		require.Error(t, err)
		assert.Equal(t, errors.ErrMissingParameter, err)
	})

	t.Run("Empty SecretString", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"Name":         "new-secret",
				"SecretString": "",
			},
		}

		result, err := CreateSecret(req)
		require.NoError(t, err)

		resp := result.(map[string]interface{})
		assert.Equal(t, "new-secret", resp["Name"])
		assert.Equal(t, "", resp["SecretString"])
	})
}

func TestListSecrets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{},
		}

		result, err := ListSecrets(req)
		require.NoError(t, err)

		resp := result.(map[string]interface{})
		secretList := resp["SecretList"].([]interface{})
		assert.Len(t, secretList, 1)

		secret := secretList[0].(map[string]interface{})
		assert.Equal(t, "arn:aws:secretsmanager:us-east-1:123456789012:secret:test-secret-abcdef", secret["ARN"])
		assert.Equal(t, "test-secret", secret["Name"])
		assert.Equal(t, "", resp["NextToken"])
	})

	t.Run("Returns default list", func(t *testing.T) {
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"MaxResults": 10,
			},
		}

		result, err := ListSecrets(req)
		require.NoError(t, err)

		resp := result.(map[string]interface{})
		assert.NotNil(t, resp["SecretList"])
		assert.NotNil(t, resp["NextToken"])
	})
}
