package cognitoidentityprovider

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// The user import data plane. CreateUserImportJob hands the caller an
// upload URL backed by the internal S3 bucket below; StartUserImportJob
// spawns runUserImportJob, which reads the CSV back through the eventbus
// S3Invoker, creates users row by row, records the per-line outcomes in
// CloudWatch Logs by publishing logs:PutLogEvents events on the service
// bus (the logs service handler owns the store instance, so the events are
// immediately readable through the CloudWatch Logs API), and deletes the
// CSV object when the job finishes — mirroring the AWS workflow documented
// in "Importing users into user pools from a CSV file".

// importBucketName is the internal S3 bucket that receives the CSV
// uploads. AWS fronts an internal bucket with the cognito-import
// virtual host; this platform serves the same role with a path-style
// bucket on the shared S3 plane.
const importBucketName = "cognito-import"

// Import job status values, matching the UserImportJobType Status
// documentation in the Smithy model.
const (
	importJobStatusCreated    = "Created"
	importJobStatusPending    = "Pending"
	importJobStatusInProgress = "InProgress"
	importJobStatusStopping   = "Stopping"
	importJobStatusStopped    = "Stopped"
	importJobStatusSucceeded  = "Succeeded"
	importJobStatusFailed     = "Failed"
	importJobStatusExpired    = "Expired"
)

// stoppedImportCompletionMessage is the completion message AWS reports for
// developer-stopped import jobs.
const stoppedImportCompletionMessage = "The Import Job was stopped by the developer."

// importObjectKey is the S3 object key a job's CSV is uploaded to and read
// back from; it doubles as the path portion of the upload URL.
func importObjectKey(userPoolID, jobID string) string {
	return userPoolID + "/" + jobID
}

// userImportLogGroup is the CloudWatch Logs group the per-line import
// outcomes are written to, per the developer guide:
// /aws/cognito/userpools/{USER_POOL_ID}/{USER_POOL_NAME}.
func userImportLogGroup(userPoolID, userPoolName string) string {
	return fmt.Sprintf("/aws/cognito/userpools/%s/%s", userPoolID, userPoolName)
}

func (s *CognitoService) importS3Invoker() eventbus.S3Invoker {
	if s.bus == nil {
		return nil
	}
	return s.bus.S3Invoker()
}

// publishImportOutcome records one per-line import outcome in CloudWatch
// Logs. The bus event path publishes a logs:PutLogEvents event that the
// logs service handler applies through its own regional store instance, so
// the outcome is immediately readable through the CloudWatch Logs API.
func (s *CognitoService) publishImportOutcome(region, logGroup, logStream string, lineNumber int, outcome, message string) {
	if s.bus == nil {
		return
	}
	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{{
			Timestamp: time.Now().UnixMilli(),
			Message:   fmt.Sprintf("[%s] Line Number %d - %s", outcome, lineNumber, message),
		}},
	}
	evt.Region = region
	evt.AccountID = s.accountID
	if err := s.bus.Publish(s.bgCtx, evt); err != nil {
		logs.Warn("Failed to publish user import log outcome",
			logs.String("logStream", logStream), logs.Err(err))
	}
}

// startUserImportWorker launches the asynchronous import worker for a job
// that has just entered the Pending state.
func (s *CognitoService) startUserImportWorker(region string, job *cognitostore.UserImportJob) {
	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
		defer s.ensureTerminalImportState(region, job.UserPoolID, job.JobID)
		defer resilience.RecoverPanic("cognito user import job")
		s.runUserImportJob(region, job.UserPoolID, job.JobID)
	}()
}

