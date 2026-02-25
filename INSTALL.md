# AI AgentGuard - 快速安装指南

## 推荐安装方式

### macOS / Linux 用户（最简单）

```bash
# 方式 1: 一键安装脚本
curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash
```

### Homebrew 用户

```bash
brew tap imdlan/AIAgentGuard
brew install agent-guard
```

### Go 开发者

```bash
go install github.com/imdlan/AIAgentGuard@latest
```

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/imdlan/AIAgentGuard.git
cd agent-guard

# 使用 Makefile（推荐）
make build
make install  # 需要 sudo 权限

# 或手动构建
go build -o agent-guard
sudo mv agent-guard /usr/local/bin/
```

## 验证安装

```bash
agent-guard --help
agent-guard scan
```

## 快速开始

```bash
# 1. 扫描当前环境的安全风险
agent-guard scan

# 2. 在沙箱中运行命令
agent-guard run "curl https://api.example.com"

# 3. 生成安全报告
agent-guard report

# 4. 初始化配置文件
agent-guard init
```

## 需要帮助？

查看完整文档: [README.md](README.md)

报告问题: [GitHub Issues](https://github.com/imdlan/AIAgentGuard/issues)
