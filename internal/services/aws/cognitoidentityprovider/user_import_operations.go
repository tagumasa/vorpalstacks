package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateUserImportJob creates a new user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserImportJob.html
func (s *CognitoService) CreateUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scheme := "http"
	if req.IsTLS {
		scheme = "https"
	}
	return s.createUserImportJobCore(ctx, reqCtx, CreateUserImportJobInput{
		UserPoolID:               req.GetParam("UserPoolId"),
		JobName:                  req.GetParam("JobName"),
		CloudWatchLogsRoleArn:    req.GetParam("CloudWatchLogsRoleArn"),
		PasswordHashingAlgorithm: req.GetParam("PasswordHashingAlgorithm"),
		UploadScheme:             scheme,
		UploadHost:               req.Host,
	})
}

// DescribeUserImportJob describes a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeUserImportJob.html
func (s *CognitoService) DescribeUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeUserImportJobCore(reqCtx, DescribeUserImportJobInput{
		UserPoolID: req.GetParam("UserPoolId"),
		JobID:      req.GetParam("JobId"),
	})
}

// ListUserImportJobs lists user import jobs for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserImportJobs.html
func (s *CognitoService) ListUserImportJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listUserImportJobsCore(reqCtx, ListUserImportJobsInput{
		UserPoolID:      req.GetParam("UserPoolId"),
		PaginationToken: req.GetParam("PaginationToken"),
		Params:          req.Parameters,
	})
}

// StartUserImportJob starts a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartUserImportJob.html
func (s *CognitoService) StartUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.startUserImportJobCore(reqCtx, StartUserImportJobInput{
		UserPoolID: req.GetParam("UserPoolId"),
		JobID:      req.GetParam("JobId"),
	})
}

// StopUserImportJob stops a user import job.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StopUserImportJob.html
func (s *CognitoService) StopUserImportJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.stopUserImportJobCore(reqCtx, StopUserImportJobInput{
		UserPoolID: req.GetParam("UserPoolId"),
		JobID:      req.GetParam("JobId"),
	})
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
	if job.PasswordHashingAlgorithm != "" {
		result["PasswordHashingAlgorithm"] = job.PasswordHashingAlgorithm
	}
	// AWS example responses always report the three counters, including
	// zero values.
	result["ImportedUsers"] = job.ImportedUsers
	result["SkippedUsers"] = job.SkippedUsers
	result["FailedUsers"] = job.FailedUsers
	if job.CompletionMessage != "" {
		result["CompletionMessage"] = job.CompletionMessage
	}
	return result
}

func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
