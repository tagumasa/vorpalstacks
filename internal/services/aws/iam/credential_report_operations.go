package iam

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
)

const reportExpiry = 4 * time.Hour

var (
	// ErrReportNotPresent indicates that no credential report has been generated yet.
	ErrReportNotPresent = errors.NewAWSError("ReportNotPresent", "Credential report not present. Use GenerateCredentialReport to generate one.", http.StatusGone)
	// ErrReportInProgress indicates that a credential report generation is already in progress.
	ErrReportInProgress = errors.NewAWSError("ReportInProgress", "Credential report is in progress. Please try again later.", http.StatusNotFound)
)

// GenerateCredentialReport generates a credential report for the account.
func (s *IAMService) GenerateCredentialReport(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	s.credentialReportMu.Lock()
	defer s.credentialReportMu.Unlock()

	now := time.Now().UTC()

	if s.credentialReportState == "COMPLETE" && s.credentialReportTime.Add(reportExpiry).After(now) {
		return map[string]interface{}{
			"Description": "Report already exists. No action taken.",
			"State":       "COMPLETE",
		}, nil
	}

	s.credentialReportState = "STARTED"

	store, err := s.store(reqCtx)
	if err != nil {
		s.credentialReportState = ""
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	s.reportWg.Add(1)
	go func() {
		defer s.reportWg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("PANIC in IAM credential report generation", logs.Any("panic", r))
				s.credentialReportMu.Lock()
				s.credentialReportState = ""
				s.credentialReportMu.Unlock()
			}
		}()
		time.Sleep(500 * time.Millisecond)

		s.credentialReportMu.Lock()
		defer s.credentialReportMu.Unlock()

		s.credentialReportState = "COMPLETE"
		s.credentialReportTime = time.Now().UTC()
		s.credentialReportData = generateReportContentFromStore(store)
	}()

	return map[string]interface{}{
		"Description": "No report exists. Starting a new report generation task",
		"State":       "STARTED",
	}, nil
}

// GetCredentialReport retrieves the most recently generated credential report for the account.
func (s *IAMService) GetCredentialReport(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	s.credentialReportMu.RLock()
	state := s.credentialReportState
	data := s.credentialReportData
	genTime := s.credentialReportTime
	s.credentialReportMu.RUnlock()

	switch state {
	case "":
		return nil, ErrReportNotPresent
	case "STARTED":
		return nil, ErrReportInProgress
	case "COMPLETE":
		if data == "" {
			return nil, ErrReportNotPresent
		}
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return map[string]interface{}{
		"Content":       encoded,
		"GeneratedTime": genTime.Format(time.RFC3339),
		"ReportFormat":  "text/csv",
	}, nil
}
