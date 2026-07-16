package dispatcher

import (
	"errors"
	"net/http"

	"vorpalstacks/internal/core/logs"
)

func (d *Dispatcher) executeWithRecovery(w http.ResponseWriter, r *http.Request, serviceName, opName string, fn func()) {
	defer func() {
		if rec := recover(); rec != nil {
			logs.Error("PANIC in handler", logs.String("service", serviceName), logs.String("operation", opName), logs.Any("panic", rec))
			d.handleErrorForRequest(w, r, errors.New("panic in handler"))
		}
	}()
	fn()
}
