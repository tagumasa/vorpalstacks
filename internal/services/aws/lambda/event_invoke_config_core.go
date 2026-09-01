package lambda

import (
	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// EventInvokeConfigInput carries the wire members of a Put or Update
// FunctionEventInvokeConfig request. The Has* flags distinguish an
// explicitly provided value from an omitted member so that Put can apply
// the AWS defaults and Update can preserve the stored value.
type EventInvokeConfigInput struct {
	FunctionName string
	Qualifier    string

	HasMaximumEventAgeInSeconds bool
	MaximumEventAgeInSeconds    int32
	HasMaximumRetryAttempts     bool
	MaximumRetryAttempts        int32

	DestinationConfig *lambdastore.DestinationConfig
}

// putFunctionEventInvokeConfigCore creates or replaces the asynchronous
// invocation configuration for a function qualifier. Put replaces the whole
// configuration; omitted members fall back to the AWS defaults rather than
// zero.
func (s *LambdaService) putFunctionEventInvokeConfigCore(reqCtx *request.RequestContext, in *EventInvokeConfigInput) (*lambdastore.EventInvokeConfig, error) {
	if in.FunctionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	config := &lambdastore.EventInvokeConfig{}

	if in.HasMaximumEventAgeInSeconds {
		if err := validateMaximumEventAgeInSeconds(in.MaximumEventAgeInSeconds); err != nil {
			return nil, err
		}
		config.MaximumEventAgeInSeconds = in.MaximumEventAgeInSeconds
	} else {
		config.MaximumEventAgeInSeconds = lambdastore.DefaultMaximumEventAgeInSeconds
	}
	if in.HasMaximumRetryAttempts {
		if err := validateMaximumRetryAttempts(in.MaximumRetryAttempts); err != nil {
			return nil, err
		}
		config.MaximumRetryAttempts = in.MaximumRetryAttempts
	} else {
		config.MaximumRetryAttempts = lambdastore.DefaultMaximumRetryAttempts
	}
	config.DestinationConfig = in.DestinationConfig

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.Functions.SetEventInvokeConfig(in.FunctionName, in.Qualifier, config); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return config, nil
}

// getFunctionEventInvokeConfigCore retrieves the asynchronous invocation
// configuration for a function qualifier.
func (s *LambdaService) getFunctionEventInvokeConfigCore(stores *lambdaStore, functionName, qualifier string) (*lambdastore.EventInvokeConfig, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	config, err := stores.Functions.GetEventInvokeConfig(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrEventInvokeConfigNotFound || err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return config, nil
}

// deleteFunctionEventInvokeConfigCore deletes the asynchronous invocation
// configuration for a function qualifier.
func (s *LambdaService) deleteFunctionEventInvokeConfigCore(stores *lambdaStore, functionName, qualifier string) error {
	if functionName == "" {
		return NewInvalidParameter("FunctionName", "Function name is required")
	}
	if err := stores.Functions.DeleteEventInvokeConfig(functionName, qualifier); err != nil {
		if err == lambdastore.ErrEventInvokeConfigNotFound || err == lambdastore.ErrFunctionNotFound {
			return ErrResourceNotFound
		}
		return err
	}
	return nil
}

// listFunctionEventInvokeConfigsCore lists every asynchronous invocation
// configuration of a function.
func (s *LambdaService) listFunctionEventInvokeConfigsCore(stores *lambdaStore, functionName string) ([]lambdastore.EventInvokeConfig, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	configs, err := stores.Functions.ListEventInvokeConfigs(functionName)
	if err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return configs, nil
}

// updateFunctionEventInvokeConfigCore updates the asynchronous invocation
// configuration for a function qualifier. Only fields provided in the input
// are modified; existing values for unprovided fields are preserved.
// DestinationConfig is replaced atomically when provided. If no config
// exists for the qualifier, a new one is created.
func (s *LambdaService) updateFunctionEventInvokeConfigCore(stores *lambdaStore, in *EventInvokeConfigInput) (*lambdastore.EventInvokeConfig, error) {
	if in.FunctionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	// Attempt to load existing config; if not found, create a blank one.
	config, err := stores.Functions.GetEventInvokeConfig(in.FunctionName, in.Qualifier)
	if err != nil && err != lambdastore.ErrEventInvokeConfigNotFound && err != lambdastore.ErrFunctionNotFound {
		return nil, err
	}
	if config == nil {
		// Creating the configuration through an update: fields the update
		// does not specify take the AWS defaults.
		config = &lambdastore.EventInvokeConfig{
			MaximumEventAgeInSeconds: lambdastore.DefaultMaximumEventAgeInSeconds,
			MaximumRetryAttempts:     lambdastore.DefaultMaximumRetryAttempts,
		}
	}

	// Overwrite only the fields that were explicitly provided in the request.
	if in.HasMaximumEventAgeInSeconds {
		if err := validateMaximumEventAgeInSeconds(in.MaximumEventAgeInSeconds); err != nil {
			return nil, err
		}
		config.MaximumEventAgeInSeconds = in.MaximumEventAgeInSeconds
	}
	if in.HasMaximumRetryAttempts {
		if err := validateMaximumRetryAttempts(in.MaximumRetryAttempts); err != nil {
			return nil, err
		}
		config.MaximumRetryAttempts = in.MaximumRetryAttempts
	}
	if in.DestinationConfig != nil {
		config.DestinationConfig = in.DestinationConfig
	}

	if err := stores.Functions.SetEventInvokeConfig(in.FunctionName, in.Qualifier, config); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return config, nil
}
