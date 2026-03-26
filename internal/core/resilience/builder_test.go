package resilience

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOperationResilienceBuilder_BuildWrapper(t *testing.T) {
	config := &ServiceResilienceConfig{
		ServiceName: "TestService",
		Operations: []*OperationConfig{
			{
				OperationName: "TestOperation",
			},
		},
	}

	builder := NewOperationResilienceBuilder(config)
	wrapper := builder.BuildWrapper("TestService", "TestOperation")

	assert.NotNil(t, wrapper)
	assert.NotNil(t, wrapper.circuitBreakerWrapper)
	assert.NotNil(t, wrapper.bulkheadWrapper)
	assert.NotNil(t, wrapper.adaptiveTimeoutWrapper)
	assert.NotNil(t, wrapper.retryPolicy)
}

func TestOperationResilienceBuilder_BuildWrapper_WithDefaults(t *testing.T) {
	config := &ServiceResilienceConfig{
		ServiceName: "TestService",
	}

	builder := NewOperationResilienceBuilder(config)
	wrapper := builder.BuildWrapper("TestService", "TestOperation")

	assert.NotNil(t, wrapper)
}

func TestOperationResilienceBuilder_BuildAllWrappers(t *testing.T) {
	config := &ServiceResilienceConfig{
		ServiceName: "TestService",
		Operations: []*OperationConfig{
			{OperationName: "Op1"},
			{OperationName: "Op2"},
		},
	}

	builder := NewOperationResilienceBuilder(config)
	wrappers := builder.BuildAllWrappers()

	assert.Len(t, wrappers, 2)
	assert.Contains(t, wrappers, "Op1")
	assert.Contains(t, wrappers, "Op2")
}

func TestOperationResilienceBuilder_WithCustomConfigs(t *testing.T) {
	customCB := &CircuitBreakerConfig{
		Name:         "CustomCB",
		MaxFailures:  10,
		ResetTimeout: 30 * time.Second,
	}

	config := &ServiceResilienceConfig{
		ServiceName: "TestService",
		Operations: []*OperationConfig{
			{
				OperationName:  "CustomOp",
				CircuitBreaker: customCB,
			},
		},
	}

	builder := NewOperationResilienceBuilder(config)
	wrapper := builder.BuildWrapper("TestService", "CustomOp")

	assert.NotNil(t, wrapper.circuitBreakerWrapper)
}

func TestResolveServiceName(t *testing.T) {
	config := &ServiceResilienceConfig{
		ServiceName: "TestService",
	}

	builder := NewOperationResilienceBuilder(config)

	result := builder.resolveServiceName(&OperationConfig{ServiceName: "OverrideService"})
	assert.Equal(t, "OverrideService", result)

	result = builder.resolveServiceName(nil)
	assert.Equal(t, "TestService", result)
}
