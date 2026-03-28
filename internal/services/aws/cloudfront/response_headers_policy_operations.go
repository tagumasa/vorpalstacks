package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateResponseHeadersPolicy creates a new response headers policy.
func (s *CloudFrontService) CreateResponseHeadersPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if configMap := request.GetMapParam(params, "ResponseHeadersPolicyConfig"); configMap != nil {
		params = configMap
	}

	name := request.GetStringParam(params, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	config := &cloudfrontstore.ResponseHeadersPolicyConfig{
		Name:    name,
		Comment: request.GetStringParam(params, "Comment"),
	}

	if corsMap := request.GetMapParam(params, "CorsConfig"); corsMap != nil {
		config.CorsConfig = parseCorsConfig(corsMap)
	}
	if customMap := request.GetMapParam(params, "CustomHeadersConfig"); customMap != nil {
		config.CustomHeadersConfig = parseCustomHeadersConfig(customMap)
	}
	if removeMap := request.GetMapParam(params, "RemoveHeadersConfig"); removeMap != nil {
		config.RemoveHeadersConfig = parseRemoveHeadersConfig(removeMap)
	}
	if securityMap := request.GetMapParam(params, "SecurityHeadersConfig"); securityMap != nil {
		config.SecurityHeadersConfig = parseSecurityHeadersConfig(securityMap)
	}
	if serverTimingMap := request.GetMapParam(params, "ServerTimingHeadersConfig"); serverTimingMap != nil {
		config.ServerTimingHeadersConfig = parseServerTimingHeadersConfig(serverTimingMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.responseHeadersPolicies.GetByName(name)
	if existing != nil {
		return nil, NewAPIError("ResponseHeadersPolicyAlreadyExists", "Response headers policy with this name already exists", 409)
	}

	policy, err := store.responseHeadersPolicies.Create(config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ResponseHeadersPolicy": buildResponseHeadersPolicyResponse(policy),
		"Location":              policy.ARN,
		"ETag":                  policy.ETag,
	}, nil
}

// GetResponseHeadersPolicy gets a response headers policy by ID.
func (s *CloudFrontService) GetResponseHeadersPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.responseHeadersPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ResponseHeadersPolicy": buildResponseHeadersPolicyResponse(policy),
		"ETag":                  policy.ETag,
	}, nil
}

// GetResponseHeadersPolicyConfig gets the configuration for a response headers policy.
func (s *CloudFrontService) GetResponseHeadersPolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.responseHeadersPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ResponseHeadersPolicyConfig": buildResponseHeadersPolicyConfigResponse(policy),
		"ETag":                        policy.ETag,
	}, nil
}

