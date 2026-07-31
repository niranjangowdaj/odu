package cmd

import (
	"fmt"

	"github.com/odu-cli/odu/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered namespaces",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Namespaces) == 0 {
			fmt.Println("No namespaces registered.")
			fmt.Println("Add one with: odu add <namespace> <github-url>")
			return nil
		}

		fmt.Printf("%-20s  %s\n", "NAMESPACE", "URL")
		fmt.Printf("%-20s  %s\n", "---------", "---")
		for ns, entry := range cfg.Namespaces {
			fmt.Printf("%-20s  %s\n", ns, entry.URL)
		}
		return nil
	},
}
