// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
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

// Payload size limits per AWS Lambda specification.
const (
	syncMaxPayloadSize  = 6 * 1024 * 1024 // 6 MB for synchronous invocation
	asyncMaxPayloadSize = 256 * 1024      // 256 KB for asynchronous invocation
)

// Invoke synchronously invokes a Lambda function with the given payload.
// Returns the function output, status code, and executed version.
// If InvocationType is "Event", the function is invoked asynchronously and HTTP 202 is returned.
func (s *LambdaService) Invoke(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, ver, alias, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Resolve alias to its target version, applying weighted routing if
	// RoutingConfig is configured.  When no alias is provided, ver is
	// already the correct version (or nil for $LATEST).
	if alias != nil {
		ver = resolveAliasTargetVersion(function, alias)
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	invocationType := request.GetStringParam(req.Parameters, "InvocationType")
	if err := validateInvocationType(invocationType); err != nil {
		return nil, err
	}

	// DryRun: AWS returns 204 without invoking the function. The
	// caller has already been authorised by this point, so we simply
	// acknowledge the request.
	if invocationType == "DryRun" {
		return &lambdaInvokeResponse{result: &lambdastore.InvocationResult{StatusCode: 204}}, nil
	}

	if invocationType == "Event" {
		if len(payload) > asyncMaxPayloadSize {
			return nil, ErrRequestTooLarge
		}
		region := reqCtx.GetRegion()
		functionCopy := deepCopyFunction(function)
		verCopy := deepCopyVersion(ver)
		qualifier := request.GetStringParam(req.Parameters, "Qualifier")
		if alias != nil {
			qualifier = alias.Name
		}
		s.asyncWg.Add(1)
		go func() {
			defer s.asyncWg.Done()
			defer func() {
				if r := recover(); r != nil {
					logs.Error("Invoke Event panic", logs.String("function", functionCopy.FunctionName), logs.Any("panic", r))
				}
			}()
			// Use a detached context so that retry and destination
			// delivery survive after the HTTP 202 response is sent and
			// the request context is cancelled.
			asyncCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			asyncStore := s.getOrCreateFunctionStore(region)
			s.invokeAsyncWithRetry(asyncCtx, functionCopy, verCopy, asyncStore, region, payload, qualifier)
		}()
		return &lambdaInvokeResponse{result: &lambdastore.InvocationResult{StatusCode: 202}}, nil
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Synchronous invocation payload limit: 6 MB.
	if len(payload) > syncMaxPayloadSize {
		return nil, ErrRequestTooLarge
	}

	logType := request.GetStringParam(req.Parameters, "LogType")
	if err := validateLogType(logType); err != nil {
		return nil, err
	}

	result, err := s.invokeFunction(function, ver, store.Functions, reqCtx.GetRegion(), payload, logType)
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
	function, ver, alias, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if alias != nil {
		ver = resolveAliasTargetVersion(function, alias)
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	if len(payload) > syncMaxPayloadSize {
		return nil, ErrRequestTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.invokeFunction(function, ver, store.Functions, reqCtx.GetRegion(), payload, "")
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("invocation returned nil result")
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		defer func() { recover() }()
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
	function, ver, alias, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if alias != nil {
		ver = resolveAliasTargetVersion(function, alias)
	}

	var payload []byte
	if payloadStr, ok := req.Parameters["Payload"].(string); ok {
		payload = []byte(payloadStr)
	}

	if len(payload) > asyncMaxPayloadSize {
		return nil, ErrRequestTooLarge
	}

	region := reqCtx.GetRegion()

	functionCopy := deepCopyFunction(function)
	var verCopy *lambdastore.Version
	if ver != nil {
		verCopy = deepCopyVersion(ver)
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if alias != nil {
		qualifier = alias.Name
	}

	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("InvokeAsync panic", logs.String("function", functionCopy.FunctionName), logs.Any("panic", r))
			}
		}()
		// Use a detached context so that retry and destination
		// delivery survive after the HTTP 202 response is sent and
		// the request context is cancelled.
		asyncCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		asyncStore := s.getOrCreateFunctionStore(region)
		s.invokeAsyncWithRetry(asyncCtx, functionCopy, verCopy, asyncStore, region, payload, qualifier)
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
	version, err := s.publishVersionWithCode(store, function, description, reqCtx.GetRegion())
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
	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))

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

// parseRoutingConfig parses the RoutingConfig member into the store shape.
func parseRoutingConfig(routingMap map[string]interface{}) *lambdastore.RoutingConfig {
	weights, ok := routingMap["AdditionalVersionWeights"].(map[string]interface{})
	if !ok {
		return nil
	}
	routingConfig := &lambdastore.RoutingConfig{
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
	return routingConfig
}

// CreateAlias creates an alias for a Lambda function.
// An alias points to a specific version and can be used for traffic shifting.
func (s *LambdaService) CreateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	in := &AliasCreateInput{
		Name:            request.GetStringParam(req.Parameters, "Name"),
		FunctionVersion: request.GetStringParam(req.Parameters, "FunctionVersion"),
		Description:     request.GetStringParam(req.Parameters, "Description"),
	}
	if routingMap := request.GetMapParam(req.Parameters, "RoutingConfig"); routingMap != nil {
		in.RoutingConfig = parseRoutingConfig(routingMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := s.createAliasCore(store, function, in)
	if err != nil {
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

	if err := s.deleteAliasCore(store, functionName, aliasName); err != nil {
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

	alias, err := s.getAliasCore(store, functionName, aliasName)
	if err != nil {
		return nil, err
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
	if err := validateAliasName(aliasName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &AliasUpdateInput{
		Description:     request.GetStringParam(req.Parameters, "Description"),
		FunctionVersion: request.GetStringParam(req.Parameters, "FunctionVersion"),
	}
	_, in.HasDescription = req.Parameters["Description"]
	if routingMap := request.GetMapParam(req.Parameters, "RoutingConfig"); routingMap != nil {
		in.RoutingConfig = parseRoutingConfig(routingMap)
	}

	alias, err := s.updateAliasCore(store, functionName, aliasName, in)
	if err != nil {
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

	allAliases, err := s.listAliasesCore(store, functionName)
	if err != nil {
		return nil, err
	}

	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	marker := request.GetStringParam(req.Parameters, "Marker")

	pageResult := pagination.PaginateSlice(allAliases, marker, maxItems, func(a lambdastore.Alias) string {
		return a.Name
	})

	aliases := make([]interface{}, 0, len(pageResult.Items))
	for _, a := range pageResult.Items {
		aliases = append(aliases, s.toAliasResponse(&a))
	}

	resp := map[string]interface{}{
		"Aliases": aliases,
	}
	if pageResult.IsTruncated {
		resp["NextMarker"] = pageResult.NextMarker
	}

	return resp, nil
}

func (s *LambdaService) toAliasResponse(a *lambdastore.Alias) map[string]interface{} {
	resp := map[string]interface{}{
		"AliasArn":        a.AliasArn,
		"Name":            a.Name,
		"FunctionVersion": a.FunctionVersion,
		"Description":     a.Description,
		"RevisionId":      a.RevisionId,
	}
	if a.RoutingConfig != nil && len(a.RoutingConfig.AdditionalVersionWeights) > 0 {
		resp["RoutingConfig"] = map[string]interface{}{
			"AdditionalVersionWeights": a.RoutingConfig.AdditionalVersionWeights,
		}
	}
	return resp
}
