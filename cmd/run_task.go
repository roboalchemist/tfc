package cmd

import "github.com/spf13/cobra"

var runTaskCmd = &cobra.Command{
	Use:     "run-task",
	Aliases: []string{"rt"},
	Short:   "Manage run tasks",
}

var runTaskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List run tasks in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run-task list")
	},
}

var runTaskShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show run task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run-task show")
	},
}

var runTaskCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new run task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run-task create")
	},
}

var runTaskUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a run task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run-task update")
	},
}

var runTaskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a run task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("run-task delete")
	},
}

func init() {
	// Create flags
	runTaskCreateCmd.Flags().String("url", "", "Callback URL (required)")
	runTaskCreateCmd.Flags().String("description", "", "Description")
	runTaskCreateCmd.Flags().String("hmac-key", "", "HMAC key for verification")
	runTaskCreateCmd.Flags().Bool("enabled", true, "Enable the run task")

	// Update flags
	runTaskUpdateCmd.Flags().String("name", "", "New name")
	runTaskUpdateCmd.Flags().String("url", "", "Callback URL")
	runTaskUpdateCmd.Flags().String("description", "", "Description")
	runTaskUpdateCmd.Flags().Bool("enabled", true, "Enable the run task")

	runTaskCmd.AddCommand(
		runTaskListCmd,
		runTaskShowCmd,
		runTaskCreateCmd,
		runTaskUpdateCmd,
		runTaskDeleteCmd,
	)
	rootCmd.AddCommand(runTaskCmd)
}
