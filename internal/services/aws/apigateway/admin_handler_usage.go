package apigateway

import (
	"context"
	"google.golang.org/protobuf/proto"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// CreateApiKey creates a new API key.
func (h *AdminHandler) CreateApiKey(ctx context.Context, req *connect.Request[pb.CreateApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &ApiKeyInput{
		Name:               req.Msg.GetName(),
		Description:        req.Msg.GetDescription(),
		CustomerId:         req.Msg.GetCustomerid(),
		Value:              req.Msg.GetValue(),
		Enabled:            req.Msg.Enabled,
		GenerateDistinctId: req.Msg.Generatedistinctid,
	}

	if len(req.Msg.Tags) > 0 {
		in.Tags = tagutil.MapToTags(req.Msg.Tags)
	}

	in.StageKeys = stageKeysFromPb(req.Msg.Stagekeys)

	created, err := h.service.createApiKeyCore(stores, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbApiKey(created, true)), nil
}

// GetApiKeys returns all API keys.
func (h *AdminHandler) GetApiKeys(ctx context.Context, req *connect.Request[pb.GetApiKeysRequest]) (*connect.Response[pb.ApiKeys], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.listApiKeysCore(stores, int(req.Msg.GetLimit()), req.Msg.GetPosition(), "", "")
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	items := make([]*pb.ApiKey, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, toPbApiKey(k, false))
	}
	resp := &pb.ApiKeys{Items: items}
	if result.NextMarker != "" {
		resp.Position = proto.String(result.NextMarker)
	}
	return connect.NewResponse(resp), nil
}

// GetApiKey returns a single API key by ID.
func (h *AdminHandler) GetApiKey(ctx context.Context, req *connect.Request[pb.GetApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	k, err := h.service.getApiKeyCore(stores, req.Msg.Apikey)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbApiKey(k, req.Msg.GetIncludevalue())), nil
}

// DeleteApiKey removes an API key.
func (h *AdminHandler) DeleteApiKey(ctx context.Context, req *connect.Request[pb.DeleteApiKeyRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteApiKeyCore(stores, req.Msg.Apikey); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateUsagePlan creates a new usage plan.
func (h *AdminHandler) CreateUsagePlan(ctx context.Context, req *connect.Request[pb.CreateUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &UsagePlanInput{
		Name:        req.Msg.GetName(),
		Description: req.Msg.GetDescription(),
	}
	for _, as := range req.Msg.Apistages {
		if as != nil {
			in.ApiStages = append(in.ApiStages, apiStageInputFromPb(as))
		}
	}
	if req.Msg.Quota != nil {
		in.Quota = &QuotaInput{
			Limit:  int64(req.Msg.Quota.GetLimit()),
			Offset: int64(req.Msg.Quota.GetOffset()),
			Period: toQuotaPeriodString(req.Msg.Quota.Period),
		}
	}
	if req.Msg.Throttle != nil {
		in.Throttle = &ThrottleInput{
			BurstLimit: int64(req.Msg.Throttle.GetBurstlimit()),
			RateLimit:  req.Msg.Throttle.GetRatelimit(),
		}
	}

	created, err := h.service.createUsagePlanCore(stores, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbUsagePlan(created)), nil
}

// GetUsagePlans returns all usage plans.
func (h *AdminHandler) GetUsagePlans(ctx context.Context, req *connect.Request[pb.GetUsagePlansRequest]) (*connect.Response[pb.UsagePlans], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.listUsagePlansCore(stores, int(req.Msg.GetLimit()), req.Msg.GetPosition(), "")
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	items := make([]*pb.UsagePlan, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, toPbUsagePlan(p))
	}
	resp := &pb.UsagePlans{Items: items}
	if result.NextMarker != "" {
		resp.Position = proto.String(result.NextMarker)
	}
	return connect.NewResponse(resp), nil
}

// GetUsagePlan returns a single usage plan by ID.
func (h *AdminHandler) GetUsagePlan(ctx context.Context, req *connect.Request[pb.GetUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	p, err := h.service.getUsagePlanCore(stores, req.Msg.Usageplanid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbUsagePlan(p)), nil
}

// DeleteUsagePlan removes a usage plan.
func (h *AdminHandler) DeleteUsagePlan(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteUsagePlanCore(stores, req.Msg.Usageplanid); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateUsagePlanKey adds an API key to a usage plan.
func (h *AdminHandler) CreateUsagePlanKey(ctx context.Context, req *connect.Request[pb.CreateUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &UsagePlanKeyInput{
		KeyId:   req.Msg.Keyid,
		KeyType: req.Msg.Keytype,
	}
	created, err := h.service.createUsagePlanKeyCore(stores, req.Msg.Usageplanid, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbUsagePlanKey(created)), nil
}

// GetUsagePlanKeys returns all keys in a usage plan.
func (h *AdminHandler) GetUsagePlanKeys(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeysRequest]) (*connect.Response[pb.UsagePlanKeys], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.listUsagePlanKeysCore(stores, req.Msg.Usageplanid, int(req.Msg.GetLimit()), req.Msg.GetPosition())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	items := make([]*pb.UsagePlanKey, 0, len(result.Items))
	for _, k := range result.Items {
		items = append(items, toPbUsagePlanKey(k))
	}
	resp := &pb.UsagePlanKeys{Items: items}
	if result.NextMarker != "" {
		resp.Position = proto.String(result.NextMarker)
	}
	return connect.NewResponse(resp), nil
}

// GetUsagePlanKey returns a single key from a usage plan.
func (h *AdminHandler) GetUsagePlanKey(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	k, err := h.service.getUsagePlanKeyCore(stores, req.Msg.Usageplanid, req.Msg.Keyid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbUsagePlanKey(k)), nil
}

// DeleteUsagePlanKey removes a key from a usage plan.
func (h *AdminHandler) DeleteUsagePlanKey(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanKeyRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteUsagePlanKeyCore(stores, req.Msg.Usageplanid, req.Msg.Keyid); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
