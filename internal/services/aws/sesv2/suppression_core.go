package sesv2

import (
	"context"
	"strconv"
	"time"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — suppression-list family
// ---------------------------------------------------------------------------

// PutSuppressedDestinationInput carries the suppress-request members.
type PutSuppressedDestinationInput struct {
	EmailAddress string
	Reason       string
}

// ListSuppressedDestinationsInput carries the list members including the
// Reasons / StartDate / EndDate filter window (zero times mean absent).
type ListSuppressedDestinationsInput struct {
	MaxItems  int
	NextToken string
	Reasons   []string
	StartDate time.Time
	EndDate   time.Time
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseStoredTime converts a stored Unix-seconds string into the float the
// wire timestamps are serialised as.
func parseStoredTime(s string) float64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return float64(v)
	}
	if s != "" {
		// Silent zero-on-parse-error masked bad data in the past. Log so
		// future incidents surface in operator dashboards instead of
		// causing LastUpdateTime to silently read as 1970-01-01.
		logs.Warn("sesv2: malformed stored timestamp, returning 0",
			logs.String("value", s))
	}
	return 0
}

// parseTimestampParam parses a numeric (Unix seconds) or RFC3339 timestamp
// parameter. Returns the zero time when the parameter is absent or invalid;
// callers use IsZero to detect absence.
func parseTimestampParam(params map[string]interface{}, key string) time.Time {
	raw := request.GetStringParam(params, key)
	if raw == "" {
		return time.Time{}
	}
	if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return time.Unix(n, 0).UTC()
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t
	}
	return time.Time{}
}

