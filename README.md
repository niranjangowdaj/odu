# odu

A namespace-based script runner. Register GitHub repos as namespaces and run their scripts by name from anywhere.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/niranjangowdaj/odu/main/install.sh | bash
```

Supports macOS (Intel + Apple Silicon) and Linux. Requires `git` and `bash`.

## Quick start

```bash
# Register a repo as a namespace
odu add bpi github.com/your-org/bpi-scripts

# List available scripts in that namespace
odu bpi

# Run a script
odu bpi install

# Pass arguments to the script
odu bpi deploy --env prod
```

## Commands

| Command | Description |
|---|---|
| `odu add <namespace> <github-url>` | Clone a repo and register it as a namespace |
| `odu remove <namespace>` | Remove a namespace |
| `odu list` | List all registered namespaces |
| `odu update [namespace]` | Pull latest scripts (all namespaces if none specified) |
| `odu <namespace>` | List available scripts in a namespace |
| `odu <namespace> <script> [args...]` | Run a script |

## Setting up a script repo

Create a GitHub repo with an `odu.yaml` in the root:

```yaml
scripts:
  install:
    path: scripts/install.sh
    description: Install all dependencies
  deploy:
    path: scripts/deploy.sh
    description: Deploy to an environment
```

No `odu.yaml`? No problem — odu automatically discovers any `.sh` files in the repo root and `scripts/` folder. Add a `# Description: ...` comment on the first non-shebang line to show a description in the help output.

```bash
#!/bin/bash
# Description: Install all dependencies
...
```

## How it works

- Repos are cloned to `~/.odu/repos/<namespace>/` on `odu add`
- On every script run, odu silently does a `git pull` to keep scripts fresh
- Scripts run with their repo root as the working directory, so they can reference other files in the repo using relative paths
- `~/.odu/config.json` stores the namespace registry

## Releasing a new version

```bash
./release.sh v0.2.0
```

GitHub Actions builds binaries for macOS, Linux, and Windows and publishes them as a GitHub Release automatically.
