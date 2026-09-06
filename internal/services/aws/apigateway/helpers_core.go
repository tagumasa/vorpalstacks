package apigateway

import (
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
)

// tagsFromMap converts a flat tag map into the []tags.Tag slice used by
// the store layer. Returns nil when the map is empty so that the JSON
// serialiser skips the field entirely.
func tagsFromMap(m map[string]string) []tags.Tag {
	if len(m) == 0 {
		return nil
	}
	return tags.MapToTags(m)
}

// canarySettingsFromInput copies the transport-agnostic canary settings
// into the store-layer struct, validating the stage variable overrides
// against the documented stage variable constraints. deploymentIdOverride,
// when non-empty, forces the embedded deploymentId (used by
// CreateDeployment to attach the new deployment to its canary); pass ""
// to honour the input value.
func canarySettingsFromInput(in *CanarySettingsInput, deploymentIdOverride string) (*apigateway.CanarySettings, error) {
	if err := validateStageVariables(in.StageVariableOverrides); err != nil {
		return nil, err
	}
	deploymentId := in.DeploymentId
	if deploymentIdOverride != "" {
		deploymentId = deploymentIdOverride
	}
	return &apigateway.CanarySettings{
		PercentTraffic:         in.PercentTraffic,
		DeploymentId:           deploymentId,
		StageVariableOverrides: in.StageVariableOverrides,
		UseStageCache:          in.UseStageCache,
	}, nil
}
