package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateIAMUserID(t *testing.T) {
	id, err := GenerateIAMUserID()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "AID"))
	assert.Greater(t, len(id), 3)
}

func TestGenerateIAMAccessKeyID(t *testing.T) {
	id, err := GenerateIAMAccessKeyID()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "AKIA"))
	assert.Greater(t, len(id), 4)
}

func TestGenerateSQSMessageID(t *testing.T) {
	id, err := GenerateSQSMessageID()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "MSG"))
	assert.Greater(t, len(id), 3)
}

func TestGenerateS3UploadID(t *testing.T) {
	id, err := GenerateS3UploadID()
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Greater(t, len(id), 2)
}

func TestGenerateWAFID(t *testing.T) {
	id, err := GenerateWAFID()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "WAF"))
	assert.Greater(t, len(id), 3)
}

func TestGenerateACMCertificateID(t *testing.T) {
	id, err := GenerateACMCertificateID()
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Greater(t, len(id), 10)
}

func TestGenerateCloudFrontDistributionID(t *testing.T) {
	id, err := GenerateCloudFrontDistributionID()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "E"))
	assert.Greater(t, len(id), 1)
}

func TestGenerateSQSReceiptHandle(t *testing.T) {
	handle, err := GenerateSQSReceiptHandle()
	require.NoError(t, err)
	assert.NotEmpty(t, handle)
	assert.Greater(t, len(handle), 10)
}

func TestGenerateIAMGroupID(t *testing.T) {
	id, err := GenerateIAMGroupID()
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Greater(t, len(id), 4)
}

func TestGenerateIAMRoleID(t *testing.T) {
	id, err := GenerateIAMRoleID()
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Greater(t, len(id), 4)
}

func TestGenerateIAMPolicyID(t *testing.T) {
	id, err := GenerateIAMPolicyID()
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Greater(t, len(id), 4)
}

func TestGenerateACMCertificateSerial(t *testing.T) {
	serial, err := GenerateACMCertificateSerial()
	require.NoError(t, err)
	assert.NotEmpty(t, serial)
	assert.Greater(t, len(serial), 5)
}

func TestGenerateCloudFrontDomainName(t *testing.T) {
	t.Run("Generates valid domain", func(t *testing.T) {
		domain, err := GenerateCloudFrontDomainName("ABC123")
		require.NoError(t, err)
		assert.Equal(t, "ABC123.cloudfront.net", domain)
	})

	t.Run("Error on empty ID", func(t *testing.T) {
		domain, err := GenerateCloudFrontDomainName("")
		assert.Empty(t, domain)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestGenerateIAMSecretAccessKey(t *testing.T) {
	key, err := GenerateIAMSecretAccessKey()
	require.NoError(t, err)
	assert.Len(t, key, 40)
}
