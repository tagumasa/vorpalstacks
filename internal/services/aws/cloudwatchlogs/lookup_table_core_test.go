package cloudwatchlogs

import (
	"errors"
	"fmt"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// seedLookupTables stores count lookup tables with ascending names.
func seedLookupTables(t *testing.T, store *logsstore.Store, count int) {
	t.Helper()
	for i := 0; i < count; i++ {
		if err := store.PutLookupTable(&logsstore.LookupTable{
			Name: fmt.Sprintf("t_%03d", i),
		}); err != nil {
			t.Fatal(err)
		}
	}
}

// An omitted maxResults serves the documented default page of 50 entries
// and a next token covering the remainder.
func TestDescribeLookupTablesDefaultPageSize(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	seedLookupTables(t, store, 60)

	page, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.LookupTables) != 50 {
		t.Fatalf("default page returned %d tables, want 50", len(page.LookupTables))
	}
	if page.NextToken == "" {
		t.Fatal("default page must carry a next token when more tables remain")
	}
	rest, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{
		NextToken: page.NextToken,
		Region:    "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rest.LookupTables) != 10 {
		t.Fatalf("second page returned %d tables, want 10", len(rest.LookupTables))
	}
	if rest.NextToken != "" {
		t.Fatalf("second page must not carry a next token, got %q", rest.NextToken)
	}
}

// maxResults above the documented maximum and negative values are rejected
// with InvalidParameterException over HTTP 400; the maximum is accepted.
func TestDescribeLookupTablesMaxResultsBounds(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	seedLookupTables(t, store, 1)

	for _, maxResults := range []int32{-1, 101} {
		_, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{
			MaxResults: maxResults,
			Region:     "us-east-1",
		})
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) {
			t.Fatalf("maxResults %d rejected with %v, want an AWS error", maxResults, err)
		}
		if awsErr.Code != "InvalidParameterException" || awsErr.HTTPStatus != 400 {
			t.Fatalf("maxResults %d rejected with %s over HTTP %d, want InvalidParameterException over HTTP 400",
				maxResults, awsErr.Code, awsErr.HTTPStatus)
		}
	}
	page, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{
		MaxResults: 100,
		Region:     "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.LookupTables) != 1 {
		t.Fatalf("maxResults at the maximum returned %d tables, want the 1 stored", len(page.LookupTables))
	}
}

// A lookup table sourced from an unknown queryId fails with
// ResourceNotFoundException over HTTP 400 — the awsJson1_1 status every
// CloudWatch Logs client error carries when the model defines no
// httpError trait.
func TestCreateLookupTableQueryNotFoundStatus(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	_, _, err := svc.createLookupTableCore(store, &LookupTableInput{
		Name:    "valid_name",
		QueryId: "no-such-query",
	}, "us-east-1")
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) {
		t.Fatalf("create with unknown queryId returned %v, want an AWS error", err)
	}
	if awsErr.Code != "ResourceNotFoundException" {
		t.Fatalf("code = %s, want ResourceNotFoundException", awsErr.Code)
	}
	if awsErr.HTTPStatus != 400 {
		t.Fatalf("HTTP status = %d, want 400", awsErr.HTTPStatus)
	}
}
