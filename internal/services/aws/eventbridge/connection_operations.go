package eventbridge

import (
	"context"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// connectionToDescribeMap serialises a Connection for the
// DescribeConnectionResponse shape. Includes all fields per Smithy:
// ConnectionArn, Name, Description, InvocationConnectivityParameters,
// ConnectionState, StateReason, AuthorizationType, SecretArn,
// KmsKeyIdentifier, AuthParameters (redacted), CreationTime,
// LastModifiedTime, LastAuthorizedTime.
func connectionToDescribeMap(c *eventsstore.Connection) map[string]interface{} {
	result := map[string]interface{}{
		"ConnectionArn":     c.ARN,
		"Name":              c.Name,
		"ConnectionState":   string(c.State),
		"AuthorizationType": c.AuthorizationType,
		"CreationTime":      c.CreatedAt.Unix(),
	}
	if c.Description != "" {
		result["Description"] = c.Description
	}
	if c.StateReason != "" {
		result["StateReason"] = c.StateReason
	}
	if !c.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = c.LastModifiedAt.Unix()
	}
	if !c.LastAuthorizedAt.IsZero() {
		result["LastAuthorizedTime"] = c.LastAuthorizedAt.Unix()
	}
	if c.KmsKeyIdentifier != "" {
		result["KmsKeyIdentifier"] = c.KmsKeyIdentifier
	}
	if c.SecretArn != "" {
		result["SecretArn"] = c.SecretArn
	}
	if c.InvocationConnectivityParameters != nil && c.InvocationConnectivityParameters.ResourceParameters != nil {
		rp := c.InvocationConnectivityParameters.ResourceParameters
		if rp.ResourceConfigurationArn != "" || rp.ResourceAssociationArn != "" {
			inner := map[string]interface{}{}
			if rp.ResourceConfigurationArn != "" {
				inner["ResourceConfigurationArn"] = rp.ResourceConfigurationArn
			}
			if rp.ResourceAssociationArn != "" {
				inner["ResourceAssociationArn"] = rp.ResourceAssociationArn
			}
			result["InvocationConnectivityParameters"] = map[string]interface{}{
				"ResourceParameters": inner,
			}
		}
	}
	// AuthParameters are returned in redacted form: AWS returns the
	// parameter shape with credential fields blanked, so that callers can
	// confirm which auth configuration was applied without exposing
	// secrets. We follow the same contract.
	if c.AuthParameters != nil {
		result["AuthParameters"] = redactAuthParameters(c.AuthParameters)
	}
	return result
}

// connectionToStatusMap serialises a Connection for the
// DeleteConnectionResponse and DeauthorizeConnectionResponse shapes per
// Smithy: ConnectionArn, ConnectionState, CreationTime, LastModifiedTime,
// LastAuthorizedTime.
func connectionToStatusMap(c *eventsstore.Connection) map[string]interface{} {
	result := map[string]interface{}{
		"ConnectionArn":   c.ARN,
		"ConnectionState": string(c.State),
		"CreationTime":    c.CreatedAt.Unix(),
	}
	if !c.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = c.LastModifiedAt.Unix()
	}
	if !c.LastAuthorizedAt.IsZero() {
		result["LastAuthorizedTime"] = c.LastAuthorizedAt.Unix()
	}
	return result
}

// connectionToListMap serialises a Connection for the Connection list-item
// shape (Smithy: ConnectionArn, Name, ConnectionState, StateReason,
// AuthorizationType, CreationTime, LastModifiedTime, LastAuthorizedTime).
func connectionToListMap(c *eventsstore.Connection) map[string]interface{} {
	result := map[string]interface{}{
		"ConnectionArn":     c.ARN,
		"Name":              c.Name,
		"ConnectionState":   string(c.State),
		"AuthorizationType": c.AuthorizationType,
		"CreationTime":      c.CreatedAt.Unix(),
	}
	if c.StateReason != "" {
		result["StateReason"] = c.StateReason
	}
	if !c.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = c.LastModifiedAt.Unix()
	}
	if !c.LastAuthorizedAt.IsZero() {
		result["LastAuthorizedTime"] = c.LastAuthorizedAt.Unix()
	}
	return result
}

// redactAuthParameters mirrors AWS DescribeConnection semantics: the
// AuthParameters shape is preserved but every credential-bearing field is
// replaced with an empty string. This proves which sub-parameters were set
// without leaking stored secrets.
func redactAuthParameters(p *eventsstore.AuthParameters) map[string]interface{} {
	out := map[string]interface{}{}
	if p.BasicAuthParameters != nil {
		out["BasicAuthParameters"] = map[string]interface{}{
			"Username": "",
			"Password": "",
		}
	}
	if p.OAuthParameters != nil {
		oauth := map[string]interface{}{
			"AuthorizationEndpoint": p.OAuthParameters.AuthorizationEndpoint,
			"HttpMethod":            p.OAuthParameters.HttpMethod,
		}
		if p.OAuthParameters.ClientParameters != nil {
			oauth["ClientParameters"] = map[string]interface{}{
				"ClientID":     "",
				"ClientSecret": "",
			}
		}
		if p.OAuthParameters.OAuthHttpParameters != nil {
			oauth["OAuthHttpParameters"] = redactHttpParameters(p.OAuthParameters.OAuthHttpParameters)
		}
		out["OAuthParameters"] = oauth
	}
	if p.ApiKeyAuthParameters != nil {
		out["ApiKeyAuthParameters"] = map[string]interface{}{
			"ApiKeyName":  "",
			"ApiKeyValue": "",
		}
	}
	if p.InvocationHttpParameters != nil {
		out["InvocationHttpParameters"] = redactHttpParameters(p.InvocationHttpParameters)
	}
	return out
}

