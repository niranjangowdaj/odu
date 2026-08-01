package manifest

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Script struct {
	Path        string `yaml:"path"`
	Description string `yaml:"description"`
	Runner      string `yaml:"runner"`
}

type Manifest struct {
	Scripts map[string]Script `yaml:"scripts"`
}

func Load(repoPath string) (*Manifest, error) {
	manifestPath := filepath.Join(repoPath, "odu.yaml")
	data, err := os.ReadFile(manifestPath)
	if err == nil {
		var m Manifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		if m.Scripts == nil {
			m.Scripts = map[string]Script{}
		}
		return &m, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	// fallback: discover .sh files
	return discover(repoPath)
}

func discover(repoPath string) (*Manifest, error) {
	m := &Manifest{Scripts: map[string]Script{}}

	patterns := []string{
		filepath.Join(repoPath, "*.sh"),
		filepath.Join(repoPath, "*.py"),
		filepath.Join(repoPath, "scripts", "*.sh"),
		filepath.Join(repoPath, "scripts", "*.py"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			ext := filepath.Ext(match)
			name := strings.TrimSuffix(filepath.Base(match), ext)
			relPath, _ := filepath.Rel(repoPath, match)
			m.Scripts[name] = Script{
				Path:        relPath,
				Description: readDescription(match),
			}
		}
	}

	return m, nil
}

func readDescription(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "(no description)"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// skip shebang
		if strings.HasPrefix(line, "#!") {
			continue
		}
		if strings.HasPrefix(line, "# Description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# Description:"))
		}
		break
	}
	return "(no description)"
}