// UpdateResponseHeadersPolicy updates an existing response headers policy.
func (s *CloudFrontService) UpdateResponseHeadersPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	params := req.Parameters
	if configMap := request.GetMapParam(params, "ResponseHeadersPolicyConfig"); configMap != nil {
		params = configMap
	}

	name := request.GetStringParam(params, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	config := &cloudfrontstore.ResponseHeadersPolicyConfig{
		Name:    name,
		Comment: request.GetStringParam(params, "Comment"),
	}

	if corsMap := request.GetMapParam(params, "CorsConfig"); corsMap != nil {
		config.CorsConfig = parseCorsConfig(corsMap)
	}
	if customMap := request.GetMapParam(params, "CustomHeadersConfig"); customMap != nil {
		config.CustomHeadersConfig = parseCustomHeadersConfig(customMap)
	}
	if removeMap := request.GetMapParam(params, "RemoveHeadersConfig"); removeMap != nil {
		config.RemoveHeadersConfig = parseRemoveHeadersConfig(removeMap)
	}
	if securityMap := request.GetMapParam(params, "SecurityHeadersConfig"); securityMap != nil {
		config.SecurityHeadersConfig = parseSecurityHeadersConfig(securityMap)
	}
	if serverTimingMap := request.GetMapParam(params, "ServerTimingHeadersConfig"); serverTimingMap != nil {
		config.ServerTimingHeadersConfig = parseServerTimingHeadersConfig(serverTimingMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.responseHeadersPolicies.Update(id, config)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ResponseHeadersPolicy": buildResponseHeadersPolicyResponse(policy),
		"ETag":                  policy.ETag,
	}, nil
}

// DeleteResponseHeadersPolicy deletes a response headers policy.
func (s *CloudFrontService) DeleteResponseHeadersPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	err = store.responseHeadersPolicies.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListResponseHeadersPolicies lists all response headers policies.
func (s *CloudFrontService) ListResponseHeadersPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.responseHeadersPolicies.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.ResponseHeadersPolicies))
	for _, policy := range result.ResponseHeadersPolicies {
		items = append(items, map[string]interface{}{
			"ResponseHeadersPolicy": map[string]interface{}{
				"Id":               policy.ID,
				"LastModifiedTime": policy.LastModifiedAt.Format(time.RFC3339),
				"ResponseHeadersPolicyConfig": map[string]interface{}{
					"Name":    policy.Name,
					"Comment": policy.Comment,
				},
			},
			"Type": "custom",
		})
	}

	return map[string]interface{}{
		"ResponseHeadersPolicyList": map[string]interface{}{
			"Marker":      marker,
			"MaxItems":    maxItems,
			"IsTruncated": result.IsTruncated,
			"Quantity":    len(items),
			"NextMarker":  result.NextMarker,
			"Items":       protocol.XMLElements{ElementName: "ResponseHeadersPolicySummary", Items: items},
		},
	}, nil
}

func buildResponseHeadersPolicyResponse(policy *cloudfrontstore.ResponseHeadersPolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                          policy.ID,
		"LastModifiedTime":            policy.LastModifiedAt.Format(time.RFC3339),
		"ResponseHeadersPolicyConfig": buildResponseHeadersPolicyConfigResponse(policy),
	}
}

func buildResponseHeadersPolicyConfigResponse(policy *cloudfrontstore.ResponseHeadersPolicy) map[string]interface{} {
	config := map[string]interface{}{
		"Name":    policy.Name,
		"Comment": policy.Comment,
	}

	if policy.CorsConfig != nil {
		config["CorsConfig"] = buildCorsConfigResponse(policy.CorsConfig)
	}
	if policy.CustomHeadersConfig != nil {
		config["CustomHeadersConfig"] = buildCustomHeadersConfigResponse(policy.CustomHeadersConfig)
	}
	if policy.RemoveHeadersConfig != nil {
		config["RemoveHeadersConfig"] = buildRemoveHeadersConfigResponse(policy.RemoveHeadersConfig)
	}
	if policy.SecurityHeadersConfig != nil {
		config["SecurityHeadersConfig"] = buildSecurityHeadersConfigResponse(policy.SecurityHeadersConfig)
	}
	if policy.ServerTimingHeadersConfig != nil {
		config["ServerTimingHeadersConfig"] = buildServerTimingHeadersConfigResponse(policy.ServerTimingHeadersConfig)
	}

	return config
}

