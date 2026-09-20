package apps

// Pin for the eventbus SNS topic-policy store resolution: the topic ARN
// names its region, and the policy evaluation must read that region's
// store through the service's region resolver — a cross-region topic
// evaluates its own policy, not an empty lookup in the default region's
// store.

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
	storesns "vorpalstacks/internal/store/aws/sns"
)

// fakeRegionResolver serves a fixed region→store map and fails for every
// other region — the shape the service's resolver has when a region has
// no store yet.
type fakeRegionResolver struct {
	stores map[string]storesns.SNSStoreInterface
}

func (f *fakeRegionResolver) GetSNSStoreForRegion(region string) (storesns.SNSStoreInterface, error) {
	if store, ok := f.stores[region]; ok {
		return store, nil
	}
	return nil, errors.New("no store for region")
}

// markerPolicyStore answers GetTopic with a policy naming the store that
// served it, so the test can tell which regional store a resolution chose.
type markerPolicyStore struct {
	storesns.SNSStoreInterface
	marker string
}

func (m *markerPolicyStore) GetTopic(topicARN string) (*storesns.Topic, error) {
	return &storesns.Topic{Arn: topicARN, Attributes: map[string]string{"Policy": m.marker}}, nil
}

func TestResolveSNSTopicPolicyStoreServesTheARNRegion(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	fallback := storesns.NewSNSStore(st, "123456789012", "us-east-1")
	fallbackTopic, err := fallback.CreateTopic(&storesns.Topic{Name: "fallback-topic"}, nil)
	if err != nil {
		t.Fatalf("seed fallback topic: %v", err)
	}
	if err := fallback.SetTopicAttributes(fallbackTopic.Arn, map[string]string{"Policy": "fallback-policy"}); err != nil {
		t.Fatalf("seed fallback policy: %v", err)
	}

	regional := &markerPolicyStore{marker: "regional-policy"}
	service := &fakeRegionResolver{stores: map[string]storesns.SNSStoreInterface{"eu-west-1": regional}}

	// A cross-region topic ARN resolves its own region's store.
	topic, err := resolveSNSTopicPolicyStore(fallback, service, "arn:aws:sns:eu-west-1:123456789012:topic").GetTopic("arn:aws:sns:eu-west-1:123456789012:topic")
	if err != nil {
		t.Fatalf("cross-region GetTopic: %v", err)
	}
	if topic.Attributes["Policy"] != "regional-policy" {
		t.Fatalf("cross-region policy = %q, want the regional store's marker", topic.Attributes["Policy"])
	}

	// A region the resolver cannot serve falls back to the default
	// instance: the ARN's topic is absent there and the lookup answers
	// not-found — never the regional store's marker.
	if _, err := resolveSNSTopicPolicyStore(fallback, service, "arn:aws:sns:ap-southeast-2:123456789012:topic").GetTopic("arn:aws:sns:ap-southeast-2:123456789012:topic"); !errors.Is(err, storesns.ErrTopicNotFound) {
		t.Fatalf("unresolvable region error = %v, want the fallback store's not-found", err)
	}

	// Service-less wiring reads the default instance.
	topic, err = resolveSNSTopicPolicyStore(fallback, nil, fallbackTopic.Arn).GetTopic(fallbackTopic.Arn)
	if err != nil {
		t.Fatalf("service-less GetTopic: %v", err)
	}
	if topic.Attributes["Policy"] != "fallback-policy" {
		t.Fatalf("service-less policy = %q, want the fallback store's policy", topic.Attributes["Policy"])
	}
}
