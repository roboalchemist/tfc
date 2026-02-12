package cmd

import "github.com/spf13/cobra"

var policySetCmd = &cobra.Command{
	Use:     "policy-set",
	Aliases: []string{"ps"},
	Short:   "View policy sets",
}

var policySetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List policy sets in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy-set list")
	},
}

var policySetShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show policy set details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy-set show")
	},
}

func init() {
	policySetCmd.AddCommand(policySetListCmd, policySetShowCmd)
	rootCmd.AddCommand(policySetCmd)
}
