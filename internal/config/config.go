package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Namespace struct {
	URL       string `json:"url"`
	LocalPath string `json:"local_path"`
}

type Config struct {
	Namespaces map[string]Namespace `json:"namespaces"`
}

func oduDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".odu"), nil
}

func configPath() (string, error) {
	dir, err := oduDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func ReposDir() (string, error) {
	dir, err := oduDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "repos"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Namespaces: map[string]Namespace{}}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Namespaces == nil {
		cfg.Namespaces = map[string]Namespace{}
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir, err := oduDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
