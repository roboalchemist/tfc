package cmd

import (
	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj"},
	Short:   "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects in an organization",
	RunE:  runProjectList,
}

var projectShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show project details",
	Args:  cobra.ExactArgs(1),
	RunE:  runProjectShow,
}

var projectCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("project create")
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("project update")
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("project delete")
	},
}

func init() {
	projectCreateCmd.Flags().String("description", "", "Project description")

	projectUpdateCmd.Flags().String("name", "", "New name")
	projectUpdateCmd.Flags().String("description", "", "Description")

	projectCmd.AddCommand(
		projectListCmd,
		projectShowCmd,
		projectCreateCmd,
		projectUpdateCmd,
		projectDeleteCmd,
	)
	rootCmd.AddCommand(projectCmd)
}

type projectAttrs struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created-at"`
}

func runProjectList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	org, err := requireOrg()
	if err != nil {
		return err
	}

	var doc jsonapi.Document
	if err := client.Get("/organizations/"+org+"/projects", &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type projJSON struct {
		ID    string       `json:"id"`
		Attrs projectAttrs `json:"attributes"`
	}
	var jsonData []projJSON
	td := output.TableData{
		Headers: []string{"ID", "NAME", "DESCRIPTION", "CREATED"},
	}

	for _, r := range resources {
		var a projectAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		td.Rows = append(td.Rows, []string{
			r.ID, a.Name, truncateStr(defaultStr(a.Description, "-"), 50), shortDate(a.CreatedAt),
		})
		jsonData = append(jsonData, projJSON{ID: r.ID, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runProjectShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	projectID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/projects/"+projectID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a projectAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type projDetail struct {
		ID    string       `json:"id"`
		Attrs projectAttrs `json:"attributes"`
	}
	data := projDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Name", a.Name},
			{"Description", defaultStr(a.Description, "-")},
			{"Created", a.CreatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}
