package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/invokers"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// CreateConnectionInput carries the parameters for CreateConnection.
type CreateConnectionInput struct {
	Name                             string
	DescriptionSet                   bool
	Description                      string
	AuthorizationType                string
	AuthParameters                   *eventsstore.AuthParameters
	KmsKeyIdentifierSet              bool
	KmsKeyIdentifier                 string
	InvocationConnectivityParameters *eventsstore.ConnectivityResourceParameters
}

// UpdateConnectionInput carries the parameters for UpdateConnection.
// AuthParametersSet distinguishes an omitted AuthParameters member from a
// provided one that parsed to nil (a non-map wire value), because the
// validator rejects the latter.
type UpdateConnectionInput struct {
	Name                             string
	DescriptionSet                   bool
	Description                      string
	AuthorizationTypeSet             bool
	AuthorizationType                string
	AuthParametersSet                bool
	AuthParameters                   *eventsstore.AuthParameters
	KmsKeyIdentifierSet              bool
	KmsKeyIdentifier                 string
	InvocationConnectivityParameters *eventsstore.ConnectivityResourceParameters
}

// ListConnectionsInput carries the parameters for ListConnections.
type ListConnectionsInput struct {
	NamePrefix string
	State      string
	Limit      int32
	NextToken  string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// createConnectionCore validates input and creates the connection. The
// full authorization parameters are stored in a service-owned Secrets
// Manager secret (the AWS design: "When you create a connection and add
// authorization parameters, EventBridge creates a secret in AWS Secrets
// Manager", user guide) and the record keeps only the non-credential
// halves plus the secret's ARN — the plaintext never lands in the events
// keyspace.
func (s *EventsService) createConnectionCore(ctx context.Context, store *eventsstore.EventsStore, input CreateConnectionInput) (*eventsstore.Connection, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}
	if !validateResourceName(input.Name, "connection") {
		return nil, awserrors.NewValidationException("Connection name must match ^[.\\-_A-Za-z0-9]+$ and be 1-64 characters")
	}

	if input.AuthorizationType == "" {
		return nil, awserrors.NewValidationException("AuthorizationType is required")
	}
	if !validAuthTypes[input.AuthorizationType] {
		return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
	}

	if err := validateAuthParameters(input.AuthorizationType, input.AuthParameters); err != nil {
		return nil, err
	}

	// Every input validation runs while the request is still side-effect
	// free: the credential secret below is an external resource, and the
	// compensation delete covers a store-write failure alone — a validation
	// failure after the creation would leak the secret on every retry.
	connection := &eventsstore.Connection{
		Name:              input.Name,
		AuthorizationType: input.AuthorizationType,
		AuthParameters:    redactCredentialsForStorage(input.AuthParameters),
	}
	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		connection.Description = input.Description
	}
	if input.KmsKeyIdentifierSet && input.KmsKeyIdentifier != "" {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a key ARN, key ID, key alias or key alias ARN")
		}
		connection.KmsKeyIdentifier = input.KmsKeyIdentifier
	}
	if input.InvocationConnectivityParameters != nil {
		connection.InvocationConnectivityParameters = input.InvocationConnectivityParameters
	}

	invoker := s.secretsInvoker()
	if invoker == nil {
		return nil, errSecretsManagerUnavailable()
	}
	secretJSON, err := marshalAuthParameters(input.AuthParameters)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalException", "encode connection authorization parameters: "+err.Error(), 500)
	}
	secretName := "events/connection/" + input.Name + "/" + uuid.NewString()
	secretARN, err := invoker.CreateServiceSecret(ctx, store.Region(), secretName, secretJSON, "EventBridge connection "+input.Name)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalException", "store connection credentials: "+err.Error(), 500)
	}
	connection.SecretArn = secretARN
	connection.State = eventsstore.ConnectionStateAuthorized
	// Establishing the credentials is an authorization event: the
	// CreateConnectionResponse shape carries LastAuthorizedTime.
	connection.LastAuthorizedAt = time.Now().UTC()

	if err := store.CreateConnection(ctx, connection); err != nil {
		// The record never landed; the credential secret must not outlive it.
		_ = invoker.DeleteServiceSecret(ctx, store.Region(), secretARN)
		return nil, mapStoreError(err, input.Name)
	}
	return connection, nil
}

