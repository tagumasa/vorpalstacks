package route53

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// CreateHealthCheck creates a new health check in Route 53.
func (s *Route53Service) CreateHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	callerRef := request.GetStringParam(req.Parameters, "CallerReference")
	if callerRef == "" {
		callerRef = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
	}

	healthCheckConfigMap := request.GetMapParam(req.Parameters, "HealthCheckConfig")
	if healthCheckConfigMap == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "HealthCheckConfig is required", 400)
	}

	config := parseHealthCheckConfig(healthCheckConfigMap)

	healthCheck := &route53store.HealthCheck{
		ID:                 generateHealthCheckId(),
		CallerReference:    callerRef,
		HealthCheckConfig:  config,
		HealthCheckVersion: "1",
		Region:             reqCtx.GetRegion(),
		AccountID:          s.accountID,
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HealthChecks().Create(healthCheck); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, NewHealthCheckAlreadyExistsError()
		}
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"HealthCheck": map[string]interface{}{
			"Id":                 healthCheck.ID,
			"CallerReference":    healthCheck.CallerReference,
			"HealthCheckConfig":  s.healthCheckConfigToResponse(healthCheck.HealthCheckConfig),
			"HealthCheckVersion": healthCheck.HealthCheckVersion,
		},
	}, nil
}

// GetHealthCheck returns details about a health check.
func (s *Route53Service) GetHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	healthCheckId, err := extractHealthCheckId(req.Parameters, "HealthCheckId")
	if err != nil {
		return nil, err
	}

	healthCheck, err := s.getHealthCheckById(reqCtx, healthCheckId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HealthCheck": map[string]interface{}{
			"Id":                 healthCheck.ID,
			"CallerReference":    healthCheck.CallerReference,
			"HealthCheckConfig":  s.healthCheckConfigToResponse(healthCheck.HealthCheckConfig),
			"HealthCheckVersion": healthCheck.HealthCheckVersion,
		},
	}, nil
}

// ListHealthChecks returns a list of health checks with pagination support.
func (s *Route53Service) ListHealthChecks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := st.HealthChecks().List(marker, maxItems)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return s.buildHealthChecksListResponse(result.HealthChecks, result.IsTruncated, result.Marker), nil
}

// DeleteHealthCheck deletes a health check by its ID.
func (s *Route53Service) DeleteHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	healthCheckId, err := extractHealthCheckId(req.Parameters, "HealthCheckId")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HealthChecks().Delete(healthCheckId); err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHealthCheckError(healthCheckId)
		}
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UpdateHealthCheck updates an existing health check's configuration.
func (s *Route53Service) UpdateHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	healthCheckId, err := extractHealthCheckId(req.Parameters, "HealthCheckId")
	if err != nil {
		return nil, err
	}

	healthCheck, err := s.getHealthCheckById(reqCtx, healthCheckId)
	if err != nil {
		return nil, err
	}

	healthCheckConfig := request.GetMapParam(req.Parameters, "HealthCheckConfig")
	if healthCheckConfig != nil {
		applyHealthCheckConfigUpdates(healthCheck.HealthCheckConfig, healthCheckConfig)
	} else {
		applyHealthCheckConfigUpdates(healthCheck.HealthCheckConfig, req.Parameters)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HealthChecks().Update(healthCheck); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"HealthCheck": map[string]interface{}{
			"Id":                 healthCheck.ID,
			"CallerReference":    healthCheck.CallerReference,
			"HealthCheckConfig":  s.healthCheckConfigToResponse(healthCheck.HealthCheckConfig),
			"HealthCheckVersion": healthCheck.HealthCheckVersion,
		},
	}, nil
}

func (s *Route53Service) healthCheckConfigToResponse(config *route53store.HealthCheckConfig) map[string]interface{} {
	result := map[string]interface{}{
		"Type": config.Type,
		"Port": config.Port,
	}

	if config.IPAddress != "" {
		result["IPAddress"] = config.IPAddress
	}
	if config.ResourcePath != "" {
		result["ResourcePath"] = config.ResourcePath
	}
	if config.FullyQualifiedDomainName != "" {
		result["FullyQualifiedDomainName"] = config.FullyQualifiedDomainName
	}
	if config.RequestInterval > 0 {
		result["RequestInterval"] = config.RequestInterval
	}
	if config.FailureThreshold > 0 {
		result["FailureThreshold"] = config.FailureThreshold
	}
	if config.MeasureLatency {
		result["MeasureLatency"] = true
	}
	if config.Inverted {
		result["Inverted"] = true
	}
	if config.Disabled {
		result["Disabled"] = true
	}
	if config.EnableSNI {
		result["EnableSNI"] = true
	}

	return result
}

// AssociateVPCWithHostedZone associates a VPC with a private hosted zone.
func (s *Route53Service) AssociateVPCWithHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	vpcMap := request.GetMapParam(req.Parameters, "VPC")
	if vpcMap == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "VPC is required", 400)
	}

	zone, err := s.getHostedZoneById(reqCtx, hostedZoneId)
	if err != nil {
		return nil, err
	}

	if !zone.Private {
		return nil, awserrors.NewAWSError("InvalidInput", "Cannot associate VPC with a public hosted zone", 400)
	}

	vpc := parseVPC(vpcMap)
	for _, existing := range zone.VPCs {
		if existing.VPCID == vpc.VPCID && existing.VPCRegion == vpc.VPCRegion {
			return nil, awserrors.NewAWSError("VPCAlreadyAssociated", "VPC is already associated with the hosted zone", 400)
		}
	}
	zone.VPCs = append(zone.VPCs, vpc)

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + generateChangeId(),
			"Status":      "INSYNC",
			"SubmittedAt": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// DisassociateVPCFromHostedZone disassociates a VPC from a private hosted zone.
func (s *Route53Service) DisassociateVPCFromHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	vpcMap := request.GetMapParam(req.Parameters, "VPC")
	if vpcMap == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "VPC is required", 400)
	}

	zone, err := s.getHostedZoneById(reqCtx, hostedZoneId)
	if err != nil {
		return nil, err
	}

	vpcRegion := request.GetStringParam(vpcMap, "VPCRegion")
	vpcId := request.GetStringParam(vpcMap, "VPCId")

	var newVPCs []*route53store.VPC
	for _, v := range zone.VPCs {
		if v.VPCID != vpcId || v.VPCRegion != vpcRegion {
			newVPCs = append(newVPCs, v)
		}
	}

	if len(newVPCs) == len(zone.VPCs) {
		return nil, awserrors.NewAWSError("InvalidInput", "VPC is not associated with the hosted zone", 400)
	}

	zone.VPCs = newVPCs
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + generateChangeId(),
			"Status":      "INSYNC",
			"SubmittedAt": time.Now().Format(time.RFC3339),
		},
	}, nil
}
