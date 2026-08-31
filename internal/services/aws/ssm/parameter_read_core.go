package ssm

import (
	"context"
	"errors"
	"strconv"
	"strings"

	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result DTOs — parameter reads
// ---------------------------------------------------------------------------

// GetParameterInput carries the fields for GetParameter.
type GetParameterInput struct {
	Name           string
	WithDecryption bool
}

// GetParameterResult carries the fetched parameter plus the selector string
// the response Selector field must echo.
type GetParameterResult struct {
	Parameter    *ssmstore.Parameter
	UsedSelector string
}

// GetParametersInput carries the fields for GetParameters.
type GetParametersInput struct {
	Names          []string
	WithDecryption bool
}

// GetParametersEntry pairs a fetched parameter with the raw requested name
// so the response serialisation can echo the original selector form.
type GetParametersEntry struct {
	Parameter *ssmstore.Parameter
	RawName   string
}

// GetParametersResult carries the per-name outcomes of GetParameters.
type GetParametersResult struct {
	Parameters        []GetParametersEntry
	InvalidParameters []string
}

// GetParametersByPathInput carries the fields for GetParametersByPath. The
// raw wire parameters travel with the request because ParameterFilters
// arrive in both a JSON-body shape and a flat query-string shape, and the
// parsing that distinguishes them belongs to the Core.
type GetParametersByPathInput struct {
	Path           string
	Recursive      bool
	WithDecryption bool
	MaxResults     int32
	Parameters     map[string]interface{}
	NextToken      string
}

// GetParametersByPathResult carries the parameters under the path and the
// continuation marker.
type GetParametersByPathResult struct {
	Parameters []*ssmstore.Parameter
	NextToken  string
}

// GetParameterHistoryInput carries the fields for GetParameterHistory.
type GetParameterHistoryInput struct {
	Name       string
	MaxResults int32
	NextToken  string
}

// GetParameterHistoryResult carries the version history and the effective
// continuation marker (empty when the history is not truncated).
type GetParameterHistoryResult struct {
	History   []*ssmstore.ParameterVersion
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getParameterCore is the single entry point for single-parameter reads,
// including version and label selectors and SecureString decryption.
func (s *SSMService) getParameterCore(ctx context.Context, store ssmstore.SSMStoreInterface, in GetParameterInput) (*GetParameterResult, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}

	baseName, selector := parseParameterSelector(in.Name)

	var param *ssmstore.Parameter
	var err error
	var usedSelector string

	// When a selector is provided (version or label), the Selector response
	// field must include the full "name:selector" string, matching AWS behaviour.
	if selector == "" {
		param, err = store.GetParameter(baseName, false)
		usedSelector = ""
	} else if version, parseErr := strconv.ParseInt(selector, 10, 64); parseErr == nil {
		param, err = store.GetParameterByVersion(baseName, version)
		usedSelector = in.Name
	} else {
		param, err = store.GetParameterByLabel(baseName, selector)
		usedSelector = in.Name
	}

	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, ErrParameterNotFound
	}

	if in.WithDecryption && param.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
		decryptedValue, err := s.decryptValue(ctx, param.Value, param.KeyID)
		if err != nil {
			return nil, err
		}
		param.Value = decryptedValue
	}

	return &GetParameterResult{Parameter: param, UsedSelector: usedSelector}, nil
}

// getParametersCore is the single entry point for multi-name reads. Names
// is a required member with a modelled length floor, so an empty list is
// rejected by validateParameterNameList before any store access. Names
// that fail to resolve (or fail decryption) are reported as invalid rather
// than failing the whole call, matching the AWS partial-success contract.
func (s *SSMService) getParametersCore(ctx context.Context, store ssmstore.SSMStoreInterface, in GetParametersInput) (*GetParametersResult, error) {
	if err := validateParameterNameList(in.Names); err != nil {
		return nil, err
	}

	result := &GetParametersResult{}
	for _, rawName := range in.Names {
		baseName, selector := parseParameterSelector(rawName)
		var (
			param *ssmstore.Parameter
			err   error
		)
		switch {
		case selector == "":
			param, err = store.GetParameter(baseName, false)
		case isNumericSelector(selector):
			version, _ := strconv.ParseInt(selector, 10, 64)
			param, err = store.GetParameterByVersion(baseName, version)
		default:
			param, err = store.GetParameterByLabel(baseName, selector)
		}
		if err != nil {
			result.InvalidParameters = append(result.InvalidParameters, rawName)
			continue
		}
		if in.WithDecryption && param.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
			decryptedValue, decErr := s.decryptValue(ctx, param.Value, param.KeyID)
			if decErr != nil {
				result.InvalidParameters = append(result.InvalidParameters, rawName)
				continue
			}
			param.Value = decryptedValue
		}
		result.Parameters = append(result.Parameters, GetParametersEntry{Parameter: param, RawName: rawName})
	}

	return result, nil
}

