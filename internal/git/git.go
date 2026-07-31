package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func checkGit() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH — please install git first")
	}
	return nil
}

func Clone(url, dest string) error {
	if err := checkGit(); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", url, dest)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %s", stderr.String())
	}
	return nil
}

func Pull(repoPath string) error {
	if err := checkGit(); err != nil {
		return err
	}
	cmd := exec.Command("git", "-C", repoPath, "pull", "--ff-only", "--quiet")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// pull failure is non-fatal; caller decides whether to surface the error
	return cmd.Run()
}
