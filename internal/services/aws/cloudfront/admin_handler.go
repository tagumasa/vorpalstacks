package cloudfront

import (
	"context"
	"fmt"
	"net/http"
	"time"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// AdminHandler implements the CloudFront admin console gRPC-Web handler.
type AdminHandler struct {
	cloudfrontconnect.UnimplementedCloudFrontServiceHandler
	store     storage.BasicStorage
	accountId string
}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudFront admin handler with the given storage and account ID.
func NewAdminHandler(store storage.BasicStorage, accountId string) *AdminHandler {
	return &AdminHandler{
		store:     store,
		accountId: accountId,
	}
}

// CreateDistribution creates a new CloudFront distribution with minimal required configuration.
func (h *AdminHandler) CreateDistribution(ctx context.Context, req *connect.Request[pb.CreateDistributionRequest]) (*connect.Response[pb.CreateDistributionResult], error) {
	cfg := req.Msg.GetDistributionconfig()
	if cfg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("distributionconfig is required"))
	}

	comment := cfg.GetComment()
	enabled := cfg.GetEnabled()

	originDomain := ""
	originID := "default-origin"
	if cfg.GetOrigins() != nil && len(cfg.GetOrigins().GetItems()) > 0 {
		firstOrigin := cfg.GetOrigins().GetItems()[0]
		originDomain = firstOrigin.GetDomainname()
		if firstOrigin.GetId() != "" {
			originID = firstOrigin.GetId()
		}
	}
	if originDomain == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("at least one origin with a domain name is required"))
	}

	callerRef := fmt.Sprintf("%d", time.Now().UnixNano())

	storeConfig := &cloudfrontstore.DistributionConfig{
		CallerReference: callerRef,
		Comment:         comment,
		Enabled:         enabled,
		Origins: cloudfrontstore.Origins{
			Quantity: 1,
			Items: []*cloudfrontstore.Origin{
				{
					ID:         originID,
					DomainName: originDomain,
					CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{
						HTTPPort:             80,
						HTTPSPort:            443,
						OriginProtocolPolicy: "https-only",
					},
					ConnectionAttempts: 3,
					ConnectionTimeout:  10,
				},
			},
		},
		DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
			TargetOriginId:       originID,
			ViewerProtocolPolicy: "allow-all",
			AllowedMethods: &cloudfrontstore.AllowedMethods{
				Quantity: 2,
				Items:    []string{"GET", "HEAD"},
			},
			ForwardedValues: &cloudfrontstore.ForwardedValues{
				QueryString: false,
				Cookies:     &cloudfrontstore.CookiePreferences{Forward: "none"},
			},
			MinTTL: 0,
		},
		PriceClass:    "PriceClass_All",
		HttpVersion:   "http2and3",
		IsIPV6Enabled: true,
		ViewerCertificate: &cloudfrontstore.ViewerCertificate{
			CloudFrontDefaultCertificate: true,
			CertificateSource:            "cloudfront",
		},
		Restrictions: &cloudfrontstore.Restrictions{
			GeoRestriction: cloudfrontstore.GeoRestriction{
				RestrictionType: "none",
			},
		},
	}

	distStore := cloudfrontstore.NewDistributionStore(h.store, h.accountId)

	dist, err := distStore.Create(callerRef, storeConfig)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDistributionResult{
		Distribution: &pb.Distribution{
			Id:               dist.ID,
			Arn:              dist.ARN,
			Status:           dist.Status,
			Domainname:       dist.DomainName,
			Lastmodifiedtime: dist.LastModifiedAt.Format(timeutils.ISO8601UTCFormat),
		},
	}), nil
}

// ListDistributions returns all CloudFront distributions visible to the admin console.
func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	distStore := cloudfrontstore.NewDistributionStore(h.store, h.accountId)

	maxItems := int(req.Msg.Maxitems)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := distStore.List(req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var items []*pb.DistributionSummary
	for _, d := range result.Distributions {
		summary := &pb.DistributionSummary{
			Id:         d.ID,
			Arn:        d.ARN,
			Status:     d.Status,
			Enabled:    d.Enabled,
			Staging:    d.Staging,
			Etag:       d.ETag,
			Comment:    d.DistributionConfig.Comment,
			Domainname: d.DomainName,
		}
		if !d.LastModifiedAt.IsZero() {
			summary.Lastmodifiedtime = d.LastModifiedAt.Format(timeutils.ISO8601UTCFormat)
		}
		items = append(items, summary)
	}

	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    int32(len(items)),
			Items:       items,
			Istruncated: result.IsTruncated,
			Nextmarker:  result.NextMarker,
			Marker:      req.Msg.Marker,
			Maxitems:    int32(maxItems),
		},
	}), nil
}

// DeleteDistribution deletes a CloudFront distribution via the admin console.
func (h *AdminHandler) DeleteDistribution(ctx context.Context, req *connect.Request[pb.DeleteDistributionRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	distStore := cloudfrontstore.NewDistributionStore(h.store, h.accountId)

	if err := distStore.Delete(req.Msg.Id); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cloudfront admin console.
func NewConnectHandler(store storage.BasicStorage, accountID string) (string, http.Handler) {
	return cloudfrontconnect.NewCloudFrontServiceHandler(NewAdminHandler(store, accountID))
}
