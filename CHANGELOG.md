## v0.6.2 — 2026-08-02

- Fix upgrade: use proper JSON parsing instead of line scanning (85e8de0)

## v0.6.1 — 2026-08-02

- Add Node.js and Ruby script examples to README (4ea37fc)

## v0.6.0 — 2026-08-02

- Ignore .DS_Store (54e81f7)
- Add Python (and Ruby, Node) script support via shebang, extension, or runner field (d512ae4)
- Fix spelling warnings in README (81d5562)
- Add comprehensive script repo management section to README (c5ef06e)

## v0.5.0 — 2026-08-01

- Add odu upgrade command for self-updating the binary (adbdea4)

## v0.4.0 — 2026-08-01

- Add unit tests for config, manifest, git, and cmd packages (24636db)

## v0.3.2 — 2026-08-01

- Add download spinner to install.sh (0a953e7)

## v0.3.1 — 2026-08-01

- Install to ~/.local/bin to avoid sudo requirement (56658a4)

## v0.3.0 — 2026-08-01

- Add trust prompt on odu add and improved error handling for clone/pull failures (e1ead22)

# Changelog

## v0.2.0 — 2026-08-01

- Update README: sample repo usage and odu init docs (362d8ba)
- Add odu init command to scaffold new script repos (be06318)
- Update README: document patch/minor/major release workflow (731b002)
- Simplify release.sh: patch/minor/major bumping with auto CHANGELOG (8a1f130)
- Add README (7ae1111)
- Add release.sh helper script (11cb373)
