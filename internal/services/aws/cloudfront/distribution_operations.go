package cloudfront

import (
	"context"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	types "vorpalstacks/internal/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

func getIfMatch(req *request.ParsedRequest) string {
	ifMatch := req.Headers.Get("If-Match")
	return strings.TrimSpace(ifMatch)
}

// preconditionFailedETag is the standard error message returned when an
// If-Match header does not match the current ETag.
const preconditionFailedETagMsg = "The If-Match header does not match the current ETag"

func getDistributionConfigMap(req *request.ParsedRequest) map[string]interface{} {
	if configMap := request.GetMapParam(req.Parameters, "DistributionConfig"); configMap != nil {
		return configMap
	}
	return req.Parameters
}

func formatDistributionResponse(d *cloudfrontstore.Distribution, inProgressInvalidations int, activeSigners, activeKeyGroups map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Id":                            d.ID,
		"ARN":                           d.ARN,
		"ETag":                          d.ETag,
		"Status":                        d.Status,
		"DomainName":                    d.DomainName,
		"DistributionConfig":            formatDistributionConfig(d.DistributionConfig),
		"CallerReference":               d.CallerReference,
		"LastModifiedTime":              d.LastModifiedAt.Format(time.RFC3339),
		"InProgressInvalidationBatches": inProgressInvalidations,
		"ActiveTrustedSigners":          activeSigners,
		"ActiveTrustedKeyGroups":        activeKeyGroups,
	}
}

// computeActiveTrustedSigners inspects all cache behaviours in the
// distribution config for TrustedSigners with Enabled=true and produces
// the ActiveTrustedSigners output shape. KeyPairIds are empty because the
// legacy CloudFront key-pair system (account-level key pairs uploaded via
// the AWS credentials page) is not implemented on this platform.
func computeActiveTrustedSigners(d *cloudfrontstore.Distribution) map[string]interface{} {
	accounts := collectTrustedSignerAccounts(d)
	if len(accounts) == 0 {
		return map[string]interface{}{"Enabled": false, "Quantity": 0}
	}
	items := make([]interface{}, 0, len(accounts))
	for acct := range accounts {
		items = append(items, map[string]interface{}{
			"AwsAccountNumber": acct,
			"KeyPairIds":       map[string]interface{}{"Quantity": 0},
		})
	}
	return map[string]interface{}{
		"Enabled":  true,
		"Quantity": len(accounts),
		"Items":    protocol.XMLElements{ElementName: "Signer", Items: items},
	}
}

// computeActiveTrustedKeyGroups inspects all cache behaviours for
// TrustedKeyGroups with Enabled=true and produces the
// ActiveTrustedKeyGroups output shape. Each key group's PublicKey IDs are
// resolved from the KeyGroup store to populate KeyPairIds.
func computeActiveTrustedKeyGroups(d *cloudfrontstore.Distribution, stores *cloudfrontStores) map[string]interface{} {
	kgIDs := collectTrustedKeyGroupIDs(d)
	if len(kgIDs) == 0 {
		return map[string]interface{}{"Enabled": false, "Quantity": 0}
	}
	items := make([]interface{}, 0, len(kgIDs))
	for kgID := range kgIDs {
		keyPairIds := map[string]interface{}{"Quantity": 0}
		if stores != nil {
			if kg, err := stores.keyGroups.Get(kgID); err == nil && len(kg.KeyGroupConfig.Items) > 0 {
				kpItems := make([]interface{}, len(kg.KeyGroupConfig.Items))
				for i, kp := range kg.KeyGroupConfig.Items {
					kpItems[i] = kp
				}
				keyPairIds = map[string]interface{}{
					"Quantity": len(kg.KeyGroupConfig.Items),
					"Items":    protocol.XMLElements{ElementName: "KeyPairId", Items: kpItems},
				}
			}
		}
		items = append(items, map[string]interface{}{
			"KeyGroupId": kgID,
			"KeyPairIds": keyPairIds,
		})
	}
	return map[string]interface{}{
		"Enabled":  true,
		"Quantity": len(kgIDs),
		"Items":    protocol.XMLElements{ElementName: "KeyGroup", Items: items},
	}
}

