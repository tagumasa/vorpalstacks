package ssm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	pagination "vorpalstacks/internal/common/pagination"
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
		"ARN":              p.ARN,
		"DataType":         p.DataType,
	}
}

// PutParameter adds or updates a parameter in the Parameter Store.
func (s *SSMService) PutParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")
	if name == "" {
		return nil, ErrInvalidParameterName
	}

	value := req.GetParam("Value")
	paramType := req.GetParam("Type")
	if paramType == "" {
		paramType = "String"
	}

	keyID := req.GetParam("KeyId")

	if paramType == "SecureString" && s.kmsEncryptor != nil {
		encryptedValue, err := s.encryptValue(ctx, value, keyID)
		if err != nil {
			return nil, err
		}
		value = encryptedValue
	}

	param := ssmstore.NewParameter(name, value, ssmstore.ParameterType(paramType))
	param.Description = req.GetParam("Description")
	param.KeyID = keyID
	param.AllowedPattern = req.GetParam("AllowedPattern")
	param.DataType = req.GetParam("DataType")
	if param.DataType == "" {
		param.DataType = "text"
	}

	tier := req.GetParam("Tier")
	if tier != "" {
		param.Tier = ssmstore.ParameterTier(tier)
	}

	overwrite := getBoolParam(req, "Overwrite")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	version, err := store.PutParameter(param, overwrite)
	if err != nil {
		return nil, err
	}

	if !overwrite {
		if tags := tagutil.GetTags(req.Parameters, tagutil.StandardConfig); len(tags) > 0 {
			if err := store.AddTagsToResource(name, tagutil.ToMap(tags)); err != nil {
				logs.Warn("Failed to add tags to parameter", logs.String("name", name), logs.Err(err))
			}
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

	var names []string
	if namesRaw, ok := req.Parameters["Names"].([]interface{}); ok {
		for _, n := range namesRaw {
			if ns, ok := n.(string); ok {
				names = append(names, ns)
			}
		}
	} else {
		for i := 1; ; i++ {
			name := req.GetParam("Names." + strconv.Itoa(i))
			if name == "" {
				break
			}
			names = append(names, name)
		}
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
	parameters, invalidNames := store.GetParameters(names, false)

	params := make([]map[string]interface{}, 0, len(parameters))
	for _, p := range parameters {
		if withDecryption && p.Type == ssmstore.ParameterTypeSecureString && s.kmsEncryptor != nil {
			decryptedValue, err := s.decryptValue(ctx, p.Value, p.KeyID)
			if err != nil {
				invalidNames = append(invalidNames, p.Name)
				continue
			}
			p.Value = decryptedValue
		}
		params = append(params, parameterToResponse(p, ""))
	}

	return map[string]interface{}{
		"Parameters":        params,
		"InvalidParameters": invalidNames,
	}, nil
}

// GetParametersByPath retrieves parameters under a specified path.
func (s *SSMService) GetParametersByPath(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	path := req.GetParam("Path")
	if path == "" {
		return nil, ErrInvalidParameterName
	}

	recursive := getBoolParam(req, "Recursive")
	withDecryption := getBoolParam(req, "WithDecryption")
	maxResults := getIntParam(req, "MaxResults")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	parameters, nextMarker, err := store.GetParametersByPath(path, recursive, false, maxResults, nextToken)
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
		return nil, fmt.Errorf("failed to delete parameter: %w", err)
	}

	return response.EmptyResponse(), nil
}

// DeleteParameters removes multiple parameters from the Parameter Store.
func (s *SSMService) DeleteParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var names []string
	if namesRaw, ok := req.Parameters["Names"].([]interface{}); ok {
		for _, n := range namesRaw {
			if ns, ok := n.(string); ok {
				names = append(names, ns)
			}
		}
	} else {
		for i := 1; ; i++ {
			name := req.GetParam("Names." + strconv.Itoa(i))
			if name == "" {
				break
			}
			names = append(names, name)
		}
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
	maxResults := getIntParam(req, "MaxResults")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	filters := make(map[string]string)

	for _, filterKey := range []string{"ParameterFilters", "Filters"} {
		if filterArr, ok := req.Parameters[filterKey]; ok {
			if filterArr, ok := filterArr.([]interface{}); ok {
				for _, f := range filterArr {
					if fm, ok := f.(map[string]interface{}); ok {
						key, _ := fm["Key"].(string)
						values, _ := fm["Values"].([]interface{})
						if key != "" && len(values) > 0 {
							if v, ok := values[0].(string); ok {
								filters[key] = v
							}
						}
					}
				}
			}
		}
	}

	if len(filters) == 0 {
		for i := 1; ; i++ {
			filterKey := req.GetParam("ParameterFilters.member." + strconv.Itoa(i) + ".Key")
			if filterKey == "" {
				break
			}
			filterValue := req.GetParam("ParameterFilters.member." + strconv.Itoa(i) + ".Values.member.1")
			if filterValue == "" {
				filterValue = req.GetParam("Filters." + strconv.Itoa(i) + ".Key")
				if filterValue == "" {
					break
				}
				filterKey = filterValue
				filterValue = req.GetParam("Filters." + strconv.Itoa(i) + ".Value")
			}
			filters[filterKey] = filterValue
		}
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
			"LastModifiedUser": "N/A",
			"Description":      m.Description,
			"AllowedPattern":   m.AllowedPattern,
			"Version":          m.Version,
			"Tier":             string(m.Tier),
			"Policies":         []interface{}{},
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

	maxResults := getIntParam(req, "MaxResults")
	withDecryption := getBoolParam(req, "WithDecryption")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	history, nextMarker, isTruncated, err := store.GetParameterHistory(name, maxResults, nextToken)
	if err != nil {
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
			"LastModifiedUser": "N/A",
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

	parameterVersion := int64(getIntParam(req, "ParameterVersion"))
	if parameterVersion == 0 {
		return nil, ErrInvalidParameterVersion
	}

	var labels []string
	if labelsRaw, ok := req.Parameters["Labels"].([]interface{}); ok {
		for _, l := range labelsRaw {
			if ls, ok := l.(string); ok {
				labels = append(labels, ls)
			}
		}
	} else {
		for i := 1; ; i++ {
			label := req.GetParam("Labels.member." + strconv.Itoa(i))
			if label == "" {
				break
			}
			labels = append(labels, label)
		}
	}

	if len(labels) == 0 {
		return nil, ErrInvalidParameterLabel
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.LabelParameterVersion(name, parameterVersion, labels); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"InvalidLabels":    []string{},
		"ParameterVersion": parameterVersion,
	}, nil
}
