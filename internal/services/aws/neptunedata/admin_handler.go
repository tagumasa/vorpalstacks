package neptunedata

import (
	neptunedataconnect "vorpalstacks/internal/pb/aws/neptunedata/neptunedataconnect"
)

type AdminHandler struct {
	neptunedataconnect.UnimplementedNeptunedataServiceHandler
}

var _ neptunedataconnect.NeptunedataServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}