// collectTrustedSignerAccounts returns the set of unique AwsAccountNumbers
// from TrustedSigners across all cache behaviours where Enabled=true.
func collectTrustedSignerAccounts(d *cloudfrontstore.Distribution) map[string]bool {
	accounts := make(map[string]bool)
	if d.DistributionConfig == nil {
		return accounts
	}
	cfg := d.DistributionConfig
	if cfg.DefaultCacheBehavior != nil && cfg.DefaultCacheBehavior.TrustedSigners != nil && cfg.DefaultCacheBehavior.TrustedSigners.Enabled {
		for _, acct := range cfg.DefaultCacheBehavior.TrustedSigners.Items {
			accounts[acct] = true
		}
	}
	if cfg.CacheBehaviors != nil {
		for _, cb := range cfg.CacheBehaviors.Items {
			if cb != nil && cb.TrustedSigners != nil && cb.TrustedSigners.Enabled {
				for _, acct := range cb.TrustedSigners.Items {
					accounts[acct] = true
				}
			}
		}
	}
	return accounts
}

// collectTrustedKeyGroupIDs returns the set of unique KeyGroupId values
// from TrustedKeyGroups across all cache behaviours where Enabled=true.
func collectTrustedKeyGroupIDs(d *cloudfrontstore.Distribution) map[string]bool {
	kgIDs := make(map[string]bool)
	if d.DistributionConfig == nil {
		return kgIDs
	}
	cfg := d.DistributionConfig
	if cfg.DefaultCacheBehavior != nil && cfg.DefaultCacheBehavior.TrustedKeyGroups != nil && cfg.DefaultCacheBehavior.TrustedKeyGroups.Enabled {
		for _, kgID := range cfg.DefaultCacheBehavior.TrustedKeyGroups.Items {
			kgIDs[kgID] = true
		}
	}
	if cfg.CacheBehaviors != nil {
		for _, cb := range cfg.CacheBehaviors.Items {
			if cb != nil && cb.TrustedKeyGroups != nil && cb.TrustedKeyGroups.Enabled {
				for _, kgID := range cb.TrustedKeyGroups.Items {
					kgIDs[kgID] = true
				}
			}
		}
	}
	return kgIDs
}

