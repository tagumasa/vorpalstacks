package cloudfront

import (
	"context"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	"vorpalstacks/internal/store/aws/common"
)

func newCloudFrontError(typeName, message string, statusCode int) *awserrors.AWSError {
	return awserrors.NewAWSError(typeName, message, statusCode)
}

// NewAPIError creates a new CloudFront API error.
func NewAPIError(typeName, message string, statusCode int) *awserrors.AWSError {
	return newCloudFrontError(typeName, message, statusCode)
}

func getIfMatch(req *request.ParsedRequest) string {
	ifMatch := req.Headers.Get("If-Match")
	return strings.TrimSpace(ifMatch)
}

func getDistributionConfigMap(req *request.ParsedRequest) map[string]interface{} {
	if configMap := request.GetMapParam(req.Parameters, "DistributionConfig"); configMap != nil {
		return configMap
	}
	return req.Parameters
}

func formatDistributionResponse(d *cloudfrontstore.Distribution) map[string]interface{} {
	return map[string]interface{}{
		"Id":                            d.ID,
		"ARN":                           d.ARN,
		"ETag":                          d.ETag,
		"Status":                        d.Status,
		"DomainName":                    d.DomainName,
		"DistributionConfig":            formatDistributionConfig(d.DistributionConfig),
		"CallerReference":               d.CallerReference,
		"LastModifiedTime":              d.LastModifiedAt.Format(time.RFC3339),
		"InProgressInvalidationBatches": 0,
		"ActiveTrustedSigners":          map[string]interface{}{"Enabled": false, "Quantity": 0},
		"ActiveTrustedKeyGroups":        map[string]interface{}{"Enabled": false, "Quantity": 0},
	}
}

func formatDistributionConfig(config *cloudfrontstore.DistributionConfig) map[string]interface{} {
	if config == nil {
		return nil
	}
	m := map[string]interface{}{
		"CallerReference":   config.CallerReference,
		"Comment":           config.Comment,
		"DefaultRootObject": config.DefaultRootObject,
		"Enabled":           config.Enabled,
		"PriceClass":        config.PriceClass,
		"HttpVersion":       config.HttpVersion,
		"IsIPV6Enabled":     config.IsIPV6Enabled,
		"Staging":           config.Staging,
		"WebACLId":          config.WebACLId,
	}

	if config.Origins.Quantity > 0 || len(config.Origins.Items) > 0 {
		m["Origins"] = formatOrigins(config.Origins)
	}

	if config.DefaultCacheBehavior != nil {
		m["DefaultCacheBehavior"] = formatCacheBehavior(config.DefaultCacheBehavior)
	}

	if config.CacheBehaviors != nil && (config.CacheBehaviors.Quantity > 0 || len(config.CacheBehaviors.Items) > 0) {
		m["CacheBehaviors"] = formatCacheBehaviors(config.CacheBehaviors)
	}

	if config.CustomErrorResponses != nil {
		m["CustomErrorResponses"] = map[string]interface{}{
			"Quantity": config.CustomErrorResponses.Quantity,
		}
	}

	if config.Logging != nil {
		m["Logging"] = map[string]interface{}{
			"Enabled":        config.Logging.Enabled,
			"IncludeCookies": config.Logging.IncludeCookies,
			"Bucket":         config.Logging.Bucket,
			"Prefix":         config.Logging.Prefix,
		}
	}

	if config.ViewerCertificate != nil {
		m["ViewerCertificate"] = formatViewerCertificate(config.ViewerCertificate)
	}

	if config.Restrictions != nil {
		m["Restrictions"] = map[string]interface{}{
			"GeoRestriction": map[string]interface{}{
				"RestrictionType": config.Restrictions.GeoRestriction.RestrictionType,
				"Quantity":        config.Restrictions.GeoRestriction.Quantity,
			},
		}
	}

	m["OriginGroups"] = map[string]interface{}{
		"Quantity": 0,
	}

	m["AnycastIpListId"] = ""
	m["ContinuousDeploymentPolicyId"] = ""

	if config.Aliases != nil {
		aliasMap := map[string]interface{}{
			"Quantity": config.Aliases.Quantity,
		}
		if len(config.Aliases.Items) > 0 {
			items := make([]interface{}, len(config.Aliases.Items))
			for i, a := range config.Aliases.Items {
				items[i] = a
			}
			aliasMap["Items"] = protocol.XMLElements{ElementName: "CNAME", Items: items}
		}
		m["Aliases"] = aliasMap
	} else {
		m["Aliases"] = map[string]interface{}{
			"Quantity": 0,
		}
	}

	return m
}

