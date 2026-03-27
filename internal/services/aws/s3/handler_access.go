package s3

import (
	"net/http"

	"vorpalstacks/internal/services/aws/common/request"
)

func (h *S3Handler) checkAccess(ctx *request.RequestContext, r *http.Request, stores *s3Stores, action, bucket, key string) error {
	secureTransport := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	check := &AccessCheck{
		Principal:       ctx.Principal,
		PrincipalID:     ctx.PrincipalID,
		PrincipalType:   ctx.PrincipalType,
		Action:          action,
		Resource:        buildResource(h.svc.accountID, h.region, bucket, key),
		Bucket:          bucket,
		Key:             key,
		SourceIP:        ctx.SourceIP,
		Referer:         r.Referer(),
		SecureTransport: secureTransport,
	}

	if key != "" {
		return h.svc.accessController.CheckObjectAccess(ctx, stores, check)
	}
	return h.svc.accessController.CheckAccess(ctx, stores, check)
}
