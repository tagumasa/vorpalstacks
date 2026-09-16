package eventbridge

import (
	"strings"
	"testing"

	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestTagResourceLimits pins the documented tag limits on the tag path:
// at most 50 tags per resource, keys of up to 128 Unicode characters,
// values of up to 256, and the aws: prefix reserved (eb-tagging page,
// fetched 2026-09-15).
func TestTagResourceLimits(t *testing.T) {
	svc := newEventBusCoreTestService()
	ctx := t.Context()

	cases := []struct {
		name    string
		tags    []tagutil.Tag
		wantErr string
	}{
		{
			name:    "51st tag rejected",
			tags:    makeTags(51, "k", "v"),
			wantErr: "too many tags",
		},
		{
			name:    "50 tags accepted",
			tags:    makeTags(50, "k", "v"),
			wantErr: "",
		},
		{
			name:    "aws: prefixed key rejected",
			tags:    []tagutil.Tag{{Key: "aws:reserved", Value: "v"}},
			wantErr: "reserved prefix",
		},
		{
			name:    "129-character key rejected",
			tags:    []tagutil.Tag{{Key: strings.Repeat("k", 129), Value: "v"}},
			wantErr: "cannot exceed 128 characters",
		},
		{
			name:    "128-character key accepted",
			tags:    []tagutil.Tag{{Key: strings.Repeat("k", 128), Value: "v"}},
			wantErr: "",
		},
		{
			name:    "257-character value rejected",
			tags:    []tagutil.Tag{{Key: "k", Value: strings.Repeat("v", 257)}},
			wantErr: "cannot exceed 256 characters",
		},
	}
	for _, tc := range cases {
		// A fresh bus per case: the 50-tag ceiling counts total tags on
		// the resource, so accepted cases must not accumulate onto the
		// next case's bus.
		store := newEventBusCoreTestStore(t)
		if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "tagged-bus"}); err != nil {
			t.Fatalf("create bus: %v", err)
		}
		busARN := BuildEventBusARN("000000000000", "us-east-1", "tagged-bus")

		err := svc.tagResourceCore(ctx, store, busARN, tc.tags)
		if tc.wantErr == "" {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("%s: accepted, want rejection containing %q", tc.name, tc.wantErr)
			continue
		}
		if !strings.Contains(err.Error(), tc.wantErr) {
			t.Errorf("%s: error = %q, want it to contain %q", tc.name, err.Error(), tc.wantErr)
		}
	}
}

// TestTaggableScopeIsBusAndRuleOnly pins the taggable-scope decision: the
// tagging page enumerates exactly rules and event buses ("In EventBridge,
// you can assign tags to rule and event buses", fetched 2026-09-15), so
// archive, connection and API-destination ARNs are not taggable — they
// answer with the resource-not-found shape even when the resource exists.
func TestTaggableScopeIsBusAndRuleOnly(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := t.Context()

	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{Name: "arch", EventBusName: "default"}); err != nil {
		t.Fatalf("create archive: %v", err)
	}

	arns := []string{
		"arn:aws:events:us-east-1:000000000000:archive/arch",
		"arn:aws:events:us-east-1:000000000000:connection/conn",
		"arn:aws:events:us-east-1:000000000000:api-destination/dest",
	}
	for _, arn := range arns {
		err := svc.tagResourceCore(ctx, store, arn, []tagutil.Tag{{Key: "k", Value: "v"}})
		if err == nil {
			t.Errorf("tagging %s accepted, want resource-not-found", arn)
			continue
		}
		if !strings.Contains(err.Error(), "ResourceNotFoundException") {
			t.Errorf("tagging %s: error = %v, want the resource-not-found shape", arn, err)
		}
	}
}

func makeTags(n int, key, value string) []tagutil.Tag {
	tags := make([]tagutil.Tag, n)
	for i := range tags {
		tags[i] = tagutil.Tag{Key: key, Value: value}
		if n > 1 {
			tags[i].Key = key + string(rune('a'+i%26)) + string(rune('a'+i/26))
		}
	}
	return tags
}

// TestTagRequiredMembersCoreOwned pins that the tag family's required-member
// rejections live in the Core functions, not the handler plane: the
// ResourceARN requirement in the shared resolver, Tags in tagResourceCore
// and TagKeys in untagResourceCore — so the admin plane calling the Cores
// sees the same contract as the HTTP plane.
func TestTagRequiredMembersCoreOwned(t *testing.T) {
	svc := newEventBusCoreTestService()
	store := newEventBusCoreTestStore(t)
	ctx := t.Context()
	busARN := BuildEventBusARN("000000000000", "us-east-1", "req-bus")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "req-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	cases := []struct {
		name    string
		call    func() error
		wantErr string
	}{
		{"empty ResourceARN on tag", func() error {
			return svc.tagResourceCore(ctx, store, "", []tagutil.Tag{{Key: "k", Value: "v"}})
		}, "ResourceARN is required"},
		{"empty Tags", func() error {
			return svc.tagResourceCore(ctx, store, busARN, nil)
		}, "Tags are required"},
		{"empty ResourceARN on untag", func() error {
			return svc.untagResourceCore(ctx, store, "", map[string]bool{"k": true})
		}, "ResourceARN is required"},
		{"empty TagKeys", func() error {
			return svc.untagResourceCore(ctx, store, busARN, map[string]bool{})
		}, "TagKeys are required"},
		{"empty ResourceARN on list", func() error {
			_, err := svc.listTagsForResourceCore(ctx, store, "")
			return err
		}, "ResourceARN is required"},
	}
	for _, tc := range cases {
		err := tc.call()
		if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
			t.Errorf("%s: got %v, want ValidationException containing %q", tc.name, err, tc.wantErr)
		}
	}
}