func formatOrigins(origins cloudfrontstore.Origins) map[string]interface{} {
	m := map[string]interface{}{
		"Quantity": origins.Quantity,
	}
	if len(origins.Items) > 0 {
		items := make([]interface{}, len(origins.Items))
		for i, o := range origins.Items {
			items[i] = formatOriginMap(o)
		}
		m["Items"] = protocol.XMLElements{ElementName: "Origin", Items: items}
	}
	return m
}

func formatOriginMap(o *cloudfrontstore.Origin) map[string]interface{} {
	om := map[string]interface{}{
		"Id":                    o.ID,
		"DomainName":            o.DomainName,
		"OriginPath":            o.OriginPath,
		"ConnectionAttempts":    o.ConnectionAttempts,
		"ConnectionTimeout":     o.ConnectionTimeout,
		"OriginAccessControlId": o.OriginAccessControlId,
		"CustomHeaders":         map[string]interface{}{"Quantity": 0},
		"OriginShield":          map[string]interface{}{"Enabled": false},
	}
	if o.CustomOriginConfig != nil {
		coc := map[string]interface{}{
			"HTTPPort":               o.CustomOriginConfig.HTTPPort,
			"HTTPSPort":              o.CustomOriginConfig.HTTPSPort,
			"OriginProtocolPolicy":   o.CustomOriginConfig.OriginProtocolPolicy,
			"OriginReadTimeout":      o.CustomOriginConfig.OriginReadTimeout,
			"OriginKeepaliveTimeout": o.CustomOriginConfig.OriginKeepaliveTimeout,
		}
		if len(o.CustomOriginConfig.OriginSslProtocols) > 0 {
			items := make([]interface{}, len(o.CustomOriginConfig.OriginSslProtocols))
			for i, proto := range o.CustomOriginConfig.OriginSslProtocols {
				items[i] = proto
			}
			coc["OriginSslProtocols"] = map[string]interface{}{
				"Quantity": len(items),
				"Items":    protocol.XMLElements{ElementName: "SslProtocol", Items: items},
			}
		}
		om["CustomOriginConfig"] = coc
	}
	if o.S3OriginConfig != nil {
		om["S3OriginConfig"] = map[string]interface{}{
			"OriginAccessIdentity": o.S3OriginConfig.OriginAccessIdentity,
		}
	}
	return om
}

