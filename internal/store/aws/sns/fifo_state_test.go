package sns

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestFifoStateSurvivesRestart pins the durability decision: the
// deduplication window and the per-group sequence counters live in storage,
// so a fresh store over the same data directory — the process-restart
// equivalent — still hits the window with the original publish's
// identifiers and keeps the group's SequenceNumber increasing.
func TestFifoStateSurvivesRestart(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:restart.fifo"

	first, err := st.AllocateFifoSequence(topicArn, "group-a")
	if err != nil {
		t.Fatalf("first allocation: %v", err)
	}
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "dedup-1", "message-1", first, false); err != nil {
		t.Fatalf("record dedup: %v", err)
	} else if hit {
		t.Fatal("first record reported a hit on an empty window")
	}

	// A restart is a fresh store instance over the same durable buckets.
	st.Close()
	restored := NewSNSStore(st.storage, "123456789012", "us-east-1")
	t.Cleanup(restored.Close)

	messageID, sequenceNumber, hit, err := restored.CheckFifoDeduplication(topicArn, "group-a", "dedup-1", false)
	if err != nil {
		t.Fatalf("restored check: %v", err)
	}
	if !hit {
		t.Fatal("the deduplication window was lost across the restart")
	}
	if messageID != "message-1" || sequenceNumber != first {
		t.Fatalf("restored entry = (%q, %q), want the original publish's (message-1, %s)", messageID, sequenceNumber, first)
	}

	second, err := restored.AllocateFifoSequence(topicArn, "group-a")
	if err != nil {
		t.Fatalf("restored allocation: %v", err)
	}
	firstVal, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		t.Fatalf("first sequence %q: %v", first, err)
	}
	secondVal, err := strconv.ParseInt(second, 10, 64)
	if err != nil {
		t.Fatalf("second sequence %q: %v", second, err)
	}
	if secondVal <= firstVal {
		t.Fatalf("SequenceNumber regressed across the restart: %s then %s", first, second)
	}
}

// TestCheckAndRecordFifoDeduplicationExactlyOneWinner pins the atomicity of
// the durable check-and-record: concurrent publishers with the same
// deduplication identifier admit exactly one copy — one goroutine records,
// every other goroutine observes the winner's identifiers.
func TestCheckAndRecordFifoDeduplicationExactlyOneWinner(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:winner.fifo"

	const publishers = 8
	type publisherResult struct {
		winner    bool
		messageID string
		sequence  string
	}
	results := make(chan publisherResult, publishers)
	var start sync.WaitGroup
	start.Add(1)
	var wg sync.WaitGroup
	for i := 0; i < publishers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start.Wait()
			messageID := fmt.Sprintf("message-%d", i)
			sequenceNumber := fmt.Sprintf("seq-%d", i)
			existingID, existingSeq, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "same-dedup", messageID, sequenceNumber, false)
			if err != nil {
				t.Errorf("publisher %d: %v", i, err)
				return
			}
			if hit {
				results <- publisherResult{winner: false, messageID: existingID, sequence: existingSeq}
				return
			}
			results <- publisherResult{winner: true, messageID: messageID, sequence: sequenceNumber}
		}(i)
	}
	start.Done()
	wg.Wait()
	close(results)

	winners := 0
	var winner publisherResult
	for result := range results {
		if result.winner {
			winners++
			winner = result
		}
	}
	if winners != 1 {
		t.Fatalf("exactly one publisher must record the dedup entry, got %d", winners)
	}
	messageID, sequenceNumber, hit, err := st.CheckFifoDeduplication(topicArn, "group-a", "same-dedup", false)
	if err != nil || !hit {
		t.Fatalf("post-race check = (hit %v, err %v), want the recorded entry", hit, err)
	}
	if messageID != winner.messageID || sequenceNumber != winner.sequence {
		t.Fatalf("recorded entry = (%q, %q), want the winner's (%q, %q)", messageID, sequenceNumber, winner.messageID, winner.sequence)
	}
}

