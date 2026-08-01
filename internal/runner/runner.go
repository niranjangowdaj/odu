package runner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(repoPath, scriptRelPath, explicitRunner string, args []string) error {
	scriptPath := filepath.Join(repoPath, scriptRelPath)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script not found: %s", scriptRelPath)
	}

	interpreter := resolveRunner(scriptPath, explicitRunner)

	display := interpreter + " " + scriptRelPath
	if len(args) > 0 {
		display += " " + strings.Join(args, " ")
	}
	fmt.Fprintf(os.Stderr, "→ Running: %s\n", display)

	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(interpreter, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoPath

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}

// resolveRunner picks the interpreter in priority order:
// 1. explicit runner from odu.yaml
// 2. shebang line in the script file
// 3. file extension
// 4. fallback to bash
func resolveRunner(scriptPath, explicitRunner string) string {
	if explicitRunner != "" {
		return explicitRunner
	}

	if runner := readShebang(scriptPath); runner != "" {
		return runner
	}

	switch strings.ToLower(filepath.Ext(scriptPath)) {
	case ".py":
		return "python3"
	case ".rb":
		return "ruby"
	case ".js":
		return "node"
	case ".ts":
		return "ts-node"
	}

	return "bash"
}

func readShebang(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return ""
	}
	line := scanner.Text()
	if !strings.HasPrefix(line, "#!") {
		return ""
	}

	// #!/usr/bin/env python3  →  python3
	// #!/usr/bin/python3      →  python3
	shebang := strings.TrimPrefix(line, "#!")
	parts := strings.Fields(shebang)
	if len(parts) == 0 {
		return ""
	}

	// /usr/bin/env python3 → use the argument after env
	if strings.HasSuffix(parts[0], "env") && len(parts) > 1 {
		return parts[1]
	}

	// /usr/bin/python3 → basename
	return filepath.Base(parts[0])
}