// ensureTerminalImportState guarantees the job record never stays in a
// non-terminal state after the worker exits — including after a recovered
// panic. A job left in Stopping is finalised as developer-stopped; anything
// else non-terminal is failed with an abnormal-termination message. AWS
// removes the uploaded CSV once a job reaches a terminal state, so every
// transition made here deletes the file as well.
func (s *CognitoService) ensureTerminalImportState(region, userPoolID, jobID string) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		logs.Error("Failed to resolve store to finalise user import job state", logs.String("jobId", jobID), logs.Err(err))
		return
	}
	job, err := store.GetUserImportJob(userPoolID, jobID)
	if err != nil {
		return
	}
	switch job.Status {
	case importJobStatusStopping:
		s.finalizeStoppedImport(store, job)
		s.deleteImportCSV(s.importS3Invoker(), region, userPoolID, jobID)
	case importJobStatusPending, importJobStatusInProgress:
		job.Status = importJobStatusFailed
		job.CompletionDate = time.Now().UTC()
		job.CompletionMessage = "The import job terminated before completion."
		if err := store.UpdateUserImportJob(job); err != nil {
			logs.Error("Failed to finalise an interrupted user import job", logs.String("jobId", jobID), logs.Err(err))
		}
		s.deleteImportCSV(s.importS3Invoker(), region, userPoolID, jobID)
	}
}

// finalizeStoppedImport persists the terminal Stopped state with the
// completion message AWS reports for developer-stopped jobs. The Stopping
// state is only ever observed when a stop request raced this call; the
// StopUserImportJob handler itself finalises Stopped synchronously.
func (s *CognitoService) finalizeStoppedImport(store cognitostore.CognitoStoreInterface, job *cognitostore.UserImportJob) {
	completionDate := time.Now().UTC()
	_, err := store.TransitionUserImportJobStatus(job.UserPoolID, job.JobID, importJobStatusStopping, importJobStatusStopped, func(j *cognitostore.UserImportJob) {
		j.CompletionDate = completionDate
		j.CompletionMessage = stoppedImportCompletionMessage
	})
	if err != nil && !errors.Is(err, cognitostore.ErrImportJobStatusConflict) {
		logs.Error("Failed to persist stopped user import job", logs.String("jobId", job.JobID), logs.Err(err))
	}
}

// deleteImportCSV removes the uploaded CSV after the job reaches a
// terminal state, mirroring AWS's post-completion cleanup.
func (s *CognitoService) deleteImportCSV(s3 eventbus.S3Invoker, region, userPoolID, jobID string) {
	if s3 == nil {
		return
	}
	if err := s3.DeleteObject(s.bgCtx, region, importBucketName, importObjectKey(userPoolID, jobID)); err != nil {
		logs.Warn("Failed to delete the imported CSV object",
			logs.String("bucket", importBucketName),
			logs.String("key", importObjectKey(userPoolID, jobID)), logs.Err(err))
	}
}

// expireStaleImportJob flips a Created job that was never started within
// the expiry window to Expired, lazily at read time, and deletes the job's
// uploaded data: AWS documents that an Expired job has all data associated
// with it deleted.
func (s *CognitoService) expireStaleImportJob(store cognitostore.CognitoStoreInterface, job *cognitostore.UserImportJob, region string) {
	if job.Status != importJobStatusCreated {
		return
	}
	if time.Since(job.CreationDate) <= time.Duration(cognitostore.ImportJobExpiryHours)*time.Hour {
		return
	}
	job.Status = importJobStatusExpired
	job.CompletionMessage = "The job was not started within the expiry window and its data was deleted."
	if err := store.UpdateUserImportJob(job); err != nil {
		logs.Error("Failed to persist expired user import job", logs.String("jobId", job.JobID), logs.Err(err))
	}
	s.deleteImportCSV(s.importS3Invoker(), region, job.UserPoolID, job.JobID)
}

