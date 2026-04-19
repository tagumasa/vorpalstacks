package cloudwatchlogs

import (
	"time"

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
		if createErr := logsStore.CreateLogGroup(lg); createErr != nil {
			logs.Error("Failed to create log group",
				logs.String("logGroup", logGroup),
				logs.Err(createErr))
			return nil
		}
	}

	ls := logsstore.NewLogStream(logStream, logGroup)
	if createErr := logsStore.CreateLogStream(ls); createErr != nil {
		logs.Error("Failed to create log stream",
			logs.String("logGroup", logGroup),
			logs.String("logStream", logStream),
			logs.Err(createErr))
		return nil
	}

	return logsStore
}

// writeLogEvents converts bus LogEntry values to store LogEntry values and
// writes them via PutLogEvents. Returns false on failure.
func (s *LogsService) writeLogEvents(logsStore *logsstore.Store, logGroup, logStream string, entries []logsstore.LogEntry) bool {
	if _, err := logsStore.PutLogEvents(logGroup, logStream, entries); err != nil {
		logs.Error("Failed to write log events",
			logs.String("logGroup", logGroup),
			logs.String("logStream", logStream),
			logs.Err(err))
		return false
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

// writeSingleLogMessage is a convenience wrapper that writes one log entry
// with the current timestamp.
func (s *LogsService) writeSingleLogMessage(region, logGroup, logStream, accountID, message string) {
	logsStore := s.ensureLogGroupAndStream(region, logGroup, logStream, accountID)
	if logsStore == nil {
		return
	}
	if !s.writeLogEvents(logsStore, logGroup, logStream, []logsstore.LogEntry{
		{Timestamp: time.Now().UnixMilli(), Message: message},
	}) {
		logs.Error("Failed to write log event",
			logs.String("logGroup", logGroup), logs.String("logStream", logStream))
	}
}
