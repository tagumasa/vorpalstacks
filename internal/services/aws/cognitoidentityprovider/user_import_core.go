package cognitoidentityprovider

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/crypto"
)

// CreateUserImportJobInput carries the wire parameters of
// CreateUserImportJob. UploadScheme and UploadHost carry the wire transport
// context the presigned CSV upload URL is derived from.
type CreateUserImportJobInput struct {
	UserPoolID               string
	JobName                  string
	CloudWatchLogsRoleArn    string
	PasswordHashingAlgorithm string
	UploadScheme             string
	UploadHost               string
}

// DescribeUserImportJobInput carries the wire parameters of
// DescribeUserImportJob.
type DescribeUserImportJobInput struct {
	UserPoolID string
	JobID      string
}

// ListUserImportJobsInput carries the wire parameters of
// ListUserImportJobs. Params holds the raw request parameter map for the
// MaxResults member.
type ListUserImportJobsInput struct {
	UserPoolID      string
	PaginationToken string
	Params          map[string]interface{}
}

// StartUserImportJobInput carries the wire parameters of StartUserImportJob.
type StartUserImportJobInput struct {
	UserPoolID string
	JobID      string
}

// StopUserImportJobInput carries the wire parameters of StopUserImportJob.
type StopUserImportJobInput struct {
	UserPoolID string
	JobID      string
}