// runUserImportJob executes one import: it loads the CSV from the internal
// bucket, validates and applies each row, writes the per-line outcome to
// CloudWatch Logs, and records the counters and terminal status on the job.
func (s *CognitoService) runUserImportJob(region, userPoolID, jobID string) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		logs.Error("Failed to resolve store for user import job", logs.String("jobId", jobID), logs.Err(err))
		return
	}

	failJob := func(reason string) {
		// AWS deletes the uploaded CSV when the job completes in any
		// terminal state, including failure.
		s.deleteImportCSV(s.importS3Invoker(), region, userPoolID, jobID)
		_, err := store.TransitionUserImportJobStatus(userPoolID, jobID, importJobStatusInProgress, importJobStatusFailed, func(j *cognitostore.UserImportJob) {
			j.CompletionDate = time.Now().UTC()
			j.CompletionMessage = reason
		})
		if err != nil && !errors.Is(err, cognitostore.ErrImportJobStatusConflict) {
			logs.Error("Failed to persist failed user import job", logs.String("jobId", jobID), logs.Err(err))
		}
	}

	job, err := store.GetUserImportJob(userPoolID, jobID)
	if err != nil {
		logs.Error("Failed to load user import job", logs.String("jobId", jobID), logs.Err(err))
		return
	}
	if job.Status == importJobStatusStopping {
		s.finalizeStoppedImport(store, job)
		s.deleteImportCSV(s.importS3Invoker(), region, userPoolID, jobID)
		return
	}
	if _, err := store.TransitionUserImportJobStatus(userPoolID, jobID, importJobStatusPending, importJobStatusInProgress, nil); err != nil {
		// The job was stopped (or already claimed) between the start
		// request and this worker; nothing to import.
		if !errors.Is(err, cognitostore.ErrImportJobStatusConflict) {
			logs.Error("Failed to mark user import job in progress", logs.String("jobId", jobID), logs.Err(err))
		}
		return
	}

	s3 := s.importS3Invoker()
	if s3 == nil {
		failJob("the S3 service is not available for user import")
		return
	}
	data, err := s3.GetObject(s.bgCtx, region, importBucketName, importObjectKey(userPoolID, jobID), cognitostore.MaxImportCSVFileSizeBytes+1)
	if err != nil {
		failJob(fmt.Sprintf("failed to read the CSV file from the import bucket: %v", err))
		return
	}
	if int64(len(data)) > cognitostore.MaxImportCSVFileSizeBytes {
		failJob(fmt.Sprintf("the CSV file exceeds the maximum size of %d bytes", cognitostore.MaxImportCSVFileSizeBytes))
		return
	}

	header, rows, err := parseImportCSV(data)
	if err != nil {
		failJob(fmt.Sprintf("unable to parse the CSV file: %v", err))
		return
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		failJob("the user pool was deleted during the import")
		return
	}
	if err := validateImportCSVHeader(pool, header); err != nil {
		failJob(fmt.Sprintf("unable to parse the CSV file: %v", err))
		return
	}
	// A job that specifies a password hashing algorithm requires the
	// password_hash column; the whole job fails without it.
	if job.PasswordHashingAlgorithm != "" && !hasColumn(header, importColumnHash) {
		failJob("the CSV file does not include the password_hash column required by the job's password hashing algorithm")
		return
	}

	// Usernames are unique case-insensitively unless the pool opts into
	// case-sensitive usernames via UsernameConfiguration.
	caseSensitiveUsernames := pool.UsernameConfiguration != nil && pool.UsernameConfiguration.CaseSensitive
	usernameKey := func(username string) string {
		if caseSensitiveUsernames {
			return username
		}
		return strings.ToLower(username)
	}

	logGroup := userImportLogGroup(userPoolID, pool.Name)
	logStream := jobID + "/" + job.JobName
	writeOutcome := func(lineNumber int, outcome, message string) {
		s.publishImportOutcome(region, logGroup, logStream, lineNumber, outcome, message)
	}

	seen := make(map[string]bool, len(rows))
	for _, csvRow := range rows {
		// A cancelled background context (server shutdown) must not
		// leave the worker draining a large file while Close waits.
		if err := s.bgCtx.Err(); err != nil {
			failJob("the import job was interrupted before completion")
			return
		}

		current, err := store.GetUserImportJob(userPoolID, jobID)
		if err == nil && current.Status != importJobStatusInProgress {
			// Stopped or otherwise finalised concurrently; exit
			// without overwriting the terminal state.
			s.deleteImportCSV(s3, region, userPoolID, jobID)
			return
		}

		// Deleting the user pool mid-import fails the whole job; rows
		// must not keep being attempted against a vanished pool.
		if _, err := store.GetUserPool(userPoolID); err != nil {
			failJob("the user pool was deleted during the import")
			return
		}

		parsed, applyErr := applyImportRow(pool, job.PasswordHashingAlgorithm, header, csvRow.Fields)
		var outcome, message string
		switch {
		case applyErr != nil:
			outcome, message = "FAILED", applyErr.Error()
		case seen[usernameKey(parsed.Username)]:
			outcome, message = "SKIPPED", "The user already exists."
		default:
			seen[usernameKey(parsed.Username)] = true
			user := cognitostore.NewUser(userPoolID, parsed.Username)
			user.Attributes = parsed.Attributes
			user.UserStatus = parsed.UserStatus
			user.PasswordHash = parsed.PasswordHash
			user.PasswordHashAlgo = parsed.PasswordHashAlgo
			user.MFAOptions = parsed.MFAOptions
			if err := store.CreateUser(user); err != nil {
				if errors.Is(err, cognitostore.ErrUserAlreadyExists) {
					outcome, message = "SKIPPED", "The user already exists."
				} else {
					outcome, message = "FAILED", fmt.Sprintf("failed to create the user: %v", err)
				}
			} else {
				outcome, message = "SUCCEEDED", "The import succeeded."
			}
		}
		if err := store.UpdateUserImportJobProgress(userPoolID, jobID, func(j *cognitostore.UserImportJob) {
			switch outcome {
			case "FAILED":
				j.FailedUsers++
			case "SKIPPED":
				j.SkippedUsers++
			default:
				j.ImportedUsers++
			}
		}); err != nil {
			logs.Error("Failed to persist user import progress", logs.String("jobId", jobID), logs.Err(err))
		}
		writeOutcome(csvRow.LineNumber, outcome, message)
	}

	// AWS publishes no numeric threshold for the "Too many users have
	// failed or been skipped during the import." failure; the one
	// official ListUserImportJobs example showing it reports
	// ImportedUsers 0 with failures, so an import that produced no users
	// ends Failed, and any import that produced at least one user ends
	// Succeeded regardless of individual row failures.
	finalStatus := importJobStatusSucceeded
	completionMessage := ""
	if counted, err := store.GetUserImportJob(userPoolID, jobID); err == nil {
		completionMessage = fmt.Sprintf("%d users were imported, %d users were skipped, and %d users failed.",
			counted.ImportedUsers, counted.SkippedUsers, counted.FailedUsers)
		if counted.ImportedUsers == 0 && counted.FailedUsers+counted.SkippedUsers > 0 {
			finalStatus = importJobStatusFailed
			completionMessage = "Too many users have failed or been skipped during the import."
		}
	}
	completionDate := time.Now().UTC()
	if _, err := store.TransitionUserImportJobStatus(userPoolID, jobID, importJobStatusInProgress, finalStatus, func(j *cognitostore.UserImportJob) {
		j.CompletionDate = completionDate
		j.CompletionMessage = completionMessage
	}); err != nil {
		if current, getErr := store.GetUserImportJob(userPoolID, jobID); getErr == nil && current.Status == importJobStatusStopping {
			s.finalizeStoppedImport(store, current)
		} else if !errors.Is(err, cognitostore.ErrImportJobStatusConflict) {
			logs.Error("Failed to persist completed user import job", logs.String("jobId", jobID), logs.Err(err))
		}
	}
	s.deleteImportCSV(s3, region, userPoolID, jobID)
}