func formatCacheBehavior(cb *cloudfrontstore.CacheBehavior) map[string]interface{} {
	if cb == nil {
		return nil
	}
	m := map[string]interface{}{
		"TargetOriginId":             cb.TargetOriginId,
		"ViewerProtocolPolicy":       cb.ViewerProtocolPolicy,
		"Compress":                   cb.Compress,
		"FieldLevelEncryptionId":     "",
		"RealtimeLogConfigArn":       "",
		"SmoothStreaming":            cb.SmoothStreaming,
		"FunctionAssociations":       map[string]interface{}{"Quantity": 0},
		"LambdaFunctionAssociations": map[string]interface{}{"Quantity": 0},
		"GrpcConfig":                 map[string]interface{}{"Enabled": false},
	}
	if cb.PathPattern != "" {
		m["PathPattern"] = cb.PathPattern
	}
	if cb.CachePolicyId != "" {
		m["CachePolicyId"] = cb.CachePolicyId
	}
	if cb.OriginRequestPolicyId != "" {
		m["OriginRequestPolicyId"] = cb.OriginRequestPolicyId
	}
	if cb.ResponseHeadersPolicyId != "" {
		m["ResponseHeadersPolicyId"] = cb.ResponseHeadersPolicyId
	}
	if cb.AllowedMethods != nil {
		am := map[string]interface{}{
			"Quantity": cb.AllowedMethods.Quantity,
		}
		if len(cb.AllowedMethods.Items) > 0 {
			items := make([]interface{}, len(cb.AllowedMethods.Items))
			for i, item := range cb.AllowedMethods.Items {
				items[i] = item
			}
			am["Items"] = protocol.XMLElements{ElementName: "Method", Items: items}
		}
		if len(cb.AllowedMethods.CachedMethods) > 0 {
			cmItems := make([]interface{}, len(cb.AllowedMethods.CachedMethods))
			for i, item := range cb.AllowedMethods.CachedMethods {
				cmItems[i] = item
			}
			am["CachedMethods"] = map[string]interface{}{
				"Quantity": len(cmItems),
				"Items":    protocol.XMLElements{ElementName: "Method", Items: cmItems},
			}
		}
		m["AllowedMethods"] = am
	}
	if cb.TrustedSigners != nil {
		m["TrustedSigners"] = map[string]interface{}{
			"Enabled":  cb.TrustedSigners.Enabled,
			"Quantity": cb.TrustedSigners.Quantity,
		}
	} else {
		m["TrustedSigners"] = map[string]interface{}{
			"Enabled":  false,
			"Quantity": 0,
		}
	}
	if cb.TrustedKeyGroups != nil {
		m["TrustedKeyGroups"] = map[string]interface{}{
			"Enabled":  cb.TrustedKeyGroups.Enabled,
			"Quantity": cb.TrustedKeyGroups.Quantity,
		}
	} else {
		m["TrustedKeyGroups"] = map[string]interface{}{
			"Enabled":  false,
			"Quantity": 0,
		}
	}
	if cb.ForwardedValues != nil {
		fv := map[string]interface{}{
			"QueryString": cb.ForwardedValues.QueryString,
		}
		if cb.ForwardedValues.Cookies != nil {
			fv["Cookies"] = map[string]interface{}{
				"Forward": cb.ForwardedValues.Cookies.Forward,
			}
		}
		m["ForwardedValues"] = fv
	}
	if cb.MinTTL > 0 || cb.DefaultTTL > 0 || cb.MaxTTL > 0 {
		m["MinTTL"] = cb.MinTTL
		m["DefaultTTL"] = cb.DefaultTTL
		m["MaxTTL"] = cb.MaxTTL
	}
	return m
}

func formatCacheBehaviors(cbs *cloudfrontstore.CacheBehaviors) map[string]interface{} {
	m := map[string]interface{}{
		"Quantity": cbs.Quantity,
	}
	if len(cbs.Items) > 0 {
		items := make([]interface{}, len(cbs.Items))
		for i, cb := range cbs.Items {
			items[i] = formatCacheBehavior(cb)
		}
		m["Items"] = protocol.XMLElements{ElementName: "CacheBehavior", Items: items}
	}
	return m
}

func formatViewerCertificate(vc *cloudfrontstore.ViewerCertificate) map[string]interface{} {
	m := map[string]interface{}{}
	if vc.CloudFrontDefaultCertificate {
		m["CloudFrontDefaultCertificate"] = true
	}
	if vc.ACMCertificateArn != "" {
		m["ACMCertificateArn"] = vc.ACMCertificateArn
	}
	if vc.IAMCertificateId != "" {
		m["IAMCertificateId"] = vc.IAMCertificateId
	}
	if vc.SSLSupportMethod != "" {
		m["SSLSupportMethod"] = vc.SSLSupportMethod
	}
	if vc.MinimumProtocolVersion != "" {
		m["MinimumProtocolVersion"] = vc.MinimumProtocolVersion
	} else if vc.CloudFrontDefaultCertificate {
		m["MinimumProtocolVersion"] = "TLSv1"
	}
	if vc.CertificateSource != "" {
		m["CertificateSource"] = vc.CertificateSource
	}
	return m
}

