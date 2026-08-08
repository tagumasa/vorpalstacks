package scheduler

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// getStore resolves the regional SchedulerStore from the connect request
// headers. This is the sole file in the admin handler layer that imports
// the store package (sole exception: getStore and toPb helpers).
func (h *AdminHandler) getStore(headers http.Header) (*schedulerstore.SchedulerStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}
