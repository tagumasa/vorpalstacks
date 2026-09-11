// Transport-agnostic Core functions for the IAM credential report: the
// generation state machine and the store aggregation shared by the
// AWS-compatible HTTP API handlers and any admin plane paths (the xxxCore
// pattern).
package iam

import (
	"bytes"
	"cmp"
	stderrors "errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

const reportExpiry = 4 * time.Hour

var (
	// ErrReportNotPresent indicates that no credential report has been generated yet.
	ErrReportNotPresent = errors.NewAWSError("ReportNotPresent", "Credential report not present. Use GenerateCredentialReport to generate one.", http.StatusGone)
	// ErrReportInProgress indicates that a credential report generation is already in progress.
	ErrReportInProgress = errors.NewAWSError("ReportInProgress", "Credential report is in progress. Please try again later.", http.StatusNotFound)
)

// generateCredentialReportCore starts a credential report generation, or
// reports the state of an existing one. It returns the wire state:
// "STARTED" when a generation task has been launched, "INPROGRESS" when a
// generation is already running (a second Generate reports the running
// task instead of starting another one), "COMPLETE" when the existing
// report is still fresh and no generation was started. The generation
// itself runs on a background goroutine; this call returns as soon as the
// state machine has moved to STARTED, without waiting for the report
// content. The state mutex guards state transitions only — the content
// build runs outside it — so a GetCredentialReport arriving during a
// generation observes STARTED and receives the ReportInProgress fault
// instead of blocking until the build finishes. A generation whose store
// reads fail moves the state to FAILED with the reason (surfaced by
// getCredentialReportCore, and never a silently partial report); the next
// GenerateCredentialReport restarts from FAILED as from any non-COMPLETE
// state.
func (s *IAMService) generateCredentialReportCore(store *iamstore.IAMStore) (string, error) {
	s.credentialReportMu.Lock()
	defer s.credentialReportMu.Unlock()

	if s.credentialReportState == "COMPLETE" && s.credentialReportTime.Add(reportExpiry).After(time.Now().UTC()) {
		return "COMPLETE", nil
	}
	if s.credentialReportState == "STARTED" {
		return "INPROGRESS", nil
	}

	s.credentialReportState = "STARTED"

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

		content, err := generateReportContentFromStore(store)

		s.credentialReportMu.Lock()
		defer s.credentialReportMu.Unlock()
		if err != nil {
			logs.Error("IAM credential report generation failed", logs.Err(err))
			s.credentialReportState = "FAILED"
			s.credentialReportErr = err.Error()
			s.credentialReportData = ""
			return
		}
		s.credentialReportState = "COMPLETE"
		s.credentialReportErr = ""
		s.credentialReportTime = time.Now().UTC()
		s.credentialReportData = content
	}()

	return "STARTED", nil
}

// getCredentialReportCore returns the completed report content and its
// generation time. No report yet (or an empty one) maps to
// ErrReportNotPresent; a still-running generation maps to
// ErrReportInProgress; a failed generation maps to a ServiceFailure
// carrying the generation failure reason (the operations' modelled
// server-fault vocabulary); a completed report older than the four-hour
// validity window maps to ErrReportExpired.
func (s *IAMService) getCredentialReportCore() (string, time.Time, error) {
	s.credentialReportMu.RLock()
	state := s.credentialReportState
	data := s.credentialReportData
	genTime := s.credentialReportTime
	failReason := s.credentialReportErr
	s.credentialReportMu.RUnlock()

	switch state {
	case "":
		return "", time.Time{}, ErrReportNotPresent
	case "STARTED":
		return "", time.Time{}, ErrReportInProgress
	case "FAILED":
		return "", time.Time{}, NewServiceFailureException("Credential report generation failed: " + failReason)
	case "COMPLETE":
		if data == "" {
			return "", time.Time{}, ErrReportNotPresent
		}
		// A report is valid for four hours; an older one is expired and
		// the caller must regenerate (GenerateCredentialReport restarts
		// from an expired COMPLETE state).
		if !genTime.Add(reportExpiry).After(time.Now().UTC()) {
			return "", time.Time{}, ErrReportExpired
		}
	}

	return data, genTime, nil
}

