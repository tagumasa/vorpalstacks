package cloudwatchlogs

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
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

// The tags member carries the standard Tags entry traits (per-entry
// TagKey/TagValue patterns, lengths and the reserved aws: prefix) and
// UpdateLookupTable revalidates the description it merges — the create
// path's spec validator applies to both halves of the family.
func TestLookupTableTagsAndUpdateSpecTraits(t *testing.T) {
	svc, store := newReadTestService(t, "lookup-tags-group")

	for name, tags := range map[string]map[string]string{
		"reserved prefix":   {"aws:reserved": "v"},
		"bad key alphabet":  {"bad key!": "v"},
		"bad value pattern": {"k": "bad\nvalue"},
	} {
		_, _, err := svc.createLookupTableCore(store, &LookupTableInput{
			Name: "tags_reject", TableBody: "a,b\n1,2\n", Tags: tags,
		}, "us-east-1")
		var logsErr *awserrors.AWSError
		if !errors.As(err, &logsErr) || logsErr.Code != "InvalidParameterException" {
			t.Fatalf("%s: err=%v, want InvalidParameterException", name, err)
		}
	}
	if _, _, err := svc.createLookupTableCore(store, &LookupTableInput{
		Name: "tags_accept", TableBody: "a,b\n1,2\n", Tags: map[string]string{"team": "logs"},
	}, "us-east-1"); err != nil {
		t.Fatalf("well-formed tags rejected: %v", err)
	}

	// The merged description rides the same 0-1024 bound create enforces.
	_, err := svc.updateLookupTableCore(store, &UpdateLookupTableInput{
		Identifier:     "tags_accept",
		Description:    strings.Repeat("d", logsstore.MaxLookupTableDescriptionLength+1),
		DescriptionSet: true,
		Region:         "us-east-1",
	})
	var logsErr *awserrors.AWSError
	if !errors.As(err, &logsErr) || logsErr.Code != "InvalidParameterException" {
		t.Fatalf("oversize description on update: err=%v, want InvalidParameterException", err)
	}
	// The in-bound description rides a full-replacement update ("You must
	// specify either tableBody or queryId" — the operation replaces all
	// existing content, so the body travels with the description).
	if _, err := svc.updateLookupTableCore(store, &UpdateLookupTableInput{
		Identifier:     "tags_accept",
		Description:    "short",
		DescriptionSet: true,
		TableBody:      "a,b\n1,2\n",
		TableBodySet:   true,
		Region:         "us-east-1",
	}); err != nil {
		t.Fatalf("in-bound description on update: %v", err)
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

// Concurrent creations of the same lookup-table name admit exactly one:
// the exists-check and the write run as one critical section, so the
// loser observes the winner's record instead of writing a second one.
func TestCreateLookupTableSameNameAdmitsExactlyOne(t *testing.T) {
	svc, store := newReadTestService(t, "lookup-admission-group")
	const racers = 8
	var wg sync.WaitGroup
	var succeeded int32
	var alreadyExists int32
	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := svc.createLookupTableCore(store, &LookupTableInput{
				Name:      "raced_table",
				TableBody: "a,b\n1,2\n",
			}, "us-east-1")
			switch {
			case err == nil:
				atomic.AddInt32(&succeeded, 1)
			case strings.Contains(err.Error(), "already exists"):
				atomic.AddInt32(&alreadyExists, 1)
			}
		}()
	}
	wg.Wait()
	if succeeded != 1 || alreadyExists != racers-1 {
		t.Fatalf("admission must admit exactly one creation: succeeded=%d alreadyExists=%d", succeeded, alreadyExists)
	}
}

// The per-account quota census runs inside the creation admission window:
// two creations racing at a full-1 census admit exactly one — the loser
// reads the winner's table in its own census and answers
// LimitExceededException, never a second admission.
func TestCreateLookupTableQuotaCensusAdmitsExactlyOne(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	seedLookupTables(t, store, logsstore.MaxLookupTables-1)

	var succeeded, rejected atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		name := fmt.Sprintf("raced_%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := svc.createLookupTableCore(store, &LookupTableInput{
				Name:      name,
				TableBody: "h1,h2\n1,2\n",
			}, "us-east-1")
			var logsErr *awserrors.AWSError
			switch {
			case err == nil:
				succeeded.Add(1)
			case errors.As(err, &logsErr) && logsErr.Code == "LimitExceededException":
				rejected.Add(1)
			default:
				t.Errorf("create %s: %v, want success or LimitExceededException", name, err)
			}
		}()
	}
	wg.Wait()

	if succeeded.Load() != 1 || rejected.Load() != 1 {
		t.Fatalf("quota census admitted %d and rejected %d, want exactly one of each",
			succeeded.Load(), rejected.Load())
	}
}

// The tableBody content rule — "The content must use UTF-8 encoding and
// not exceed 10 MB" (PutLookupTable) — and the DescribeLookupTables prefix
// member's LookupTableName traits (1-256 characters of the name alphabet):
// an out-of-trait value rejects instead of storing unreadable bytes or
// silently matching nothing.
func TestLookupTableBodyAndPrefixTraits(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)

	lt := &logsstore.LookupTable{Name: "utf8_table"}
	if err := svc.applyLookupTableBody(lt, string([]byte{0xff, 0xfe, 'b'}), "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("non-UTF-8 body: %v", err)
	}
	if err := svc.applyLookupTableBody(lt, "k,v\n", "us-east-1"); err != nil {
		t.Fatalf("valid UTF-8 body rejected: %v", err)
	}

	if _, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{Prefix: "bad prefix!", Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("invalid prefix: %v", err)
	}
	if _, err := svc.describeLookupTablesCore(store, &DescribeLookupTablesInput{Prefix: strings.Repeat("p", logsstore.MaxLookupTableNameLength+1), Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("oversize prefix: %v", err)
	}
}

// A read whose decryption the substrate cannot complete — the KMS leg
// absent, or the key refusing the decrypt — is a server-side condition of
// the table's storage, not a fault of the valid request: the declared
// ServiceUnavailableException answers instead of InvalidParameterException.
func TestLookupTableKMSReadFailureClass(t *testing.T) {
	encrypted := &logsstore.LookupTable{
		Name:             "kms_table",
		KmsKeyId:         "arn:aws:kms:us-east-1:000000000000:key/unknown",
		EncryptedBody:    []byte("ciphertext"),
		EncryptedDataKey: []byte("wrapped"),
		ContentNonce:     []byte("nonce1234"),
	}

	unwired := &LogsService{accountID: "000000000000"}
	if _, err := unwired.lookupTablePlainBody(encrypted, "us-east-1"); logsErrorCode(err) != "ServiceUnavailableException" {
		t.Fatalf("unwired KMS read: %v", err)
	}

	wired := newDeliveryTestService() // the fake knows no keys
	if _, err := wired.lookupTablePlainBody(encrypted, "us-east-1"); logsErrorCode(err) != "ServiceUnavailableException" {
		t.Fatalf("failed decrypt read: %v", err)
	}
}
