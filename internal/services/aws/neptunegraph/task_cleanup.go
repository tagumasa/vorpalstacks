package neptunegraph

import (
	"time"

	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
)

var terminalTaskStatuses = map[string]bool{
	"SUCCEEDED": true,
	"FAILED":    true,
	"CANCELLED": true,
}

func (s *NeptuneGraphService) cleanupExpiredTasks(store *ngstore.NeptuneGraphStore) {
	ttl := taskTTL
	if s.testMode {
		ttl = taskTTLTestMode
	}
	now := time.Now()

	s.cleanupImportTasks(store, ttl, now)
	s.cleanupExportTasks(store, ttl, now)
}

func (s *NeptuneGraphService) cleanupImportTasks(store *ngstore.NeptuneGraphStore, ttl time.Duration, now time.Time) {
	var marker string
	for {
		tasks, nextToken, _, err := store.ListImportTasks(storecommon.ListOptions{
			MaxItems: 100,
			Marker:   marker,
		})
		if err != nil {
			logs.Warn("neptunegraph: failed to list import tasks for cleanup", logs.Err(err))
			return
		}
		for _, t := range tasks {
			if !terminalTaskStatuses[t.Status] {
				continue
			}
			if t.StartTime != nil && now.Sub(*t.StartTime) >= ttl {
				if err := store.DeleteImportTask(t.TaskId); err != nil {
					logs.Warn("neptunegraph: failed to delete expired import task",
						logs.String("taskId", t.TaskId), logs.Err(err))
				}
			}
		}
		if nextToken == "" {
			break
		}
		marker = nextToken
	}
}

func (s *NeptuneGraphService) cleanupExportTasks(store *ngstore.NeptuneGraphStore, ttl time.Duration, now time.Time) {
	var marker string
	for {
		tasks, nextToken, _, err := store.ListExportTasks(storecommon.ListOptions{
			MaxItems: 100,
			Marker:   marker,
		}, "")
		if err != nil {
			logs.Warn("neptunegraph: failed to list export tasks for cleanup", logs.Err(err))
			return
		}
		for _, t := range tasks {
			if !terminalTaskStatuses[t.Status] {
				continue
			}
			if t.StartTime != nil && now.Sub(*t.StartTime) >= ttl {
				if err := store.DeleteExportTask(t.TaskId); err != nil {
					logs.Warn("neptunegraph: failed to delete expired export task",
						logs.String("taskId", t.TaskId), logs.Err(err))
				}
			}
		}
		if nextToken == "" {
			break
		}
		marker = nextToken
	}
}
