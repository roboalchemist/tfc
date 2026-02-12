package cmd

import "github.com/spf13/cobra"

var varsetCmd = &cobra.Command{
	Use:     "varset",
	Aliases: []string{"vs"},
	Short:   "Manage variable sets",
}

var varsetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variable sets for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset list")
	},
}

var varsetShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show variable set details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset show")
	},
}

var varsetCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new variable set",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset create")
	},
}

var varsetUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a variable set",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset update")
	},
}

var varsetDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a variable set",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset delete")
	},
}

var varsetApplyCmd = &cobra.Command{
	Use:   "apply [id]",
	Short: "Apply a variable set to workspaces",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset apply")
	},
}

var varsetRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a variable set from workspaces",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("varset remove")
	},
}

func init() {
	// Create flags
	varsetCreateCmd.Flags().String("description", "", "Description")
	varsetCreateCmd.Flags().Bool("global", false, "Apply to all workspaces in the organization")

	// Update flags
	varsetUpdateCmd.Flags().String("name", "", "New name")
	varsetUpdateCmd.Flags().String("description", "", "Description")
	varsetUpdateCmd.Flags().Bool("global", false, "Apply to all workspaces")

	// Apply/Remove flags
	varsetApplyCmd.Flags().StringSlice("workspace", nil, "Workspace IDs to apply to")
	varsetRemoveCmd.Flags().StringSlice("workspace", nil, "Workspace IDs to remove from")

	varsetCmd.AddCommand(
		varsetListCmd,
		varsetShowCmd,
		varsetCreateCmd,
		varsetUpdateCmd,
		varsetDeleteCmd,
		varsetApplyCmd,
		varsetRemoveCmd,
	)
	rootCmd.AddCommand(varsetCmd)
}
