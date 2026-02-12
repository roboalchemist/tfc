package cmd

import "github.com/spf13/cobra"

var notificationCmd = &cobra.Command{
	Use:     "notification",
	Aliases: []string{"notif"},
	Short:   "Manage notification configurations",
}

var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notification configurations for a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("notification list")
	},
}

var notificationShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show notification configuration details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("notification show")
	},
}

var notificationCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a notification configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("notification create")
	},
}

var notificationUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a notification configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("notification update")
	},
}

var notificationDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a notification configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("notification delete")
	},
}

func init() {
	// List flags
	notificationListCmd.Flags().String("workspace", "", "Workspace ID (required)")

	// Create flags
	notificationCreateCmd.Flags().String("workspace", "", "Workspace ID (required)")
	notificationCreateCmd.Flags().String("destination-type", "", "Type: generic, slack, email, microsoft-teams")
	notificationCreateCmd.Flags().String("url", "", "Webhook URL (for generic/slack/microsoft-teams)")
	notificationCreateCmd.Flags().StringSlice("triggers", nil, "Trigger events")
	notificationCreateCmd.Flags().Bool("enabled", true, "Enable notification")

	// Update flags
	notificationUpdateCmd.Flags().String("name", "", "New name")
	notificationUpdateCmd.Flags().String("url", "", "Webhook URL")
	notificationUpdateCmd.Flags().StringSlice("triggers", nil, "Trigger events")
	notificationUpdateCmd.Flags().Bool("enabled", true, "Enable notification")

	notificationCmd.AddCommand(
		notificationListCmd,
		notificationShowCmd,
		notificationCreateCmd,
		notificationUpdateCmd,
		notificationDeleteCmd,
	)
	rootCmd.AddCommand(notificationCmd)
}
