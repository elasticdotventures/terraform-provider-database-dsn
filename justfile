# Terraform Provider Database DSN - Release Automation
# Usage: just <command>

set shell := ["bash", "-c"]
set dotenv-load := true

# Default recipe - show available commands
default:
    @just --list

# Development commands
alias t := test
alias b := build
alias l := lint

# Release commands  
alias rp := release-patch
alias rm := release-minor
alias rM := release-major

# Variables
go_version := `go version | awk '{print $3}' | sed 's/go//'`
git_branch := `git branch --show-current`
git_dirty := `git status --porcelain | wc -l | tr -d ' '`

# Display current status
status:
    @echo "🔍 Repository Status"
    @echo "Go version: {{go_version}}"
    @echo "Git branch: {{git_branch}}"
    @echo "Working tree: $([ {{git_dirty}} -eq 0 ] && echo 'clean ✅' || echo 'dirty ❌ ({{git_dirty}} files)')"
    @echo "Latest tag: $(git describe --tags --abbrev=0 2>/dev/null || echo 'none')"

# Run tests
test:
    @echo "🧪 Running tests..."
    go test -v ./...
    go test -race ./internal/provider/

# Build the provider
build:
    @echo "🔨 Building provider..."
    go build -o terraform-provider-database-dsn

# Lint and format code
lint:
    @echo "🧹 Linting code..."
    go fmt ./...
    go vet ./...
    go mod tidy

# Pre-release checks
_pre_release_checks:
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "🔍 Running pre-release checks..."
    
    # Check if on main/master branch
    if [[ "{{git_branch}}" != "main" && "{{git_branch}}" != "master" ]]; then
        echo "❌ Must be on main/master branch (currently on {{git_branch}})"
        exit 1
    fi
    
    # Check if working tree is clean
    if [[ {{git_dirty}} -ne 0 ]]; then
        echo "❌ Working tree is dirty ({{git_dirty}} files). Commit or stash changes first."
        git status --porcelain
        exit 1
    fi
    
    # Check if we can reach origin
    if ! git ls-remote origin >/dev/null 2>&1; then
        echo "❌ Cannot reach git origin. Check network connection."
        exit 1
    fi
    
    # Pull latest changes
    echo "📥 Pulling latest changes..."
    git pull origin "{{git_branch}}"
    
    # Run tests
    echo "🧪 Running tests..."
    if ! go test ./...; then
        echo "❌ Tests failed"
        exit 1
    fi
    
    echo "✅ Pre-release checks passed"

# Get current version from git tags
_get_current_version:
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Get the latest tag, default to v0.0.0 if none exists
    current=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo "$current"

# Bump version (patch, minor, or major)
_bump_version TYPE:
    #!/usr/bin/env bash
    set -euo pipefail
    
    bump_type="{{TYPE}}"
    current=$(just _get_current_version)
    echo "📋 Current version: $current"
    
    # Remove 'v' prefix for processing
    version_no_v=${current#v}
    
    # Split version into components
    IFS='.' read -r major minor patch <<< "$version_no_v"
    
    # Bump the appropriate component
    case "$bump_type" in
        "patch")
            patch=$((patch + 1))
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        *)
            echo "❌ Invalid bump type: $bump_type. Use patch, minor, or major."
            exit 1
            ;;
    esac
    
    new_version="v${major}.${minor}.${patch}"
    echo "🎯 New version: $new_version"
    echo "$new_version"

