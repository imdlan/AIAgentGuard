# Homebrew Tap 仓库设置指南

本文档说明如何设置 Homebrew tap 仓库以支持 `brew install agent-guard`。

## 背景说明

`goreleaser` 配置中指定的 tap 仓库是 `imdlan/homebrew-AIAgentGuard`。这个仓库需要单独创建，用于存放 Homebrew formula 文件。

**重要**：goreleaser 会在每次 release 时自动向这个仓库提交 formula 文件，你只需要创建仓库即可。

## 步骤 1：创建 Homebrew Tap 仓库

1. 访问 GitHub 并创建新仓库：`https://github.com/new`
2. 仓库名称：`homebrew-AIAgentGuard`
3. 设置为 Public（公开仓库）
4. **不要**初始化 README、.gitignore 或 license
5. 点击 "Create repository"

创建完成后，仓库地址应该是：`https://github.com/imdlan/homebrew-AIAgentGuard`

## 步骤 2：配置 GitHub Secrets

在主仓库 `imdlan/AIAgentGuard` 中添加 GitHub Secret：

1. 进入主仓库的 Settings 页面
2. 左侧菜单选择 "Secrets and variables" → "Actions"
3. 点击 "New repository secret"
4. 添加以下 secret：

```
Name: HOMEBREW_TAP_GITHUB_TOKEN
Secret: 你的 GitHub Personal Access Token
```

**创建 GitHub Token**：
1. 访问：https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 勾选权限：
   - `repo` (完整仓库访问权限)
   - `workflow` (如果需要)
4. 生成并复制 token
5. 粘贴到上面的 secret 字段

## 步骤 3：测试自动发布

完成上述步骤后，每次推送 version tag 时，goreleaser 会自动：

1. 构建多平台二进制文件
2. 创建 GitHub Release
3. **自动创建/更新 Homebrew formula 文件**
4. 自动提交到 `homebrew-AIAgentGuard` 仓库

### 测试流程

```bash
# 1. 确保在主分支
git checkout main

# 2. 创建测试 tag
git tag v1.0.1

# 3. 推送 tag（触发 release workflow）
git push origin v1.0.1

# 4. 等待 workflow 完成
# 访问：https://github.com/imdlan/AIAgentGuard/actions

# 5. 检查 homebrew tap 仓库是否有新的 formula 文件
# 访问：https://github.com/imdlan/homebrew-AIAgentGuard

# 6. 测试安装
brew tap imdlan/AIAgentGuard
brew install agent-guard
```

## Homebrew Formula 文件结构

goreleaser 会自动生成类似这样的 formula 文件：

```ruby
# Formula/agent-guard.rb
class AgentGuard < Formula
  desc "AI Agent, CLI tools, and MCP server security scanning tool"
  homepage "https://github.com/imdlan/AIAgentGuard"
  version "v1.0.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/imdlan/AIAgentGuard/releases/download/v1.0.1/agent-guard_1.0.1_darwin_arm64.tar.gz"
      sha256 "darwin_arm64_sha256_hash"
    else
      url "https://github.com/imdlan/AIAgentGuard/releases/download/v1.0.1/agent-guard_1.0.1_darwin_amd64.tar.gz"
      sha256 "darwin_amd64_sha256_hash"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/imdlan/AIAgentGuard/releases/download/v1.0.1/agent-guard_1.0.1_linux_arm64.tar.gz"
      sha256 "linux_arm64_sha256_hash"
    else
      url "https://github.com/imdlan/AIAgentGuard/releases/download/v1.0.1/agent-guard_1.0.1_linux_amd64.tar.gz"
      sha256 "linux_amd64_sha256_hash"
    end
  end

  def install
    bin.install "agent-guard"
  end

  test do
    system "#{bin}/agent-guard --help"
  end
end
```

**注意**：你不需要手动创建这个文件，goreleaser 会自动生成并提交到 tap 仓库。

## 故障排查

### 问题 1：goreleaser 提示 "repository not found"

**原因**：homebrew tap 仓库不存在
**解决**：按照步骤 1 创建 `homebrew-AIAgentGuard` 仓库

### 问题 2：goreleaser 提示 "permission denied"

**原因**：GitHub Token 权限不足或未配置
**解决**：按照步骤 2 配置 `HOMEBREW_TAP_GITHUB_TOKEN`，确保 token 有 `repo` 权限

### 问题 3：formula 文件生成但没有提交到 tap 仓库

**原因**：可能是第一次运行，goreleaser 需要克隆空仓库
**解决**：确保 tap 仓库是完全空的（没有 README），如果 GitHub 自动创建了 README，需要删除它

```bash
cd /tmp
git clone https://github.com/imdlan/homebrew-AIAgentGuard.git
cd homebrew-AIAgentGuard
rm -rf README.md .gitignore
git commit -am "Initial commit" || true
git push origin main
```

## 自动化流程总结

```
推送 tag (v1.0.1)
    ↓
触发 GitHub Actions
    ↓
运行测试
    ↓
goreleaser 构建
    ├─ 构建 4 个平台二进制文件
    ├─ 创建 GitHub Release
    └─ 自动更新 Homebrew tap
         ↓
    用户可以直接安装：
    brew install agent-guard
```

## 相关文档

- [Goreleaser Homebrew 配置](https://goreleaser.com/customization/homebrew/)
- [Homebrew Tap 仓库最佳实践](https://docs.brew.sh/Taps)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
