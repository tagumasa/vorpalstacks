package iot

import (
	"testing"
)

func TestPaginatedStrings(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}

	t.Run("no pagination returns all", func(t *testing.T) {
		resp := paginatedStrings("items", items, nil)
		list := resp["items"].([]string)
		if len(list) != 5 {
			t.Fatalf("expected 5 items, got %d", len(list))
		}
		if _, ok := resp["nextToken"]; ok {
			t.Fatal("unexpected nextToken")
		}
	})

	t.Run("maxResults limits page", func(t *testing.T) {
		params := map[string]interface{}{"maxResults": float64(2)}
		resp := paginatedStrings("items", items, params)
		list := resp["items"].([]string)
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if resp["nextToken"] != "2" {
			t.Fatalf("expected nextToken 2, got %v", resp["nextToken"])
		}
	})

	t.Run("nextToken offsets correctly", func(t *testing.T) {
		params := map[string]interface{}{
			"nextToken":  "3",
			"maxResults": float64(2),
		}
		resp := paginatedStrings("items", items, params)
		list := resp["items"].([]string)
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if list[0] != "d" || list[1] != "e" {
			t.Fatalf("expected [d, e], got %v", list)
		}
		if _, ok := resp["nextToken"]; ok {
			t.Fatal("unexpected nextToken for last page")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		resp := paginatedStrings("items", nil, nil)
		list := resp["items"].([]string)
		if len(list) != 0 {
			t.Fatalf("expected 0 items, got %d", len(list))
		}
	})
}

func TestPaginatedMaps(t *testing.T) {
	items := []map[string]interface{}{
		{"key": "a"}, {"key": "b"}, {"key": "c"},
	}

	t.Run("maxResults limits page", func(t *testing.T) {
		params := map[string]interface{}{"maxResults": float64(2)}
		resp := paginatedMaps("groups", items, params)
		list := resp["groups"].([]map[string]interface{})
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if resp["nextToken"] != "2" {
			t.Fatalf("expected nextToken 2, got %v", resp["nextToken"])
		}
	})
}
