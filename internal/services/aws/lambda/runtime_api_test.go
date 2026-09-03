package lambda

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRuntimeAPINextServesEventAndHeaders(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST",
		"arn:aws:lambda:us-east-1:123456789012:function:fn", 256, 10,
		`{"custom":{"k":"v"}}`)
	api, err := startRuntimeAPI(rec, []byte(`{"hello":"world"}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", api.Port(), runtimeAPINextPath))
	if err != nil {
		t.Fatalf("GET next: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on next, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"hello":"world"}` {
		t.Fatalf("next body = %q, want the event", body)
	}
	if got := resp.Header.Get("Lambda-Runtime-Aws-Request-Id"); got != rec.RequestID {
		t.Fatalf("request id header = %q, want %q", got, rec.RequestID)
	}
	wantDeadline := strconv.FormatInt(rec.DeadlineUnixMS(), 10)
	if got := resp.Header.Get("Lambda-Runtime-Deadline-Ms"); got != wantDeadline {
		t.Fatalf("deadline header = %q, want %q", got, wantDeadline)
	}
	if got := resp.Header.Get("Lambda-Runtime-Invoked-Function-Arn"); got != rec.InvokedARN {
		t.Fatalf("invoked arn header = %q, want %q", got, rec.InvokedARN)
	}
	if got := resp.Header.Get("Lambda-Runtime-Invocation-Id"); got != rec.RequestID {
		t.Fatalf("invocation id header = %q, want %q", got, rec.RequestID)
	}
	wantCC := base64.StdEncoding.EncodeToString([]byte(rec.ClientContextJSON))
	if got := resp.Header.Get("Lambda-Runtime-Client-Context"); got != wantCC {
		t.Fatalf("client context header = %q, want %q", got, wantCC)
	}
}

func TestRuntimeAPINextOmitsClientContextWhenAbsent(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", api.Port(), runtimeAPINextPath))
	if err != nil {
		t.Fatalf("GET next: %v", err)
	}
	defer resp.Body.Close()
	if got := resp.Header.Get("Lambda-Runtime-Client-Context"); got != "" {
		t.Fatalf("client context header = %q, want absent", got)
	}
}

func TestRuntimeAPIResponseCapture(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d%s%s/response", api.Port(), runtimeAPIInvokePath, rec.RequestID)
	resp, err := http.Post(url, "application/json", strings.NewReader(`{"ok":1}`))
	if err != nil {
		t.Fatalf("POST response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 on response, got %d", resp.StatusCode)
	}
	body, kind, ok := api.Captured()
	if !ok || kind != runtimeAPIResponseKind || string(body) != `{"ok":1}` {
		t.Fatalf("captured = (%q, %q, %v), want response with body", body, kind, ok)
	}
}

func TestRuntimeAPIErrorCapture(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d%s%s/error", api.Port(), runtimeAPIInvokePath, rec.RequestID)
	resp, err := http.Post(url, "application/json",
		strings.NewReader(`{"errorMessage":"boom","errorType":"Runtime.Boom"}`))
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 on error, got %d", resp.StatusCode)
	}
	body, kind, ok := api.Captured()
	if !ok || kind != runtimeAPIErrorKind || !strings.Contains(string(body), "boom") {
		t.Fatalf("captured = (%q, %q, %v), want error kind with body", body, kind, ok)
	}
}

func TestRuntimeAPIInitErrorCapture(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d%s", api.Port(), runtimeAPIInitErrPath),
		"application/json", strings.NewReader(`{"errorMessage":"no handler"}`))
	if err != nil {
		t.Fatalf("POST init error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 on init error, got %d", resp.StatusCode)
	}
	_, kind, ok := api.Captured()
	if !ok || kind != runtimeAPIInitErrKind {
		t.Fatalf("captured kind = %q ok=%v, want init-error", kind, ok)
	}
}

func TestRuntimeAPIFirstAnswerWins(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	base := fmt.Sprintf("http://127.0.0.1:%d%s%s", api.Port(), runtimeAPIInvokePath, rec.RequestID)
	for _, suffix := range []string{"/response", "/error"} {
		resp, err := http.Post(base+suffix, "application/json", strings.NewReader("x"))
		if err != nil {
			t.Fatalf("POST %s: %v", suffix, err)
		}
		resp.Body.Close()
	}
	body, kind, ok := api.Captured()
	if !ok || kind != runtimeAPIResponseKind || string(body) != "x" {
		t.Fatalf("captured = (%q, %q, %v), want the first answer to win", body, kind, ok)
	}
}

