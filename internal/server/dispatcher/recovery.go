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

func (d *Dispatcher) executeWithRecoveryAndResult(w http.ResponseWriter, r *http.Request, serviceName, opName string, fn func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	d.executeWithRecovery(w, r, serviceName, opName, func() {
		result, err = fn()
	})
	return result, err
}
