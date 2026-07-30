package ssm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

func getIntParam(req *request.ParsedRequest, key string) int32 {
	return int32(request.GetIntParam(req.Parameters, key))
}

func getBoolParam(req *request.ParsedRequest, key string) bool {
	return request.GetBoolParam(req.Parameters, key)
}

func parseParameterSelector(name string) (baseName string, selector string) {
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		return name[:idx], name[idx+1:]
	}
	return name, ""
}

func parseStringList(params map[string]interface{}, field, memberPrefix string) []string {
	if raw, ok := params[field].([]interface{}); ok {
		var result []string
		for _, v := range raw {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	var result []string
	for i := 1; ; i++ {
		val := request.GetStringParam(params, memberPrefix+strconv.Itoa(i))
		if val == "" {
			break
		}
		result = append(result, val)
	}
	return result
}

func parameterToResponse(p *ssmstore.Parameter, selector string) map[string]interface{} {
	if selector == "" {
		selector = fmt.Sprintf("%s:%d", p.Name, p.Version)
	}
	return map[string]interface{}{
		"Name":             p.Name,
		"Type":             string(p.Type),
		"Value":            p.Value,
		"Version":          p.Version,
		"Selector":         selector,
		"SourceResult":     "",
		"LastModifiedDate": p.LastModifiedDate.Unix(),
		"LastModifiedUser": lastModifiedUser(p.LastModifiedBy),
		"ARN":              p.ARN,
		"DataType":         p.DataType,
	}
}

// lastModifiedUser maps the stored access-key identity to a principal-shaped
// string. When a key is recorded we expose it directly; otherwise we fall
// back to the platform identifier so the response field is always populated
// (AWS returns a non-empty value for this field).
func lastModifiedUser(by string) string {
	if by != "" {
		return by
	}
	return "vorpalstacks:admin"
}

// PutParameter adds or updates a parameter in the Parameter Store.
func (s *SSMService) PutParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")

	param, err := normalisePutParameter(ParameterPutFields{
		Name:           name,
		Value:          req.GetParam("Value"),
		Type:           req.GetParam("Type"),
		Description:    req.GetParam("Description"),
		KeyID:          req.GetParam("KeyId"),
		AllowedPattern: req.GetParam("AllowedPattern"),
		DataType:       req.GetParam("DataType"),
		Tier:           req.GetParam("Tier"),
		Policies:       req.GetParam("Policies"),
	})
	if err != nil {
		return nil, err
	}

	overwrite := getBoolParam(req, "Overwrite")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	version, err := s.putParameterWithEncryption(ctx, store, param, overwrite, reqCtx.Principal)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterAlreadyExists) {
			return nil, ErrParameterAlreadyExists
		}
		if errors.Is(err, ssmstore.ErrReservedParameterName) {
			return nil, ErrParameterPatternMismatch
		}
		if errors.Is(err, ssmstore.ErrInvalidAllowedPattern) {
			return nil, ErrInvalidAllowedPattern
		}
		if errors.Is(err, ssmstore.ErrParameterPatternMismatch) {
			return nil, ErrParameterPatternMismatch
		}
		return nil, err
	}

	if tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig); len(tags) > 0 {
		if err := store.AddTagsToResource(name, tagutil.ToMap(tags)); err != nil {
			logs.Warn("Failed to add tags to parameter", logs.String("name", name), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Version": version,
		"Tier":    string(param.Tier),
	}, nil
}

// GetParameter retrieves a parameter from the Parameter Store.
func (s *SSMService) GetParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	withDecryption := getBoolParam(req, "WithDecryption")

	baseName, selector := parseParameterSelector(name)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var param *ssmstore.Parameter
	var usedSelector string

	// When a selector is provided (version or label), the Selector response
	// field must include the full "name:selector" string, matching AWS behaviour.
	if selector == "" {
		param, err = store.GetParameter(baseName, false)
		usedSelector = ""
	} else if version, parseErr := strconv.ParseInt(selector, 10, 64); parseErr == nil {
		param, err = store.GetParameterByVersion(baseName, version)
		usedSelector = name
	} else {
		param, err = store.GetParameterByLabel(baseName, selector)
		usedSelector = name
	}

	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, ErrParameterNotFound
	}

	if withDecryption && param.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
		decryptedValue, err := s.decryptValue(ctx, param.Value, param.KeyID)
		if err != nil {
			return nil, err
		}
		param.Value = decryptedValue
	}

	return map[string]interface{}{
		"Parameter": parameterToResponse(param, usedSelector),
	}, nil
}

