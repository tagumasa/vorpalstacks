package lambda

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// GetFunctionInput carries the fields needed for GetFunction.
type GetFunctionInput struct {
	FunctionName string
	Qualifier    string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// validateAndGetFunction resolves the FunctionName wire reference and
// fetches the function, failing with ResourceNotFound when it does not
// exist.
func (s *LambdaService) validateAndGetFunction(ctx *request.RequestContext, params map[string]interface{}) (*lambdastore.Function, error) {
	functionName := request.GetStringParam(params, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(ctx)
	if err != nil {
		return nil, err
	}
	function, err := store.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return function, nil
}

// validateAndGetFunctionWithQualifier resolves the FunctionName wire
// reference together with its qualifier and resolves the addressed
// version or alias.
func (s *LambdaService) validateAndGetFunctionWithQualifier(ctx *request.RequestContext, params map[string]interface{}) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	functionNameRaw := request.GetStringParam(params, "FunctionName")
	if functionNameRaw == "" {
		return nil, nil, nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName, embeddedQualifier := resolveFunctionRef(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return nil, nil, nil, err
	}

	qualifier := mergeQualifier(request.GetStringParam(params, "Qualifier"), embeddedQualifier)
	store, err := s.store(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return s.resolveQualifier(store.Functions, functionName, qualifier)
}

// resolveQualifier resolves a function qualifier to its function, version
// and alias, mapping the store's not-found sentinels onto the API's
// ResourceNotFound shapes.
func (s *LambdaService) resolveQualifier(store *lambdastore.FunctionStore, functionName, qualifier string) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	function, version, alias, err := store.ResolveQualifier(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, nil, nil, ErrResourceNotFound
		}
		if err == lambdastore.ErrVersionNotFound {
			return nil, nil, nil, NewLambdaError("ResourceNotFoundException",
				fmt.Sprintf("Qualifier '%s' not found for function '%s'.", qualifier, functionName),
				http.StatusNotFound)
		}
		return nil, nil, nil, err
	}
	return function, version, alias, nil
}

// getFunctionCore retrieves a function (optionally by qualifier) along with
// its tags. It is the single entry point shared by the HTTP API handler and
// the admin gRPC handler. The raw function name or ARN is resolved and
// validated internally so that all callers share a single validation path.
func (s *LambdaService) getFunctionCore(stores *lambdaStore, in *GetFunctionInput) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, map[string]string, error) {
	functionName, embeddedQualifier := resolveFunctionRef(in.FunctionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, nil, nil, nil, err
	}

	function, version, alias, err := s.resolveQualifier(stores.Functions, functionName, mergeQualifier(in.Qualifier, embeddedQualifier))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// An alias qualifier addresses the published version it points to;
	// Get operations always report the alias's primary version (weighted
	// routing affects invocation only).
	if alias != nil {
		version = findVersion(function, alias.FunctionVersion)
	}

	tags, err := stores.Functions.TagStore.List(function.FunctionName)
	if err != nil {
		logs.Warn("Failed to fetch tags for function",
			logs.String("function", function.FunctionName),
			logs.Err(err))
		tags = map[string]string{}
	}

	return function, version, alias, tags, nil
}

// getFunctionConfigurationCore retrieves a function configuration (optionally
// by qualifier). It is the single entry point shared by the HTTP API handler
// and the admin gRPC handler.
func (s *LambdaService) getFunctionConfigurationCore(stores *lambdaStore, in *GetFunctionInput) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	functionName, embeddedQualifier := resolveFunctionRef(in.FunctionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, nil, nil, err
	}

	function, version, alias, err := s.resolveQualifier(stores.Functions, functionName, mergeQualifier(in.Qualifier, embeddedQualifier))
	if err != nil {
		return nil, nil, nil, err
	}
	// An alias qualifier addresses the published version it points to;
	// Get operations always report the alias's primary version (weighted
	// routing affects invocation only).
	if alias != nil {
		version = findVersion(function, alias.FunctionVersion)
	}
	return function, version, alias, nil
}

// versionListEntry is one FunctionVersion=ALL listing entry: the
// function's own configuration entry (versionIndex -1) or one of its
// published versions.
type versionListEntry struct {
	fn           *lambdastore.Function
	versionIndex int
}

// versionListMarkerSep separates a function name from the number of the
// function's entries already emitted. Function names match
// ^[a-zA-Z0-9-_]+$, so the separator cannot occur inside one.
const versionListMarkerSep = '|'

// listFunctionVersionEntries assembles one FunctionVersion=ALL page with
// entry-level pagination. The marker is opaque to clients; internally it is
// "name|emitted" where emitted counts the entries already listed for that
// function (0 = the function's own entry comes next), or a plain function
// name to resume strictly after it (the store's marker semantics).
func (s *LambdaService) listFunctionVersionEntries(store *lambdaStore, marker string, maxItems int) ([]versionListEntry, string, error) {
	entries := make([]versionListEntry, 0, maxItems)
	listAfter := marker
	emit := func(entry versionListEntry) bool {
		if len(entries) >= maxItems {
			return false
		}
		entries = append(entries, entry)
		return true
	}

	// Resume inside a partially listed function when the marker carries an
	// entry offset.
	if idx := strings.IndexByte(marker, versionListMarkerSep); idx >= 0 {
		name := marker[:idx]
		emitted, err := strconv.Atoi(marker[idx+1:])
		if err != nil || emitted < 0 || name == "" {
			return nil, "", NewInvalidParameter("Marker", fmt.Sprintf("Invalid pagination token: %s", marker))
		}
		listAfter = name
		fn, gerr := store.Functions.Get(name)
		if gerr != nil {
			return nil, "", mapStoreError(gerr)
		}
		if emitted == 0 {
			if !emit(versionListEntry{fn: fn, versionIndex: -1}) {
				return entries, fmt.Sprintf("%s%c%d", name, versionListMarkerSep, 0), nil
			}
			emitted = 1
		}
		for v := emitted - 1; v < len(fn.Versions); v++ {
			if !emit(versionListEntry{fn: fn, versionIndex: v}) {
				return entries, fmt.Sprintf("%s%c%d", name, versionListMarkerSep, v+1), nil
			}
		}
	}

	// Fresh functions from the store: every function contributes at least
	// its own entry, so one store page of maxItems functions always fills
	// the entry page.
	result, err := store.Functions.List(storecommon.ListOptions{
		Marker:   listAfter,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, "", err
	}
	for _, fn := range result.Items {
		if !emit(versionListEntry{fn: fn, versionIndex: -1}) {
			return entries, fmt.Sprintf("%s%c%d", fn.FunctionName, versionListMarkerSep, 0), nil
		}
		for v := range fn.Versions {
			if !emit(versionListEntry{fn: fn, versionIndex: v}) {
				return entries, fmt.Sprintf("%s%c%d", fn.FunctionName, versionListMarkerSep, v+1), nil
			}
		}
	}
	if result.NextMarker == "" {
		return entries, "", nil
	}
	// The store page ended before the entry page filled (every function
	// contributed exactly one entry); resume after its last function.
	return entries, result.NextMarker, nil
}
