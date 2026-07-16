package dispatcher

// RegisterHandler registers a handler for an operation.
//
// Parameters:
//   - operationName: The name of the operation
//   - handler: The handler function to register
func (d *Dispatcher) RegisterHandler(operationName string, handler Handler) {
	d.handlersMu.Lock()
	d.handlers[operationName] = handler
	d.handlerExists[operationName] = true
	d.handlersMu.Unlock()
}

// RegisterHandlerForService registers a handler for a specific service and operation.
//
// Parameters:
//   - serviceName: The name of the service
//   - operationName: The name of the operation
//   - handler: The handler function to register
func (d *Dispatcher) RegisterHandlerForService(serviceName, operationName string, handler Handler) {
	key := serviceName + ":" + operationName
	d.handlersMu.Lock()
	d.handlers[key] = handler
	d.handlerExists[key] = true
	d.handlersMu.Unlock()
}

func (d *Dispatcher) hasHandlerForService(serviceName, operationName string) bool {
	d.handlersMu.RLock()
	defer d.handlersMu.RUnlock()
	if serviceName != "" {
		key := serviceName + ":" + operationName
		if d.handlerExists[key] {
			return true
		}
	}
	return d.handlerExists[operationName]
}

func (d *Dispatcher) getHandler(serviceName, operationName string) (Handler, bool) {
	d.handlersMu.RLock()
	defer d.handlersMu.RUnlock()
	if serviceName != "" {
		key := serviceName + ":" + operationName
		if h, ok := d.handlers[key]; ok {
			return h, true
		}
	}
	h, ok := d.handlers[operationName]
	return h, ok
}
