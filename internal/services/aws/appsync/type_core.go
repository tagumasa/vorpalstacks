package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// createTypeInput carries the parsed CreateType request payload.
type createTypeInput struct {
	ApiId       string
	Definition  string
	Format      string
	Description string
}

// updateTypeInput carries the parsed UpdateType request payload. The Has*
// flags distinguish an explicitly supplied member from an omitted one.
type updateTypeInput struct {
	ApiId          string
	TypeName       string
	Format         string
	Definition     string
	HasDefinition  bool
	Description    string
	HasDescription bool
}

// createTypeCore validates the request and persists a new type definition.
func (s *AppSyncService) createTypeCore(store *appsyncstore.AppSyncStore, in createTypeInput) (*appsyncstore.Type, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if err := validateGraphqlApiExists(store, in.ApiId); err != nil {
		return nil, err
	}

	if in.Definition == "" {
		return nil, NewBadRequestException("definition is required")
	}

	if in.Format == "" {
		return nil, NewBadRequestException("format is required")
	}

	if !validateTypeFormat(in.Format) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid format: %s. Valid values: SDL, JSON", in.Format))
	}

	t := &appsyncstore.Type{
		ApiId:       in.ApiId,
		Definition:  in.Definition,
		Format:      in.Format,
		Description: in.Description,
	}

	created, err := store.CreateType(t)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return created, nil
}

// getTypeCore fetches a type definition by API ID and type name.
func (s *AppSyncService) getTypeCore(store *appsyncstore.AppSyncStore, apiId, typeName string) (*appsyncstore.Type, error) {
	if apiId == "" || typeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	t, err := store.GetType(apiId, typeName)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return t, nil
}

// updateTypeCore applies an update to an existing type definition.
func (s *AppSyncService) updateTypeCore(store *appsyncstore.AppSyncStore, in updateTypeInput) (*appsyncstore.Type, error) {
	if in.ApiId == "" || in.TypeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	if in.Format == "" {
		return nil, NewBadRequestException("format is required")
	}

	if !validateTypeFormat(in.Format) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid format: %s. Valid values: SDL, JSON", in.Format))
	}

	t := &appsyncstore.Type{
		ApiId:  in.ApiId,
		Name:   in.TypeName,
		Format: in.Format,
	}

	if in.HasDefinition {
		t.Definition = in.Definition
		t.DefinitionSet = true
	}
	if in.HasDescription {
		t.Description = in.Description
		t.DescriptionSet = true
	}

	updated, err := store.UpdateType(t)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return updated, nil
}

// deleteTypeCore removes a type definition.
func (s *AppSyncService) deleteTypeCore(store *appsyncstore.AppSyncStore, apiId, typeName string) error {
	if apiId == "" || typeName == "" {
		return NewBadRequestException("apiId and typeName are required")
	}

	if err := store.DeleteType(apiId, typeName); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// listTypesCore lists the type definitions of a GraphQL API. The format
// member is required on ListTypes and validated against the SDL and JSON
// enum values; each definition is returned in the requested serialisation,
// converted from the stored one when they differ.
func (s *AppSyncService) listTypesCore(store *appsyncstore.AppSyncStore, apiId, format string, maxResults int, nextToken string) ([]*appsyncstore.Type, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	if format == "" {
		return nil, "", NewBadRequestException("format is required")
	}
	if !validateTypeFormat(format) {
		return nil, "", NewBadRequestException(fmt.Sprintf("Invalid format: %s. Valid values: SDL, JSON", format))
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	types, nextToken, err := store.ListTypes(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return typesInFormat(types, format), nextToken, nil
}

// listTypesByAssociationCore lists the type definitions of the source API
// behind one merged-API association. The format member is required and
// validated against the SDL and JSON enum values.
func (s *AppSyncService) listTypesByAssociationCore(store *appsyncstore.AppSyncStore, mergedApiId, associationId, format string, maxResults int, nextToken string) ([]*appsyncstore.Type, string, error) {
	if mergedApiId == "" {
		return nil, "", NewBadRequestException("mergedApiIdentifier is required")
	}

	if associationId == "" {
		return nil, "", NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	if format == "" {
		return nil, "", NewBadRequestException("format is required")
	}
	if !validateTypeFormat(format) {
		return nil, "", NewBadRequestException(fmt.Sprintf("Invalid format: %s. Valid values: SDL, JSON", format))
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	types, nextToken, err := store.ListTypes(assoc.SourceApiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return typesInFormat(types, format), nextToken, nil
}

// typesInFormat projects stored type records into the requested output
// serialisation, converting definitions whose stored format differs.
func typesInFormat(types []*appsyncstore.Type, format string) []*appsyncstore.Type {
	out := make([]*appsyncstore.Type, len(types))
	for i, t := range types {
		converted := *t
		converted.Definition, converted.Format = typeInRequestedFormat(t, format)
		out[i] = &converted
	}
	return out
}
