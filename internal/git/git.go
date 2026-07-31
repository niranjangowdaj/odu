package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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
		return friendlyCloneError(url, stderr.String())
	}
	return nil
}

// PullResult describes what happened during a pull.
type PullResult int

const (
	PullOK      PullResult = iota
	PullNoUpdate           // already up to date
	PullFailed             // network or access error
)

func Pull(repoPath string) (PullResult, error) {
	if err := checkGit(); err != nil {
		return PullFailed, err
	}
	cmd := exec.Command("git", "-C", repoPath, "pull", "--ff-only", "--quiet")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return PullFailed, friendlyPullError(stderr.String())
	}
	return PullOK, nil
}

func friendlyCloneError(url, stderr string) error {
	s := strings.ToLower(stderr)
	switch {
	case strings.Contains(s, "repository not found") || strings.Contains(s, "does not exist"):
		return fmt.Errorf("repo not found: %s\nCheck the URL and make sure the repository exists", url)
	case strings.Contains(s, "authentication failed") || strings.Contains(s, "could not read username") || strings.Contains(s, "access denied") || strings.Contains(s, "permission denied"):
		return fmt.Errorf("access denied: %s\nMake sure you have read access and your credentials are configured", url)
	case strings.Contains(s, "could not resolve host") || strings.Contains(s, "unable to access") || strings.Contains(s, "network"):
		return fmt.Errorf("network error: could not reach %s\nCheck your internet connection and try again", url)
	case strings.Contains(s, "already exists") && strings.Contains(s, "not an empty directory"):
		return fmt.Errorf("destination directory already exists and is not empty")
	default:
		return fmt.Errorf("clone failed: %s", strings.TrimSpace(stderr))
	}
}

func friendlyPullError(stderr string) error {
	s := strings.ToLower(stderr)
	switch {
	case strings.Contains(s, "authentication failed") || strings.Contains(s, "access denied") || strings.Contains(s, "permission denied"):
		return fmt.Errorf("access denied — could not pull latest scripts (using cached version)")
	case strings.Contains(s, "could not resolve host") || strings.Contains(s, "unable to access") || strings.Contains(s, "network"):
		return fmt.Errorf("no internet connection — using cached scripts")
	default:
		return fmt.Errorf("update failed — using cached scripts")
	}
}
