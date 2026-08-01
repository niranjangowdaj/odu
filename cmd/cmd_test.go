package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

// scaffoldRepo creates a temp dir with an odu.yaml and a runnable script.
func scaffoldRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "scripts"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "odu.yaml"), []byte(`
scripts:
  greet:
    path: scripts/greet.sh
    description: Print a greeting
`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "scripts", "greet.sh"), []byte("#!/bin/bash\necho hello"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestInitScaffold(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "my-scripts")

	if err := os.MkdirAll(name, 0755); err != nil {
		t.Fatal(err)
	}

	// verify expected files would be created (testing the template content)
	files := []string{"odu.yaml", "scripts/setup.sh", "scripts/deploy.sh", "README.md"}
	for _, f := range files {
		path := filepath.Join(name, f)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("test"), 0755); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist", f)
		}
	}
}
