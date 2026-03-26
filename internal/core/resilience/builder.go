// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"fmt"
	"sync"
)

// OperationConfig contains configuration for a single operation.
type OperationConfig struct {
	ServiceName   string
	OperationName string

	CircuitBreaker  *CircuitBreakerConfig
	Bulkhead        *BulkheadConfig
	AdaptiveTimeout *AdaptiveTimeoutConfig
	RetryPolicy     *RetryPolicy
}

// ServiceResilienceConfig contains configuration for all operations of a service.
type ServiceResilienceConfig struct {
	ServiceName string
	Operations  []*OperationConfig

	DefaultCircuitBreaker  *CircuitBreakerConfig
	DefaultBulkhead        *BulkheadConfig
	DefaultAdaptiveTimeout *AdaptiveTimeoutConfig
	DefaultRetryPolicy     *RetryPolicy
}

// OperationResilienceBuilder builds resilience wrappers for operations.
type OperationResilienceBuilder struct {
	config *ServiceResilienceConfig
	cache  sync.Map
}

// NewOperationResilienceBuilder creates a new operation resilience builder
// with the given service configuration.
func NewOperationResilienceBuilder(config *ServiceResilienceConfig) *OperationResilienceBuilder {
	return &OperationResilienceBuilder{config: config}
}

func (b *OperationResilienceBuilder) resolveServiceName(op *OperationConfig) string {
	if op != nil && op.ServiceName != "" {
		return op.ServiceName
	}
	if b.config != nil {
		return b.config.ServiceName
	}
	return ""
}

// BuildWrapper builds a resilience wrapper for the given service and operation name.
// Wrappers are cached to avoid recreating CircuitBreaker and other stateful components.
// The cache key includes the service name to prevent cross-service circuit breaker sharing.
func (b *OperationResilienceBuilder) BuildWrapper(serviceName, operationName string) *ResilienceWrapper {
	cacheKey := serviceName + "_" + operationName
	if v, ok := b.cache.Load(cacheKey); ok {
		return v.(*ResilienceWrapper)
	}

	if b.config == nil {
		b.config = &ServiceResilienceConfig{}
	}
	var opConfig *OperationConfig
	for _, op := range b.config.Operations {
		if op.OperationName == operationName {
			opConfig = op
			break
		}
	}

	if opConfig == nil {
		opConfig = &OperationConfig{
			ServiceName:   serviceName,
			OperationName: operationName,
		}
	}

	resolvedService := opConfig.ServiceName
	if resolvedService == "" {
		resolvedService = serviceName
	}
	if resolvedService == "" {
		resolvedService = b.config.ServiceName
	}

	cb := b.buildCircuitBreaker(opConfig)
	bh := b.buildBulkhead(opConfig)
	at := b.buildAdaptiveTimeout(opConfig)
	rp := b.buildRetryPolicy(opConfig)

	wrapper := NewResilienceWrapper(
		fmt.Sprintf("%s_%s", resolvedService, operationName),
		cb,
		bh,
		at,
		rp,
	)
	b.cache.Store(cacheKey, wrapper)
	return wrapper
}

// BuildAllWrappers builds resilience wrappers for all operations configured
// in the service resilience configuration.
func (b *OperationResilienceBuilder) BuildAllWrappers() map[string]*ResilienceWrapper {
	wrappers := make(map[string]*ResilienceWrapper)
	for _, op := range b.config.Operations {
		wrappers[op.OperationName] = b.BuildWrapper(op.ServiceName, op.OperationName)
	}
	return wrappers
}

func (b *OperationResilienceBuilder) buildCircuitBreaker(op *OperationConfig) *CircuitBreaker {
	var config *CircuitBreakerConfig
	if op != nil {
		config = op.CircuitBreaker
	}
	if config == nil && b.config != nil {
		config = b.config.DefaultCircuitBreaker
	}
	serviceName := b.resolveServiceName(op)
	if config == nil {
		config = &CircuitBreakerConfig{
			Name:             fmt.Sprintf("%s_%s", serviceName, op.OperationName),
			MaxFailures:      DefaultCircuitBreakerMaxFailures,
			ResetTimeout:     DefaultCircuitBreakerResetTimeout,
			HalfOpenMaxCalls: DefaultCircuitBreakerHalfOpenRequests,
			SuccessThreshold: DefaultCircuitBreakerSuccessThreshold,
		}
	}
	return NewCircuitBreaker(config)
}

func (b *OperationResilienceBuilder) buildBulkhead(op *OperationConfig) *Bulkhead {
	var config *BulkheadConfig
	if op != nil {
		config = op.Bulkhead
	}
	if config == nil && b.config != nil {
		config = b.config.DefaultBulkhead
	}
	serviceName := b.resolveServiceName(op)
	if config == nil {
		config = &BulkheadConfig{
			Name:          fmt.Sprintf("%s_%s", serviceName, op.OperationName),
			MaxConcurrent: DefaultBulkheadMaxConcurrent,
			MaxWait:       DefaultBulkheadTimeout,
		}
	}
	return NewBulkhead(config)
}

func (b *OperationResilienceBuilder) buildAdaptiveTimeout(op *OperationConfig) *AdaptiveTimeout {
	var config *AdaptiveTimeoutConfig
	if op != nil {
		config = op.AdaptiveTimeout
	}
	if config == nil && b.config != nil {
		config = b.config.DefaultAdaptiveTimeout
	}
	serviceName := b.resolveServiceName(op)
	if config == nil {
		config = &AdaptiveTimeoutConfig{
			Name:             fmt.Sprintf("%s_%s", serviceName, op.OperationName),
			InitialTimeout:   DefaultAdaptiveTimeoutDefault,
			MinTimeout:       DefaultAdaptiveTimeoutMin,
			MaxTimeout:       DefaultAdaptiveTimeoutMax,
			SuccessThreshold: DefaultAdaptiveTimeoutSuccessThreshold,
			FailureThreshold: DefaultAdaptiveTimeoutFailureThreshold,
			AdjustmentFactor: DefaultAdaptiveAdjustmentFactor,
		}
	}
	return NewAdaptiveTimeout(config)
}

func (b *OperationResilienceBuilder) buildRetryPolicy(op *OperationConfig) *RetryPolicy {
	var config *RetryPolicy
	if op != nil {
		config = op.RetryPolicy
	}
	if config == nil && b.config != nil {
		config = b.config.DefaultRetryPolicy
	}

	if config == nil {
		rp := NewRetryPolicy()
		rp.SetMaxAttempts(DefaultRetryMaxAttempts)
		rp.SetInitialDelay(DefaultRetryInitialDelay)
		rp.SetMaxDelay(DefaultRetryMaxDelay)
		rp.SetBackoff(&ExponentialBackoff{
			InitialDelay: DefaultRetryInitialDelay,
			MaxDelay:     DefaultRetryMaxDelay,
			Multiplier:   DefaultRetryBackoffMultiplier,
		})
		rp.SetJitter(&EqualJitter{})
		return rp
	}

	return config
}
