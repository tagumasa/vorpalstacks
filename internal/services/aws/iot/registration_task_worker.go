package iot

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/services/aws/iot/ca"
	iotstore "vorpalstacks/internal/store/aws/iot"
	"vorpalstacks/internal/utils/crypto"
)

// The bulk thing-registration data plane. StartThingRegistrationTask
// records the task and launches runThingRegistrationTask, which reads the
// newline-delimited JSON input file back through the eventbus S3Invoker
// (each line carries the parameter values to provision one device),
// provisions every line through the shared provisioning-template engine,
// tracks success/failure counts and progress, honours StopThingRegistrationTask's
// Cancelling handshake, and writes the RESULTS/ERRORS CSV reports that
// ListThingRegistrationTaskReports hands out as presigned URLs.

// registrationReportBucket is the platform-managed bucket the task reports
// are written to. AWS generates the reports in a service-managed bucket and
// returns presigned links; this platform serves the same role with a
// path-style bucket on the shared S3 plane.
const registrationReportBucket = "iot-thing-registration-reports"

// MaxRegistrationInputFileBytes caps the S3 input file the worker reads.
// AWS documents no size bound for the input file; the cap keeps a
// mispointed task from unbounded memory use.
const MaxRegistrationInputFileBytes = 16 << 20

// registrationReportLinkExpiry is the validity window of the presigned
// report URLs handed out by ListThingRegistrationTaskReports.
const registrationReportLinkExpiry = 15 * time.Minute

// registrationTaskWg tracks the in-flight registration workers so
// Shutdown can wait for their task records to reach a terminal state
// before the storage layer closes.
var registrationTaskWg sync.WaitGroup

// registrationShutdownWait bounds how long Shutdown waits for in-flight
// registration workers. Tasks are line-bounded by the input-file cap and
// normally finish in milliseconds; the bound keeps a stalled worker from
// holding server shutdown hostage. A worker that overruns it is
// abandoned with its record left InProgress — the same state a hard
// kill leaves.
const registrationShutdownWait = 10 * time.Second

