package lambda

// Sandbox pool unit tests run against an in-memory ContainerLifecycle
// fake; the Runtime API servers they exercise are the real ones, so the
// idle signal (a parked /next long-poll) is driven by actual HTTP, exactly
// as a container RIC drives it.

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/client/mobyclient"
)

// fakeContainer is one in-memory container of the fakeLifecycle.
type fakeContainer struct {
	id      string
	cfg     mobyclient.AdvancedContainerConfig
	running bool
	logs    []fakeLogLine
}

type fakeLogLine struct {
	ts   time.Time
	text string
}

// fakeLifecycle is an in-memory ContainerLifecycle: containers are born
// stopped, start marks them running, remove deletes them, and the log
// store answers ReadLogWindow with the net window contract (strictly
// after `since`, per-line nanosecond timestamps). Docker's own since
// filter is inclusive — it re-delivers the record at `since` — and the
// client compensates for that, so the fake models the compensated result
// rather than docker's raw filter.
type fakeLifecycle struct {
	mu         sync.Mutex
	containers map[string]*fakeContainer
	nextID     int
	failCreate bool
}

func newFakeLifecycle() *fakeLifecycle {
	return &fakeLifecycle{containers: make(map[string]*fakeContainer)}
}

func (f *fakeLifecycle) GetContainerStatus(_ context.Context, containerID string) (mobyclient.ContainerStatus, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.containers[containerID]
	if !ok {
		return "", fmt.Errorf("no such container %s", containerID)
	}
	if c.running {
		return mobyclient.ContainerStatusRunning, nil
	}
	return "exited", nil
}

func (f *fakeLifecycle) CreateContainerFromConfig(_ context.Context, cfg mobyclient.AdvancedContainerConfig) (*mobyclient.CreateContainerResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failCreate {
		return nil, fmt.Errorf("injected create failure")
	}
	f.nextID++
	id := fmt.Sprintf("fake-%d", f.nextID)
	f.containers[id] = &fakeContainer{id: id, cfg: cfg}
	return &mobyclient.CreateContainerResult{ID: id}, nil
}

func (f *fakeLifecycle) StartContainer(_ context.Context, containerID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.containers[containerID]
	if !ok {
		return fmt.Errorf("no such container %s", containerID)
	}
	c.running = true
	return nil
}

func (f *fakeLifecycle) RemoveContainer(_ context.Context, containerID string, _ bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.containers, containerID)
	return nil
}

func (f *fakeLifecycle) KillContainer(_ context.Context, containerID string, _ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.containers[containerID]; ok {
		c.running = false
	}
	return nil
}

func (f *fakeLifecycle) ListContainers(_ context.Context, _ bool) ([]mobyclient.ContainerInfo, error) {
	return nil, nil
}

func (f *fakeLifecycle) CreateFileInContainer(_ context.Context, _ string, _ string, _ []byte, _ os.FileMode) error {
	return nil
}

func (f *fakeLifecycle) Exec(_ context.Context, _ string, _ mobyclient.ExecConfig) (*mobyclient.ExecResult, error) {
	return &mobyclient.ExecResult{}, nil
}

func (f *fakeLifecycle) ReadLogWindow(_ context.Context, containerID string, since time.Time) (string, time.Time, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.containers[containerID]
	if !ok {
		return "", since, fmt.Errorf("no such container %s", containerID)
	}
	var b strings.Builder
	cursor := since
	for _, line := range c.logs {
		if !line.ts.After(since) {
			continue
		}
		b.WriteString(line.ts.Format(time.RFC3339Nano))
		b.WriteByte(' ')
		b.WriteString(line.text)
		b.WriteByte('\n')
		if line.ts.After(cursor) {
			cursor = line.ts
		}
	}
	return b.String(), cursor, nil
}

func (f *fakeLifecycle) Close() error { return nil }

// appendLogAt records a log line on a container at a chosen time.
func (f *fakeLifecycle) appendLogAt(containerID string, ts time.Time, text string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.containers[containerID]; ok {
		c.logs = append(c.logs, fakeLogLine{ts: ts, text: text})
	}
}

// countRunning returns the number of containers still present.
func (f *fakeLifecycle) countContainers() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.containers)
}

