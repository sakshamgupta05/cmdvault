# Release

## Building locally

```bash
# Debug build
cargo build

# Release build
cargo build --release

# Install to ~/.cargo/bin
cargo install --path .
```

## Publishing to crates.io

```bash
cargo login
cargo publish
```

## GitHub Releases

Pushing a version tag triggers the GitHub Actions workflow which:

1. Cross-compiles for macOS (amd64, arm64) and Linux (amd64)
2. Creates a GitHub release with the binaries attached

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Homebrew

After a GitHub release is created, update the Homebrew tap formula at
`sakshamgupta05/homebrew-tap` with the new version and SHA256 hashes.
