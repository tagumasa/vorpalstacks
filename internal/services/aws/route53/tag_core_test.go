package route53

import (
	"errors"
	"fmt"
	"os"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	route53store "vorpalstacks/internal/store/aws/route53"
)

func newTagCoreTestService(t *testing.T) (*Route53Service, *route53store.Route53Stores) {
	t.Helper()
	dir := "./tmp/tag-core-test-" + t.Name()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() {
		st.Close()
		os.RemoveAll(dir)
	})
	svc := NewRoute53Service(st, "123456789012")
	return svc, svc.dnsStores
}

// TestChangeTagsForResourceCoreMergedOverflowAnswersInvalidInput pins the
// merged-total overflow identity: a zone already carrying eight tags that
// receives five more (thirteen against the API's ten-tag quota) answers
// the operation's declared InvalidInput with a client status — the same
// condition family the per-request check covers — never a server-error
// wrap with an invented code.
func TestChangeTagsForResourceCoreMergedOverflowAnswersInvalidInput(t *testing.T) {
	_, stores := newTagCoreTestService(t)

	zone := &route53store.HostedZone{
		ID:              "Z1234567890ABC",
		Name:            "example.com.",
		CallerReference: "tag-core-test",
		Region:          "us-east-1",
		AccountID:       "123456789012",
	}
	if err := stores.HostedZones().Create(zone); err != nil {
		t.Fatalf("create zone: %v", err)
	}

	eight := make([]tags.Tag, 8)
	for i := range eight {
		eight[i] = tags.Tag{Key: fmt.Sprintf("seed-%d", i), Value: "v"}
	}
	if err := changeTagsForResourceCore(stores, ChangeTagsForResourceInput{
		ResourceType: "hostedzone",
		ResourceId:   zone.ID,
		AddTags:      eight,
	}); err != nil {
		t.Fatalf("seed eight tags: %v", err)
	}

	five := make([]tags.Tag, 5)
	for i := range five {
		five[i] = tags.Tag{Key: fmt.Sprintf("extra-%d", i), Value: "v"}
	}
	err := changeTagsForResourceCore(stores, ChangeTagsForResourceInput{
		ResourceType: "hostedzone",
		ResourceId:   zone.ID,
		AddTags:      five,
	})

	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) {
		t.Fatalf("merged overflow: want an AWS error, got %v", err)
	}
	if awsErr.Code != "InvalidInput" {
		t.Fatalf("merged overflow code: got %q, want InvalidInput", awsErr.Code)
	}
	if awsErr.HTTPStatus != 400 {
		t.Fatalf("merged overflow status: got %d, want 400", awsErr.HTTPStatus)
	}

	// The rejection left the tag set at its seeded size.
	list, err := stores.Tags().ListTagsForResource("hostedzone/" + zone.ID)
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	if len(list) != 8 {
		t.Fatalf("tag set after rejected overflow: %d entries, want the 8 seeded", len(list))
	}
}
