package apigateway

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// maxRestApiBlobSize is the upper bound on the serialised size of a single
// RestApi JSON blob. The RestApi struct embeds all sub-collections
// (resources, deployments, stages, models, authorizers, validators) and is
// loaded in its entirety by Get and every sub-resource List operation.
// This guard prevents unbounded growth from causing excessive memory
// consumption during (de)serialisation.
const maxRestApiBlobSize = 5 * 1024 * 1024 // 5 MiB

// RestApiStore provides storage operations for REST APIs.
type RestApiStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	accountId  string
	region     string
	mu         sync.Mutex
}

func restApiBucketName(region string) string {
	return "apigateway-restapis-" + region
}

// NewRestApiStore creates a new RestApiStore instance.
func NewRestApiStore(store storage.BasicStorage, accountId, region string) *RestApiStore {
	bucket := store.Bucket(restApiBucketName(region))
	return &RestApiStore{
		BaseStore:  common.NewBaseStore(bucket, "apigateway-restapis"),
		arnBuilder: NewARNBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

// ensureRestApiMaps initialises nil maps on a RestApi so that serialised
// output is always a valid empty map rather than null.
func ensureRestApiMaps(api *RestApi) {
	if api.Resources == nil {
		api.Resources = make(map[string]*Resource)
	}
	if api.Deployments == nil {
		api.Deployments = make(map[string]*Deployment)
	}
	if api.Stages == nil {
		api.Stages = make(map[string]*Stage)
	}
	if api.RequestValidators == nil {
		api.RequestValidators = make(map[string]*RequestValidator)
	}
	if api.Models == nil {
		api.Models = make(map[string]*Model)
	}
	if api.Authorizers == nil {
		api.Authorizers = make(map[string]*Authorizer)
	}
}

// updateLocked persists a RestApi under the caller's lock.
// Validates existence and ensures maps before writing.
func (s *RestApiStore) updateLocked(api *RestApi) error {
	if !s.Exists(api.Id) {
		return ErrRestApiNotFound
	}
	ensureRestApiMaps(api)

	// Reject blobs that exceed the size guard. The entire RestApi — including
	// all sub-collections — is serialised as one JSON value and loaded in full
	// by Get and List operations. An oversized blob risks excessive memory
	// consumption and slow (de)serialisation on every access.
	data, err := json.Marshal(api)
	if err != nil {
		return common.NewStoreErrorWithKey("apigateway", "marshal", api.Id, err)
	}
	if len(data) > maxRestApiBlobSize {
		return common.NewStoreErrorWithKey("apigateway", "update", api.Id,
			fmt.Errorf("rest API %s exceeds maximum storage size (%d bytes); consider splitting resources across multiple APIs", api.Id, maxRestApiBlobSize))
	}
	return s.BaseStore.PutRaw(api.Id, data)
}

// Create creates a new REST API.
func (s *RestApiStore) Create(api *RestApi) (*RestApi, error) {
	if api.Id == "" {
		api.Id = s.arnBuilder.GenerateApiId()
	}
	if api.Name == "" {
		return nil, ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(api.Id) {
		return nil, ErrRestApiAlreadyExists
	}

	now := time.Now().UTC()
	api.CreatedDate = now
	ensureRestApiMaps(api)

	rootResource := &Resource{
		Id:              s.arnBuilder.GenerateResourceId(),
		RestApiId:       api.Id,
		Path:            "/",
		ResourceMethods: make(map[string]*Method),
	}
	api.Resources[rootResource.Id] = rootResource

	if err := s.Put(api.Id, api); err != nil {
		return nil, err
	}

	return api, nil
}

// Get retrieves a REST API by its ID.
func (s *RestApiStore) Get(apiId string) (*RestApi, error) {
	var api RestApi
	if err := s.BaseStore.Get(apiId, &api); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrRestApiNotFound
		}
		return nil, err
	}
	return &api, nil
}

// Update updates an existing REST API. Acquires the store lock so that
// callers from the service layer get an atomic write.
func (s *RestApiStore) Update(api *RestApi) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateLocked(api)
}

