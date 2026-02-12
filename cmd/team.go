package cmd

import (
	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams",
}

var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List teams in an organization",
	RunE:  runTeamList,
}

var teamShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show team details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamShow,
}

var teamCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team create")
	},
}

var teamUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team update")
	},
}

var teamDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team delete")
	},
}

func init() {
	teamCreateCmd.Flags().String("visibility", "secret", "Visibility: secret or organization")
	teamCreateCmd.Flags().Bool("manage-workspaces", false, "Can manage workspaces")
	teamCreateCmd.Flags().Bool("manage-modules", false, "Can manage modules")
	teamCreateCmd.Flags().Bool("manage-providers", false, "Can manage providers")
	teamCreateCmd.Flags().Bool("manage-policies", false, "Can manage policies")

	teamUpdateCmd.Flags().String("name", "", "New name")
	teamUpdateCmd.Flags().String("visibility", "", "Visibility: secret or organization")

	teamCmd.AddCommand(
		teamListCmd,
		teamShowCmd,
		teamCreateCmd,
		teamUpdateCmd,
		teamDeleteCmd,
	)
	rootCmd.AddCommand(teamCmd)
}

type teamAttrs struct {
	Name               string `json:"name"`
	Visibility         string `json:"visibility"`
	UsersCount         int    `json:"users-count"`
	OrganizationAccess struct {
		ManageWorkspaces bool `json:"manage-workspaces"`
		ManageModules    bool `json:"manage-modules"`
		ManageProviders  bool `json:"manage-providers"`
		ManagePolicies   bool `json:"manage-policies"`
	} `json:"organization-access"`
}

func runTeamList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	org, err := requireOrg()
	if err != nil {
		return err
	}

	var doc jsonapi.Document
	if err := client.Get("/organizations/"+org+"/teams", &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type teamJSON struct {
		ID    string    `json:"id"`
		Attrs teamAttrs `json:"attributes"`
	}
	var jsonData []teamJSON
	td := output.TableData{
		Headers: []string{"ID", "NAME", "VISIBILITY", "USERS"},
	}

	for _, r := range resources {
		var a teamAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		td.Rows = append(td.Rows, []string{
			r.ID, a.Name, a.Visibility, itoa(a.UsersCount),
		})
		jsonData = append(jsonData, teamJSON{ID: r.ID, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runTeamShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	teamID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/teams/"+teamID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a teamAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type teamDetail struct {
		ID    string    `json:"id"`
		Attrs teamAttrs `json:"attributes"`
	}
	data := teamDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Name", a.Name},
			{"Visibility", a.Visibility},
			{"Users", itoa(a.UsersCount)},
			{"Manage Workspaces", boolStr(a.OrganizationAccess.ManageWorkspaces)},
			{"Manage Modules", boolStr(a.OrganizationAccess.ManageModules)},
			{"Manage Providers", boolStr(a.OrganizationAccess.ManageProviders)},
			{"Manage Policies", boolStr(a.OrganizationAccess.ManagePolicies)},
		},
	}

	return output.RenderTable(td, data, opts)
}
