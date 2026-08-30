package lambda

// This file carries the invocation-plane Core helpers: the container
// lifecycle that persists container IDs through the function store and the
// asynchronous retry engine that reads the function's event-invoke
// configuration. Relocated verbatim from service.go and async_invoker.go
// so that handler closures no longer reach the store packages directly.

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

func (s *LambdaService) getRuntimeImage(runtime lambdastore.Runtime) string {
	return lambdastore.GetImageForRuntime(runtime)
}

func (s *LambdaService) ensureFunctionContainer(function *lambdastore.Function, ver *lambdastore.Version, store *lambdastore.FunctionStore, region string) (string, error) {
	ctx := context.Background()

	version := "$LATEST"
	if ver != nil {
		version = ver.Version
	}

	containerName := fmt.Sprintf("lambda-%s-%s-%s", region, function.FunctionName, sanitizeForContainerName(version))

	containerID := function.ContainerID
	if ver != nil && ver.ContainerID != "" {
		containerID = ver.ContainerID
	}

	if containerID != "" {
		status, err := s.dockerClient.GetContainerStatus(ctx, containerID)
		if err == nil && status == mobyclient.ContainerStatusRunning {
			return containerID, nil
		}
	}

	// Serialise creation per container name. Waiters re-check the cache
	// after acquiring: the goroutine that held the lock created the
	// container and recorded its ID.
	muAny, _ := s.containerEnsureMu.LoadOrStore(containerName, &sync.Mutex{})
	mu := muAny.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	if hintAny, ok := s.containerIDs.Load(containerName); ok {
		if hint, _ := hintAny.(string); hint != "" {
			if status, err := s.dockerClient.GetContainerStatus(ctx, hint); err == nil && status == mobyclient.ContainerStatusRunning {
				return hint, nil
			}
		}
	}

	execCfg := executionConfigFor(function, ver)

	image := s.getRuntimeImage(execCfg.Runtime)
	if execCfg.ImageUri != "" {
		image = execCfg.ImageUri
	}

	envVars := map[string]string{
		"AWS_LAMBDA_FUNCTION_TIMEOUT":     fmt.Sprintf("%d", execCfg.Timeout),
		"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": fmt.Sprintf("%d", execCfg.MemorySize),
		"AWS_LAMBDA_FUNCTION_HANDLER":     execCfg.Handler,
		"AWS_LAMBDA_FUNCTION_NAME":        function.FunctionName,
		"AWS_LAMBDA_FUNCTION_VERSION":     version,
		"AWS_REGION":                      region,
	}

	if execCfg.Environment != nil {
		for k, v := range execCfg.Environment.Variables {
			envVars[k] = v
		}
	}

	if _, ok := envVars["AWS_ENDPOINT_URL"]; !ok && s.hostEndpoint != "" {
		envVars["AWS_ENDPOINT_URL"] = s.hostEndpoint
	}

	cfg := mobyclient.AdvancedContainerConfig{
		Name:       containerName,
		Image:      image,
		PullImage:  true,
		Env:        envVars,
		Entrypoint: []string{"/lambda-entrypoint.sh"},
		Cmd:        []string{execCfg.Handler},
		Network:    "bridge",
		AutoRemove: false,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}

	result, err := s.dockerClient.CreateContainerFromConfig(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.dockerClient.StartContainer(ctx, result.ID); err != nil {
		if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
			logs.Warn("Failed to remove container after start failure", logs.String("containerID", result.ID), logs.Err(rmErr))
		}
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	if version == "$LATEST" {
		originalContainerID := function.ContainerID
		function.ContainerID = result.ID
		function.ContainerImageID = result.ID
		if err := store.Update(function); err != nil {
			function.ContainerID = originalContainerID
			function.ContainerImageID = originalContainerID
			if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
				logs.Warn("Failed to remove container during rollback", logs.String("containerID", result.ID), logs.Err(rmErr))
			}
			return "", fmt.Errorf("failed to update function: %w", err)
		}
		// Clean up the previous container to prevent resource leaks when the
		// function is updated or the old container has crashed.
		if originalContainerID != "" && originalContainerID != result.ID {
			if rmErr := s.dockerClient.RemoveContainer(ctx, originalContainerID, true); rmErr != nil {
				logs.Warn("Failed to remove previous container", logs.String("containerID", originalContainerID), logs.Err(rmErr))
			}
		}
	} else if ver != nil {
		originalVerContainerID := ver.ContainerID
		ver.ContainerID = result.ID
		ver.ContainerImageID = result.ID
		if err := store.Update(function); err != nil {
			ver.ContainerID = originalVerContainerID
			ver.ContainerImageID = originalVerContainerID
			if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
				logs.Warn("Failed to remove container during rollback", logs.String("containerID", result.ID), logs.Err(rmErr))
			}
			return "", fmt.Errorf("failed to update function version: %w", err)
		}
		// Clean up the previous version container to prevent resource leaks.
		if originalVerContainerID != "" && originalVerContainerID != result.ID {
			if rmErr := s.dockerClient.RemoveContainer(ctx, originalVerContainerID, true); rmErr != nil {
				logs.Warn("Failed to remove previous version container", logs.String("containerID", originalVerContainerID), logs.Err(rmErr))
			}
		}
	}

	s.containerIDs.Store(containerName, result.ID)

	return result.ID, nil
}