// suppressedMatchesFilter applies the ListSuppressedDestinations filter
// predicate: Reasons restricts the SuppressionListReason, and
// StartDate/EndDate form an inclusive Unix-time window over LastUpdateTime.
func suppressedMatchesFilter(dest *sesv2store.SuppressedDestination, reasons []string, startDate, endDate time.Time) bool {
	if len(reasons) > 0 {
		matched := false
		for _, r := range reasons {
			if dest.Reason == r {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if !startDate.IsZero() || !endDate.IsZero() {
		t := parseStoredTime(dest.LastUpdateTime)
		ts := time.Unix(int64(t), 0).UTC()
		if !startDate.IsZero() && ts.Before(startDate) {
			return false
		}
		if !endDate.IsZero() && ts.After(endDate) {
			return false
		}
	}
	return true
}

// suppressedDestinationSummary renders the JSON shape returned by both
// GetSuppressedDestination and ListSuppressedDestinations per Smithy
// com.amazonaws.sesv2#SuppressedDestination / SuppressedDestinationSummary.
func suppressedDestinationSummary(dest *sesv2store.SuppressedDestination) map[string]interface{} {
	return map[string]interface{}{
		"EmailAddress":   dest.EmailAddress,
		"Reason":         dest.Reason,
		"LastUpdateTime": parseStoredTime(dest.LastUpdateTime),
	}
}

// ---------------------------------------------------------------------------
// Core functions — suppression-list family
// ---------------------------------------------------------------------------

// getSuppressedDestinationCore is the single entry point for reading a
// suppressed destination.
func (s *SESv2Service) getSuppressedDestinationCore(store sesv2store.SESv2StoreInterface, emailAddress string) (map[string]interface{}, error) {
	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	dest, err := store.GetSuppressedDestination(emailAddress)
	if err != nil {
		if common.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result := suppressedDestinationSummary(dest)
	if len(dest.Attributes) > 0 {
		result["Attributes"] = dest.Attributes
	}

	return map[string]interface{}{
		"SuppressedDestination": result,
	}, nil
}

// putSuppressedDestinationCore is the single entry point for adding or
// updating a suppressed destination.
func (s *SESv2Service) putSuppressedDestinationCore(store sesv2store.SESv2StoreInterface, in PutSuppressedDestinationInput) error {
	if in.EmailAddress == "" {
		return ErrMissingParameter
	}
	if !validateSuppressionEmailAddress(in.EmailAddress) {
		return ErrBadRequest
	}
	// Reason must be a valid SuppressionListReason enum value.
	if !validateSuppressionListReason(in.Reason) {
		return ErrBadRequest
	}

	dest := &sesv2store.SuppressedDestination{
		EmailAddress:   in.EmailAddress,
		Reason:         in.Reason,
		LastUpdateTime: strconv.FormatInt(time.Now().UTC().Unix(), 10),
	}

	return store.PutSuppressedDestination(dest)
}

// deleteSuppressedDestinationCore is the single entry point for removing an
// address from the suppression list.
func (s *SESv2Service) deleteSuppressedDestinationCore(store sesv2store.SESv2StoreInterface, emailAddress string) error {
	if emailAddress == "" {
		return ErrMissingParameter
	}
	return store.DeleteSuppressedDestination(emailAddress)
}

// listSuppressedDestinationsCore is the single entry point for listing the
// suppression list. Per Smithy
// com.amazonaws.sesv2#ListSuppressedDestinationsRequest the caller may
// narrow the result set by Reasons (BOUNCE/COMPLAINT) and a
// StartDate/EndDate Unix-time window; a filtered walk pages the whole set
// with a safety cap and layers pagination on top of the filtered set via
// an offset token so NextToken stays stable across pages.
func (s *SESv2Service) listSuppressedDestinationsCore(store sesv2store.SESv2StoreInterface, in ListSuppressedDestinationsInput) (map[string]interface{}, error) {
	hasFilter := len(in.Reasons) > 0 || !in.StartDate.IsZero() || !in.EndDate.IsZero()

	// Reject StartDate > EndDate.
	if !in.StartDate.IsZero() && !in.EndDate.IsZero() && in.StartDate.After(in.EndDate) {
		return nil, ErrBadRequest
	}

	if !hasFilter {
		result, err := store.ListSuppressedDestinations(common.ListOptions{
			MaxItems: in.MaxItems,
			Marker:   in.NextToken,
		})
		if err != nil {
			return nil, err
		}
		destinations := make([]map[string]interface{}, 0, len(result.Items))
		for _, dest := range result.Items {
			destinations = append(destinations, suppressedDestinationSummary(dest))
		}
		resp := map[string]interface{}{
			"SuppressedDestinationSummaries": destinations,
		}
		pagination.SetNextToken(resp, "NextToken", result.NextMarker)
		return resp, nil
	}

	// Filter is supplied: walk all pages with a safety cap and
	// apply the predicate in-memory, then layer pagination on top of the
	// filtered set via an offset token so NextToken stays stable across
	// pages.
	walkMarker := ""
	var filtered []*sesv2store.SuppressedDestination
	for {
		walkResult, err := store.ListSuppressedDestinations(common.ListOptions{MaxItems: 1000, Marker: walkMarker})
		if err != nil {
			return nil, err
		}
		for _, dest := range walkResult.Items {
			if !suppressedMatchesFilter(dest, in.Reasons, in.StartDate, in.EndDate) {
				continue
			}
			filtered = append(filtered, dest)
		}
		if len(filtered) > maxFilteredSuppressionScan {
			return nil, ErrLimitExceeded
		}
		if !walkResult.IsTruncated || walkResult.NextMarker == "" {
			break
		}
		walkMarker = walkResult.NextMarker
	}

	start := 0
	if in.NextToken != "" {
		idx, ok := decodeContactOffset(in.NextToken)
		if !ok || idx < 0 || idx >= len(filtered) {
			return nil, ErrBadRequest
		}
		start = idx
	}
	end := start + in.MaxItems
	if end > len(filtered) {
		end = len(filtered)
	}
	destinations := make([]map[string]interface{}, 0, end-start)
	for _, dest := range filtered[start:end] {
		destinations = append(destinations, suppressedDestinationSummary(dest))
	}
	resp := map[string]interface{}{
		"SuppressedDestinationSummaries": destinations,
	}
	if end < len(filtered) {
		resp["NextToken"] = encodeContactOffset(end)
	}
	return resp, nil
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------

// GetSuppressedDestination retrieves details about a suppressed email destination.
func (s *SESv2Service) GetSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getSuppressedDestinationCore(store, request.GetStringParam(req.Parameters, "EmailAddress"))
}

// PutSuppressedDestination adds or updates a suppressed destination.
func (s *SESv2Service) PutSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putSuppressedDestinationCore(store, PutSuppressedDestinationInput{
		EmailAddress: request.GetStringParam(req.Parameters, "EmailAddress"),
		Reason:       request.GetStringParam(req.Parameters, "Reason"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteSuppressedDestination removes an email address from the suppression list.
func (s *SESv2Service) DeleteSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSuppressedDestinationCore(store, request.GetStringParam(req.Parameters, "EmailAddress")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListSuppressedDestinations returns a list of suppressed destinations.
func (s *SESv2Service) ListSuppressedDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listSuppressedDestinationsCore(store, ListSuppressedDestinationsInput{
		MaxItems:  pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		NextToken: pagination.GetMarker(req.Parameters, "NextToken"),
		Reasons:   request.GetStringList(req.Parameters, "Reasons"),
		StartDate: parseTimestampParam(req.Parameters, "StartDate"),
		EndDate:   parseTimestampParam(req.Parameters, "EndDate"),
	})
}
