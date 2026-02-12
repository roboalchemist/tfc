package cmd

import "github.com/spf13/cobra"

var configVersionCmd = &cobra.Command{
	Use:     "config-version",
	Aliases: []string{"cv"},
	Short:   "Manage configuration versions",
}

var configVersionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration versions for a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("config-version list")
	},
}

var configVersionShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show configuration version details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("config-version show")
	},
}

var configVersionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("config-version create")
	},
}

var configVersionUploadCmd = &cobra.Command{
	Use:   "upload [id] [file]",
	Short: "Upload a tarball to a configuration version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("config-version upload")
	},
}

func init() {
	// List flags
	configVersionListCmd.Flags().String("workspace", "", "Workspace ID (required)")
	configVersionListCmd.Flags().Int("page-size", 20, "Results per page")

	// Create flags
	configVersionCreateCmd.Flags().String("workspace", "", "Workspace ID (required)")
	configVersionCreateCmd.Flags().Bool("auto-queue-runs", true, "Auto-queue runs on upload")
	configVersionCreateCmd.Flags().Bool("speculative", false, "Speculative plan only")

	configVersionCmd.AddCommand(
		configVersionListCmd,
		configVersionShowCmd,
		configVersionCreateCmd,
		configVersionUploadCmd,
	)
	rootCmd.AddCommand(configVersionCmd)
}