// CloneFromSource deep-copies resources, models, authorizers and request
// validators from the source API into target, restamping every RestApiId
// (including those nested under Method and MethodIntegration) to target.Id.
// The caller is responsible for having already persisted target via Create.
// Returns ErrRestApiNotFound if sourceApiId does not exist.
func (s *RestApiStore) CloneFromSource(target *RestApi, sourceApiId string) error {
	source, err := s.Get(sourceApiId)
	if err != nil {
		return err
	}

	resources, err := deepCopyResources(source.Resources)
	if err != nil {
		return fmt.Errorf("clone resources: %w", err)
	}
	models, err := deepCopyModels(source.Models)
	if err != nil {
		return fmt.Errorf("clone models: %w", err)
	}
	authorizers, err := deepCopyAuthorizers(source.Authorizers)
	if err != nil {
		return fmt.Errorf("clone authorizers: %w", err)
	}
	validators, err := deepCopyRequestValidators(source.RequestValidators)
	if err != nil {
		return fmt.Errorf("clone validators: %w", err)
	}

	// Restamp RestApiId on all cloned child elements. AWS assigns the new
	// API id to every nested resource so that GetResources, GetModels etc.
	// associate them with the clone rather than the source.
	for _, r := range resources {
		r.RestApiId = target.Id
		for _, m := range r.ResourceMethods {
			m.RestApiId = target.Id
			if m.MethodIntegration != nil {
				m.MethodIntegration.RestApiId = target.Id
			}
		}
	}
	for _, m := range models {
		m.RestApiId = target.Id
	}
	for _, a := range authorizers {
		a.RestApiId = target.Id
	}
	for _, v := range validators {
		v.RestApiId = target.Id
	}

	target.Resources = resources
	target.Models = models
	target.Authorizers = authorizers
	target.RequestValidators = validators

	return s.Update(target)
}

// Delete deletes a REST API.
func (s *RestApiStore) Delete(apiId string) error {
	if !s.Exists(apiId) {
		return ErrRestApiNotFound
	}
	return s.BaseStore.Delete(apiId)
}

// List returns all REST APIs.
func (s *RestApiStore) List(opts common.ListOptions) (*common.ListResult[RestApi], error) {
	return common.List[RestApi](s.BaseStore, opts, nil)
}

// TagResource adds tags to a REST API.
func (s *RestApiStore) Tag(apiId string, inputTags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if api.Tags == nil {
		api.Tags = []tags.Tag{}
	}

	api.Tags = tags.Apply(api.Tags, tags.MapToTags(inputTags))

	return s.updateLocked(api)
}

// UntagResource removes tags from a REST API.
func (s *RestApiStore) Untag(apiId string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	api.Tags = tags.RemoveByTagKeys(api.Tags, tagKeys)

	return s.updateLocked(api)
}

// TagStage adds or updates tags on a specific stage of a REST API.
func (s *RestApiStore) TagStage(apiId, stageName string, inputTags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	stage, ok := api.Stages[stageName]
	if !ok {
		return ErrStageNotFound
	}

	if stage.Tags == nil {
		stage.Tags = []tags.Tag{}
	}

	stage.Tags = tags.Apply(stage.Tags, tags.MapToTags(inputTags))

	return s.updateLocked(api)
}

// UntagStage removes the specified tag keys from a stage of a REST API.
func (s *RestApiStore) UntagStage(apiId, stageName string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	stage, ok := api.Stages[stageName]
	if !ok {
		return ErrStageNotFound
	}

	stage.Tags = tags.RemoveByTagKeys(stage.Tags, tagKeys)

	return s.updateLocked(api)
}

// GetStageTags returns the tags associated with a specific stage of a REST API.
func (s *RestApiStore) GetStageTags(apiId, stageName string) ([]tags.Tag, error) {
	stage, err := s.GetStage(apiId, stageName)
	if err != nil {
		return nil, err
	}
	return stage.Tags, nil
}

// GetResourceTags returns the tags associated with a REST API resource.
func (s *RestApiStore) GetResourceTags(apiId string) ([]tags.Tag, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}
	return api.Tags, nil
}

// CreateRequestValidator creates a new request validator for a REST API.
func (s *RestApiStore) CreateRequestValidator(apiId string, validator *RequestValidator) (*RequestValidator, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	if validator.Id == "" {
		validator.Id = s.arnBuilder.GenerateValidatorId()
	}
	validator.RestApiId = apiId

	if api.RequestValidators == nil {
		api.RequestValidators = make(map[string]*RequestValidator)
	}
	api.RequestValidators[validator.Id] = validator

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return validator, nil
}

// GetRequestValidator retrieves a request validator by its ID.
func (s *RestApiStore) GetRequestValidator(apiId, validatorId string) (*RequestValidator, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	validator, ok := api.RequestValidators[validatorId]
	if !ok {
		return nil, ErrRequestValidatorNotFound
	}
	return validator, nil
}

