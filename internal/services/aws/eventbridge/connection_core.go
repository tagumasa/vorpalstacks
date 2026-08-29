package eventbridge

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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

// createConnectionCore validates input and creates the connection.
func (s *EventsService) createConnectionCore(ctx context.Context, store *eventsstore.EventsStore, input CreateConnectionInput) (*eventsstore.Connection, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}
	if !validateResourceName(input.Name, "connection") {
		return nil, awserrors.NewValidationException("Connection name must match ^[.\\-_A-Za-z0-9]+$ and be 1-64 characters")
	}

	if input.Description != "" && !validateDescription(input.Description) {
		return nil, errDescriptionTooLong()
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

	connection := &eventsstore.Connection{
		Name:              input.Name,
		AuthorizationType: input.AuthorizationType,
		AuthParameters:    input.AuthParameters,
		State:             eventsstore.ConnectionStateAuthorized,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		connection.Description = input.Description
	}
	if input.KmsKeyIdentifierSet && input.KmsKeyIdentifier != "" {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
		}
		connection.KmsKeyIdentifier = input.KmsKeyIdentifier
	}
	if input.InvocationConnectivityParameters != nil {
		connection.InvocationConnectivityParameters = input.InvocationConnectivityParameters
	}

	if err := store.CreateConnection(ctx, connection); err != nil {
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

	allDests, err := store.ListApiDestinations(ctx, "", "", 1000, "")
	if err != nil {
		return nil, err
	}
	for _, d := range allDests.ApiDestinations {
		if d.ConnectionARN == connection.ARN {
			return nil, awserrors.NewValidationException("Connection '" + name + "' is in use by API destination '" + d.Name + "'")
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
// stored connection and persists the update.
func (s *EventsService) updateConnectionCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateConnectionInput) (*eventsstore.Connection, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	connection, err := store.GetConnection(ctx, input.Name)
	if err != nil {
		return nil, mapStoreError(err, input.Name)
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		connection.Description = input.Description
	}
	authChanged := false
	if input.AuthorizationTypeSet && input.AuthorizationType != "" {
		if !validAuthTypes[input.AuthorizationType] {
			return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
		}
		connection.AuthorizationType = input.AuthorizationType
		authChanged = true
	}
	// Re-validate AuthParameters when supplied. AuthorizationType may be
	// omitted on update (in which case the existing type is retained but the
	// caller must still supply credentials consistent with it).
	if input.AuthParametersSet {
		if err := validateAuthParameters(connection.AuthorizationType, input.AuthParameters); err != nil {
			return nil, err
		}
		connection.AuthParameters = input.AuthParameters
		connection.LastAuthorizedAt = time.Now().UTC()
		authChanged = true
	}
	if input.KmsKeyIdentifierSet && input.KmsKeyIdentifier != "" {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
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

	if err := store.UpdateConnection(ctx, connection); err != nil {
		return nil, err
	}
	return connection, nil
}

// deauthorizeConnectionCore validates input, revokes the connection's
// authorisation and returns the re-read record for the response.
func (s *EventsService) deauthorizeConnectionCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.Connection, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	if err := store.DeauthorizeConnection(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	return connection, nil
}

// listConnectionsCore applies the documented limit window and lists the
// connections.
func (s *EventsService) listConnectionsCore(ctx context.Context, store *eventsstore.EventsStore, input ListConnectionsInput) (*eventsstore.ConnectionListResult, error) {
	limit := input.Limit
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
	}
	return store.ListConnections(ctx, input.NamePrefix, input.State, limit, input.NextToken)
}
