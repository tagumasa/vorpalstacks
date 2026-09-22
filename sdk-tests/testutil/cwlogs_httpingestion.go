package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
)

// The HTTP ingestion endpoint pins (E6): ACWL bearer tokens issued as IAM
// service-specific credentials of logs.amazonaws.com, the four documented
// request shapes (OTLP JSON, HLC, ND-JSON, Structured JSON) through the
// plain-HTTP endpoints, and the authentication and validation rows.

// ingestionPost sends one plain-HTTP ingestion request; a nil bearerToken
// sends no Authorization header (the platform runs with signature
// verification disabled, the SigV4 leg of the documented auth pair).
func (tc *cwlogsTestCtx) ingestionPost(path, contentType, body, bearerToken string) (int, string, error) {
	url := tc.runner.endpoint + path
	var req *http.Request
	var err error
	if body == "" {
		req, err = http.NewRequest(http.MethodPost, url, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	}
	if err != nil {
		return 0, "", err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

// ingestionReadEvents reads the whole stream back through GetLogEvents
// (the SDK pin's own read loop; tokens page per the re-offer contract).
func (tc *cwlogsTestCtx) ingestionReadEvents(group, stream string) ([]string, error) {
	var messages []string
	var token *string
	for page := 0; page < 20; page++ {
		resp, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(group),
			LogStreamName: aws.String(stream),
			StartFromHead: aws.Bool(true),
			NextToken:     token,
		})
		if err != nil {
			return nil, fmt.Errorf("get page %d: %v", page, err)
		}
		for _, e := range resp.Events {
			messages = append(messages, aws.ToString(e.Message))
		}
		if token != nil && aws.ToString(resp.NextForwardToken) == *token {
			return messages, nil
		}
		token = resp.NextForwardToken
	}
	return nil, fmt.Errorf("too many pages reading %s/%s", group, stream)
}

func (tc *cwlogsTestCtx) httpIngestionTests() []TestResult {
	var results []TestResult

	// The bearer-token lifecycle end to end: the token is an IAM
	// service-specific credential of logs.amazonaws.com (the guide's own
	// issuance path); bearer requests against a group without the switch
	// answer 403; PutBearerTokenAuthentication enables them; disabling or
	// deactivating the credential closes the endpoint again.
	results = append(results, tc.runner.RunTest("logs", "HTTPIngestion_BearerTokenLifecycle", func() error {
		iamCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: tc.runner.endpoint,
			Region:   tc.runner.region,
		})
		if err != nil {
			return err
		}
		iamClient := iam.NewFromConfig(iamCfg)
		userName := tc.uniquePrefix("cwl-bearer-user")
		if _, err := iamClient.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(userName)}); err != nil {
			return fmt.Errorf("create IAM user: %v", err)
		}
		defer func() {
			_, _ = iamClient.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)})
		}()
		created, err := iamClient.CreateServiceSpecificCredential(tc.ctx, &iam.CreateServiceSpecificCredentialInput{
			UserName:    aws.String(userName),
			ServiceName: aws.String("logs.amazonaws.com"),
		})
		if err != nil {
			return fmt.Errorf("create service-specific credential: %v", err)
		}
		credential := created.ServiceSpecificCredential
		defer func() {
			_, _ = iamClient.DeleteServiceSpecificCredential(tc.ctx, &iam.DeleteServiceSpecificCredentialInput{
				UserName:                    aws.String(userName),
				ServiceSpecificCredentialId: credential.ServiceSpecificCredentialId,
			})
		}()
		token := aws.ToString(credential.ServicePassword)

		group, cleanupGroup, err := tc.newGroupStreamFixture("cwl-bearer-group", "bearer")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		const body = `{"message":"bearer hello"}`
		reqPath := fmt.Sprintf("/ingest/bulk?logGroup=%s&logStream=bearer", group)

		// An invalid token answers 401 before anything else.
		code, _, err := tc.ingestionPost(reqPath, "application/x-ndjson", body, "not-a-real-token")
		if err != nil {
			return err
		}
		if code != http.StatusUnauthorized {
			return fmt.Errorf("invalid token status %d, want 401", code)
		}

		// A valid token against a group without the switch answers 403.
		code, _, err = tc.ingestionPost(reqPath, "application/x-ndjson", body, token)
		if err != nil {
			return err
		}
		if code != http.StatusForbidden {
			return fmt.Errorf("flag-disabled status %d, want 403", code)
		}

		// PutBearerTokenAuthentication (the SDK surface of the op) opens
		// the group to bearer tokens; the same request ingests.
		if _, err := tc.client.PutBearerTokenAuthentication(tc.ctx, &cloudwatchlogs.PutBearerTokenAuthenticationInput{
			LogGroupIdentifier:               aws.String(group),
			BearerTokenAuthenticationEnabled: aws.Bool(true),
		}); err != nil {
			return fmt.Errorf("put bearer token authentication: %v", err)
		}
		code, respBody, err := tc.ingestionPost(reqPath, "application/x-ndjson", body, token)
		if err != nil {
			return err
		}
		if code != http.StatusOK || strings.TrimSpace(respBody) != "{}" {
			return fmt.Errorf("enabled status %d body %s", code, respBody)
		}

		// An inactive credential is an invalid token: 401.
		if _, err := iamClient.UpdateServiceSpecificCredential(tc.ctx, &iam.UpdateServiceSpecificCredentialInput{
			UserName:                    aws.String(userName),
			ServiceSpecificCredentialId: credential.ServiceSpecificCredentialId,
			Status:                      iamtypes.StatusTypeInactive,
		}); err != nil {
			return fmt.Errorf("deactivate credential: %v", err)
		}
		code, _, err = tc.ingestionPost(reqPath, "application/x-ndjson", body, token)
		if err != nil {
			return err
		}
		if code != http.StatusUnauthorized {
			return fmt.Errorf("inactive credential status %d, want 401", code)
		}

		// Disabling the switch closes the group even to a reactivated
		// token.
		if _, err := iamClient.UpdateServiceSpecificCredential(tc.ctx, &iam.UpdateServiceSpecificCredentialInput{
			UserName:                    aws.String(userName),
			ServiceSpecificCredentialId: credential.ServiceSpecificCredentialId,
			Status:                      iamtypes.StatusTypeActive,
		}); err != nil {
			return fmt.Errorf("reactivate credential: %v", err)
		}
		if _, err := tc.client.PutBearerTokenAuthentication(tc.ctx, &cloudwatchlogs.PutBearerTokenAuthenticationInput{
			LogGroupIdentifier:               aws.String(group),
			BearerTokenAuthenticationEnabled: aws.Bool(false),
		}); err != nil {
			return fmt.Errorf("disable bearer token authentication: %v", err)
		}
		code, _, err = tc.ingestionPost(reqPath, "application/x-ndjson", body, token)
		if err != nil {
			return err
		}
		if code != http.StatusForbidden {
			return fmt.Errorf("disabled status %d, want 403", code)
		}

		messages, err := tc.ingestionReadEvents(group, "bearer")
		if err != nil {
			return err
		}
		if len(messages) != 1 || messages[0] != `{"message":"bearer hello"}` {
			return fmt.Errorf("ingested %v", messages)
		}
		return nil
	}))

	// The four documented request shapes on one group, unauthenticated
	// (the platform runs with signature verification disabled): each
	// format's parsing rules and timestamps verified through GetLogEvents.
	results = append(results, tc.runner.RunTest("logs", "HTTPIngestion_FourFormats", func() error {
		group, cleanupGroup, err := tc.newLogGroupFixture("cwl-http-formats")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		for _, stream := range []string{"ndjson", "json", "hlc", "otlp"} {
			if err := tc.createLogStream(group, stream); err != nil {
				return err
			}
		}

		// ND-JSON: an object with its own timestamp, an invalid line
		// (skipped and reported), a bare string, and a number.
		ts := time.Now().UnixMilli() - 5000
		ndjsonBody := fmt.Sprintf("{\"timestamp\":%d,\"message\":\"nd one\"}\nnot json at all\n\"nd string\"\n42\n", ts)
		code, respBody, err := tc.ingestionPost(
			"/ingest/bulk?logGroup="+group+"&logStream=ndjson", "application/x-ndjson", ndjsonBody, "")
		if err != nil {
			return err
		}
		if code != http.StatusOK {
			return fmt.Errorf("ndjson status %d body %s", code, respBody)
		}
		var partial struct {
			PartialSuccess struct {
				RejectedLogRecords int    `json:"rejectedLogRecords"`
				ErrorMessage       string `json:"errorMessage"`
			} `json:"partialSuccess"`
		}
		if err := json.Unmarshal([]byte(respBody), &partial); err != nil {
			return fmt.Errorf("ndjson body %s: %v", respBody, err)
		}
		if partial.PartialSuccess.RejectedLogRecords != 1 {
			return fmt.Errorf("ndjson rejectedLogRecords = %d, want 1", partial.PartialSuccess.RejectedLogRecords)
		}

		// Structured JSON: an array whose non-object elements skip.
		code, respBody, err = tc.ingestionPost(
			"/ingest/json?logGroup="+group+"&logStream=json", "application/json",
			`[{"message":"sj one"},"skip me",{"message":"sj two"}]`, "")
		if err != nil {
			return err
		}
		if code != http.StatusOK {
			return fmt.Errorf("json status %d body %s", code, respBody)
		}
		partial = struct {
			PartialSuccess struct {
				RejectedLogRecords int    `json:"rejectedLogRecords"`
				ErrorMessage       string `json:"errorMessage"`
			} `json:"partialSuccess"`
		}{}
		_ = json.Unmarshal([]byte(respBody), &partial)
		if partial.PartialSuccess.RejectedLogRecords != 1 {
			return fmt.Errorf("json rejectedLogRecords = %d, want 1", partial.PartialSuccess.RejectedLogRecords)
		}

		// HLC: fractional epoch-seconds time; a message-less object
		// skips without any report (the format has no partial success).
		hlcBase := time.Now().Unix() - 5
		hlcBody := fmt.Sprintf(`{"event":"hlc one","time":%d.5}{"message":"no event"}`, hlcBase)
		code, respBody, err = tc.ingestionPost(
			"/services/collector/event?logGroup="+group+"&logStream=hlc", "application/json", hlcBody, "")
		if err != nil {
			return err
		}
		if code != http.StatusOK || strings.TrimSpace(respBody) != "{}" {
			return fmt.Errorf("hlc status %d body %s", code, respBody)
		}

		// OTLP JSON: header addressing only; nanosecond timestamps.
		nanos := time.Now().Add(-3 * time.Second).UnixNano()
		otlpBody := fmt.Sprintf(`{"resourceLogs":[{"scopeLogs":[{"logRecords":[{"timeUnixNano":"%d","body":{"stringValue":"otlp one"}}]}]}]}`, nanos)
		otlpURL := "/v1/logs?logGroup=ignored&logStream=ignored"
		req, err := http.NewRequest(http.MethodPost, tc.runner.endpoint+otlpURL, strings.NewReader(otlpBody))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-aws-log-group", group)
		req.Header.Set("x-aws-log-stream", "otlp")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		otlpBodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK || strings.TrimSpace(string(otlpBodyBytes)) != "{}" {
			return fmt.Errorf("otlp status %d body %s", resp.StatusCode, string(otlpBodyBytes))
		}

		expected := map[string][]string{
			"ndjson": {fmt.Sprintf(`{"message":"nd one","timestamp":%d}`, ts), "nd string", "42"},
			"json":   {`{"message":"sj one"}`, `{"message":"sj two"}`},
			"hlc":    {"hlc one"},
			"otlp":   {"otlp one"},
		}
		for stream, wants := range expected {
			messages, err := tc.ingestionReadEvents(group, stream)
			if err != nil {
				return err
			}
			if len(messages) != len(wants) {
				return fmt.Errorf("stream %s ingested %v, want %v", stream, messages, wants)
			}
			for i, want := range wants {
				if messages[i] != want {
					return fmt.Errorf("stream %s event %d = %q, want %q", stream, i, messages[i], want)
				}
			}
		}
		return nil
	}))

	// The validation rows: missing or duplicated target parameters, an
	// unknown group or stream, an oversized request, and the ND-JSON
	// all-invalid request.
	results = append(results, tc.runner.RunTest("logs", "HTTPIngestion_RejectionRows", func() error {
		group, cleanupGroup, err := tc.newGroupStreamFixture("cwl-http-reject", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		code, body, err := tc.ingestionPost("/ingest/bulk?logStream=s1", "application/x-ndjson", `{"message":"x"}`, "")
		if err != nil {
			return err
		}
		if code != http.StatusBadRequest {
			return fmt.Errorf("missing logGroup status %d body %s", code, body)
		}

		code, body, err = tc.ingestionPost("/ingest/bulk?logGroup="+group+"&logStream=s1", "application/x-ndjson", `{"message":"x"}`, "")
		if err != nil {
			return err
		}
		if code != http.StatusOK {
			return fmt.Errorf("baseline request status %d body %s", code, body)
		}
		headerTarget := fmt.Sprintf("/ingest/bulk?logGroup=%s&logStream=s1", group)
		req, err := http.NewRequest(http.MethodPost, tc.runner.endpoint+headerTarget, strings.NewReader(`{"message":"x"}`))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-ndjson")
		req.Header.Set("x-aws-log-group", group)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("both-ways logGroup status %d, want 400", resp.StatusCode)
		}

		code, body, err = tc.ingestionPost("/ingest/bulk?logGroup="+group+"&logStream=no-such-stream", "application/x-ndjson", `{"message":"x"}`, "")
		if err != nil {
			return err
		}
		if code != http.StatusNotFound {
			return fmt.Errorf("unknown stream status %d body %s", code, body)
		}

		code, body, err = tc.ingestionPost("/ingest/bulk?logGroup="+group+"&logStream=s1", "application/x-ndjson", "not json\nalso not\n", "")
		if err != nil {
			return err
		}
		if code != http.StatusBadRequest || !strings.Contains(body, "All events were invalid") {
			return fmt.Errorf("all-invalid status %d body %s", code, body)
		}

		// An empty-message record rejects the whole request before any
		// chunk commits: the records deliberately span more than the
		// event-time span, so the chunker would otherwise flush the first
		// (valid) record before the seam rejects the empty one — the
		// baseline event above must remain the stream's only content.
		// The empty message is a bare empty JSON string line (an ND-JSON
		// object's message is its own full serialisation, never empty).
		spanOld := time.Now().Add(-25*time.Hour - 2*time.Second).UnixMilli()
		spanNew := time.Now().Add(-2 * time.Second).UnixMilli()
		emptyMsg := fmt.Sprintf("{\"timestamp\":%d,\"message\":\"kept\"}\n\"\"\n{\"timestamp\":%d,\"message\":\"late\"}\n",
			spanOld, spanNew)
		code, body, err = tc.ingestionPost("/ingest/bulk?logGroup="+group+"&logStream=s1", "application/x-ndjson", emptyMsg, "")
		if err != nil {
			return err
		}
		if code != http.StatusBadRequest || !strings.Contains(body, "message member") {
			return fmt.Errorf("empty-message status %d body %s", code, body)
		}
		messages, err := tc.ingestionReadEvents(group, "s1")
		if err != nil {
			return err
		}
		if len(messages) != 1 || messages[0] != `{"message":"x"}` {
			return fmt.Errorf("after empty-message rejection the stream holds %v, want only the baseline event", messages)
		}

		oversized := "[" + strings.Repeat(`{"message":"pad"},`, 210000) + `{}]`
		code, body, err = tc.ingestionPost("/ingest/json?logGroup="+group+"&logStream=s1", "application/json", oversized, "")
		if err != nil {
			return err
		}
		if code != http.StatusBadRequest {
			return fmt.Errorf("oversized request status %d body %s", code, body)
		}
		return nil
	}))

	return results
}
