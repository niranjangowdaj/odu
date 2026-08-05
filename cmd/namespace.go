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

// SmartDispatch handles: odu <script> [args...] — finds which namespace has it.
// If exactly one namespace has the script, runs it directly.
// If multiple namespaces have it, asks the user to disambiguate.
func SmartDispatch(scriptName string, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	type match struct {
		ns     string
		script manifest.Script
	}
	var matches []match

	for ns, entry := range cfg.Namespaces {
		m, err := manifest.Load(entry.LocalPath)
		if err != nil {
			continue
		}
		if script, ok := m.Scripts[scriptName]; ok {
			matches = append(matches, match{ns, script})
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf("unknown command %q — not a namespace or script name\nRun 'odu' to see available scripts", scriptName)
	}

	if len(matches) == 1 {
		entry := cfg.Namespaces[matches[0].ns]
		_, _ = git.Pull(entry.LocalPath)
		return runner.Run(entry.LocalPath, matches[0].script.Path, matches[0].script.Runner, args)
	}

	// ambiguous — list matches and ask user to be specific
	fmt.Printf("'%s' exists in multiple namespaces:\n\n", scriptName)
	sort.Slice(matches, func(i, j int) bool { return matches[i].ns < matches[j].ns })
	for _, m := range matches {
		fmt.Printf("  odu %s %s\n", m.ns, scriptName)
	}
	fmt.Println("\nPlease specify the namespace.")
	return nil
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
