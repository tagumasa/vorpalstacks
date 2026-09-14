package sfn

import "time"

// This file holds the activity model: the activity resource and the
// activity-task records the worker poll and the task-result waits use.

// Activity represents an AWS Step Functions activity used by Task states.
type Activity struct {
	ActivityArn             string                   `json:"activityArn"`
	Name                    string                   `json:"name"`
	CreationDate            time.Time                `json:"creationDate"`
	EncryptionConfiguration *EncryptionConfiguration `json:"encryptionConfiguration,omitempty"`
}

// ActivityTask represents a single activity task dispatched to a worker.
type ActivityTask struct {
	TaskToken       string    `json:"taskToken"`
	ActivityArn     string    `json:"activityArn"`
	ExecutionArn    string    `json:"executionArn"`
	Input           string    `json:"input"`
	Status          string    `json:"status"`
	Output          string    `json:"output,omitempty"`
	Error           string    `json:"error,omitempty"`
	Cause           string    `json:"cause,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	CompletedAt     time.Time `json:"completedAt,omitempty"`
	LastHeartbeatAt time.Time `json:"lastHeartbeatAt,omitempty"`
	WorkerName      string    `json:"workerName,omitempty"`
}

// ActivityTaskResult contains the outcome of an activity task. Error keeps
// the worker's error name and Cause its cause as separate members, matching
// the SendTaskFailure pair; WorkerName travels with the outcome so the
// executor can record the ActivityStarted event when the result arrives.
type ActivityTaskResult struct {
	TaskToken  string
	Output     string
	Error      error
	Cause      string
	WorkerName string
}