func formatDistributionSummary(d *cloudfrontstore.Distribution) map[string]interface{} {
	m := map[string]interface{}{
		"Id":               d.ID,
		"ARN":              d.ARN,
		"ETag":             d.ETag,
		"Status":           d.Status,
		"LastModifiedTime": d.LastModifiedAt.Format(time.RFC3339),
		"DomainName":       d.DomainName,
		"Enabled":          d.Enabled,
		"Staging":          d.Staging,
	}

	if d.DistributionConfig != nil {
		cfg := d.DistributionConfig
		m["Comment"] = cfg.Comment
		m["PriceClass"] = cfg.PriceClass
		m["HttpVersion"] = cfg.HttpVersion
		m["IsIPV6Enabled"] = cfg.IsIPV6Enabled
		m["WebACLId"] = cfg.WebACLId

		if cfg.Aliases != nil {
			aliasMap := map[string]interface{}{"Quantity": cfg.Aliases.Quantity}
			if len(cfg.Aliases.Items) > 0 {
				items := make([]interface{}, len(cfg.Aliases.Items))
				for i, a := range cfg.Aliases.Items {
					items[i] = a
				}
				aliasMap["Items"] = protocol.XMLElements{ElementName: "CNAME", Items: items}
			}
			m["Aliases"] = aliasMap
		} else {
			m["Aliases"] = map[string]interface{}{"Quantity": 0}
		}

		if cfg.Origins.Quantity > 0 || len(cfg.Origins.Items) > 0 {
			m["Origins"] = formatOrigins(cfg.Origins)
		}

		if cfg.DefaultCacheBehavior != nil {
			m["DefaultCacheBehavior"] = formatCacheBehavior(cfg.DefaultCacheBehavior)
		}

		if cfg.CacheBehaviors != nil && (cfg.CacheBehaviors.Quantity > 0 || len(cfg.CacheBehaviors.Items) > 0) {
			m["CacheBehaviors"] = formatCacheBehaviors(cfg.CacheBehaviors)
		}

		if cfg.CustomErrorResponses != nil {
			cerMap := map[string]interface{}{"Quantity": cfg.CustomErrorResponses.Quantity}
			if len(cfg.CustomErrorResponses.Items) > 0 {
				items := make([]interface{}, len(cfg.CustomErrorResponses.Items))
				for i, r := range cfg.CustomErrorResponses.Items {
					itemMap := map[string]interface{}{"ErrorCode": r.ErrorCode}
					if r.ResponsePagePath != "" {
						itemMap["ResponsePagePath"] = r.ResponsePagePath
					}
					if r.ResponseCode != "" {
						itemMap["ResponseCode"] = r.ResponseCode
					}
					if r.ErrorCachingMinTTL != 0 {
						itemMap["ErrorCachingMinTTL"] = r.ErrorCachingMinTTL
					}
					items[i] = itemMap
				}
				cerMap["Items"] = protocol.XMLElements{ElementName: "CustomErrorResponse", Items: items}
			}
			m["CustomErrorResponses"] = cerMap
		}

		if cfg.ViewerCertificate != nil {
			m["ViewerCertificate"] = formatViewerCertificate(cfg.ViewerCertificate)
		}

		if cfg.Restrictions != nil {
			geoMap := map[string]interface{}{
				"RestrictionType": cfg.Restrictions.GeoRestriction.RestrictionType,
				"Quantity":        cfg.Restrictions.GeoRestriction.Quantity,
			}
			if len(cfg.Restrictions.GeoRestriction.Items) > 0 {
				items := make([]interface{}, len(cfg.Restrictions.GeoRestriction.Items))
				for i, loc := range cfg.Restrictions.GeoRestriction.Items {
					items[i] = loc
				}
				geoMap["Items"] = protocol.XMLElements{ElementName: "Location", Items: items}
			}
			m["Restrictions"] = map[string]interface{}{"GeoRestriction": geoMap}
		}
	}

	m["OriginGroups"] = map[string]interface{}{"Quantity": 0}

	return m
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
		cerMap := map[string]interface{}{
			"Quantity": config.CustomErrorResponses.Quantity,
		}
		if len(config.CustomErrorResponses.Items) > 0 {
			items := make([]interface{}, len(config.CustomErrorResponses.Items))
			for i, r := range config.CustomErrorResponses.Items {
				itemMap := map[string]interface{}{
					"ErrorCode": r.ErrorCode,
				}
				if r.ResponsePagePath != "" {
					itemMap["ResponsePagePath"] = r.ResponsePagePath
				}
				if r.ResponseCode != "" {
					itemMap["ResponseCode"] = r.ResponseCode
				}
				if r.ErrorCachingMinTTL != 0 {
					itemMap["ErrorCachingMinTTL"] = r.ErrorCachingMinTTL
				}
				items[i] = itemMap
			}
			cerMap["Items"] = protocol.XMLElements{ElementName: "CustomErrorResponse", Items: items}
		}
		m["CustomErrorResponses"] = cerMap
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
		geoMap := map[string]interface{}{
			"RestrictionType": config.Restrictions.GeoRestriction.RestrictionType,
			"Quantity":        config.Restrictions.GeoRestriction.Quantity,
		}
		if len(config.Restrictions.GeoRestriction.Items) > 0 {
			items := make([]interface{}, len(config.Restrictions.GeoRestriction.Items))
			for i, loc := range config.Restrictions.GeoRestriction.Items {
				items[i] = loc
			}
			geoMap["Items"] = protocol.XMLElements{ElementName: "Location", Items: items}
		}
		m["Restrictions"] = map[string]interface{}{
			"GeoRestriction": geoMap,
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
	}
	if o.CustomHeaders != nil {
		chm := map[string]interface{}{
			"Quantity": o.CustomHeaders.Quantity,
		}
		if len(o.CustomHeaders.Items) > 0 {
			items := make([]interface{}, len(o.CustomHeaders.Items))
			for i, h := range o.CustomHeaders.Items {
				items[i] = map[string]interface{}{
					"HeaderName":  h.HeaderName,
					"HeaderValue": h.HeaderValue,
				}
			}
			chm["Items"] = protocol.XMLElements{ElementName: "OriginCustomHeader", Items: items}
		}
		om["CustomHeaders"] = chm
	} else {
		om["CustomHeaders"] = map[string]interface{}{"Quantity": 0}
	}
	if o.OriginShield != nil {
		osm := map[string]interface{}{
			"Enabled": o.OriginShield.Enabled,
		}
		if o.OriginShield.OriginShieldRegion != "" {
			osm["OriginShieldRegion"] = o.OriginShield.OriginShieldRegion
		}
		om["OriginShield"] = osm
	} else {
		om["OriginShield"] = map[string]interface{}{"Enabled": false}
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
		"TargetOriginId":         cb.TargetOriginId,
		"ViewerProtocolPolicy":   cb.ViewerProtocolPolicy,
		"Compress":               cb.Compress,
		"FieldLevelEncryptionId": "",
		"RealtimeLogConfigArn":   "",
		"SmoothStreaming":        cb.SmoothStreaming,
		"GrpcConfig":             map[string]interface{}{"Enabled": false},
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
		tsm := map[string]interface{}{
			"Enabled":  cb.TrustedSigners.Enabled,
			"Quantity": cb.TrustedSigners.Quantity,
		}
		if len(cb.TrustedSigners.Items) > 0 {
			items := make([]interface{}, len(cb.TrustedSigners.Items))
			for i, item := range cb.TrustedSigners.Items {
				items[i] = item
			}
			tsm["Items"] = protocol.XMLElements{ElementName: "AwsAccountNumber", Items: items}
		}
		m["TrustedSigners"] = tsm
	} else {
		m["TrustedSigners"] = map[string]interface{}{
			"Enabled":  false,
			"Quantity": 0,
		}
	}
	if cb.TrustedKeyGroups != nil {
		tkgm := map[string]interface{}{
			"Enabled":  cb.TrustedKeyGroups.Enabled,
			"Quantity": cb.TrustedKeyGroups.Quantity,
		}
		if len(cb.TrustedKeyGroups.Items) > 0 {
			items := make([]interface{}, len(cb.TrustedKeyGroups.Items))
			for i, item := range cb.TrustedKeyGroups.Items {
				items[i] = item
			}
			tkgm["Items"] = protocol.XMLElements{ElementName: "KeyGroup", Items: items}
		}
		m["TrustedKeyGroups"] = tkgm
	} else {
		m["TrustedKeyGroups"] = map[string]interface{}{
			"Enabled":  false,
			"Quantity": 0,
		}
	}
	if cb.ForwardedValues != nil {
		m["ForwardedValues"] = formatForwardedValues(cb.ForwardedValues)
	}
	if cb.LambdaFunctionAssociations != nil {
		lfam := map[string]interface{}{
			"Quantity": cb.LambdaFunctionAssociations.Quantity,
		}
		if len(cb.LambdaFunctionAssociations.Items) > 0 {
			items := make([]interface{}, len(cb.LambdaFunctionAssociations.Items))
			for i, item := range cb.LambdaFunctionAssociations.Items {
				items[i] = map[string]interface{}{
					"LambdaFunctionARN": item.LambdaFunctionARN,
					"EventType":         item.EventType,
					"IncludeBody":       item.IncludeBody,
				}
			}
			lfam["Items"] = protocol.XMLElements{ElementName: "LambdaFunctionAssociation", Items: items}
		}
		m["LambdaFunctionAssociations"] = lfam
	} else {
		m["LambdaFunctionAssociations"] = map[string]interface{}{"Quantity": 0}
	}
	if cb.FunctionAssociations != nil {
		fam := map[string]interface{}{
			"Quantity": cb.FunctionAssociations.Quantity,
		}
		if len(cb.FunctionAssociations.Items) > 0 {
			items := make([]interface{}, len(cb.FunctionAssociations.Items))
			for i, item := range cb.FunctionAssociations.Items {
				items[i] = map[string]interface{}{
					"FunctionARN": item.FunctionARN,
					"EventType":   item.EventType,
				}
			}
			fam["Items"] = protocol.XMLElements{ElementName: "FunctionAssociation", Items: items}
		}
		m["FunctionAssociations"] = fam
	} else {
		m["FunctionAssociations"] = map[string]interface{}{"Quantity": 0}
	}
	if cb.MinTTL > 0 || cb.DefaultTTL > 0 || cb.MaxTTL > 0 {
		m["MinTTL"] = cb.MinTTL
		m["DefaultTTL"] = cb.DefaultTTL
		m["MaxTTL"] = cb.MaxTTL
	}
	return m
}

