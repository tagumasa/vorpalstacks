package lambda

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRuntime(t *testing.T) {
	t.Run("valid Node.js runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("nodejs24.x"))
		assert.True(t, ValidateRuntime("nodejs22.x"))
	})

	t.Run("valid Python runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("python3.14"))
		assert.True(t, ValidateRuntime("python3.13"))
		assert.True(t, ValidateRuntime("python3.12"))
		assert.True(t, ValidateRuntime("python3.11"))
		assert.True(t, ValidateRuntime("python3.10"))
	})

	t.Run("valid Java runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("java25"))
		assert.True(t, ValidateRuntime("java21"))
		assert.True(t, ValidateRuntime("java17"))
		assert.True(t, ValidateRuntime("java11"))
		assert.True(t, ValidateRuntime("java8.al2"))
	})

	t.Run("valid .NET runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("dotnet10"))
		assert.True(t, ValidateRuntime("dotnet8"))
	})

	t.Run("valid Ruby runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("ruby4.0"))
		assert.True(t, ValidateRuntime("ruby3.4"))
		assert.True(t, ValidateRuntime("ruby3.3"))
	})

	t.Run("valid custom runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("provided.al2023"))
		assert.True(t, ValidateRuntime("provided.al2"))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, ValidateRuntime("PYTHON3.12"))
		assert.True(t, ValidateRuntime("NodeJS22.x"))
	})

	t.Run("invalid runtime", func(t *testing.T) {
		assert.False(t, ValidateRuntime("python3.7"))
		assert.False(t, ValidateRuntime("nodejs14.x"))
		assert.False(t, ValidateRuntime("invalid"))
		assert.False(t, ValidateRuntime(""))
	})
}

func TestValidateHandler(t *testing.T) {
	t.Run("valid Python handler", func(t *testing.T) {
		err := ValidateHandler("python3.12", "myhandler.handle")
		assert.NoError(t, err)
	})

	t.Run("valid Node.js handler", func(t *testing.T) {
		err := ValidateHandler("nodejs20.x", "index.handler")
		assert.NoError(t, err)
	})

	t.Run("valid Java handler", func(t *testing.T) {
		err := ValidateHandler("java17", "com.example.MyHandler::handleRequest")
		assert.NoError(t, err)
	})

	t.Run("valid Java handler with package only", func(t *testing.T) {
		err := ValidateHandler("java17", "com.example.MyHandler.handleRequest")
		assert.NoError(t, err)
	})

	t.Run("empty handler returns error", func(t *testing.T) {
		err := ValidateHandler("python3.12", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Handler cannot be empty")
	})

	t.Run("Python handler without dot returns error", func(t *testing.T) {
		err := ValidateHandler("python3.12", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Python handler must be in the format module.function")
	})

	t.Run("Node.js handler without dot returns error", func(t *testing.T) {
		err := ValidateHandler("nodejs20.x", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Node.js handler must be in the format file.function")
	})

	t.Run("Java handler without proper format returns error", func(t *testing.T) {
		err := ValidateHandler("java17", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Java handler must be in the format package.Class::method")
	})
}

func TestValidateFunctionName(t *testing.T) {
	t.Run("valid names", func(t *testing.T) {
		assert.NoError(t, validateFunctionName("my-function"))
		assert.NoError(t, validateFunctionName("my_function"))
		assert.NoError(t, validateFunctionName("MyFunction123"))
	})

	t.Run("too long", func(t *testing.T) {
		assert.Error(t, validateFunctionName("a"+strings.Repeat("b", 64)))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Error(t, validateFunctionName(""))
	})

	t.Run("invalid characters", func(t *testing.T) {
		assert.Error(t, validateFunctionName("my.function"))
		assert.Error(t, validateFunctionName("my function"))
	})
}

func TestValidateAuthType(t *testing.T) {
	t.Run("NONE is valid", func(t *testing.T) {
		assert.NoError(t, validateAuthType("NONE"))
	})

	t.Run("AWS_IAM is valid", func(t *testing.T) {
		assert.NoError(t, validateAuthType("AWS_IAM"))
	})

	t.Run("empty is rejected (H3 fix)", func(t *testing.T) {
		err := validateAuthType("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AuthType is required")
	})

	t.Run("invalid value rejected", func(t *testing.T) {
		err := validateAuthType("BASIC")
		assert.Error(t, err)
	})
}

func TestValidateInvokeMode(t *testing.T) {
	t.Run("BUFFERED is valid", func(t *testing.T) {
		assert.NoError(t, validateInvokeMode("BUFFERED"))
	})

	t.Run("RESPONSE_STREAM is valid", func(t *testing.T) {
		assert.NoError(t, validateInvokeMode("RESPONSE_STREAM"))
	})

	t.Run("empty defaults to BUFFERED (no error)", func(t *testing.T) {
		assert.NoError(t, validateInvokeMode(""))
	})

	t.Run("invalid rejected", func(t *testing.T) {
		assert.Error(t, validateInvokeMode("SYNC"))
	})
}

func TestValidateMaximumEventAgeInSeconds(t *testing.T) {
	assert.NoError(t, validateMaximumEventAgeInSeconds(60))
	assert.NoError(t, validateMaximumEventAgeInSeconds(21600))
	assert.Error(t, validateMaximumEventAgeInSeconds(59))
	assert.Error(t, validateMaximumEventAgeInSeconds(21601))
}

func TestValidateMaximumRetryAttempts(t *testing.T) {
	assert.NoError(t, validateMaximumRetryAttempts(0))
	assert.NoError(t, validateMaximumRetryAttempts(2))
	assert.Error(t, validateMaximumRetryAttempts(-1))
	assert.Error(t, validateMaximumRetryAttempts(3))
}

func TestValidateCodeSigningConfigArn(t *testing.T) {
	assert.NoError(t, validateCodeSigningConfigArn(""))
	assert.Error(t, validateCodeSigningConfigArn("arn:aws:lambda:us-east-1:123:code-signing-config:csc-abc"))
}

func TestIsValidPrincipal(t *testing.T) {
	t.Run("wildcard", func(t *testing.T) {
		assert.True(t, isValidPrincipal("*"))
	})

	t.Run("IAM ARN", func(t *testing.T) {
		assert.True(t, isValidPrincipal("arn:aws:iam::123:root"))
	})

	t.Run("known service principal", func(t *testing.T) {
		assert.True(t, isValidPrincipal("lambda.amazonaws.com"))
		assert.True(t, isValidPrincipal("s3.amazonaws.com"))
		assert.True(t, isValidPrincipal("events.amazonaws.com"))
	})

	t.Run("typo rejected (M5 fix)", func(t *testing.T) {
		assert.False(t, isValidPrincipal("lamda.amazonaws.com"))
	})

	t.Run("spoof rejected (M5 fix)", func(t *testing.T) {
		assert.False(t, isValidPrincipal("evil.amazonaws.com"))
	})

	t.Run("unknown suffix rejected", func(t *testing.T) {
		assert.False(t, isValidPrincipal("fake.amazonaws.com"))
	})
}

func TestPrincipalType(t *testing.T) {
	assert.Equal(t, "", principalType("*"))
	assert.Equal(t, "AWS", principalType("arn:aws:iam::123:root"))
	assert.Equal(t, "Service", principalType("lambda.amazonaws.com"))
}
