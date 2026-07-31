#!/bin/bash
set -e

BUMP="${1:-patch}"

if [[ "$BUMP" != "patch" && "$BUMP" != "minor" && "$BUMP" != "major" ]]; then
  echo "Usage: ./release.sh [patch|minor|major]"
  echo "  patch  — v0.1.0 → v0.1.1 (default)"
  echo "  minor  — v0.1.0 → v0.2.0"
  echo "  major  — v0.1.0 → v1.0.0"
  exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
  echo "Error: uncommitted changes present — commit or stash them before releasing"
  git status --short
  exit 1
fi

# Determine current version
LATEST=$(git tag --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -1)
if [ -z "$LATEST" ]; then
  LATEST="v0.0.0"
fi

MAJOR=$(echo "$LATEST" | sed 's/v\([0-9]*\)\..*/\1/')
MINOR=$(echo "$LATEST" | sed 's/v[0-9]*\.\([0-9]*\)\..*/\1/')
PATCH=$(echo "$LATEST" | sed 's/v[0-9]*\.[0-9]*\.\([0-9]*\)/\1/')

case "$BUMP" in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
esac

VERSION="v${MAJOR}.${MINOR}.${PATCH}"
DATE=$(date +"%Y-%m-%d")

echo "Bumping $LATEST → $VERSION ($BUMP)"

# Collect commits since last tag
if [ "$LATEST" = "v0.0.0" ]; then
  COMMITS=$(git log --pretty=format:"- %s (%h)" 2>/dev/null)
else
  COMMITS=$(git log "${LATEST}..HEAD" --pretty=format:"- %s (%h)" 2>/dev/null)
fi

if [ -z "$COMMITS" ]; then
  echo "Error: no commits since $LATEST — nothing to release"
  exit 1
fi

# Prepend new entry to CHANGELOG.md
CHANGELOG="CHANGELOG.md"
ENTRY="## $VERSION — $DATE

$COMMITS"

if [ -f "$CHANGELOG" ]; then
  EXISTING=$(cat "$CHANGELOG")
  printf '%s\n\n%s\n' "$ENTRY" "$EXISTING" > "$CHANGELOG"
else
  printf '# Changelog\n\n%s\n' "$ENTRY" > "$CHANGELOG"
fi

git add "$CHANGELOG"
git commit -m "Release $VERSION"
git push

git tag "$VERSION"
git push origin "$VERSION"

echo ""
echo "✓ $VERSION released — changelog updated, binaries building."
echo "  Track progress: https://github.com/niranjangowdaj/odu/actions"
echo "  Release page:   https://github.com/niranjangowdaj/odu/releases/tag/$VERSION"
