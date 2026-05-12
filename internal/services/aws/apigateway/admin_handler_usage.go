package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
)

// CreateApiKey creates a new API key.
func (h *AdminHandler) CreateApiKey(ctx context.Context, req *connect.Request[pb.CreateApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	apiKey := &apigatewaystore.ApiKey{
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		Enabled:     true,
		CustomerId:  req.Msg.Customerid,
		Value:       req.Msg.Value,
	}
	for _, sk := range req.Msg.Stagekeys {
		if sk != nil {
			apiKey.StageKeys = append(apiKey.StageKeys, sk.Restapiid+"/"+sk.Stagename)
		}
	}

	created, err := stores.usage.CreateApiKey(apiKey)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbApiKey(created, true)), nil
}

// GetApiKeys returns all API keys.
func (h *AdminHandler) GetApiKeys(ctx context.Context, req *connect.Request[pb.GetApiKeysRequest]) (*connect.Response[pb.ApiKeys], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 100
	}
	result, err := stores.usage.ListApiKeys(common.ListOptions{
		Marker:   req.Msg.Position,
		MaxItems: limit,
	})
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.ApiKey, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, toPbApiKey(k, false))
	}
	resp := &pb.ApiKeys{Items: items}
	if result.NextMarker != "" {
		resp.Position = result.NextMarker
	}
	return connect.NewResponse(resp), nil
}

// GetApiKey returns a single API key by ID.
func (h *AdminHandler) GetApiKey(ctx context.Context, req *connect.Request[pb.GetApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	if req.Msg.Apikey == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("api_key is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	k, err := stores.usage.GetApiKey(req.Msg.Apikey)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbApiKey(k, req.Msg.Includevalue)), nil
}

// DeleteApiKey removes an API key.
func (h *AdminHandler) DeleteApiKey(ctx context.Context, req *connect.Request[pb.DeleteApiKeyRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Apikey == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("api_key is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.usage.DeleteApiKey(req.Msg.Apikey); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateUsagePlan creates a new usage plan.
func (h *AdminHandler) CreateUsagePlan(ctx context.Context, req *connect.Request[pb.CreateUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	usagePlan := &apigatewaystore.UsagePlan{
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	}
	for _, as := range req.Msg.Apistages {
		if as != nil {
			stage := apigatewaystore.ApiStage{
				ApiId: as.Apiid,
				Stage: as.Stage,
			}
			usagePlan.ApiStages = append(usagePlan.ApiStages, stage)
		}
	}
	if req.Msg.Quota != nil {
		usagePlan.Quota = &apigatewaystore.Quota{
			Limit:  int64(req.Msg.Quota.Limit),
			Offset: int64(req.Msg.Quota.Offset),
			Period: req.Msg.Quota.Period.String(),
		}
	}
	if req.Msg.Throttle != nil {
		usagePlan.Throttle = &apigatewaystore.Throttle{
			BurstLimit: int64(req.Msg.Throttle.Burstlimit),
			RateLimit:  req.Msg.Throttle.Ratelimit,
		}
	}

	created, err := stores.usage.CreateUsagePlan(usagePlan)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbUsagePlan(created)), nil
}

// GetUsagePlans returns all usage plans.
func (h *AdminHandler) GetUsagePlans(ctx context.Context, req *connect.Request[pb.GetUsagePlansRequest]) (*connect.Response[pb.UsagePlans], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 100
	}
	result, err := stores.usage.ListUsagePlans(common.ListOptions{
		Marker:   req.Msg.Position,
		MaxItems: limit,
	})
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.UsagePlan, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, toPbUsagePlan(p))
	}
	resp := &pb.UsagePlans{Items: items}
	if result.NextMarker != "" {
		resp.Position = result.NextMarker
	}
	return connect.NewResponse(resp), nil
}

// GetUsagePlan returns a single usage plan by ID.
func (h *AdminHandler) GetUsagePlan(ctx context.Context, req *connect.Request[pb.GetUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	if req.Msg.Usageplanid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	p, err := stores.usage.GetUsagePlan(req.Msg.Usageplanid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbUsagePlan(p)), nil
}

// DeleteUsagePlan removes a usage plan.
func (h *AdminHandler) DeleteUsagePlan(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Usageplanid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.usage.DeleteUsagePlan(req.Msg.Usageplanid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateUsagePlanKey adds an API key to a usage plan.
func (h *AdminHandler) CreateUsagePlanKey(ctx context.Context, req *connect.Request[pb.CreateUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	if req.Msg.Usageplanid == "" || req.Msg.Keyid == "" || req.Msg.Keytype == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id, key_id, and key_type are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	apiKey, err := stores.usage.GetApiKey(req.Msg.Keyid)
	if err != nil {
		return nil, storeErr(err)
	}

	key := &apigatewaystore.UsagePlanKey{
		Id:    req.Msg.Keyid,
		Type:  req.Msg.Keytype,
		Value: apiKey.Value,
		Name:  apiKey.Name,
	}

	created, err := stores.usage.CreateUsagePlanKey(req.Msg.Usageplanid, key)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbUsagePlanKey(created)), nil
}

// GetUsagePlanKeys returns all keys in a usage plan.
func (h *AdminHandler) GetUsagePlanKeys(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeysRequest]) (*connect.Response[pb.UsagePlanKeys], error) {
	if req.Msg.Usageplanid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 100
	}
	result, err := stores.usage.ListUsagePlanKeys(req.Msg.Usageplanid, common.ListOptions{
		Marker:   req.Msg.Position,
		MaxItems: limit,
	})
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.UsagePlanKey, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, toPbUsagePlanKey(k))
	}
	resp := &pb.UsagePlanKeys{Items: items}
	if result.NextMarker != "" {
		resp.Position = result.NextMarker
	}
	return connect.NewResponse(resp), nil
}

// GetUsagePlanKey returns a single key from a usage plan.
func (h *AdminHandler) GetUsagePlanKey(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	if req.Msg.Usageplanid == "" || req.Msg.Keyid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id and key_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	k, err := stores.usage.GetUsagePlanKey(req.Msg.Usageplanid, req.Msg.Keyid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbUsagePlanKey(k)), nil
}

// DeleteUsagePlanKey removes a key from a usage plan.
func (h *AdminHandler) DeleteUsagePlanKey(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanKeyRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Usageplanid == "" || req.Msg.Keyid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("usage_plan_id and key_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.usage.DeleteUsagePlanKey(req.Msg.Usageplanid, req.Msg.Keyid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
