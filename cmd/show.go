package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/odu-cli/odu/internal/config"
	"github.com/odu-cli/odu/internal/manifest"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <namespace> <script>",
	Short: "Show the source of a script (press q to quit)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := args[0]
		scriptName := args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		entry, exists := cfg.Namespaces[ns]
		if !exists {
			return fmt.Errorf("unknown namespace %q — run 'odu list' to see registered namespaces", ns)
		}

		m, err := manifest.Load(entry.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to load scripts for '%s': %w", ns, err)
		}

		script, ok := m.Scripts[scriptName]
		if !ok {
			return fmt.Errorf("script '%s' not found in namespace '%s' — run 'odu %s' to list scripts", scriptName, ns, ns)
		}

		scriptPath := filepath.Join(entry.LocalPath, script.Path)

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return fmt.Errorf("script file not found: %s", script.Path)
		}

		return page(scriptPath)
	},
}

func page(path string) error {
	// try bat first (syntax highlighting), fall back to less, then more
	for _, pager := range []string{"bat", "less", "more"} {
		bin, err := exec.LookPath(pager)
		if err != nil {
			continue
		}

		var cmd *exec.Cmd
		switch pager {
		case "bat":
			cmd = exec.Command(bin, "--style=numbers,header", "--paging=always", path)
		case "less":
			cmd = exec.Command(bin, "-N", "--quit-if-one-screen", path)
		default:
			cmd = exec.Command(bin, path)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// last resort: just print it
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
