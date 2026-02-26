# 贡献指南

感谢你对 AI AgentGuard 项目的关注！我们欢迎社区的各种贡献。

[English](CONTRIBUTING.md) | [简体中文](CONTRIBUTING_zh.md)

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [代码规范](#代码规范)
- [提交信息规范](#提交信息规范)
- [Pull Request 流程](#pull-request-流程)
- [报告问题](#报告问题)
- [功能建议](#功能建议)

## 行为准则

请在所有交流中保持尊重和建设性。我们致力于维护一个热情和包容的社区。

## 如何贡献

有多种方式可以做出贡献：

- 🐛 **报告错误** - 帮助我们发现和修复问题
- 💡 **建议功能** - 分享你的改进想法
- 📝 **改进文档** - 帮助完善文档内容
- 🔧 **提交 Pull Request** - 修复错误、添加功能或改进代码
- 💬 **参与讨论** - 分享知识并帮助他人
- 🎨 **设计改进** - 贡献 UI/UX 想法

## 开发环境设置

### 前置要求

- Go 1.25 或更高版本
- Git
- Make（可选，用于构建自动化）

### Fork 和克隆

```bash
# 在 GitHub 上 Fork 仓库
# 克隆你的 fork
git clone https://github.com/YOUR_USERNAME/AIAgentGuard.git
cd AIAgentGuard

# 添加上游远程仓库
git remote add upstream https://github.com/imdlan/AIAgentGuard.git
```

### 构建和测试

```bash
# 构建二进制文件
go build -o agent-guard

# 运行测试
go test ./...

# 本地安装
go install
```

### 开发工作流

```bash
# 创建新分支
git checkout -b feature/your-feature-name

# 进行修改
# ...

# 运行测试
go test ./...

# 提交更改
git commit -m "feat: 添加你的功能"

# 推送到你的 fork
git push origin feature/your-feature-name
```

## 代码规范

### Go 代码风格

遵循标准的 Go 惯例：
- 使用 `gofmt` 格式化代码
- 运行 `go vet` 检查问题
- 为导出的函数编写清晰、描述性的注释
- 保持函数专注和简洁

### 示例

```go
// ScanFilesystem 扫描文件系统的安全风险
// 并返回包含详细发现的 ScanResult。
func ScanFilesystem() ScanResult {
    // 实现
}
```

### 项目结构

```
agent-guard/
├── cmd/              # CLI 命令
├── internal/         # 内部包
│   ├── scanner/     # 安全扫描引擎
│   ├── risk/        # 风险分析
│   ├── sandbox/     # 沙箱执行
│   ├── policy/      # 策略管理
│   ├── security/    # 安全防护
│   └── report/      # 报告生成
├── pkg/model/       # 公共数据模型
├── configs/         # 默认配置
└── scripts/         # 安装脚本
```

### 命名规范

- **包名**: 小写，单词（`scanner`、`risk`、`policy`）
- **导出函数**: PascalCase（`RunAllScans`、`Analyze`、`CalculateRisk`）
- **私有函数**: camelCase（`calculateOverallRisk`、`generateReport`）
- **常量**: PascalCase（`Low`、`Medium`、`High`、`Critical`）

## 提交信息规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

### 格式

```
<type>(<scope>): <description>

[可选的正文]

[可选的脚注]
```

### 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更改
- `style`: 代码风格更改（格式化等）
- `refactor`: 代码重构
- `perf`: 性能改进
- `test`: 添加或更新测试
- `chore`: 维护任务
- `ci`: CI/CD 更改

### 示例

```bash
feat(scanner): 添加网络端口扫描
fix(risk): 修正风险分数计算
docs(readme): 更新安装说明
refactor(policy): 简化策略加载逻辑
test(sandbox): 添加沙箱隔离测试
```

## Pull Request 流程

### 提交前

1. **搜索现有 PR** - 避免重复工作
2. **更新文档** - 包含相关文档
3. **添加测试** - 确保新代码有测试
4. **运行 linters** - 检查代码质量
5. **更新 CHANGELOG** - 记录更改

### 提交 PR

1. 创建描述性标题
2. 在正文中描述你的更改
3. 链接相关问题
4. 确保所有检查通过
5. 请求维护者审查

### PR 模板

```markdown
## 描述
更改的简要描述

## 更改类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 破坏性更改
- [ ] 文档更新

## 测试
- [ ] 已添加/更新单元测试
- [ ] 已完成手动测试
- [ ] 所有测试通过

## 检查清单
- [ ] 代码遵循项目风格指南
- [ ] 已更新文档
- [ ] 提交信息遵循规范
- [ ] 无合并冲突
```

### 审查流程

- 维护者将审查你的 PR
- 回应审查意见
- 保持讨论专注和建设性
- 请耐心等待 - 审查可能需要时间

## 报告问题

### 报告前

- 检查现有 issue
- 搜索类似问题
- 确认尚未修复

### Bug 报告模板

```markdown
## 描述
错误的清晰描述

## 复现步骤
1. 第一步
2. 第二步
3. 第三步

## 预期行为
应该发生什么

## 实际行为
实际发生了什么

## 环境
- 操作系统: [例如 macOS 14.0]
- Go 版本: [例如 1.25.5]
- AgentGuard 版本: [例如 v1.0.1]

## 日志
相关的错误消息或日志
```

## 功能建议

### 功能请求模板

```markdown
## 描述
功能描述

## 问题陈述
这解决了什么问题？

## 建议的解决方案
应该如何工作？

## 替代方案
考虑过的其他方法

## 附加信息
截图、示例等
```

## 获取帮助

- 💬 **讨论区**: https://github.com/imdlan/AIAgentGuard/discussions
- 🐛 **问题追踪**: https://github.com/imdlan/AIAgentGuard/issues
- 📧 **邮件**: imdlan@users.noreply.github.com

## 致谢

贡献者将在以下地方被认可：
- CONTRIBUTORS.md 文件
- 发布说明
- 项目文档

感谢你对 AI AgentGuard 的贡献！🛡️

## 许可证

通过贡献，你同意你的贡献将在 MIT 许可证下授权。
