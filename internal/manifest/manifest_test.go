package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/odu-cli/odu/internal/manifest"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "odu.yaml"), `
scripts:
  install:
    path: scripts/install.sh
    description: Install dependencies
  deploy:
    path: scripts/deploy.sh
    description: Deploy to environment
`)

	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Scripts) != 2 {
		t.Fatalf("expected 2 scripts, got %d", len(m.Scripts))
	}
	if m.Scripts["install"].Description != "Install dependencies" {
		t.Errorf("unexpected description: %s", m.Scripts["install"].Description)
	}
	if m.Scripts["deploy"].Path != "scripts/deploy.sh" {
		t.Errorf("unexpected path: %s", m.Scripts["deploy"].Path)
	}
}

func TestFallbackDiscoveryRoot(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "setup.sh"), "#!/bin/bash\necho hi")
	writeFile(t, filepath.Join(dir, "clean.sh"), "#!/bin/bash\necho clean")

	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.Scripts["setup"]; !ok {
		t.Error("expected 'setup' script from setup.sh")
	}
	if _, ok := m.Scripts["clean"]; !ok {
		t.Error("expected 'clean' script from clean.sh")
	}
}

func TestFallbackDiscoveryScriptsDir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "scripts", "deploy.sh"), "#!/bin/bash\necho deploy")

	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.Scripts["deploy"]; !ok {
		t.Error("expected 'deploy' script from scripts/deploy.sh")
	}
}

func TestFallbackDescriptionComment(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "setup.sh"), "#!/bin/bash\n# Description: My setup script\necho hi")

	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Scripts["setup"].Description != "My setup script" {
		t.Errorf("unexpected description: %q", m.Scripts["setup"].Description)
	}
}

func TestEmptyRepo(t *testing.T) {
	dir := t.TempDir()
	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Scripts) != 0 {
		t.Errorf("expected no scripts in empty repo, got %d", len(m.Scripts))
	}
}

func TestManifestTakesPriorityOverDiscovery(t *testing.T) {
	dir := t.TempDir()
	// both odu.yaml and .sh files exist — manifest should win
	writeFile(t, filepath.Join(dir, "odu.yaml"), `
scripts:
  install:
    path: scripts/install.sh
    description: From manifest
`)
	writeFile(t, filepath.Join(dir, "other.sh"), "#!/bin/bash\necho other")

	m, err := manifest.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.Scripts["other"]; ok {
		t.Error("expected odu.yaml to take priority — 'other' should not appear")
	}
	if _, ok := m.Scripts["install"]; !ok {
		t.Error("expected 'install' from odu.yaml manifest")
	}
}
