package dispatcher

import (
	"errors"
	"log"
	"net/http"
)

func (d *Dispatcher) executeWithRecovery(w http.ResponseWriter, r *http.Request, serviceName, opName string, fn func()) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("PANIC in %s.%s: %v", serviceName, opName, rec)
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
