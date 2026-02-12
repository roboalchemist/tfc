package cmd

import (
	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage organizations",
}

var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations the token has access to",
	RunE:  runOrgList,
}

var orgShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show organization details",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runOrgShow,
}

func init() {
	orgCmd.AddCommand(orgListCmd, orgShowCmd)
	rootCmd.AddCommand(orgCmd)
}

type orgAttrs struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	ExternalID           string `json:"external-id"`
	CreatedAt            string `json:"created-at"`
	PlanIdentifier       string `json:"plan-identifier"`
	CostEstimation       bool   `json:"cost-estimation-enabled"`
	ManagedResourceCount int    `json:"managed-resource-count"`
	DefaultExecMode      string `json:"default-execution-mode"`
	SSOEnabled           bool   `json:"sso-enabled"`
	TwoFactorConformant  bool   `json:"two-factor-conformant"`
}

func runOrgList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	var doc jsonapi.Document
	if err := client.Get("/organizations", &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	// Build JSON data for --json mode
	type orgJSON struct {
		ID    string   `json:"id"`
		Name  string   `json:"name"`
		Attrs orgAttrs `json:"attributes"`
	}
	var jsonData []orgJSON
	td := output.TableData{
		Headers: []string{"NAME", "EMAIL", "PLAN", "RESOURCES", "EXEC MODE", "CREATED"},
	}

	for _, r := range resources {
		var a orgAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		td.Rows = append(td.Rows, []string{
			a.Name, a.Email, a.PlanIdentifier,
			itoa(a.ManagedResourceCount), a.DefaultExecMode,
			shortDate(a.CreatedAt),
		})
		jsonData = append(jsonData, orgJSON{ID: r.ID, Name: a.Name, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runOrgShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	var orgName string
	if len(args) > 0 {
		orgName = args[0]
	} else {
		orgName, err = requireOrg()
		if err != nil {
			return err
		}
	}

	var doc jsonapi.Document
	if err := client.Get("/organizations/"+orgName, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a orgAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type orgDetail struct {
		ID    string   `json:"id"`
		Name  string   `json:"name"`
		Attrs orgAttrs `json:"attributes"`
	}
	data := orgDetail{ID: res.ID, Name: a.Name, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"Name", a.Name},
			{"External ID", a.ExternalID},
			{"Email", a.Email},
			{"Plan", a.PlanIdentifier},
			{"Managed Resources", itoa(a.ManagedResourceCount)},
			{"Execution Mode", a.DefaultExecMode},
			{"SSO Enabled", boolStr(a.SSOEnabled)},
			{"Created", a.CreatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}
