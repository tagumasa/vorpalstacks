package iot

import (
	"errors"
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

func TestPaginatedStrings(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}

	t.Run("no pagination returns all", func(t *testing.T) {
		resp, err := paginatedStrings("items", items, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list := resp.(map[string]interface{})["items"].([]string)
		if len(list) != 5 {
			t.Fatalf("expected 5 items, got %d", len(list))
		}
		if _, ok := resp.(map[string]interface{})["nextToken"]; ok {
			t.Fatal("unexpected nextToken")
		}
	})

	t.Run("maxResults limits page", func(t *testing.T) {
		params := map[string]interface{}{"maxResults": float64(2)}
		resp, err := paginatedStrings("items", items, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list := resp.(map[string]interface{})["items"].([]string)
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if resp.(map[string]interface{})["nextToken"] != "2" {
			t.Fatalf("expected nextToken 2, got %v", resp.(map[string]interface{})["nextToken"])
		}
	})

	t.Run("nextToken offsets correctly", func(t *testing.T) {
		params := map[string]interface{}{
			"nextToken":  "3",
			"maxResults": float64(2),
		}
		resp, err := paginatedStrings("items", items, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list := resp.(map[string]interface{})["items"].([]string)
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if list[0] != "d" || list[1] != "e" {
			t.Fatalf("expected [d, e], got %v", list)
		}
		if _, ok := resp.(map[string]interface{})["nextToken"]; ok {
			t.Fatal("unexpected nextToken for last page")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		resp, err := paginatedStrings("items", nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list := resp.(map[string]interface{})["items"].([]string)
		if len(list) != 0 {
			t.Fatalf("expected 0 items, got %d", len(list))
		}
	})

	// A token the server could never have issued must not silently
	// restart the walk at the first page; the iot list operations report
	// it as InvalidRequestException.
	invalidTokens := []struct {
		name  string
		token string
	}{
		{"non-numeric", "abc"},
		{"negative", "-1"},
		{"beyond the result set", "9999"},
		{"equal to the result-set length", "5"},
	}
	for _, tt := range invalidTokens {
		t.Run("invalid token "+tt.name, func(t *testing.T) {
			params := map[string]interface{}{"nextToken": tt.token}
			if _, err := paginatedStrings("items", items, params); !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
			if _, err := paginatedMaps("items", mapsFixture(), params); !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest from paginatedMaps, got %v", err)
			}
		})
	}
}

func TestPaginatedMaps(t *testing.T) {
	items := mapsFixture()

	t.Run("maxResults limits page", func(t *testing.T) {
		params := map[string]interface{}{"maxResults": float64(2)}
		resp, err := paginatedMaps("groups", items, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list := resp.(map[string]interface{})["groups"].([]map[string]interface{})
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if resp.(map[string]interface{})["nextToken"] != "2" {
			t.Fatalf("expected nextToken 2, got %v", resp.(map[string]interface{})["nextToken"])
		}
	})
}

func mapsFixture() []map[string]interface{} {
	return []map[string]interface{}{
		{"key": "a"}, {"key": "b"}, {"key": "c"},
	}
}
