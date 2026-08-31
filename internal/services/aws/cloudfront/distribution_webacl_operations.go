package cloudfront

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
)

// AssociateDistributionWebACL associates a WAF Web ACL with a
// distribution.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_AssociateDistributionWebACL.html
func (s *CloudFrontService) AssociateDistributionWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	distribution, err := s.associateDistributionWebACLCore(ctx, stores, AssociateDistributionWebACLInput{
		Id:        request.GetStringParam(req.Parameters, "Id"),
		WebACLArn: request.GetStringParam(req.Parameters, "WebACLArn"),
		IfMatch:   getIfMatch(req),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AssociateDistributionWebACLResult": map[string]interface{}{
			"Id":        distribution.ID,
			"WebACLArn": request.GetStringParam(req.Parameters, "WebACLArn"),
		},
		"ETag": distribution.ETag,
	}, nil
}

// DisassociateDistributionWebACL removes the WAF Web ACL association from
// a distribution.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DisassociateDistributionWebACL.html
func (s *CloudFrontService) DisassociateDistributionWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	distribution, err := s.disassociateDistributionWebACLCore(ctx, stores,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DisassociateDistributionWebACLResult": map[string]interface{}{
			"Id": distribution.ID,
		},
		"ETag": distribution.ETag,
	}, nil
}

// CopyDistribution creates a staging copy of an existing distribution.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CopyDistribution.html
func (s *CloudFrontService) CopyDistribution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := CopyDistributionInput{
		PrimaryId:       request.GetStringParam(req.Parameters, "Id"),
		CallerReference: request.GetStringParam(req.Parameters, "CallerReference"),
		IfMatch:         getIfMatch(req),
		ACMRegion:       reqCtx.GetRegion(),
	}
	if _, ok := req.Parameters["Enabled"]; ok {
		in.Enabled = request.GetBoolParam(req.Parameters, "Enabled")
		in.EnabledProvided = true
	}
	// Staging travels as an HTTP header, not a body or query parameter.
	if stagingHdr := req.Headers.Get("Staging"); stagingHdr != "" {
		in.Staging = strings.EqualFold(strings.TrimSpace(stagingHdr), "true")
		in.StagingProvided = true
	}

	distribution, err := s.copyDistributionCore(ctx, stores, in)
	if err != nil {
		return nil, err
	}

	detail := s.distributionDetailCore(stores, distribution)
	return map[string]interface{}{
		"Distribution": formatDistributionResponse(distribution, detail.InProgressInvalidations, detail.ActiveSigners, detail.ActiveKeyGroups),
		"ETag":         distribution.ETag,
		"Location":     distribution.DomainName,
	}, nil
}
