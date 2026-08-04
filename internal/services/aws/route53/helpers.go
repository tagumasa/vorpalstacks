package route53

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	route53store "vorpalstacks/internal/store/aws/route53"
)

func generateHostedZoneId() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 13)
	b[0] = 'Z'
	for i := 1; i < 13; i++ {
		n, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}
	return string(b)
}

func generateHealthCheckId() string {
	const letters = "0123456789abcdef"
	b := make([]byte, 10)
	for i := range b {
		n, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}
	return string(b)
}

var changeCounter uint64 = 0

func generateChangeId() string {
	counter := atomic.AddUint64(&changeCounter, 1)
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%d", counter, time.Now().UnixNano())))
	return fmt.Sprintf("%X", hash)[:8]
}

func extractHostedZoneId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", awserrors.NewAWSError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
	}
	return strings.TrimPrefix(id, "/hostedzone/"), nil
}

func extractHealthCheckId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", awserrors.NewAWSError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
	}
	return strings.TrimPrefix(id, "/healthcheck/"), nil
}

func extractChangeId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", awserrors.NewAWSError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
	}
	return strings.TrimPrefix(id, "/change/"), nil
}

func (s *Route53Service) getHostedZoneById(reqCtx *request.RequestContext, id string) (*route53store.HostedZone, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	zone, err := st.HostedZones().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHostedZoneError(id)
		}
		return nil, mapStoreError(err)
	}
	return zone, nil
}

func (s *Route53Service) getHealthCheckById(reqCtx *request.RequestContext, id string) (*route53store.HealthCheck, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	hc, err := st.HealthChecks().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHealthCheckError(id)
		}
		return nil, mapStoreError(err)
	}
	return hc, nil
}

func (s *Route53Service) getChangeById(reqCtx *request.RequestContext, id string) (*route53store.ChangeInfo, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	change, err := st.Changes().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchChangeError(id)
		}
		return nil, mapStoreError(err)
	}
	return change, nil
}

func parseHealthCheckConfig(configMap map[string]interface{}, defaultPort int64) *route53store.HealthCheckConfig {
	if configMap == nil {
		return nil
	}

	config := &route53store.HealthCheckConfig{
		Type: request.GetStringParam(configMap, "Type"),
	}

	hcType := request.GetStringParam(configMap, "Type")
	if hcType != "CALCULATED" && hcType != "CLOUDWATCH_METRIC" {
		config.Port = defaultPort
	}

	if v := request.GetIntParam(configMap, "Port"); v > 0 {
		config.Port = int64(v)
	}
	if v, ok := configMap["IPAddress"].(string); ok {
		config.IPAddress = v
	}
	if v, ok := configMap["ResourcePath"].(string); ok {
		config.ResourcePath = v
	}
	if v, ok := configMap["FullyQualifiedDomainName"].(string); ok {
		config.FullyQualifiedDomainName = v
	}
	if v, ok := configMap["SearchString"].(string); ok {
		config.SearchString = v
	}
	if v := request.GetIntParam(configMap, "RequestInterval"); v > 0 {
		config.RequestInterval = int64(v)
	}
	if v := request.GetIntParam(configMap, "FailureThreshold"); v > 0 {
		config.FailureThreshold = int64(v)
	}
	config.MeasureLatency = request.GetBoolParam(configMap, "MeasureLatency")
	config.Inverted = request.GetBoolParam(configMap, "Inverted")
	config.Disabled = request.GetBoolParam(configMap, "Disabled")
	config.EnableSNI = request.GetBoolParam(configMap, "EnableSNI")

	if v, ok := configMap["InsufficientDataHealthStatus"].(string); ok {
		config.InsufficientDataHealthStatus = v
	}

	if v := request.GetIntParam(configMap, "HealthThreshold"); v > 0 {
		config.HealthThreshold = int64(v)
	}

	if v, ok := configMap["RoutingControlArn"].(string); ok {
		config.RoutingControlArn = v
	}

	if regionsRaw, ok := configMap["Regions"].([]interface{}); ok {
		for _, r := range regionsRaw {
			if region, ok := r.(string); ok {
				config.Regions = append(config.Regions, region)
			}
		}
	}

	if alarmMap, ok := configMap["AlarmIdentifier"].(map[string]interface{}); ok {
		config.AlarmIdentifier = &route53store.AlarmIdentifier{
			Region: request.GetStringParam(alarmMap, "Region"),
			Name:   request.GetStringParam(alarmMap, "Name"),
		}
	}

	if childrenRaw, ok := configMap["ChildHealthChecks"].([]interface{}); ok {
		for _, c := range childrenRaw {
			if child, ok := c.(string); ok {
				config.ChildHealthChecks = append(config.ChildHealthChecks, child)
			}
		}
	}

	return config
}

