package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/odu-cli/odu/internal/config"
	"github.com/odu-cli/odu/internal/git"
	"github.com/odu-cli/odu/internal/manifest"
	"github.com/odu-cli/odu/internal/runner"
)

// Dispatch handles: odu <namespace> [script] [args...]
func Dispatch(ns string, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	entry, exists := cfg.Namespaces[ns]
	if !exists {
		return fmt.Errorf("unknown namespace %q — run 'odu list' to see registered namespaces", ns)
	}

	// silently pull; warn on failure but continue with cached version
	if _, err := git.Pull(entry.LocalPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s\n\n", err)
	}

	m, err := manifest.Load(entry.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to load scripts for '%s': %w", ns, err)
	}

	// no script name — show help
	if len(args) == 0 || args[0] == "help" {
		return printHelp(ns, m)
	}

	scriptName := args[0]
	scriptArgs := args[1:]

	script, ok := m.Scripts[scriptName]
	if !ok {
		fmt.Printf("Script '%s' not found in namespace '%s'.\n\n", scriptName, ns)
		return printHelp(ns, m)
	}

	return runner.Run(entry.LocalPath, script.Path, script.Runner, scriptArgs)
}

func printHelp(ns string, m *manifest.Manifest) error {
	if len(m.Scripts) == 0 {
		fmt.Printf("No scripts found in namespace '%s'.\n", ns)
		fmt.Println("Make sure the repo has an odu.yaml or .sh files in the root/scripts/ folder.")
		return nil
	}

	fmt.Printf("Available scripts in '%s':\n\n", ns)

	names := make([]string, 0, len(m.Scripts))
	for name := range m.Scripts {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		script := m.Scripts[name]
		fmt.Printf("  %-20s  %s\n", name, script.Description)
	}
	fmt.Printf("\nUsage: odu %s <script> [args...]\n", ns)
	return nil
}
