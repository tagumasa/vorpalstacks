package sqs

import (
	"time"

	"github.com/google/uuid"

	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// StartMessageMoveTask starts a task to move messages from a source queue to a destination queue.
func (s *SQSStore) StartMessageMoveTask(sourceARN, destARN string, maxMessages int32) (*MessageMoveTask, error) {
	task := &MessageMoveTask{
		TaskId:              uuid.New().String(),
		SourceQueueARN:      sourceARN,
		DestinationQueueARN: destARN,
		Status:              "RUNNING",
		MaxNumberOfMessages: maxMessages,
		StartTime:           time.Now().UTC(),
	}
	if err := s.tasksStore.PutProto(task.TaskId, MessageMoveTaskToProto(task)); err != nil {
		return nil, err
	}
	return task, nil
}

// CancelMessageMoveTask cancels a running message move task.
func (s *SQSStore) CancelMessageMoveTask(taskId string) (*MessageMoveTask, error) {
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskId, &taskPb); err != nil {
		return nil, ErrTaskNotFound
	}
	task := ProtoToMessageMoveTask(&taskPb)
	task.Status = "CANCELLING"
	task.EndTime = time.Now().UTC()
	if err := s.tasksStore.PutProto(taskId, MessageMoveTaskToProto(task)); err != nil {
		return nil, err
	}
	return task, nil
}

// ListMessageMoveTasks lists all message move tasks for a source queue.
func (s *SQSStore) ListMessageMoveTasks(sourceARN string) ([]*MessageMoveTask, error) {
	opts := common.ListOptions{MaxItems: 100}
	result, err := common.ListProto[*pb.MessageMoveTask](s.tasksStore, opts, func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} }, func(t *pb.MessageMoveTask) bool {
		return t.SourceQueueArn == sourceARN
	})
	if err != nil {
		return nil, err
	}
	tasks := make([]*MessageMoveTask, len(result.Items))
	for i, t := range result.Items {
		tasks[i] = ProtoToMessageMoveTask(t)
	}
	return tasks, nil
}

// GetMessageMoveTask retrieves a specific message move task by its ID.
func (s *SQSStore) GetMessageMoveTask(taskId string) (*MessageMoveTask, error) {
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskId, &taskPb); err != nil {
		return nil, ErrTaskNotFound
	}
	return ProtoToMessageMoveTask(&taskPb), nil
}
