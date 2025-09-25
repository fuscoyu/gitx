#cmd=patch
cmd=go run main.go
p=doc #项目
j=VM-1888 #jiraID
b=dev #分支 可以逗号分隔

build:
	go build -o $(shell echo $$GOPATH/bin/gitx)  main.go

.PHONY: bin
bin:
	@mkdir -p bin
	@echo "Building for all platforms..."
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ); \
	GIT_COMMIT=$$(git rev-parse HEAD 2>/dev/null || echo "unknown"); \
	echo "Version: $$VERSION"; \
	echo "Build Time: $$BUILD_TIME"; \
	echo "Git Commit: $$GIT_COMMIT"; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$$VERSION -X main.buildTime=$$BUILD_TIME -X main.gitCommit=$$GIT_COMMIT" -o bin/gitx_linux_amd64 main.go; \
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$$VERSION -X main.buildTime=$$BUILD_TIME -X main.gitCommit=$$GIT_COMMIT" -o bin/gitx_linux_arm64 main.go; \
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$$VERSION -X main.buildTime=$$BUILD_TIME -X main.gitCommit=$$GIT_COMMIT" -o bin/gitx_darwin_amd64 main.go; \
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$$VERSION -X main.buildTime=$$BUILD_TIME -X main.gitCommit=$$GIT_COMMIT" -o bin/gitx_darwin_arm64 main.go; \
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$$VERSION -X main.buildTime=$$BUILD_TIME -X main.gitCommit=$$GIT_COMMIT" -o bin/gitx_windows_amd64.exe main.go
	@echo "Build completed! Binaries are in bin/ directory"

# Build for specific platform
.PHONY: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
build-linux-amd64:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gitx_linux_amd64 main.go

build-linux-arm64:
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gitx_linux_arm64 main.go

build-darwin-amd64:
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gitx_darwin_amd64 main.go

build-darwin-arm64:
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gitx_darwin_arm64 main.go

build-windows-amd64:
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gitx_windows_amd64.exe main.go


add:
	$(cmd) jira -a=add -p=$(p) -j=$(j) -b=$(b)

del:
	$(cmd) jira -a=del -p=$(p) -j=$(j)

print:
	$(cmd) jira -a=print

clear:
	$(cmd) jira -a=clear

push:
	$(cmd) push -p=$(p) -j=$(j)  -b=$(b)


git_clear:
	git filter-branch --force --prune-empty --index-filter 'git rm -rf --cached --ignore-unmatch bin/*' --tag-name-filter cat -- --all

# Release management
.PHONY: release-patch release-minor release-major release-version
release-patch:
	@chmod +x scripts/release.sh
	@./scripts/release.sh patch

release-minor:
	@chmod +x scripts/release.sh
	@./scripts/release.sh minor

release-major:
	@chmod +x scripts/release.sh
	@./scripts/release.sh major

release-version:
	@chmod +x scripts/release.sh
	@./scripts/release.sh version $(VERSION)

# Show version info
.PHONY: version
version:
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'v0.0.0')"
	@echo "Latest commit: $$(git rev-parse --short HEAD)"
	@echo "Branch: $$(git branch --show-current)"