func parseDistributionConfig(configMap map[string]interface{}) *cloudfrontstore.DistributionConfig {
	if configMap == nil {
		return nil
	}
	config := &cloudfrontstore.DistributionConfig{
		CallerReference:      request.GetStringParam(configMap, "CallerReference"),
		Comment:              request.GetStringParam(configMap, "Comment"),
		Enabled:              request.GetBoolParam(configMap, "Enabled"),
		PriceClass:           request.GetStringParam(configMap, "PriceClass"),
		DefaultRootObject:    request.GetStringParam(configMap, "DefaultRootObject"),
		HttpVersion:          request.GetStringParam(configMap, "HttpVersion"),
		IsIPV6Enabled:        request.GetBoolParam(configMap, "IsIPV6Enabled"),
		Staging:              request.GetBoolParam(configMap, "Staging"),
		WebACLId:             request.GetStringParam(configMap, "WebACLId"),
		Origins:              parseOrigins(request.GetMapParam(configMap, "Origins")),
		DefaultCacheBehavior: parseCacheBehavior(request.GetMapParam(configMap, "DefaultCacheBehavior")),
		CacheBehaviors:       parseCacheBehaviors(request.GetMapParam(configMap, "CacheBehaviors")),
		Aliases:              parseAliases(request.GetMapParam(configMap, "Aliases")),
		Logging:              parseLoggingConfig(request.GetMapParam(configMap, "Logging")),
		ViewerCertificate:    parseViewerCertificate(request.GetMapParam(configMap, "ViewerCertificate")),
		Restrictions:         parseRestrictions(request.GetMapParam(configMap, "Restrictions")),
	}

	if cer := request.GetMapParam(configMap, "CustomErrorResponses"); cer != nil {
		config.CustomErrorResponses = &cloudfrontstore.CustomErrorResponses{
			Quantity: request.GetIntParam(cer, "Quantity"),
		}
	}

	return config
}

func parseAliases(m map[string]interface{}) *cloudfrontstore.Aliases {
	if m == nil {
		return nil
	}
	aliases := &cloudfrontstore.Aliases{
		Quantity: request.GetIntParam(m, "Quantity"),
	}
	if items, ok := m["Items"].([]interface{}); ok {
		for _, item := range items {
			if s, ok := item.(string); ok {
				aliases.Items = append(aliases.Items, s)
			}
		}
	}
	if len(aliases.Items) > 0 && aliases.Quantity == 0 {
		aliases.Quantity = len(aliases.Items)
	}
	return aliases
}

func parseLoggingConfig(m map[string]interface{}) *cloudfrontstore.LoggingConfig {
	if m == nil {
		return nil
	}
	return &cloudfrontstore.LoggingConfig{
		Enabled:        request.GetBoolParam(m, "Enabled"),
		IncludeCookies: request.GetBoolParam(m, "IncludeCookies"),
		Bucket:         request.GetStringParam(m, "Bucket"),
		Prefix:         request.GetStringParam(m, "Prefix"),
	}
}

func parseViewerCertificate(m map[string]interface{}) *cloudfrontstore.ViewerCertificate {
	if m == nil {
		return nil
	}
	vc := &cloudfrontstore.ViewerCertificate{
		ACMCertificateArn:      request.GetStringParam(m, "ACMCertificateArn"),
		IAMCertificateId:       request.GetStringParam(m, "IAMCertificateId"),
		SSLSupportMethod:       request.GetStringParam(m, "SSLSupportMethod"),
		MinimumProtocolVersion: request.GetStringParam(m, "MinimumProtocolVersion"),
		CertificateSource:      request.GetStringParam(m, "CertificateSource"),
	}
	if _, ok := m["CloudFrontDefaultCertificate"]; ok {
		vc.CloudFrontDefaultCertificate = true
	}
	return vc
}

func parseRestrictions(m map[string]interface{}) *cloudfrontstore.Restrictions {
	if m == nil {
		return nil
	}
	restrictions := &cloudfrontstore.Restrictions{}
	if gr := request.GetMapParam(m, "GeoRestriction"); gr != nil {
		restrictions.GeoRestriction = cloudfrontstore.GeoRestriction{
			RestrictionType: request.GetStringParam(gr, "RestrictionType"),
			Quantity:        request.GetIntParam(gr, "Quantity"),
		}
	}
	return restrictions
}

