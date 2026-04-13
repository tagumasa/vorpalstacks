// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Compile-time interface compliance checks for StreamableResponse and StatusCodeResponse.
var (
	_ response.StreamableResponse = (*lambdaInvokeResponse)(nil)
	_ response.StatusCodeResponse = (*lambdaInvokeResponse)(nil)

	_ response.StreamableResponse = (*invokeWithResponseStreamResponse)(nil)
	_ response.StatusCodeResponse = (*invokeWithResponseStreamResponse)(nil)

	_ response.StreamableResponse = (*invokeAsyncResponse)(nil)
	_ response.StatusCodeResponse = (*invokeAsyncResponse)(nil)
)

// Invoke synchronously invokes a Lambda function with the given payload.
// Returns the function output, status code, and executed version.
// If InvocationType is "Event", the function is invoked asynchronously and HTTP 202 is returned.
func (s *LambdaService) Invoke(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, ver, _, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	invocationType := request.GetStringParam(req.Parameters, "InvocationType")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if invocationType == "Event" {
		region := reqCtx.GetRegion()
		functionCopy := deepCopyFunction(function)
		s.asyncWg.Add(1)
		go func() {
			defer s.asyncWg.Done()
			defer func() {
				if r := recover(); r != nil {
					logs.Error("Invoke Event panic", logs.String("function", functionCopy.FunctionName), logs.Any("panic", r))
				}
			}()
			select {
			case <-ctx.Done():
				logs.Info("Invoke Event cancelled", logs.String("function", functionCopy.FunctionName))
				return
			default:
			}
			asyncStore := s.getOrCreateFunctionStore(region)
			if _, err := s.invokeFunction(functionCopy, nil, asyncStore, region, payload); err != nil {
				logs.Error("Invoke Event failed", logs.String("function", functionCopy.FunctionName), logs.Err(err))
			}
		}()
		return &lambdaInvokeResponse{result: &lambdastore.InvocationResult{StatusCode: 202}}, nil
	}

	result, err := s.invokeFunction(function, ver, store.Functions, reqCtx.GetRegion(), payload)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("invocation returned nil result")
	}

	return &lambdaInvokeResponse{result: result}, nil
}

// lambdaInvokeResponse wraps an invocation result to satisfy the StreamableResponse interface.
type lambdaInvokeResponse struct {
	result *lambdastore.InvocationResult
}

// GetStream returns the invocation payload as an io.Reader.
func (r *lambdaInvokeResponse) GetStream() io.Reader {
	if r.result == nil || len(r.result.Payload) == 0 {
		return bytes.NewReader(nil)
	}
	return bytes.NewReader(r.result.Payload)
}

// GetStreamHeaders returns the Lambda invocation response headers.
func (r *lambdaInvokeResponse) GetStreamHeaders() http.Header {
	h := make(http.Header)
	if r.result != nil {
		if r.result.ExecutedVersion != "" {
			h.Set("x-amz-executed-version", r.result.ExecutedVersion)
		}
		if r.result.FunctionError != "" {
			h.Set("x-amz-function-error", r.result.FunctionError)
		}
		if r.result.LogResult != "" {
			h.Set("x-amz-log-result", r.result.LogResult)
		}
	}
	return h
}

// GetStreamStatusCode returns the HTTP status code for the invocation response.
func (r *lambdaInvokeResponse) GetStreamStatusCode() int {
	if r.result == nil {
		return http.StatusOK
	}
	return int(r.result.StatusCode)
}

// InvokeWithResponseStream invokes a Lambda function with response streaming.
func (s *LambdaService) InvokeWithResponseStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, ver, _, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.invokeFunction(function, ver, store.Functions, reqCtx.GetRegion(), payload)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("invocation returned nil result")
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		w := NewInvokeResponseStreamWriter(pw)

		if result.FunctionError != "" {
			if err := w.WriteInvokeCompleteError(result.FunctionError, string(result.Payload)); err != nil {
				logs.Error("InvokeWithResponseStream: failed to write error InvokeComplete event", logs.Err(err))
			}
			return
		}

		if len(result.Payload) > 0 {
			if err := w.WritePayloadChunk(result.Payload); err != nil {
				logs.Error("InvokeWithResponseStream: failed to write PayloadChunk event", logs.Err(err))
				return
			}
		}

		if err := w.WriteInvokeComplete(int(result.StatusCode), result.ExecutedVersion, result.LogResult); err != nil {
			logs.Error("InvokeWithResponseStream: failed to write InvokeComplete event", logs.Err(err))
		}
	}()

	return &invokeWithResponseStreamResponse{
		reader:          pr,
		executedVersion: result.ExecutedVersion,
		statusCode:      int(result.StatusCode),
		functionError:   result.FunctionError,
		contentType:     "application/vnd.amazon.eventstream",
	}, nil
}

