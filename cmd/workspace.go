package cmd

import (
	"fmt"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var (
	flagWsSearch   string
	flagWsPageSize int
)

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspaces",
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspaces in an organization",
	RunE:  runWorkspaceList,
}

var workspaceShowCmd = &cobra.Command{
	Use:   "show [name-or-id]",
	Short: "Show workspace details",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkspaceShow,
}

var workspaceCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("workspace create")
	},
}

var workspaceUpdateCmd = &cobra.Command{
	Use:   "update [name-or-id]",
	Short: "Update workspace settings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("workspace update")
	},
}

var workspaceDeleteCmd = &cobra.Command{
	Use:   "delete [name-or-id]",
	Short: "Delete a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("workspace delete")
	},
}

var workspaceLockCmd = &cobra.Command{
	Use:   "lock [id]",
	Short: "Lock a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("workspace lock")
	},
}

var workspaceUnlockCmd = &cobra.Command{
	Use:   "unlock [id]",
	Short: "Unlock a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("workspace unlock")
	},
}

func init() {
	workspaceListCmd.Flags().StringVar(&flagWsSearch, "search", "", "Filter workspaces by name")
	workspaceListCmd.Flags().IntVar(&flagWsPageSize, "page-size", 20, "Results per page")

	workspaceCreateCmd.Flags().String("description", "", "Workspace description")
	workspaceCreateCmd.Flags().String("terraform-version", "", "Terraform version")
	workspaceCreateCmd.Flags().String("working-directory", "", "Working directory")
	workspaceCreateCmd.Flags().Bool("auto-apply", false, "Auto-apply successful plans")
	workspaceCreateCmd.Flags().String("vcs-repo", "", "VCS repository identifier")
	workspaceCreateCmd.Flags().String("project-id", "", "Project ID to associate with")

	workspaceUpdateCmd.Flags().String("description", "", "Workspace description")
	workspaceUpdateCmd.Flags().String("terraform-version", "", "Terraform version")
	workspaceUpdateCmd.Flags().String("working-directory", "", "Working directory")
	workspaceUpdateCmd.Flags().Bool("auto-apply", false, "Auto-apply successful plans")

	workspaceLockCmd.Flags().String("reason", "", "Reason for locking")

	workspaceCmd.AddCommand(
		workspaceListCmd,
		workspaceShowCmd,
		workspaceCreateCmd,
		workspaceUpdateCmd,
		workspaceDeleteCmd,
		workspaceLockCmd,
		workspaceUnlockCmd,
	)
	rootCmd.AddCommand(workspaceCmd)
}

type wsAttrs struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	TerraformVersion string   `json:"terraform-version"`
	AutoApply        bool     `json:"auto-apply"`
	WorkingDirectory string   `json:"working-directory"`
	ExecutionMode    string   `json:"execution-mode"`
	ResourceCount    int      `json:"resource-count"`
	Locked           bool     `json:"locked"`
	CreatedAt        string   `json:"created-at"`
	UpdatedAt        string   `json:"updated-at"`
	TagNames         []string `json:"tag-names"`
}

func runWorkspaceList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	org, err := requireOrg()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/organizations/%s/workspaces?page[size]=%d", org, flagWsPageSize)
	if flagWsSearch != "" {
		path += "&search[name]=" + flagWsSearch
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

	type wsJSON struct {
		ID    string  `json:"id"`
		Attrs wsAttrs `json:"attributes"`
	}
	var jsonData []wsJSON
	td := output.TableData{
		Headers: []string{"ID", "NAME", "TF VERSION", "EXEC MODE", "RESOURCES", "LOCKED", "UPDATED"},
	}

	for _, r := range resources {
		var a wsAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		locked := ""
		if a.Locked {
			locked = "LOCKED"
		}
		td.Rows = append(td.Rows, []string{
			r.ID, a.Name, a.TerraformVersion, a.ExecutionMode,
			itoa(a.ResourceCount), locked, shortDate(a.UpdatedAt),
		})
		jsonData = append(jsonData, wsJSON{ID: r.ID, Attrs: a})
	}

	if doc.Meta != nil && doc.Meta.Pagination != nil {
		p := doc.Meta.Pagination
		fmt.Fprintf(cmd.ErrOrStderr(), "Page %d/%d (%d total)\n", p.CurrentPage, p.TotalPages, p.TotalCount)
	}

	return output.RenderTable(td, jsonData, opts)
}

func runWorkspaceShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	org, err := requireOrg()
	if err != nil {
		return err
	}

	name := args[0]
	var doc jsonapi.Document
	if err := client.Get(fmt.Sprintf("/organizations/%s/workspaces/%s", org, name), &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a wsAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type wsDetail struct {
		ID    string  `json:"id"`
		Attrs wsAttrs `json:"attributes"`
	}
	data := wsDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Name", a.Name},
			{"Description", defaultStr(a.Description, "-")},
			{"Terraform Version", a.TerraformVersion},
			{"Execution Mode", a.ExecutionMode},
			{"Auto Apply", boolStr(a.AutoApply)},
			{"Working Directory", defaultStr(a.WorkingDirectory, "/")},
			{"Resource Count", itoa(a.ResourceCount)},
			{"Locked", boolStr(a.Locked)},
			{"Tags", defaultStr(joinTags(a.TagNames), "-")},
			{"Created", a.CreatedAt},
			{"Updated", a.UpdatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}
