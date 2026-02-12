package cmd

import "github.com/spf13/cobra"

var agentPoolCmd = &cobra.Command{
	Use:     "agent-pool",
	Aliases: []string{"ap"},
	Short:   "View agent pools",
}

var agentPoolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List agent pools in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("agent-pool list")
	},
}

var agentPoolShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show agent pool details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("agent-pool show")
	},
}

func init() {
	agentPoolCmd.AddCommand(agentPoolListCmd, agentPoolShowCmd)
	rootCmd.AddCommand(agentPoolCmd)
}
