package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/odu-cli/odu/internal/config"
	"github.com/odu-cli/odu/internal/git"
	"github.com/spf13/cobra"
)

var validNamespace = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var addCmd = &cobra.Command{
	Use:   "add <namespace> <github-url>",
	Short: "Register a GitHub repo under a namespace",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := args[0]
		url := args[1]

		if !validNamespace.MatchString(ns) {
			return fmt.Errorf("invalid namespace %q — use only letters, numbers, hyphens, underscores", ns)
		}

		if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") &&
			!strings.HasPrefix(url, "git@") && !strings.HasPrefix(url, "/") {
			url = "https://" + url
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if _, exists := cfg.Namespaces[ns]; exists {
			return fmt.Errorf("namespace %q already exists — run 'odu remove %s' first", ns, ns)
		}

		reposDir, err := config.ReposDir()
		if err != nil {
			return err
		}
		localPath := filepath.Join(reposDir, ns)

		if _, err := os.Stat(localPath); err == nil {
			return fmt.Errorf("directory already exists at %s — remove it or choose a different namespace", localPath)
		}

		fmt.Printf("Cloning %s...\n", url)
		if err := git.Clone(url, localPath); err != nil {
			return err
		}

		cfg.Namespaces[ns] = config.Namespace{URL: url, LocalPath: localPath}
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("✓ Added namespace '%s' → %s\n", ns, url)
		return nil
	},
}
