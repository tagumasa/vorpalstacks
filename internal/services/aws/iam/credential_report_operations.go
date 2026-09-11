package iam

import (
	"context"
	"encoding/base64"
	"time"

	"vorpalstacks/internal/common/request"
)

// GenerateCredentialReport generates a credential report for the account.
func (s *IAMService) GenerateCredentialReport(_ context.Context, reqCtx *request.RequestContext, _ *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	state, err := s.generateCredentialReportCore(store)
	if err != nil {
		return nil, err
	}

	// The STARTED and COMPLETE descriptions follow the documented CLI
	// example response; the model defines no per-state Description values,
	// so the INPROGRESS wording is ours.
	description := "No report exists. Starting a new report generation task"
	switch state {
	case "COMPLETE":
		description = "Report already exists. No action taken."
	case "INPROGRESS":
		description = "Report generation is in progress."
	}

	return map[string]interface{}{
		"Description": description,
		"State":       state,
	}, nil
}

// GetCredentialReport retrieves the most recently generated credential report for the account.
func (s *IAMService) GetCredentialReport(_ context.Context, _ *request.RequestContext, _ *request.ParsedRequest) (interface{}, error) {
	data, genTime, err := s.getCredentialReportCore()
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return map[string]interface{}{
		"Content":       encoded,
		"GeneratedTime": genTime.Format(time.RFC3339),
		"ReportFormat":  "text/csv",
	}, nil
}
