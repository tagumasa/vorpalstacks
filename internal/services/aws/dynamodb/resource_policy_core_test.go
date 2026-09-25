package dynamodb

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"vorpalstacks/internal/common/request"
)

// The ExpectedRevisionId optimistic lock is one locked read-modify-write:
// concurrent writers that all observed the same revision cannot all apply
// — exactly one put and exactly one delete win, the losers answer the
// policy-not-found shape the wire contract maps a moved revision to, and
// the revision advances by exactly the winners' steps.
func TestResourcePolicyRevisionCheckIsAtomic(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	arn := "arn:aws:dynamodb:us-east-1:123456789012:table/LegacyTable"

	first, err := svc.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{ResourceArn: arn, Policy: `{"Version":"2012-10-17"}`})
	if err != nil {
		t.Fatalf("initial put: %v", err)
	}
	if first.RevisionId != "v1" {
		t.Fatalf("initial revision = %s, want v1", first.RevisionId)
	}

	const writers = 8
	var wg sync.WaitGroup
	var winners atomic.Int32
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, perr := svc.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{
				ResourceArn:        arn,
				Policy:             `{"Version":"2012-10-17","w":` + strconv.Itoa(i) + `}`,
				ExpectedRevisionId: "v1",
			})
			if perr == nil {
				winners.Add(1)
			} else if !errors.Is(perr, ErrPolicyNotFound) {
				t.Errorf("losing writer: err = %v, want the policy-not-found shape", perr)
			}
		}(i)
	}
	wg.Wait()
	if n := winners.Load(); n != 1 {
		t.Fatalf("concurrent same-revision puts: %d winners, want exactly 1", n)
	}
	got, err := svc.getResourcePolicyCore(ctx, reqCtx, GetResourcePolicyInput{ResourceArn: arn})
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if got.RevisionId != "v2" {
		t.Fatalf("revision after the race = %s, want v2", got.RevisionId)
	}

	// The delete competes under the same lock — and advances the revision,
	// so the token moves and the second same-revision delete is rejected:
	// deletion is not a free pass past the optimistic lock.
	var delWinners atomic.Int32
	var delWG sync.WaitGroup
	for i := 0; i < 2; i++ {
		delWG.Add(1)
		go func() {
			defer delWG.Done()
			_, derr := svc.deleteResourcePolicyCore(ctx, reqCtx, DeleteResourcePolicyInput{ResourceArn: arn, ExpectedRevisionId: "v2"})
			if derr == nil {
				delWinners.Add(1)
			} else if !errors.Is(derr, ErrPolicyNotFound) {
				t.Errorf("losing delete: err = %v, want the policy-not-found shape", derr)
			}
		}()
	}
	delWG.Wait()
	if n := delWinners.Load(); n != 1 {
		t.Fatalf("concurrent same-revision deletes: %d winners, want exactly 1", n)
	}

	// A put and a delete holding the same expected revision compete under
	// the one lock as well: exactly one applies, the loser is rejected.
	var crossWinners atomic.Int32
	var crossWG sync.WaitGroup
	for i := 0; i < 2; i++ {
		crossWG.Add(1)
		go func(i int) {
			defer crossWG.Done()
			if i == 0 {
				_, perr := svc.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{
					ResourceArn:        arn,
					Policy:             `{"Version":"2012-10-17","cross":"put"}`,
					ExpectedRevisionId: "v3",
				})
				if perr == nil {
					crossWinners.Add(1)
				} else if !errors.Is(perr, ErrPolicyNotFound) {
					t.Errorf("cross put loser: err = %v, want the policy-not-found shape", perr)
				}
				return
			}
			_, derr := svc.deleteResourcePolicyCore(ctx, reqCtx, DeleteResourcePolicyInput{ResourceArn: arn, ExpectedRevisionId: "v3"})
			if derr == nil {
				crossWinners.Add(1)
			} else if !errors.Is(derr, ErrPolicyNotFound) {
				t.Errorf("cross delete loser: err = %v, want the policy-not-found shape", derr)
			}
		}(i)
	}
	crossWG.Wait()
	if n := crossWinners.Load(); n != 1 {
		t.Fatalf("concurrent put-vs-delete at one revision: %d winners, want exactly 1", n)
	}

	// The revision advanced exactly once per winner: v1, then the put-race
	// winner to v2, the delete-race winner to v3, the cross winner to v4.
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	finalRev, err := store.Tables().GetResourcePolicyRevisionId("LegacyTable")
	if err != nil {
		t.Fatalf("final revision read: %v", err)
	}
	if finalRev != 4 {
		t.Fatalf("final revision = %d, want 4 (one advance per race winner)", finalRev)
	}
}

// The delete operation's wire response carries the fresh revision the
// delete minted (DeleteResourcePolicyOutput.RevisionId): the token a later
// put must expect once the policy is gone — never an empty body.
func TestDeleteResourcePolicyResponseCarriesRevision(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	arn := "arn:aws:dynamodb:us-east-1:123456789012:table/LegacyTable"

	if _, err := svc.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{ResourceArn: arn, Policy: `{"Version":"2012-10-17"}`}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	resp, err := svc.DeleteResourcePolicy(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ResourceArn":        arn,
		"ExpectedRevisionId": "v1",
	}})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	revision, ok := resp.(map[string]interface{})["RevisionId"].(string)
	if !ok || revision != "v2" {
		t.Fatalf("delete response revision = %#v, want the fresh v2 the delete minted", resp)
	}
}
