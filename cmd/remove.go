package cmd

import (
	"fmt"
	"os"

	"github.com/odu-cli/odu/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <namespace>",
	Aliases: []string{"rm"},
	Short:   "Remove a registered namespace",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		entry, exists := cfg.Namespaces[ns]
		if !exists {
			return fmt.Errorf("namespace %q not found — run 'odu list' to see registered namespaces", ns)
		}

		if err := os.RemoveAll(entry.LocalPath); err != nil {
			return fmt.Errorf("failed to remove repo directory: %w", err)
		}

		delete(cfg.Namespaces, ns)
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("✓ Removed namespace '%s'\n", ns)
		return nil
	},
}
