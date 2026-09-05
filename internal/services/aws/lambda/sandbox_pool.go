package lambda

// This file carries the sandbox pool: the persistent execution
// environments of image-package functions. One sandbox is one container
// running the image's own ENTRYPOINT as PID1, paired with a host-side
// Runtime API server. The AWS execution-environment model maps directly —
// one sandbox serves one invoke at a time, concurrency is the sandbox
// count, and an idle sandbox is one whose runtime is parked on /next
// waiting for the next event.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Per-sandbox log rotation bounds: a chatty warm sandbox must not grow the
// json-file log without limit, and the daemon-side files disappear with
// the container at destroy.
const (
	sandboxLogDriver   = "json-file"
	sandboxLogMaxSize  = "10m"
	sandboxLogMaxFiles = "5"
)

// reapInterval is the default reaper cadence; well under the idle TTL so a
// sandbox whose runtime never returned to /next frees its concurrency slot
// quickly. Tests tighten it directly on the pool.
const reapInterval = 10 * time.Second

// lambdaSandbox is one persistent execution environment. Every mutable
// field is guarded by the owning pool's mutex; the Runtime API server
// carries its own.
type lambdaSandbox struct {
	id            string
	functionArn   string // unqualified function ARN
	version       string // "$LATEST" or a numeric version
	containerName string
	containerID   string
	api           *runtimeAPIServer

	createdAt  time.Time
	lastUsedAt time.Time
	logCursor  time.Time // log watermark: everything at or before it was read

	busy   bool // owned by an in-flight invoke
	doomed bool // drain marked; destroyed when its invoke releases it
	dead   bool // destroy ran; idempotence flag
}

// reusable reports whether the router may hand this sandbox to an invoke:
// free, not drain-marked, and with its runtime parked on /next — the
// signal that the previous round finished and the runtime asked for the
// next event. A container being Running is not enough: an ENTRYPOINT that
// wraps the runtime without speaking the Runtime API loop never parks.
func (sb *lambdaSandbox) reusable() bool {
	return !sb.busy && !sb.doomed && !sb.dead && sb.api != nil && sb.api.Idle()
}

// sandboxSpec is everything a sandbox bakes in at creation; configuration
// changes drain the affected sandboxes instead of mutating them.
type sandboxSpec struct {
	imageURI    string
	env         map[string]string // includes AWS_LAMBDA_RUNTIME_API
	memoryBytes int64
	// ImageConfig overrides; zero values keep the image's own config.
	entrypoint []string
	cmd        []string
	workingDir string
}

// sandboxPool owns every image-function sandbox: routing, concurrency
// caps, the idle reaper, and destruction on lifecycle events.
type sandboxPool struct {
	docker mobyclient.ContainerLifecycle

	mu        sync.Mutex
	sandboxes map[string]*lambdaSandbox

	globalLimit  int
	idleTTL      time.Duration
	reapInterval time.Duration

	stop chan struct{}
	done chan struct{}
}

// newSandboxPool creates the pool with the platform-default bounds.
func newSandboxPool(docker mobyclient.ContainerLifecycle) *sandboxPool {
	return &sandboxPool{
		docker:       docker,
		sandboxes:    make(map[string]*lambdaSandbox),
		globalLimit:  lambdastore.DefaultGlobalSandboxLimit,
		idleTTL:      lambdastore.DefaultSandboxIdleTTL,
		reapInterval: reapInterval,
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
	}
}

// start launches the idle reaper. Call once after construction.
func (p *sandboxPool) start() {
	go p.reapLoop()
}

// shutdown stops the reaper and destroys every sandbox, in-flight ones
// included — the server is going away and their containers must not
// outlive it (the lambda- name prefix also catches them at next startup).
// Safe to call more than once.
func (p *sandboxPool) shutdown() {
	select {
	case <-p.stop:
		return
	default:
	}
	close(p.stop)
	<-p.done
	p.mu.Lock()
	victims := make([]*lambdaSandbox, 0, len(p.sandboxes))
	for _, sb := range p.sandboxes {
		victims = append(victims, sb)
	}
	p.sandboxes = make(map[string]*lambdaSandbox)
	p.mu.Unlock()
	for _, sb := range victims {
		p.teardown(sb)
	}
}

