package lambda

// This file builds the per-invocation command that executes the handler
// inside the function container, and owns the result framing contract
// between the wrapper and splitResultPayload. The node/python wrappers
// construct the AWS context object — the handler's second argument — from
// the invocation record, including the getRemainingTimeInMillis closure
// backed by the record's deadline.

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	lambdastore "vorpalstacks/internal/store/aws/lambda"

	"github.com/google/uuid"
)

// nodeWrapperTemplate is executed with `node -e`. The {PLACEHOLDER} tokens
// are substituted by buildInvokeCommand; the event and client context are
// embedded base64-encoded so that payloads containing newlines, quotes or
// backslashes cannot break out of the script's string literals. The done
// flag makes the promise resolution and the callback idempotent: whichever
// finishes first with a value frames it exactly once. A callback-style
// handler (arity ≥ 3) resolves its promise with undefined immediately, so
// an undefined result alone must not consume the invocation — the wrapper
// keeps waiting for the callback the handler registered. A promise or
// synchronous return is framed the moment it settles and the wrapper
// exits at once: the platform freezes the execution at the response, so
// pending timers must not hold the invocation open until the timeout. A
// callback response follows the AWS callback contract instead — the
// response is not sent until all event loop tasks have finished — so cb
// only records the result and a beforeExit handler frames it once the
// loop is empty. A callback handler whose loop never drains (a lingering
// interval) reaches the timeout and the caller receives the timeout
// error instead of the recorded response, which is the documented AWS
// outcome; only an explicit callbackWaitsForEmptyEventLoop=false would
// opt out, which this wrapper does not expose.
// A failing handler is framed the same way: the E helper writes the
// Runtime API error document — errorMessage, errorType and the stack
// trace split into lines — inside the markers and exits non-zero, so the
// payload carries the error object while the exit status carries the
// classification: exit 3 for a failure the handler signalled through the
// callback (the AWS Node.js runtime reports callback-signalled errors as
// Handled function errors) and exit 1 for an uncaught failure
// (Unhandled). The stack trace also lands on stderr, which is the
// channel the execution logs read.
const nodeWrapperTemplate = `const m=require('/var/task/{MOD}');const h=typeof m==='function'?m:m['{FN}'];const ev=JSON.parse(Buffer.from('{EVB64}','base64').toString('utf8'));const cc={CC};const ctx={functionName:'{FNAME}',functionVersion:'{VER}',invokedFunctionArn:'{ARN}',memoryLimitInMB:{MEM},awsRequestId:'{RID}',logGroupName:'{LG}',logStreamName:'{LS}',clientContext:cc,identity:null,callbackWaitsForEmptyEventLoop:true,getRemainingTimeInMillis:()=>{DEADLINE}-Date.now()};let done=false,cbFired=false,cbRes;const usesCb=h.length>=3;const W=r=>{if(done)return;if(usesCb&&cbFired)return;if(r===undefined&&usesCb)return;done=true;if(r!==undefined)process.stdout.write('{MK}'+JSON.stringify(r)+'{MK}');process.exit(0);};const E=(x,h)=>{if(done)return;done=true;process.stderr.write(String(x&&x.stack||x));process.stdout.write('{MK}'+JSON.stringify({errorMessage:x&&x.message||String(x),errorType:x&&x.name||'Error',stackTrace:String(x&&x.stack||'').split('\n')})+'{MK}');process.exit(h?{EXIT_HANDLED}:1);};process.on('beforeExit',()=>{if(!usesCb||done||!cbFired)return;done=true;if(cbRes!==undefined)process.stdout.write('{MK}'+JSON.stringify(cbRes)+'{MK}');});const cb=(e,r)=>{if(cbFired)return;cbFired=true;if(e){E(e,true);return;}cbRes=r;};Promise.resolve(h(ev,ctx,cb)).then(W).catch(e=>{E(e);});`

// pythonWrapperTemplate is executed with `python3 -c`. It mirrors the node
// wrapper: base64-embedded event and client context, a SimpleNamespace
// context object with snake_case members, and the marker-framed return
// value. It exits once the response is printed and flushed, matching the
// node wrapper's prompt exit. A failing handler — or a return value that
// cannot be serialised — is framed as the Runtime API error document
// (errorMessage, errorType, stackTrace) with a non-zero exit, so the
// payload carries the error object and the classification stays
// Unhandled; the traceback is echoed to stderr for the execution logs.
const pythonWrapperTemplate = `import json,sys,time,types,base64,os,traceback
mod=__import__('{MOD}');h=getattr(mod,'{FN}',mod)
ev=json.loads(base64.b64decode('{EVB64}').decode('utf-8'))
cc={CC}
ctx=types.SimpleNamespace(function_name='{FNAME}',function_version='{VER}',invoked_function_arn='{ARN}',memory_limit_in_mb={MEM},aws_request_id='{RID}',log_group_name='{LG}',log_stream_name='{LS}',client_context=cc,identity=None,get_remaining_time_in_millis=lambda:{DEADLINE}-int(time.time()*1000))
try:
    r=h(ev,ctx)
    print('{MK}'+json.dumps(r)+'{MK}')
except BaseException as ex:
    print('{MK}'+json.dumps({"errorMessage":str(ex),"errorType":type(ex).__name__,"stackTrace":traceback.format_exc().splitlines()})+'{MK}')
    sys.stdout.flush()
    sys.stderr.write(traceback.format_exc())
    sys.exit(1)
sys.stdout.flush()
os._exit(0)`

