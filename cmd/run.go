package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var (
	flagRunWorkspace string
	flagRunStatus    string
	flagRunPageSize  int
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage runs",
}

var runListCmd = &cobra.Command{
	Use:   "list",
	Short: "List runs for a workspace",
	RunE:  runRunList,
}

var runShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show run details",
	Args:  cobra.ExactArgs(1),
	RunE:  runRunShow,
}

var runCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new run",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run create")
	},
}

var runApplyCmd = &cobra.Command{
	Use:   "apply [id]",
	Short: "Apply a run that is paused waiting for confirmation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run apply")
	},
}

var runDiscardCmd = &cobra.Command{
	Use:   "discard [id]",
	Short: "Discard a run that is paused waiting for confirmation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run discard")
	},
}

var runCancelCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel a run that is currently planning or applying",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run cancel")
	},
}

func init() {
	runListCmd.Flags().StringVar(&flagRunWorkspace, "workspace", "", "Workspace ID (required)")
	runListCmd.Flags().StringVar(&flagRunStatus, "status", "", "Filter by status")
	runListCmd.Flags().IntVar(&flagRunPageSize, "page-size", 20, "Results per page")

	runCreateCmd.Flags().String("workspace", "", "Workspace name or ID (required)")
	runCreateCmd.Flags().String("message", "", "Run message")
	runCreateCmd.Flags().Bool("is-destroy", false, "Plan a destroy operation")
	runCreateCmd.Flags().Bool("auto-apply", false, "Auto-apply if plan succeeds")
	runCreateCmd.Flags().String("target", "", "Comma-separated resource targets")

	runApplyCmd.Flags().String("comment", "", "Comment for the apply")
	runDiscardCmd.Flags().String("comment", "", "Comment for the discard")
	runCancelCmd.Flags().Bool("force", false, "Force cancel")

	runCmd.AddCommand(
		runListCmd,
		runShowCmd,
		runCreateCmd,
		runApplyCmd,
		runDiscardCmd,
		runCancelCmd,
	)
	rootCmd.AddCommand(runCmd)
}

type runAttrs struct {
	Status           string `json:"status"`
	Source           string `json:"source"`
	Message          string `json:"message"`
	IsDestroy        bool   `json:"is-destroy"`
	HasChanges       bool   `json:"has-changes"`
	AutoApply        bool   `json:"auto-apply"`
	CreatedAt        string `json:"created-at"`
	StatusTimestamps json.RawMessage `json:"status-timestamps"`
}

func runRunList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	if flagRunWorkspace == "" {
		return output.NewUsageError("--workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/runs?page[size]=%d", flagRunWorkspace, flagRunPageSize)
	if flagRunStatus != "" {
		path += "&filter[status]=" + flagRunStatus
	}

	var doc jsonapi.Document
	if err := client.Get(path, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type runJSON struct {
		ID    string   `json:"id"`
		Attrs runAttrs `json:"attributes"`
	}
	var jsonData []runJSON
	td := output.TableData{
		Headers: []string{"ID", "STATUS", "MESSAGE", "SOURCE", "CHANGES", "CREATED"},
	}

	for _, r := range resources {
		var a runAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		changes := ""
		if a.HasChanges {
			changes = "yes"
		}
		msg := truncateStr(a.Message, 50)
		td.Rows = append(td.Rows, []string{
			r.ID, a.Status, msg, a.Source, changes, shortDate(a.CreatedAt),
		})
		jsonData = append(jsonData, runJSON{ID: r.ID, Attrs: a})
	}

	if doc.Meta != nil && doc.Meta.Pagination != nil {
		p := doc.Meta.Pagination
		fmt.Fprintf(cmd.ErrOrStderr(), "Page %d/%d (%d total)\n", p.CurrentPage, p.TotalPages, p.TotalCount)
	}

	return output.RenderTable(td, jsonData, opts)
}

func runRunShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	runID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/runs/"+runID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a runAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	// Extract plan/apply IDs from relationships
	planID := extractRelationshipID(res, "plan")
	applyID := extractRelationshipID(res, "apply")

	opts := GetOutputOptions()

	type runDetail struct {
		ID      string   `json:"id"`
		PlanID  string   `json:"plan_id,omitempty"`
		ApplyID string   `json:"apply_id,omitempty"`
		Attrs   runAttrs `json:"attributes"`
	}
	data := runDetail{ID: res.ID, PlanID: planID, ApplyID: applyID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Status", a.Status},
			{"Message", a.Message},
			{"Source", a.Source},
			{"Is Destroy", boolStr(a.IsDestroy)},
			{"Has Changes", boolStr(a.HasChanges)},
			{"Auto Apply", boolStr(a.AutoApply)},
			{"Plan ID", defaultStr(planID, "-")},
			{"Apply ID", defaultStr(applyID, "-")},
			{"Created", a.CreatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}

// extractRelationshipID pulls the "id" from a relationship's data object.
func extractRelationshipID(res *jsonapi.Resource, relName string) string {
	raw, ok := res.Relationships[relName]
	if !ok {
		return ""
	}
	var rel struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &rel) == nil {
		return rel.Data.ID
	}
	return ""
}