// acquire returns a sandbox for the function version, reusing an idle one
// when its runtime is parked on /next, else creating a new one. reserved
// is the function's ReservedConcurrency: a set value caps this function's
// sandbox count (the reserved/unreserved concurrency split of AWS — a
// reserved function draws against its own reserve, not the shared pool);
// nil draws against the global limit. At the cap the acquire fails with
// ErrTooManyRequests — the AWS throttling contract rejects instead of
// queueing a synchronous invoke.
func (p *sandboxPool) acquire(ctx context.Context, functionArn, version string, reserved *int64, spec sandboxSpec) (*lambdaSandbox, error) {
	now := time.Now()

	p.mu.Lock()
	for _, sb := range p.sandboxes {
		if sb.functionArn == functionArn && sb.version == version && sb.reusable() {
			sb.busy = true
			sb.lastUsedAt = now
			p.mu.Unlock()
			return sb, nil
		}
	}
	if err := p.checkCapLocked(functionArn, reserved); err != nil {
		p.mu.Unlock()
		return nil, err
	}
	// Reserve the slot before the slow docker calls: the placeholder counts
	// toward the caps while it is being built, so a burst cannot overshoot.
	id, err := newSandboxID()
	if err != nil {
		p.mu.Unlock()
		return nil, err
	}
	sb := &lambdaSandbox{
		id:            id,
		functionArn:   functionArn,
		version:       version,
		containerName: "lambda-sb-" + id,
		createdAt:     now,
		lastUsedAt:    now,
		busy:          true,
	}
	p.sandboxes[sb.id] = sb
	p.mu.Unlock()

	if err := p.startSandbox(ctx, sb, spec); err != nil {
		p.mu.Lock()
		delete(p.sandboxes, sb.id)
		p.mu.Unlock()
		return nil, err
	}
	return sb, nil
}

// checkCapLocked enforces the sandbox caps; the caller holds the mutex.
func (p *sandboxPool) checkCapLocked(functionArn string, reserved *int64) error {
	if reserved != nil {
		if *reserved <= 0 {
			// ReservedConcurrency zero pauses the function; the invoke
			// plane's function-level gate rejects first, this is the
			// pool-side backstop.
			return ErrTooManyRequests
		}
		count := 0
		for _, sb := range p.sandboxes {
			if sb.functionArn == functionArn {
				count++
			}
		}
		if count >= int(*reserved) {
			return ErrTooManyRequests
		}
		return nil
	}
	if len(p.sandboxes) >= p.globalLimit {
		return ErrTooManyRequests
	}
	return nil
}

// startSandbox builds the container and its Runtime API server. The
// container is created without Entrypoint/Cmd unless the function's
// ImageConfig overrides them, so the image's own ENTRYPOINT+CMD become
// PID1 — the images-create contract's only requirement is that the image
// implements the Runtime API, which it discovers through
// AWS_LAMBDA_RUNTIME_API injected here at creation.
func (p *sandboxPool) startSandbox(ctx context.Context, sb *lambdaSandbox, spec sandboxSpec) error {
	api, err := newRuntimeAPIServer()
	if err != nil {
		return err
	}

	env := make(map[string]string, len(spec.env)+1)
	for k, v := range spec.env {
		env[k] = v
	}
	env["AWS_LAMBDA_RUNTIME_API"] = api.Addr()

	cfg := mobyclient.AdvancedContainerConfig{
		Name:       sb.containerName,
		Image:      spec.imageURI,
		PullImage:  true,
		Env:        env,
		Network:    "bridge",
		AutoRemove: false,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Memory:     spec.memoryBytes,
		// MemorySwap equal to Memory disables swap — the AWS execution
		// environment has none.
		MemorySwap: spec.memoryBytes,
		LogConfig: mobyclient.ContainerLogConfig{
			Driver: sandboxLogDriver,
			Options: map[string]string{
				"max-size": sandboxLogMaxSize,
				"max-file": sandboxLogMaxFiles,
			},
		},
	}
	if len(spec.entrypoint) > 0 {
		cfg.Entrypoint = spec.entrypoint
	}
	if len(spec.cmd) > 0 {
		cfg.Cmd = spec.cmd
	}
	if spec.workingDir != "" {
		cfg.WorkingDir = spec.workingDir
	}

	result, err := p.docker.CreateContainerFromConfig(ctx, cfg)
	if err != nil {
		api.Close()
		return fmt.Errorf("failed to create sandbox container: %w", err)
	}
	if err := p.docker.StartContainer(ctx, result.ID); err != nil {
		api.Close()
		if rmErr := p.docker.RemoveContainer(ctx, result.ID, true); rmErr != nil {
			logs.Warn("Failed to remove sandbox container after start failure",
				logs.String("containerID", result.ID), logs.Err(rmErr))
		}
		return fmt.Errorf("failed to start sandbox container: %w", err)
	}

	sb.containerID = result.ID
	sb.api = api
	return nil
}

