package route53

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
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

func generateDelegationSetId() string {
	return fmt.Sprintf("N%s", generateId())
}

func generateId() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 12)
	for i := range b {
		n, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}
	return string(b)
}

func extractHostedZoneId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", NewAPIError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
	}
	return strings.TrimPrefix(id, "/hostedzone/"), nil
}

func extractHealthCheckId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", NewAPIError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
	}
	return id, nil
}

func extractChangeId(params map[string]interface{}, paramName string) (string, error) {
	id := request.GetStringParam(params, paramName)
	if id == "" {
		return "", NewAPIError("InvalidInput", fmt.Sprintf("%s is required", paramName), 400)
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

func parseHealthCheckConfig(configMap map[string]interface{}) *route53store.HealthCheckConfig {
	if configMap == nil {
		return nil
	}

	config := &route53store.HealthCheckConfig{
		Type: request.GetStringParam(configMap, "Type"),
	}

	if v, ok := configMap["IPAddress"].(string); ok {
		config.IPAddress = v
	}
	if v, ok := configMap["Port"].(float64); ok {
		config.Port = int64(v)
	}
	if v, ok := configMap["ResourcePath"].(string); ok {
		config.ResourcePath = v
	}
	if v, ok := configMap["FullyQualifiedDomainName"].(string); ok {
		config.FullyQualifiedDomainName = v
	}
	if v, ok := configMap["RequestInterval"].(float64); ok {
		config.RequestInterval = int64(v)
	}
	if v, ok := configMap["FailureThreshold"].(float64); ok {
		config.FailureThreshold = int64(v)
	}
	config.MeasureLatency = request.GetBoolParam(configMap, "MeasureLatency")
	config.Inverted = request.GetBoolParam(configMap, "Inverted")
	config.Disabled = request.GetBoolParam(configMap, "Disabled")
	config.EnableSNI = request.GetBoolParam(configMap, "EnableSNI")

	return config
}

func applyHealthCheckConfigUpdates(config *route53store.HealthCheckConfig, updates map[string]interface{}) {
	if config == nil || updates == nil {
		return
	}

	if v, ok := updates["IPAddress"].(string); ok {
		config.IPAddress = v
	}
	if v, ok := updates["Port"].(float64); ok {
		config.Port = int64(v)
	}
	if v, ok := updates["ResourcePath"].(string); ok {
		config.ResourcePath = v
	}
	if v, ok := updates["FullyQualifiedDomainName"].(string); ok {
		config.FullyQualifiedDomainName = v
	}
	if v, ok := updates["RequestInterval"].(float64); ok {
		config.RequestInterval = int64(v)
	}
	if v, ok := updates["FailureThreshold"].(float64); ok {
		config.FailureThreshold = int64(v)
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

func (s *Route53Service) buildHostedZonesListResponse(zones []*route53store.HostedZone, isTruncated bool, nextMarker string, maxItems int) map[string]interface{} {
	result := make([]interface{}, len(zones))
	for i, z := range zones {
		result[i] = s.hostedZoneToResponse(z)
	}
	response := map[string]interface{}{
		"HostedZones": result,
		"IsTruncated": isTruncated,
		"Marker":      nextMarker,
		"MaxItems":    maxItems,
	}
	if isTruncated && nextMarker != "" {
		response["NextMarker"] = nextMarker
	}
	return response
}

func (s *Route53Service) buildHealthChecksListResponse(healthChecks []*route53store.HealthCheck, isTruncated bool, nextMarker string) map[string]interface{} {
	result := make([]interface{}, len(healthChecks))
	for i, hc := range healthChecks {
		result[i] = map[string]interface{}{
			"Id":                 hc.ID,
			"CallerReference":    hc.CallerReference,
			"HealthCheckConfig":  s.healthCheckConfigToResponse(hc.HealthCheckConfig),
			"HealthCheckVersion": hc.HealthCheckVersion,
		}
	}
	response := map[string]interface{}{
		"HealthChecks": result,
		"IsTruncated":  isTruncated,
	}
	if isTruncated && nextMarker != "" {
		response["NextMarker"] = nextMarker
	}
	return response
}