// redactHttpParameters returns ConnectionHttpParameters with header/query
// values blanked, preserving only the keys. Body parameter keys are removed
// entirely (their values are themselves the secret payload).
func redactHttpParameters(p *eventsstore.ConnectionHttpParameters) map[string]interface{} {
	out := map[string]interface{}{}
	if len(p.HeaderParameters) > 0 {
		hdrs := make(map[string]string, len(p.HeaderParameters))
		for k := range p.HeaderParameters {
			hdrs[k] = ""
		}
		out["HeaderParameters"] = hdrs
	}
	if len(p.QueryStringParameters) > 0 {
		qs := make(map[string]string, len(p.QueryStringParameters))
		for k := range p.QueryStringParameters {
			qs[k] = ""
		}
		out["QueryStringParameters"] = qs
	}
	if len(p.BodyParameters) > 0 {
		out["BodyParameters"] = []string{}
	}
	return out
}

// parseAuthParameters extracts the AuthParameters shape from a PutConnection /
// UpdateConnection request payload. Returns nil if no auth parameters are
// present (the caller validates that this is consistent with the declared
// AuthorizationType).
func parseAuthParameters(raw interface{}) *eventsstore.AuthParameters {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	out := &eventsstore.AuthParameters{}

	if basic, ok := m["BasicAuthParameters"].(map[string]interface{}); ok {
		out.BasicAuthParameters = &eventsstore.BasicAuthParameters{
			Username: getStringField(basic, "Username"),
			Password: getStringField(basic, "Password"),
		}
	}
	if oauth, ok := m["OAuthParameters"].(map[string]interface{}); ok {
		o := &eventsstore.OAuthParameters{
			AuthorizationEndpoint: getStringField(oauth, "AuthorizationEndpoint"),
			HttpMethod:            getStringField(oauth, "HttpMethod"),
		}
		if cp, ok := oauth["ClientParameters"].(map[string]interface{}); ok {
			o.ClientParameters = &eventsstore.OAuthClientParameters{
				ClientID:     getStringField(cp, "ClientID"),
				ClientSecret: getStringField(cp, "ClientSecret"),
			}
		}
		if hp, ok := oauth["OAuthHttpParameters"].(map[string]interface{}); ok {
			o.OAuthHttpParameters = parseConnectionHttpParameters(hp)
		}
		out.OAuthParameters = o
	}
	if api, ok := m["ApiKeyAuthParameters"].(map[string]interface{}); ok {
		out.ApiKeyAuthParameters = &eventsstore.ApiKeyAuthParameters{
			ApiKeyName:  getStringField(api, "ApiKeyName"),
			ApiKeyValue: getStringField(api, "ApiKeyValue"),
		}
	}
	if hp, ok := m["InvocationHttpParameters"].(map[string]interface{}); ok {
		out.InvocationHttpParameters = parseConnectionHttpParameters(hp)
	}
	return out
}

// parseConnectionHttpParameters builds a ConnectionHttpParameters from a
// request map. HeaderParameters and QueryStringParameters are map[string]string
// in the Smithy model; BodyParameters is a list of opaque strings.
func parseConnectionHttpParameters(m map[string]interface{}) *eventsstore.ConnectionHttpParameters {
	out := &eventsstore.ConnectionHttpParameters{}
	if hdrs, ok := m["HeaderParameters"].(map[string]interface{}); ok {
		out.HeaderParameters = make(map[string]string, len(hdrs))
		for k, v := range hdrs {
			if s, ok := v.(string); ok {
				out.HeaderParameters[k] = s
			}
		}
	}
	if qs, ok := m["QueryStringParameters"].(map[string]interface{}); ok {
		out.QueryStringParameters = make(map[string]string, len(qs))
		for k, v := range qs {
			if s, ok := v.(string); ok {
				out.QueryStringParameters[k] = s
			}
		}
	}
	if bodies, ok := m["BodyParameters"].([]interface{}); ok {
		for _, b := range bodies {
			if s, ok := b.(string); ok {
				out.BodyParameters = append(out.BodyParameters, s)
			}
		}
	}
	return out
}

