package apigateway

import (
	"fmt"
	"sort"
	"time"
)

// CreateDeployment creates a new deployment for a REST API.
func (s *RestApiStore) CreateDeployment(apiId string, deployment *Deployment) (*Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	if deployment.Id == "" {
		deployment.Id = s.arnBuilder.GenerateDeploymentId()
	}
	deployment.RestApiId = apiId
	deployment.CreatedDate = time.Now().UTC()
	snapshot, err := CloneSnapshot(api)
	if err != nil {
		return nil, fmt.Errorf("clone snapshot: %w", err)
	}
	deployment.Snapshot = snapshot
	deployment.ApiSummary = apiSummaryFromSnapshot(snapshot)

	if api.Deployments == nil {
		api.Deployments = make(map[string]*Deployment)
	}
	api.Deployments[deployment.Id] = deployment

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return deployment, nil
}

// apiSummaryFromSnapshot renders the PathToMapOfMethodSnapshot view of a
// deployment snapshot: each resource path maps to its methods, each method
// to the authorizationType/apiKeyRequired snapshot a deployment reader
// inspects. The summary is derived from the same clone the deployment
// persists, so the two can never disagree.
func apiSummaryFromSnapshot(snapshot *DeploymentSnapshot) map[string]interface{} {
	summary := make(map[string]interface{}, len(snapshot.Resources))
	for _, resource := range snapshot.Resources {
		if len(resource.ResourceMethods) == 0 {
			continue
		}
		methods := make(map[string]interface{}, len(resource.ResourceMethods))
		for verb, method := range resource.ResourceMethods {
			methods[verb] = map[string]interface{}{
				"authorizationType": method.AuthorizationType,
				"apiKeyRequired":    method.ApiKeyRequired,
			}
		}
		summary[resource.Path] = methods
	}
	return summary
}

// GetDeployment retrieves a deployment by API ID and deployment ID.
func (s *RestApiStore) GetDeployment(apiId, deploymentId string) (*Deployment, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	deployment, ok := api.Deployments[deploymentId]
	if !ok {
		return nil, ErrDeploymentNotFound
	}
	return deployment, nil
}

// DeleteDeployment deletes a deployment from a REST API.
func (s *RestApiStore) DeleteDeployment(apiId, deploymentId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Deployments[deploymentId]; !ok {
		return ErrDeploymentNotFound
	}

	for _, stage := range api.Stages {
		if stage.DeploymentId == deploymentId {
			return ErrDeploymentInUse
		}
	}

	delete(api.Deployments, deploymentId)
	return s.updateLocked(api)
}

// ListDeployments returns all deployments for a REST API in a deterministic
// total order (newest first, id as the tiebreak): a page token is only
// meaningful when two listing calls order equal-timestamp deployments the
// same way.
func (s *RestApiStore) ListDeployments(apiId string) ([]*Deployment, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	deployments := make([]*Deployment, 0, len(api.Deployments))
	for _, d := range api.Deployments {
		deployments = append(deployments, d)
	}
	sort.Slice(deployments, func(i, j int) bool {
		if !deployments[i].CreatedDate.Equal(deployments[j].CreatedDate) {
			return deployments[i].CreatedDate.After(deployments[j].CreatedDate)
		}
		return deployments[i].Id > deployments[j].Id
	})
	return deployments, nil
}

// UpdateDeployment updates an existing deployment for a REST API.
func (s *RestApiStore) UpdateDeployment(apiId string, deployment *Deployment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Deployments[deployment.Id]; !ok {
		return ErrDeploymentNotFound
	}

	deployment.RestApiId = apiId
	api.Deployments[deployment.Id] = deployment
	return s.updateLocked(api)
}

// CreateStage creates a new stage for a REST API.
func (s *RestApiStore) CreateStage(apiId string, stage *Stage) (*Stage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	if api.Stages == nil {
		api.Stages = make(map[string]*Stage)
	}

	if _, ok := api.Stages[stage.StageName]; ok {
		return nil, ErrStageAlreadyExists
	}

	if stage.DeploymentId != "" {
		if _, ok := api.Deployments[stage.DeploymentId]; !ok {
			return nil, ErrDeploymentNotFound
		}
	}

	stage.RestApiId = apiId
	now := time.Now().UTC()
	stage.CreatedDate = now
	stage.LastUpdatedDate = now
	if stage.Variables == nil {
		stage.Variables = make(map[string]string)
	}
	if stage.MethodSettings == nil {
		stage.MethodSettings = make(map[string]*MethodSetting)
	}

	api.Stages[stage.StageName] = stage

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return stage, nil
}

// GetStage retrieves a stage by API ID and stage name.
func (s *RestApiStore) GetStage(apiId, stageName string) (*Stage, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	stage, ok := api.Stages[stageName]
	if !ok {
		return nil, ErrStageNotFound
	}
	return stage, nil
}

// UpdateStage updates an existing stage for a REST API.
func (s *RestApiStore) UpdateStage(apiId string, stage *Stage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Stages[stage.StageName]; !ok {
		return ErrStageNotFound
	}

	stage.LastUpdatedDate = time.Now().UTC()
	api.Stages[stage.StageName] = stage
	return s.updateLocked(api)
}

// DeleteStage deletes a stage from a REST API.
func (s *RestApiStore) DeleteStage(apiId, stageName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Stages[stageName]; !ok {
		return ErrStageNotFound
	}

	delete(api.Stages, stageName)
	return s.updateLocked(api)
}

// ListStages returns all stages for a REST API ordered by stage name — a
// deterministic total order for stable pagination.
func (s *RestApiStore) ListStages(apiId string) ([]*Stage, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	stages := make([]*Stage, 0, len(api.Stages))
	for _, st := range api.Stages {
		stages = append(stages, st)
	}
	sort.Slice(stages, func(i, j int) bool { return stages[i].StageName < stages[j].StageName })
	return stages, nil
}