// TestExpiredDedupEntryReadsAsMiss pins the window boundary: an entry whose
// expiry has passed reads as a miss even before the reclamation sweep
// deletes it.
func TestExpiredDedupEntryReadsAsMiss(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:expired.fifo"
	key := []byte(topicArn + ":stale")

	entry, err := json.Marshal(&fifoDedupRecord{
		MessageID:      "old-message",
		SequenceNumber: "1000000000000000001",
		ExpiresAt:      time.Now().Add(-time.Minute),
	})
	if err != nil {
		t.Fatalf("marshal stale entry: %v", err)
	}
	if err := st.fifoDedupStore.Put(key, entry); err != nil {
		t.Fatalf("seed stale entry: %v", err)
	}

	if _, _, hit, err := st.CheckFifoDeduplication(topicArn, "group-a", "stale", false); err != nil {
		t.Fatalf("check: %v", err)
	} else if hit {
		t.Fatal("an expired dedup entry read as a hit")
	}
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "stale", "new-message", "1000000000000000002", false); err != nil {
		t.Fatalf("record over stale entry: %v", err)
	} else if hit {
		t.Fatal("check-and-record treated the expired entry as a live duplicate")
	}
}

// TestPerGroupDedupScope pins the FifoThroughputScope=MessageGroup scope:
// "message deduplication applies to each individual message group" — the
// same deduplication ID in two message groups records twice (each group
// delivers its own copy), while the topic scope admits only the first.
func TestPerGroupDedupScope(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:scope.fifo"

	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "shared-dedup", "message-a", "1000000000000000001", true); err != nil || hit {
		t.Fatalf("per-group record group-a (hit %v, err %v)", hit, err)
	}
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-b", "shared-dedup", "message-b", "1000000000000000002", true); err != nil {
		t.Fatalf("per-group record group-b: %v", err)
	} else if hit {
		t.Fatal("a second message group must not hit the first group's per-group window")
	}

	// The topic scope keys the window without the group: the same
	// identifier in another group hits.
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "topic-dedup", "message-1", "1000000000000000003", false); err != nil || hit {
		t.Fatalf("topic-scope record (hit %v, err %v)", hit, err)
	}
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-b", "topic-dedup", "message-2", "1000000000000000004", false); err != nil {
		t.Fatalf("topic-scope cross-group record: %v", err)
	} else if !hit {
		t.Fatal("the topic scope must admit only the first copy across groups")
	}
}

// TestFifoDedupKeyScopeSeparator pins the key's collision safety: the
// identifiers may contain colons (the documented punctuation set includes
// one), so the scope encoding must not be a flat colon concatenation — a
// topic-scope window for the id "G:D" must not collide with a group-scope
// window for the pair (G, D), a collision the switchable
// FifoThroughputScope would turn into a false dedup hit.
func TestFifoDedupKeyScopeSeparator(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:separator.fifo"

	// A group-scope window for (G, D)...
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "G", "D", "message-g", "1000000000000000001", true); err != nil || hit {
		t.Fatalf("group-scope record (hit %v, err %v)", hit, err)
	}
	// ...must not satisfy a topic-scope check for the id "G:D".
	if _, _, hit, err := st.CheckFifoDeduplication(topicArn, "any-group", "G:D", false); err != nil {
		t.Fatalf("topic-scope check: %v", err)
	} else if hit {
		t.Fatal("topic-scope window for the id G:D collided with the group-scope window for (G, D)")
	}
	// The mirrored pair on its own topic (so the arm-1 window cannot
	// satisfy it): a topic-scope window for "G:D"...
	mirrorTopic := "arn:aws:sns:us-east-1:123456789012:separator-mirror.fifo"
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(mirrorTopic, "any-group", "G:D", "message-t", "1000000000000000002", false); err != nil || hit {
		t.Fatalf("topic-scope record (hit %v, err %v)", hit, err)
	}
	// ...must not satisfy a group-scope check for (G, D) either.
	if _, _, hit, err := st.CheckFifoDeduplication(mirrorTopic, "G", "D", true); err != nil {
		t.Fatalf("group-scope check: %v", err)
	} else if hit {
		t.Fatal("group-scope window for (G, D) collided with the topic-scope window for the id G:D")
	}

	// Two group pairs the flat concatenation also conflates — ("a:b", "c")
	// versus ("a", "b:c") — stay distinct windows.
	pairTopic := "arn:aws:sns:us-east-1:123456789012:separator-pair.fifo"
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(pairTopic, "a:b", "c", "message-p", "1000000000000000003", true); err != nil || hit {
		t.Fatalf("pair record a:b/c (hit %v, err %v)", hit, err)
	}
	if _, _, hit, err := st.CheckFifoDeduplication(pairTopic, "a", "b:c", true); err != nil {
		t.Fatalf("pair check a/b:c: %v", err)
	} else if hit {
		t.Fatal("group pairs (a:b, c) and (a, b:c) collided into one window")
	}
}

