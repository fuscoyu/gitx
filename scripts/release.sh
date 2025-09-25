#!/bin/bash

# GitX Release Script
# This script helps create and manage releases

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_NAME="gitx"
MAIN_BRANCH="main"

# Get current version from git
get_current_version() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Get next version
get_next_version() {
    local current_version="$1"
    local version_type="$2"
    
    # Remove 'v' prefix
    local version=$(echo "$current_version" | sed 's/^v//')
    
    # Split version into parts
    local major=$(echo "$version" | cut -d. -f1)
    local minor=$(echo "$version" | cut -d. -f2)
    local patch=$(echo "$version" | cut -d. -f3)
    
    case "$version_type" in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}Error: Invalid version type. Use major, minor, or patch${NC}"
            exit 1
            ;;
    esac
    
    echo "v$major.$minor.$patch"
}

# Validate version format
validate_version() {
    local version="$1"
    if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        echo -e "${RED}Error: Invalid version format. Use semantic versioning (e.g., v1.0.0)${NC}"
        exit 1
    fi
}

# Check if working directory is clean
check_clean_working_dir() {
    if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}Error: Working directory is not clean. Please commit or stash changes.${NC}"
        exit 1
    fi
}

# Check if on main branch
check_main_branch() {
    local current_branch
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "$MAIN_BRANCH" ]; then
        echo -e "${RED}Error: Not on $MAIN_BRANCH branch. Current branch: $current_branch${NC}"
        exit 1
    fi
}

# Run tests
run_tests() {
    echo -e "${BLUE}Running tests...${NC}"
    go test -v ./...
    echo -e "${GREEN}✅ Tests passed${NC}"
}

# Build binaries
build_binaries() {
    echo -e "${BLUE}Building binaries...${NC}"
    make bin
    echo -e "${GREEN}✅ Binaries built${NC}"
}

# Create git tag
create_tag() {
    local version="$1"
    local message="$2"
    
    echo -e "${BLUE}Creating tag $version...${NC}"
    git tag -a "$version" -m "$message"
    echo -e "${GREEN}✅ Tag $version created${NC}"
}

# Push to remote
push_to_remote() {
    local version="$1"
    
    echo -e "${BLUE}Pushing to remote...${NC}"
    git push origin "$MAIN_BRANCH"
    git push origin "$version"
    echo -e "${GREEN}✅ Pushed to remote${NC}"
}

# Create release
create_release() {
    local version="$1"
    local message="$2"
    
    echo -e "${BLUE}Creating release $version...${NC}"
    echo -e "${YELLOW}Release will be created automatically by GitHub Actions${NC}"
    echo -e "${BLUE}You can also create it manually at: https://github.com/goeoeo/$REPO_NAME/releases/new${NC}"
}

# Show help
show_help() {
    cat << EOF
GitX Release Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
  patch     Create a patch release (0.0.1 -> 0.0.2)
  minor     Create a minor release (0.1.0 -> 0.2.0)
  major     Create a major release (1.0.0 -> 2.0.0)
  version   Create a release with specific version
  help      Show this help message

Options:
  -m, --message MESSAGE    Release message (default: "Release VERSION")
  --dry-run                Show what would be done without executing
  --skip-tests             Skip running tests
  --skip-build             Skip building binaries

Examples:
  $0 patch
  $0 minor -m "Add new features"
  $0 version v1.2.3 -m "Bug fixes"
  $0 major --dry-run

EOF
}

# Main function
main() {
    local command="$1"
    local version=""
    local message=""
    local dry_run=false
    local skip_tests=false
    local skip_build=false
    
    # Parse arguments
    shift
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--message)
                message="$2"
                shift 2
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            --skip-tests)
                skip_tests=true
                shift
                ;;
            --skip-build)
                skip_build=true
                shift
                ;;
            *)
                if [ -z "$version" ] && [[ "$1" =~ ^v[0-9] ]]; then
                    version="$1"
                else
                    echo -e "${RED}Error: Unknown option $1${NC}"
                    show_help
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Handle commands
    case "$command" in
        "patch"|"minor"|"major")
            local current_version
            current_version=$(get_current_version)
            version=$(get_next_version "$current_version" "$command")
            ;;
        "version")
            if [ -z "$version" ]; then
                echo -e "${RED}Error: Version required for 'version' command${NC}"
                show_help
                exit 1
            fi
            ;;
        "help"|"-h"|"--help")
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown command '$command'${NC}"
            show_help
            exit 1
            ;;
    esac
    
    # Validate version
    validate_version "$version"
    
    # Set default message
    if [ -z "$message" ]; then
        message="Release $version"
    fi
    
    echo -e "${BLUE}GitX Release Script${NC}"
    echo -e "${BLUE}==================${NC}"
    echo -e "Version: $version"
    echo -e "Message: $message"
    echo -e "Dry run: $dry_run"
    echo ""
    
    if [ "$dry_run" = true ]; then
        echo -e "${YELLOW}DRY RUN - No changes will be made${NC}"
        echo -e "${BLUE}Would execute:${NC}"
        echo -e "  1. Check working directory is clean"
        echo -e "  2. Check on $MAIN_BRANCH branch"
        echo -e "  3. Run tests"
        echo -e "  4. Build binaries"
        echo -e "  5. Create tag $version"
        echo -e "  6. Push to remote"
        echo -e "  7. Create release"
        exit 0
    fi
    
    # Pre-release checks
    check_clean_working_dir
    check_main_branch
    
    # Run tests
    if [ "$skip_tests" = false ]; then
        run_tests
    fi
    
    # Build binaries
    if [ "$skip_build" = false ]; then
        build_binaries
    fi
    
    # Create tag
    create_tag "$version" "$message"
    
    # Push to remote
    push_to_remote "$version"
    
    # Create release
    create_release "$version" "$message"
    
    echo -e "${GREEN}🎉 Release $version created successfully!${NC}"
    echo -e "${BLUE}Check the GitHub Actions tab for build progress${NC}"
}

# Run main function
main "$@"