func parseOrigin(m map[string]interface{}) *cloudfrontstore.Origin {
	origin := &cloudfrontstore.Origin{
		ID:                    request.GetStringParam(m, "Id"),
		DomainName:            request.GetStringParam(m, "DomainName"),
		OriginPath:            request.GetStringParam(m, "OriginPath"),
		ConnectionAttempts:    request.GetIntParam(m, "ConnectionAttempts"),
		ConnectionTimeout:     request.GetIntParam(m, "ConnectionTimeout"),
		OriginAccessControlId: request.GetStringParam(m, "OriginAccessControlId"),
	}
	if s3Config := request.GetMapParam(m, "S3OriginConfig"); s3Config != nil {
		origin.S3OriginConfig = &cloudfrontstore.S3OriginConfig{
			OriginAccessIdentity: request.GetStringParam(s3Config, "OriginAccessIdentity"),
		}
	}
	if customConfig := request.GetMapParam(m, "CustomOriginConfig"); customConfig != nil {
		coc := &cloudfrontstore.CustomOriginConfig{
			HTTPPort:               request.GetIntParam(customConfig, "HTTPPort"),
			HTTPSPort:              request.GetIntParam(customConfig, "HTTPSPort"),
			OriginProtocolPolicy:   request.GetStringParam(customConfig, "OriginProtocolPolicy"),
			OriginReadTimeout:      request.GetIntParam(customConfig, "OriginReadTimeout"),
			OriginKeepaliveTimeout: request.GetIntParam(customConfig, "OriginKeepaliveTimeout"),
		}
		if ospMap := request.GetMapParam(customConfig, "OriginSslProtocols"); ospMap != nil {
			if items := ospMap["Items"]; items != nil {
				switch v := items.(type) {
				case []interface{}:
					for _, item := range v {
						if s, ok := item.(string); ok {
							coc.OriginSslProtocols = append(coc.OriginSslProtocols, s)
						}
					}
				case map[string]interface{}:
					if protocols, ok := v["SslProtocol"]; ok {
						switch pv := protocols.(type) {
						case string:
							coc.OriginSslProtocols = append(coc.OriginSslProtocols, pv)
						case []interface{}:
							for _, item := range pv {
								if s, ok := item.(string); ok {
									coc.OriginSslProtocols = append(coc.OriginSslProtocols, s)
								}
							}
						}
					}
				case string:
					coc.OriginSslProtocols = append(coc.OriginSslProtocols, v)
				}
			}
		}
		if coc.HTTPPort == 0 {
			coc.HTTPPort = 80
		}
		if coc.HTTPSPort == 0 {
			coc.HTTPSPort = 443
		}
		if coc.OriginReadTimeout == 0 {
			coc.OriginReadTimeout = 30
		}
		if coc.OriginKeepaliveTimeout == 0 {
			coc.OriginKeepaliveTimeout = 5
		}
		if len(coc.OriginSslProtocols) == 0 {
			coc.OriginSslProtocols = []string{"TLSv1", "TLSv1.1", "TLSv1.2"}
		}
		origin.CustomOriginConfig = coc
	}
	if origin.ConnectionAttempts == 0 {
		origin.ConnectionAttempts = 3
	}
	if origin.ConnectionTimeout == 0 {
		origin.ConnectionTimeout = 10
	}
	return origin
}

func xmlItemsToSlice(items interface{}, elemName string) []map[string]interface{} {
	switch v := items.(type) {
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	case map[string]interface{}:
		if single, ok := v[elemName].(map[string]interface{}); ok {
			return []map[string]interface{}{single}
		}
		if single, ok := v[elemName]; ok {
			if m, ok := single.(map[string]interface{}); ok {
				return []map[string]interface{}{m}
			}
		}
	}
	return nil
}

func parseOrigins(originsMap map[string]interface{}) cloudfrontstore.Origins {
	if originsMap == nil {
		return cloudfrontstore.Origins{}
	}
	origins := cloudfrontstore.Origins{}
	if qty, ok := originsMap["Quantity"].(float64); ok {
		origins.Quantity = int(qty)
	}
	if items, ok := originsMap["Items"]; ok {
		for _, m := range xmlItemsToSlice(items, "Origin") {
			origins.Items = append(origins.Items, parseOrigin(m))
		}
	}
	if len(origins.Items) > 0 && origins.Quantity == 0 {
		origins.Quantity = len(origins.Items)
	}
	return origins
}