// deleteConnectionCore validates input, rejects deletion while an API
// destination still references the connection, and deletes it. The
// pre-delete record is returned so the caller can report the removed
// connection's status fields.
func (s *EventsService) deleteConnectionCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.Connection, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
		return nil, err
	}

	// The in-use scan paginates: a single capped page would let a
	// connection past the page edge be deleted while API destinations
	// still reference it.
	destToken := ""
	for {
		dests, err := store.ListApiDestinations(ctx, "", "", 100, destToken)
		if err != nil {
			return nil, err
		}
		for _, d := range dests.ApiDestinations {
			if d.ConnectionARN == connection.ARN {
				return nil, awserrors.NewValidationException("Connection '" + name + "' is in use by API destination '" + d.Name + "'")
			}
		}
		if dests.NextToken == "" {
			break
		}
		destToken = dests.NextToken
	}

	// The credential secret follows the connection's lifecycle. It is
	// removed before the record: a mid-way failure leaves a retryable
	// record whose secret reference no longer resolves (the delete
	// tolerates an already-missing secret).
	if connection.SecretArn != "" {
		if invoker := s.secretsInvoker(); invoker != nil {
			if err := invoker.DeleteServiceSecret(ctx, store.Region(), connection.SecretArn); err != nil {
				return nil, awserrors.NewAWSError("InternalException", "delete connection credentials: "+err.Error(), 500)
			}
		}
	}

	if err := store.DeleteConnection(ctx, name); err != nil {
		return nil, err
	}

	return connection, nil
}

// getConnectionCore validates input and fetches the connection.
func (s *EventsService) getConnectionCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.Connection, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}
	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	return connection, nil
}

// updateConnectionCore validates input, merges the provided members onto the
// stored connection and persists the update through the atomic record
// mutation (a concurrent deauthorization cannot be overwritten by a stale
// pre-read record, and vice versa).
func (s *EventsService) updateConnectionCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateConnectionInput) (*eventsstore.Connection, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	var updated *eventsstore.Connection
	if err := store.MutateConnection(ctx, input.Name, func(connection *eventsstore.Connection) error {
		if input.DescriptionSet {
			if !validateDescription(input.Description) {
				return errDescriptionTooLong()
			}
			connection.Description = input.Description
		}
		authChanged := false
		if input.AuthorizationTypeSet && input.AuthorizationType != "" {
			if !validAuthTypes[input.AuthorizationType] {
				return awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
			}
			// "For connections to public APIs, EventBridge re-authorizes the
			// connection if you: ... Update the connection authorization
			// method and authorization parameters" (user guide, Updating
			// connections) — the method's credentials arrive with it. A
			// type-only change would flip the record to the new type while
			// the credential secret and the stored AuthParameters still
			// describe the previous type, so it is rejected instead of
			// marked authorized.
			if input.AuthorizationType != connection.AuthorizationType && !input.AuthParametersSet {
				return awserrors.NewValidationException(
					"Changing AuthorizationType requires AuthParameters with the credentials for the new authorization type")
			}
			connection.AuthorizationType = input.AuthorizationType
		}
		// Re-validate AuthParameters when supplied. AuthorizationType may be
		// omitted on update (in which case the existing type is retained but the
		// caller must still supply credentials consistent with it). Only the
		// credential write re-marks the connection AUTHORIZED: the
		// re-authorization trigger is the parameter update, so an update that
		// changes no credentials (including an AuthorizationType echo) never
		// re-authorizes a deauthorized connection.
		if input.AuthParametersSet {
			if err := validateAuthParameters(connection.AuthorizationType, input.AuthParameters); err != nil {
				return err
			}
			// The full credentials live in the connection's secret; the
			// record keeps the non-credential halves. Writing the secret
			// inside the mutation keeps an unavailable Secrets Manager an
			// atomic abort: the record is left untouched.
			if err := s.storeConnectionCredentials(ctx, store, connection, input.AuthParameters); err != nil {
				return err
			}
			connection.AuthParameters = redactCredentialsForStorage(input.AuthParameters)
			connection.LastAuthorizedAt = time.Now().UTC()
			authChanged = true
		}
		if input.KmsKeyIdentifierSet && input.KmsKeyIdentifier != "" {
			if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
				return awserrors.NewValidationException("KmsKeyIdentifier must be a key ARN, key ID, key alias or key alias ARN")
			}
			connection.KmsKeyIdentifier = input.KmsKeyIdentifier
		}
		if input.InvocationConnectivityParameters != nil {
			connection.InvocationConnectivityParameters = input.InvocationConnectivityParameters
		}

		connection.LastModifiedAt = time.Now().UTC()
		if authChanged {
			connection.State = eventsstore.ConnectionStateAuthorized
		}
		updated = connection
		return nil
	}); err != nil {
		return nil, mapStoreError(err, input.Name)
	}
	return updated, nil
}