// validateHealthCheckConfig validates the parsed HealthCheckConfig against
// AWS constraints: Type enum, numeric ranges, and string length limits.
func validateHealthCheckConfig(config *route53store.HealthCheckConfig) error {
	if config == nil {
		return awserrors.NewAWSError("InvalidInput", "HealthCheckConfig is required", 400)
	}

	// Type must be one of the 8 valid values.
	validTypes := map[string]bool{
		"HTTP": true, "HTTPS": true, "HTTP_STR_MATCH": true,
		"HTTPS_STR_MATCH": true, "TCP": true,
		"CALCULATED": true, "CLOUDWATCH_METRIC": true,
		"RECOVERY_CONTROL": true,
	}
	if !validTypes[config.Type] {
		return awserrors.NewAWSError("InvalidInput",
			fmt.Sprintf("Invalid or missing health check type: %q. Must be one of: HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP, CALCULATED, CLOUDWATCH_METRIC, RECOVERY_CONTROL", config.Type), 400)
	}

	// Numeric range validation (only validates when the field is set).
	if config.Port > 65535 {
		return awserrors.NewAWSError("InvalidInput", "Port must be between 1 and 65535", 400)
	}
	if config.FailureThreshold > 10 {
		return awserrors.NewAWSError("InvalidInput", "FailureThreshold must be between 1 and 10", 400)
	}
	if config.RequestInterval > 0 && (config.RequestInterval < 10 || config.RequestInterval > 30) {
		return awserrors.NewAWSError("InvalidInput", "RequestInterval must be between 10 and 30", 400)
	}
	if config.HealthThreshold > 256 {
		return awserrors.NewAWSError("InvalidInput", "HealthThreshold must be between 0 and 256", 400)
	}

	// String length validation (AWS docs constraints).
	if len(config.ResourcePath) > 255 {
		return awserrors.NewAWSError("InvalidInput", "ResourcePath must not exceed 255 characters", 400)
	}
	if len(config.SearchString) > 255 {
		return awserrors.NewAWSError("InvalidInput", "SearchString must not exceed 255 characters", 400)
	}
	if len(config.FullyQualifiedDomainName) > 255 {
		return awserrors.NewAWSError("InvalidInput", "FullyQualifiedDomainName must not exceed 255 characters", 400)
	}
	if len(config.IPAddress) > 45 {
		return awserrors.NewAWSError("InvalidInput", "IPAddress must not exceed 45 characters", 400)
	}
	if len(config.RoutingControlArn) > 255 {
		return awserrors.NewAWSError("InvalidInput", "RoutingControlArn must not exceed 255 characters", 400)
	}

	return nil
}

func applyHealthCheckConfigUpdates(config *route53store.HealthCheckConfig, updates map[string]interface{}) {
	if config == nil || updates == nil {
		return
	}

	if v, ok := updates["IPAddress"].(string); ok {
		config.IPAddress = v
	}
	if _, ok := updates["Port"]; ok {
		config.Port = int64(request.GetIntParam(updates, "Port"))
	}
	if v, ok := updates["ResourcePath"].(string); ok {
		config.ResourcePath = v
	}
	if v, ok := updates["FullyQualifiedDomainName"].(string); ok {
		config.FullyQualifiedDomainName = v
	}
	if v, ok := updates["SearchString"].(string); ok {
		config.SearchString = v
	}
	if _, ok := updates["RequestInterval"]; ok {
		config.RequestInterval = int64(request.GetIntParam(updates, "RequestInterval"))
	}
	if _, ok := updates["FailureThreshold"]; ok {
		config.FailureThreshold = int64(request.GetIntParam(updates, "FailureThreshold"))
	}
	if updates["MeasureLatency"] != nil {
		config.MeasureLatency = request.GetBoolParam(updates, "MeasureLatency")
	}
	if updates["Inverted"] != nil {
		config.Inverted = request.GetBoolParam(updates, "Inverted")
	}
	if updates["Disabled"] != nil {
		config.Disabled = request.GetBoolParam(updates, "Disabled")
	}
	if updates["EnableSNI"] != nil {
		config.EnableSNI = request.GetBoolParam(updates, "EnableSNI")
	}
	if v, ok := updates["InsufficientDataHealthStatus"].(string); ok {
		config.InsufficientDataHealthStatus = v
	}
	if _, ok := updates["HealthThreshold"]; ok {
		config.HealthThreshold = int64(request.GetIntParam(updates, "HealthThreshold"))
	}
	if v, ok := updates["RoutingControlArn"].(string); ok {
		config.RoutingControlArn = v
	}
	if regionsRaw, ok := updates["Regions"].([]interface{}); ok {
		config.Regions = nil
		for _, r := range regionsRaw {
			if region, ok := r.(string); ok {
				config.Regions = append(config.Regions, region)
			}
		}
	}
	if alarmMap, ok := updates["AlarmIdentifier"].(map[string]interface{}); ok {
		config.AlarmIdentifier = &route53store.AlarmIdentifier{
			Region: request.GetStringParam(alarmMap, "Region"),
			Name:   request.GetStringParam(alarmMap, "Name"),
		}
	}
	if childrenRaw, ok := updates["ChildHealthChecks"].([]interface{}); ok {
		config.ChildHealthChecks = nil
		for _, c := range childrenRaw {
			if child, ok := c.(string); ok {
				config.ChildHealthChecks = append(config.ChildHealthChecks, child)
			}
		}
	}

	// Process ResetElements — reset specified fields to their default
	// values. Applied AFTER field updates so that resets take precedence.
	if resetsRaw, ok := updates["ResetElements"].([]interface{}); ok {
		for _, r := range resetsRaw {
			field, ok := r.(string)
			if !ok {
				continue
			}
			switch field {
			case "FullyQualifiedDomainName":
				config.FullyQualifiedDomainName = ""
			case "Regions":
				config.Regions = nil
			case "ResourcePath":
				config.ResourcePath = ""
			case "ChildHealthChecks":
				config.ChildHealthChecks = nil
			case "IPAddress":
				config.IPAddress = ""
			case "Port":
				config.Port = 0
			case "SearchString":
				config.SearchString = ""
			case "Type":
				config.Type = ""
			}
		}
	}
}