type invokeWithResponseStreamResponse struct {
	reader          io.Reader
	executedVersion string
	statusCode      int
	functionError   string
	contentType     string
}

// GetStream returns the response stream payload for a Lambda invocation
// with response streaming enabled.
func (r *invokeWithResponseStreamResponse) GetStream() io.Reader {
	return r.reader
}

// GetStreamHeaders returns the HTTP headers for a Lambda invocation
// with response streaming enabled.
func (r *invokeWithResponseStreamResponse) GetStreamHeaders() http.Header {
	h := make(http.Header)
	h.Set("Content-Type", r.contentType)
	h.Set("x-amz-executed-version", r.executedVersion)
	if r.functionError != "" {
		h.Set("x-amz-function-error", r.functionError)
	}
	return h
}

// GetStreamStatusCode returns the HTTP status code for a Lambda invocation
// with response streaming enabled.
func (r *invokeWithResponseStreamResponse) GetStreamStatusCode() int {
	if r.statusCode == 0 {
		return http.StatusOK
	}
	return r.statusCode
}

// InvokeAsync asynchronously invokes a Lambda function with the given payload.
// The function is invoked in the background and returns immediately with HTTP 202.
func (s *LambdaService) InvokeAsync(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	region := reqCtx.GetRegion()

	functionCopy := deepCopyFunction(function)

	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("InvokeAsync panic", logs.String("function", functionCopy.FunctionName), logs.Any("panic", r))
			}
		}()
		select {
		case <-ctx.Done():
			logs.Info("InvokeAsync cancelled", logs.String("function", functionCopy.FunctionName))
			return
		default:
		}
		asyncStore := s.getOrCreateFunctionStore(region)
		if _, err := s.invokeFunction(functionCopy, nil, asyncStore, region, payload); err != nil {
			logs.Error("InvokeAsync failed", logs.String("function", functionCopy.FunctionName), logs.Err(err))
		}
	}()

	return &invokeAsyncResponse{Status: 202}, nil
}

// invokeAsyncResponse represents the response from an asynchronous Lambda invocation.
type invokeAsyncResponse struct {
	Status int
}

// GetStreamStatusCode returns the HTTP status code for the asynchronous invocation response.
func (r *invokeAsyncResponse) GetStreamStatusCode() int {
	return r.Status
}

// GetStream returns an empty reader for the asynchronous invocation response.
func (r *invokeAsyncResponse) GetStream() io.Reader {
	return bytes.NewReader(nil)
}

// GetStreamHeaders returns empty headers for the asynchronous invocation response.
func (r *invokeAsyncResponse) GetStreamHeaders() http.Header {
	return http.Header{}
}

// PublishVersion creates a new version of the Lambda function from the current $LATEST version.
// The new version is immutable and can be used for deployments.
func (s *LambdaService) PublishVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "Description")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	version, err := store.Functions.PublishVersion(function, description)
	if err != nil {
		return nil, err
	}

	return s.toVersionConfiguration(version), nil
}

// ListVersionsByFunction returns all versions of a Lambda function,
// including the $LATEST version and all published versions.
func (s *LambdaService) ListVersionsByFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 50
	}

	allVersions := make([]map[string]interface{}, 0)
	latestConfig := s.toFunctionConfiguration(function)
	latestConfig["Version"] = "$LATEST"
	allVersions = append(allVersions, latestConfig)

	for i := range function.Versions {
		allVersions = append(allVersions, s.toVersionConfiguration(&function.Versions[i]))
	}

	pageResult := pagination.PaginateSlice(allVersions, marker, maxItems, func(v map[string]interface{}) string {
		if version, ok := v["Version"].(string); ok {
			return version
		}
		return ""
	})

	versions := make([]interface{}, 0, len(pageResult.Items))
	for _, v := range pageResult.Items {
		versions = append(versions, v)
	}

	response := map[string]interface{}{
		"Versions": versions,
	}
	if pageResult.IsTruncated {
		response["NextMarker"] = pageResult.NextMarker
	}

	return response, nil
}

