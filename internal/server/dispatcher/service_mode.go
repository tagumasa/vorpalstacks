package dispatcher

import (
	"vorpalstacks/internal/store/api"
)

func (d *Dispatcher) warmupCache() {
	d.serviceModeCacheMu.Lock()
	defer d.serviceModeCacheMu.Unlock()
	_ = d.configStore.ForEach(func(cfg *api.ServiceConfig) error {
		d.serviceModeCache[cfg.ServiceName] = cfg.Mode
		return nil
	})
}

func (d *Dispatcher) getServiceMode(serviceName string) (api.ServiceMode, bool) {
	d.serviceModeCacheMu.RLock()
	defer d.serviceModeCacheMu.RUnlock()
	m, ok := d.serviceModeCache[serviceName]
	return m, ok
}

func (d *Dispatcher) setServiceMode(serviceName string, mode api.ServiceMode) {
	d.serviceModeCacheMu.Lock()
	defer d.serviceModeCacheMu.Unlock()
	d.serviceModeCache[serviceName] = mode
}
