package eventbridge

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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

// validateAuthParameters enforces that the supplied AuthParameters shape
// matches the declared AuthorizationType. AWS rejects the API call with a
// ValidationException when the wrong sub-object is populated or required
// credentials are missing.
func validateAuthParameters(authType string, p *eventsstore.AuthParameters) error {
	if p == nil {
		return awserrors.NewValidationException("AuthParameters are required for AuthorizationType " + authType)
	}
	switch authType {
	case "BASIC":
		if p.BasicAuthParameters == nil {
			return awserrors.NewValidationException("BasicAuthParameters are required when AuthorizationType is BASIC")
		}
		if p.BasicAuthParameters.Username == "" || p.BasicAuthParameters.Password == "" {
			return awserrors.NewValidationException("BasicAuthParameters.Username and BasicAuthParameters.Password are required")
		}
	case "OAUTH_CLIENT_CREDENTIALS":
		if p.OAuthParameters == nil {
			return awserrors.NewValidationException("OAuthParameters are required when AuthorizationType is OAUTH_CLIENT_CREDENTIALS")
		}
		if p.OAuthParameters.ClientParameters == nil {
			return awserrors.NewValidationException("OAuthParameters.ClientParameters are required")
		}
		if p.OAuthParameters.ClientParameters.ClientID == "" || p.OAuthParameters.ClientParameters.ClientSecret == "" {
			return awserrors.NewValidationException("ClientParameters.ClientID and ClientParameters.ClientSecret are required")
		}
		if p.OAuthParameters.AuthorizationEndpoint == "" {
			return awserrors.NewValidationException("OAuthParameters.AuthorizationEndpoint is required")
		}
		if p.OAuthParameters.HttpMethod == "" {
			return awserrors.NewValidationException("OAuthParameters.HttpMethod is required")
		}
	case "API_KEY":
		if p.ApiKeyAuthParameters == nil {
			return awserrors.NewValidationException("ApiKeyAuthParameters are required when AuthorizationType is API_KEY")
		}
		if p.ApiKeyAuthParameters.ApiKeyName == "" || p.ApiKeyAuthParameters.ApiKeyValue == "" {
			return awserrors.NewValidationException("ApiKeyAuthParameters.ApiKeyName and ApiKeyAuthParameters.ApiKeyValue are required")
		}
	}
	return nil
}

// CreateConnection creates a new EventBridge connection.
func (s *EventsService) CreateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}
	if !validateResourceName(name, "connection") {
		return nil, awserrors.NewValidationException("Connection name must match ^[.\\-_A-Za-z0-9]+$ and be 1-64 characters")
	}

	desc := request.GetStringParam(req.Parameters, "Description")
	if desc != "" && !validateDescription(desc) {
		return nil, awserrors.NewValidationException("Description must be at most 512 characters")
	}

	authType := request.GetParamLowerFirst(req.Parameters, "AuthorizationType")
	if authType == "" {
		return nil, awserrors.NewValidationException("AuthorizationType is required")
	}
	if !validAuthTypes[authType] {
		return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
	}

	authParams := parseAuthParameters(req.Parameters["AuthParameters"])
	if err := validateAuthParameters(authType, authParams); err != nil {
		return nil, err
	}

	connection := &eventsstore.Connection{
		Name:              name,
		AuthorizationType: authType,
		AuthParameters:    authParams,
		State:             eventsstore.ConnectionStateAuthorized,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, awserrors.NewValidationException("Description must be at most 512 characters")
		}
		connection.Description = desc
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok && kms != "" {
		if !validateKmsKeyIdentifier(kms) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
		}
		connection.KmsKeyIdentifier = kms
	}
	if icp := parseConnectivityParameters(req.Parameters, "InvocationConnectivityParameters"); icp != nil {
		connection.InvocationConnectivityParameters = icp
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateConnection(ctx, connection); err != nil {
		return nil, mapStoreError(err, name)
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
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
		return nil, err
	}

	allDests, err := store.ListApiDestinations(ctx, "", "", 1000, "")
	if err != nil {
		return nil, err
	}
	for _, d := range allDests.ApiDestinations {
		if d.ConnectionARN == connection.ARN {
			return nil, awserrors.NewValidationException("Connection '" + name + "' is in use by API destination '" + d.Name + "'")
		}
	}

	// Capture status fields before deletion so callers can confirm the
	// connection that was removed.
	snapshot := connectionToStatusMap(connection)

	if err := store.DeleteConnection(ctx, name); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// DescribeConnection returns information about a connection.
func (s *EventsService) DescribeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	return connectionToDescribeMap(connection), nil
}

// UpdateConnection updates an existing EventBridge connection.
func (s *EventsService) UpdateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, awserrors.NewValidationException("Description must be at most 512 characters")
		}
		connection.Description = desc
	}
	authChanged := false
	if authType, ok := req.Parameters["AuthorizationType"].(string); ok && authType != "" {
		if !validAuthTypes[authType] {
			return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
		}
		connection.AuthorizationType = authType
		authChanged = true
	}
	// Re-parse and validate AuthParameters when supplied. AuthorizationType
	// may be omitted on update (in which case the existing type is retained
	// but the caller must still supply credentials consistent with it).
	if rawAuth, ok := req.Parameters["AuthParameters"]; ok {
		authParams := parseAuthParameters(rawAuth)
		if err := validateAuthParameters(connection.AuthorizationType, authParams); err != nil {
			return nil, err
		}
		connection.AuthParameters = authParams
		connection.LastAuthorizedAt = time.Now().UTC()
		authChanged = true
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok && kms != "" {
		if !validateKmsKeyIdentifier(kms) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
		}
		connection.KmsKeyIdentifier = kms
	}
	if icp := parseConnectivityParameters(req.Parameters, "InvocationConnectivityParameters"); icp != nil {
		connection.InvocationConnectivityParameters = icp
	}

	connection.LastModifiedAt = time.Now().UTC()
	if authChanged {
		connection.State = eventsstore.ConnectionStateAuthorized
	}

	if err := store.UpdateConnection(ctx, connection); err != nil {
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
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeauthorizeConnection(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	connArn := connection.ARN
	createdAt := connection.CreatedAt.Unix()
	modifiedAt := connection.LastModifiedAt.Unix()
	authorizedAt := connection.LastAuthorizedAt.Unix()

	resp := map[string]interface{}{
		"ConnectionArn":   connArn,
		"ConnectionState": string(eventsstore.ConnectionStateDeauthorized),
	}
	resp["CreationTime"] = createdAt
	resp["LastModifiedTime"] = modifiedAt
	if !connection.LastAuthorizedAt.IsZero() {
		resp["LastAuthorizedTime"] = authorizedAt
	}
	return resp, nil
}

// ListConnections lists connections with optional filtering.
func (s *EventsService) ListConnections(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	stateStr := request.GetParamLowerFirst(req.Parameters, "ConnectionState")

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
	}

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListConnections(ctx, namePrefix, stateStr, limit, nextToken)
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
