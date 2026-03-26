package dispatcher

import (
	"log"

	cloudtrailaudit "vorpalstacks/internal/services/aws/cloudtrail/audit"
	"vorpalstacks/internal/services/aws/common/audit"
	"vorpalstacks/internal/services/aws/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

func (d *Dispatcher) recordAudit(serviceName, operation string, reqCtx *request.RequestContext, req *request.ParsedRequest, response interface{}, err error) {
	if !d.auditConfig.Enabled {
		return
	}

	if !reqCtx.HasAuditRecorder() {
		store, storeErr := reqCtx.GetStorage()
		if storeErr != nil {
			return
		}
		cloudtrailStore := cloudtrailstore.NewCloudTrailStore(store, reqCtx.GetAccountID(), reqCtx.GetRegion())
		recorder := cloudtrailaudit.NewCloudTrailRecorder(cloudtrailStore)
		reqCtx.SetAuditRecorder(recorder)
	}

	builder := audit.NewEventBuilder(serviceName, operation, reqCtx, req)
	event := builder.Build(response, err)

	if recorder := reqCtx.GetAuditRecorder(); recorder != nil {
		if auditRecorder, ok := recorder.(audit.Recorder); ok {
			if err := auditRecorder.RecordEvent(event); err != nil {
				log.Printf("Failed to record audit event: %v", err)
			}
		}
	}
}