func TestRuntimeAPIResponseTooLarge(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	oversized := strings.Repeat("a", syncMaxPayloadSize+1)
	url := fmt.Sprintf("http://127.0.0.1:%d%s%s/response", api.Port(), runtimeAPIInvokePath, rec.RequestID)
	resp, err := http.Post(url, "application/json", strings.NewReader(oversized))
	if err != nil {
		t.Fatalf("POST oversized response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("an oversized response must be rejected with 413, got %d", resp.StatusCode)
	}
	body, kind, ok := api.Captured()
	if !ok || kind != runtimeAPIErrorKind {
		t.Fatalf("captured = (%q, %q, %v), want an error-kind answer for the oversized response", body, kind, ok)
	}
	if !strings.Contains(string(body), "exceeded") {
		t.Fatalf("the oversize answer must explain the rejection, got %s", body)
	}

	withinLimit := strings.Repeat("a", 1024)
	resp2, err := http.Post(url, "application/json", strings.NewReader(withinLimit))
	if err != nil {
		t.Fatalf("POST in-limit response: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusAccepted {
		t.Fatalf("an in-limit response must be accepted, got %d", resp2.StatusCode)
	}
	body2, kind2, ok2 := api.Captured()
	if !ok2 || kind2 != runtimeAPIErrorKind || len(body2) == len(withinLimit) {
		t.Fatalf("the first (oversize) answer must stand, got (%q, %q, %v)", body2, kind2, ok2)
	}
}

func TestRuntimeAPIAnsweredChannel(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		payload string
	}{
		{name: "response", path: runtimeAPIInvokePath + "rid/response", payload: `{"ok":1}`},
		{name: "error", path: runtimeAPIInvokePath + "rid/error", payload: `{"errorMessage":"boom"}`},
		{name: "init error", path: runtimeAPIInitErrPath, payload: `{"errorMessage":"no handler"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
			api, err := startRuntimeAPI(rec, []byte(`{}`))
			if err != nil {
				t.Fatalf("startRuntimeAPI: %v", err)
			}
			defer api.Close()

			select {
			case <-api.Answered():
				t.Fatalf("Answered must stay open before any POST lands")
			default:
			}

			resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d%s", api.Port(), tc.path),
				"application/json", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("POST %s: %v", tc.path, err)
			}
			resp.Body.Close()

			select {
			case <-api.Answered():
			case <-time.After(2 * time.Second):
				t.Fatalf("Answered did not close after the %s POST", tc.name)
			}
		})
	}
}

func TestRuntimeAPISecondNextWaitsUntilClose(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}

	first := httptest.NewRecorder()
	api.ServeHTTP(first, httptest.NewRequest(http.MethodGet, runtimeAPINextPath, nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first next = %d, want 200", first.Code)
	}

	done := make(chan int, 1)
	go func() {
		second := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, runtimeAPINextPath, nil)
		api.ServeHTTP(second, req)
		done <- second.Code
	}()
	select {
	case code := <-done:
		t.Fatalf("second next returned %d before Close, want it to wait", code)
	case <-time.After(100 * time.Millisecond):
	}
	api.Close()
	select {
	case code := <-done:
		if code != http.StatusServiceUnavailable {
			t.Fatalf("second next after Close = %d, want 503", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("second next still blocked after Close")
	}
}

func TestRuntimeAPINotFound(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 3, "")
	api, err := startRuntimeAPI(rec, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet,
		"/2018-06-01/runtime/restore/next", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("unknown path = %d, want 404", recorder.Code)
	}
}

func TestRuntimeAPIAddrCarriesHostGatewayName(t *testing.T) {
	api, err := startRuntimeAPI(invocationRecord{RequestID: "r"}, []byte(`{}`))
	if err != nil {
		t.Fatalf("startRuntimeAPI: %v", err)
	}
	defer api.Close()

	want := fmt.Sprintf("host.docker.internal:%d", api.Port())
	if api.Addr() != want {
		t.Fatalf("Addr() = %q, want %q", api.Addr(), want)
	}
}
