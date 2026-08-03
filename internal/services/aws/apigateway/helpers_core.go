package apigateway

import (
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/aws/types"
)

// tagsFromMap converts a flat tag map into the []types.Tag slice used by
// the store layer. Returns nil when the map is empty so that the JSON
// serialiser skips the field entirely.
func tagsFromMap(m map[string]string) []types.Tag {
	if len(m) == 0 {
		return nil
	}
	return tags.MapToTags(m)
}

// canarySettingsFromInput copies the transport-agnostic canary settings
// into the store-layer struct. deploymentIdOverride, when non-empty,
// forces the embedded deploymentId (used by CreateDeployment to attach
// the new deployment to its canary); pass "" to honour the input value.
func canarySettingsFromInput(in *CanarySettingsInput, deploymentIdOverride string) *apigateway.CanarySettings {
	deploymentId := in.DeploymentId
	if deploymentIdOverride != "" {
		deploymentId = deploymentIdOverride
	}
	return &apigateway.CanarySettings{
		PercentTraffic:         in.PercentTraffic,
		DeploymentId:           deploymentId,
		StageVariableOverrides: in.StageVariableOverrides,
		UseStageCache:          in.UseStageCache,
	}
}