func formatForwardedValues(fv *cloudfrontstore.ForwardedValues) map[string]interface{} {
	m := map[string]interface{}{
		"QueryString": fv.QueryString,
	}
	if fv.Cookies != nil {
		cm := map[string]interface{}{
			"Forward": fv.Cookies.Forward,
		}
		if fv.Cookies.WhitelistedNames != nil {
			wnm := map[string]interface{}{
				"Quantity": fv.Cookies.WhitelistedNames.Quantity,
			}
			if len(fv.Cookies.WhitelistedNames.Items) > 0 {
				items := make([]interface{}, len(fv.Cookies.WhitelistedNames.Items))
				for i, name := range fv.Cookies.WhitelistedNames.Items {
					items[i] = name
				}
				wnm["Items"] = protocol.XMLElements{ElementName: "Name", Items: items}
			}
			cm["WhitelistedNames"] = wnm
		}
		m["Cookies"] = cm
	}
	if fv.Headers != nil {
		hm := map[string]interface{}{
			"Quantity": fv.Headers.Quantity,
		}
		if len(fv.Headers.Items) > 0 {
			items := make([]interface{}, len(fv.Headers.Items))
			for i, h := range fv.Headers.Items {
				items[i] = h
			}
			hm["Items"] = protocol.XMLElements{ElementName: "Name", Items: items}
		}
		m["Headers"] = hm
	}
	if fv.QueryStringCacheKeys != nil {
		qm := map[string]interface{}{
			"Quantity": fv.QueryStringCacheKeys.Quantity,
		}
		if len(fv.QueryStringCacheKeys.Items) > 0 {
			items := make([]interface{}, len(fv.QueryStringCacheKeys.Items))
			for i, k := range fv.QueryStringCacheKeys.Items {
				items[i] = k
			}
			qm["Items"] = protocol.XMLElements{ElementName: "Name", Items: items}
		}
		m["QueryStringCacheKeys"] = qm
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
		m["MinimumProtocolVersion"] = "TLSv1.2_2021"
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
		config.CustomErrorResponses = parseCustomErrorResponses(cer)
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
	} else if itemsMap, ok := m["Items"].(map[string]interface{}); ok {
		xmlStringItemsToSlice(itemsMap, "CNAME", &aliases.Items)
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

func parseCustomErrorResponses(m map[string]interface{}) *cloudfrontstore.CustomErrorResponses {
	if m == nil {
		return nil
	}
	cer := &cloudfrontstore.CustomErrorResponses{
		Quantity: request.GetIntParam(m, "Quantity"),
	}
	for _, item := range xmlItemsToSlice(m["Items"], "CustomErrorResponse") {
		resp := cloudfrontstore.CustomErrorResponse{
			ErrorCode: request.GetIntParam(item, "ErrorCode"),
		}
		if v, ok := item["ResponsePagePath"]; ok {
			if s, ok := v.(string); ok {
				resp.ResponsePagePath = s
			}
		}
		if v, ok := item["ResponseCode"]; ok {
			if s, ok := v.(string); ok {
				resp.ResponseCode = s
			}
		}
		resp.ErrorCachingMinTTL = request.GetIntParam(item, "ErrorCachingMinTTL")
		cer.Items = append(cer.Items, resp)
	}
	if len(cer.Items) > 0 && cer.Quantity == 0 {
		cer.Quantity = len(cer.Items)
	}
	return cer
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
	vc.CloudFrontDefaultCertificate = request.GetBoolParam(m, "CloudFrontDefaultCertificate")
	return vc
}

func parseRestrictions(m map[string]interface{}) *cloudfrontstore.Restrictions {
	if m == nil {
		return nil
	}
	restrictions := &cloudfrontstore.Restrictions{}
	if gr := request.GetMapParam(m, "GeoRestriction"); gr != nil {
		geo := cloudfrontstore.GeoRestriction{
			RestrictionType: request.GetStringParam(gr, "RestrictionType"),
			Quantity:        request.GetIntParam(gr, "Quantity"),
		}
		geo.Items = parseStringItemList(gr, "Items", "Location")
		if len(geo.Items) > 0 && geo.Quantity == 0 {
			geo.Quantity = len(geo.Items)
		}
		restrictions.GeoRestriction = geo
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
	if chMap := request.GetMapParam(m, "CustomHeaders"); chMap != nil {
		ch := &cloudfrontstore.CustomHeaders{
			Quantity: request.GetIntParam(chMap, "Quantity"),
		}
		for _, item := range xmlItemsToSlice(chMap["Items"], "OriginCustomHeader") {
			ch.Items = append(ch.Items, cloudfrontstore.OriginCustomHeader{
				HeaderName:  request.GetStringParam(item, "HeaderName"),
				HeaderValue: request.GetStringParam(item, "HeaderValue"),
			})
		}
		if len(ch.Items) > 0 && ch.Quantity == 0 {
			ch.Quantity = len(ch.Items)
		}
		origin.CustomHeaders = ch
	}
	if osMap := request.GetMapParam(m, "OriginShield"); osMap != nil {
		origin.OriginShield = &cloudfrontstore.OriginShield{
			Enabled:            request.GetBoolParam(osMap, "Enabled"),
			OriginShieldRegion: request.GetStringParam(osMap, "OriginShieldRegion"),
		}
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
			coc.OriginSslProtocols = []string{"TLSv1.2"}
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

func xmlStringItemsToSlice(itemsMap map[string]interface{}, elemName string, out *[]string) {
	if arr, ok := itemsMap[elemName].([]interface{}); ok {
		for _, v := range arr {
			if s, ok := v.(string); ok {
				*out = append(*out, s)
			}
		}
	} else if s, ok := itemsMap[elemName].(string); ok {
		*out = append(*out, s)
	}
}

// parseStringItemList extracts a list of strings from a map's "Items" key,
// handling both JSON array and XML nested-map forms.
func parseStringItemList(m map[string]interface{}, key string, elemName string) []string {
	if m == nil {
		return nil
	}
	itemsVal, ok := m[key]
	if !ok || itemsVal == nil {
		return nil
	}
	switch v := itemsVal.(type) {
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case map[string]interface{}:
		var result []string
		xmlStringItemsToSlice(v, elemName, &result)
		return result
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
			Items:    parseStringItemList(ts, "Items", "AwsAccountNumber"),
		}
	}
	if tkg := request.GetMapParam(cbMap, "TrustedKeyGroups"); tkg != nil {
		cb.TrustedKeyGroups = &cloudfrontstore.TrustedKeyGroups{
			Enabled:  request.GetBoolParam(tkg, "Enabled"),
			Quantity: request.GetIntParam(tkg, "Quantity"),
			Items:    parseStringItemList(tkg, "Items", "KeyGroup"),
		}
	}
	if fv := request.GetMapParam(cbMap, "ForwardedValues"); fv != nil {
		cb.ForwardedValues = parseForwardedValues(fv)
	}
	if lfa := request.GetMapParam(cbMap, "LambdaFunctionAssociations"); lfa != nil {
		lfaResult := &cloudfrontstore.LambdaFunctionAssociations{
			Quantity: request.GetIntParam(lfa, "Quantity"),
		}
		for _, item := range xmlItemsToSlice(lfa["Items"], "LambdaFunctionAssociation") {
			lfaResult.Items = append(lfaResult.Items, cloudfrontstore.LambdaFunctionAssociation{
				LambdaFunctionARN: request.GetStringParam(item, "LambdaFunctionARN"),
				EventType:         request.GetStringParam(item, "EventType"),
				IncludeBody:       request.GetBoolParam(item, "IncludeBody"),
			})
		}
		if len(lfaResult.Items) > 0 && lfaResult.Quantity == 0 {
			lfaResult.Quantity = len(lfaResult.Items)
		}
		cb.LambdaFunctionAssociations = lfaResult
	}
	if fa := request.GetMapParam(cbMap, "FunctionAssociations"); fa != nil {
		faResult := &cloudfrontstore.FunctionAssociations{
			Quantity: request.GetIntParam(fa, "Quantity"),
		}
		for _, item := range xmlItemsToSlice(fa["Items"], "FunctionAssociation") {
			faResult.Items = append(faResult.Items, cloudfrontstore.FunctionAssociation{
				FunctionARN: request.GetStringParam(item, "FunctionARN"),
				EventType:   request.GetStringParam(item, "EventType"),
			})
		}
		if len(faResult.Items) > 0 && faResult.Quantity == 0 {
			faResult.Quantity = len(faResult.Items)
		}
		cb.FunctionAssociations = faResult
	}
	return cb
}

func parseForwardedValues(fv map[string]interface{}) *cloudfrontstore.ForwardedValues {
	if fv == nil {
		return nil
	}
	result := &cloudfrontstore.ForwardedValues{
		QueryString: request.GetBoolParam(fv, "QueryString"),
	}
	if cookies := request.GetMapParam(fv, "Cookies"); cookies != nil {
		cp := &cloudfrontstore.CookiePreferences{
			Forward: request.GetStringParam(cookies, "Forward"),
		}
		if wn := request.GetMapParam(cookies, "WhitelistedNames"); wn != nil {
			cp.WhitelistedNames = &cloudfrontstore.WhitelistedNames{
				Quantity: request.GetIntParam(wn, "Quantity"),
				Items:    parseStringItemList(wn, "Items", "Name"),
			}
			if len(cp.WhitelistedNames.Items) > 0 && cp.WhitelistedNames.Quantity == 0 {
				cp.WhitelistedNames.Quantity = len(cp.WhitelistedNames.Items)
			}
		}
		result.Cookies = cp
	}
	if headers := request.GetMapParam(fv, "Headers"); headers != nil {
		h := &cloudfrontstore.Headers{
			Quantity: request.GetIntParam(headers, "Quantity"),
			Items:    parseStringItemList(headers, "Items", "Name"),
		}
		if len(h.Items) > 0 && h.Quantity == 0 {
			h.Quantity = len(h.Items)
		}
		result.Headers = h
	}
	if qsck := request.GetMapParam(fv, "QueryStringCacheKeys"); qsck != nil {
		k := &cloudfrontstore.QueryStringCacheKeys{
			Quantity: request.GetIntParam(qsck, "Quantity"),
			Items:    parseStringItemList(qsck, "Items", "Name"),
		}
		if len(k.Items) > 0 && k.Quantity == 0 {
			k.Quantity = len(k.Items)
		}
		result.QueryStringCacheKeys = k
	}
	return result
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
	config := parseDistributionConfig(configMap)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createDistributionCore(ctx, store, CreateDistributionInput{
		CallerReference: request.GetStringParam(configMap, "CallerReference"),
		Config:          config,
		ACMRegion:       reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	inProgressCount, _ := store.invalidations.CountInProgress(result.Distribution.ID)
	activeSigners := computeActiveTrustedSigners(result.Distribution)
	activeKeyGroups := computeActiveTrustedKeyGroups(result.Distribution, store)
	return map[string]interface{}{
		"Distribution": formatDistributionResponse(result.Distribution, inProgressCount, activeSigners, activeKeyGroups),
	}, nil
}

// CreateDistributionWithTags creates a new CloudFront distribution with the specified tags.
func (s *CloudFrontService) CreateDistributionWithTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "DistributionConfig")
	if configMap == nil {
		configMap = req.Parameters
	}
	config := parseDistributionConfig(configMap)

	var tags []types.Tag
	tagsMap := request.GetMapParam(req.Parameters, "Tags")
	if tagsMap != nil {
		tags = parseXMLTags(tagsMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createDistributionCore(ctx, store, CreateDistributionInput{
		CallerReference: request.GetStringParam(configMap, "CallerReference"),
		Config:          config,
		Tags:            tags,
		TagsProvided:    len(tags) > 0,
		ACMRegion:       reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}
	inProgressCount, _ := store.invalidations.CountInProgress(result.Distribution.ID)
	activeSigners := computeActiveTrustedSigners(result.Distribution)
	activeKeyGroups := computeActiveTrustedKeyGroups(result.Distribution, store)
	return map[string]interface{}{
		"Distribution": formatDistributionResponse(result.Distribution, inProgressCount, activeSigners, activeKeyGroups),
	}, nil
}

// GetDistribution retrieves a CloudFront distribution by ID.
func (s *CloudFrontService) GetDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}
	inProgressCount, _ := store.invalidations.CountInProgress(distribution.ID)
	activeSigners := computeActiveTrustedSigners(distribution)
	activeKeyGroups := computeActiveTrustedKeyGroups(distribution, store)
	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution, inProgressCount, activeSigners, activeKeyGroups),
	}, nil
}

// GetDistributionConfig retrieves the configuration of a CloudFront distribution.
func (s *CloudFrontService) GetDistributionConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	distribution, err := store.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := resolveListMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	result, err := s.listDistributionsCore(store, ListDistributionsInput{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Distributions))
	for _, d := range result.Distributions {
		items = append(items, formatDistributionSummary(d))
	}

	distList := map[string]interface{}{
		"Marker":      marker,
		"MaxItems":    maxItems,
		"IsTruncated": result.IsTruncated,
		"Quantity":    len(items),
		"Items":       protocol.XMLElements{ElementName: "DistributionSummary", Items: items},
	}
	if result.NextMarker != "" {
		distList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"DistributionList": distList}, nil
}

// UpdateDistribution updates the configuration of a CloudFront distribution.
func (s *CloudFrontService) UpdateDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateDistributionCore(ctx, store, UpdateDistributionInput{
		Id:        request.GetStringParam(req.Parameters, "Id"),
		IfMatch:   getIfMatch(req),
		Config:    parseDistributionConfig(getDistributionConfigMap(req)),
		ACMRegion: reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	distribution := result.Distribution
	inProgressCount, _ := store.invalidations.CountInProgress(distribution.ID)
	activeSigners := computeActiveTrustedSigners(distribution)
	activeKeyGroups := computeActiveTrustedKeyGroups(distribution, store)
	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution, inProgressCount, activeSigners, activeKeyGroups),
	}, nil
}

// DeleteDistribution deletes a CloudFront distribution. The distribution must be disabled first.
func (s *CloudFrontService) DeleteDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteDistributionCore(ctx, store, DeleteDistributionInput{
		Id:        request.GetStringParam(req.Parameters, "Id"),
		IfMatch:   getIfMatch(req),
		ACMRegion: reqCtx.GetRegion(),
	}); err != nil {
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
	maxItems := resolveListMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	// This operation documents a maximum of 100 for MaxItems; larger
	// requests are served with the maximum page size.
	if maxItems > cloudfrontstore.MaxListDistributionsByWebACLIdItems {
		maxItems = cloudfrontstore.MaxListDistributionsByWebACLIdItems
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// This operation models InvalidWebACLId, returned when the specified
	// Web ACL does not exist; the listing must not silently succeed with
	// an empty result for an unknown ACL.
	if s.wafInvoker != nil && !s.wafInvoker.WebACLExists(ctx, webACLId) {
		return nil, awserrors.NewAWSError("InvalidWebACLId", "The specified Web ACL does not exist: "+webACLId, 400)
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

	distList := map[string]interface{}{
		"Marker":      marker,
		"MaxItems":    maxItems,
		"IsTruncated": isTruncated,
		"Quantity":    len(paged),
		"Items":       protocol.XMLElements{ElementName: "DistributionSummary", Items: paged},
	}
	if nextMarker != "" {
		distList["NextMarker"] = nextMarker
	}
	return map[string]interface{}{"DistributionList": distList}, nil
}

// parseXMLTags extracts tags from an XML-style map, handling both array and
// nested map forms produced by XML serialisation.
func parseXMLTags(tagsMap map[string]interface{}) []types.Tag {
	var tags []types.Tag
	parseSingle := func(m map[string]interface{}) {
		tags = append(tags, types.Tag{
			Key:   request.GetStringParam(m, "Key"),
			Value: request.GetStringParam(m, "Value"),
		})
	}
	parseSlice := func(items []interface{}) {
		for _, t := range items {
			if m, ok := t.(map[string]interface{}); ok {
				parseSingle(m)
			}
		}
	}

	if itemsVal := tagsMap["Items"]; itemsVal != nil {
		switch v := itemsVal.(type) {
		case []interface{}:
			parseSlice(v)
		case map[string]interface{}:
			if tagVal, ok := v["Tag"]; ok {
				switch tv := tagVal.(type) {
				case []interface{}:
					parseSlice(tv)
				case map[string]interface{}:
					parseSingle(tv)
				}
			}
		}
	}
	if tagVal := tagsMap["Tag"]; tagVal != nil && len(tags) == 0 {
		switch v := tagVal.(type) {
		case []interface{}:
			parseSlice(v)
		case map[string]interface{}:
			parseSingle(v)
		}
	}
	return tags
}
