package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// RunNeptuneDirectProtocolTests runs integration tests that connect directly
// to a Neptune cluster's data plane via raw HTTP POST and WebSocket, bypassing
// the AWS SDK. These tests verify the protocol-level behaviour that SDK tests
// cannot exercise (e.g. TinkerPop Gremlin Server WebSocket protocol).
func (r *TestRunner) RunNeptuneDirectProtocolTests() []TestResult {
	var results []TestResult

	clusterID := fmt.Sprintf("integ-ws-%d", time.Now().UnixNano())
	clusterPort, cleanup := r.createNeptuneClusterForDataTests(clusterID)
	defer cleanup()

	if clusterPort == 0 {
		return append(results, SetupFailResult("integration",
			"NeptuneDirectProtocol: failed to create cluster"))
	}

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", clusterPort)

	results = append(results, r.runNeptuneCypherHTTPTests(baseURL)...)
	results = append(results, r.runNeptuneGremlinHTTPTests(baseURL)...)
	results = append(results, r.runNeptuneGremlinWSTests(baseURL)...)

	return results
}

// --- Cypher HTTP direct tests ---

func (r *TestRunner) runNeptuneCypherHTTPTests(baseURL string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("integration", "CypherHTTP_CreateNode", func() error {
		return postCypher(baseURL, "CREATE (n:Person {name: 'alice'}) RETURN n")
	}))

	results = append(results, r.RunTest("integration", "CypherHTTP_MatchNode", func() error {
		return postCypher(baseURL, "MATCH (n:Person) WHERE n.name = 'alice' RETURN n")
	}))

	results = append(results, r.RunTest("integration", "CypherHTTP_CreateEdge", func() error {
		if err := postCypher(baseURL, "CREATE (a:Person {name: 'bob'})"); err != nil {
			return err
		}
		if err := postCypher(baseURL, "CREATE (b:Person {name: 'carol'})"); err != nil {
			return err
		}
		return postCypher(baseURL, "MATCH (a:Person {name: 'bob'}), (b:Person {name: 'carol'}) CREATE (a)-[:KNOWS]->(b)")
	}))

	results = append(results, r.RunTest("integration", "CypherHTTP_InvalidQuery", func() error {
		_, err := postRaw(baseURL+"/opencypher", map[string]string{"query": "INVALID CYPHER !!!!", "Accept": "application/json"})
		if err == nil {
			return fmt.Errorf("expected error for invalid cypher query")
		}
		return nil
	}))

	return results
}

// --- Gremlin HTTP direct tests ---

func (r *TestRunner) runNeptuneGremlinHTTPTests(baseURL string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("integration", "GremlinHTTP_Count", func() error {
		return postGremlin(baseURL, "g.V().count()")
	}))

	results = append(results, r.RunTest("integration", "GremlinHTTP_AddV", func() error {
		return postGremlin(baseURL, "g.addV('httptest').property('name','direct')")
	}))

	results = append(results, r.RunTest("integration", "GremlinHTTP_HasName", func() error {
		return postGremlin(baseURL, "g.V().has('name','direct')")
	}))

	results = append(results, r.RunTest("integration", "GremlinHTTP_InvalidQuery", func() error {
		_, err := postRaw(baseURL+"/gremlin", map[string]string{"gremlin": "INVALID GREMLIN !!!!", "Accept": "application/json"})
		if err == nil {
			return fmt.Errorf("expected error for invalid gremlin query")
		}
		return nil
	}))

	return results
}

// --- Gremlin WebSocket tests ---

