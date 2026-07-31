package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Scaffold a new odu-compatible script repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if _, err := os.Stat(name); err == nil {
			return fmt.Errorf("directory %q already exists", name)
		}

		dirs := []string{
			name,
			filepath.Join(name, "scripts"),
		}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				return err
			}
		}

		files := map[string]string{
			filepath.Join(name, "odu.yaml"): oduYamlTemplate,
			filepath.Join(name, "scripts", "setup.sh"):  setupScript,
			filepath.Join(name, "scripts", "deploy.sh"): deployScript,
			filepath.Join(name, "README.md"):             readmeTemplate(name),
		}

		for path, content := range files {
			if err := os.WriteFile(path, []byte(content), 0755); err != nil {
				return err
			}
		}

		fmt.Printf("✓ Created %s/\n", name)
		fmt.Println()
		fmt.Printf("  %-30s %s\n", name+"/odu.yaml", "script manifest")
		fmt.Printf("  %-30s %s\n", name+"/scripts/setup.sh", "example script")
		fmt.Printf("  %-30s %s\n", name+"/scripts/deploy.sh", "example script")
		fmt.Printf("  %-30s %s\n", name+"/README.md", "documentation")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. cd %s\n", name)
		fmt.Println("  2. Edit odu.yaml and scripts/ to fit your team")
		fmt.Println("  3. Push to GitHub")
		fmt.Printf("  4. odu add %s github.com/<org>/%s\n", name, name)
		return nil
	},
}

var oduYamlTemplate = `scripts:
  setup:
    path: scripts/setup.sh
    description: Install all project dependencies
  deploy:
    path: scripts/deploy.sh
    description: Deploy to an environment (usage: deploy --env prod)
`

var setupScript = `#!/bin/bash
# Description: Install all project dependencies
set -e

echo "Installing dependencies..."
# Add your setup steps here
echo "✓ Setup complete"
`

var deployScript = `#!/bin/bash
# Description: Deploy to an environment (usage: deploy --env prod)
set -e

ENV=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --env) ENV="$2"; shift 2 ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

if [ -z "$ENV" ]; then
  echo "Error: --env is required"
  echo "Usage: odu <namespace> deploy --env <prod|staging|dev>"
  exit 1
fi

echo "Deploying to $ENV..."
# Add your deploy logic here
echo "✓ Deployed to $ENV"
`

func readmeTemplate(name string) string {
	return fmt.Sprintf(`# %s

An odu script repo. Add it as a namespace with:

`+"```"+`bash
odu add %s github.com/<org>/%s
`+"```"+`

## Available scripts

| Script | Description |
|---|---|
| `+"`setup`"+` | Install all project dependencies |
| `+"`deploy`"+` | Deploy to an environment |

## Usage

`+"```"+`bash
odu %s setup
odu %s deploy --env prod
`+"```"+`

## Structure

`+"```"+`
.
├── odu.yaml
└── scripts/
    ├── setup.sh
    └── deploy.sh
`+"```"+`
`, name, name, name, name, name)
}