func buildCorsConfigResponse(cors *cloudfrontstore.ResponseHeadersPolicyCorsConfig) map[string]interface{} {
	resp := map[string]interface{}{
		"AccessControlAllowCredentials": cors.AccessControlAllowCredentials,
		"OriginOverride":                cors.OriginOverride,
	}

	if cors.AccessControlAllowHeaders != nil {
		items := make([]interface{}, len(cors.AccessControlAllowHeaders.Items))
		for i, h := range cors.AccessControlAllowHeaders.Items {
			items[i] = h
		}
		resp["AccessControlAllowHeaders"] = map[string]interface{}{
			"Quantity": cors.AccessControlAllowHeaders.Quantity,
			"Items":    protocol.XMLElements{ElementName: "Header", Items: items},
		}
	}
	if cors.AccessControlAllowMethods != nil {
		items := make([]interface{}, len(cors.AccessControlAllowMethods.Items))
		for i, m := range cors.AccessControlAllowMethods.Items {
			items[i] = m
		}
		resp["AccessControlAllowMethods"] = map[string]interface{}{
			"Quantity": cors.AccessControlAllowMethods.Quantity,
			"Items":    protocol.XMLElements{ElementName: "Method", Items: items},
		}
	}
	if cors.AccessControlAllowOrigins != nil {
		items := make([]interface{}, len(cors.AccessControlAllowOrigins.Items))
		for i, o := range cors.AccessControlAllowOrigins.Items {
			items[i] = o
		}
		resp["AccessControlAllowOrigins"] = map[string]interface{}{
			"Quantity": cors.AccessControlAllowOrigins.Quantity,
			"Items":    protocol.XMLElements{ElementName: "Origin", Items: items},
		}
	}
	if cors.AccessControlExposeHeaders != nil {
		items := make([]interface{}, len(cors.AccessControlExposeHeaders.Items))
		for i, h := range cors.AccessControlExposeHeaders.Items {
			items[i] = h
		}
		resp["AccessControlExposeHeaders"] = map[string]interface{}{
			"Quantity": cors.AccessControlExposeHeaders.Quantity,
			"Items":    protocol.XMLElements{ElementName: "ExposeHeader", Items: items},
		}
	}
	if cors.AccessControlMaxAgeSec > 0 {
		resp["AccessControlMaxAgeSec"] = cors.AccessControlMaxAgeSec
	}

	return resp
}

func buildCustomHeadersConfigResponse(custom *cloudfrontstore.ResponseHeadersPolicyCustomHeadersConfig) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(custom.Items))
	for _, h := range custom.Items {
		items = append(items, map[string]interface{}{
			"Header":   h.Header,
			"Override": h.Override,
			"Value":    h.Value,
		})
	}
	itemsI := make([]interface{}, len(items))
	for i, item := range items {
		itemsI[i] = item
	}
	return map[string]interface{}{
		"Quantity": custom.Quantity,
		"Items":    protocol.XMLElements{ElementName: "CustomHeader", Items: itemsI},
	}
}

func buildRemoveHeadersConfigResponse(remove *cloudfrontstore.ResponseHeadersPolicyRemoveHeadersConfig) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(remove.Items))
	for _, h := range remove.Items {
		items = append(items, map[string]interface{}{
			"Header": h.Header,
		})
	}
	itemsI := make([]interface{}, len(items))
	for i, item := range items {
		itemsI[i] = item
	}
	return map[string]interface{}{
		"Quantity": remove.Quantity,
		"Items":    protocol.XMLElements{ElementName: "RemoveHeader", Items: itemsI},
	}
}

func buildSecurityHeadersConfigResponse(security *cloudfrontstore.ResponseHeadersPolicySecurityHeadersConfig) map[string]interface{} {
	resp := map[string]interface{}{}

	if security.ContentSecurityPolicy != nil {
		resp["ContentSecurityPolicy"] = map[string]interface{}{
			"ContentSecurityPolicy": security.ContentSecurityPolicy.ContentSecurityPolicy,
			"Override":              security.ContentSecurityPolicy.Override,
		}
	}
	if security.ContentTypeOptions != nil {
		resp["ContentTypeOptions"] = map[string]interface{}{
			"Override": security.ContentTypeOptions.Override,
		}
	}
	if security.FrameOptions != nil {
		resp["FrameOptions"] = map[string]interface{}{
			"FrameOption": security.FrameOptions.FrameOption,
			"Override":    security.FrameOptions.Override,
		}
	}
	if security.ReferrerPolicy != nil {
		resp["ReferrerPolicy"] = map[string]interface{}{
			"Override":       security.ReferrerPolicy.Override,
			"ReferrerPolicy": security.ReferrerPolicy.ReferrerPolicy,
		}
	}
	if security.StrictTransportSecurity != nil {
		resp["StrictTransportSecurity"] = map[string]interface{}{
			"AccessControlMaxAgeSec": security.StrictTransportSecurity.AccessControlMaxAgeSec,
			"Override":               security.StrictTransportSecurity.Override,
			"IncludeSubdomains":      security.StrictTransportSecurity.IncludeSubdomains,
			"Preload":                security.StrictTransportSecurity.Preload,
		}
	}
	if security.XSSProtection != nil {
		resp["XSSProtection"] = map[string]interface{}{
			"Override":   security.XSSProtection.Override,
			"Protection": security.XSSProtection.Protection,
			"ModeBlock":  security.XSSProtection.ModeBlock,
			"ReportUri":  security.XSSProtection.ReportUri,
		}
	}

	return resp
}

