package lambda

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

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

func TestValidateNamespacedFunctionName(t *testing.T) {
	t.Run("pattern-valid references", func(t *testing.T) {
		valid := []string{
			"my-function",
			"my.function",
			"arn:aws:lambda:us-west-2:123456789012:function:my-function",
			"arn:aws:lambda:us-west-2:123456789012:function:my-function:prod",
			"arn:aws:lambda:us-west-2:123456789012:function:fn:$LATEST.PUBLISHED",
			"123456789012:function:my-function",
			"arn:aws-us-gov:lambda:us-gov-west-1:123456789012:function:fn",
			strings.Repeat("a", 256),
		}
		for _, ref := range valid {
			assert.NoError(t, validateNamespacedFunctionName(ref), "ref %q", ref)
		}
	})

	t.Run("pattern-invalid references", func(t *testing.T) {
		invalid := []string{
			"",
			strings.Repeat("a", 257),
			"arn:aws:lambda:NOT_A_REGION:123456789012:function:fn",
			"my function",
			"arn:aws:sqs:us-east-1:123456789012:queue",
		}
		for _, ref := range invalid {
			assert.Error(t, validateNamespacedFunctionName(ref), "ref %q", ref)
		}
	})
}

func TestValidateFileSystemConfigs(t *testing.T) {
	const s3FilesAP = "arn:aws:s3files:us-east-1:123456789012:file-system/fs-0000000000000000000000000000000/access-point/fsap-0123456789abcdef0"
	const efsAP = "arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-0123456789abcdef0"

	t.Run("accepts EFS entries and absent member", func(t *testing.T) {
		assert.NoError(t, validateFileSystemConfigs(nil))
		assert.NoError(t, validateFileSystemConfigs([]lambdastore.FileSystemConfig{
			{Arn: efsAP, LocalMountPath: "/mnt/efs"},
		}))
	})

	t.Run("accepts S3FilesConfig on an S3 Files access point", func(t *testing.T) {
		assert.NoError(t, validateFileSystemConfigs([]lambdastore.FileSystemConfig{
			{Arn: s3FilesAP, LocalMountPath: "/mnt/s3", S3FilesConfig: &lambdastore.S3FilesConfig{
				DirectS3Read: lambdastore.DirectS3ReadEnabled,
			}},
		}))
	})

	t.Run("rejects unknown DirectS3Read values", func(t *testing.T) {
		assert.Error(t, validateFileSystemConfigs([]lambdastore.FileSystemConfig{
			{Arn: s3FilesAP, LocalMountPath: "/mnt/s3", S3FilesConfig: &lambdastore.S3FilesConfig{
				DirectS3Read: "FAST",
			}},
		}))
	})

	t.Run("rejects S3FilesConfig on a non-S3-Files access point", func(t *testing.T) {
		assert.Error(t, validateFileSystemConfigs([]lambdastore.FileSystemConfig{
			{Arn: efsAP, LocalMountPath: "/mnt/efs", S3FilesConfig: &lambdastore.S3FilesConfig{
				DirectS3Read: lambdastore.DirectS3ReadAuto,
			}},
		}))
	})
}

func TestValidateTimeoutModelRange(t *testing.T) {
	// CreateFunction model @range on Timeout: min 1, max 5400 seconds.
	assert.NoError(t, validateTimeout(1))
	assert.NoError(t, validateTimeout(5400))
	assert.Error(t, validateTimeout(0), "timeout below the modelled minimum must be rejected")
	assert.Error(t, validateTimeout(5401), "timeout above the modelled maximum must be rejected")
}

func TestValidateMemorySizeModelRange(t *testing.T) {
	// CreateFunction model @range on MemorySize: min 128, max 32768 MB.
	assert.NoError(t, validateMemorySize(128))
	assert.NoError(t, validateMemorySize(32768))
	assert.Error(t, validateMemorySize(127), "memory size below the modelled minimum must be rejected")
	assert.Error(t, validateMemorySize(32769), "memory size above the modelled maximum must be rejected")
}

func TestValidateEphemeralStorageSizeModelRange(t *testing.T) {
	// CreateFunction model @range on EphemeralStorageSize: min 512, max 32768 MB.
	assert.NoError(t, validateEphemeralStorageSize(512))
	assert.NoError(t, validateEphemeralStorageSize(32768))
	assert.Error(t, validateEphemeralStorageSize(511), "ephemeral storage below the modelled minimum must be rejected")
	assert.Error(t, validateEphemeralStorageSize(32769), "ephemeral storage above the modelled maximum must be rejected")
}