// parkNext opens a real /next long-poll against the sandbox's Runtime API
// server, as the container RIC does between rounds.
func parkNext(t *testing.T, sb *lambdaSandbox) {
	t.Helper()
	go func() {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", sb.api.Port(), runtimeAPINextPath))
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()
	waitForCondition(t, 2*time.Second, "the sandbox runtime to park on /next", func() bool {
		return sb.api.Idle()
	})
}

func testSpec() sandboxSpec {
	return sandboxSpec{
		imageURI:    "localhost/fake/image:1",
		env:         map[string]string{"AWS_REGION": "us-east-1"},
		memoryBytes: 128 * 1024 * 1024,
	}
}

const testArn = "arn:aws:lambda:us-east-1:123456789012:function:img"

func newTestPool(t *testing.T, f *fakeLifecycle) *sandboxPool {
	t.Helper()
	p := newSandboxPool(f)
	// The idle TTL stays comfortably above the tests' release-acquire
	// gaps; the TTL-specific test tightens it for its own timing.
	p.idleTTL = 250 * time.Millisecond
	p.reapInterval = 10 * time.Millisecond
	p.start()
	t.Cleanup(p.shutdown)
	return p
}

func TestSandboxPoolAcquireCreatesAndReusesIdle(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)

	sb1, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	if sb1.containerID == "" || sb1.api == nil {
		t.Fatalf("acquire must build the container and the runtime api server")
	}
	if f.countContainers() != 1 {
		t.Fatalf("first acquire must create exactly one container, got %d", f.countContainers())
	}

	// A freshly acquired sandbox is busy and not reusable; a parallel
	// acquire must create a second sandbox, not hand out the busy one.
	sb2, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("parallel acquire: %v", err)
	}
	if sb2.id == sb1.id {
		t.Fatalf("a busy sandbox must never be handed to a second invoke")
	}
	if f.countContainers() != 2 {
		t.Fatalf("parallel acquires must create two containers, got %d", f.countContainers())
	}

	// Release sb1 with its runtime parked on /next: the next acquire for
	// the same function version reuses the warm sandbox.
	parkNext(t, sb1)
	p.release(sb1)
	sb3, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("reuse acquire: %v", err)
	}
	if sb3.id != sb1.id || sb3.containerID != sb1.containerID {
		t.Fatalf("an idle parked sandbox must be reused, got a new one (%s vs %s)", sb3.id, sb1.id)
	}
	if f.countContainers() != 2 {
		t.Fatalf("reuse must not create containers, got %d", f.countContainers())
	}

	// A different version never reuses another version's sandbox.
	sb4, err := p.acquire(context.Background(), testArn, "3", nil, testSpec())
	if err != nil {
		t.Fatalf("version acquire: %v", err)
	}
	if sb4.id == sb1.id || sb4.id == sb2.id {
		t.Fatalf("a versioned acquire must not reuse a $LATEST sandbox")
	}
}

func TestSandboxPoolReservedConcurrencyCap(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)
	reserved := int64(1)

	sb1, err := p.acquire(context.Background(), testArn, "$LATEST", &reserved, testSpec())
	if err != nil {
		t.Fatalf("first acquire under reserve: %v", err)
	}
	if _, err := p.acquire(context.Background(), testArn, "$LATEST", &reserved, testSpec()); err != ErrTooManyRequests {
		t.Fatalf("acquire at the reserve cap = %v, want ErrTooManyRequests", err)
	}

	// A busy sandbox blocks; once released and parked, the same reserve
	// slot serves the next invoke by reuse.
	parkNext(t, sb1)
	p.release(sb1)
	sb2, err := p.acquire(context.Background(), testArn, "$LATEST", &reserved, testSpec())
	if err != nil {
		t.Fatalf("reuse acquire under reserve: %v", err)
	}
	if sb2.id != sb1.id {
		t.Fatalf("the reserved slot must be reused, got %s want %s", sb2.id, sb1.id)
	}
}

func TestSandboxPoolGlobalLimitCap(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)
	p.globalLimit = 1

	otherArn := "arn:aws:lambda:us-east-1:123456789012:function:other"
	if _, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec()); err != nil {
		t.Fatalf("first global acquire: %v", err)
	}
	// A reserved function draws against its own reserve, not the global
	// pool — the AWS reserved/unreserved split.
	reserved := int64(4)
	if _, err := p.acquire(context.Background(), otherArn, "$LATEST", &reserved, testSpec()); err != nil {
		t.Fatalf("reserved acquire at the global cap: %v", err)
	}
	if _, err := p.acquire(context.Background(), otherArn, "$LATEST", nil, testSpec()); err != ErrTooManyRequests {
		t.Fatalf("unreserved acquire at the global cap = %v, want ErrTooManyRequests", err)
	}
}

func TestSandboxPoolTTLReapsIncludingLast(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)
	p.idleTTL = 50 * time.Millisecond

	sb, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	parkNext(t, sb)
	p.release(sb)

	// The idle TTL reaps every idle sandbox — the last one included; the
	// platform keeps none permanently warm.
	waitForCondition(t, 2*time.Second, "the idle sandbox to be reaped by the TTL", func() bool {
		return f.countContainers() == 0
	})
	p.mu.Lock()
	remaining := len(p.sandboxes)
	p.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("the pool must be empty after the TTL reaped the last sandbox, got %d", remaining)
	}
}