func buildServerTimingHeadersConfigResponse(serverTiming *cloudfrontstore.ResponseHeadersPolicyServerTimingHeadersConfig) map[string]interface{} {
	resp := map[string]interface{}{
		"Enabled": serverTiming.Enabled,
	}
	if serverTiming.SamplingRate > 0 {
		resp["SamplingRate"] = serverTiming.SamplingRate
	}
	return resp
}

func parseCorsConfig(m map[string]interface{}) *cloudfrontstore.ResponseHeadersPolicyCorsConfig {
	cors := &cloudfrontstore.ResponseHeadersPolicyCorsConfig{
		AccessControlAllowCredentials: request.GetBoolParam(m, "AccessControlAllowCredentials"),
		OriginOverride:                request.GetBoolParam(m, "OriginOverride"),
		AccessControlMaxAgeSec:        int32(request.GetIntParam(m, "AccessControlMaxAgeSec")),
	}

	if allowHeaders := request.GetMapParam(m, "AccessControlAllowHeaders"); allowHeaders != nil {
		cors.AccessControlAllowHeaders = &cloudfrontstore.ResponseHeadersPolicyAccessControlAllowHeaders{
			Quantity: int32(request.GetIntParam(allowHeaders, "Quantity")),
			Items:    parseStringSlice(allowHeaders, "Items"),
		}
	}

	if allowMethods := request.GetMapParam(m, "AccessControlAllowMethods"); allowMethods != nil {
		cors.AccessControlAllowMethods = &cloudfrontstore.ResponseHeadersPolicyAccessControlAllowMethods{
			Quantity: int32(request.GetIntParam(allowMethods, "Quantity")),
			Items:    parseStringSlice(allowMethods, "Items"),
		}
	}

	if allowOrigins := request.GetMapParam(m, "AccessControlAllowOrigins"); allowOrigins != nil {
		cors.AccessControlAllowOrigins = &cloudfrontstore.ResponseHeadersPolicyAccessControlAllowOrigins{
			Quantity: int32(request.GetIntParam(allowOrigins, "Quantity")),
			Items:    parseStringSlice(allowOrigins, "Items"),
		}
	}

	if exposeHeaders := request.GetMapParam(m, "AccessControlExposeHeaders"); exposeHeaders != nil {
		cors.AccessControlExposeHeaders = &cloudfrontstore.ResponseHeadersPolicyAccessControlExposeHeaders{
			Quantity: int32(request.GetIntParam(exposeHeaders, "Quantity")),
			Items:    parseStringSlice(exposeHeaders, "Items"),
		}
	}

	return cors
}

func parseCustomHeadersConfig(m map[string]interface{}) *cloudfrontstore.ResponseHeadersPolicyCustomHeadersConfig {
	config := &cloudfrontstore.ResponseHeadersPolicyCustomHeadersConfig{
		Quantity: int32(request.GetIntParam(m, "Quantity")),
	}

	itemsSlice := request.GetSliceParam(m, "Items")
	for _, item := range itemsSlice {
		if itemMap, ok := item.(map[string]interface{}); ok {
			config.Items = append(config.Items, cloudfrontstore.ResponseHeadersPolicyCustomHeader{
				Header:   request.GetStringParam(itemMap, "Header"),
				Override: request.GetBoolParam(itemMap, "Override"),
				Value:    request.GetStringParam(itemMap, "Value"),
			})
		}
	}

	return config
}

