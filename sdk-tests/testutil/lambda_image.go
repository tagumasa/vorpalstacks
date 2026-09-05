package testutil

// Container-image function tests: the image's own ENTRYPOINT runs as PID1
// and speaks the Runtime API against the platform's host-side server. Two
// fixtures cover both halves of the contract — a locally built Alpine image
// with a shell RIC behind its own ENTRYPOINT (the arbitrary-entrypoint
// case) and the AWS nodejs base image (a real runtime interface client).

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// shellRICDockerfile builds the arbitrary-entrypoint fixture: a shell
// runtime interface client as the image's ENTRYPOINT, the shape the
// platform must run without rewriting.
const shellRICDockerfile = `FROM alpine:3.20
RUN apk add --no-cache curl
COPY ric.sh /app/ric
RUN chmod +x /app/ric
ENTRYPOINT ["/app/ric"]
`

// shellRIC loops the Runtime API the documented way: park on
// /runtime/invocation/next, echo one console line, POST the event back
// unchanged. Every non-hang round sleeps long enough that two serialized
// invokes cannot hide behind warm reuse; an event carrying "hang" parks
// the round until the platform's deadline reaps the sandbox.
const shellRIC = `#!/bin/sh
set -u
API="${AWS_LAMBDA_RUNTIME_API:?AWS_LAMBDA_RUNTIME_API not set}"
BASE="http://$API/2018-06-01/runtime"
while true; do
  rm -f /tmp/h /tmp/b
  if ! curl -sS -D /tmp/h -o /tmp/b "$BASE/invocation/next"; then
    exit 0
  fi
  RID=$(sed -n 's/^Lambda-Runtime-Aws-Request-Id:[[:space:]]*//p' /tmp/h | tr -d '\r' | head -n 1)
  echo "ric round $RID body $(cat /tmp/b)"
  case "$(cat /tmp/b)" in
    *hang*)
      echo "ric hanging until the deadline"
      sleep 600
      ;;
    *)
      sleep 2
      ;;
  esac
  curl -sS -X POST --data-binary @/tmp/b "$BASE/invocation/$RID/response" >/dev/null || exit 0
done
`

// awsBaseDockerfile builds the real-runtime fixture from the AWS nodejs
// base image: its own ENTRYPOINT (the AWS runtime interface client) stays
// untouched and the handler ships inside the image, the canonical
// container-image function shape.
const awsBaseDockerfile = `FROM public.ecr.aws/lambda/nodejs:24
COPY app.js ${LAMBDA_TASK_ROOT}/app.js
CMD ["app.handler"]
`

const awsBaseHandler = `exports.handler = async (event) => {
    console.log("aws-base image handler received:", JSON.stringify(event));
    return event;
};
`

// buildLambdaImageFixture builds one docker image from inline files and
// returns its tag. The build context is a throwaway directory.
func buildLambdaImageFixture(tag, dockerfile string, files map[string]string) error {
	dir, err := os.MkdirTemp("", "vs-lambda-image-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(dockerfile), 0o644); err != nil {
		return err
	}
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker CLI not available: %w", err)
	}
	cmd := exec.Command(dockerPath, "build", "-t", tag, dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker build %s: %w\n%s", tag, err, out)
	}
	return nil
}

// removeLambdaImageFixture drops a built fixture image, best effort.
func removeLambdaImageFixture(tag string) {
	if dockerPath, err := exec.LookPath("docker"); err == nil {
		exec.Command(dockerPath, "rmi", "-f", tag).Run()
	}
}

// createImageFunction creates one PackageType=Image function on the given
// locally built fixture tag and defers its cleanup (function + log group).
func (tc *lambdaTestContext) createImageFunction(name, roleARN, imageTag string, opts ...func(*lambda.CreateFunctionInput)) func() {
	input := &lambda.CreateFunctionInput{
		FunctionName: aws.String(name),
		Role:         aws.String(roleARN),
		PackageType:  types.PackageTypeImage,
		Code:         &types.FunctionCode{ImageUri: aws.String(imageTag)},
		Timeout:      aws.Int32(30),
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := tc.client.CreateFunction(tc.ctx, input); err != nil {
		// The invoke tests surface this as their failure; the cleanup stays
		// a no-op so a failed create cannot delete an unrelated function.
		return func() {}
	}
	return func() { tc.deleteFunctionAndLogs(name) }
}

// decodedLogResult base64-decodes an Invoke LogType=Tail result.
func decodedLogResult(resp *lambda.InvokeOutput) (string, error) {
	if resp.LogResult == nil {
		return "", fmt.Errorf("LogResult is nil")
	}
	raw, err := base64.StdEncoding.DecodeString(*resp.LogResult)
	if err != nil {
		return "", fmt.Errorf("decode LogResult: %w", err)
	}
	return string(raw), nil
}

// countFunctionSandboxes counts the platform's running sandboxes whose
// environment names the given function. The pool stamps every container
// with AWS_LAMBDA_FUNCTION_NAME, which is what makes a function's
// sandboxes distinguishable in docker — sandbox names carry no function
// identity. Containers that exit between the list and the inspect are
// simply not counted.
func countFunctionSandboxes(fn string) (int, error) {
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return 0, err
	}
	ids, err := exec.Command(dockerPath, "ps", "--filter", "name=lambda-",
		"--format", "{{.ID}}").Output()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, id := range strings.Fields(string(ids)) {
		env, err := exec.Command(dockerPath, "inspect", "-f",
			"{{range .Config.Env}}{{println .}}{{end}}", id).Output()
		if err != nil {
			continue
		}
		for _, entry := range strings.Split(string(env), "\n") {
			if strings.TrimSpace(entry) == "AWS_LAMBDA_FUNCTION_NAME="+fn {
				count++
				break
			}
		}
	}
	return count, nil
}

func runLambdaImageTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	roleARN, cleanupRole, err := tc.createRole(tc.unique("ImageExecRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Image_InvokeRoundtrip", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupRole()

	shellTag := "vs-lst-image-shell-" + tc.ts
	awsTag := "vs-lst-image-aws-" + tc.ts
	fixtureErr := buildLambdaImageFixture(shellTag, shellRICDockerfile, map[string]string{"ric.sh": shellRIC})
	if fixtureErr == nil {
		fixtureErr = buildLambdaImageFixture(awsTag, awsBaseDockerfile, map[string]string{"app.js": awsBaseHandler})
	}
	defer removeLambdaImageFixture(shellTag)
	defer removeLambdaImageFixture(awsTag)

	// A real runtime interface client (the AWS base image's own ENTRYPOINT)
	// must round-trip an invoke: the payload comes back, no function error,
	// and the tailed log carries the handler's console output captured from
	// the container's log stream. Two consecutive invokes pin the log-window
	// boundary: the pool serves the second round on the released warm
	// sandbox, whose window cursor sits on the first round's final console
	// line, so the second tail must carry its own line without re-delivering
	// the first round's.
	results = append(results, tc.r.RunTest("lambda", "Image_InvokeRoundtrip", func() error {
		if fixtureErr != nil {
			return fmt.Errorf("image fixture build failed: %v", fixtureErr)
		}
		fn := tc.unique("ImageFnRoundtrip")
		cleanup := tc.createImageFunction(fn, roleARN, awsTag)
		defer cleanup()

		invokeWithTail := func(payload []byte) (string, error) {
			resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
				FunctionName: aws.String(fn),
				Payload:      payload,
				LogType:      types.LogTypeTail,
			})
			if err != nil {
				return "", err
			}
			if resp.StatusCode != 200 {
				return "", fmt.Errorf("expected status 200, got %d (payload %s)", resp.StatusCode, string(resp.Payload))
			}
			if aws.ToString(resp.FunctionError) != "" {
				return "", fmt.Errorf("unexpected function error %q (payload %s)", aws.ToString(resp.FunctionError), string(resp.Payload))
			}
			var got, want map[string]interface{}
			if err := json.Unmarshal(resp.Payload, &got); err != nil {
				return "", fmt.Errorf("parse echoed payload: %w (%s)", err, string(resp.Payload))
			}
			if err := json.Unmarshal(payload, &want); err != nil {
				return "", fmt.Errorf("parse request payload: %w", err)
			}
			if fmt.Sprint(got) != fmt.Sprint(want) {
				return "", fmt.Errorf("payload not echoed back: got %s, want %s", string(resp.Payload), string(payload))
			}
			return decodedLogResult(resp)
		}

		firstTail, err := invokeWithTail([]byte(`{"echo":"image-roundtrip-a1"}`))
		if err != nil {
			return err
		}
		if !strings.Contains(firstTail, "aws-base image handler received") ||
			!strings.Contains(firstTail, "image-roundtrip-a1") {
			return fmt.Errorf("tailed log misses the handler console output: %s", firstTail)
		}

		secondTail, err := invokeWithTail([]byte(`{"echo":"image-roundtrip-b2"}`))
		if err != nil {
			return err
		}
		if !strings.Contains(secondTail, "image-roundtrip-b2") {
			return fmt.Errorf("second tailed log misses this round's console output: %s", secondTail)
		}
		if strings.Contains(secondTail, "image-roundtrip-a1") {
			return fmt.Errorf("second tailed log re-delivers the first round's boundary line: %s", secondTail)
		}
		return nil
	}))

	// Two concurrent invokes of the shell-RIC image must not serialize.
	// Each round parks 2s in the fixture RIC, and while both are in
	// flight the test observes docker directly: a serialized pair queues
	// both rounds on the single sandbox, so only one container ever
	// exists, while a parallel pair must grow a second sandbox the moment
	// its acquire finds the first one busy. Two distinct sandboxes of the
	// function parked simultaneously is the direct observation of
	// non-serialization, with no elapsed-time bound and no exposure to
	// cold-start jitter.
	results = append(results, tc.r.RunTest("lambda", "Image_ParallelInvokes", func() error {
		if fixtureErr != nil {
			return fmt.Errorf("image fixture build failed: %v", fixtureErr)
		}
		fn := tc.unique("ImageFnParallel")
		cleanup := tc.createImageFunction(fn, roleARN, shellTag, func(in *lambda.CreateFunctionInput) {
			in.Timeout = aws.Int32(30)
		})
		defer cleanup()

		const workers = 2
		errs := make([]error, workers)
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
					FunctionName: aws.String(fn),
					Payload:      []byte(fmt.Sprintf(`{"worker":%d}`, i)),
				})
				if err != nil {
					errs[i] = err
					return
				}
				if resp.StatusCode != 200 || aws.ToString(resp.FunctionError) != "" {
					errs[i] = fmt.Errorf("worker %d: status %d error %q payload %s", i, resp.StatusCode, aws.ToString(resp.FunctionError), string(resp.Payload))
					return
				}
				if string(resp.Payload) != fmt.Sprintf(`{"worker":%d}`, i) {
					errs[i] = fmt.Errorf("worker %d: payload not echoed back: %s", i, string(resp.Payload))
				}
			}(i)
		}

		parallel := false
		var dockerErr error
		for i := 0; i < 40; i++ {
			n, err := countFunctionSandboxes(fn)
			if err != nil {
				dockerErr = err
				break
			}
			if n >= 2 {
				parallel = true
				break
			}
			time.Sleep(250 * time.Millisecond)
		}
		wg.Wait()
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
		if dockerErr != nil {
			return fmt.Errorf("counting the function's sandboxes: %v", dockerErr)
		}
		if !parallel {
			return fmt.Errorf("no second sandbox of %s was ever created while both invokes were in flight — the invokes were serialized on one sandbox", fn)
		}
		return nil
	}))

	// A hanging round must hit the invocation deadline: the 200 response
	// carries the timeout envelope with an Unhandled function error, the
	// tailed log shows the hang, and the next invoke recovers on a fresh
	// sandbox. The 4s timeout leaves the fixture's 2s work round enough
	// headroom to answer on the warm-up and the recovery.
	results = append(results, tc.r.RunTest("lambda", "Image_TimeoutAndRecovery", func() error {
		if fixtureErr != nil {
			return fmt.Errorf("image fixture build failed: %v", fixtureErr)
		}
		fn := tc.unique("ImageFnTimeout")
		cleanup := tc.createImageFunction(fn, roleARN, shellTag, func(in *lambda.CreateFunctionInput) {
			in.Timeout = aws.Int32(4)
		})
		defer cleanup()

		// Warm the sandbox first so the hang round measures the deadline,
		// not the cold start. The warm-up itself must succeed — a silent
		// timeout here would invalidate everything after it.
		warm, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(fn),
			Payload:      []byte(`{"mode":"warm"}`),
		})
		if err != nil {
			return fmt.Errorf("warm-up invoke: %v", err)
		}
		if aws.ToString(warm.FunctionError) != "" {
			return fmt.Errorf("warm-up invoke did not succeed: error %q payload %s", aws.ToString(warm.FunctionError), string(warm.Payload))
		}
		if string(warm.Payload) != `{"mode":"warm"}` {
			return fmt.Errorf("warm-up payload not echoed back: %s", string(warm.Payload))
		}

		hangStart := time.Now()
		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(fn),
			Payload:      []byte(`{"mode":"hang"}`),
			LogType:      types.LogTypeTail,
		})
		hangElapsed := time.Since(hangStart)
		if err != nil {
			return fmt.Errorf("hang invoke: %v", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200 on a timed-out invoke, got %d", resp.StatusCode)
		}
		if aws.ToString(resp.FunctionError) != "Unhandled" {
			return fmt.Errorf("expected Unhandled on timeout, got %q", aws.ToString(resp.FunctionError))
		}
		if !strings.Contains(string(resp.Payload), "Task timed out after 4.00 seconds") {
			return fmt.Errorf("payload misses the timeout envelope: %s", string(resp.Payload))
		}
		if hangElapsed < 3900*time.Millisecond {
			return fmt.Errorf("the hang invoke returned after %v, before the 4s deadline", hangElapsed)
		}
		tail, err := decodedLogResult(resp)
		if err != nil {
			return err
		}
		if !strings.Contains(tail, "ric hanging until the deadline") {
			return fmt.Errorf("tailed log misses the hang output: %s", tail)
		}

		// The timed-out sandbox is destroyed; the next invoke builds a
		// fresh one and succeeds.
		recover, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
			FunctionName: aws.String(fn),
			Payload:      []byte(`{"mode":"after"}`),
		})
		if err != nil {
			return fmt.Errorf("recovery invoke: %v", err)
		}
		if recover.StatusCode != 200 || aws.ToString(recover.FunctionError) != "" {
			return fmt.Errorf("recovery invoke failed: status %d error %q payload %s", recover.StatusCode, aws.ToString(recover.FunctionError), string(recover.Payload))
		}
		if string(recover.Payload) != `{"mode":"after"}` {
			return fmt.Errorf("recovery payload not echoed back: %s", string(recover.Payload))
		}
		return nil
	}))

	return results
}
