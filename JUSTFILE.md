# Justfile Usage Guide

This project uses [Just](https://just.systems/) as a command runner for development and release automation.

## Installation

If you don't have `just` installed:

```bash
just install-just
```

Or install manually:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ~/.local/bin
```

## Available Commands

### Development Commands

```bash
# Show repository status
just status

# Run tests
just test        # or: just t

# Build the provider
just build       # or: just b

# Lint and format code
just lint        # or: just l

# Clean build artifacts
just clean
```

### Release Commands

⚠️ **Prerequisites**: Releases require a clean working tree on the main/master branch.

```bash
# Release patch version (0.1.0 → 0.1.1)
just release-patch    # or: just rp

# Release minor version (0.1.0 → 0.2.0)  
just release-minor    # or: just rm

# Release major version (0.1.0 → 1.0.0)
just release-major    # or: just rM
```

### Preview Commands

```bash
# Show what the next version would be
just next-version patch
just next-version minor
just next-version major

# Preview changelog for next release
just preview-changelog patch
just preview-changelog minor
just preview-changelog major
```

## Release Process

The automated release process performs these steps:

### 1. Pre-release Checks
- ✅ Verify on main/master branch
- ✅ Check working tree is clean
- ✅ Pull latest changes
- ✅ Run tests

### 2. Version Management
- 📋 Determine next version using semantic versioning
- 📝 Update CHANGELOG.md with commit history
- 🏷️ Create annotated Git tag

### 3. Automated Deployment
- 🚀 Push tag to GitHub
- 📦 Trigger GitHub Actions workflow
- 🔄 Build multi-platform binaries with GoReleaser
- 📤 Publish to Terraform Registry

## Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **Patch** (`x.y.Z+1`): Bug fixes, documentation updates
- **Minor** (`x.Y+1.0`): New features, backwards-compatible changes
- **Major** (`X+1.0.0`): Breaking changes, API modifications

## Example Workflows

### Making a Patch Release

```bash
# Check current status
just status

# Preview what would be released
just preview-changelog patch

# If everything looks good, release
just release-patch
```

### Development Cycle

```bash
# Before working
just status
just test

# During development
just lint
just build
just test

# When ready to release
just preview-changelog minor
just release-minor
```

## Troubleshooting

### "Working tree is dirty" Error
```bash
# Check what files are modified
git status

# Either commit changes or stash them
git add . && git commit -m "feat: your changes"
# OR
git stash
```

### "Must be on main/master branch" Error
```bash
# Switch to main branch
git checkout main
git pull origin main
```

### "Tests failed" Error
```bash
# Run tests to see failures
just test

# Fix issues, then try release again
just release-patch
```

## GitHub Secrets Required

For the automated release to work, ensure these secrets are configured in your GitHub repository:

- `GPG_PRIVATE_KEY`: Your GPG private key for signing
- `PASSPHRASE`: Your GPG key passphrase
- `TF_CLOUD_API_TOKEN`: Terraform Cloud API token
- `TF_CLOUD_ORGANIZATION`: Your Terraform Cloud organization
- `TF_CLOUD_WORKSPACE`: Your Terraform Cloud workspace

See `RELEASE.md` for detailed setup instructions.