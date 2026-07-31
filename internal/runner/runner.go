package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(repoPath, scriptRelPath string, args []string) error {
	scriptPath := filepath.Join(repoPath, scriptRelPath)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script not found: %s", scriptRelPath)
	}

	display := "bash " + scriptRelPath
	if len(args) > 0 {
		display += " " + strings.Join(args, " ")
	}
	fmt.Fprintf(os.Stderr, "→ Running: %s\n", display)

	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command("bash", cmdArgs...)
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
