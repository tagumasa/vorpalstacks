package iam

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
)

// A store read failure during report generation moves the credential report
// to the FAILED state, and GetCredentialReport surfaces the reason as the
// operations' modelled ServiceFailure fault — the report is never silently
// partial.
func TestCredentialReportFailsOnStoreReadFailure(t *testing.T) {
	t.Run("user pagination failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_users", &readFault{
			failScan: true, err: errors.New("simulated user listing failure"),
		})
		if _, err := s.generateCredentialReportCore(store); err != nil {
			t.Fatalf("generate: %v", err)
		}
		s.WaitForReport()
		if s.credentialReportState != "FAILED" {
			t.Fatalf("state: got %s, want FAILED", s.credentialReportState)
		}
		_, _, err := s.getCredentialReportCore()
		assertServiceFailure(t, err, "simulated user listing failure")
	})

	t.Run("per-user signing certificate listing failure", func(t *testing.T) {
		s, store := faultTestStore(t, "iam_signing_certificates", &readFault{
			failScan: true, err: errors.New("simulated certificate listing failure"),
		})
		if _, err := store.Users().Create("report-user", "/", "123456789012", nil); err != nil {
			t.Fatalf("create user: %v", err)
		}
		if _, err := s.generateCredentialReportCore(store); err != nil {
			t.Fatalf("generate: %v", err)
		}
		s.WaitForReport()
		if s.credentialReportState != "FAILED" {
			t.Fatalf("state: got %s, want FAILED", s.credentialReportState)
		}
		_, _, err := s.getCredentialReportCore()
		assertServiceFailure(t, err, "simulated certificate listing failure")
	})
}

// assertServiceFailure pins the credential-report operations' server-fault
// shape: HTTP 500 with the modelled ServiceFailure code, carrying the
// generation failure cause.
func assertServiceFailure(t *testing.T, err error, causeSubstr string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error from a failed report generation")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, http.StatusInternalServerError, awsErr.GetHTTPStatusCode())
	assert.Contains(t, awsErr.Error(), "ServiceFailure")
	assert.Contains(t, err.Error(), "Credential report generation failed")
	assert.Contains(t, err.Error(), causeSubstr)
}

// A completed report is valid for four hours: GetCredentialReport serves
// it within the window and rejects it as ReportExpired (410) afterwards,
// leaving regeneration to GenerateCredentialReport.
func TestGetCredentialReportExpiredAfterFourHours(t *testing.T) {
	s, store := simulationTestStore(t)
	if _, err := store.Users().Create("report-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if _, err := s.generateCredentialReportCore(store); err != nil {
		t.Fatalf("generate: %v", err)
	}
	s.WaitForReport()

	content, _, err := s.getCredentialReportCore()
	if err != nil {
		t.Fatalf("a fresh report must be served: %v", err)
	}
	assert.Contains(t, content, "user,arn,user_creation_time")

	s.credentialReportMu.Lock()
	s.credentialReportTime = time.Now().UTC().Add(-reportExpiry - time.Minute)
	s.credentialReportMu.Unlock()

	_, _, err = s.getCredentialReportCore()
	if err == nil {
		t.Fatal("an expired report must be rejected")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, "ReportExpired", awsErr.Code)
	assert.Equal(t, http.StatusGone, awsErr.GetHTTPStatusCode())

	// Regeneration restarts from an expired COMPLETE state.
	state, err := s.generateCredentialReportCore(store)
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	assert.Equal(t, "STARTED", state)
	s.WaitForReport()
	if _, _, err := s.getCredentialReportCore(); err != nil {
		t.Fatalf("the regenerated report must be served: %v", err)
	}
}