// release returns a sandbox after its invoke settled. A sandbox marked
// doomed by a drain is destroyed here — its in-flight invoke was allowed
// to finish on the old environment, the AWS behaviour on code and
// configuration updates.
func (p *sandboxPool) release(sb *lambdaSandbox) {
	p.mu.Lock()
	if _, ok := p.sandboxes[sb.id]; !ok || sb.dead {
		p.mu.Unlock()
		return
	}
	if sb.doomed {
		p.mu.Unlock()
		p.destroy(sb)
		return
	}
	sb.busy = false
	sb.lastUsedAt = time.Now()
	p.mu.Unlock()
}

// destroy removes a sandbox from the pool and tears down its container and
// Runtime API server. Safe to call more than once.
func (p *sandboxPool) destroy(sb *lambdaSandbox) {
	p.mu.Lock()
	if sb.dead {
		p.mu.Unlock()
		return
	}
	sb.dead = true
	delete(p.sandboxes, sb.id)
	p.mu.Unlock()
	p.teardown(sb)
}

// teardown performs the docker-side cleanup without touching pool state.
func (p *sandboxPool) teardown(sb *lambdaSandbox) {
	if sb.api != nil {
		sb.api.Close()
	}
	if sb.containerID == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.docker.RemoveContainer(ctx, sb.containerID, true); err != nil {
		logs.Warn("Failed to remove sandbox container",
			logs.String("containerID", sb.containerID), logs.Err(err))
	}
}

// drainFunction destroys every sandbox of the function — all versions —
// for the delete cores. Busy sandboxes are marked doomed and destroyed by
// their owning invoke's release.
func (p *sandboxPool) drainFunction(functionArn string) {
	p.drainMatching(func(sb *lambdaSandbox) bool { return sb.functionArn == functionArn })
}

// drainVersion destroys a function version's sandboxes for code and
// configuration updates: a sandbox bakes its image, memory, timeout, and
// environment at creation, so a $LATEST update must invalidate the warm
// ones instead of letting them run stale.
func (p *sandboxPool) drainVersion(functionArn, version string) {
	p.drainMatching(func(sb *lambdaSandbox) bool {
		return sb.functionArn == functionArn && sb.version == version
	})
}

func (p *sandboxPool) drainMatching(match func(*lambdaSandbox) bool) {
	p.mu.Lock()
	var victims []*lambdaSandbox
	for _, sb := range p.sandboxes {
		if !match(sb) {
			continue
		}
		if sb.busy {
			sb.doomed = true
			continue
		}
		victims = append(victims, sb)
	}
	p.mu.Unlock()
	for _, sb := range victims {
		p.destroy(sb)
	}
}

// reapLoop reaps idle sandboxes past the TTL — the last one included; AWS
// terminates inactive environments rather than keeping one permanently
// warm — and probes idle-but-unparked sandboxes whose runtime never came
// back to /next, freeing their concurrency slot without waiting for the
// TTL.
func (p *sandboxPool) reapLoop() {
	defer close(p.done)
	ticker := time.NewTicker(p.reapInterval)
	defer ticker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
			p.reapOnce()
		}
	}
}

func (p *sandboxPool) reapOnce() {
	now := time.Now()

	p.mu.Lock()
	var ttlVictims, probeVictims []*lambdaSandbox
	for _, sb := range p.sandboxes {
		if sb.busy || sb.dead {
			continue
		}
		if now.Sub(sb.lastUsedAt) > p.idleTTL {
			ttlVictims = append(ttlVictims, sb)
			continue
		}
		if sb.api != nil && !sb.api.Idle() {
			probeVictims = append(probeVictims, sb)
		}
	}
	p.mu.Unlock()

	for _, sb := range probeVictims {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		status, err := p.docker.GetContainerStatus(ctx, sb.containerID)
		cancel()
		if err != nil || status != mobyclient.ContainerStatusRunning {
			p.destroy(sb)
		}
	}
	for _, sb := range ttlVictims {
		p.destroy(sb)
	}
}