func TestSandboxPoolReapsDeadUnparkedSandbox(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)

	sb, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	// Released without a parked /next — the runtime never looped back —
	// and its container died: the reaper must free the slot without
	// waiting for the idle TTL.
	p.release(sb)
	_ = f.KillContainer(context.Background(), sb.containerID, "SIGKILL")
	_ = f.RemoveContainer(context.Background(), sb.containerID, true)

	waitForCondition(t, 2*time.Second, "the dead sandbox to be reaped", func() bool {
		p.mu.Lock()
		n := len(p.sandboxes)
		p.mu.Unlock()
		return n == 0
	})
}

func TestSandboxPoolLogWindowPartitionsWithoutLoss(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)

	sb, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}

	t0 := time.Now().Add(-time.Minute)
	t1 := t0.Add(10 * time.Second)
	t2 := t0.Add(20 * time.Second)
	f.appendLogAt(sb.containerID, t0, "init phase output")
	f.appendLogAt(sb.containerID, t1, "round one output")

	// The first window starts at the container's beginning, so the init
	// output attaches to the first invoke.
	first, err := p.readLogWindow(context.Background(), sb)
	if err != nil {
		t.Fatalf("first readLogWindow: %v", err)
	}
	if !strings.Contains(first, "init phase output") || !strings.Contains(first, "round one output") {
		t.Fatalf("first window = %q, want both init and round-one output", first)
	}

	f.appendLogAt(sb.containerID, t2, "round two output")
	second, err := p.readLogWindow(context.Background(), sb)
	if err != nil {
		t.Fatalf("second readLogWindow: %v", err)
	}
	if !strings.Contains(second, "round two output") {
		t.Fatalf("second window = %q, want the round-two output", second)
	}
	if strings.Contains(second, "round one output") || strings.Contains(second, "init phase output") {
		t.Fatalf("second window = %q, must not re-deliver already-read lines", second)
	}
}

func TestSandboxPoolDrainSparesBusyUntilRelease(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)

	sbBusy, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("busy acquire: %v", err)
	}
	sbIdle, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("idle acquire: %v", err)
	}
	sbVer, err := p.acquire(context.Background(), testArn, "7", nil, testSpec())
	if err != nil {
		t.Fatalf("version acquire: %v", err)
	}
	parkNext(t, sbIdle)
	p.release(sbIdle)

	// A $LATEST code update drains only $LATEST: the idle one goes at
	// once, the busy one survives its in-flight invoke and dies on
	// release, and the published version's sandbox is untouched.
	p.drainVersion(testArn, "$LATEST")
	if _, err := f.GetContainerStatus(context.Background(), sbIdle.containerID); err == nil {
		t.Fatalf("the idle $LATEST sandbox's container must be removed by drainVersion")
	}
	if _, err := f.GetContainerStatus(context.Background(), sbBusy.containerID); err != nil {
		t.Fatalf("the busy sandbox's container must survive its in-flight invoke")
	}
	p.release(sbBusy)
	if _, err := f.GetContainerStatus(context.Background(), sbBusy.containerID); err == nil {
		t.Fatalf("the doomed sandbox must be destroyed on release")
	}
	if _, err := f.GetContainerStatus(context.Background(), sbVer.containerID); err != nil {
		t.Fatalf("a published version's sandbox must not be drained by a $LATEST update")
	}

	// drainFunction sweeps every version.
	parkNext(t, sbVer)
	p.release(sbVer)
	p.drainFunction(testArn)
	if f.countContainers() != 0 {
		t.Fatalf("drainFunction must remove every version's sandbox, got %d containers", f.countContainers())
	}
}

func TestSandboxPoolShutdownDestroysEverything(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)

	if _, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec()); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if _, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec()); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	p.shutdown()
	if f.countContainers() != 0 {
		t.Fatalf("shutdown must destroy every sandbox including busy ones, got %d", f.countContainers())
	}
}

func TestSandboxPoolCreateFailureFreesSlot(t *testing.T) {
	f := newFakeLifecycle()
	p := newTestPool(t, f)
	p.globalLimit = 1

	_, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec())
	if err != nil {
		t.Fatalf("initial acquire: %v", err)
	}
	// Detach the first sandbox from the pool so the next acquire is not
	// capped by it; only the create-failure path is under test here.
	p.mu.Lock()
	p.sandboxes = map[string]*lambdaSandbox{}
	p.mu.Unlock()

	f.failCreate = true
	if _, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec()); err == nil {
		t.Fatalf("an injected create failure must surface")
	}
	f.failCreate = false
	// The failed reservation must have freed the global slot.
	if _, err := p.acquire(context.Background(), testArn, "$LATEST", nil, testSpec()); err != nil {
		t.Fatalf("acquire after a failed create must find the slot free: %v", err)
	}
}
