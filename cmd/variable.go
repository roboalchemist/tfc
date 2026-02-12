package cmd

import (
	"fmt"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var flagVarWorkspace string

var variableCmd = &cobra.Command{
	Use:   "var",
	Short: "Manage workspace variables",
}

var variableListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables for a workspace",
	RunE:  runVariableList,
}

var variableShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show variable details",
	Args:  cobra.ExactArgs(1),
	RunE:  runVariableShow,
}

var variableCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new variable",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("var create")
	},
}

var variableUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("var update")
	},
}

var variableDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("var delete")
	},
}

func init() {
	variableListCmd.Flags().StringVar(&flagVarWorkspace, "workspace", "", "Workspace ID (required)")

	variableCreateCmd.Flags().String("workspace", "", "Workspace ID (required)")
	variableCreateCmd.Flags().String("key", "", "Variable key (required)")
	variableCreateCmd.Flags().String("value", "", "Variable value")
	variableCreateCmd.Flags().String("description", "", "Variable description")
	variableCreateCmd.Flags().String("category", "terraform", "Category: terraform or env")
	variableCreateCmd.Flags().Bool("hcl", false, "Parse value as HCL")
	variableCreateCmd.Flags().Bool("sensitive", false, "Mark as sensitive")

	variableUpdateCmd.Flags().String("workspace", "", "Workspace ID (required)")
	variableUpdateCmd.Flags().String("key", "", "Variable key")
	variableUpdateCmd.Flags().String("value", "", "Variable value")
	variableUpdateCmd.Flags().String("description", "", "Variable description")
	variableUpdateCmd.Flags().Bool("hcl", false, "Parse value as HCL")
	variableUpdateCmd.Flags().Bool("sensitive", false, "Mark as sensitive")

	variableDeleteCmd.Flags().String("workspace", "", "Workspace ID (required)")

	variableCmd.AddCommand(
		variableListCmd,
		variableShowCmd,
		variableCreateCmd,
		variableUpdateCmd,
		variableDeleteCmd,
	)
	rootCmd.AddCommand(variableCmd)
}

type varAttrs struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Category    string `json:"category"`
	HCL         bool   `json:"hcl"`
	Sensitive   bool   `json:"sensitive"`
	CreatedAt   string `json:"created-at"`
}

func runVariableList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	if flagVarWorkspace == "" {
		return output.NewUsageError("--workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/vars", flagVarWorkspace)

	var doc jsonapi.Document
	if err := client.Get(path, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type varJSON struct {
		ID    string   `json:"id"`
		Attrs varAttrs `json:"attributes"`
	}
	var jsonData []varJSON
	td := output.TableData{
		Headers: []string{"ID", "KEY", "VALUE", "CATEGORY", "HCL", "SENSITIVE"},
	}

	for _, r := range resources {
		var a varAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		value := a.Value
		if a.Sensitive {
			value = "(sensitive)"
		}
		td.Rows = append(td.Rows, []string{
			r.ID, a.Key, truncateStr(value, 40), a.Category,
			boolStr(a.HCL), boolStr(a.Sensitive),
		})
		jsonData = append(jsonData, varJSON{ID: r.ID, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runVariableShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	// The vars endpoint needs workspace context for show via /workspaces/{id}/vars/{id}
	// But we can also use the direct /vars/{id} endpoint
	varID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/vars/"+varID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a varAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type varDetail struct {
		ID    string   `json:"id"`
		Attrs varAttrs `json:"attributes"`
	}
	data := varDetail{ID: res.ID, Attrs: a}

	value := a.Value
	if a.Sensitive {
		value = "(sensitive)"
	}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Key", a.Key},
			{"Value", value},
			{"Description", defaultStr(a.Description, "-")},
			{"Category", a.Category},
			{"HCL", boolStr(a.HCL)},
			{"Sensitive", boolStr(a.Sensitive)},
			{"Created", a.CreatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}
