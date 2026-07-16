package iot

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
)

func TestDescribeEndpointMQTTTypes(t *testing.T) {
	s := NewIoTService("test-account")
	req := &request.ParsedRequest{Parameters: map[string]interface{}{}}
	expected := "localhost:" + strconv.Itoa(serviceports.IotMQTT)

	for _, ep := range []string{"iot:Data", "iot:Data-ATS", "iot:Data-ALPN", "iot:Jobs"} {
		req.Parameters["endpointType"] = ep
		resp, err := s.DescribeEndpoint(context.Background(), nil, req)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", ep, err)
		}
		m, ok := resp.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: unexpected response type %T", ep, resp)
		}
		addr, ok := m["endpointAddress"].(string)
		if !ok {
			t.Fatalf("%s: endpointAddress missing or not string", ep)
		}
		if addr != expected {
			t.Errorf("%s: endpointAddress = %q, want %q", ep, addr, expected)
		}
	}
}

func TestDescribeEndpointDefault(t *testing.T) {
	s := NewIoTService("test-account")
	req := &request.ParsedRequest{Parameters: map[string]interface{}{}}

	resp, err := s.DescribeEndpoint(context.Background(), nil, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	addr := resp.(map[string]interface{})["endpointAddress"].(string)
	if !strings.HasPrefix(addr, "localhost:") {
		t.Fatalf("default endpoint should be localhost:{port}, got %q", addr)
	}
}

func TestDescribeEndpointCredentialProvider(t *testing.T) {
	s := NewIoTService("test-account")
	req := &request.ParsedRequest{
		Parameters: map[string]interface{}{"endpointType": "iot:CredentialProvider"},
	}

	resp, err := s.DescribeEndpoint(context.Background(), nil, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	addr := resp.(map[string]interface{})["endpointAddress"].(string)
	if !strings.HasPrefix(addr, "localhost:") {
		t.Fatalf("CredentialProvider endpoint should be localhost:{port}, got %q", addr)
	}
}

func TestDescribeEndpointUnknownType(t *testing.T) {
	s := NewIoTService("test-account")
	req := &request.ParsedRequest{
		Parameters: map[string]interface{}{"endpointType": "iot:Bogus"},
	}

	resp, err := s.DescribeEndpoint(context.Background(), nil, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	addr := resp.(map[string]interface{})["endpointAddress"].(string)
	if !strings.HasPrefix(addr, "localhost:") {
		t.Fatalf("unknown endpoint type should still return localhost:{port}, got %q", addr)
	}
}
