package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const wsTestEndpoint = "ws://127.0.0.1:8086/event/realtime"
const wsTestHTTPEndpoint = "http://127.0.0.1:8086/event"

func (r *TestRunner) RunAppSyncWSTests() []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("appsync-ws", "WebSocket_ConnectionAck", func() error {
		conn, err := dialWS()
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		defer conn.Close()

		msg, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return fmt.Errorf("failed to read connection_ack: %w", err)
		}
		if msg["type"] != "connection_ack" {
			return fmt.Errorf("expected connection_ack, got %v", msg["type"])
		}
		if _, ok := msg["connectionTimeoutMs"]; !ok {
			return fmt.Errorf("connection_ack missing connectionTimeoutMs")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_SubscribeSuccess", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		subMsg := map[string]interface{}{
			"type":    "subscribe",
			"id":      "sub-1",
			"channel": "/default/test",
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			return fmt.Errorf("write subscribe failed: %w", err)
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return fmt.Errorf("failed to read subscribe response: %w", err)
		}
		if resp["type"] != "subscribe_success" {
			return fmt.Errorf("expected subscribe_success, got %v", resp["type"])
		}
		if resp["id"] != "sub-1" {
			return fmt.Errorf("expected id sub-1, got %v", resp["id"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_SubscribeError_InvalidChannel", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		subMsg := map[string]interface{}{
			"type":    "subscribe",
			"id":      "sub-bad-ch",
			"channel": "",
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "subscribe_error" {
			return fmt.Errorf("expected subscribe_error, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_SubscribeError_DuplicateId", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		subMsg := map[string]interface{}{
			"type":    "subscribe",
			"id":      "sub-dup",
			"channel": "/default/ch1",
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		if err := conn.WriteJSON(subMsg); err != nil {
			return err
		}
		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "subscribe_error" {
			return fmt.Errorf("expected subscribe_error for duplicate, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_SubscribeError_InvalidSubId", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		subMsg := map[string]interface{}{
			"type":    "subscribe",
			"id":      "",
			"channel": "/default/test",
		}
		if err := conn.WriteJSON(subMsg); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "subscribe_error" {
			return fmt.Errorf("expected subscribe_error, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_PublishAndReceiveData", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		ch := fmt.Sprintf("/default/pubtest-%d", time.Now().UnixNano())

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "sub-pub",
			"channel": ch,
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return fmt.Errorf("subscribe failed: %w", err)
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "publish",
			"id":      "pub-1",
			"channel": ch,
			"events":  []string{`{"msg":"hello"}`},
		}); err != nil {
			return err
		}

		msgs, err := readMessages(conn, 3*time.Second, 2)
		if err != nil {
			return fmt.Errorf("failed to read responses: %w", err)
		}

		var pubResp, dataMsg map[string]interface{}
		for _, m := range msgs {
			switch m["type"] {
			case "publish_success":
				pubResp = m
			case "data":
				dataMsg = m
			}
		}
		if pubResp == nil {
			return fmt.Errorf("expected publish_success, got types: %v", messageTypes(msgs))
		}
		if dataMsg == nil {
			return fmt.Errorf("expected data message, got types: %v", messageTypes(msgs))
		}
		if dataMsg["type"] != "data" {
			return fmt.Errorf("expected data message, got type %v", dataMsg["type"])
		}
		if dataMsg["id"] != "sub-pub" {
			return fmt.Errorf("expected data id sub-pub, got %v", dataMsg["id"])
		}
		eventStr, ok := dataMsg["event"].(string)
		if !ok {
			return fmt.Errorf("event field is not a string: %T", dataMsg["event"])
		}
		var events []json.RawMessage
		if err := json.Unmarshal([]byte(eventStr), &events); err != nil {
			return fmt.Errorf("event field is not valid JSON array: %w", err)
		}
		if len(events) != 1 {
			return fmt.Errorf("expected 1 event, got %d", len(events))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_PublishError_EmptyEvents", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "publish",
			"id":      "pub-empty",
			"channel": "/default/test",
			"events":  []string{},
		}); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "publish_error" {
			return fmt.Errorf("expected publish_error, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_UnsubscribeSuccess", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "sub-unsub",
			"channel": "/default/unsub-test",
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type": "unsubscribe",
			"id":   "sub-unsub",
		}); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "unsubscribe_success" {
			return fmt.Errorf("expected unsubscribe_success, got %v", resp["type"])
		}
		if resp["id"] != "sub-unsub" {
			return fmt.Errorf("expected id sub-unsub, got %v", resp["id"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_UnsubscribeError_UnknownId", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type": "unsubscribe",
			"id":   "nonexistent",
		}); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "unsubscribe_error" {
			return fmt.Errorf("expected unsubscribe_error, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_ConnectionInit_Accepted", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type": "connection_init",
		}); err != nil {
			return err
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, err = conn.ReadMessage()
		if err == nil {
			return fmt.Errorf("expected no response to connection_init (server should accept silently), but got a message")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_UnknownMessageType", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type": "bogus_type",
		}); err != nil {
			return err
		}

		resp, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return err
		}
		if resp["type"] != "error" {
			return fmt.Errorf("expected error, got %v", resp["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_MultiSubscriberFanOut", func() error {
		conn1, err := dialWS()
		if err != nil {
			return err
		}
		defer conn1.Close()

		conn2, err := dialWS()
		if err != nil {
			return err
		}
		defer conn2.Close()

		if err := drainAck(conn1); err != nil {
			return err
		}
		if err := drainAck(conn2); err != nil {
			return err
		}

		ch := fmt.Sprintf("/default/fanout-%d", time.Now().UnixNano())

		if err := conn1.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "fan-sub-1",
			"channel": ch,
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn1, 3*time.Second); err != nil {
			return err
		}

		if err := conn2.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "fan-sub-2",
			"channel": ch,
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn2, 3*time.Second); err != nil {
			return err
		}

		if err := conn1.WriteJSON(map[string]interface{}{
			"type":    "publish",
			"id":      "fan-pub",
			"channel": ch,
			"events":  []string{`{"fan":"out"}`},
		}); err != nil {
			return err
		}

		msgs1, err := readMessages(conn1, 3*time.Second, 2)
		if err != nil {
			return fmt.Errorf("conn1 failed to receive messages: %w", err)
		}
		var foundPubAck, foundData1 bool
		for _, m := range msgs1 {
			switch m["type"] {
			case "publish_success":
				foundPubAck = true
			case "data":
				foundData1 = true
			}
		}
		if !foundPubAck {
			return fmt.Errorf("conn1: expected publish_success, got types: %v", messageTypes(msgs1))
		}
		if !foundData1 {
			return fmt.Errorf("conn1: expected data, got types: %v", messageTypes(msgs1))
		}

		data2, err := readMessage(conn2, 3*time.Second)
		if err != nil {
			return fmt.Errorf("conn2 failed to receive data: %w", err)
		}
		if data2["type"] != "data" {
			return fmt.Errorf("conn2: expected data, got %v", data2["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_WildcardChannel", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		ns := fmt.Sprintf("wildns-%d", time.Now().UnixNano())

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "wild-sub",
			"channel": "/" + ns + "/*",
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "publish",
			"id":      "wild-pub",
			"channel": "/" + ns + "/topic1",
			"events":  []string{`{"wildcard":true}`},
		}); err != nil {
			return err
		}

		msgs, err := readMessages(conn, 3*time.Second, 2)
		if err != nil {
			return fmt.Errorf("wildcard: failed to read responses: %w", err)
		}

		var foundPublish, foundData bool
		for _, m := range msgs {
			switch m["type"] {
			case "publish_success":
				foundPublish = true
			case "data":
				foundData = true
			}
		}
		if !foundPublish {
			return fmt.Errorf("wildcard: expected publish_success, got types: %v", messageTypes(msgs))
		}
		if !foundData {
			return fmt.Errorf("wildcard: expected data, got types: %v", messageTypes(msgs))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "WebSocket_UnsubscribeStopsDelivery", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		ch := fmt.Sprintf("/default/unsubstop-%d", time.Now().UnixNano())

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "stop-sub",
			"channel": ch,
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type": "unsubscribe",
			"id":   "stop-sub",
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "publish",
			"id":      "stop-pub",
			"channel": ch,
			"events":  []string{`{"after":"unsub"}`},
		}); err != nil {
			return err
		}

		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return fmt.Errorf("should get publish_success: %w", err)
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, err = conn.ReadMessage()
		if err == nil {
			return fmt.Errorf("expected timeout after unsubscribe, but got a message")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "HTTP_Publish_Success", func() error {
		conn, err := dialWS()
		if err != nil {
			return err
		}
		defer conn.Close()

		if err := drainAck(conn); err != nil {
			return err
		}

		ch := fmt.Sprintf("/default/httppub-%d", time.Now().UnixNano())

		if err := conn.WriteJSON(map[string]interface{}{
			"type":    "subscribe",
			"id":      "http-sub",
			"channel": ch,
		}); err != nil {
			return err
		}
		if _, err := readMessage(conn, 3*time.Second); err != nil {
			return err
		}

		bodyBytes, _ := json.Marshal(map[string]interface{}{
			"channel": ch,
			"events":  []string{`{"via":"http"}`},
		})
		httpResp, err := r.client.Post(wsTestHTTPEndpoint, "application/json", strings.NewReader(string(bodyBytes)))
		if err != nil {
			return fmt.Errorf("HTTP publish failed: %w", err)
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP publish returned status %d", httpResp.StatusCode)
		}

		var httpResult map[string]interface{}
		if err := json.NewDecoder(httpResp.Body).Decode(&httpResult); err != nil {
			return fmt.Errorf("failed to decode HTTP response: %w", err)
		}
		successful, ok := httpResult["successful"].([]interface{})
		if !ok || len(successful) != 1 {
			return fmt.Errorf("expected 1 successful event, got %v", httpResult["successful"])
		}

		dataMsg, err := readMessage(conn, 3*time.Second)
		if err != nil {
			return fmt.Errorf("WS subscriber failed to receive HTTP-published event: %w", err)
		}
		if dataMsg["type"] != "data" {
			return fmt.Errorf("expected data, got %v", dataMsg["type"])
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "HTTP_Publish_Error_EmptyChannel", func() error {
		httpResp, err := r.client.Post(wsTestHTTPEndpoint, "application/json", strings.NewReader(`{"channel":"","events":["{}"]}`))
		if err != nil {
			return err
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("expected 400, got %d", httpResp.StatusCode)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync-ws", "HTTP_Publish_Error_TooManyEvents", func() error {
		evts := make([]string, 6)
		for i := range evts {
			evts[i] = "{}"
		}
		body, _ := json.Marshal(map[string]interface{}{
			"channel": "/default/test",
			"events":  evts,
		})
		httpResp, err := r.client.Post(wsTestHTTPEndpoint, "application/json", strings.NewReader(string(body)))
		if err != nil {
			return err
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("expected 400 for too many events, got %d", httpResp.StatusCode)
		}
		return nil
	}))

	return results
}

func dialWS() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsTestEndpoint, http.Header{
		"Sec-WebSocket-Protocol": []string{"aws-appsync-event-ws"},
	})
	return conn, err
}

func drainAck(conn *websocket.Conn) error {
	msg, err := readMessage(conn, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to drain connection_ack: %w", err)
	}
	if msg["type"] != "connection_ack" {
		return fmt.Errorf("expected connection_ack during drain, got %v", msg["type"])
	}
	return nil
}

func readMessage(conn *websocket.Conn, timeout time.Duration) (map[string]interface{}, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var msg map[string]interface{}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return msg, nil
}

func readMessages(conn *websocket.Conn, timeout time.Duration, count int) ([]map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	var msgs []map[string]interface{}
	for len(msgs) < count {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		msg, err := readMessage(conn, remaining)
		if err != nil {
			break
		}
		msgs = append(msgs, msg)
	}
	if len(msgs) < count {
		return msgs, fmt.Errorf("expected %d messages, got %d", count, len(msgs))
	}
	return msgs, nil
}

func messageTypes(msgs []map[string]interface{}) []string {
	types := make([]string, len(msgs))
	for i, m := range msgs {
		if t, ok := m["type"].(string); ok {
			types[i] = t
		}
	}
	return types
}
