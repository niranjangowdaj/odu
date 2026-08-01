package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odu",
	Short: "Namespace-based script runner",
	Long: `odu lets you register GitHub repos as namespaces and run their scripts by name.

Examples:
  odu add bpi github.com/org/bpi-scripts   # register a namespace
  odu bpi                                   # list available scripts
  odu bpi install                           # run install.sh
  odu bpi deploy --env prod                # run deploy.sh with args`,
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return Dispatch(args[0], args[1:])
	},
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(showCmd)
}