// deauthorizeConnectionCore validates input, revokes the connection's
// authorisation — credentials and their secret go with it ("Removes all
// authorization parameters from the connection") — and returns the
// deauthorized record for the response.
func (s *EventsService) deauthorizeConnectionCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.Connection, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	var deauthorized *eventsstore.Connection
	if err := store.MutateConnection(ctx, name, func(connection *eventsstore.Connection) error {
		// The secret is deleted inside the mutation, mirroring the update
		// path's atomicity ("writing the secret inside the mutation keeps
		// an unavailable Secrets Manager an atomic abort"): a deletion
		// failure leaves the record — and its claim on the secret —
		// untouched, and the read that drives the deletion cannot
		// interleave with a concurrent update minting a fresh secret the
		// record would then no longer name.
		if connection.SecretArn != "" {
			if invoker := s.secretsInvoker(); invoker != nil {
				if err := invoker.DeleteServiceSecret(ctx, store.Region(), connection.SecretArn); err != nil {
					return awserrors.NewAWSError("InternalException", "remove connection credentials: "+err.Error(), 500)
				}
			}
		}
		connection.State = eventsstore.ConnectionStateDeauthorized
		connection.StateReason = "User initiated deauthorization"
		connection.AuthParameters = nil
		connection.SecretArn = ""
		deauthorized = connection
		return nil
	}); err != nil {
		return nil, mapStoreError(err, name)
	}
	return deauthorized, nil
}

// listConnectionsCore validates the query and lists the connections: the
// limit window, the NamePrefix charset (both per the model) and the State
// filter's enum membership.
func (s *EventsService) listConnectionsCore(ctx context.Context, store *eventsstore.EventsStore, input ListConnectionsInput) (*eventsstore.ConnectionListResult, error) {
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "connection"); err != nil {
			return nil, err
		}
	}
	if input.State != "" && !eventsstore.IsValidConnectionState(eventsstore.ConnectionState(input.State)) {
		return nil, awserrors.NewValidationException(
			"State must be a member of the ConnectionState enum: " + strings.Join(eventsstore.ConnectionStateVocabulary(), ", "))
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}
	return store.ListConnections(ctx, input.NamePrefix, input.State, limit, input.NextToken)
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
		// The OAuth token request supports GET, POST and PUT only (the
		// ConnectionOAuthHttpMethod enum is narrower than the API
		// destination method enum).
		if !validOAuthHttpMethods[p.OAuthParameters.HttpMethod] {
			return awserrors.NewValidationException("OAuthParameters.HttpMethod must be one of: GET, POST, PUT")
		}
	case "API_KEY":
		if p.ApiKeyAuthParameters == nil {
			return awserrors.NewValidationException("ApiKeyAuthParameters are required when AuthorizationType is API_KEY")
		}
		if p.ApiKeyAuthParameters.ApiKeyName == "" || p.ApiKeyAuthParameters.ApiKeyValue == "" {
			return awserrors.NewValidationException("ApiKeyAuthParameters.ApiKeyName and ApiKeyAuthParameters.ApiKeyValue are required")
		}
	}
	// The HTTP-parameter constraints are family-wide: the same
	// ConnectionHttpParameters shape backs the invocation parameters of
	// every authorization type and the OAuth token-request parameters,
	// so the bounds hang off the switch rather than a single arm.
	if p.InvocationHttpParameters != nil {
		if err := validateConnectionHttpParameters(p.InvocationHttpParameters); err != nil {
			return err
		}
	}
	if p.OAuthParameters != nil && p.OAuthParameters.OAuthHttpParameters != nil {
		if err := validateConnectionHttpParameters(p.OAuthParameters.OAuthHttpParameters); err != nil {
			return err
		}
	}
	return nil
}

