package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/odu-cli/odu/internal/git"
)

// makeLocalRepo creates a minimal git repo at dir with one commit.
func makeLocalRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s", args, out)
		}
	}

	run("init", "-q")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-q", "-m", "init")
	return dir
}

func TestCloneLocal(t *testing.T) {
	src := makeLocalRepo(t)
	dest := filepath.Join(t.TempDir(), "cloned")

	if err := git.Clone(src, dest); err != nil {
		t.Fatalf("clone failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "README.md")); err != nil {
		t.Error("cloned repo missing README.md")
	}
}

func TestCloneNonExistentRepo(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "cloned")
	err := git.Clone("/nonexistent/path/repo", dest)
	if err == nil {
		t.Fatal("expected error cloning nonexistent repo")
	}
}

func TestPullLocal(t *testing.T) {
	src := makeLocalRepo(t)
	dest := filepath.Join(t.TempDir(), "cloned")

	if err := git.Clone(src, dest); err != nil {
		t.Fatalf("clone failed: %v", err)
	}
	if _, err := git.Pull(dest); err != nil {
		t.Errorf("pull failed: %v", err)
	}
}

func TestPullInvalidPath(t *testing.T) {
	_, err := git.Pull("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error pulling from nonexistent path")
	}
}
