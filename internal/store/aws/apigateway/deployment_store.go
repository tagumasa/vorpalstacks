package apigateway

import (
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
	if deployment.ApiSummary == nil {
		deployment.ApiSummary = make(map[string]interface{})
	}
	deployment.Snapshot = CloneSnapshot(api)

	if api.Deployments == nil {
		api.Deployments = make(map[string]*Deployment)
	}
	api.Deployments[deployment.Id] = deployment

	if err := s.Update(api); err != nil {
		return nil, err
	}

	return deployment, nil
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
	return s.Update(api)
}

// ListDeployments returns all deployments for a REST API.
func (s *RestApiStore) ListDeployments(apiId string) ([]*Deployment, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	deployments := make([]*Deployment, 0, len(api.Deployments))
	for _, d := range api.Deployments {
		deployments = append(deployments, d)
	}
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
	return s.Update(api)
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

	if err := s.Update(api); err != nil {
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
	return s.Update(api)
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
	return s.Update(api)
}

// ListStages returns all stages for a REST API.
func (s *RestApiStore) ListStages(apiId string) ([]*Stage, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	stages := make([]*Stage, 0, len(api.Stages))
	for _, st := range api.Stages {
		stages = append(stages, st)
	}
	return stages, nil
}