func generateNameServers(count int) []string {
	servers := make([]string, count)
	for i := 0; i < count; i++ {
		servers[i] = fmt.Sprintf("ns-%d.vorpalstacks.local.", i+1)
	}
	return servers
}

func parseVPC(vpcMap map[string]interface{}) *route53store.VPC {
	if vpcMap == nil {
		return nil
	}
	return &route53store.VPC{
		VPCRegion: request.GetStringParam(vpcMap, "VPCRegion"),
		VPCID:     request.GetStringParam(vpcMap, "VPCId"),
	}
}

func (s *Route53Service) buildHostedZonesListResponse(zones []*route53store.HostedZone, isTruncated bool, requestMarker, nextMarker string, maxItems int) map[string]interface{} {
	result := make([]interface{}, len(zones))
	for i, z := range zones {
		result[i] = s.hostedZoneToResponse(z)
	}
	response := map[string]interface{}{
		"HostedZones": protocol.XMLElements{ElementName: "HostedZone", Items: result},
		"IsTruncated": isTruncated,
		"Marker":      requestMarker,
		"MaxItems":    maxItems,
	}
	if isTruncated && nextMarker != "" {
		response["NextMarker"] = nextMarker
	}
	return response
}

func (s *Route53Service) buildHealthChecksListResponse(healthChecks []*route53store.HealthCheck, isTruncated bool, requestMarker, nextMarker string, maxItems int) map[string]interface{} {
	result := make([]interface{}, len(healthChecks))
	for i, hc := range healthChecks {
		result[i] = s.healthCheckToResponse(hc)
	}
	response := map[string]interface{}{
		"HealthChecks": protocol.XMLElements{ElementName: "HealthCheck", Items: result},
		"IsTruncated":  isTruncated,
		"Marker":       requestMarker,
		"MaxItems":     maxItems,
	}
	if isTruncated && nextMarker != "" {
		response["NextMarker"] = nextMarker
	}
	return response
}

func buildDelegationSetResponse(nameServers []string, delegationSetID string) delegationSetResponse {
	nsItems := make([]interface{}, len(nameServers))
	for i, ns := range nameServers {
		nsItems[i] = ns
	}
	dsResp := delegationSetResponse{
		NameServers: protocol.XMLElements{
			ElementName: "NameServer",
			Items:       nsItems,
		},
	}
	if delegationSetID != "" {
		dsResp.ID = "/delegationset/" + delegationSetID
	}
	return dsResp
}

func (s *Route53Service) healthCheckToResponse(hc *route53store.HealthCheck) map[string]interface{} {
	result := map[string]interface{}{
		"Id":                 hc.ID,
		"CallerReference":    hc.CallerReference,
		"HealthCheckConfig":  s.healthCheckConfigToResponse(hc.HealthCheckConfig),
		"HealthCheckVersion": hc.HealthCheckVersion,
	}

	// Output CloudWatchAlarmConfiguration when populated on the store
	// object, or derive minimal info from AlarmIdentifier for
	// CLOUDWATCH_METRIC health checks.
	if hc.CloudWatchAlarmConfiguration != nil {
		cwa := hc.CloudWatchAlarmConfiguration
		cwaMap := map[string]interface{}{
			"AlarmName": cwa.AlarmName,
		}
		if cwa.AlarmRegion != "" {
			cwaMap["AlarmRegion"] = cwa.AlarmRegion
		}
		if len(cwa.Dimensions) > 0 {
			dims := make([]interface{}, len(cwa.Dimensions))
			for i, d := range cwa.Dimensions {
				dims[i] = map[string]interface{}{
					"Name":  d.Name,
					"Value": d.Value,
				}
			}
			cwaMap["Dimensions"] = protocol.XMLElements{ElementName: "Dimension", Items: dims}
		}
		result["CloudWatchAlarmConfiguration"] = cwaMap
	} else if hc.HealthCheckConfig != nil && hc.HealthCheckConfig.AlarmIdentifier != nil {
		result["CloudWatchAlarmConfiguration"] = map[string]interface{}{
			"AlarmName":   hc.HealthCheckConfig.AlarmIdentifier.Name,
			"AlarmRegion": hc.HealthCheckConfig.AlarmIdentifier.Region,
		}
	}

	return result
}