// GetParameters retrieves multiple parameters from the Parameter Store.
func (s *SSMService) GetParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	withDecryption := getBoolParam(req, "WithDecryption")
	names := parseStringList(req.Parameters, "Names", "Names.")
	if err := validateParameterNameList(names); err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return map[string]interface{}{
			"Parameters":        []interface{}{},
			"InvalidParameters": []string{},
		}, nil
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(names))
	var invalidNames []string
	for _, rawName := range names {
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
			invalidNames = append(invalidNames, rawName)
			continue
		}
		if withDecryption && param.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
			decryptedValue, decErr := s.decryptValue(ctx, param.Value, param.KeyID)
			if decErr != nil {
				invalidNames = append(invalidNames, rawName)
				continue
			}
			param.Value = decryptedValue
		}
		params = append(params, parameterToResponse(param, rawName))
	}

	return map[string]interface{}{
		"Parameters":        params,
		"InvalidParameters": invalidNames,
	}, nil
}

// isNumericSelector reports whether a parameter selector is a version number
// (e.g. ":5") rather than a label (e.g. ":production").
func isNumericSelector(selector string) bool {
	_, err := strconv.ParseInt(selector, 10, 64)
	return err == nil
}

// GetParametersByPath retrieves parameters under a specified path.
func (s *SSMService) GetParametersByPath(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	path := req.GetParam("Path")
	if err := validateHierarchyPath(path); err != nil {
		return nil, err
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	recursive := getBoolParam(req, "Recursive")
	withDecryption := getBoolParam(req, "WithDecryption")
	maxResults, err := validateMaxResultsForPath(getIntParam(req, "MaxResults"))
	if err != nil {
		return nil, err
	}
	filters, err := parseParameterFilters(req.Parameters)
	if err != nil {
		return nil, err
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	parameters, nextMarker, err := store.GetParametersByPath(path, recursive, false, filters, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(parameters))
	for _, p := range parameters {
		if withDecryption && p.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
			decryptedValue, err := s.decryptValue(ctx, p.Value, p.KeyID)
			if err != nil {
				continue
			}
			p.Value = decryptedValue
		}
		params = append(params, parameterToResponse(p, ""))
	}

	response := map[string]interface{}{
		"Parameters": params,
	}
	pagination.SetNextToken(response, "NextToken", nextMarker)

	return response, nil
}

// DeleteParameter removes a parameter from the Parameter Store.
func (s *SSMService) DeleteParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteParameter(name); err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		logs.Error("Failed to delete parameter from store",
			logs.String("name", name), logs.Err(err))
		return nil, awserrors.NewInternalErrorException("failed to delete parameter")
	}

	return response.EmptyResponse(), nil
}

// DeleteParameters removes multiple parameters from the Parameter Store.
func (s *SSMService) DeleteParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	names := parseStringList(req.Parameters, "Names", "Names.")
	if err := validateParameterNameList(names); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	deleted, invalid := store.DeleteParameters(names)

	return map[string]interface{}{
		"DeletedParameters": deleted,
		"InvalidParameters": invalid,
	}, nil
}