// getParametersByPathCore is the single entry point for hierarchy reads:
// path validation, MaxResults range enforcement, filter parsing, the store
// query, and SecureString decryption all live here. A parameter whose
// decryption fails is dropped from the visible set, not reported.
func (s *SSMService) getParametersByPathCore(ctx context.Context, store ssmstore.SSMStoreInterface, in GetParametersByPathInput) (*GetParametersByPathResult, error) {
	if err := validateHierarchyPath(in.Path); err != nil {
		return nil, err
	}
	path := in.Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	maxResults, err := validateMaxResultsForPath(in.MaxResults)
	if err != nil {
		return nil, err
	}
	filters, err := parseParameterFilters(in.Parameters)
	if err != nil {
		return nil, err
	}

	parameters, nextMarker, err := store.GetParametersByPath(path, in.Recursive, false, filters, maxResults, in.NextToken)
	if err != nil {
		return nil, err
	}

	visible := make([]*ssmstore.Parameter, 0, len(parameters))
	for _, p := range parameters {
		if in.WithDecryption && p.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
			decryptedValue, err := s.decryptValue(ctx, p.Value, p.KeyID)
			if err != nil {
				continue
			}
			p.Value = decryptedValue
		}
		visible = append(visible, p)
	}

	return &GetParametersByPathResult{Parameters: visible, NextToken: nextMarker}, nil
}

// getParameterHistoryCore is the single entry point for version-history
// reads: name validation, MaxResults range enforcement, the store query,
// and continuation-token semantics all live here. The continuation marker
// is only surfaced when the store reports truncation.
func (s *SSMService) getParameterHistoryCore(store ssmstore.SSMStoreInterface, in GetParameterHistoryInput) (*GetParameterHistoryResult, error) {
	if in.Name == "" {
		return nil, ErrInvalidParameterName
	}

	maxResults, err := validateMaxResultsForPage(in.MaxResults)
	if err != nil {
		return nil, err
	}

	history, nextMarker, isTruncated, err := store.GetParameterHistory(in.Name, maxResults, in.NextToken)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		if errors.Is(err, ssmstore.ErrInvalidNextToken) {
			return nil, ErrInvalidNextToken
		}
		return nil, err
	}

	respNextToken := ""
	if isTruncated && nextMarker != "" {
		respNextToken = nextMarker
	}

	return &GetParameterHistoryResult{History: history, NextToken: respNextToken}, nil
}

// ---------------------------------------------------------------------------
// Wire-level parsers — shared by the read Cores
// ---------------------------------------------------------------------------

// parseParameterSelector splits an optional ":selector" suffix off a
// parameter name. The selector is a version number or a label.
func parseParameterSelector(name string) (baseName string, selector string) {
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		return name[:idx], name[idx+1:]
	}
	return name, ""
}

// isNumericSelector reports whether a parameter selector is a version number
// (e.g. ":5") rather than a label (e.g. ":production").
func isNumericSelector(selector string) bool {
	_, err := strconv.ParseInt(selector, 10, 64)
	return err == nil
}

// parseParameterFilters extracts ParameterFilters from a request in either
// shape: a JSON-decoded "ParameterFilters" key holding a list of maps (the
// AWS SDK JSON protocol), or flat query-string keys with the
// "ParameterFilters.N.Key" / "ParameterFilters.N.Option" /
// "ParameterFilters.N.Values.M" pattern (the AWS SDK query protocol).
//
// The deprecated "Filters" field uses the same query-string shape but with
// no Option member; this parser accepts both field names so SDK clients that
// still use the deprecated ParametersFilter keep working.
func parseParameterFilters(params map[string]interface{}) ([]ssmstore.ParameterFilter, error) {
	if filters, err := filtersFromList(params["ParameterFilters"]); filters != nil || err != nil {
		return filters, err
	}
	if filters, err := filtersFromList(params["Filters"]); filters != nil || err != nil {
		return filters, err
	}
	return filtersFromQueryParams(params, "ParameterFilters")
}

