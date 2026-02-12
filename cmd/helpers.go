package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/roboalchemist/tfc/pkg/api"
	"github.com/roboalchemist/tfc/pkg/auth"
	"github.com/roboalchemist/tfc/pkg/output"
)

// newClient creates an authenticated Terraform Cloud API client.
func newClient() (*api.Client, error) {
	token, err := auth.GetToken()
	if err != nil {
		return nil, output.NewAuthError(err.Error())
	}
	baseURL := auth.GetAddress()
	client := api.NewClient(baseURL, token)
	if flagDebug {
		client.SetDebug(DebugLog)
	}
	return client, nil
}

// requireOrg returns the organization name from --org flag or TFC_ORG env var.
func requireOrg() (string, error) {
	if flagOrg == "" {
		return "", output.NewUsageError("organization required: use --org flag or set TFC_ORG env var")
	}
	return flagOrg, nil
}

// truncateStr truncates a string to max length with ellipsis.
func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// notImplemented returns a structured "not yet implemented" error for stub commands.
func notImplemented(cmdName string) error {
	return fmt.Errorf("command %q is not yet implemented", cmdName)
}

// shortDate extracts the date portion from an RFC3339 timestamp.
func shortDate(ts string) string {
	if len(ts) >= 10 {
		return ts[:10]
	}
	return ts
}

// boolStr returns "yes" or "no" for a boolean.
func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// itoa converts an int to string.
func itoa(n int) string {
	return strconv.Itoa(n)
}

// defaultStr returns the value or a fallback if empty.
func defaultStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// statusColor returns a status string (no color in table, color optional later).
func statusColor(s string) string {
	return s
}

// joinTags joins a slice of strings with commas.
func joinTags(tags []string) string {
	return strings.Join(tags, ", ")
}
