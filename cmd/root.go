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
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := config.Load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// first arg: suggest namespace names with their URLs as descriptions
		if len(args) == 0 {
			var completions []string
			for ns, entry := range cfg.Namespaces {
				completions = append(completions, fmt.Sprintf("%s\t%s", ns, entry.URL))
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}

		// second arg: suggest script names if first arg is a known namespace
		ns := args[0]
		entry, exists := cfg.Namespaces[ns]
		if !exists {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		m, err := manifest.Load(entry.LocalPath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var completions []string
		for name, script := range m.Scripts {
			completions = append(completions, fmt.Sprintf("%s\t%s", name, script.Description))
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
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
	fmt.Println("  Available Scripts")
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