# Update changelog for release
_update_changelog VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    
    version="{{VERSION}}"
    echo "📝 Updating changelog for $version..."
    
    # Create changelog if it doesn't exist
    if [[ ! -f CHANGELOG.md ]]; then
        echo "# Changelog" > CHANGELOG.md
        echo "" >> CHANGELOG.md
        echo "All notable changes to this project will be documented in this file." >> CHANGELOG.md
        echo "" >> CHANGELOG.md
        echo "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)," >> CHANGELOG.md
        echo "and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)." >> CHANGELOG.md
        echo "" >> CHANGELOG.md
    fi
    
    # Get commits since last tag
    last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [[ -n "$last_tag" ]]; then
        commits=$(git log ${last_tag}..HEAD --oneline --pretty=format:"- %s" | head -20)
    else
        commits=$(git log --oneline --pretty=format:"- %s" | head -20)
    fi
    
    # Prepare changelog entry
    date=$(date '+%Y-%m-%d')
    temp_file=$(mktemp)
    
    # Insert new version at the top of changelog
    {
        head -n 6 CHANGELOG.md  # Keep header
        echo ""
        echo "## [$version] - $date"
        echo ""
        if [[ -n "$commits" ]]; then
            echo "### Changed"
            echo "$commits"
        else
            echo "### Changed"
            echo "- Initial release"
        fi
        echo ""
        tail -n +7 CHANGELOG.md  # Rest of changelog
    } > "$temp_file"
    
    mv "$temp_file" CHANGELOG.md
    echo "✅ Changelog updated"

# Create and push release tag
_create_release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    
    version="{{VERSION}}"
    branch="{{git_branch}}"
    echo "🏷️  Creating release $version..."
    
    # Stage and commit changelog if it was modified
    if git diff --name-only | grep -q CHANGELOG.md; then
        git add CHANGELOG.md
        git commit -m "chore: update changelog for $version"
    fi
    
    # Create annotated tag
    git tag -a "$version" -m "Release $version"
    
    # Push branch and tags
    git push origin "$branch"
    git push origin "$version"
    
    echo "🚀 Release $version created and pushed!"
    echo "📦 GitHub Actions will now build and publish to Terraform Registry"
    echo "🔗 Monitor the release at: https://github.com/elasticdotventures/terraform-provider-database-dsn/actions"

# Release patch version (x.y.Z+1)
release-patch: _pre_release_checks
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "🔄 Creating patch release..."
    new_version=$(just _bump_version patch)
    just _update_changelog "$new_version"
    just _create_release "$new_version"

# Release minor version (x.Y+1.0)
release-minor: _pre_release_checks
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "🔄 Creating minor release..."
    new_version=$(just _bump_version minor)
    just _update_changelog "$new_version"
    just _create_release "$new_version"

# Release major version (X+1.0.0)
release-major: _pre_release_checks
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "🔄 Creating major release..."
    echo "⚠️  MAJOR VERSION BUMP - This indicates breaking changes!"
    read -p "Are you sure you want to create a major release? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Major release cancelled"
        exit 1
    fi
    
    new_version=$(just _bump_version major)
    just _update_changelog "$new_version"
    just _create_release "$new_version"

# Show what the next version would be without creating it
next-version TYPE="patch":
    @just _bump_version {{TYPE}} | tail -1

# Preview changelog entry for next release
preview-changelog TYPE="patch":
    #!/usr/bin/env bash
    set -euo pipefail
    
    new_version=$(just _bump_version {{TYPE}} | tail -1)
    echo "📋 Changelog preview for $new_version:"
    echo ""
    
    # Get commits since last tag
    last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [[ -n "$last_tag" ]]; then
        commits=$(git log ${last_tag}..HEAD --oneline --pretty=format:"- %s" | head -20)
    else
        commits=$(git log --oneline --pretty=format:"- %s" | head -20)
    fi
    
    date=$(date '+%Y-%m-%d')
    echo "## [$new_version] - $date"
    echo ""
    echo "### Changed"
    if [[ -n "$commits" ]]; then
        echo "$commits"
    else
        echo "- No changes since last release"
    fi

# Clean up build artifacts
clean:
    @echo "🧹 Cleaning up..."
    rm -f terraform-provider-database-dsn
    go clean

# Install just if not present
install-just:
    #!/usr/bin/env bash
    if ! command -v just &> /dev/null; then
        echo "📦 Installing just..."
        curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ~/.local/bin
        echo "✅ just installed to ~/.local/bin/just"
        echo "Add ~/.local/bin to your PATH if not already present"
    else
        echo "✅ just is already installed"
    fi