// usesRuntimeWrapper reports whether the runtime runs the handler through
// the generated wrapper script (nodejs*, python*). Every other runtime —
// the provided.* custom runtimes and the RIC-based managed images — enters
// through /var/runtime/bootstrap and exchanges the event and the response
// over the Runtime API.
func usesRuntimeWrapper(runtime lambdastore.Runtime) bool {
	return strings.HasPrefix(string(runtime), "nodejs") || strings.HasPrefix(string(runtime), "python")
}

// escapeForSingleQuotes makes a value safe to embed inside a single-quoted
// literal in either wrapper script.
func escapeForSingleQuotes(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", `\'`)
}

// execExitTimedOut is the exit status GNU coreutils timeout(1) reports when
// it killed the command on expiry.
const execExitTimedOut = 124

// execExitHandledError is the exit status the node wrapper reports when the
// handler signalled its failure through the callback. The AWS Node.js
// runtime reports callback-signalled errors as Handled function errors, so
// the wrapper distinguishes them from uncaught failures (exit 1,
// Unhandled) by exit status — the only channel the host observes directly.
const execExitHandledError = 3

func (s *LambdaService) buildInvokeCommand(runtime lambdastore.Runtime, moduleFile, handlerFunc, eventJSON string, rec invocationRecord) ([]string, string) {
	// The node/python wrappers frame the serialised return value with a
	// per-invocation marker. Handler code cannot emit the marker (a fresh
	// one is generated per invocation), so the framed region is exactly the
	// return value and the preceding output is console logging — mirroring
	// the AWS runtime contract where the response payload and the logs are
	// separate channels. A custom runtime answers over the Runtime API
	// instead of stdout, so it gets no marker.
	marker := ""
	if usesRuntimeWrapper(runtime) {
		marker = fmt.Sprintf("__VORPALSTACKS_RESULT_%s__", uuid.New().String())
	}
	var cmd []string
	if strings.HasPrefix(string(runtime), "nodejs") {
		script := wrapperReplacer(moduleFile, handlerFunc, eventJSON, rec, marker, "node").Replace(nodeWrapperTemplate)
		cmd = []string{"node", "-e", script}
	} else if strings.HasPrefix(string(runtime), "python") {
		script := wrapperReplacer(moduleFile, handlerFunc, eventJSON, rec, marker, "python").Replace(pythonWrapperTemplate)
		cmd = []string{"python3", "-c", script}
	} else {
		cmd = []string{"/var/runtime/bootstrap"}
	}
	// GNU timeout(1) enforces the function timeout inside the container:
	// on expiry it kills the handler with exit status execExitTimedOut and
	// invokeFunction answers with the timeout envelope. The -k escalation
	// follows the initial SIGTERM with a SIGKILL two seconds later, so a
	// handler that traps and ignores the signal cannot outlive its
	// deadline. The record carries the already-normalised timeout. Only
	// the wrapper runtimes get the prefix: they always execute on the
	// AWS-managed zip images, which carry coreutils. The bootstrap
	// runtimes must not be handed image requirements beyond the runtime
	// itself, so their deadline is enforced host-side instead
	// (see invocation_exec.go).
	if usesRuntimeWrapper(runtime) {
		return append([]string{"timeout", "-k", "2", strconv.FormatInt(int64(rec.TimeoutSeconds), 10)}, cmd...), marker
	}
	return cmd, marker
}

// wrapperReplacer builds the token substitution shared by both wrapper
// templates. Every value landing inside a single-quoted script literal is
// escaped — the handler tokens included, because the handler validation
// constrains the format but not the character set, and an unescaped quote
// in a module or function name would break out of the literal and inject
// script.
func wrapperReplacer(moduleFile, handlerFunc, eventJSON string, rec invocationRecord, marker, language string) *strings.Replacer {
	return strings.NewReplacer(
		"{MOD}", escapeForSingleQuotes(moduleFile),
		"{FN}", escapeForSingleQuotes(handlerFunc),
		"{EVB64}", base64.StdEncoding.EncodeToString([]byte(eventJSON)),
		"{CC}", clientContextExpr(rec.ClientContextJSON, language),
		"{FNAME}", escapeForSingleQuotes(rec.FunctionName),
		"{VER}", escapeForSingleQuotes(rec.Version),
		"{ARN}", escapeForSingleQuotes(rec.InvokedARN),
		"{MEM}", strconv.FormatInt(int64(rec.MemorySize), 10),
		"{RID}", escapeForSingleQuotes(rec.RequestID),
		"{LG}", escapeForSingleQuotes(rec.LogGroupName),
		"{LS}", escapeForSingleQuotes(rec.LogStreamName),
		"{DEADLINE}", strconv.FormatInt(rec.DeadlineUnixMS(), 10),
		"{MK}", marker,
		"{EXIT_HANDLED}", strconv.Itoa(execExitHandledError),
	)
}

// clientContextExpr renders the decoded ClientContext document as the
// literal the wrapper assigns to the context's clientContext member: a
// parsed value when present, the language's null literal otherwise.
func clientContextExpr(clientContextJSON, language string) string {
	if clientContextJSON == "" {
		if language == "python" {
			return "None"
		}
		return "null"
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(clientContextJSON))
	if language == "python" {
		return "json.loads(base64.b64decode('" + encoded + "').decode('utf-8'))"
	}
	return "JSON.parse(Buffer.from('" + encoded + "','base64').toString('utf8'))"
}
