package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/odu-cli/odu/internal/config"
	"github.com/odu-cli/odu/internal/manifest"
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
			cmd.Help()
			printNamespaceSummary()
			return nil
		}
		// check if first arg is a known namespace
		cfg, err := config.Load()
		if err == nil {
			if _, exists := cfg.Namespaces[args[0]]; exists {
				return Dispatch(args[0], args[1:])
			}
		}
		// try smart dispatch: find which namespace has this script
		return SmartDispatch(args[0], args[1:])
	},
}

// printNamespaceSummary prints all namespaces and their scripts below the help output.
func printNamespaceSummary() {
	cfg, err := config.Load()
	if err != nil || len(cfg.Namespaces) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("─────────────────────────────────────────")

	nsList := make([]string, 0, len(cfg.Namespaces))
	for ns := range cfg.Namespaces {
		nsList = append(nsList, ns)
	}
	sort.Strings(nsList)

	for _, ns := range nsList {
		entry := cfg.Namespaces[ns]
		m, err := manifest.Load(entry.LocalPath)
		if err != nil || len(m.Scripts) == 0 {
			fmt.Printf("  %-16s  (no scripts found)\n", ns)
			continue
		}

		names := make([]string, 0, len(m.Scripts))
		for name := range m.Scripts {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Printf("  %-16s  %s\n", ns, strings.Join(names, ", "))
	}
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