// getStringField returns the string value at key in m, or empty string if
// absent or non-string.
func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// parseCreateConnectionInput reads the CreateConnection wire request into
// the transport-agnostic Core input.
func parseCreateConnectionInput(req *request.ParsedRequest) CreateConnectionInput {
	input := CreateConnectionInput{
		Name:              request.GetParamLowerFirst(req.Parameters, "Name"),
		Description:       request.GetStringParam(req.Parameters, "Description"),
		AuthorizationType: request.GetParamLowerFirst(req.Parameters, "AuthorizationType"),
		AuthParameters:    parseAuthParameters(req.Parameters["AuthParameters"]),
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		input.KmsKeyIdentifierSet = true
		input.KmsKeyIdentifier = kms
	}
	input.InvocationConnectivityParameters = parseConnectivityParameters(req.Parameters, "InvocationConnectivityParameters")
	return input
}

// parseUpdateConnectionInput reads the UpdateConnection wire request into
// the transport-agnostic Core input.
func parseUpdateConnectionInput(req *request.ParsedRequest) UpdateConnectionInput {
	input := UpdateConnectionInput{
		Name: request.GetParamLowerFirst(req.Parameters, "Name"),
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	if authType, ok := req.Parameters["AuthorizationType"].(string); ok {
		input.AuthorizationTypeSet = true
		input.AuthorizationType = authType
	}
	if rawAuth, ok := req.Parameters["AuthParameters"]; ok {
		input.AuthParametersSet = true
		input.AuthParameters = parseAuthParameters(rawAuth)
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		input.KmsKeyIdentifierSet = true
		input.KmsKeyIdentifier = kms
	}
	input.InvocationConnectivityParameters = parseConnectivityParameters(req.Parameters, "InvocationConnectivityParameters")
	return input
}

// CreateConnection creates a new EventBridge connection.
func (s *EventsService) CreateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseCreateConnectionInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := s.createConnectionCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return map[string]interface{}{
		"ConnectionArn":    connection.ARN,
		"ConnectionState":  string(connection.State),
		"CreationTime":     connection.CreatedAt.Unix(),
		"LastModifiedTime": now.Unix(),
	}, nil
}

// DeleteConnection deletes an EventBridge connection.
func (s *EventsService) DeleteConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The status fields of the removed connection are reported back to the
	// caller so they can confirm which connection was deleted.
	connection, err := s.deleteConnectionCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	return connectionToStatusMap(connection), nil
}

// DescribeConnection returns information about a connection.
func (s *EventsService) DescribeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := s.getConnectionCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	return connectionToDescribeMap(connection), nil
}

// UpdateConnection updates an existing EventBridge connection.
func (s *EventsService) UpdateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseUpdateConnectionInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := s.updateConnectionCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"ConnectionArn":    connection.ARN,
		"ConnectionState":  string(connection.State),
		"CreationTime":     connection.CreatedAt.Unix(),
		"LastModifiedTime": connection.LastModifiedAt.Unix(),
	}
	if !connection.LastAuthorizedAt.IsZero() {
		resp["LastAuthorizedTime"] = connection.LastAuthorizedAt.Unix()
	}
	return resp, nil
}

// DeauthorizeConnection deauthorises an EventBridge connection, revoking its authorisation.
func (s *EventsService) DeauthorizeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := s.deauthorizeConnectionCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"ConnectionArn":   connection.ARN,
		"ConnectionState": string(eventsstore.ConnectionStateDeauthorized),
	}
	resp["CreationTime"] = connection.CreatedAt.Unix()
	resp["LastModifiedTime"] = connection.LastModifiedAt.Unix()
	if !connection.LastAuthorizedAt.IsZero() {
		resp["LastAuthorizedTime"] = connection.LastAuthorizedAt.Unix()
	}
	return resp, nil
}

// ListConnections lists connections with optional filtering.
func (s *EventsService) ListConnections(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := ListConnectionsInput{
		NamePrefix: request.GetParamLowerFirst(req.Parameters, "NamePrefix"),
		State:      request.GetParamLowerFirst(req.Parameters, "ConnectionState"),
		Limit:      int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listConnectionsCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	connections := make([]map[string]interface{}, 0, len(result.Connections))
	for _, conn := range result.Connections {
		connections = append(connections, connectionToListMap(conn))
	}

	resp := map[string]interface{}{
		"Connections": connections,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)

	return resp, nil
}

// parseConnectivityParameters reads an InvocationConnectivityParameters
// (top-level request) or ConnectivityParameters (AuthParameters child)
// member from the request. Both map to the Smithy
// ConnectivityResourceParameters shape, which nests ResourceParameters →
// ResourceConfigurationArn. Returns nil if absent or empty.
func parseConnectivityParameters(params map[string]interface{}, key string) *eventsstore.ConnectivityResourceParameters {
	conn, ok := params[key].(map[string]interface{})
	if !ok {
		return nil
	}
	rp, ok := conn["ResourceParameters"].(map[string]interface{})
	if !ok {
		return nil
	}
	rcArn := getStringField(rp, "ResourceConfigurationArn")
	if rcArn == "" {
		return nil
	}
	return &eventsstore.ConnectivityResourceParameters{
		ResourceParameters: &eventsstore.ResourceConfiguration{
			ResourceConfigurationArn: rcArn,
		},
	}
}
