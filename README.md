# odu

A namespace-based script runner. Register GitHub repos as namespaces and run their scripts by name from anywhere.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/niranjangowdaj/odu/main/install.sh | bash
```

Supports macOS (Intel + Apple Silicon) and Linux. Requires `git` and `bash`.

## Upgrade

```bash
odu upgrade
```

Checks for the latest release and replaces the binary in-place. No sudo needed.

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

Try the sample repo to see odu in action:

```bash
odu add sample https://github.com/<org>/odu-sample-scripts.git
odu sample
odu sample setup
odu sample deploy --env prod
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
| `odu init <name>` | Scaffold a new script repo |
| `odu upgrade` | Upgrade odu to the latest version |

## Setting up a script repo

Scaffold a new repo instantly with:

```bash
odu init my-scripts
cd my-scripts
# edit odu.yaml and scripts/, then push to GitHub
odu add my-scripts github.com/<org>/my-scripts
```

Or look at the sample repo for a reference on structure and `odu.yaml` format.

Manually, create a GitHub repo with an `odu.yaml` in the root:

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
./release.sh           # patch: v0.1.0 → v0.1.1 (default)
./release.sh minor     # minor: v0.1.0 → v0.2.0
./release.sh major     # major: v0.1.0 → v1.0.0
```

Each release automatically updates `CHANGELOG.md` with all commits since the last release, commits it, then tags and pushes — triggering GitHub Actions to build binaries for macOS, Linux, and Windows.
