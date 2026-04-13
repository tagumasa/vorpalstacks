package dispatcher

import (
	"vorpalstacks/internal/common/audit"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

func (d *Dispatcher) recordAudit(serviceName, operation string, reqCtx *request.RequestContext, req *request.ParsedRequest, response interface{}, err error) {
	if !d.auditConfig.Enabled {
		return
	}

	if reqCtx.Principal == "" && req != nil && req.AccessKeyID != "" {
		d.resolvePrincipal(reqCtx, req.AccessKeyID)
	}

	if !reqCtx.HasAuditRecorder() {
		store, storeErr := reqCtx.GetStorage()
		if storeErr != nil {
			return
		}
		ctStore := cloudtrailstore.NewCloudTrailStore(store, reqCtx.GetAccountID(), reqCtx.GetRegion())
		recorder := audit.NewCloudTrailRecorder(ctStore)
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
	globalStorage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return
	}
	iamStore := iamstore.NewIAMStore(globalStorage, reqCtx.GetAccountID())
	accessKey, err := iamStore.AccessKeys().Get(accessKeyID)
	if err != nil {
		return
	}
	if accessKey != nil && accessKey.UserName != "" {
		user, err := iamStore.Users().Get(accessKey.UserName)
		if err != nil || user == nil {
			reqCtx.Principal = accessKey.UserName
		} else {
			reqCtx.Principal = user.UserName
			reqCtx.PrincipalID = user.ID
			reqCtx.PrincipalType = request.PrincipalTypeUser
		}
	}
}
