package route53

import (
	"context"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// CreateHealthCheck creates a new health check in Route 53.
func (s *Route53Service) CreateHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	callerRef := request.GetStringParam(req.Parameters, "CallerReference")

	healthCheckConfigMap := request.GetMapParam(req.Parameters, "HealthCheckConfig")
	config := parseHealthCheckConfig(healthCheckConfigMap, s.defaultHCPort)

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	healthCheck, err := s.createHealthCheckCore(st, CreateHealthCheckInput{
		CallerReference: callerRef,
		Config:          config,
		Region:          reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HealthCheck": s.healthCheckToResponse(healthCheck),
	}, nil
}

// GetHealthCheck returns details about a health check.
func (s *Route53Service) GetHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	healthCheckId, err := extractHealthCheckId(req.Parameters, "HealthCheckId")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	healthCheck, err := getHealthCheckCore(st, healthCheckId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HealthCheck": s.healthCheckToResponse(healthCheck),
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

	result, err := listHealthChecksCore(st, ListHealthChecksInput{Marker: marker, MaxItems: maxItems})
	if err != nil {
		return nil, err
	}

	return s.buildHealthChecksListResponse(result.HealthChecks, result.IsTruncated, marker, result.Marker, maxItems), nil
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

	if err := deleteHealthCheckCore(st, DeleteHealthCheckInput{Id: healthCheckId}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateHealthCheck updates an existing health check's configuration.
func (s *Route53Service) UpdateHealthCheck(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	healthCheckId, err := extractHealthCheckId(req.Parameters, "HealthCheckId")
	if err != nil {
		return nil, err
	}

	updates := request.GetMapParam(req.Parameters, "HealthCheckConfig")
	if updates == nil {
		updates = req.Parameters
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	healthCheck, err := updateHealthCheckCore(st, UpdateHealthCheckInput{
		HealthCheckId:      healthCheckId,
		HealthCheckVersion: request.GetIntParam(req.Parameters, "HealthCheckVersion"),
		Updates:            updates,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HealthCheck": s.healthCheckToResponse(healthCheck),
	}, nil
}

func (s *Route53Service) healthCheckConfigToResponse(config *route53store.HealthCheckConfig) map[string]interface{} {
	result := map[string]interface{}{
		"Type": config.Type,
	}

	if config.IPAddress != "" {
		result["IPAddress"] = config.IPAddress
	}
	if config.Port > 0 {
		result["Port"] = config.Port
	}
	if config.ResourcePath != "" {
		result["ResourcePath"] = config.ResourcePath
	}
	if config.FullyQualifiedDomainName != "" {
		result["FullyQualifiedDomainName"] = config.FullyQualifiedDomainName
	}
	if config.SearchString != "" {
		result["SearchString"] = config.SearchString
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
	if config.InsufficientDataHealthStatus != "" {
		result["InsufficientDataHealthStatus"] = config.InsufficientDataHealthStatus
	}
	if config.HealthThreshold > 0 {
		result["HealthThreshold"] = config.HealthThreshold
	}
	if config.RoutingControlArn != "" {
		result["RoutingControlArn"] = config.RoutingControlArn
	}
	if len(config.Regions) > 0 {
		result["Regions"] = protocol.XMLElements{ElementName: "Region", Items: func() []interface{} {
			items := make([]interface{}, len(config.Regions))
			for i, r := range config.Regions {
				items[i] = r
			}
			return items
		}()}
	}
	if config.AlarmIdentifier != nil {
		result["AlarmIdentifier"] = map[string]interface{}{
			"Region": config.AlarmIdentifier.Region,
			"Name":   config.AlarmIdentifier.Name,
		}
	}
	if len(config.ChildHealthChecks) > 0 {
		result["ChildHealthChecks"] = protocol.XMLElements{ElementName: "ChildHealthCheck", Items: func() []interface{} {
			items := make([]interface{}, len(config.ChildHealthChecks))
			for i, c := range config.ChildHealthChecks {
				items[i] = c
			}
			return items
		}()}
	}

	return result
}

// AssociateVPCWithHostedZone associates a VPC with a private hosted zone.
func (s *Route53Service) AssociateVPCWithHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.associateVPCWithHostedZoneCore(ctx, st, AssociateVPCWithHostedZoneInput{
		HostedZoneId: hostedZoneId,
		VPC:          parseVPC(request.GetMapParam(req.Parameters, "VPC")),
		Comment:      request.GetStringParam(req.Parameters, "Comment"),
		Region:       reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + result.ChangeId,
			"Status":      result.Status,
			"SubmittedAt": result.SubmittedAt.Format(time.RFC3339),
			"Comment":     result.Comment,
		},
	}, nil
}

// DisassociateVPCFromHostedZone disassociates a VPC from a private hosted zone.
func (s *Route53Service) DisassociateVPCFromHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := disassociateVPCFromHostedZoneCore(st, DisassociateVPCFromHostedZoneInput{
		HostedZoneId: hostedZoneId,
		VPC:          parseVPC(request.GetMapParam(req.Parameters, "VPC")),
		Comment:      request.GetStringParam(req.Parameters, "Comment"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + result.ChangeId,
			"Status":      result.Status,
			"SubmittedAt": result.SubmittedAt.Format(time.RFC3339),
			"Comment":     result.Comment,
		},
	}, nil
}
