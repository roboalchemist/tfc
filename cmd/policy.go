package cmd

import "github.com/spf13/cobra"

var policyCmd = &cobra.Command{
	Use:     "policy",
	Aliases: []string{"pol"},
	Short:   "View policies",
}

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List policies in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy list")
	},
}

var policyShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show policy details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("policy show")
	},
}

func init() {
	policyCmd.AddCommand(policyListCmd, policyShowCmd)
	rootCmd.AddCommand(policyCmd)
}
