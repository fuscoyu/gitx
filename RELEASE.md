# GitX 发布管理指南

## 🚀 快速开始

### 创建发布

```bash
# 创建补丁版本 (1.0.0 -> 1.0.1)
make release-patch

# 创建次要版本 (1.0.0 -> 1.1.0)
make release-minor

# 创建主要版本 (1.0.0 -> 2.0.0)
make release-major

# 创建指定版本
make release-version VERSION=v1.2.3
```

### 查看版本信息

```bash
# 查看当前版本信息
make version

# 查看二进制文件版本
./bin/gitx_linux_amd64 version
```

## 📋 发布流程

### 1. 自动发布流程

当推送 Git tag 时，GitHub Actions 会自动：

1. **运行测试** - 确保代码质量
2. **多平台构建** - 支持 Linux (AMD64/ARM64), macOS (Intel/Apple Silicon), Windows (AMD64)
3. **创建发布** - 自动生成 GitHub Release
4. **上传文件** - 包含所有平台的二进制文件和校验和

### 2. 支持的平台

- **Linux AMD64** - `gitx_linux_amd64`
- **Linux ARM64** - `gitx_linux_arm64`
- **macOS Intel** - `gitx_darwin_amd64`
- **macOS Apple Silicon** - `gitx_darwin_arm64`
- **Windows AMD64** - `gitx_windows_amd64.exe`

### 3. 版本信息注入

每个二进制文件都包含以下版本信息：

- **版本号** - 从 Git tag 获取
- **构建时间** - UTC 时间戳
- **Git 提交** - 完整的 commit SHA
- **Go 版本** - 构建时使用的 Go 版本

## 🛠️ 手动发布

### 使用发布脚本

```bash
# 基本用法
./scripts/release.sh patch
./scripts/release.sh minor
./scripts/release.sh major
./scripts/release.sh version v1.2.3

# 带自定义消息
./scripts/release.sh patch -m "修复重要 bug"

# 预览模式（不执行实际操作）
./scripts/release.sh patch --dry-run

# 跳过测试
./scripts/release.sh patch --skip-tests

# 跳过构建
./scripts/release.sh patch --skip-build
```

### 手动步骤

```bash
# 1. 确保工作目录干净
git status

# 2. 切换到主分支
git checkout main

# 3. 运行测试
go test -v ./...

# 4. 构建所有平台
make bin

# 5. 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 6. 推送到远程
git push origin main
git push origin v1.0.0
```

## 📦 发布文件

每次发布会生成以下文件：

### 压缩包
- `gitx-{version}-binaries.tar.gz` - 包含所有平台二进制文件
- `gitx-{version}-linux_amd64.tar.gz` - Linux AMD64
- `gitx-{version}-linux_arm64.tar.gz` - Linux ARM64
- `gitx-{version}-darwin_amd64.tar.gz` - macOS Intel
- `gitx-{version}-darwin_arm64.tar.gz` - macOS Apple Silicon
- `gitx-{version}-windows_amd64.tar.gz` - Windows AMD64

### 校验文件
- `checksums.txt` - 所有二进制文件的 SHA256 校验和

## 🔧 本地构建

### 构建所有平台

```bash
make bin
```

### 构建特定平台

```bash
make build-linux-amd64
make build-linux-arm64
make build-darwin-amd64
make build-darwin-arm64
make build-windows-amd64
```

### 构建单个平台

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o gitx_linux_amd64 main.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o gitx_darwin_arm64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o gitx_windows_amd64.exe main.go
```

## 🎯 版本命名规范

### 语义化版本控制

使用 [Semantic Versioning](https://semver.org/) 规范：

- **主版本号** (MAJOR) - 不兼容的 API 修改
- **次版本号** (MINOR) - 向下兼容的功能性新增
- **修订号** (PATCH) - 向下兼容的问题修正

### 预发布版本

- `v1.0.0-alpha.1` - Alpha 版本
- `v1.0.0-beta.1` - Beta 版本
- `v1.0.0-rc.1` - Release Candidate

### 开发版本

- `dev-{commit}` - 开发版本（非 tag 构建）

## 📋 检查清单

发布前请确认：

- [ ] 代码已提交到主分支
- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] CHANGELOG 已更新
- [ ] 版本号符合语义化版本规范
- [ ] 工作目录干净（无未提交更改）

## 🐛 故障排除

### 常见问题

1. **标签已存在**
   ```bash
   git tag -d v1.0.0  # 删除本地标签
   git push origin :refs/tags/v1.0.0  # 删除远程标签
   ```

2. **构建失败**
   - 检查 Go 版本兼容性
   - 确保所有依赖已安装
   - 检查网络连接（下载依赖）

3. **发布失败**
   - 检查 GitHub Actions 权限
   - 确保 GITHUB_TOKEN 有效
   - 检查仓库设置

### 调试命令

```bash
# 检查版本信息
git describe --tags --always --dirty

# 检查构建信息
go version
go env

# 检查二进制文件信息
file bin/gitx_linux_amd64
ldd bin/gitx_linux_amd64  # Linux
otool -L bin/gitx_darwin_amd64  # macOS
```

## 📚 相关文档

- [GitHub Actions 工作流](.github/workflows/build.yml)
- [Makefile 构建配置](Makefile)
- [发布脚本](scripts/release.sh)
- [安装脚本](install.sh)
