package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// HTTP handlers for the continuous deployment policy operations. Each
// handler parses the wire request into a service-layer DTO, delegates to
// the matching Core function, and serialises the store result.

// CreateContinuousDeploymentPolicy creates a new continuous deployment
// policy.
func (s *CloudFrontService) CreateContinuousDeploymentPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := parseContinuousDeploymentPolicyConfig(req)
	if err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.createContinuousDeploymentPolicyCore(store, CreateContinuousDeploymentPolicyInput{Config: config})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ContinuousDeploymentPolicy": formatContinuousDeploymentPolicy(policy),
		"Location":                   policy.ID,
		"ETag":                       policy.ETag,
	}, nil
}

// GetContinuousDeploymentPolicy retrieves a policy by ID.
func (s *CloudFrontService) GetContinuousDeploymentPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getContinuousDeploymentPolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ContinuousDeploymentPolicy": formatContinuousDeploymentPolicy(policy),
		"ETag":                       policy.ETag,
	}, nil
}

// GetContinuousDeploymentPolicyConfig retrieves the configuration of a
// policy.
func (s *CloudFrontService) GetContinuousDeploymentPolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getContinuousDeploymentPolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ContinuousDeploymentPolicyConfig": formatContinuousDeploymentPolicyConfig(policy.ContinuousDeploymentPolicyConfig),
		"ETag":                             policy.ETag,
	}, nil
}

// UpdateContinuousDeploymentPolicy updates an existing policy.
func (s *CloudFrontService) UpdateContinuousDeploymentPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := parseContinuousDeploymentPolicyConfig(req)
	if err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.updateContinuousDeploymentPolicyCore(store, UpdateContinuousDeploymentPolicyInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config:  config,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ContinuousDeploymentPolicy": formatContinuousDeploymentPolicy(policy),
		"ETag":                       policy.ETag,
	}, nil
}

// DeleteContinuousDeploymentPolicy deletes a policy.
func (s *CloudFrontService) DeleteContinuousDeploymentPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteContinuousDeploymentPolicyCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListContinuousDeploymentPolicies lists continuous deployment policies.
func (s *CloudFrontService) ListContinuousDeploymentPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listContinuousDeploymentPoliciesCore(store, ListContinuousDeploymentPoliciesInput{
		Marker:   request.GetStringParam(req.Parameters, "Marker"),
		MaxItems: request.GetIntParam(req.Parameters, "MaxItems"),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Policies))
	for _, policy := range result.Policies {
		items = append(items, map[string]interface{}{
			"ContinuousDeploymentPolicy": formatContinuousDeploymentPolicy(policy),
		})
	}
	policyList := map[string]interface{}{
		"MaxItems": resolveListMaxItems(request.GetIntParam(req.Parameters, "MaxItems")),
		"Quantity": len(items),
		"Items":    protocol.XMLElements{ElementName: "ContinuousDeploymentPolicySummary", Items: items},
	}
	if result.NextMarker != "" {
		policyList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"ContinuousDeploymentPolicyList": policyList}, nil
}

// parseContinuousDeploymentPolicyConfig builds the store configuration
// from the wire payload.
func parseContinuousDeploymentPolicyConfig(req *request.ParsedRequest) (*cloudfrontstore.ContinuousDeploymentPolicyConfig, error) {
	configMap := request.GetMapParam(req.Parameters, "ContinuousDeploymentPolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	dnsMap := request.GetMapParam(configMap, "StagingDistributionDnsNames")
	config := &cloudfrontstore.ContinuousDeploymentPolicyConfig{
		StagingDistributionDnsNames: &cloudfrontstore.StagingDistributionDnsNames{
			Items: parseStringItemList(dnsMap, "Items", "DnsName"),
		},
		Enabled: request.GetBoolParam(configMap, "Enabled"),
	}
	config.StagingDistributionDnsNames.Quantity = len(config.StagingDistributionDnsNames.Items)

	if trafficMap := request.GetMapParam(configMap, "TrafficConfig"); trafficMap != nil {
		traffic := &cloudfrontstore.TrafficConfig{
			Type: request.GetStringParam(trafficMap, "Type"),
		}
		if weightMap := request.GetMapParam(trafficMap, "SingleWeightConfig"); weightMap != nil {
			weight := &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{
				Weight: request.GetFloatParam(weightMap, "Weight"),
			}
			if stickinessMap := request.GetMapParam(weightMap, "SessionStickinessConfig"); stickinessMap != nil {
				weight.SessionStickinessConfig = &cloudfrontstore.SessionStickinessConfig{
					IdleTTL:    int64(request.GetIntParam(stickinessMap, "IdleTTL")),
					MaximumTTL: int64(request.GetIntParam(stickinessMap, "MaximumTTL")),
				}
			}
			traffic.SingleWeightConfig = weight
		}
		if headerMap := request.GetMapParam(trafficMap, "SingleHeaderConfig"); headerMap != nil {
			traffic.SingleHeaderConfig = &cloudfrontstore.ContinuousDeploymentSingleHeaderConfig{
				Header: request.GetStringParam(headerMap, "Header"),
				Value:  request.GetStringParam(headerMap, "Value"),
			}
		}
		config.TrafficConfig = traffic
	}
	return config, nil
}

func formatContinuousDeploymentPolicy(policy *cloudfrontstore.ContinuousDeploymentPolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                               policy.ID,
		"LastModifiedTime":                 policy.LastModifiedTime.Format(time.RFC3339),
		"ContinuousDeploymentPolicyConfig": formatContinuousDeploymentPolicyConfig(policy.ContinuousDeploymentPolicyConfig),
	}
}

func formatContinuousDeploymentPolicyConfig(config *cloudfrontstore.ContinuousDeploymentPolicyConfig) map[string]interface{} {
	if config == nil {
		return map[string]interface{}{}
	}
	m := map[string]interface{}{
		"Enabled": config.Enabled,
	}
	if names := config.StagingDistributionDnsNames; names != nil {
		nameMap := map[string]interface{}{"Quantity": len(names.Items)}
		if len(names.Items) > 0 {
			items := make([]interface{}, len(names.Items))
			for i, item := range names.Items {
				items[i] = item
			}
			nameMap["Items"] = protocol.XMLElements{ElementName: "DnsName", Items: items}
		}
		m["StagingDistributionDnsNames"] = nameMap
	}
	if traffic := config.TrafficConfig; traffic != nil {
		trafficMap := map[string]interface{}{"Type": traffic.Type}
		if weight := traffic.SingleWeightConfig; weight != nil {
			weightMap := map[string]interface{}{"Weight": weight.Weight}
			if stickiness := weight.SessionStickinessConfig; stickiness != nil {
				weightMap["SessionStickinessConfig"] = map[string]interface{}{
					"IdleTTL":    stickiness.IdleTTL,
					"MaximumTTL": stickiness.MaximumTTL,
				}
			}
			trafficMap["SingleWeightConfig"] = weightMap
		}
		if header := traffic.SingleHeaderConfig; header != nil {
			trafficMap["SingleHeaderConfig"] = map[string]interface{}{
				"Header": header.Header,
				"Value":  header.Value,
			}
		}
		m["TrafficConfig"] = trafficMap
	}
	return m
}