// validateConnectionHttpParameters enforces the Smithy traits on the
// ConnectionHttpParameters shape: each family list is bounded at 100
// entries, and the header and query-string keys and values carry
// @length(0,512) plus their character patterns. The length basis is
// Unicode code points (utf8.RuneCountInString), matching how the model's
// @length trait measures strings; the header family's patterns are ASCII
// anyway, so the basis only distinguishes multibyte query-string members.
// The body family's Key and Value carry no modelled traits — its list cap
// is the only bound.
func validateConnectionHttpParameters(p *eventsstore.ConnectionHttpParameters) error {
	if len(p.HeaderParameters) > eventsstore.ConnectionHttpParametersMaxEntries {
		return awserrors.NewValidationException(fmt.Sprintf("HeaderParameters must have at most %d entries", eventsstore.ConnectionHttpParametersMaxEntries))
	}
	if len(p.QueryStringParameters) > eventsstore.ConnectionHttpParametersMaxEntries {
		return awserrors.NewValidationException(fmt.Sprintf("QueryStringParameters must have at most %d entries", eventsstore.ConnectionHttpParametersMaxEntries))
	}
	if len(p.BodyParameters) > eventsstore.ConnectionHttpParametersMaxEntries {
		return awserrors.NewValidationException(fmt.Sprintf("BodyParameters must have at most %d entries", eventsstore.ConnectionHttpParametersMaxEntries))
	}
	for i, hp := range p.HeaderParameters {
		if utf8.RuneCountInString(hp.Key) > eventsstore.ConnectionParameterMaxLength {
			return awserrors.NewValidationException(fmt.Sprintf("HeaderParameters[%d].Key must be at most %d characters", i, eventsstore.ConnectionParameterMaxLength))
		}
		if !connectionHeaderKeyPattern.MatchString(hp.Key) {
			return awserrors.NewValidationException(fmt.Sprintf("HeaderParameters[%d].Key contains invalid characters", i))
		}
		if utf8.RuneCountInString(hp.Value) > eventsstore.ConnectionParameterMaxLength {
			return awserrors.NewValidationException(fmt.Sprintf("HeaderParameters[%d].Value must be at most %d characters", i, eventsstore.ConnectionParameterMaxLength))
		}
		if !connectionHeaderValuePattern.MatchString(hp.Value) {
			return awserrors.NewValidationException(fmt.Sprintf("HeaderParameters[%d].Value contains invalid characters", i))
		}
	}
	for i, qp := range p.QueryStringParameters {
		if utf8.RuneCountInString(qp.Key) > eventsstore.ConnectionParameterMaxLength {
			return awserrors.NewValidationException(fmt.Sprintf("QueryStringParameters[%d].Key must be at most %d characters", i, eventsstore.ConnectionParameterMaxLength))
		}
		if !connectionQueryStringKeyPattern.MatchString(qp.Key) {
			return awserrors.NewValidationException(fmt.Sprintf("QueryStringParameters[%d].Key contains invalid characters", i))
		}
		if utf8.RuneCountInString(qp.Value) > eventsstore.ConnectionParameterMaxLength {
			return awserrors.NewValidationException(fmt.Sprintf("QueryStringParameters[%d].Value must be at most %d characters", i, eventsstore.ConnectionParameterMaxLength))
		}
		if !connectionQueryStringValuePattern.MatchString(qp.Value) {
			return awserrors.NewValidationException(fmt.Sprintf("QueryStringParameters[%d].Value contains invalid characters", i))
		}
	}
	return nil
}

// secretsInvoker returns the Secrets Manager invoker for connection
// credentials, or nil when the bus or the providing service is absent.
func (s *EventsService) secretsInvoker() invokers.SecretsManagerInvoker {
	if s.bus == nil {
		return nil
	}
	return s.bus.SecretsManagerInvoker()
}

// errSecretsManagerUnavailable is the error raised when connection
// credentials cannot be stored: AWS likewise requires Secrets Manager
// access to create or update a connection ("To successfully create or
// update a connection, you must use an account that has permission to use
// Secrets Manager", user guide).
func errSecretsManagerUnavailable() error {
	return awserrors.NewAWSError("InternalException",
		"Storing connection credentials requires the Secrets Manager service", 500)
}