// createUserImportJobCore creates a new user import job.
func (s *CognitoService) createUserImportJobCore(ctx context.Context, reqCtx *request.RequestContext, in CreateUserImportJobInput) (interface{}, error) {
	if in.UserPoolID == "" || in.JobName == "" || in.CloudWatchLogsRoleArn == "" {
		return nil, ErrInvalidParameter
	}
	if !validateImportJobName(in.JobName) {
		return nil, ErrInvalidParameter
	}
	if _, err := svcarn.ParseARN(in.CloudWatchLogsRoleArn); err != nil {
		return nil, ErrInvalidParameter
	}
	if in.PasswordHashingAlgorithm != "" {
		if !validatePasswordHashingAlgorithm(in.PasswordHashingAlgorithm) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	id, err := generateID()
	if err != nil {
		logs.Error("failed to generate import job ID", logs.Err(err))
		return nil, ErrInternalError
	}
	jobID := "import-" + id

	if s3 := s.importS3Invoker(); s3 != nil {
		if err := s3.EnsureBucket(ctx, reqCtx.GetRegion(), importBucketName); err != nil {
			logs.Warn("failed to ensure the user import bucket",
				logs.String("bucket", importBucketName), logs.Err(err))
		}
	}

	job := &cognitostore.UserImportJob{
		JobID:                    jobID,
		JobName:                  in.JobName,
		UserPoolID:               in.UserPoolID,
		PreSignedUrl:             s.importUploadURL(in.UploadScheme, in.UploadHost, in.UserPoolID, jobID),
		CreationDate:             time.Now().UTC(),
		Status:                   importJobStatusCreated,
		CloudWatchLogsRoleArn:    in.CloudWatchLogsRoleArn,
		PasswordHashingAlgorithm: in.PasswordHashingAlgorithm,
	}

	if err := store.CreateUserImportJob(job); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

// importUploadURLExpiry is the validity window of the CSV upload URL:
// "The {{PRE_SIGNED_URL}} returned in the response is valid for 15
// minutes. After that time, it will expire and you must create a new user
// import job to get a new URL."
const importUploadURLExpiry = 15 * time.Minute

// importUploadURL builds the presigned CSV upload URL handed to the
// caller. AWS fronts its internal import bucket with the cognito-import
// virtual host and hands out a SigV4-presigned URL valid for 15 minutes;
// this platform serves the upload from the shared S3 plane, so the URL is
// a SigV4 presigned PUT target derived from the request's own host and
// scheme, carrying the same expiry.
func (s *CognitoService) importUploadURL(scheme, host, userPoolID, jobID string) string {
	if host == "" {
		// Reachable only when the request reached the handler without a
		// wire host (synthetic dispatch); live traffic carries Host from
		// the dispatch boundary.
		host = "localhost:" + strconv.Itoa(serviceports.HTTP)
	}
	accessKey, secret := "", ""
	if s.importCredentials != nil {
		if creds, err := s.importCredentials.GetCredentials(); err == nil {
			accessKey, secret = creds.AccessKeyID, creds.SecretAccessKey
		}
	}
	return crypto.PresignS3URL(http.MethodPut, scheme, host, importBucketName,
		importObjectKey(userPoolID, jobID), s.region, importUploadURLExpiry, accessKey, secret)
}

// describeUserImportJobCore describes a user import job.
func (s *CognitoService) describeUserImportJobCore(reqCtx *request.RequestContext, in DescribeUserImportJobInput) (interface{}, error) {
	if in.UserPoolID == "" || in.JobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(in.UserPoolID, in.JobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	s.expireStaleImportJob(store, job, reqCtx.GetRegion())

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

// listUserImportJobsCore lists user import jobs for a user pool.
func (s *CognitoService) listUserImportJobsCore(reqCtx *request.RequestContext, in ListUserImportJobsInput) (interface{}, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Smithy PoolQueryLimitType: range {min: 1, max: 60}
	maxResults, err := parseStrictListLimit(in.Params, "MaxResults", 60)
	if err != nil {
		return nil, err
	}

	// "Jobs are listed in chronological order from last created to first
	// created." The store lists by key, so the full set is collected,
	// sorted by creation date, and paginated here.
	var all []*cognitostore.UserImportJob
	opts := storecommon.ListOptions{MaxItems: 1000}
	for {
		result, err := store.ListUserImportJobsPaginated(in.UserPoolID, opts)
		if err != nil {
			return nil, ErrInternalError
		}
		all = append(all, result.Items...)
		if !result.IsTruncated {
			break
		}
		opts.Marker = result.NextMarker
	}
	sort.Slice(all, func(i, j int) bool { return all[i].CreationDate.After(all[j].CreationDate) })

	start := 0
	if in.PaginationToken != "" {
		offset, convErr := strconv.Atoi(in.PaginationToken)
		if convErr != nil || offset < 0 || offset > len(all) {
			return nil, ErrInvalidParameter
		}
		start = offset
	}
	end := len(all)
	if maxResults > 0 && start+maxResults < len(all) {
		end = start + maxResults
	}

	formatted := make([]map[string]interface{}, 0, end-start)
	for _, j := range all[start:end] {
		s.expireStaleImportJob(store, j, reqCtx.GetRegion())
		formatted = append(formatted, formatUserImportJob(j))
	}

	resp := map[string]interface{}{"UserImportJobs": formatted}
	if end < len(all) {
		resp["PaginationToken"] = strconv.Itoa(end)
	}
	return resp, nil
}

// startUserImportJobCore starts a user import job.
func (s *CognitoService) startUserImportJobCore(reqCtx *request.RequestContext, in StartUserImportJobInput) (interface{}, error) {
	if in.UserPoolID == "" || in.JobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(in.UserPoolID, in.JobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	s.expireStaleImportJob(store, job, reqCtx.GetRegion())
	if job.Status != importJobStatusCreated {
		return nil, ErrInvalidParameter
	}

	pool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	// "If the user pool has no auto-verified attributes, the import job
	// will not start." The error code is not documented; this platform
	// reports PreconditionNotMetException.
	if len(pool.AutoVerifiedAttributes) == 0 {
		return nil, ErrPreconditionNotMet
	}

	// "Only one import job can be active at a time per account." The
	// eligibility check and the Pending transition happen atomically in
	// the store so two concurrent starts cannot launch two workers. The
	// error code is not documented; this platform reports
	// PreconditionNotMetException.
	started, err := store.StartUserImportJobIfEligible(in.UserPoolID, in.JobID)
	if err != nil {
		if errors.Is(err, cognitostore.ErrImportJobActiveExists) || errors.Is(err, cognitostore.ErrImportJobStatusConflict) {
			return nil, ErrPreconditionNotMet
		}
		return nil, ErrInternalError
	}

	s.startUserImportWorker(reqCtx.GetRegion(), started)

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(started),
	}, nil
}

// stopUserImportJobCore stops a user import job.
func (s *CognitoService) stopUserImportJobCore(reqCtx *request.RequestContext, in StopUserImportJobInput) (interface{}, error) {
	if in.UserPoolID == "" || in.JobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(in.UserPoolID, in.JobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// Only a running job can be stopped. AWS's StopUserImportJob example
	// response reports the terminal Stopped status together with the
	// completion date and message, so the transition is finalised here;
	// the worker observes the terminal state and exits without
	// overwriting it.
	stopped := false
	for _, from := range []string{importJobStatusPending, importJobStatusInProgress, importJobStatusStopping} {
		final, transitionErr := store.TransitionUserImportJobStatus(in.UserPoolID, in.JobID, from, importJobStatusStopped, func(j *cognitostore.UserImportJob) {
			j.CompletionDate = time.Now().UTC()
			j.CompletionMessage = stoppedImportCompletionMessage
		})
		if transitionErr == nil {
			job = final
			stopped = true
			break
		}
		if !errors.Is(transitionErr, cognitostore.ErrImportJobStatusConflict) {
			return nil, ErrInternalError
		}
	}
	if !stopped {
		return nil, ErrInvalidParameter
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}
