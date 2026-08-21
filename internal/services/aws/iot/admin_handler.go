package iot

import (
	"context"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/iot"
	iotconnect "vorpalstacks/internal/pb/aws/iot/iotconnect"
)

// AdminHandler implements the IoT admin console gRPC-Web handler.
// It delegates to the shared IoTService Core methods so that both the
// HTTP API and the admin console follow the same validation and
// persistence code path.
type AdminHandler struct {
	iotconnect.UnimplementedIoTServiceHandler
	service *IoTService
}

var _ iotconnect.IoTServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new IoT admin handler.
func NewAdminHandler(svc *IoTService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// --- Thing operations ---

// CreateThing creates a new IoT thing via the admin console.
func (h *AdminHandler) CreateThing(ctx context.Context, req *connect.Request[pb.CreateThingRequest]) (*connect.Response[pb.CreateThingResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var attributes map[string]string
	if ap := req.Msg.GetAttributepayload(); ap != nil {
		attributes = ap.GetAttributes()
	}

	result, err := h.service.createThingCore(stores, CreateThingInput{
		ThingName:        req.Msg.GetThingname(),
		ThingTypeName:    req.Msg.GetThingtypename(),
		BillingGroupName: req.Msg.GetBillinggroupname(),
		Attributes:       attributes,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbCreateThingResponse(result.Thing)), nil
}

// DescribeThing retrieves thing details via the admin console.
func (h *AdminHandler) DescribeThing(ctx context.Context, req *connect.Request[pb.DescribeThingRequest]) (*connect.Response[pb.DescribeThingResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.describeThingCore(stores, req.Msg.GetThingname())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbDescribeThingResponse(result.Thing)), nil
}

// ListThings returns a paginated list of things via the admin console.
func (h *AdminHandler) ListThings(ctx context.Context, req *connect.Request[pb.ListThingsRequest]) (*connect.Response[pb.ListThingsResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listThingsCore(stores, ListThingsInput{
		AttributeName:  req.Msg.GetAttributename(),
		AttributeValue: req.Msg.GetAttributevalue(),
		ThingTypeName:  req.Msg.GetThingtypename(),
		NextToken:      req.Msg.GetNexttoken(),
		MaxItems:       int(req.Msg.GetMaxresults()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	things := make([]*pb.ThingAttribute, 0, len(result.Things))
	for _, t := range result.Things {
		things = append(things, toPbThingAttribute(t))
	}

	return connect.NewResponse(&pb.ListThingsResponse{
		Things:    things,
		Nexttoken: result.NextToken,
	}), nil
}

// DeleteThing removes a thing via the admin console.
func (h *AdminHandler) DeleteThing(ctx context.Context, req *connect.Request[pb.DeleteThingRequest]) (*connect.Response[pb.DeleteThingResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	region := defaults.GetRegionFromHeader(req.Header())

	if err := h.service.deleteThingCore(stores, req.Msg.GetThingname(), h.service.accountID, region); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteThingResponse{}), nil
}

// --- Policy operations ---

// CreatePolicy creates a new IoT policy via the admin console.
func (h *AdminHandler) CreatePolicy(ctx context.Context, req *connect.Request[pb.CreatePolicyRequest]) (*connect.Response[pb.CreatePolicyResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.createPolicyCore(stores, CreatePolicyInput{
		PolicyName:     req.Msg.GetPolicyname(),
		PolicyDocument: req.Msg.GetPolicydocument(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbCreatePolicyResponse(result.Policy)), nil
}

// GetPolicy retrieves policy details via the admin console.
func (h *AdminHandler) GetPolicy(ctx context.Context, req *connect.Request[pb.GetPolicyRequest]) (*connect.Response[pb.GetPolicyResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.getPolicyCore(stores, req.Msg.GetPolicyname())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbGetPolicyResponse(result.Policy)), nil
}

// ListPolicies returns a paginated list of policies via the admin console.
func (h *AdminHandler) ListPolicies(ctx context.Context, req *connect.Request[pb.ListPoliciesRequest]) (*connect.Response[pb.ListPoliciesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxItems := int(req.Msg.GetPagesize())
	result, err := h.service.listPoliciesCore(stores, req.Msg.GetMarker(), maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	policies := make([]*pb.Policy, 0, len(result.Policies))
	for _, p := range result.Policies {
		policies = append(policies, toPbPolicy(p))
	}

	return connect.NewResponse(&pb.ListPoliciesResponse{
		Policies:   policies,
		Nextmarker: result.NextToken,
	}), nil
}

// DeletePolicy removes a policy via the admin console.
func (h *AdminHandler) DeletePolicy(ctx context.Context, req *connect.Request[pb.DeletePolicyRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	region := defaults.GetRegionFromHeader(req.Header())

	if err := h.service.deletePolicyCore(stores, req.Msg.GetPolicyname(), h.service.accountID, region); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Certificate operations ---

// ListCertificates returns a paginated list of certificates via the admin console.
func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxItems := int(req.Msg.GetPagesize())
	result, err := h.service.listCertificatesCore(stores, req.Msg.GetMarker(), maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	certs := make([]*pb.Certificate, 0, len(result.Certificates))
	for _, c := range result.Certificates {
		certs = append(certs, toPbCertificate(c))
	}

	return connect.NewResponse(&pb.ListCertificatesResponse{
		Certificates: certs,
		Nextmarker:   result.NextToken,
	}), nil
}

// DescribeCertificate retrieves certificate details via the admin console.
func (h *AdminHandler) DescribeCertificate(ctx context.Context, req *connect.Request[pb.DescribeCertificateRequest]) (*connect.Response[pb.DescribeCertificateResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.describeCertificateCore(stores, req.Msg.GetCertificateid())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeCertificateResponse{
		Certificatedescription: toPbCertificateDescription(result.Certificate),
	}), nil
}

// UpdateCertificate changes a certificate's status via the admin console.
func (h *AdminHandler) UpdateCertificate(ctx context.Context, req *connect.Request[pb.UpdateCertificateRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.updateCertificateCore(stores, UpdateCertificateInput{
		CertificateID: req.Msg.GetCertificateid(),
		NewStatus:     certStatusFromProto(req.Msg.GetNewstatus()),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteCertificate removes a certificate via the admin console.
func (h *AdminHandler) DeleteCertificate(ctx context.Context, req *connect.Request[pb.DeleteCertificateRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	region := defaults.GetRegionFromHeader(req.Header())

	if err := h.service.deleteCertificateCore(stores, req.Msg.GetCertificateid(), h.service.accountID, region); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Topic Rule operations ---

// ListTopicRules returns a paginated list of topic rules via the admin console.
func (h *AdminHandler) ListTopicRules(ctx context.Context, req *connect.Request[pb.ListTopicRulesRequest]) (*connect.Response[pb.ListTopicRulesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listTopicRulesCore(stores, req.Msg.GetNexttoken(), int(req.Msg.GetMaxresults()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	rules := make([]*pb.TopicRuleListItem, 0, len(result.Rules))
	for _, r := range result.Rules {
		rules = append(rules, toPbTopicRuleListItem(r))
	}

	return connect.NewResponse(&pb.ListTopicRulesResponse{
		Rules:     rules,
		Nexttoken: result.NextToken,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the IoT admin console.
func NewConnectHandler(svc *IoTService) (string, http.Handler) {
	return iotconnect.NewIoTServiceHandler(NewAdminHandler(svc))
}
