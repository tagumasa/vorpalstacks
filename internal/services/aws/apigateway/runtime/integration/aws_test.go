package integration

import "testing"

// TestParseLambdaResponseRejectsMalformedProxyPayload pins the
// malformed-proxy-response contract: a Lambda payload that parses as
// JSON but carries no valid status code must be rejected so the caller
// answers 502 instead of writing a zero status code.
func TestParseLambdaResponseRejectsMalformedProxyPayload(t *testing.T) {
	cases := map[string]string{
		"raw echo payload": `{"headers":{"x":"y"},"body":"unrelated"}`,
		"empty object":     `{}`,
		"out of range":     `{"statusCode":600}`,
	}
	for name, payload := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := parseLambdaResponse([]byte(payload)); err == nil {
				t.Fatalf("expected malformed proxy response %q to be rejected", payload)
			}
		})
	}

	resp, err := parseLambdaResponse([]byte(`{"statusCode":200,"body":"ok"}`))
	if err != nil {
		t.Fatalf("valid proxy response rejected: %v", err)
	}
	if resp.StatusCode != 200 || string(resp.Body) != "ok" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
