package iam

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

var (
	credentialReportMu    sync.RWMutex
	credentialReportState string
	credentialReportData  string
	credentialReportTime  time.Time
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
	credentialReportMu.Lock()
	defer credentialReportMu.Unlock()

	now := time.Now().UTC()

	if credentialReportState == "COMPLETE" && credentialReportTime.Add(reportExpiry).After(now) {
		return map[string]interface{}{
			"Description": "No report exists. Starting a new report generation task",
			"State":       "COMPLETE",
		}, nil
	}

	credentialReportState = "STARTED"

	store, err := s.store(reqCtx)
	if err != nil {
		credentialReportState = ""
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	s.reportWg.Add(1)
	go func() {
		defer s.reportWg.Done()
		time.Sleep(500 * time.Millisecond)

		credentialReportMu.Lock()
		defer credentialReportMu.Unlock()

		credentialReportState = "COMPLETE"
		credentialReportTime = time.Now().UTC()
		credentialReportData = generateReportContentFromStore(store)
	}()

	return map[string]interface{}{
		"Description": "No report exists. Starting a new report generation task",
		"State":       "STARTED",
	}, nil
}

// GetCredentialReport retrieves the most recently generated credential report for the account.
func (s *IAMService) GetCredentialReport(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	credentialReportMu.RLock()
	state := credentialReportState
	data := credentialReportData
	genTime := credentialReportTime
	credentialReportMu.RUnlock()

	switch state {
	case "":
		return nil, ErrReportNotPresent
	case "STARTED", "INPROGRESS":
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

func generateReportContentFromStore(store *iamstore.IAMStore) string {
	var allUsers []*iamstore.User
	marker := ""
	for {
		result, err := store.Users().List("", marker, 1000)
		if err != nil {
			break
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
		mfaCount, _ := store.MFADevices().CountForUser(user.UserName)
		mfaActive := "FALSE"
		if mfaCount > 0 {
			mfaActive = "TRUE"
		}

		keys, _ := store.AccessKeys().ListByUserName(user.UserName)

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
		if store.LoginProfiles().Exists(user.UserName) {
			passwordEnabled = "TRUE"
			passwordLastUsed = "no_information"
			if user.PasswordLastUsed != nil {
				passwordLastUsed = user.PasswordLastUsed.Format(timeutils.ISO8601SimpleFormat)
			}
			passwordLastChanged = "N/A"
		}

		certs, _ := store.SigningCertificates().ListByUserName(user.UserName)
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

		row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			csvEscape(user.UserName),
			csvEscape(user.Arn),
			csvEscape(user.CreateDate.Format(timeutils.ISO8601SimpleFormat)),
			passwordEnabled,
			passwordLastUsed,
			passwordLastChanged,
			"not_supported",
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
	return report
}

func csvEscape(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") || strings.Contains(s, "\n") || strings.Contains(s, "\r") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
