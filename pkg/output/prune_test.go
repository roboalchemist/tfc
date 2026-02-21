package output

import (
	"encoding/json"
	"testing"
)

func TestPruneEmpty_PreservesEmptyStrings(t *testing.T) {
	input := map[string]interface{}{
		"id":       "run-123",
		"plan_id":  "",
		"apply_id": "",
		"attributes": map[string]interface{}{
			"status": "pending",
		},
	}

	result := PruneEmpty(input)
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}

	if _, exists := m["plan_id"]; !exists {
		t.Error("plan_id should be preserved even when empty")
	}
	if _, exists := m["apply_id"]; !exists {
		t.Error("apply_id should be preserved even when empty")
	}
}

func TestPruneEmpty_RemovesNullAndEmptyContainers(t *testing.T) {
	input := map[string]interface{}{
		"id":         "run-123",
		"empty_map":  map[string]interface{}{},
		"empty_list": []interface{}{},
		"null_val":   nil,
		"kept":       "value",
	}

	result := PruneEmpty(input)
	m := result.(map[string]interface{})

	if _, exists := m["empty_map"]; exists {
		t.Error("empty map should be pruned")
	}
	if _, exists := m["empty_list"]; exists {
		t.Error("empty list should be pruned")
	}
	if _, exists := m["null_val"]; exists {
		t.Error("null should be pruned")
	}
	if m["kept"] != "value" {
		t.Error("non-empty values should be kept")
	}
}

func TestRenderJSON_PreservesRelationshipIDs(t *testing.T) {
	type detail struct {
		ID      string `json:"id"`
		PlanID  string `json:"plan_id"`
		ApplyID string `json:"apply_id"`
	}

	data := detail{ID: "run-123", PlanID: "plan-456", ApplyID: ""}

	// Marshal -> unmarshal -> prune (same as renderJSON)
	raw, _ := json.Marshal(data)
	var generic interface{}
	json.Unmarshal(raw, &generic)
	pruned := PruneEmpty(generic)

	m := pruned.(map[string]interface{})
	if m["plan_id"] != "plan-456" {
		t.Errorf("expected plan_id 'plan-456', got %v", m["plan_id"])
	}
	if _, exists := m["apply_id"]; !exists {
		t.Error("apply_id should be present even when empty string")
	}
}