// CreateAlias creates an alias for a Lambda function.
// An alias points to a specific version and can be used for traffic shifting.
func (s *LambdaService) CreateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	aliasName := request.GetStringParam(req.Parameters, "Name")
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}

	functionVersion := request.GetStringParam(req.Parameters, "FunctionVersion")
	if functionVersion == "" {
		functionVersion = "$LATEST"
	}

	if functionVersion != "$LATEST" {
		versionExists := false
		for _, v := range function.Versions {
			if v.Version == functionVersion {
				versionExists = true
				break
			}
		}
		if !versionExists {
			return nil, NewResourceNotFound("FunctionVersion", functionVersion)
		}
	}

	alias := &lambdastore.Alias{
		Name:            aliasName,
		FunctionVersion: functionVersion,
		Description:     request.GetStringParam(req.Parameters, "Description"),
	}

	if routingMap := request.GetMapParam(req.Parameters, "RoutingConfig"); routingMap != nil {
		if weights, ok := routingMap["AdditionalVersionWeights"].(map[string]interface{}); ok {
			alias.RoutingConfig = &lambdastore.RoutingConfig{
				AdditionalVersionWeights: make(map[string]float64),
			}
			for version, weight := range weights {
				switch w := weight.(type) {
				case float64:
					alias.RoutingConfig.AdditionalVersionWeights[version] = w
				case int:
					alias.RoutingConfig.AdditionalVersionWeights[version] = float64(w)
				}
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.Functions.CreateAliasAtomically(function.FunctionName, func(fn *lambdastore.Function) (*lambdastore.Alias, error) {
		return alias, nil
	})
	if err != nil {
		if errors.Is(err, lambdastore.ErrAliasAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("Alias already exists: %s", alias.Name))
		}
		return nil, err
	}

	return s.toAliasResponse(created), nil
}

// DeleteAlias deletes an alias from a Lambda function.
func (s *LambdaService) DeleteAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName = extractFunctionName(functionName)

	aliasName := request.GetStringParam(req.Parameters, "Name")
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Functions.DeleteAlias(functionName, aliasName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetAlias retrieves the configuration of an alias for a Lambda function.
func (s *LambdaService) GetAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName = extractFunctionName(functionName)

	aliasName := request.GetStringParam(req.Parameters, "Name")
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	alias, err := store.Functions.GetAlias(functionName, aliasName)
	if err != nil {
		return nil, NewResourceNotFound("Alias", aliasName)
	}

	return s.toAliasResponse(alias), nil
}

// UpdateAlias updates the configuration of an alias for a Lambda function.
// Allows modifying the function version, description, and routing configuration.
func (s *LambdaService) UpdateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName = extractFunctionName(functionName)

	aliasName := request.GetStringParam(req.Parameters, "Name")
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "Description")
	functionVersion := request.GetStringParam(req.Parameters, "FunctionVersion")

	var routingConfig *lambdastore.RoutingConfig
	if routingMap := request.GetMapParam(req.Parameters, "RoutingConfig"); routingMap != nil {
		if weights, ok := routingMap["AdditionalVersionWeights"].(map[string]interface{}); ok {
			routingConfig = &lambdastore.RoutingConfig{
				AdditionalVersionWeights: make(map[string]float64),
			}
			for version, weight := range weights {
				switch w := weight.(type) {
				case float64:
					routingConfig.AdditionalVersionWeights[version] = w
				case int:
					routingConfig.AdditionalVersionWeights[version] = float64(w)
				}
			}
		}
	}

	alias, err := store.Functions.UpdateAliasAtomically(functionName, aliasName, func(fn *lambdastore.Function, existing *lambdastore.Alias) error {
		if description != "" {
			existing.Description = description
		}
		if functionVersion != "" {
			existing.FunctionVersion = functionVersion
		}
		if routingConfig != nil {
			existing.RoutingConfig = routingConfig
		}
		return nil
	})
	if err != nil {
		if err == lambdastore.ErrAliasNotFound {
			return nil, NewResourceNotFound("Alias", aliasName)
		}
		return nil, err
	}

	return s.toAliasResponse(alias), nil
}

// ListAliases returns all aliases for a Lambda function.
func (s *LambdaService) ListAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName = extractFunctionName(functionName)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	function, err := store.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	aliases := make([]interface{}, 0)
	for _, a := range function.Aliases {
		aliases = append(aliases, s.toAliasResponse(&a))
	}

	return map[string]interface{}{
		"Aliases": aliases,
	}, nil
}

func (s *LambdaService) toAliasResponse(a *lambdastore.Alias) map[string]interface{} {
	resp := map[string]interface{}{
		"AliasArn":        a.AliasArn,
		"Name":            a.Name,
		"FunctionVersion": a.FunctionVersion,
		"Description":     a.Description,
		"FunctionName":    a.FunctionName,
		"RevisionId":      a.RevisionId,
	}
	if a.RoutingConfig != nil && len(a.RoutingConfig.AdditionalVersionWeights) > 0 {
		resp["RoutingConfig"] = map[string]interface{}{
			"AdditionalVersionWeights": a.RoutingConfig.AdditionalVersionWeights,
		}
	}
	return resp
}