// storeConnectionCredentials persists the supplied authorization
// parameters as the connection's credential secret: an existing secret is
// updated in place; a connection without one (its credentials were
// established by an update after deauthorization) gets a fresh secret and
// carries the new ARN on the record.
func (s *EventsService) storeConnectionCredentials(ctx context.Context, store *eventsstore.EventsStore, connection *eventsstore.Connection, params *eventsstore.AuthParameters) error {
	invoker := s.secretsInvoker()
	if invoker == nil {
		return errSecretsManagerUnavailable()
	}
	secretJSON, err := marshalAuthParameters(params)
	if err != nil {
		return awserrors.NewAWSError("InternalException", "encode connection authorization parameters: "+err.Error(), 500)
	}
	if connection.SecretArn != "" {
		if err := invoker.UpdateServiceSecretString(ctx, store.Region(), connection.SecretArn, secretJSON); err != nil {
			return awserrors.NewAWSError("InternalException", "update connection credentials: "+err.Error(), 500)
		}
		return nil
	}
	secretName := "events/connection/" + connection.Name + "/" + uuid.NewString()
	secretARN, err := invoker.CreateServiceSecret(ctx, store.Region(), secretName, secretJSON, "EventBridge connection "+connection.Name)
	if err != nil {
		return awserrors.NewAWSError("InternalException", "store connection credentials: "+err.Error(), 500)
	}
	connection.SecretArn = secretARN
	return nil
}

// marshalAuthParameters serialises the full authorization parameters into
// the credential secret's payload.
func marshalAuthParameters(p *eventsstore.AuthParameters) (string, error) {
	if p == nil {
		return "{}", nil
	}
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// redactCredentialsForStorage copies the authorization parameters with
// every credential blanked: the Basic password, the OAuth client secret,
// the API key value, and the values of parameter entries marked
// IsValueSecret. The result is what the connection record keeps — the
// non-credential halves Describe responses are built from; the full tree
// lives only in the credential secret.
func redactCredentialsForStorage(p *eventsstore.AuthParameters) *eventsstore.AuthParameters {
	if p == nil {
		return nil
	}
	out := &eventsstore.AuthParameters{}
	if p.BasicAuthParameters != nil {
		out.BasicAuthParameters = &eventsstore.BasicAuthParameters{
			Username: p.BasicAuthParameters.Username,
		}
	}
	if p.OAuthParameters != nil {
		oauth := &eventsstore.OAuthParameters{
			AuthorizationEndpoint: p.OAuthParameters.AuthorizationEndpoint,
			HttpMethod:            p.OAuthParameters.HttpMethod,
		}
		if p.OAuthParameters.ClientParameters != nil {
			oauth.ClientParameters = &eventsstore.OAuthClientParameters{
				ClientID: p.OAuthParameters.ClientParameters.ClientID,
			}
		}
		if p.OAuthParameters.OAuthHttpParameters != nil {
			oauth.OAuthHttpParameters = redactHttpParametersForStorage(p.OAuthParameters.OAuthHttpParameters)
		}
		out.OAuthParameters = oauth
	}
	if p.ApiKeyAuthParameters != nil {
		out.ApiKeyAuthParameters = &eventsstore.ApiKeyAuthParameters{
			ApiKeyName: p.ApiKeyAuthParameters.ApiKeyName,
		}
	}
	if p.InvocationHttpParameters != nil {
		out.InvocationHttpParameters = redactHttpParametersForStorage(p.InvocationHttpParameters)
	}
	return out
}

// redactHttpParametersForStorage blanks the values of secret-marked
// parameter entries, preserving keys, non-secret values and the
// IsValueSecret flags.
func redactHttpParametersForStorage(p *eventsstore.ConnectionHttpParameters) *eventsstore.ConnectionHttpParameters {
	if p == nil {
		return nil
	}
	out := &eventsstore.ConnectionHttpParameters{}
	for _, h := range p.HeaderParameters {
		if h.IsValueSecret {
			h.Value = ""
		}
		out.HeaderParameters = append(out.HeaderParameters, h)
	}
	for _, q := range p.QueryStringParameters {
		if q.IsValueSecret {
			q.Value = ""
		}
		out.QueryStringParameters = append(out.QueryStringParameters, q)
	}
	for _, b := range p.BodyParameters {
		if b.IsValueSecret {
			b.Value = ""
		}
		out.BodyParameters = append(out.BodyParameters, b)
	}
	return out
}
