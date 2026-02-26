# 发布流程指南

本文档说明如何使用 Goreleaser 自动化发布 AI AgentGuard 的新版本。

## 🚀 快速发布

```bash
# 1. 确保在主分支
git checkout main
git pull origin main

# 2. 更新版本号（如果需要）
# 编辑 Makefile 中的 VERSION 变量

# 3. 创建版本 tag
git tag v1.0.1

# 4. 推送 tag（自动触发 release workflow）
git push origin v1.0.1

# 5. 等待 GitHub Actions 完成
# 访问：https://github.com/imdlan/AIAgentGuard/actions
```

**就这么简单！** 剩下的工作由 Goreleaser 自动完成：
- ✅ 运行测试
- ✅ 构建多平台二进制文件
- ✅ 创建 GitHub Release
- ✅ 生成 checksums.txt
- ✅ 自动更新 Homebrew formula

## 📋 发布前检查清单

在发布新版本前，请确保：

- [ ] 所有测试通过：`go test ./...`
- [ ] 代码已提交到 main 分支
- [ ] 更新了 CHANGELOG.md（如果有重大变更）
- [ ] 版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)
- [ ] 已设置 Homebrew tap 仓库（首次发布）
- [ ] 已配置 `HOMEBREW_TAP_GITHUB_TOKEN` secret

## 🔄 版本号规范

项目使用语义化版本（Semantic Versioning）：`MAJOR.MINOR.PATCH`

- **MAJOR**：不兼容的 API 变更
- **MINOR**：向后兼容的功能新增
- **PATCH**：向后兼容的问题修复

示例：
- `v1.0.0` - 第一个稳定版本
- `v1.0.1` - Bug 修复
- `v1.1.0` - 新增功能
- `v2.0.0` - 重大变更

## 🛠️ Goreleaser 自动化内容

### 构建产物

每次发布会生成以下文件：

```
agent-guard_1.0.1_darwin_amd64.tar.gz
agent-guard_1.0.1_darwin_arm64.tar.gz
agent-guard_1.0.1_linux_amd64.tar.gz
agent-guard_1.0.1_linux_arm64.tar.gz
checksums.txt
```

### Homebrew 自动更新

Goreleaser 会自动：
1. 生成 Homebrew formula 文件
2. 提交到 `imdlan/homebrew-AIAgentGuard` 仓库
3. 用户可以直接使用 `brew install agent-guard` 安装

### GitHub Release

自动创建的 Release 包含：
- 所有构建的二进制文件
- checksums.txt（文件校验和）
- README.md
- LICENSE
- INSTALL.md

## 🧪 测试 Release

发布后，验证安装是否正常：

### 测试二进制文件

```bash
# 下载并测试 macOS ARM64
curl -LO https://github.com/imdlan/AIAgentGuard/releases/download/v1.0.1/agent-guard_1.0.1_darwin_arm64.tar.gz
tar -xzf agent-guard_1.0.1_darwin_arm64.tar.gz
chmod +x agent-guard
./agent-guard --version
# 应输出：agent-guard version 1.0.1
```

### 测试 Homebrew 安装

```bash
# 更新 tap
brew tap imdlan/AIAgentGuard

# 安装
brew install agent-guard

# 验证
agent-guard --version
agent-guard scan
```

## 🔧 本地构建（测试用）

如果想在本地测试构建流程（不创建 Release）：

```bash
# 安装 goreleaser
brew install goreleaser

# 测试构建（不发布）
goreleaser build --clean --snapshot

# 或者跳过发布步骤
goreleaser release --clean --snapshot --skip-publish
```

## 🐛 回滚 Release

如果发现问题需要回滚：

```bash
# 1. 删除 GitHub Release（在 GitHub 网页操作）

# 2. 删除本地 tag
git tag -d v1.0.1

# 3. 删除远程 tag
git push origin :refs/tags/v1.0.1

# 4. 修复问题后重新发布
git tag v1.0.2
git push origin v1.0.2
```

**注意**：如果 Homebrew formula 已经发布，用户可能已经安装了旧版本。考虑发布新的 PATCH 版本而不是回滚。

## 📊 发布历史

查看所有发布版本：
- GitHub Releases: https://github.com/imdlan/AIAgentGuard/releases
- Change Log: (待创建 CHANGELOG.md)

## 🆘 常见问题

### Q: goreleaser workflow 失败怎么办？

A: 检查 Actions 日志：
1. 访问 https://github.com/imdlan/AIAgentGuard/actions
2. 查看失败的任务日志
3. 常见问题：
   - 测试失败：修复测试后重新推送 tag
   - Homebrew token 无效：更新 `HOMEBREW_TAP_GITHUB_TOKEN`
   - 网络问题：重新触发 workflow

### Q: 如何修改已发布的版本？

A: Git tag 不应该修改。如果发现严重问题：
1. 发布新的 PATCH 版本（如 v1.0.1 → v1.0.2）
2. 在 Release 说明中标注修复内容
3. 考虑发布安全公告（如果是安全问题）

### Q: Homebrew formula 没有更新？

A: 检查：
1. `HOMEBREW_TAP_GITHUB_TOKEN` 是否正确配置
2. `homebrew-AIAgentGuard` 仓库是否存在
3. 查看 goreleaser 日志确认是否推送成功

### Q: 如何创建 Pre-release？

A: 使用带 prerelease 标签的版本号：
```bash
git tag v1.1.0-rc.1
git push origin v1.1.0-rc.1
```

Goreleaser 会自动将其标记为 pre-release。

## 📚 相关文档

- [Goreleaser 官方文档](https://goreleaser.com/)
- [Homebrew Tap 设置](./HOMEBREW_SETUP.md)
- [项目贡献指南](../CONTRIBUTING.md)

## 🎯 下一步

发布完成后：

1. **发布公告**：在 Discussions、Twitter 或其他渠道分享
2. **更新文档**：如有新功能，更新 README.md
3. **监控反馈**：关注 Issues 和 Discussions
4. **准备下一版本**：创建 Milestone 跟踪计划

---

**祝发布顺利！** 🎉
