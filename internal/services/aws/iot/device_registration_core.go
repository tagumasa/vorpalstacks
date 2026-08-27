package iot

import (
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/services/aws/iot/ca"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the bulk device-registration family (RegisterThing
// and the thing-registration task lifecycle). Handlers on both protocol
// planes are thin adapters; validation and persistence live here only.

// RegisterThingInput carries the register-thing parameters: the
// provisioning template body (required, TemplateBody shape) and the
// caller-supplied template parameters.
type RegisterThingInput struct {
	TemplateBody string
	Parameters   map[string]string
}

// RegisterThingResult reports the generated resources: the ARNs keyed by
// the template's logical resource names (the model's resourceArns map) and
// the PEM of the provisioned certificate when the template declares one.
type RegisterThingResult struct {
	ResourceArns   map[string]string
	CertificatePem string
}

// registerThingCore provisions every resource declared by the provisioning
// template (thing, certificate, policy) through the shared template
// engine. Partial failures are not rolled back.
func (s *IoTService) registerThingCore(store iotstore.IotStoreInterface, authority *ca.CertificateAuthority, in RegisterThingInput) (*RegisterThingResult, error) {
	tpl, err := parseProvisioningTemplate(in.TemplateBody)
	if err != nil {
		return nil, err
	}
	outcome, err := s.provisionFromTemplate(store, authority, tpl, in.Parameters)
	if err != nil {
		return nil, err
	}
	return &RegisterThingResult{
		ResourceArns:   outcome.resourceArns,
		CertificatePem: outcome.certificatePem,
	}, nil
}

// thingRegistrationTaskKey prefixes every registration-task record in the
// generic-KV namespace.
const thingRegistrationTaskKey = "registrationTask/"

// Registration task status values, matching the API_DescribeThingRegistrationTask
// status enumeration: InProgress | Completed | Failed | Cancelled | Cancelling.
const (
	taskStatusInProgress = "InProgress"
	taskStatusCompleted  = "Completed"
	taskStatusFailed     = "Failed"
	taskStatusCancelled  = "Cancelled"
	taskStatusCancelling = "Cancelling"
)

// StartThingRegistrationTaskInput carries the bulk-registration task
// parameters; all four members are required by the model.
type StartThingRegistrationTaskInput struct {
	TemplateBody    string
	InputFileBucket string
	InputFileKey    string
	RoleArn         string
	Region          string
}

// validateThingRegistrationTaskInput enforces the documented member
// contract: all four members required, templateBody within the shared
// TemplateBody length bound, the input-file bucket/key within their
// length+pattern bounds, roleArn within its length bounds, and a
// structurally valid provisioning template (the task provisions things
// from it).
func validateThingRegistrationTaskInput(in StartThingRegistrationTaskInput) (*provisioningTemplate, error) {
	if in.TemplateBody == "" || in.InputFileBucket == "" || in.InputFileKey == "" || in.RoleArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	if len(in.TemplateBody) > MaxTemplateBodyLength ||
		len(in.InputFileBucket) < MinInputFileBucketLength || len(in.InputFileBucket) > MaxInputFileBucketLength ||
		!inputFileBucketPattern.MatchString(in.InputFileBucket) ||
		len(in.InputFileKey) < MinInputFileKeyLength || len(in.InputFileKey) > MaxInputFileKeyLength ||
		!inputFileKeyPattern.MatchString(in.InputFileKey) ||
		len(in.RoleArn) < MinRoleArnLength || len(in.RoleArn) > MaxRoleArnLength {
		return nil, iotstore.ErrInvalidRequest
	}
	return parseProvisioningTemplate(in.TemplateBody)
}

// startThingRegistrationTaskCore records a new registration task in the
// InProgress state and launches the asynchronous worker that reads the S3
// input file and provisions one device per line.
func (s *IoTService) startThingRegistrationTaskCore(store iotstore.IotStoreInterface, in StartThingRegistrationTaskInput) (string, error) {
	tpl, err := validateThingRegistrationTaskInput(in)
	if err != nil {
		return "", err
	}
	taskID := uuid.New().String()
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"taskId":             taskID,
		"templateBody":       in.TemplateBody,
		"inputFileBucket":    in.InputFileBucket,
		"inputFileKey":       in.InputFileKey,
		"roleArn":            in.RoleArn,
		"status":             taskStatusInProgress,
		"creationDate":       now,
		"lastModifiedDate":   now,
		"successCount":       0,
		"failureCount":       0,
		"percentageProgress": 0,
	}
	if err := store.PutGeneric(thingRegistrationTaskKey+taskID, rec); err != nil {
		return "", err
	}
	s.startThingRegistrationWorker(store, in.Region, taskID, in.InputFileBucket, in.InputFileKey, tpl)
	return taskID, nil
}

