package cmd

import "github.com/spf13/cobra"

var auditTrailCmd = &cobra.Command{
	Use:     "audit-trail",
	Aliases: []string{"audit"},
	Short:   "View audit trail events",
}

var auditTrailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audit trail events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("audit-trail list")
	},
}

func init() {
	// List flags
	auditTrailListCmd.Flags().String("since", "", "Filter events after this timestamp (RFC3339)")
	auditTrailListCmd.Flags().Int("page-size", 100, "Results per page")

	auditTrailCmd.AddCommand(auditTrailListCmd)
	rootCmd.AddCommand(auditTrailCmd)
}
