package cmd

import (
	"fmt"

	"github.com/odu-cli/odu/internal/config"
	"github.com/odu-cli/odu/internal/git"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [namespace]",
	Short: "Pull latest scripts for one or all namespaces",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(args) == 1 {
			return pullOne(cfg, args[0])
		}

		if len(cfg.Namespaces) == 0 {
			fmt.Println("No namespaces registered.")
			return nil
		}
		for ns := range cfg.Namespaces {
			if err := pullOne(cfg, ns); err != nil {
				fmt.Printf("✗ %s: %v\n", ns, err)
			}
		}
		return nil
	},
}

func pullOne(cfg *config.Config, ns string) error {
	entry, exists := cfg.Namespaces[ns]
	if !exists {
		return fmt.Errorf("namespace %q not found — run 'odu list' to see registered namespaces", ns)
	}
	fmt.Printf("Updating '%s'... ", ns)
	if _, err := git.Pull(entry.LocalPath); err != nil {
		fmt.Printf("✗ %s\n", err)
		return nil
	}
	fmt.Println("✓")
	return nil
}
