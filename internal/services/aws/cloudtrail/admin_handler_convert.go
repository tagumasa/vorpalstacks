package cloudtrail

import (
	"net/http"

	"vorpalstacks/internal/common/defaults"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// This file is the admin plane's single store-typed bridge: it resolves the
// per-region store for the admin handlers, spelling the store interface in
// its full form so the rest of the package keeps one spelling and the admin
// handler files stay free of store imports.

// getStoreFromHeader resolves the region-scoped CloudTrail store from the
// admin request's region header.
func (h *AdminHandler) getStoreFromHeader(header http.Header) (cloudtrailstore.CloudTrailStoreInterface, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}