// readLogWindow drains the sandbox's container log since the last read
// and advances the per-sandbox cursor. The first read starts from the
// container's beginning, so the init-phase output attaches to the first
// invoke's log window, the AWS behaviour.
func (p *sandboxPool) readLogWindow(ctx context.Context, sb *lambdaSandbox) (string, error) {
	p.mu.Lock()
	cursor := sb.logCursor
	p.mu.Unlock()

	text, next, err := p.docker.ReadLogWindow(ctx, sb.containerID, cursor)
	if err != nil {
		return "", err
	}

	p.mu.Lock()
	sb.logCursor = next
	p.mu.Unlock()
	return text, nil
}

// newSandboxID mints a random name suffix; sandbox names carry no
// function identity because one function version owns many of them and
// the lambda- prefix is what the orphan cleanups match on.
func newSandboxID() (string, error) {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("mint sandbox id: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// sandboxExitPollInterval is the cadence of the container-exit watch; a
// sandbox whose PID1 dies without an answer settles within one interval
// instead of waiting out its deadline.
const sandboxExitPollInterval = 250 * time.Millisecond

// runImageExecution is the image-package execution strategy: acquire a
// sandbox from the pool, push this invocation as one Runtime API round,
// and settle at the answer, the sandbox's death, or the deadline —
// whichever comes first. The shared invoke tail consumes the outcome.
func (s *LambdaService) runImageExecution(function *lambdastore.Function, ver *lambdastore.Version, execCfg executionConfig, region, version string, rec invocationRecord, eventJSON string) (executionOutcome, error) {
	ctx := context.Background()

	spec := s.sandboxSpecFor(function, ver, execCfg, region, version)
	sb, err := s.sandboxes.acquire(ctx, function.FunctionArn, version, function.ReservedConcurrency, spec)
	if err != nil {
		return executionOutcome{}, err
	}

	execStart := time.Now()
	round := s.driveSandboxRound(ctx, sb, rec, eventJSON)
	round.outcome.duration = time.Since(execStart)

	// The log window is read before any destruction: docker keeps a stopped
	// container's logs until removal, and the hang or crash output is part
	// of this invocation's log.
	if text, logErr := s.sandboxes.readLogWindow(ctx, sb); logErr != nil {
		logs.Warn("Failed to read the sandbox log window",
			logs.String("containerID", sb.containerID), logs.Err(logErr))
	} else {
		round.outcome.logOut = text
	}

	switch round.settled {
	case sandboxRoundAnswered:
		// The runtime answered; it loops back to /next on its own and the
		// released sandbox returns to the warm pool.
		s.sandboxes.release(sb)
	case sandboxRoundDead, sandboxRoundTimedOut:
		// A sandbox that died or blew its deadline is destroyed — AWS
		// recycles the whole execution environment on timeout — and the
		// next invoke builds a fresh one.
		s.sandboxes.destroy(sb)
	}
	return round.outcome, nil
}

// sandboxRound is one settled sandbox round: the strategy-neutral outcome
// plus how it settled, which only the sandbox strategy consumes.
type sandboxRound struct {
	outcome executionOutcome
	settled sandboxRoundResult
}

// sandboxRoundResult is how one sandbox round settled.
type sandboxRoundResult int

const (
	sandboxRoundAnswered sandboxRoundResult = iota
	sandboxRoundDead
	sandboxRoundTimedOut
)

// driveSandboxRound pushes the event as the sandbox's next Runtime API
// round and waits for the runtime's answer, the container's exit, or the
// invocation deadline.
func (s *LambdaService) driveSandboxRound(ctx context.Context, sb *lambdaSandbox, rec invocationRecord, eventJSON string) sandboxRound {
	sb.api.BeginRound(rec, []byte(eventJSON))

	exited, stopWatch := s.watchSandboxExit(ctx, sb)
	defer stopWatch()

	deadlineTimer := time.NewTimer(time.Until(rec.Deadline))
	defer deadlineTimer.Stop()

	select {
	case <-sb.api.Answered():
		body, kind, _ := sb.api.Captured()
		return sandboxRound{
			outcome: executionOutcome{
				payload:     string(body),
				apiAnswered: true,
				apiKind:     kind,
			},
			settled: sandboxRoundAnswered,
		}

	case <-exited:
		// PID1 ended without answering — an init failure or a crash. A
		// captured /init/error outranks the bare exit; otherwise the
		// invocation is Unhandled, the classification a non-zero exit
		// carries.
		if body, kind, captured := sb.api.Captured(); captured {
			return sandboxRound{
				outcome: executionOutcome{
					payload:     string(body),
					apiAnswered: true,
					apiKind:     kind,
				},
				settled: sandboxRoundDead,
			}
		}
		return sandboxRound{
			outcome: executionOutcome{exitCode: 1},
			settled: sandboxRoundDead,
		}

	case <-deadlineTimer.C:
		// The deadline expired with no answer; the kill-and-remove happens
		// in the pool's destroy, driven by the caller.
		return sandboxRound{
			outcome: executionOutcome{timedOut: true, exitCode: execExitTimedOut},
			settled: sandboxRoundTimedOut,
		}
	}
}

// watchSandboxExit reports the sandbox container's death. A cancelled
// watch (the round settled another way) leaves the channel unclosed — the
// goroutine returns without signalling.
func (s *LambdaService) watchSandboxExit(ctx context.Context, sb *lambdaSandbox) (<-chan struct{}, func()) {
	wctx, cancel := context.WithCancel(ctx)
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		ticker := time.NewTicker(sandboxExitPollInterval)
		defer ticker.Stop()
		for {
			status, err := s.dockerClient.GetContainerStatus(wctx, sb.containerID)
			if wctx.Err() != nil {
				return
			}
			if err != nil || status != mobyclient.ContainerStatusRunning {
				return
			}
			<-ticker.C
		}
	}()
	return exited, cancel
}

// sandboxSpecFor bakes the sandbox creation parameters from the effective
// function configuration: the image, the Lambda environment variables (the
// same set the zip model's containers carry), the memory ceiling with swap
// disabled, and the function's ImageConfig overrides — set fields replace
// the image's ENTRYPOINT/CMD/WORKDIR, unset fields keep them, the AWS
// container-image contract.
func (s *LambdaService) sandboxSpecFor(function *lambdastore.Function, ver *lambdastore.Version, execCfg executionConfig, region, version string) sandboxSpec {
	envVars := map[string]string{
		"AWS_LAMBDA_FUNCTION_TIMEOUT":     fmt.Sprintf("%d", execCfg.Timeout),
		"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": fmt.Sprintf("%d", execCfg.MemorySize),
		"AWS_LAMBDA_FUNCTION_HANDLER":     execCfg.Handler,
		"AWS_LAMBDA_FUNCTION_NAME":        function.FunctionName,
		"AWS_LAMBDA_FUNCTION_VERSION":     version,
		"AWS_REGION":                      region,
		// The documented execution-environment members an AWS runtime
		// interface client validates once its first /next answers: one log
		// stream per execution environment, the naming convention the
		// invocation records already use.
		"AWS_LAMBDA_LOG_GROUP_NAME":  lambdaLogGroupName(function.FunctionName),
		"AWS_LAMBDA_LOG_STREAM_NAME": lambdaLogStreamName(time.Now().UTC(), version, sbLogStreamSeed()),
	}
	if execCfg.Environment != nil {
		for k, v := range execCfg.Environment.Variables {
			envVars[k] = v
		}
	}
	if _, ok := envVars["AWS_ENDPOINT_URL"]; !ok && s.hostEndpoint != "" {
		envVars["AWS_ENDPOINT_URL"] = s.hostEndpoint
	}

	spec := sandboxSpec{
		imageURI:    execCfg.ImageUri,
		env:         envVars,
		memoryBytes: int64(execCfg.MemorySize) * 1024 * 1024,
	}

	imageCfg := function.ImageConfig
	if ver != nil && ver.ImageConfig != nil {
		imageCfg = ver.ImageConfig
	}
	if imageCfg != nil {
		spec.entrypoint = imageCfg.EntryPoint
		spec.cmd = imageCfg.Command
		spec.workingDir = imageCfg.WorkingDirectory
	}
	return spec
}

// sbLogStreamSeed mints the per-environment id of a sandbox's log stream,
// standing in for the request id an invocation record's stream uses.
func sbLogStreamSeed() string {
	return uuid.New().String()
}