// DescribeParameters returns information about all parameters in the Parameter Store.
func (s *SSMService) DescribeParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, err := validateMaxResultsForPage(getIntParam(req, "MaxResults"))
	if err != nil {
		return nil, err
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	filters, err := parseParameterFilters(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	metadata, nextMarker, err := store.DescribeParameters(filters, maxResults, nextToken)
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(metadata))
	for _, m := range metadata {
		params = append(params, map[string]interface{}{
			"Name":             m.Name,
			"Type":             string(m.Type),
			"KeyId":            m.KeyID,
			"LastModifiedDate": m.LastModifiedDate.Unix(),
			"LastModifiedUser": lastModifiedUser(m.LastModifiedBy),
			"Description":      m.Description,
			"AllowedPattern":   m.AllowedPattern,
			"Version":          m.Version,
			"Tier":             string(m.Tier),
			"Policies":         policiesToResponse(m.Policies),
			"DataType":         m.DataType,
		})
	}

	response := map[string]interface{}{
		"Parameters": params,
	}
	pagination.SetNextToken(response, "NextToken", nextMarker)

	return response, nil
}

// GetParameterHistory retrieves the history of a parameter's versions.
func (s *SSMService) GetParameterHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	maxResults, err := validateMaxResultsForPage(getIntParam(req, "MaxResults"))
	if err != nil {
		return nil, err
	}
	withDecryption := getBoolParam(req, "WithDecryption")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	history, nextMarker, isTruncated, err := store.GetParameterHistory(name, maxResults, nextToken)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		return nil, err
	}

	versions := make([]map[string]interface{}, 0, len(history))
	for _, v := range history {
		value := v.Value
		if !withDecryption && v.Type == ssmstore.ParameterTypeSecureString {
			value = ""
		}
		labels := v.Labels
		if labels == nil {
			labels = []string{}
		}
		versions = append(versions, map[string]interface{}{
			"Name":             v.ParameterName,
			"Type":             string(v.Type),
			"KeyId":            v.KeyID,
			"LastModifiedDate": v.LastModifiedDate.Unix(),
			"LastModifiedUser": lastModifiedUser(v.LastModifiedBy),
			"Description":      v.Description,
			"AllowedPattern":   v.AllowedPattern,
			"Version":          v.Version,
			"Labels":           labels,
			"Tier":             string(v.Tier),
			"Value":            value,
			"DataType":         v.DataType,
		})
	}

	respNextToken := ""
	if isTruncated && nextMarker != "" {
		respNextToken = nextMarker
	}

	resp := map[string]interface{}{
		"Parameters": versions,
	}
	pagination.SetNextToken(resp, "NextToken", respNextToken)
	return resp, nil
}

// LabelParameterVersion attaches labels to a specific version of a parameter.
func (s *SSMService) LabelParameterVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	labels := parseStringList(req.Parameters, "Labels", "Labels.member.")
	if err := validateLabels(labels); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	parameterVersion := int64(getIntParam(req, "ParameterVersion"))
	if parameterVersion == 0 {
		// AWS spec: if ParameterVersion is omitted, default to the latest version.
		param, err := store.GetParameter(name, false)
		if err != nil {
			return nil, ErrParameterNotFound
		}
		parameterVersion = param.Version
	}

	invalidLabels, err := store.LabelParameterVersion(name, parameterVersion, labels)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"InvalidLabels":    invalidLabels,
		"ParameterVersion": parameterVersion,
	}, nil
}

// UnlabelParameterVersion removes labels from a specific version of a parameter.
func (s *SSMService) UnlabelParameterVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	labels := parseStringList(req.Parameters, "Labels", "Labels.member.")
	if err := validateLabels(labels); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	parameterVersion := int64(getIntParam(req, "ParameterVersion"))
	if parameterVersion == 0 {
		// AWS spec: if ParameterVersion is omitted, default to the latest version.
		param, err := store.GetParameter(name, false)
		if err != nil {
			return nil, ErrParameterNotFound
		}
		parameterVersion = param.Version
	}

	removedLabels, err := store.UnlabelParameterVersion(name, parameterVersion, labels)
	if err != nil {
		if errors.Is(err, ssmstore.ErrParameterNotFound) {
			return nil, ErrParameterNotFound
		}
		if errors.Is(err, ssmstore.ErrParameterVersionNotFound) {
			return nil, ErrParameterVersionNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"InvalidLabels": []string{},
		"RemovedLabels": removedLabels,
	}, nil
}
