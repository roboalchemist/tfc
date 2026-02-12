package cmd

import "github.com/spf13/cobra"

var teamAccessCmd = &cobra.Command{
	Use:     "team-access",
	Aliases: []string{"ta"},
	Short:   "Manage team access to workspaces",
}

var teamAccessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List team access for a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team-access list")
	},
}

var teamAccessShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show team access details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team-access show")
	},
}

var teamAccessAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add team access to a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team-access add")
	},
}

var teamAccessUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update team access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team-access update")
	},
}

var teamAccessRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove team access from a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("team-access remove")
	},
}

func init() {
	// List flags
	teamAccessListCmd.Flags().String("workspace", "", "Workspace ID (required)")

	// Add flags
	teamAccessAddCmd.Flags().String("workspace", "", "Workspace ID (required)")
	teamAccessAddCmd.Flags().String("team", "", "Team ID (required)")
	teamAccessAddCmd.Flags().String("access", "read", "Access level: read, plan, write, admin, custom")

	// Update flags
	teamAccessUpdateCmd.Flags().String("access", "", "Access level: read, plan, write, admin, custom")

	teamAccessCmd.AddCommand(
		teamAccessListCmd,
		teamAccessShowCmd,
		teamAccessAddCmd,
		teamAccessUpdateCmd,
		teamAccessRemoveCmd,
	)
	rootCmd.AddCommand(teamAccessCmd)
}
