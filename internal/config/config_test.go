package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/odu-cli/odu/internal/config"
)

func TestLoadEmpty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cfg.Namespaces) != 0 {
		t.Fatalf("expected empty namespaces, got %v", cfg.Namespaces)
	}
}

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, _ := config.Load()
	cfg.Namespaces["bpi"] = config.Namespace{
		URL:       "https://github.com/org/bpi",
		LocalPath: "/tmp/bpi",
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := config.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	ns, ok := loaded.Namespaces["bpi"]
	if !ok {
		t.Fatal("expected namespace 'bpi' to exist after reload")
	}
	if ns.URL != "https://github.com/org/bpi" {
		t.Errorf("unexpected URL: %s", ns.URL)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, _ := config.Load()
	cfg.Namespaces["x"] = config.Namespace{URL: "https://example.com", LocalPath: "/tmp/x"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, ".odu", "config.json")); err != nil {
		t.Errorf("config.json not created: %v", err)
	}
}

func TestDeleteNamespace(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, _ := config.Load()
	cfg.Namespaces["bpi"] = config.Namespace{URL: "https://github.com/org/bpi", LocalPath: "/tmp/bpi"}
	cfg.Save()

	delete(cfg.Namespaces, "bpi")
	cfg.Save()

	loaded, _ := config.Load()
	if _, ok := loaded.Namespaces["bpi"]; ok {
		t.Fatal("expected namespace 'bpi' to be deleted")
	}
}
