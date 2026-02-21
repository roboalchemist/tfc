package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestRunCreate_MissingWorkspace(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"run", "create"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --workspace")
	}
}

func TestRunCreate_WithWorkspaceID(t *testing.T) {
	var capturedBody map[string]interface{}
	var capturedPath string

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path

		if r.Method == "POST" && r.URL.Path == "/api/v2/runs" {
			json.NewDecoder(r.Body).Decode(&capturedBody)
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "run-abc123",
					"type": "runs",
					"attributes": map[string]interface{}{
						"status":     "pending",
						"message":    "test run",
						"is-destroy": false,
						"auto-apply": false,
						"created-at": "2025-01-01T00:00:00Z",
					},
					"relationships": map[string]interface{}{
						"plan": map[string]interface{}{
							"data": map[string]interface{}{
								"id":   "plan-xyz789",
								"type": "plans",
							},
						},
					},
				},
			})
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	cmd.SetArgs([]string{"run", "create", "--workspace", "ws-existing123", "--message", "test run", "--json"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedPath != "/api/v2/runs" {
		t.Errorf("expected POST to /api/v2/runs, got %s", capturedPath)
	}

	if capturedBody == nil {
		t.Fatal("expected request body to be captured")
	}
	data, ok := capturedBody["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data field in request body")
	}
	if data["type"] != "runs" {
		t.Errorf("expected type 'runs', got %v", data["type"])
	}

	rels, ok := data["relationships"].(map[string]interface{})
	if !ok {
		t.Fatal("expected relationships in request body")
	}
	ws, ok := rels["workspace"].(map[string]interface{})
	if !ok {
		t.Fatal("expected workspace relationship")
	}
	wsData, ok := ws["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected workspace data")
	}
	if wsData["id"] != "ws-existing123" {
		t.Errorf("expected workspace id 'ws-existing123', got %v", wsData["id"])
	}
}

func TestRunCreate_WithWorkspaceName(t *testing.T) {
	var postBody map[string]interface{}
	requestCount := 0

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/vnd.api+json")

		if r.Method == "GET" && r.URL.Path == "/api/v2/organizations/test-org/workspaces/my-workspace" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "ws-resolved456",
					"type": "workspaces",
					"attributes": map[string]interface{}{
						"name": "my-workspace",
					},
				},
			})
			return
		}

		if r.Method == "POST" && r.URL.Path == "/api/v2/runs" {
			json.NewDecoder(r.Body).Decode(&postBody)
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "run-newrun789",
					"type": "runs",
					"attributes": map[string]interface{}{
						"status":     "pending",
						"message":    "deploy update",
						"is-destroy": false,
						"auto-apply": true,
						"created-at": "2025-01-02T00:00:00Z",
					},
				},
			})
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	cmd.SetArgs([]string{"run", "create", "--workspace", "my-workspace", "--org", "test-org", "--message", "deploy update", "--auto-apply", "--json"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if requestCount < 2 {
		t.Errorf("expected at least 2 requests (workspace resolve + create), got %d", requestCount)
	}

	if postBody != nil {
		data := postBody["data"].(map[string]interface{})
		rels := data["relationships"].(map[string]interface{})
		ws := rels["workspace"].(map[string]interface{})
		wsData := ws["data"].(map[string]interface{})
		if wsData["id"] != "ws-resolved456" {
			t.Errorf("expected resolved workspace id 'ws-resolved456', got %v", wsData["id"])
		}
	}
}

func TestRunCreate_DestroyAndTargets(t *testing.T) {
	var capturedBody map[string]interface{}

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")

		if r.Method == "POST" && r.URL.Path == "/api/v2/runs" {
			json.NewDecoder(r.Body).Decode(&capturedBody)
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "run-destroy999",
					"type": "runs",
					"attributes": map[string]interface{}{
						"status":     "pending",
						"is-destroy": true,
						"auto-apply": false,
						"created-at": "2025-01-03T00:00:00Z",
					},
				},
			})
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	cmd.SetArgs([]string{"run", "create", "--workspace", "ws-target123", "--is-destroy", "--target", "aws_instance.web,aws_s3_bucket.data", "--json"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedBody == nil {
		t.Fatal("expected request body")
	}

	data := capturedBody["data"].(map[string]interface{})
	attrs := data["attributes"].(map[string]interface{})

	if attrs["is-destroy"] != true {
		t.Error("expected is-destroy to be true")
	}

	targets, ok := attrs["target-addrs"].([]interface{})
	if !ok {
		t.Fatal("expected target-addrs in attributes")
	}
	if len(targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(targets))
	}
	if targets[0] != "aws_instance.web" || targets[1] != "aws_s3_bucket.data" {
		t.Errorf("unexpected targets: %v", targets)
	}
}

func TestRunApply_Success(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method

		if r.Method == "POST" && r.URL.Path == "/api/v2/runs/run-abc123/actions/apply" {
			w.WriteHeader(202)
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	cmd.SetArgs([]string{"run", "apply", "run-abc123"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedMethod != "POST" {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
	if capturedPath != "/api/v2/runs/run-abc123/actions/apply" {
		t.Errorf("expected /api/v2/runs/run-abc123/actions/apply, got %s", capturedPath)
	}
}

func TestRunApply_WithComment(t *testing.T) {
	var capturedBody map[string]string

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/v2/runs/run-comment456/actions/apply" {
			json.NewDecoder(r.Body).Decode(&capturedBody)
			w.WriteHeader(202)
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	cmd.SetArgs([]string{"run", "apply", "run-comment456", "--comment", "Approved by team lead"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedBody == nil {
		t.Fatal("expected request body with comment")
	}
	if capturedBody["comment"] != "Approved by team lead" {
		t.Errorf("expected comment 'Approved by team lead', got %q", capturedBody["comment"])
	}
}

func TestRunApply_NoComment_NilBody(t *testing.T) {
	var bodyBytes []byte

	ts := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/v2/runs/run-nobody789/actions/apply" {
			var err error
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}
			w.WriteHeader(202)
			return
		}
		w.WriteHeader(404)
	})
	defer ts.Close()

	t.Setenv("TFC_TOKEN", "test-token")
	t.Setenv("TFC_ADDRESS", ts.URL)

	cmd := rootCmd
	// Explicitly set --comment to empty to reset any flag leakage from prior tests
	cmd.SetArgs([]string{"run", "apply", "run-nobody789", "--comment", ""})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(bodyBytes) != 0 {
		t.Errorf("expected empty body when no comment, got %q", string(bodyBytes))
	}
}

func TestResolveWorkspaceID_AlreadyID(t *testing.T) {
	id, err := resolveWorkspaceID("ws-abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "ws-abc123" {
		t.Errorf("expected 'ws-abc123', got %s", id)
	}
}
