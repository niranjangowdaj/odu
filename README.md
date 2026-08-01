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

Or add the community utils namespace for common developer scripts:

```bash
odu add utils https://github.com/niranjangowdaj/common-script-utils
odu utils                        # list all available scripts
odu utils kport 3000             # kill whatever is on port 3000
odu utils lsport                 # see all listening ports
odu utils cleanbranches          # clean up merged git branches
odu utils serve                  # serve current directory on port 8000
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

### Scaffold with odu init

The fastest way to create a new script repo:

```bash
odu init my-scripts
cd my-scripts
# edit odu.yaml and scripts/, then push to GitHub
odu add my-scripts github.com/<org>/my-scripts
```

### Repo structure

```
my-scripts/
├── odu.yaml          # script manifest (optional but recommended)
└── scripts/
    ├── setup.sh
    ├── deploy.sh
    └── analyze.py
```

### odu.yaml manifest

Define your scripts with names, paths, and descriptions. Descriptions containing `:` must be quoted.

```yaml
scripts:
  setup:
    path: scripts/setup.sh
    description: Install all dependencies
  deploy:
    path: scripts/deploy.sh
    description: "Deploy to an environment (usage: deploy --env prod)"
  analyze:
    path: scripts/analyze.py
    description: Analyze data
```

### Without odu.yaml

No manifest needed — odu auto-discovers any `.sh` files in the repo root and `scripts/` folder. Add a `# Description: ...` comment after the shebang for help text:

```bash
#!/bin/bash
# Description: Install all dependencies
```

### Writing scripts

odu supports bash and Python scripts out of the box. The interpreter is picked automatically:

| Priority | How | Example |
|---|---|---|
| 1 | `runner` field in `odu.yaml` | `runner: python3` |
| 2 | Shebang line in the script | `#!/usr/bin/env python3` |
| 3 | File extension | `.py` → `python3`, `.rb` → `ruby`, `.js` → `node` |
| 4 | Fallback | `bash` |

**Bash script:**
```bash
#!/bin/bash
# Description: Install all dependencies
set -e
echo "Installing..."
```

**Python script:**
```python
#!/usr/bin/env python3
# Description: Analyze data
print("Analyzing...")
```

**Node.js script:**
```js
#!/usr/bin/env node
// Description: Check environment
const os = require('os')
console.log(`Node ${process.version} on ${os.platform()}`)
```

**Ruby script:**
```ruby
#!/usr/bin/env ruby
# Description: Show Ruby info
puts "Ruby #{RUBY_VERSION} on #{RUBY_PLATFORM}"
```

**Explicit runner in odu.yaml** (useful when no shebang):
```yaml
scripts:
  analyze:
    path: scripts/analyze.py
    description: Analyze data
    runner: python3
  report:
    path: scripts/report.js
    description: Generate report
    runner: node
```

- Scripts run with the **repo root as the working directory** so you can reference other files relatively
- Pass arguments naturally — `odu my-team deploy --env prod` forwards `--env prod` to `deploy.sh`
- Use `set -e` in bash scripts to stop on first error
- Exit with a non-zero code to signal failure — odu propagates it

### Sharing with your team

1. Push your script repo to GitHub (or any git host)
2. Share the `odu add` command with teammates:
   ```bash
   odu add my-team https://github.com/<org>/my-team-scripts
   ```
3. Scripts stay up to date automatically — odu pulls the latest on every run

### Managing multiple namespaces

```bash
odu list                  # see all registered namespaces
odu update                # pull latest for all namespaces at once
odu update my-team        # pull latest for a specific namespace
odu remove my-team        # remove a namespace
```

### Adding new scripts to an existing repo

1. Add the script file to `scripts/`
2. Add an entry to `odu.yaml`
3. Commit and push — teammates get it automatically on their next `odu <namespace> <script>` run

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