// waitRegistrationWorkers waits for every in-flight registration worker
// to finish, giving up once the bound elapses. The helper goroutine
// outlives the timeout until the last worker exits; at shutdown that is
// harmless by construction.
func waitRegistrationWorkers(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		registrationTaskWg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// registrationLineOutcome records one input line's provisioning result for
// the report CSVs.
type registrationLineOutcome struct {
	Line     int
	ThingArn string
	Error    string
}

func (s *IoTService) registrationS3Invoker() invokers.S3Invoker {
	if s.deps.EventBus == nil {
		return nil
	}
	return s.deps.EventBus.S3Invoker()
}

// startThingRegistrationWorker launches the asynchronous worker for a task
// that has just been recorded in the InProgress state.
func (s *IoTService) startThingRegistrationWorker(store iotstore.IotStoreInterface, region, taskID, bucket, key string, tpl *provisioningTemplate) {
	authority := s.deps.CAs[region]
	registrationTaskWg.Add(1)
	go func() {
		defer registrationTaskWg.Done()
		defer resilience.RecoverPanic("iot thing registration task")
		s.runThingRegistrationTask(store, region, taskID, bucket, key, tpl, authority)
	}()
}

// updateThingRegistrationTask loads the record, applies mutate, and
// persists it with a refreshed lastModifiedDate.
func updateThingRegistrationTask(store iotstore.IotStoreInterface, taskID string, mutate func(rec map[string]interface{})) error {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(thingRegistrationTaskKey+taskID, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrThingRegistrationTaskNotFound
	}
	mutate(rec)
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	return store.PutGeneric(thingRegistrationTaskKey+taskID, rec)
}

func failThingRegistrationTask(store iotstore.IotStoreInterface, taskID, message string) {
	if err := updateThingRegistrationTask(store, taskID, func(rec map[string]interface{}) {
		rec["status"] = taskStatusFailed
		rec["message"] = message
	}); err != nil {
		logs.Error("Failed to persist thing-registration task failure",
			logs.String("taskId", taskID), logs.Err(err))
	}
}

// runThingRegistrationTask executes one bulk registration task end to end.
func (s *IoTService) runThingRegistrationTask(store iotstore.IotStoreInterface, region, taskID, bucket, key string, tpl *provisioningTemplate, authority *ca.CertificateAuthority) {
	ctx := context.Background()
	s3 := s.registrationS3Invoker()
	if s3 == nil {
		failThingRegistrationTask(store, taskID, "S3 access is not available for the request region")
		return
	}
	data, err := s3.GetObject(ctx, region, bucket, key, MaxRegistrationInputFileBytes)
	if err != nil {
		failThingRegistrationTask(store, taskID, fmt.Sprintf("unable to read the input file s3://%s/%s: %v", bucket, key, err))
		return
	}

	type pendingLine struct {
		lineNo int
		params map[string]string
	}
	var lines []pendingLine
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		var decoded map[string]interface{}
		if err := json.Unmarshal([]byte(line), &decoded); err != nil {
			failThingRegistrationTask(store, taskID, fmt.Sprintf("input line %d is not a JSON object: %v", i+1, err))
			return
		}
		params := make(map[string]string, len(decoded))
		for k, v := range decoded {
			if str, ok := v.(string); ok {
				params[k] = str
				continue
			}
			b, mErr := json.Marshal(v)
			if mErr != nil {
				params[k] = fmt.Sprintf("%v", v)
				continue
			}
			params[k] = string(b)
		}
		lines = append(lines, pendingLine{lineNo: i + 1, params: params})
	}

	var outcomes []registrationLineOutcome
	success, failure := 0, 0
	for _, line := range lines {
		// Honour the cancel handshake between lines.
		rec := map[string]interface{}{}
		exists, err := store.GetGenericExists(thingRegistrationTaskKey+taskID, &rec)
		if err != nil || !exists {
			logs.Error("Failed to load thing-registration task record mid-run",
				logs.String("taskId", taskID), logs.Err(err))
			return
		}
		if rec["status"] == taskStatusCancelling {
			if err := updateThingRegistrationTask(store, taskID, func(r map[string]interface{}) {
				r["status"] = taskStatusCancelled
			}); err != nil {
				logs.Error("Failed to persist cancelled thing-registration task",
					logs.String("taskId", taskID), logs.Err(err))
			}
			return
		}

		outcome := registrationLineOutcome{Line: line.lineNo}
		if result, pErr := s.provisionFromTemplate(store, authority, tpl, line.params); pErr == nil {
			success++
			outcome.ThingArn = firstThingArn(result.resourceArns)
		} else {
			failure++
			outcome.Error = pErr.Error()
		}
		outcomes = append(outcomes, outcome)
		processed := success + failure
		if err := updateThingRegistrationTask(store, taskID, func(r map[string]interface{}) {
			r["successCount"] = success
			r["failureCount"] = failure
			r["percentageProgress"] = processed * 100 / len(lines)
		}); err != nil {
			logs.Error("Failed to persist thing-registration task progress",
				logs.String("taskId", taskID), logs.Err(err))
		}
	}

	resultsKey, errorsKey := taskReportKeys(taskID)
	s3inv := s3
	if err := s3inv.EnsureBucket(ctx, region, registrationReportBucket); err != nil {
		logs.Warn("Failed to ensure the thing-registration report bucket",
			logs.String("bucket", registrationReportBucket), logs.Err(err))
	}
	if err := s3inv.PutObject(ctx, region, registrationReportBucket, resultsKey,
		[]byte(renderReportCSV(outcomes, true)), "text/csv"); err != nil {
		failThingRegistrationTask(store, taskID, fmt.Sprintf("unable to write the results report: %v", err))
		return
	}
	if err := s3inv.PutObject(ctx, region, registrationReportBucket, errorsKey,
		[]byte(renderReportCSV(outcomes, false)), "text/csv"); err != nil {
		failThingRegistrationTask(store, taskID, fmt.Sprintf("unable to write the errors report: %v", err))
		return
	}
	if err := updateThingRegistrationTask(store, taskID, func(r map[string]interface{}) {
		r["successCount"] = success
		r["failureCount"] = failure
		r["percentageProgress"] = 100
		r["reportResultsKey"] = resultsKey
		r["reportErrorsKey"] = errorsKey
		// A Stop that raced the final line must not be overwritten by the
		// completion write.
		if r["status"] != taskStatusCancelling {
			r["status"] = taskStatusCompleted
		}
	}); err != nil {
		logs.Error("Failed to persist completed thing-registration task",
			logs.String("taskId", taskID), logs.Err(err))
	}
}