// UpdateRequestValidator updates a request validator.
func (s *RestApiStore) UpdateRequestValidator(apiId string, validator *RequestValidator) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.RequestValidators[validator.Id]; !ok {
		return ErrRequestValidatorNotFound
	}

	api.RequestValidators[validator.Id] = validator
	return s.updateLocked(api)
}

// DeleteRequestValidator deletes a request validator.
func (s *RestApiStore) DeleteRequestValidator(apiId, validatorId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.RequestValidators[validatorId]; !ok {
		return ErrRequestValidatorNotFound
	}

	delete(api.RequestValidators, validatorId)
	return s.updateLocked(api)
}

// ListRequestValidators returns all request validators for a REST API.
func (s *RestApiStore) ListRequestValidators(apiId string) ([]*RequestValidator, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	validators := make([]*RequestValidator, 0, len(api.RequestValidators))
	for _, v := range api.RequestValidators {
		validators = append(validators, v)
	}
	return validators, nil
}

// CreateModel creates a new model for a REST API.
func (s *RestApiStore) CreateModel(apiId string, model *Model) (*Model, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	if api.Models == nil {
		api.Models = make(map[string]*Model)
	}

	if _, ok := api.Models[model.Name]; ok {
		return nil, ErrModelAlreadyExists
	}

	if model.Id == "" {
		model.Id = s.arnBuilder.GenerateModelId()
	}
	model.RestApiId = apiId

	api.Models[model.Name] = model

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return model, nil
}

// GetModel retrieves a model by its name.
func (s *RestApiStore) GetModel(apiId, modelName string) (*Model, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	model, ok := api.Models[modelName]
	if !ok {
		return nil, ErrModelNotFound
	}
	return model, nil
}

// DeleteModel deletes a model.
func (s *RestApiStore) DeleteModel(apiId, modelName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Models[modelName]; !ok {
		return ErrModelNotFound
	}

	delete(api.Models, modelName)
	return s.updateLocked(api)
}

// ListModels returns all models for a REST API.
func (s *RestApiStore) ListModels(apiId string) ([]*Model, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	models := make([]*Model, 0, len(api.Models))
	for _, m := range api.Models {
		models = append(models, m)
	}
	return models, nil
}

// UpdateModel updates a model.
func (s *RestApiStore) UpdateModel(apiId string, model *Model) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Models[model.Name]; !ok {
		return ErrModelNotFound
	}

	model.RestApiId = apiId
	api.Models[model.Name] = model
	return s.updateLocked(api)
}

// CreateAuthorizer creates a new authorizer for a REST API.
func (s *RestApiStore) CreateAuthorizer(apiId string, authorizer *Authorizer) (*Authorizer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	if authorizer.Id == "" {
		authorizer.Id = s.arnBuilder.GenerateAuthorizerId()
	}
	authorizer.RestApiId = apiId

	if api.Authorizers == nil {
		api.Authorizers = make(map[string]*Authorizer)
	}
	api.Authorizers[authorizer.Id] = authorizer

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return authorizer, nil
}

// GetAuthorizer retrieves an authorizer by its ID.
func (s *RestApiStore) GetAuthorizer(apiId, authorizerId string) (*Authorizer, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	authorizer, ok := api.Authorizers[authorizerId]
	if !ok {
		return nil, ErrAuthorizerNotFound
	}
	return authorizer, nil
}

// UpdateAuthorizer updates an authorizer.
func (s *RestApiStore) UpdateAuthorizer(apiId string, authorizer *Authorizer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Authorizers[authorizer.Id]; !ok {
		return ErrAuthorizerNotFound
	}

	authorizer.RestApiId = apiId
	api.Authorizers[authorizer.Id] = authorizer
	return s.updateLocked(api)
}

// DeleteAuthorizer deletes an authorizer.
func (s *RestApiStore) DeleteAuthorizer(apiId, authorizerId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Authorizers[authorizerId]; !ok {
		return ErrAuthorizerNotFound
	}

	delete(api.Authorizers, authorizerId)
	return s.updateLocked(api)
}

// ListAuthorizers returns all authorizers for a REST API.
func (s *RestApiStore) ListAuthorizers(apiId string) ([]*Authorizer, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	authorizers := make([]*Authorizer, 0, len(api.Authorizers))
	for _, a := range api.Authorizers {
		authorizers = append(authorizers, a)
	}
	return authorizers, nil
}
