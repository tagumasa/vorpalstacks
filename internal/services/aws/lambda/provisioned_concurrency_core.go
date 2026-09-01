package lambda

import (
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ProvisionedConcurrencyInput carries the wire members of the
// provisioned-concurrency operations. The function name and qualifier
// arrive already resolved from their wire reference forms by the handler.
type ProvisionedConcurrencyInput struct {
	FunctionName                    string
	Qualifier                       string
	ProvisionedConcurrentExecutions int32
	Region                          string
}

// putProvisionedConcurrencyCore configures provisioned concurrency for a
// function's published version or alias and pre-warms its container.
func (s *LambdaService) putProvisionedConcurrencyCore(reqCtx *request.RequestContext, in *ProvisionedConcurrencyInput) (*lambdastore.ProvisionedConcurrencyConfig, error) {
	if in.FunctionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if in.Qualifier == "" {
		return nil, NewInvalidParameter("Qualifier", "Qualifier is required")
	}
	if in.ProvisionedConcurrentExecutions < 1 {
		return nil, NewInvalidParameter("ProvisionedConcurrentExecutions", "Provisioned concurrent executions must be at least 1")
	}

	// Provisioned concurrency applies to a published version or alias, not
	// to $LATEST.
	if in.Qualifier == "$LATEST" {
		return nil, NewInvalidParameter("Qualifier", "Provisioned concurrency cannot be configured on the $LATEST version. Publish a version or use an alias.")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := stores.Functions.Get(in.FunctionName); err != nil {
		return nil, mapStoreError(err)
	}
	if _, version, alias, err := stores.Functions.ResolveQualifier(in.FunctionName, in.Qualifier); err != nil || (version == nil && alias == nil) {
		return nil, NewResourceNotFound("Qualifier", in.Qualifier)
	}

	if err := stores.Functions.SetProvisionedConcurrency(in.FunctionName, in.Qualifier, in.ProvisionedConcurrentExecutions); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	// Pre-warm the function container for the resolved qualifier.  This
	// eliminates cold-start latency on the first invocation.
	if s.dockerClient != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logs.Error("Panic in provisioned concurrency pre-warm goroutine",
						logs.String("function", in.FunctionName),
						logs.String("qualifier", in.Qualifier),
						logs.Any("panic", r))
				}
			}()
			fn, err := stores.Functions.Get(in.FunctionName)
			if err != nil {
				return
			}
			_, version, alias, err := stores.Functions.ResolveQualifier(in.FunctionName, in.Qualifier)
			if err != nil {
				return
			}
			ver := version
			if alias != nil {
				ver = resolveAliasTargetVersion(fn, alias)
			}
			_, _ = s.ensureFunctionContainer(fn, ver, stores.Functions, in.Region)
		}()
	}

	config, err := stores.Functions.GetProvisionedConcurrency(in.FunctionName, in.Qualifier)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// getProvisionedConcurrencyCore retrieves the provisioned concurrency
// configuration for a function qualifier.
func (s *LambdaService) getProvisionedConcurrencyCore(stores *lambdaStore, functionName, qualifier string) (*lambdastore.ProvisionedConcurrencyConfig, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if qualifier == "" {
		return nil, NewInvalidParameter("Qualifier", "Qualifier is required")
	}
	config, err := stores.Functions.GetProvisionedConcurrency(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrProvisionedConcurrencyNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return config, nil
}

// deleteProvisionedConcurrencyCore removes the provisioned concurrency
// configuration for a function qualifier.
func (s *LambdaService) deleteProvisionedConcurrencyCore(stores *lambdaStore, functionName, qualifier string) error {
	if functionName == "" {
		return NewInvalidParameter("FunctionName", "Function name is required")
	}
	if qualifier == "" {
		return NewInvalidParameter("Qualifier", "Qualifier is required")
	}
	if err := stores.Functions.DeleteProvisionedConcurrency(functionName, qualifier); err != nil {
		if err == lambdastore.ErrProvisionedConcurrencyNotFound {
			return ErrResourceNotFound
		}
		return err
	}
	return nil
}

// listProvisionedConcurrencyConfigsCore lists the provisioned concurrency
// configurations of a function.
func (s *LambdaService) listProvisionedConcurrencyConfigsCore(stores *lambdaStore, functionName string) ([]lambdastore.ProvisionedConcurrencyConfig, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	configs, err := stores.Functions.ListProvisionedConcurrency(functionName)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
