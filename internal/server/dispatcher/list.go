package dispatcher

import (
	"vorpalstacks/internal/store/api"
)

// ListServices returns a list of all registered service names.
//
// Returns:
//   - []string: A slice of service names
//   - error: Any error that occurred
func (d *Dispatcher) ListServices() ([]string, error) {
	var services []string
	err := d.serviceStore.ForEach(func(s *api.Service) error {
		services = append(services, s.Name)
		return nil
	})
	return services, err
}

// ListOperations returns a list of operations for a given service.
//
// Parameters:
//   - serviceName: The name of the service
//
// Returns:
//   - []*api.Operation: A slice of operations
//   - error: Any error that occurred
func (d *Dispatcher) ListOperations(serviceName string) ([]*api.Operation, error) {
	var operations []*api.Operation
	err := d.operationStore.ForEachByService(serviceName, func(op *api.Operation) error {
		operations = append(operations, op)
		return nil
	})
	return operations, err
}
