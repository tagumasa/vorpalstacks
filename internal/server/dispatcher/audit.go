package dispatcher

import (
	"vorpalstacks/internal/common/audit"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
)

// CloudTrailRecorderFactory creates a CloudTrail audit recorder for the given region and account.
// Injected from the composition root to decouple the dispatcher from concrete store types.
type CloudTrailRecorderFactory func(region, accountID string) request.AuditRecorder

func (d *Dispatcher) recordAudit(serviceName, operation string, reqCtx *request.RequestContext, req *request.ParsedRequest, response interface{}, err error) {
	if !d.auditConfig.Enabled {
		return
	}

	if reqCtx.Principal == "" && req != nil && req.AccessKeyID != "" {
		d.resolvePrincipal(reqCtx, req.AccessKeyID)
	}

	if !reqCtx.HasAuditRecorder() {
		if d.cloudTrailRecorderFactory == nil {
			return
		}
		recorder := d.cloudTrailRecorderFactory(reqCtx.GetRegion(), reqCtx.GetAccountID())
		reqCtx.SetAuditRecorder(recorder)
	}

	builder := audit.NewEventBuilder(serviceName, operation, reqCtx, req)
	event := builder.Build(response, err)

	if recorder := reqCtx.GetAuditRecorder(); recorder != nil {
		if auditRecorder, ok := recorder.(audit.Recorder); ok {
			if err := auditRecorder.RecordEvent(event); err != nil {
				logs.Error("Failed to record audit event", logs.Err(err))
			}
		}
	}
}

func (d *Dispatcher) resolvePrincipal(reqCtx *request.RequestContext, accessKeyID string) {
	if d.iamStore == nil {
		return
	}
	accessKey, err := d.iamStore.AccessKeys().Get(accessKeyID)
	if err != nil {
		return
	}
	if accessKey != nil && accessKey.UserName != "" {
		user, err := d.iamStore.Users().Get(accessKey.UserName)
		if err != nil || user == nil {
			reqCtx.Principal = accessKey.UserName
		} else {
			reqCtx.Principal = user.UserName
			reqCtx.PrincipalID = user.ID
			reqCtx.PrincipalType = request.PrincipalTypeUser
		}
	}
}