// TestDeleteTopicSweepsFifoState pins the lifecycle: deleting a topic
// removes its deduplication entries and sequence counters with the record,
// so a re-created topic name starts from a clean window.
func TestDeleteTopicSweepsFifoState(t *testing.T) {
	st := newTestSNSStore(t)
	topic, err := st.CreateTopic(&Topic{Name: "sweep.fifo"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	if _, err := st.AllocateFifoSequence(topic.Arn, "group-a"); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topic.Arn, "group-a", "dedup-1", "message-1", "1000000000000000001", false); err != nil || hit {
		t.Fatalf("record dedup (hit %v, err %v)", hit, err)
	}

	if err := st.DeleteTopic(topic.Arn); err != nil {
		t.Fatalf("delete topic: %v", err)
	}

	if got := st.fifoDedupStore.Count(); got != 0 {
		t.Fatalf("dedup bucket holds %d entries after DeleteTopic, want 0", got)
	}
	if got := st.fifoSequenceStore.Count(); got != 0 {
		t.Fatalf("sequence bucket holds %d counters after DeleteTopic, want 0", got)
	}
}

// TestDedupSweepReclaimsOnlyClosedWindows pins the sweep's end-to-end
// behaviour: an entry whose window closed is reclaimed, a live window
// survives the same sweep.
func TestDedupSweepReclaimsOnlyClosedWindows(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:sweep-reclaim.fifo"

	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "live-dedup", "message-live", "1000000000000000001", false); err != nil || hit {
		t.Fatalf("record live entry (hit %v, err %v)", hit, err)
	}
	expiredKey := []byte(fifoDedupKey(topicArn, "group-a", "expired-dedup", false))
	expiredValue, err := json.Marshal(&fifoDedupRecord{
		MessageID:      "message-expired",
		SequenceNumber: "1000000000000000002",
		ExpiresAt:      time.Now().Add(-time.Minute),
	})
	if err != nil {
		t.Fatalf("marshal expired record: %v", err)
	}
	if err := st.fifoDedupStore.Put(expiredKey, expiredValue); err != nil {
		t.Fatalf("seed expired record: %v", err)
	}

	st.reclaimDeduplicationEntries(st.scanExpiredDeduplicationEntries())

	if raw, err := st.fifoDedupStore.Get(expiredKey); err != nil || raw != nil {
		t.Fatalf("closed-window entry survived the sweep (raw=%v, err=%v)", raw, err)
	}
	if _, _, hit, err := st.CheckFifoDeduplication(topicArn, "group-a", "live-dedup", false); err != nil || !hit {
		t.Fatalf("live window after sweep: hit=%v err=%v — the sweep reclaimed a live entry", hit, err)
	}
}

// TestDedupSweepRecheckSparesReRecordedEntry pins the reclamation re-check
// under the deduplication mutex: a key the advisory scan collected as
// expired, then re-recorded with a fresh window by a racing publisher
// before the delete, must survive — deleting it would breach the window
// and deliver a duplicate publish's message twice.
func TestDedupSweepRecheckSparesReRecordedEntry(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:sweep-recheck.fifo"
	liveKey := []byte(fifoDedupKey(topicArn, "group-a", "re-recorded", false))

	if _, _, hit, err := st.CheckAndRecordFifoDeduplication(topicArn, "group-a", "re-recorded", "message-fresh", "1000000000000000001", false); err != nil || hit {
		t.Fatalf("record fresh entry (hit %v, err %v)", hit, err)
	}

	// The scan happened before the re-record in the racing interleaving, so
	// its collected list names what WAS expired — here the now-live key.
	st.reclaimDeduplicationEntries([][]byte{liveKey})

	if _, _, hit, err := st.CheckFifoDeduplication(topicArn, "group-a", "re-recorded", false); err != nil || !hit {
		t.Fatalf("re-recorded entry after reclamation: hit=%v err=%v — the sweep deleted a freshly re-recorded window", hit, err)
	}
}

// TestStoreCloseIsIdempotent pins the safety property the service's
// Close-sweep relies on: the shutdown wiring closes the default-region
// instance after the service already closed every cached store, so a
// second Close must be harmless — the cancel and the drained WaitGroup
// both tolerate repeats.
func TestStoreCloseIsIdempotent(t *testing.T) {
	st := newTestSNSStore(t)
	st.Close()
	st.Close()
}