// stopThingRegistrationTaskCore requests cancellation of a running task:
// InProgress moves to Cancelling and the worker finalises Cancelled; a
// task already in a terminal state is rejected as invalid. A missing task
// ID is rejected with the task-not-found error.
func (s *IoTService) stopThingRegistrationTaskCore(store iotstore.IotStoreInterface, taskID string) error {
	if taskID == "" {
		return iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(thingRegistrationTaskKey+taskID, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrThingRegistrationTaskNotFound
	}
	switch rec["status"] {
	case taskStatusInProgress, taskStatusCancelling:
		rec["status"] = taskStatusCancelling
		rec["lastModifiedDate"] = time.Now().UTC().Unix()
		return store.PutGeneric(thingRegistrationTaskKey+taskID, rec)
	default:
		return iotstore.ErrInvalidRequest
	}
}

// ThingRegistrationTaskResult is the describe view of a registration task
// record; every member maps one API_DescribeThingRegistrationTask response
// element.
type ThingRegistrationTaskResult struct {
	TaskID             string
	Status             string
	CreationDate       interface{}
	LastModifiedDate   interface{}
	TemplateBody       string
	InputFileBucket    string
	InputFileKey       string
	RoleArn            string
	Message            string
	SuccessCount       int
	FailureCount       int
	PercentageProgress int
	Exists             bool
}

// describeThingRegistrationTaskCore loads a registration task record.
// Exists is false (nil error) when the task ID is unknown.
func (s *IoTService) describeThingRegistrationTaskCore(store iotstore.IotStoreInterface, taskID string) (*ThingRegistrationTaskResult, error) {
	if taskID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(thingRegistrationTaskKey+taskID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrThingRegistrationTaskNotFound
	}
	return thingRegistrationTaskFromRecord(taskID, rec), nil
}

func thingRegistrationTaskFromRecord(taskID string, rec map[string]interface{}) *ThingRegistrationTaskResult {
	asInt := func(v interface{}) int {
		switch x := v.(type) {
		case int:
			return x
		case int64:
			return int(x)
		case float64:
			return int(x)
		default:
			return 0
		}
	}
	asString := func(v interface{}) string {
		s, _ := v.(string)
		return s
	}
	return &ThingRegistrationTaskResult{
		TaskID:             taskID,
		Status:             asString(rec["status"]),
		CreationDate:       rec["creationDate"],
		LastModifiedDate:   rec["lastModifiedDate"],
		TemplateBody:       asString(rec["templateBody"]),
		InputFileBucket:    asString(rec["inputFileBucket"]),
		InputFileKey:       asString(rec["inputFileKey"]),
		RoleArn:            asString(rec["roleArn"]),
		Message:            asString(rec["message"]),
		SuccessCount:       asInt(rec["successCount"]),
		FailureCount:       asInt(rec["failureCount"]),
		PercentageProgress: asInt(rec["percentageProgress"]),
		Exists:             true,
	}
}

// thingRegistrationReportKeyCore resolves the stored object key of a
// task's report of the requested type (ERRORS or RESULTS). The key is
// empty until the task reaches a terminal state and its reports are
// written; an unknown task ID is the documented not-found error.
func (s *IoTService) thingRegistrationReportKeyCore(store iotstore.IotStoreInterface, taskID, reportType string) (string, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(thingRegistrationTaskKey+taskID, &rec)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", iotstore.ErrThingRegistrationTaskNotFound
	}
	var key string
	switch reportType {
	case "RESULTS":
		key, _ = rec["reportResultsKey"].(string)
	case "ERRORS":
		key, _ = rec["reportErrorsKey"].(string)
	}
	return key, nil
}

// listThingRegistrationTasksCore returns every registration task ID.
// Smithy: ListThingRegistrationTasksResponse.taskIds is list<TaskId>
// (string).
func (s *IoTService) listThingRegistrationTasksCore(store iotstore.IotStoreInterface) ([]string, error) {
	items, err := store.ListGeneric(thingRegistrationTaskKey)
	if err != nil {
		return nil, err
	}
	taskIds := make([]string, 0, len(items))
	for _, item := range items {
		if id, ok := item["taskId"].(string); ok {
			taskIds = append(taskIds, id)
		}
	}
	return taskIds, nil
}