func TestValidateAuthType(t *testing.T) {
	t.Run("NONE is valid", func(t *testing.T) {
		assert.NoError(t, validateAuthType("NONE"))
	})

	t.Run("AWS_IAM is valid", func(t *testing.T) {
		assert.NoError(t, validateAuthType("AWS_IAM"))
	})

	t.Run("empty is rejected", func(t *testing.T) {
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

	t.Run("typo rejected", func(t *testing.T) {
		assert.False(t, isValidPrincipal("lamda.amazonaws.com"))
	})

	t.Run("spoof rejected", func(t *testing.T) {
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

func TestValidateEnvironmentVariableKeys(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"OK", true},
		{"Ab", true},
		{"a9_b", true},
		{"Z_9", true},
		{"a", false},         // single character
		{"1BAD", false},      // starts with a digit
		{"_x", false},        // starts with an underscore
		{"a-b", false},       // hyphen not allowed
		{"a.b", false},       // dot not allowed
		{"has space", false}, // space not allowed
		{"", false},          // empty key
	}
	for _, tc := range cases {
		env := &lambdastore.Environment{Variables: map[string]string{tc.key: "v"}}
		err := validateEnvironmentVariables(env)
		if tc.want && err != nil {
			t.Fatalf("key %q should be valid, got %v", tc.key, err)
		}
		if !tc.want && err == nil {
			t.Fatalf("key %q should be rejected", tc.key)
		}
	}

	reserved := &lambdastore.Environment{Variables: map[string]string{"AWS_LAMBDA_x": "v"}}
	if err := validateEnvironmentVariables(reserved); err == nil {
		t.Fatal("reserved AWS_LAMBDA_ prefix must be rejected")
	}
}

// TestValidateLoggingConfigModelEnums pins the modelled LoggingConfig member
// constraints: the LogFormat enum (JSON, Text), the ApplicationLogLevel
// enum (TRACE..FATAL), the SystemLogLevel enum (DEBUG..WARN), and the
// LogGroup length and pattern traits.
func TestValidateLoggingConfigModelEnums(t *testing.T) {
	valid := []*lambdastore.LoggingConfig{
		nil,
		{},
		{LogFormat: "JSON"},
		{LogFormat: "Text", ApplicationLogLevel: "TRACE", SystemLogLevel: "WARN"},
		{ApplicationLogLevel: "FATAL"},
		{SystemLogLevel: "DEBUG"},
		{LogGroup: "/aws/lambda/my-function"},
	}
	for _, lc := range valid {
		if err := validateLoggingConfig(lc); err != nil {
			t.Fatalf("logging config %+v should be valid, got %v", lc, err)
		}
	}

	invalid := []*lambdastore.LoggingConfig{
		{LogFormat: "Banana"},
		{ApplicationLogLevel: "NOTICE"},
		{SystemLogLevel: "TRACE"},
		{SystemLogLevel: "ERROR"},
		{LogGroup: "has space"},
		{LogGroup: "a:b"},
		{LogGroup: strings.Repeat("a", 513)},
	}
	for _, lc := range invalid {
		if err := validateLoggingConfig(lc); err == nil {
			t.Fatalf("logging config %+v should be rejected", lc)
		}
	}
}

// TestValidateImageConfigModelLengths pins the modelled ImageConfig member
// constraints: EntryPoint and Command carry at most 1500 entries each and
// WorkingDirectory is at most 1000 characters.
func TestValidateImageConfigModelLengths(t *testing.T) {
	entries := func(n int) []string {
		list := make([]string, n)
		for i := range list {
			list[i] = "entry"
		}
		return list
	}

	valid := []*lambdastore.ImageConfig{
		nil,
		{},
		{EntryPoint: []string{"/bin/app"}},
		{Command: entries(1500)},
		{WorkingDirectory: "/var/task"},
		{WorkingDirectory: strings.Repeat("d", 1000)},
	}
	for _, ic := range valid {
		if err := validateImageConfig(ic); err != nil {
			t.Fatalf("image config %+v should be valid, got %v", ic, err)
		}
	}

	invalid := []*lambdastore.ImageConfig{
		{EntryPoint: entries(1501)},
		{Command: entries(1501)},
		{WorkingDirectory: strings.Repeat("d", 1001)},
	}
	for _, ic := range invalid {
		if err := validateImageConfig(ic); err == nil {
			t.Fatalf("image config should be rejected: %+v", ic)
		}
	}
}
