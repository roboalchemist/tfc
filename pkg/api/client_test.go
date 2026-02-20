package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPermissionHint(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/runs/run-abc123/actions/apply", "Manage Runs"},
		{"/runs/run-abc123/actions/discard", "Manage Runs"},
		{"/runs/run-abc123/actions/cancel", "Manage Runs"},
		{"/runs", "Manage Runs"},
		{"/workspaces/ws-abc123", "Manage Workspaces"},
		{"/workspaces", "Manage Workspaces"},
		{"/policy-checks/polchk-abc/actions/override", "Override Soft-Mandatory Policies"},
		{"/state-versions", "Manage State Versions"},
		{"/state-versions/sv-abc123", "Manage State Versions"},
		{"/variables/var-abc123", "Manage Variables"},
		{"/variables", "Manage Variables"},
		{"/some/unknown/path", "Ensure your team has the required"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := permissionHint(tt.path)
			if !strings.Contains(got, tt.want) {
				t.Errorf("permissionHint(%q) = %q, want substring %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestDo_401_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errors":[{"status":"401","title":"unauthorized"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "bad-token")
	// strip the /api/v2 that NewClient appends, since our test server is flat
	client.baseURL = srv.URL

	err := client.Get("/runs", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	msg := err.Error()
	if !strings.Contains(msg, "Authentication failed (401)") {
		t.Errorf("expected '401' mention, got: %s", msg)
	}
	if !strings.Contains(msg, "TFC_TOKEN") {
		t.Errorf("expected TFC_TOKEN hint, got: %s", msg)
	}
}

func TestDo_403_Forbidden_Runs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"status":"403","title":"forbidden","detail":"insufficient permissions"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Post("/runs/run-abc/actions/apply", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	msg := err.Error()
	if !strings.Contains(msg, "Permission denied (403)") {
		t.Errorf("expected '403' mention, got: %s", msg)
	}
	if !strings.Contains(msg, "Manage Runs") {
		t.Errorf("expected 'Manage Runs' hint, got: %s", msg)
	}
}

func TestDo_403_Forbidden_PolicyOverride(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Post("/policy-checks/polchk-123/actions/override", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	msg := err.Error()
	if !strings.Contains(msg, "Override Soft-Mandatory Policies") {
		t.Errorf("expected policy override hint, got: %s", msg)
	}
}

func TestDo_403_Forbidden_StateVersions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Get("/state-versions", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Manage State Versions") {
		t.Errorf("expected state versions hint, got: %s", err.Error())
	}
}

func TestDo_403_Forbidden_Variables(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Get("/variables", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Manage Variables") {
		t.Errorf("expected variables hint, got: %s", err.Error())
	}
}

func TestDo_403_Forbidden_Workspaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Patch("/workspaces/ws-abc", []byte(`{}`), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Manage Workspaces") {
		t.Errorf("expected workspaces hint, got: %s", err.Error())
	}
}

func TestDo_403_Forbidden_GenericFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Get("/some/other/endpoint", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Ensure your team has the required") {
		t.Errorf("expected generic fallback hint, got: %s", err.Error())
	}
}

func TestGetRaw_401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "bad-token")
	client.baseURL = srv.URL

	_, err := client.GetRaw("/state-versions/sv-123/download")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "TFC_TOKEN") {
		t.Errorf("expected TFC_TOKEN hint in GetRaw 401 error, got: %s", err.Error())
	}
}

func TestGetRaw_403(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	_, err := client.GetRaw("/state-versions/sv-123/download")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Manage State Versions") {
		t.Errorf("expected state versions hint in GetRaw 403 error, got: %s", err.Error())
	}
}

func TestDo_OtherErrors_Unchanged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"status":"404","title":"not found","detail":"workspace not found"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "some-token")
	client.baseURL = srv.URL

	err := client.Get("/workspaces/ws-notexist", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "API error (status 404)") {
		t.Errorf("expected standard API error format for 404, got: %s", msg)
	}
	// Should NOT contain permission hints
	if strings.Contains(msg, "Hint:") {
		t.Errorf("non-auth errors should not have hints, got: %s", msg)
	}
}
