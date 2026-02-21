package cmd

import (
	"fmt"

	"gitea.roboalch.com/roboalchemist/tfc/pkg/jsonapi"
	"gitea.roboalch.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var (
	flagPCRunID string
)

var policyCheckCmd = &cobra.Command{
	Use:     "policy-check",
	Aliases: []string{"pc"},
	Short:   "Manage policy checks",
}

var policyCheckListCmd = &cobra.Command{
	Use:   "list",
	Short: "List policy checks for a run",
	RunE:  runPolicyCheckList,
}

var policyCheckShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show policy check details",
	Args:  cobra.ExactArgs(1),
	RunE:  runPolicyCheckShow,
}

var policyCheckOverrideCmd = &cobra.Command{
	Use:   "override [id]",
	Short: "Override a soft-mandatory policy check",
	Args:  cobra.ExactArgs(1),
	RunE:  runPolicyCheckOverride,
}

func init() {
	policyCheckListCmd.Flags().StringVar(&flagPCRunID, "run", "", "Run ID (required)")

	policyCheckCmd.AddCommand(
		policyCheckListCmd,
		policyCheckShowCmd,
		policyCheckOverrideCmd,
	)
	rootCmd.AddCommand(policyCheckCmd)
}

type policyCheckAttrs struct {
	Status     string                 `json:"status"`
	Result     map[string]interface{} `json:"result"`
	Scope      string                 `json:"scope"`
	Actions    map[string]bool        `json:"actions"`
	Permissions map[string]bool       `json:"permissions"`
	StatusTimestamps map[string]string `json:"status-timestamps"`
}

func runPolicyCheckList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	if flagPCRunID == "" {
		return output.NewUsageError("--run is required")
	}

	path := fmt.Sprintf("/runs/%s/policy-checks", flagPCRunID)

	var doc jsonapi.Document
	if err := client.Get(path, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type pcJSON struct {
		ID    string           `json:"id"`
		Attrs policyCheckAttrs `json:"attributes"`
	}
	var jsonData []pcJSON
	td := output.TableData{
		Headers: []string{"ID", "STATUS", "SCOPE", "OVERRIDABLE"},
	}

	for _, r := range resources {
		var a policyCheckAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		overridable := ""
		if a.Actions != nil && a.Actions["is-overridable"] {
			overridable = "yes"
		}
		td.Rows = append(td.Rows, []string{
			r.ID, a.Status, a.Scope, overridable,
		})
		jsonData = append(jsonData, pcJSON{ID: r.ID, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runPolicyCheckShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	pcID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/policy-checks/"+pcID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a policyCheckAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	overridable := "no"
	if a.Actions != nil && a.Actions["is-overridable"] {
		overridable = "yes"
	}

	opts := GetOutputOptions()

	type pcDetail struct {
		ID    string           `json:"id"`
		Attrs policyCheckAttrs `json:"attributes"`
	}
	data := pcDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Status", a.Status},
			{"Scope", a.Scope},
			{"Overridable", overridable},
		},
	}

	return output.RenderTable(td, data, opts)
}

func runPolicyCheckOverride(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	pcID := args[0]
	path := fmt.Sprintf("/policy-checks/%s/actions/override", pcID)

	var doc jsonapi.Document
	if err := client.Post(path, nil, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a policyCheckAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type pcDetail struct {
		ID    string           `json:"id"`
		Attrs policyCheckAttrs `json:"attributes"`
	}
	data := pcDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Status", a.Status},
			{"Scope", a.Scope},
		},
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Policy check %s overridden (status: %s)\n", res.ID, a.Status)

	return output.RenderTable(td, data, opts)
}
