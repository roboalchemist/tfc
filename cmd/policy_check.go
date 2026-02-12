package cmd

import "github.com/spf13/cobra"

var policyCheckCmd = &cobra.Command{
	Use:     "policy-check",
	Aliases: []string{"pc"},
	Short:   "Manage policy checks",
}

var policyCheckListCmd = &cobra.Command{
	Use:   "list",
	Short: "List policy checks for a run",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy-check list")
	},
}

var policyCheckShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show policy check details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy-check show")
	},
}

var policyCheckOverrideCmd = &cobra.Command{
	Use:   "override [id]",
	Short: "Override a soft-mandatory policy check",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy-check override")
	},
}

func init() {
	// List flags
	policyCheckListCmd.Flags().String("run", "", "Run ID (required)")

	policyCheckCmd.AddCommand(
		policyCheckListCmd,
		policyCheckShowCmd,
		policyCheckOverrideCmd,
	)
	rootCmd.AddCommand(policyCheckCmd)
}
