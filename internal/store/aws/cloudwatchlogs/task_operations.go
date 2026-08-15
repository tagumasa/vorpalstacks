package cloudwatchlogs

import (
	"encoding/json"
	"time"
)

// --- Export Task ---

func (s *Store) PutExportTask(task *ExportTask) error {
	return s.Put(s.exportTaskKey(task.TaskId), task)
}

func (s *Store) GetExportTask(taskId string) (*ExportTask, error) {
	var t ExportTask
	if err := s.Get(s.exportTaskKey(taskId), &t); err != nil {
		return nil, ErrResourceNotFound
	}
	return &t, nil
}

func (s *Store) DeleteExportTask(taskId string) error {
	return s.Delete(s.exportTaskKey(taskId))
}

func (s *Store) ListExportTasks(statusCode string) ([]*ExportTask, error) {
	var tasks []*ExportTask
	if err := s.ScanPrefix("export-task:", func(key string, value []byte) error {
		var t ExportTask
		if err := json.Unmarshal(value, &t); err != nil {
			return nil
		}
		if statusCode == "" || t.Status == statusCode {
			tasks = append(tasks, &t)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}

// --- Import Task ---

func (s *Store) PutImportTask(task *ImportTask) error {
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.importTaskKey(task.ImportId), task)
}

func (s *Store) GetImportTask(importId string) (*ImportTask, error) {
	var t ImportTask
	if err := s.Get(s.importTaskKey(importId), &t); err != nil {
		return nil, ErrResourceNotFound
	}
	return &t, nil
}

func (s *Store) ListImportTasks(importId, status, sourceArn string) ([]*ImportTask, error) {
	var tasks []*ImportTask
	if err := s.ScanPrefix("import-task:", func(key string, value []byte) error {
		var t ImportTask
		if err := json.Unmarshal(value, &t); err != nil {
			return nil
		}
		if importId != "" && t.ImportId != importId {
			return nil
		}
		if status != "" && t.ImportStatus != status {
			return nil
		}
		if sourceArn != "" && t.ImportSourceArn != sourceArn {
			return nil
		}
		tasks = append(tasks, &t)
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}