func parseCacheBehavior(cbMap map[string]interface{}) *cloudfrontstore.CacheBehavior {
	if cbMap == nil {
		return nil
	}
	cb := &cloudfrontstore.CacheBehavior{
		PathPattern:             request.GetStringParam(cbMap, "PathPattern"),
		TargetOriginId:          request.GetStringParam(cbMap, "TargetOriginId"),
		ViewerProtocolPolicy:    request.GetStringParam(cbMap, "ViewerProtocolPolicy"),
		Compress:                request.GetBoolParam(cbMap, "Compress"),
		CachePolicyId:           request.GetStringParam(cbMap, "CachePolicyId"),
		OriginRequestPolicyId:   request.GetStringParam(cbMap, "OriginRequestPolicyId"),
		ResponseHeadersPolicyId: request.GetStringParam(cbMap, "ResponseHeadersPolicyId"),
		MinTTL:                  request.GetIntParam(cbMap, "MinTTL"),
		DefaultTTL:              request.GetIntParam(cbMap, "DefaultTTL"),
		MaxTTL:                  request.GetIntParam(cbMap, "MaxTTL"),
		SmoothStreaming:         request.GetBoolParam(cbMap, "SmoothStreaming"),
	}
	if am := request.GetMapParam(cbMap, "AllowedMethods"); am != nil {
		cb.AllowedMethods = &cloudfrontstore.AllowedMethods{
			Quantity: request.GetIntParam(am, "Quantity"),
		}
		if items := am["Items"]; items != nil {
			switch v := items.(type) {
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok {
						cb.AllowedMethods.Items = append(cb.AllowedMethods.Items, s)
					}
				}
			case map[string]interface{}:
				if methods, ok := v["Method"].([]interface{}); ok {
					for _, item := range methods {
						if s, ok := item.(string); ok {
							cb.AllowedMethods.Items = append(cb.AllowedMethods.Items, s)
						}
					}
				}
			}
		}
		if cm := request.GetMapParam(am, "CachedMethods"); cm != nil {
			if cmItems := cm["Items"]; cmItems != nil {
				switch v := cmItems.(type) {
				case []interface{}:
					for _, item := range v {
						if s, ok := item.(string); ok {
							cb.AllowedMethods.CachedMethods = append(cb.AllowedMethods.CachedMethods, s)
						}
					}
				case map[string]interface{}:
					if methods, ok := v["Method"].([]interface{}); ok {
						for _, item := range methods {
							if s, ok := item.(string); ok {
								cb.AllowedMethods.CachedMethods = append(cb.AllowedMethods.CachedMethods, s)
							}
						}
					}
				}
			}
		}
	}
	if ts := request.GetMapParam(cbMap, "TrustedSigners"); ts != nil {
		cb.TrustedSigners = &cloudfrontstore.TrustedSigners{
			Enabled:  request.GetBoolParam(ts, "Enabled"),
			Quantity: request.GetIntParam(ts, "Quantity"),
		}
		if cb.TrustedSigners.Quantity == 0 && !cb.TrustedSigners.Enabled {
			cb.TrustedSigners.Quantity = 0
		}
	}
	if tkg := request.GetMapParam(cbMap, "TrustedKeyGroups"); tkg != nil {
		cb.TrustedKeyGroups = &cloudfrontstore.TrustedKeyGroups{
			Enabled:  request.GetBoolParam(tkg, "Enabled"),
			Quantity: request.GetIntParam(tkg, "Quantity"),
		}
	}
	if fv := request.GetMapParam(cbMap, "ForwardedValues"); fv != nil {
		cb.ForwardedValues = &cloudfrontstore.ForwardedValues{
			QueryString: request.GetBoolParam(fv, "QueryString"),
		}
		if cookies := request.GetMapParam(fv, "Cookies"); cookies != nil {
			cb.ForwardedValues.Cookies = &cloudfrontstore.CookiePreferences{
				Forward: request.GetStringParam(cookies, "Forward"),
			}
		}
	}
	return cb
}

func parseCacheBehaviors(cbsMap map[string]interface{}) *cloudfrontstore.CacheBehaviors {
	if cbsMap == nil {
		return nil
	}
	cbs := &cloudfrontstore.CacheBehaviors{}
	if qty, ok := cbsMap["Quantity"].(float64); ok {
		cbs.Quantity = int(qty)
	}
	if items, ok := cbsMap["Items"]; ok {
		for _, m := range xmlItemsToSlice(items, "CacheBehavior") {
			if cb := parseCacheBehavior(m); cb != nil {
				cbs.Items = append(cbs.Items, cb)
			}
		}
	}
	if len(cbs.Items) > 0 && cbs.Quantity == 0 {
		cbs.Quantity = len(cbs.Items)
	}
	return cbs
}