func parseRemoveHeadersConfig(m map[string]interface{}) *cloudfrontstore.ResponseHeadersPolicyRemoveHeadersConfig {
	config := &cloudfrontstore.ResponseHeadersPolicyRemoveHeadersConfig{
		Quantity: int32(request.GetIntParam(m, "Quantity")),
	}

	itemsSlice := request.GetSliceParam(m, "Items")
	for _, item := range itemsSlice {
		if itemMap, ok := item.(map[string]interface{}); ok {
			config.Items = append(config.Items, cloudfrontstore.ResponseHeadersPolicyRemoveHeader{
				Header: request.GetStringParam(itemMap, "Header"),
			})
		}
	}

	return config
}

func parseSecurityHeadersConfig(m map[string]interface{}) *cloudfrontstore.ResponseHeadersPolicySecurityHeadersConfig {
	security := &cloudfrontstore.ResponseHeadersPolicySecurityHeadersConfig{}

	if csp := request.GetMapParam(m, "ContentSecurityPolicy"); csp != nil {
		security.ContentSecurityPolicy = &cloudfrontstore.ResponseHeadersPolicyContentSecurityPolicy{
			ContentSecurityPolicy: request.GetStringParam(csp, "ContentSecurityPolicy"),
			Override:              request.GetBoolParam(csp, "Override"),
		}
	}

	if cto := request.GetMapParam(m, "ContentTypeOptions"); cto != nil {
		security.ContentTypeOptions = &cloudfrontstore.ResponseHeadersPolicyContentTypeOptions{
			Override: request.GetBoolParam(cto, "Override"),
		}
	}

	if fo := request.GetMapParam(m, "FrameOptions"); fo != nil {
		security.FrameOptions = &cloudfrontstore.ResponseHeadersPolicyFrameOptions{
			FrameOption: request.GetStringParam(fo, "FrameOption"),
			Override:    request.GetBoolParam(fo, "Override"),
		}
	}

	if rp := request.GetMapParam(m, "ReferrerPolicy"); rp != nil {
		security.ReferrerPolicy = &cloudfrontstore.ResponseHeadersPolicyReferrerPolicy{
			Override:       request.GetBoolParam(rp, "Override"),
			ReferrerPolicy: request.GetStringParam(rp, "ReferrerPolicy"),
		}
	}

	if sts := request.GetMapParam(m, "StrictTransportSecurity"); sts != nil {
		security.StrictTransportSecurity = &cloudfrontstore.ResponseHeadersPolicyStrictTransportSecurity{
			AccessControlMaxAgeSec: int32(request.GetIntParam(sts, "AccessControlMaxAgeSec")),
			Override:               request.GetBoolParam(sts, "Override"),
			IncludeSubdomains:      request.GetBoolParam(sts, "IncludeSubdomains"),
			Preload:                request.GetBoolParam(sts, "Preload"),
		}
	}

	if xss := request.GetMapParam(m, "XSSProtection"); xss != nil {
		security.XSSProtection = &cloudfrontstore.ResponseHeadersPolicyXSSProtection{
			Override:   request.GetBoolParam(xss, "Override"),
			Protection: request.GetBoolParam(xss, "Protection"),
			ModeBlock:  request.GetBoolParam(xss, "ModeBlock"),
			ReportUri:  request.GetStringParam(xss, "ReportUri"),
		}
	}

	return security
}

func parseServerTimingHeadersConfig(m map[string]interface{}) *cloudfrontstore.ResponseHeadersPolicyServerTimingHeadersConfig {
	return &cloudfrontstore.ResponseHeadersPolicyServerTimingHeadersConfig{
		Enabled:      request.GetBoolParam(m, "Enabled"),
		SamplingRate: request.GetFloatParam(m, "SamplingRate"),
	}
}

func parseStringSlice(m map[string]interface{}, key string) []string {
	itemsSlice := request.GetSliceParam(m, key)
	result := make([]string, 0, len(itemsSlice))
	for _, item := range itemsSlice {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
