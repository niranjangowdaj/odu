#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "Usage: ./release.sh <version>"
  echo "Example: ./release.sh v0.2.0"
  exit 1
fi

VERSION="$1"

if ! echo "$VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "Error: version must be in format vX.Y.Z (e.g. v0.2.0)"
  exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
  echo "Error: uncommitted changes present — commit or stash them before releasing"
  git status --short
  exit 1
fi

if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Error: tag $VERSION already exists"
  exit 1
fi

echo "Releasing $VERSION..."
git tag "$VERSION"
git push origin "$VERSION"

echo ""
echo "✓ Tag $VERSION pushed — GitHub Actions is now building the release."
echo "  Track progress: https://github.com/niranjangowdaj/odu/actions"
echo "  Release page:   https://github.com/niranjangowdaj/odu/releases/tag/$VERSION"