// CreateDistribution creates a new CloudFront distribution.
func (s *CloudFrontService) CreateDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := getDistributionConfigMap(req)

	callerRef := request.GetStringParam(configMap, "CallerReference")
	if callerRef == "" {
		callerRef = fmt.Sprintf("ref-%d", time.Now().UnixNano())
	}

	config := parseDistributionConfig(configMap)
	config.CallerReference = callerRef

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Create(callerRef, config)
	if err != nil {
		return nil, err
	}

	if err := s.syncWAFAssociation(reqCtx, config.WebACLId, distribution.ARN); err != nil {
		return nil, fmt.Errorf("failed to sync WAF association: %w", err)
	}

	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution),
	}, nil
}

// CreateDistributionWithTags creates a new CloudFront distribution with the specified tags.
func (s *CloudFrontService) CreateDistributionWithTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "DistributionConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	callerRef := request.GetStringParam(configMap, "CallerReference")
	if callerRef == "" {
		callerRef = fmt.Sprintf("ref-%d", time.Now().UnixNano())
	}

	config := parseDistributionConfig(configMap)
	config.CallerReference = callerRef

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Create(callerRef, config)
	if err != nil {
		return nil, err
	}

	if err := s.syncWAFAssociation(reqCtx, config.WebACLId, distribution.ARN); err != nil {
		return nil, fmt.Errorf("failed to sync WAF association: %w", err)
	}

	var tags []common.Tag
	tagsMap := request.GetMapParam(req.Parameters, "Tags")
	if tagsMap != nil {
		if itemsVal := tagsMap["Items"]; itemsVal != nil {
			switch v := itemsVal.(type) {
			case []interface{}:
				for _, t := range v {
					if m, ok := t.(map[string]interface{}); ok {
						tags = append(tags, common.Tag{Key: request.GetStringParam(m, "Key"), Value: request.GetStringParam(m, "Value")})
					}
				}
			case map[string]interface{}:
				if tagVal, ok := v["Tag"]; ok {
					switch tv := tagVal.(type) {
					case []interface{}:
						for _, t := range tv {
							if m, ok := t.(map[string]interface{}); ok {
								tags = append(tags, common.Tag{Key: request.GetStringParam(m, "Key"), Value: request.GetStringParam(m, "Value")})
							}
						}
					case map[string]interface{}:
						tags = append(tags, common.Tag{Key: request.GetStringParam(tv, "Key"), Value: request.GetStringParam(tv, "Value")})
					}
				}
			}
		}
		if tagVal := tagsMap["Tag"]; tagVal != nil && len(tags) == 0 {
			switch v := tagVal.(type) {
			case []interface{}:
				for _, t := range v {
					if m, ok := t.(map[string]interface{}); ok {
						tags = append(tags, common.Tag{Key: request.GetStringParam(m, "Key"), Value: request.GetStringParam(m, "Value")})
					}
				}
			case map[string]interface{}:
				tags = append(tags, common.Tag{Key: request.GetStringParam(v, "Key"), Value: request.GetStringParam(v, "Value")})
			}
		}
	}
	if len(tags) > 0 && distribution.ARN != "" {
		if err := store.tags.TagResource(distribution.ARN, tags); err != nil {
			return nil, fmt.Errorf("failed to tag distribution: %w", err)
		}
	}

	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution),
	}, nil
}

// GetDistribution retrieves a CloudFront distribution by ID.
func (s *CloudFrontService) GetDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, newCloudFrontError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution),
	}, nil
}

// GetDistributionConfig retrieves the configuration of a CloudFront distribution.
func (s *CloudFrontService) GetDistributionConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, newCloudFrontError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"DistributionConfig": formatDistributionConfig(distribution.DistributionConfig),
		"ETag":               distribution.ETag,
	}, nil
}

// ListDistributions lists all CloudFront distributions.
func (s *CloudFrontService) ListDistributions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.distributions.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Distributions))
	for _, d := range result.Distributions {
		comment := ""
		if d.DistributionConfig != nil {
			comment = d.DistributionConfig.Comment
		}
		items = append(items, map[string]interface{}{
			"Id":               d.ID,
			"ARN":              d.ARN,
			"Status":           d.Status,
			"DomainName":       d.DomainName,
			"Enabled":          d.Enabled,
			"CallerReference":  d.CallerReference,
			"Comment":          comment,
			"LastModifiedTime": d.LastModifiedAt.Format(time.RFC3339),
		})
	}

	return map[string]interface{}{
		"DistributionList": map[string]interface{}{
			"Marker":      marker,
			"MaxItems":    maxItems,
			"IsTruncated": result.IsTruncated,
			"Quantity":    len(items),
			"NextMarker":  result.NextMarker,
			"Items":       protocol.XMLElements{ElementName: "DistributionSummary", Items: items},
		},
	}, nil
}

