# Release Process

This document outlines the process for releasing the `terraform-provider-database-dsn` to the Terraform Registry.

## Prerequisites

### 1. GPG Key Setup (Required for Terraform Registry)

Generate a GPG key for signing releases:

```bash
# Generate a new GPG key (if you don't have one)
gpg --full-generate-key

# List your keys to get the fingerprint
gpg --list-secret-keys --keyid-format LONG

# Export your public key for the Terraform Registry
gpg --export --armor "your-key-id" > public-key.gpg
```

### 2. GitHub Secrets Configuration

Set up the following secrets in your GitHub repository (`Settings > Secrets and variables > Actions`):

#### Required for GoReleaser:
- `GPG_PRIVATE_KEY`: Your GPG private key (export with `gpg --export-secret-keys --armor "your-key-id"`)
- `PASSPHRASE`: Your GPG key passphrase

#### Required for Terraform Registry Push:
- `TF_CLOUD_API_TOKEN`: Your Terraform Cloud API token
- `TF_CLOUD_ORGANIZATION`: Your Terraform Cloud organization name
- `TF_CLOUD_WORKSPACE`: Your Terraform Cloud workspace name

### 3. Terraform Registry Setup

1. Sign up for a [Terraform Cloud account](https://app.terraform.io/)
2. Create an organization and workspace
3. Upload your GPG public key to the Terraform Registry
4. Configure your provider namespace

## Release Steps

### 1. Verify Module Path

Ensure your Go module path matches the GitHub repository:

```go
// go.mod should contain:
module github.com/elasticdotventures/terraform-provider-database-dsn
```

### 2. Create and Push a Release Tag

Use semantic versioning for your releases:

```bash
# Tag the release
git tag v0.1.0

# Push the tag to trigger the release workflow
git push origin v0.1.0
```

### 3. Automated Release Process

When you push a tag matching `v*`, the GitHub Actions workflow will:

1. **GoReleaser Job**:
   - Build binaries for multiple platforms (Linux, macOS, Windows, FreeBSD)
   - Generate checksums and sign them with your GPG key
   - Create a GitHub release with artifacts
   - Include the `terraform-registry-manifest.json`

2. **Terraform Registry Job**:
   - Push the provider to Terraform Cloud
   - Make it available in the Terraform Registry

## Supported Platforms

The provider is built for the following platforms:
- Linux: amd64, 386, arm, arm64
- macOS: amd64, arm64
- Windows: amd64, 386, arm, arm64
- FreeBSD: amd64, 386, arm, arm64

## Versioning

Follow [Semantic Versioning](https://semver.org/):
- `v0.1.0` - Initial release
- `v0.1.1` - Patch fixes
- `v0.2.0` - Minor features/changes
- `v1.0.0` - Major stable release

## Troubleshooting

### GPG Issues
- Ensure your GPG key doesn't expire
- Verify the `GPG_FINGERPRINT` environment variable matches your key
- Check that the private key is properly formatted in the secret

### Terraform Registry Issues
- Verify your API token has the correct permissions
- Ensure your organization and workspace names are correct
- Check that your GPG public key is uploaded to the registry

### Build Issues
- Verify all tests pass before tagging
- Check Go version compatibility
- Ensure all dependencies are properly vendored