// generateReportContentFromStore builds the account-wide report CSV. Every
// store read failure fails the generation — the credential report is a
// security-audit artefact and must never be silently partial: a truncated
// user list or an unreadable credential listing is reported as a generation
// failure, not emitted as if it were the account's true state.
func generateReportContentFromStore(store *iamstore.IAMStore) (string, error) {
	var allUsers []*iamstore.User
	marker := ""
	for {
		result, err := store.Users().List("", marker, 1000)
		if err != nil {
			return "", fmt.Errorf("list users: %w", err)
		}
		allUsers = append(allUsers, result.Users...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}

	var buf bytes.Buffer

	buf.WriteString(csvEscape("user") + "," + csvEscape("arn") + "," + csvEscape("user_creation_time") + "," +
		csvEscape("password_enabled") + "," + csvEscape("password_last_used") + "," +
		csvEscape("password_last_changed") + "," + csvEscape("password_next_rotation") + "," +
		csvEscape("mfa_active") + "," +
		csvEscape("access_key_1_active") + "," + csvEscape("access_key_1_last_rotated") + "," +
		csvEscape("access_key_1_last_used_date") + "," + csvEscape("access_key_1_last_used_region") + "," +
		csvEscape("access_key_1_last_used_service") + "," +
		csvEscape("access_key_2_active") + "," + csvEscape("access_key_2_last_rotated") + "," +
		csvEscape("access_key_2_last_used_date") + "," + csvEscape("access_key_2_last_used_region") + "," +
		csvEscape("access_key_2_last_used_service") + "," +
		csvEscape("cert_1_active") + "," + csvEscape("cert_1_last_rotated") + "," +
		csvEscape("cert_2_active") + "," + csvEscape("cert_2_last_rotated") + "\n")

	for _, user := range allUsers {
		mfaCount, err := store.MFADevices().CountForUser(user.UserName)
		if err != nil {
			return "", fmt.Errorf("count MFA devices for user %s: %w", user.UserName, err)
		}
		mfaActive := "FALSE"
		if mfaCount > 0 {
			mfaActive = "TRUE"
		}

		keys, err := store.AccessKeys().ListByUserName(user.UserName)
		if err != nil {
			return "", fmt.Errorf("list access keys for user %s: %w", user.UserName, err)
		}
		slices.SortFunc(keys, func(a, b *iamstore.AccessKey) int { return cmp.Compare(a.AccessKeyId, b.AccessKeyId) })

		ak1Active := "FALSE"
		ak1LastRotated := "N/A"
		ak1LastUsedDate := "N/A"
		ak1LastUsedRegion := "N/A"
		ak1LastUsedService := "N/A"
		if len(keys) > 0 {
			if keys[0].Status == iamstore.AccessKeyStatusActive {
				ak1Active = "TRUE"
			}
			ak1LastRotated = keys[0].CreateDate.Format(timeutils.ISO8601SimpleFormat)
			if keys[0].LastUsedDate != nil {
				ak1LastUsedDate = keys[0].LastUsedDate.Format(timeutils.ISO8601SimpleFormat)
			}
			if keys[0].LastUsedRegion != "" {
				ak1LastUsedRegion = keys[0].LastUsedRegion
			}
			if keys[0].LastUsedService != "" {
				ak1LastUsedService = keys[0].LastUsedService
			}
		}

		ak2Active := "FALSE"
		ak2LastRotated := "N/A"
		ak2LastUsedDate := "N/A"
		ak2LastUsedRegion := "N/A"
		ak2LastUsedService := "N/A"
		if len(keys) > 1 {
			if keys[1].Status == iamstore.AccessKeyStatusActive {
				ak2Active = "TRUE"
			}
			ak2LastRotated = keys[1].CreateDate.Format(timeutils.ISO8601SimpleFormat)
			if keys[1].LastUsedDate != nil {
				ak2LastUsedDate = keys[1].LastUsedDate.Format(timeutils.ISO8601SimpleFormat)
			}
			if keys[1].LastUsedRegion != "" {
				ak2LastUsedRegion = keys[1].LastUsedRegion
			}
			if keys[1].LastUsedService != "" {
				ak2LastUsedService = keys[1].LastUsedService
			}
		}

		passwordEnabled := "FALSE"
		passwordLastUsed := "no_information"
		passwordLastChanged := "N/A"
		// One resolved read per row: not-found means the user has no
		// login profile; any other failure aborts the report.
		var profile *iamstore.LoginProfile
		if p, err := store.LoginProfiles().Get(user.UserName); err == nil {
			profile = p
			passwordEnabled = "TRUE"
			if user.PasswordLastUsed != nil {
				passwordLastUsed = user.PasswordLastUsed.Format(timeutils.ISO8601SimpleFormat)
			}
			if !profile.PasswordChangedAt.IsZero() {
				passwordLastChanged = profile.PasswordChangedAt.Format(timeutils.ISO8601SimpleFormat)
			} else {
				passwordLastChanged = profile.CreateDate.Format(timeutils.ISO8601SimpleFormat)
			}
		} else if !stderrors.Is(err, iamstore.ErrLoginProfileNotFound) {
			return "", fmt.Errorf("read login profile for user %s: %w", user.UserName, err)
		}

		certs, err := store.SigningCertificates().ListByUserName(user.UserName)
		if err != nil {
			return "", fmt.Errorf("list signing certificates for user %s: %w", user.UserName, err)
		}
		cert1Active := "FALSE"
		cert1LastRotated := "N/A"
		if len(certs) > 0 {
			if certs[0].Status == "Active" {
				cert1Active = "TRUE"
			}
			cert1LastRotated = certs[0].UploadDate.Format(timeutils.ISO8601SimpleFormat)
		}
		cert2Active := "FALSE"
		cert2LastRotated := "N/A"
		if len(certs) > 1 {
			if certs[1].Status == "Active" {
				cert2Active = "TRUE"
			}
			cert2LastRotated = certs[1].UploadDate.Format(timeutils.ISO8601SimpleFormat)
		}

		passwordNextRotation := "N/A"
		if profile != nil {
			policy := store.PasswordPolicy().GetOrDefault()
			if policy.MaxPasswordAge > 0 {
				base := profile.PasswordChangedAt
				if base.IsZero() {
					base = profile.CreateDate
				}
				nextRotation := base.AddDate(0, 0, policy.MaxPasswordAge)
				passwordNextRotation = nextRotation.Format(timeutils.ISO8601SimpleFormat)
			}
		}

		row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			csvEscape(user.UserName),
			csvEscape(user.Arn),
			csvEscape(user.CreateDate.Format(timeutils.ISO8601SimpleFormat)),
			passwordEnabled,
			passwordLastUsed,
			passwordLastChanged,
			passwordNextRotation,
			mfaActive,
			ak1Active,
			ak1LastRotated,
			ak1LastUsedDate,
			ak1LastUsedRegion,
			ak1LastUsedService,
			ak2Active,
			ak2LastRotated,
			ak2LastUsedDate,
			ak2LastUsedRegion,
			ak2LastUsedService,
			cert1Active,
			cert1LastRotated,
			cert2Active,
			cert2LastRotated,
		)
		buf.WriteString(row)
	}

	report := strings.TrimRight(buf.String(), "\n")
	return report, nil
}

func csvEscape(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") || strings.Contains(s, "\n") || strings.Contains(s, "\r") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