// invokeAsyncWithRetry executes an asynchronous Lambda invocation with
// retry and destination delivery.  It reads the function's
// EventInvokeConfig (if any) to determine MaximumRetryAttempts,
// MaximumEventAgeInSeconds, and DestinationConfig. The defaults come from
// the store package constants shared with the event-invoke-config API.
//
// On success: if OnSuccess destination is configured, the invocation
// record is delivered to it.
// On failure after all retries: if OnFailure destination is configured,
// the invocation record is delivered to it.
func (s *LambdaService) invokeAsyncWithRetry(
	ctx context.Context,
	function *lambdastore.Function,
	ver *lambdastore.Version,
	store *lambdastore.FunctionStore,
	region string,
	payload []byte,
	qualifier string,
) {
	// Load EventInvokeConfig for this function+qualifier.
	maxRetries := int(lambdastore.DefaultMaximumRetryAttempts)
	maxEventAge := lambdastore.DefaultMaximumEventAgeInSeconds
	var destConfig *lambdastore.DestinationConfig

	if store != nil {
		eic, err := store.GetEventInvokeConfig(function.FunctionName, qualifier)
		if err == nil && eic != nil {
			if eic.MaximumRetryAttempts >= 0 {
				maxRetries = int(eic.MaximumRetryAttempts)
			}
			if eic.MaximumEventAgeInSeconds > 0 {
				maxEventAge = eic.MaximumEventAgeInSeconds
			}
			destConfig = eic.DestinationConfig
		}
	}

	startTime := time.Now()
	var lastResult *lambdastore.InvocationResult
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check MaximumEventAge before retrying.
		if attempt > 0 {
			elapsed := int32(time.Since(startTime).Seconds())
			if elapsed >= maxEventAge {
				break
			}
			// Exponential backoff: 1s, 2s, 4s... capped at 60s.
			backoff := time.Duration(1) << uint(attempt-1)
			if backoff > 60 {
				backoff = 60
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff * time.Second):
			}
		}

		result, err := s.invokeFunction(function, ver, store, region, payload, "")
		lastResult = result
		lastErr = err

		// Determine success: no infrastructure error and no function error.
		if err == nil && result != nil && result.FunctionError == "" {
			// Success — deliver to OnSuccess destination if configured.
			if destConfig != nil && destConfig.OnSuccess != nil && destConfig.OnSuccess.Destination != "" {
				deliverDestination(ctx, s, destConfig.OnSuccess.Destination, true,
					payload, result, function, region, attempt+1)
			}
			return
		}

		logs.Warn("async invocation failed, will retry",
			logs.String("function", function.FunctionName),
			logs.Int("attempt", attempt+1),
			logs.Int("maxRetries", maxRetries+1),
			logs.Err(err))
	}

	// All retries exhausted — deliver to OnFailure destination if configured.
	if destConfig != nil && destConfig.OnFailure != nil && destConfig.OnFailure.Destination != "" {
		errPayload := []byte(`{"errorMessage":"internal error"}`)
		if lastErr != nil {
			errPayload = []byte(fmt.Sprintf(`{"errorMessage":%q}`, lastErr.Error()))
		} else if lastResult != nil && lastResult.FunctionError != "" {
			errPayload = lastResult.Payload
		}
		deliverDestination(ctx, s, destConfig.OnFailure.Destination, false,
			payload, &lambdastore.InvocationResult{
				StatusCode:    200,
				Payload:       errPayload,
				FunctionError: "Unhandled",
			}, function, region, maxRetries+1)
	}
}