func filtersFromQueryParams(params map[string]interface{}, field string) ([]ssmstore.ParameterFilter, error) {
	var filters []ssmstore.ParameterFilter
	for i := 1; ; i++ {
		// AWS Query protocol uses ".member." between the field and the index.
		// SSM itself is awsJson1_1 and reaches this branch only when the
		// caller sends flat query-string keys; the .member. prefix restores
		// the canonical query protocol shape.
		key := paramString(params, field+".member."+strconv.Itoa(i)+".Key")
		if key == "" {
			return filters, nil
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, ErrInvalidFilterKey
		}
		option := paramString(params, field+".member."+strconv.Itoa(i)+".Option")
		if option != "" && !ssmstore.ValidateParameterFilterOption(option) {
			return nil, ErrInvalidFilterOption
		}
		var values []string
		for j := 1; ; j++ {
			v := paramString(params, field+".member."+strconv.Itoa(i)+".Values.member."+strconv.Itoa(j))
			if v == "" {
				break
			}
			values = append(values, v)
		}
		// The Values member is optional on the wire shape ("Required: No"):
		// an entry with a Key and no Values is a key-existence filter, not a
		// malformed one. The flat query encoding cannot distinguish an empty
		// Values member from an omitted member, so both arrive here as an
		// empty collection and are treated as the optional member omitted.
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: option,
			Values: values,
		})
	}
}

// filtersFromList extracts filters from a JSON-decoded list under the given
// key. Returns (nil, nil) when the key is absent so callers can fall through
// to other formats. Malformed input is rejected, never silently reduced to
// its valid subset. Wire-type violations (a non-list filter list, a
// non-object entry, a non-string Key, Option or Values entry) cannot be
// deserialised into the modelled awsJson1_1 shape and return
// SerializationException before operation-level validation runs. Well-typed
// but constraint-violating input returns the modelled filter errors: an empty
// Key or Values entry violates the ParameterStringFilterKey /
// ParameterStringFilterValue length minimum of 1, and an empty or unknown
// Option is not a valid option. An omitted Values member is accepted as a
// key-existence filter: the member is optional on the wire shape ("Required:
// No"), while a Values member that is present but empty violates the
// ParameterStringFilterValueList length minimum of 1 and is rejected.
func filtersFromList(raw interface{}) ([]ssmstore.ParameterFilter, error) {
	if raw == nil {
		return nil, nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, ErrSerializationException
	}
	var filters []ssmstore.ParameterFilter
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, ErrSerializationException
		}
		key, ok := m["Key"].(string)
		if !ok {
			return nil, ErrSerializationException
		}
		if key == "" {
			return nil, ErrInvalidFilterValue
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, ErrInvalidFilterKey
		}
		var option string
		if rawOption, hasOption := m["Option"]; hasOption && rawOption != nil {
			s, ok := rawOption.(string)
			if !ok {
				return nil, ErrSerializationException
			}
			if s == "" {
				return nil, ErrInvalidFilterOption
			}
			option = s
		}
		if option != "" && !ssmstore.ValidateParameterFilterOption(option) {
			return nil, ErrInvalidFilterOption
		}
		var values []string
		if rawValues, hasValues := m["Values"]; hasValues {
			vs, ok := rawValues.([]interface{})
			if !ok {
				return nil, ErrSerializationException
			}
			for _, v := range vs {
				s, ok := v.(string)
				if !ok {
					return nil, ErrSerializationException
				}
				if s == "" {
					return nil, ErrInvalidFilterValue
				}
				values = append(values, s)
			}
			// With every entry validated non-empty, only a present-but-empty
			// Values list reaches this check.
			if len(values) == 0 {
				return nil, ErrInvalidFilterValue
			}
		}
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: option,
			Values: values,
		})
	}
	return filters, nil
}

// paramString returns the params map value as a string, handling both
// single-value and slice encodings that come out of url.ParseQuery.
func paramString(params map[string]interface{}, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		if len(val) == 0 {
			return ""
		}
		if s, ok := val[0].(string); ok {
			return s
		}
	}
	return ""
}
