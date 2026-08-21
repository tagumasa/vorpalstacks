// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// GetFunction retrieves information about the specified Lambda function.
func (s *LambdaService) GetFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	function, version, _, tags, err := s.getFunctionCore(store, &GetFunctionInput{
		FunctionName: request.GetStringParam(req.Parameters, "FunctionName"),
		Qualifier:    request.GetStringParam(req.Parameters, "Qualifier"),
	})
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if version != nil {
		config = s.toVersionConfiguration(version)
	} else {
		config = s.toFunctionConfiguration(function)
	}

	resp := map[string]interface{}{
		"Configuration": config,
		"Code": map[string]interface{}{
			"Location":       function.CodeLocation,
			"RepositoryType": repositoryType(function),
			"ImageUri":       function.ImageUri,
		},
		"Tags": tags,
	}
	// Reserved concurrency is function-level and appears in GetFunction
	// responses whenever a limit is configured.
	if function.ReservedConcurrency != nil {
		resp["Concurrency"] = map[string]interface{}{
			"ReservedConcurrentExecutions": *function.ReservedConcurrency,
		}
	}
	return resp, nil
}

// GetFunctionConfiguration retrieves the configuration of the specified Lambda function.
func (s *LambdaService) GetFunctionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	function, version, _, err := s.getFunctionConfigurationCore(store, &GetFunctionInput{
		FunctionName: request.GetStringParam(req.Parameters, "FunctionName"),
		Qualifier:    request.GetStringParam(req.Parameters, "Qualifier"),
	})
	if err != nil {
		return nil, err
	}

	if version != nil {
		return s.toVersionConfiguration(version), nil
	}

	return s.toFunctionConfiguration(function), nil
}

// ListFunctions lists all Lambda functions in the current account.
func (s *LambdaService) ListFunctions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))

	// "Set to ALL to include entries for all published versions of each
	// function" — entries, not functions: "ListFunctions returns a maximum
	// of 50 items in each response, even if you set the number higher", so
	// the version expansion must paginate per entry, not per function.
	if request.GetStringParam(req.Parameters, "FunctionVersion") == "ALL" {
		entries, nextMarker, err := s.listFunctionVersionEntries(store, marker, maxItems)
		if err != nil {
			return nil, err
		}
		functions := make([]interface{}, 0, len(entries))
		for _, entry := range entries {
			if entry.versionIndex < 0 {
				functions = append(functions, s.toFunctionConfiguration(entry.fn))
			} else {
				functions = append(functions, s.toVersionConfiguration(&entry.fn.Versions[entry.versionIndex]))
			}
		}
		response := map[string]interface{}{"Functions": functions}
		if nextMarker != "" {
			response["NextMarker"] = nextMarker
		}
		return response, nil
	}

	items, nextMarker, err := s.listFunctionsCore(store, &ListFunctionsInput{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	functions := make([]interface{}, 0, len(items))
	for _, fn := range items {
		functions = append(functions, s.toFunctionConfiguration(fn))
	}

	response := map[string]interface{}{
		"Functions": functions,
	}

	if nextMarker != "" {
		response["NextMarker"] = nextMarker
	}

	return response, nil
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
