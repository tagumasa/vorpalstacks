// Package hooks provides hook functionality for extending component behavior.
package hooks

// Event fired before initialization begins.
const (
	EventBeforeInit    = "vorpalstack.before_init"
	EventAfterInit     = "vorpalstack.after_init"
	EventBeforeStart   = "vorpalstack.before_start"
	EventAfterStart    = "vorpalstack.after_start"
	EventBeforeStop    = "vorpalstack.before_stop"
	EventAfterStop     = "vorpalstack.after_stop"
	EventServiceReady  = "vorpalstack.service_ready"
	EventBeforeRestart = "vorpalstack.before_restart"
	EventAfterRestart  = "vorpalstack.after_restart"

	EventStorageReady    = "vorpalstack.storage.ready"
	EventConfigReady     = "vorpalstack.config.ready"
	EventServerReady     = "vorpalstack.server.ready"
	EventServicesStarted = "vorpalstack.services.started"
)

// EventNames is a list of all available hook events.
var EventNames = []string{
	EventBeforeInit,
	EventAfterInit,
	EventBeforeStart,
	EventAfterStart,
	EventBeforeStop,
	EventAfterStop,
	EventServiceReady,
	EventBeforeRestart,
	EventAfterRestart,
	EventStorageReady,
	EventConfigReady,
	EventServerReady,
	EventServicesStarted,
}