func (r *TestRunner) runNeptuneGremlinWSTests(baseURL string) []TestResult {
	var results []TestResult

	wsURL := "ws" + strings.TrimPrefix(baseURL, "http") + "/gremlin"

	results = append(results, r.RunTest("integration", "GremlinWS_ConnectAndEval", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		resp, err := evalGremlinWS(conn, "g.V().count()")
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("expected status 200, got %d: %s", resp.Status.Code, resp.Status.Message)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_AddVertex", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		resp, err := evalGremlinWS(conn, "g.addV('wstest').property('name','ws-alice')")
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("expected status 200, got %d: %s", resp.Status.Code, resp.Status.Message)
		}
		return verifyVertexInResult(resp)
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_AddEdge", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		_, err = evalGremlinWS(conn, "g.addV('ws-a').property('name','x')")
		if err != nil {
			return fmt.Errorf("addV x failed: %w", err)
		}
		_, err = evalGremlinWS(conn, "g.addV('ws-b').property('name','y')")
		if err != nil {
			return fmt.Errorf("addV y failed: %w", err)
		}

		resp, err := evalGremlinWS(conn, "g.V().has('name','x').addE('knows').to(g.V().has('name','y')).property('weight',1.5)")
		if err != nil {
			return fmt.Errorf("addE failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("expected status 200, got %d: %s", resp.Status.Code, resp.Status.Message)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_ErrorResponse", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		resp, err := evalGremlinWS(conn, "INVALID QUERY !!!!")
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}
		if resp.Status.Code != 598 {
			return fmt.Errorf("expected status 598 (script error), got %d", resp.Status.Code)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_GraphSON3Structure", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		resp, err := evalGremlinWS(conn, "g.addV('gsontest').property('key','val')")
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("expected status 200, got %d: %s", resp.Status.Code, resp.Status.Message)
		}

		resultMap, ok := resp.Result.Data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("result data is not a map, got %T", resp.Result.Data)
		}
		typeName, _ := resultMap["@type"].(string)
		if typeName != "g:List" {
			return fmt.Errorf("expected @type g:List, got %q", typeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_MultiQuery", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		for i := 0; i < 3; i++ {
			resp, err := evalGremlinWS(conn, "g.V().count()")
			if err != nil {
				return fmt.Errorf("query %d failed: %w", i, err)
			}
			if resp.Status.Code != 200 {
				return fmt.Errorf("query %d: expected status 200, got %d", i, resp.Status.Code)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_SessionAliases", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		// Send eval with processor="session" to establish a session
		req := &gremlinWSRequest{
			RequestID: "session-test-1",
			Op:        "eval",
			Processor: "session",
			Args: gremlinWSArgs{
				Gremlin: "g.addV('sesstest').property('name','sess-alice')",
				Aliases: map[string]string{"g": "g"},
			},
		}
		resp, err := sendRecvGremlinWS(conn, req)
		if err != nil {
			return fmt.Errorf("session eval 1 failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("session eval 1: expected 200, got %d", resp.Status.Code)
		}

		// Second eval in same session should work
		resp2, err := evalGremlinWS(conn, "g.V().has('name','sess-alice')")
		if err != nil {
			return fmt.Errorf("session eval 2 failed: %w", err)
		}
		if resp2.Status.Code != 200 {
			return fmt.Errorf("session eval 2: expected 200, got %d: %s", resp2.Status.Code, resp2.Status.Message)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_CloseAndReconnect", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("first dial failed: %w", err)
		}

		req := &gremlinWSRequest{
			RequestID: "close-test",
			Op:        "close",
			Processor: "",
		}
		resp, err := sendRecvGremlinWS(conn, req)
		if err != nil {
			return fmt.Errorf("close failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("close: expected 200, got %d", resp.Status.Code)
		}
		conn.Close()

		conn2, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("reconnect dial failed: %w", err)
		}
		defer conn2.Close()

		resp2, err := evalGremlinWS(conn2, "g.V().count()")
		if err != nil {
			return fmt.Errorf("reconnect eval failed: %w", err)
		}
		if resp2.Status.Code != 200 {
			return fmt.Errorf("reconnect eval: expected 200, got %d", resp2.Status.Code)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_Authentication", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		req := &gremlinWSRequest{
			RequestID: "auth-test",
			Op:        "authentication",
			Processor: "",
		}
		resp, err := sendRecvGremlinWS(conn, req)
		if err != nil {
			return fmt.Errorf("auth request failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("auth: expected 200, got %d", resp.Status.Code)
		}
		return nil
	}))

	results = append(results, r.RunTest("integration", "GremlinWS_HasFilter", func() error {
		conn, err := dialGremlinWS(wsURL)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		resp, err := evalGremlinWS(conn, "g.V().has('name','ws-alice')")
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}
		if resp.Status.Code != 200 {
			return fmt.Errorf("expected status 200, got %d: %s", resp.Status.Code, resp.Status.Message)
		}
		return verifyVertexInResult(resp)
	}))

	return results
}

// --- HTTP helpers ---

func postCypher(baseURL, query string) error {
	_, err := postRaw(baseURL+"/opencypher", map[string]string{"query": query})
	return err
}

func postGremlin(baseURL, query string) error {
	_, err := postRaw(baseURL+"/gremlin", map[string]string{"gremlin": query})
	return err
}

func postRaw(url string, body map[string]string) (map[string]interface{}, error) {
	jsonBody, _ := json.Marshal(body)
	resp, err := testHTTPClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("POST %s failed: %w", url, err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("POST %s returned %d: %s", url, resp.StatusCode, string(data))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}
	return result, nil
}

// --- WebSocket helpers ---

type gremlinWSRequest struct {
	RequestID string        `json:"requestId"`
	Op        string        `json:"op"`
	Processor string        `json:"processor"`
	Args      gremlinWSArgs `json:"args"`
}

type gremlinWSArgs struct {
	Gremlin  string                 `json:"gremlin,omitempty"`
	Language string                 `json:"language,omitempty"`
	Aliases  map[string]string      `json:"aliases,omitempty"`
	Bindings map[string]interface{} `json:"bindings,omitempty"`
	Session  string                 `json:"session,omitempty"`
}

type gremlinWSResponse struct {
	RequestID string      `json:"requestId"`
	Status    wsStatusMsg `json:"status"`
	Result    wsResultMsg `json:"result"`
}

type wsStatusMsg struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Attributes map[string]interface{} `json:"attributes"`
}

type wsResultMsg struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

func dialGremlinWS(wsURL string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ws dial %s: %w", wsURL, err)
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn, nil
}

func evalGremlinWS(conn *websocket.Conn, query string) (*gremlinWSResponse, error) {
	req := &gremlinWSRequest{
		RequestID: fmt.Sprintf("test-%d", time.Now().UnixNano()),
		Op:        "eval",
		Processor: "",
		Args: gremlinWSArgs{
			Gremlin:  query,
			Language: "gremlin-groovy",
		},
	}
	return sendRecvGremlinWS(conn, req)
}

func sendRecvGremlinWS(conn *websocket.Conn, req *gremlinWSRequest) (*gremlinWSResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var resp gremlinWSResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w (raw: %s)", err, string(msg))
	}
	return &resp, nil
}

// verifyVertexInResult checks that the GraphSON v3 result contains a g:Vertex.
func verifyVertexInResult(resp *gremlinWSResponse) error {
	resultMap, ok := resp.Result.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("result data is not a map")
	}
	typeName, _ := resultMap["@type"].(string)
	if typeName != "g:List" {
		return fmt.Errorf("expected g:List type, got %q", typeName)
	}
	val, _ := resultMap["@value"].([]interface{})
	if len(val) == 0 {
		return fmt.Errorf("expected at least one element in g:List result")
	}
	vertex, ok := val[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first element is not a map, got %T", val[0])
	}
	if vertex["@type"] != "g:Vertex" {
		return fmt.Errorf("expected g:Vertex, got %v", vertex["@type"])
	}
	return nil
}