// UpdateDistribution updates the configuration of a CloudFront distribution.
func (s *CloudFrontService) UpdateDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, newCloudFrontError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)

	configMap := getDistributionConfigMap(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	if ifMatch != "" && ifMatch != "*" && distribution.ETag != ifMatch {
		return nil, newCloudFrontError("PreconditionFailed", "The If-Match header does not match the current ETag", 412)
	}

	newConfig := parseDistributionConfig(configMap)
	if newConfig.CallerReference == "" {
		newConfig.CallerReference = distribution.CallerReference
	}
	distribution.DistributionConfig = newConfig
	distribution.Enabled = newConfig.Enabled
	distribution.ETag = s.generateETag()

	if oldWebACLId := distribution.DistributionConfig.WebACLId; oldWebACLId != newConfig.WebACLId {
		distArn := distribution.ARN
		if err := s.removeWAFAssociation(reqCtx, oldWebACLId, distArn); err != nil {
			return nil, fmt.Errorf("failed to remove old WAF association: %w", err)
		}
		if err := s.syncWAFAssociation(reqCtx, newConfig.WebACLId, distArn); err != nil {
			return nil, fmt.Errorf("failed to sync WAF association: %w", err)
		}
	}

	if err := store.distributions.UpdateWithLastModified(id, distribution); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution),
	}, nil
}

// DeleteDistribution deletes a CloudFront distribution. The distribution must be disabled first.
func (s *CloudFrontService) DeleteDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, newCloudFrontError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	if ifMatch != "" && ifMatch != "*" && distribution.ETag != ifMatch {
		return nil, newCloudFrontError("PreconditionFailed", "The If-Match header does not match the current ETag", 412)
	}

	if distribution.Enabled {
		return nil, newCloudFrontError("DistributionNotDisabled", "Distribution must be disabled before deletion", 409)
	}

	if distribution.DistributionConfig != nil {
		if err := s.removeWAFAssociation(reqCtx, distribution.DistributionConfig.WebACLId, distribution.ARN); err != nil {
			return nil, fmt.Errorf("failed to remove WAF association: %w", err)
		}
	}

	err = store.distributions.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, newCloudFrontError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListDistributionsByWebACLId lists distributions that are associated with a specified Web ACL.
// Supports pagination via Marker/MaxItems, returning IsTruncated and NextMarker
// when more results are available than the requested page size.
func (s *CloudFrontService) ListDistributionsByWebACLId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLId := request.GetStringParam(req.Parameters, "WebACLId")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allDistributions, err := store.distributions.List("", 0)
	if err != nil {
		return nil, err
	}

	var matched []interface{}
	for _, d := range allDistributions.Distributions {
		if d.DistributionConfig != nil && d.DistributionConfig.WebACLId == webACLId {
			matched = append(matched, map[string]interface{}{
				"Id":               d.ID,
				"ARN":              d.ARN,
				"Status":           d.Status,
				"DomainName":       d.DomainName,
				"Enabled":          d.Enabled,
				"CallerReference":  d.CallerReference,
				"LastModifiedTime": d.LastModifiedAt.Format(time.RFC3339),
			})
		}
	}

	skipCount := 0
	if marker != "" {
		for i, item := range matched {
			if summary, ok := item.(map[string]interface{}); ok {
				if id, ok := summary["Id"].(string); ok && id == marker {
					skipCount = i + 1
					break
				}
			}
		}
	}

	paged := matched[skipCount:]
	isTruncated := len(paged) > maxItems
	if isTruncated {
		paged = paged[:maxItems]
	}

	nextMarker := ""
	if isTruncated && len(paged) > 0 {
		if summary, ok := paged[len(paged)-1].(map[string]interface{}); ok {
			if id, ok := summary["Id"].(string); ok {
				nextMarker = id
			}
		}
	}

	return map[string]interface{}{
		"DistributionList": map[string]interface{}{
			"Marker":      marker,
			"MaxItems":    maxItems,
			"IsTruncated": isTruncated,
			"NextMarker":  nextMarker,
			"Quantity":    len(paged),
			"Items":       protocol.XMLElements{ElementName: "DistributionSummary", Items: paged},
		},
	}, nil
}
