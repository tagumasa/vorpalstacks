package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CreateUserImportJob creates a new user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserImportJob.html
func (s *CognitoService) CreateUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	jobName := req.GetParam("JobName")
	cloudWatchLogsRoleArn := req.GetParam("CloudWatchLogsRoleArn")
	if userPoolID == "" || jobName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	jobID := "import-" + generateID()
	preSignedUrl := "https://cognito-import." + s.region + ".amazonaws.com/" + userPoolID + "/" + jobID

	job := &cognitostore.UserImportJob{
		JobID:                 jobID,
		JobName:               jobName,
		UserPoolID:            userPoolID,
		PreSignedUrl:          preSignedUrl,
		CreationDate:          time.Now().UTC(),
		Status:                "Created",
		CloudWatchLogsRoleArn: cloudWatchLogsRoleArn,
	}

	if err := store.CreateUserImportJob(job); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

// DescribeUserImportJob describes a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeUserImportJob.html
func (s *CognitoService) DescribeUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	jobID := req.GetParam("JobId")
	if userPoolID == "" || jobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(userPoolID, jobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

// ListUserImportJobs lists user import jobs for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserImportJobs.html
func (s *CognitoService) ListUserImportJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults := 60
	if mr := request.GetIntParam(req.Parameters, "MaxResults"); mr > 0 {
		maxResults = mr
	}

	result, err := store.ListUserImportJobsPaginated(userPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.GetParam("PaginationToken"),
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, j := range result.Items {
		formatted = append(formatted, formatUserImportJob(j))
	}

	resp := map[string]interface{}{"UserImportJobs": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["PaginationToken"] = result.NextMarker
	}
	return resp, nil
}

// StartUserImportJob starts a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartUserImportJob.html
func (s *CognitoService) StartUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	jobID := req.GetParam("JobId")
	if userPoolID == "" || jobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(userPoolID, jobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if job.Status != "Created" {
		return nil, ErrInvalidParameter
	}

	job.Status = "InProgress"
	job.StartDate = time.Now().UTC()
	if err := store.UpdateUserImportJob(job); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

// StopUserImportJob stops a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StopUserImportJob.html
func (s *CognitoService) StopUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	jobID := req.GetParam("JobId")
	if userPoolID == "" || jobID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetUserImportJob(userPoolID, jobID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if job.Status == "Completed" || job.Status == "Stopped" || job.Status == "Failed" {
		return nil, ErrInvalidParameter
	}

	job.Status = "Stopped"
	job.CompletionDate = time.Now().UTC()
	job.CompletionMessage = "Job stopped by user request."
	if err := store.UpdateUserImportJob(job); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserImportJob": formatUserImportJob(job),
	}, nil
}

func formatUserImportJob(job *cognitostore.UserImportJob) map[string]interface{} {
	result := map[string]interface{}{
		"JobName":      job.JobName,
		"JobId":        job.JobID,
		"UserPoolId":   job.UserPoolID,
		"Status":       job.Status,
		"CreationDate": job.CreationDate.Unix(),
	}
	if job.PreSignedUrl != "" {
		result["PreSignedUrl"] = job.PreSignedUrl
	}
	if !job.StartDate.IsZero() {
		result["StartDate"] = job.StartDate.Unix()
	}
	if !job.CompletionDate.IsZero() {
		result["CompletionDate"] = job.CompletionDate.Unix()
	}
	if job.CloudWatchLogsRoleArn != "" {
		result["CloudWatchLogsRoleArn"] = job.CloudWatchLogsRoleArn
	}
	if job.ImportedUsers > 0 {
		result["ImportedUsers"] = job.ImportedUsers
	}
	if job.SkippedUsers > 0 {
		result["SkippedUsers"] = job.SkippedUsers
	}
	if job.FailedUsers > 0 {
		result["FailedUsers"] = job.FailedUsers
	}
	if job.CompletionMessage != "" {
		result["CompletionMessage"] = job.CompletionMessage
	}
	return result
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