// firstThingArn picks the thing ARN out of a resourceArns map for the
// report rows; logical names are arbitrary, so the ARN's own resource
// prefix identifies the thing.
func firstThingArn(arns map[string]string) string {
	for _, arn := range arns {
		if strings.Contains(arn, ":thing/") {
			return arn
		}
	}
	return ""
}

// taskReportKeys returns the S3 object keys of a task's two reports.
func taskReportKeys(taskID string) (resultsKey, errorsKey string) {
	return taskID + "/results.csv", taskID + "/errors.csv"
}

// renderReportCSV renders one report. AWS publishes the reports as CSV
// without documenting the column schema in the API reference; the platform
// reports carry the line number, the provisioned thing's ARN, and the
// failure message (results report) or the failure rows alone (errors
// report).
func renderReportCSV(outcomes []registrationLineOutcome, results bool) string {
	var b strings.Builder
	w := csv.NewWriter(&b)
	if results {
		_ = w.Write([]string{"LineNumber", "ThingArn", "Status"})
		for _, o := range outcomes {
			status := "Success"
			thingArn := o.ThingArn
			if o.Error != "" {
				status = "Failure: " + o.Error
				if thingArn == "" {
					thingArn = "-"
				}
			}
			_ = w.Write([]string{strconv.Itoa(o.Line), thingArn, status})
		}
	} else {
		_ = w.Write([]string{"LineNumber", "ErrorMessage"})
		for _, o := range outcomes {
			if o.Error != "" {
				_ = w.Write([]string{strconv.Itoa(o.Line), o.Error})
			}
		}
	}
	w.Flush()
	return b.String()
}

// thingRegistrationReportLinks presigns the GET URLs for a task's report of
// the requested type. AWS's links expire; the platform links carry the same
// expiry window derived from the caller's own endpoint and the service
// credentials. The expiry governs the presigned URLs only — the
// ListThingRegistrationTaskReports response shape has no member for it.
func (s *IoTService) thingRegistrationReportLinks(reqCtx *request.RequestContext, req *request.ParsedRequest, store iotstore.IotStoreInterface, taskID, reportType string) ([]string, error) {
	reportKey, err := s.thingRegistrationReportKeyCore(store, taskID, reportType)
	if err != nil {
		return nil, err
	}
	if reportKey == "" {
		// The reports exist once the task reaches a terminal state.
		return []string{}, nil
	}
	scheme := "http"
	if req.IsTLS {
		scheme = "https"
	}
	host := req.Host
	if host == "" {
		// Reachable only when the request reached the handler without a
		// wire host (synthetic dispatch); live traffic carries Host from
		// the dispatch boundary.
		host = "localhost:" + strconv.Itoa(serviceports.HTTP)
	}
	accessKey, secret := "", ""
	if s.reportCredentials != nil {
		if creds, err := s.reportCredentials.GetCredentials(); err == nil {
			accessKey, secret = creds.AccessKeyID, creds.SecretAccessKey
		}
	}
	region := ""
	if reqCtx != nil {
		region = reqCtx.GetRegion()
	}
	link := crypto.PresignS3URL("GET", scheme, host, registrationReportBucket,
		reportKey, region, registrationReportLinkExpiry, accessKey, secret)
	return []string{link}, nil
}
