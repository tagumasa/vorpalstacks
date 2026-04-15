package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchema(t *testing.T) {
	schema := GetSchema()
	assert.NotEmpty(t, schema)

	found := false
	for _, s := range schema {
		if s.Key == "server.port" {
			found = true
			assert.Equal(t, ConfigTypePort, s.Type)
			assert.Equal(t, CategoryServer, s.Category)
			assert.False(t, s.ReadOnly)
			assert.Equal(t, "PORT", s.EnvVar)
		}
	}
	assert.True(t, found, "server.port should be in schema")
}

func TestGetKeysByCategory(t *testing.T) {
	t.Run("server category", func(t *testing.T) {
		keys := GetKeysByCategory(CategoryServer)
		assert.NotEmpty(t, keys)
		assert.Contains(t, keys, "server.port")
	})

	t.Run("ports category", func(t *testing.T) {
		keys := GetKeysByCategory(CategoryPorts)
		assert.NotEmpty(t, keys)
		assert.Contains(t, keys, "ports.s3_website")
	})

	t.Run("empty category", func(t *testing.T) {
		keys := GetKeysByCategory(ConfigCategory("nonexistent"))
		assert.Empty(t, keys)
	})
}

func TestParseInt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var i int
		n, err := parseInt("42", &i)
		assert.NoError(t, err)
		assert.Equal(t, 42, n)
		assert.Equal(t, 42, i)
	})

	t.Run("empty", func(t *testing.T) {
		var i int
		_, err := parseInt("", &i)
		assert.Equal(t, ErrConfigInvalid, err)
	})

	t.Run("invalid", func(t *testing.T) {
		var i int
		_, err := parseInt("abc", &i)
		assert.Equal(t, ErrConfigInvalid, err)
	})

	t.Run("zero", func(t *testing.T) {
		var i int
		n, err := parseInt("0", &i)
		assert.NoError(t, err)
		assert.Equal(t, 0, n)
	})
}

func TestConfigTypes(t *testing.T) {
	assert.Equal(t, ConfigType("STRING"), ConfigTypeString)
	assert.Equal(t, ConfigType("INT"), ConfigTypeInt)
	assert.Equal(t, ConfigType("BOOL"), ConfigTypeBool)
	assert.Equal(t, ConfigType("PORT"), ConfigTypePort)
	assert.Equal(t, ConfigType("URL"), ConfigTypeURL)
}

func TestConfigSources(t *testing.T) {
	assert.Equal(t, ConfigSource("ENV"), ConfigSourceEnv)
	assert.Equal(t, ConfigSource("STORE"), ConfigSourceStore)
	assert.Equal(t, ConfigSource("DEFAULT"), ConfigSourceDefault)
}

func TestConfigCategories(t *testing.T) {
	assert.Equal(t, ConfigCategory("server"), CategoryServer)
	assert.Equal(t, ConfigCategory("aws"), CategoryAWS)
	assert.Equal(t, ConfigCategory("storage"), CategoryStorage)
	assert.Equal(t, ConfigCategory("features"), CategoryFeatures)
	assert.Equal(t, ConfigCategory("endpoints"), CategoryEndpoints)
	assert.Equal(t, ConfigCategory("ports"), CategoryPorts)
}

func TestConfigErrors(t *testing.T) {
	assert.Equal(t, "configuration not found", ErrConfigNotFound.Error())
	assert.Equal(t, "configuration is read-only", ErrConfigReadOnly.Error())
	assert.Equal(t, "invalid configuration value", ErrConfigInvalid.Error())
	assert.Equal(t, "config store not initialised", ErrConfigNotInitialised.Error())
}

func TestConfigEntry(t *testing.T) {
	e := ConfigEntry{
		Key:      "test.key",
		Value:    "test-value",
		Type:     ConfigTypeString,
		Source:   ConfigSourceDefault,
		Category: CategoryServer,
	}
	assert.Equal(t, "test.key", e.Key)
	assert.Equal(t, "test-value", e.Value)
}

func TestConfigSchema(t *testing.T) {
	s := ConfigSchema{
		Key:         "test",
		Type:        ConfigTypeBool,
		Default:     true,
		Description: "desc",
		ReadOnly:    true,
	}
	assert.Equal(t, "test", s.Key)
	assert.Equal(t, ConfigTypeBool, s.Type)
	assert.True(t, s.Default.(bool))
	assert.True(t, s.ReadOnly)
}
