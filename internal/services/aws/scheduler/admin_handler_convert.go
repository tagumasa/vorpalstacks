package scheduler

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// getStore resolves the regional SchedulerStore from the connect request
// headers. This is the sole file in the admin handler layer that imports
// the store package (sole exception: getStore and toPb helpers).
func (h *AdminHandler) getStore(headers http.Header) (*schedulerstore.SchedulerStore, error) {
	region := defaults.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}
