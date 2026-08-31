package ssm

import (
	"context"
	"fmt"
	"strconv"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

func getIntParam(req *request.ParsedRequest, key string) int32 {
	return int32(request.GetIntParam(req.Parameters, key))
}

func getBoolParam(req *request.ParsedRequest, key string) bool {
	return request.GetBoolParam(req.Parameters, key)
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putParameterCore(ctx, store, PutParameterInput{
		Fields: ParameterPutFields{
			Name:           req.GetParam("Name"),
			Value:          req.GetParam("Value"),
			Type:           req.GetParam("Type"),
			Description:    req.GetParam("Description"),
			KeyID:          req.GetParam("KeyId"),
			AllowedPattern: req.GetParam("AllowedPattern"),
			DataType:       req.GetParam("DataType"),
			Tier:           req.GetParam("Tier"),
			Policies:       req.GetParam("Policies"),
			Tags:           tagutil.ToMap(tagutil.GetTags(req.Parameters, tagutil.StandardConfig)),
		},
		Overwrite:  getBoolParam(req, "Overwrite"),
		ModifiedBy: reqCtx.Principal,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Version": result.Version,
		"Tier":    string(result.Tier),
	}, nil
}

// GetParameter retrieves a parameter from the Parameter Store.
func (s *SSMService) GetParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getParameterCore(ctx, store, GetParameterInput{
		Name:           req.GetParam("Name"),
		WithDecryption: getBoolParam(req, "WithDecryption"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Parameter": parameterToResponse(result.Parameter, result.UsedSelector),
	}, nil
}

// GetParameters retrieves multiple parameters from the Parameter Store.
func (s *SSMService) GetParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	names := parseStringList(req.Parameters, "Names", "Names.")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getParametersCore(ctx, store, GetParametersInput{
		Names:          names,
		WithDecryption: getBoolParam(req, "WithDecryption"),
	})
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(result.Parameters))
	for _, entry := range result.Parameters {
		params = append(params, parameterToResponse(entry.Parameter, entry.RawName))
	}

	return map[string]interface{}{
		"Parameters":        params,
		"InvalidParameters": result.InvalidParameters,
	}, nil
}

// GetParametersByPath retrieves parameters under a specified path.
func (s *SSMService) GetParametersByPath(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getParametersByPathCore(ctx, store, GetParametersByPathInput{
		Path:           req.GetParam("Path"),
		Recursive:      getBoolParam(req, "Recursive"),
		WithDecryption: getBoolParam(req, "WithDecryption"),
		MaxResults:     getIntParam(req, "MaxResults"),
		Parameters:     req.Parameters,
		NextToken:      pagination.GetMarker(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(result.Parameters))
	for _, p := range result.Parameters {
		params = append(params, parameterToResponse(p, ""))
	}

	response := map[string]interface{}{
		"Parameters": params,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// DeleteParameter removes a parameter from the Parameter Store.
func (s *SSMService) DeleteParameter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := req.GetParam("Name")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteParameterCore(store, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteParameters removes multiple parameters from the Parameter Store.
func (s *SSMService) DeleteParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	names := parseStringList(req.Parameters, "Names", "Names.")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteParametersCore(store, DeleteParametersInput{Names: names})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DeletedParameters": result.DeletedParameters,
		"InvalidParameters": result.InvalidParameters,
	}, nil
}

// DescribeParameters returns information about all parameters in the Parameter Store.
func (s *SSMService) DescribeParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeParametersWireCore(store, DescribeParametersWireInput{
		Parameters: req.Parameters,
		MaxResults: getIntParam(req, "MaxResults"),
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	params := make([]map[string]interface{}, 0, len(result.Parameters))
	for _, m := range result.Parameters {
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
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// GetParameterHistory retrieves the history of a parameter's versions.
func (s *SSMService) GetParameterHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	withDecryption := getBoolParam(req, "WithDecryption")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getParameterHistoryCore(store, GetParameterHistoryInput{
		Name:       req.GetParam("Name"),
		MaxResults: getIntParam(req, "MaxResults"),
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	versions := make([]map[string]interface{}, 0, len(result.History))
	for _, v := range result.History {
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
			"Policies":         policiesToResponse(v.Policies),
			"Tier":             string(v.Tier),
			"Value":            value,
			"DataType":         v.DataType,
		})
	}

	resp := map[string]interface{}{
		"Parameters": versions,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)
	return resp, nil
}

// LabelParameterVersion attaches labels to a specific version of a parameter.
func (s *SSMService) LabelParameterVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.labelParameterVersionCore(store, LabelParameterVersionInput{
		Name:             req.GetParam("Name"),
		ParameterVersion: int64(getIntParam(req, "ParameterVersion")),
		Labels:           parseStringList(req.Parameters, "Labels", "Labels.member."),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"InvalidLabels":    result.InvalidLabels,
		"ParameterVersion": result.ParameterVersion,
	}, nil
}

// UnlabelParameterVersion removes labels from a specific version of a parameter.
func (s *SSMService) UnlabelParameterVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.unlabelParameterVersionCore(store, UnlabelParameterVersionInput{
		Name:             req.GetParam("Name"),
		ParameterVersion: int64(getIntParam(req, "ParameterVersion")),
		Labels:           parseStringList(req.Parameters, "Labels", "Labels.member."),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"InvalidLabels": []string{},
		"RemovedLabels": result.RemovedLabels,
	}, nil
}
