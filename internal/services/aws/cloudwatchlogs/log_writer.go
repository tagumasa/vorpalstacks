package cloudwatchlogs

import (
	"errors"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// ensureLogGroupAndStream creates the log group and log stream if they do not
// already exist. Returns the resolved logs store, or nil on error.
func (s *LogsService) ensureLogGroupAndStream(region, logGroup, logStream, accountID string) *logsstore.Store {
	logsStore, err := s.getLogsStoreByRegion(region)
	if err != nil {
		logs.Error("Failed to resolve logs store",
			logs.String("logGroup", logGroup),
			logs.String("region", region),
			logs.Err(err))
		return nil
	}

	if _, err := logsStore.GetLogGroup(logGroup); err != nil {
		lg := logsstore.NewLogGroup(logGroup, region, accountID)
		// A concurrent creator winning the race is fine: the group
		// exists, which is all this function has to guarantee. Dropping
		// the shipment here would silently lose the caller's events.
		if createErr := logsStore.CreateLogGroup(lg); createErr != nil && !errors.Is(createErr, logsstore.ErrLogGroupAlreadyExists) {
			logs.Error("Failed to create log group",
				logs.String("logGroup", logGroup),
				logs.Err(createErr))
			return nil
		}
	}

	if !createLogStreamIfAbsent(logsStore, logGroup, logStream) {
		return nil
	}

	return logsStore
}

// createLogStreamIfAbsent creates the log stream in the resolved store when
// it does not already exist. The boolean reports success; a stream that
// already exists is success.
func createLogStreamIfAbsent(logsStore *logsstore.Store, logGroup, logStream string) bool {
	ls := logsstore.NewLogStream(logStream, logGroup)
	if createErr := logsStore.CreateLogStream(ls); createErr != nil {
		if !errors.Is(createErr, logsstore.ErrLogStreamAlreadyExists) {
			logs.Error("Failed to create log stream",
				logs.String("logGroup", logGroup),
				logs.String("logStream", logStream),
				logs.Err(createErr))
			return false
		}
	}
	return true
}

// convertBusLogEntries maps eventbus.LogEntry values to logsstore.LogEntry values.
func convertBusLogEntries(events []eventbus.LogEntry) []logsstore.LogEntry {
	storeEvents := make([]logsstore.LogEntry, len(events))
	for i, e := range events {
		storeEvents[i] = logsstore.LogEntry{Timestamp: e.Timestamp, Message: e.Message}
	}
	return storeEvents
}
