package scheduler

import (
	"net/http"

	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
)

// AdminHandler provides EventBridge Scheduler service administration functionality.
// It implements the SchedulerServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	schedulerconnect.UnimplementedSchedulerServiceHandler
}

// NewAdminHandler creates a new EventBridge Scheduler AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func NewConnectHandler() (string, http.Handler) {
	return schedulerconnect.NewSchedulerServiceHandler(NewAdminHandler())
}
