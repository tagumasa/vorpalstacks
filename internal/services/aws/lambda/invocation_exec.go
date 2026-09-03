package lambda

// This file drives one bootstrap-runtime execution. The invocation settles
// at the earliest of three events: the bootstrap process exiting (a
// single-shot bootstrap, an init error, a crash), the bootstrap's Runtime
// API answer arriving (the AWS contract answers the caller when the
// response POST lands — a looping runtime then waits on /next again, which
// is the next invocation's business, not this one's), and the invocation
// deadline (enforced here on the host side: the Docker API cannot kill an
// individual exec, so the container — the execution environment AWS
// recycles on timeout — is killed and removed instead).

import (
	"context"
	"time"

	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/core/logs"
)

// bootstrapLogGrace bounds how long a bootstrap execution may take to exit
// after its answer was captured. Closing the Runtime API makes a looping
// bootstrap's next /next fail, which terminates a well-behaved runtime
// within milliseconds; the grace keeps the invocation from waiting on one
// that lingers, at the cost of its trailing console output.
const bootstrapLogGrace = 2 * time.Second

// killDrainGrace bounds the wait for the exec observation to settle after
// the container kill: the kill closes the streams the observation blocks
// on, so it settles promptly. A bootstrap that ignores even that is
// abandoned with a warning — the container is gone either way.
const killDrainGrace = 5 * time.Second

// bootstrapExecOutcome is the settled result of a bootstrap execution.
type bootstrapExecOutcome struct {
	result   *mobyclient.ExecResult
	timedOut bool
	err      error
}

// runBootstrapExecution observes the exec asynchronously and settles the
// invocation without waiting for the bootstrap process to end. Known
// limitation: the deadline path kills the whole container, so a concurrent
// execution sharing the container is terminated with it — the platform runs
// one execution environment per function version, and AWS recycles the
// whole sandbox on timeout for the same reason.
func (s *LambdaService) runBootstrapExecution(ctx context.Context, containerID, containerName string, cfg mobyclient.ExecConfig, api *runtimeAPIServer, deadline time.Time) bootstrapExecOutcome {
	done := make(chan bootstrapExecOutcome, 1)
	go func() {
		result, err := s.dockerClient.Exec(ctx, containerID, cfg)
		done <- bootstrapExecOutcome{result: result, err: err}
	}()

	deadlineTimer := time.NewTimer(time.Until(deadline))
	defer deadlineTimer.Stop()

	select {
	case outcome := <-done:
		// The bootstrap process ended on its own: a single-shot runtime,
		// an init error exit, or a crash. Its captured answer (if any)
		// is merged by the caller from the Runtime API server.
		return outcome

	case <-api.Answered():
		// The answer is in — the invocation is settled. Close the Runtime
		// API so a looping bootstrap's next /next fails and it exits on
		// its own, then grant a short grace for its trailing output.
		api.Close()
		grace := time.NewTimer(bootstrapLogGrace)
		defer grace.Stop()
		select {
		case outcome := <-done:
			return outcome
		case <-grace.C:
			return bootstrapExecOutcome{result: &mobyclient.ExecResult{}}
		}

	case <-deadlineTimer.C:
		// The deadline expired with no answer. Killing the container's
		// init process closes the exec's streams and settles the
		// observation; removing the container keeps the next invocation
		// from reusing the timed-out environment. The synthesised exit
		// status carries the timeout so the caller's classification sees
		// an unsuccessful execution.
		if err := s.dockerClient.KillContainer(ctx, containerID, "SIGKILL"); err != nil {
			logs.Warn("Failed to kill the timed-out function container", logs.String("containerID", containerID), logs.Err(err))
		}
		if err := s.dockerClient.RemoveContainer(ctx, containerID, true); err != nil {
			logs.Warn("Failed to remove the timed-out function container", logs.String("containerID", containerID), logs.Err(err))
		}
		s.containerIDs.Delete(containerName)
		drain := time.NewTimer(killDrainGrace)
		defer drain.Stop()
		select {
		case outcome := <-done:
			outcome.timedOut = true
			if outcome.result != nil && outcome.result.ExitCode == 0 {
				// A stream-teardown race can report a zero exit status
				// for a killed process; the invocation still timed out.
				outcome.result.ExitCode = execExitTimedOut
			}
			return outcome
		case <-drain.C:
			logs.Warn("Timed-out bootstrap exec did not settle after the container kill", logs.String("containerID", containerID))
			return bootstrapExecOutcome{result: &mobyclient.ExecResult{ExitCode: execExitTimedOut}, timedOut: true}
		}
	}
}
