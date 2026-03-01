# AIAgentGuard v1.4.0 - 完整使用指南

> English version: [USAGE.md](USAGE.md)

## 📋 目录

1. [项目概述](#项目概述)
2. [核心使用场景](#核心使用场景)
3. [快速开始](#快速开始)
4. [CLI 命令行工具使用](#cli-命令行工具使用)
5. [Web UI 仪表板](#web-ui-仪表板)
6. [监控与告警](#监控与告警)
7. [配置与定制](#配置与定制)
8. [部署指南](#部署指南)
9. [日常维护](#日常维护)
10. [故障排除](#故障排除)

---

## 项目概述

**AIAgentGuard** 是一个企业级 AI Agent 安全扫描和监控工具，提供：

### 🔒 核心能力
- **多语言依赖漏洞扫描** - 支持 Go、npm、pip、cargo
- **系统权限扫描** - 文件系统、Shell、网络、机密信息
- **容器运行时检测** - Docker、Kubernetes、Podman、LXC、Wasm
- **沙盒隔离执行** - containerd/gVisor 容器隔离
- **实时监控指标** - Prometheus + Grafana 仪表板
- **Web 可视化界面** - React + Go RESTful API

### 📊 安全覆盖率
- **v1.2.0**: 92%+ (+多语言扫描 + 监控)
- **v1.4.0**: **95%+** (+详细报告 + 修复向导 + 趋势分析)
- **v1.1.0**: 78% (+企业特性)
- **v1.2.0**: **92%+** (+多语言扫描 + 监控)

---

## 核心使用场景

### 场景 1: CI/CD 流水线安全检查

**目标**: 在代码提交/合并前自动检测安全问题

**使用方式**:
```yaml
# .github/workflows/security-scan.yml
name: Security Scan

on: [push, pull_request]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install AIAgentGuard
        run: |
          wget https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_linux_amd64.tar.gz
          tar -xzf agent-guard_linux_amd64.tar.gz
          chmod +x agent-guard
          sudo mv agent-guard /usr/local/bin/
      
      - name: Run Security Scan
        run: |
          agent-guard scan --json > security-report.json
      
      - name: Check Results
        run: |
          CRITICAL=$(jq '.overall' security-report.json)
          if [ "$CRITICAL" = "CRITICAL" ]; then
            echo "❌ Critical security issues found!"
            exit 1
          fi
      
      - name: Upload Report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: security-report
          path: security-report.json
```

### 场景 2: 开发环境实时监控

**目标**: 开发时持续监控项目安全状态

**使用方式**:
```bash
# 终端 1: 启动监控服务
agent-guard scan --metrics-addr :9090

# 终端 2: 启动 Web UI (可选)
cd webui && docker-compose up

# 终端 3: 持续监控 (watch 模式)
watch -n 5 'agent-guard scan | jq .overall'
```

### 场景 3: 生产环境安全审计

**目标**: 定期安全审计和合规检查

**使用方式**:
```bash
# 1. 生成详细审计报告
agent-guard report --json > audit-$(date +%Y%m%d).json

# 2. 扫描特定目录
agent-guard scan --dir /path/to/project

# 3. 只扫描依赖漏洞
agent-guard scan --category dependencies
agent-guard scan --category npmdeps
agent-guard scan --category pipdeps
```

### 场景 4: 容器化环境安全扫描

**目标**: 扫描 Docker/Kubernetes 容器镜像

**使用方式**:
```bash
# 扫描容器内部
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock:ro \
  agent-guard:latest scan

# 扫描 Kubernetes Pod
kubectl exec -it <pod-name> -- /agent-guard scan

# 扫描特定镜像
docker run --rm agent-guard:latest scan \
  --dir /app
```

---

## 快速开始

### 方式 1: Homebrew 安装 (推荐)

```bash
brew tap imdlan/AIAgentGuard
brew install agent-guard

# 运行扫描
agent-guard scan
```

### 方式 2: 下载二进制

```bash
# macOS/Linux ARM64
curl -LO https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_darwin_arm64.tar.gz
tar -xzf agent-guard_darwin_arm64.tar.gz
chmod +x agent-guard
sudo mv agent-guard /usr/local/bin/

# Linux AMD64
curl -LO https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_linux_amd64.tar.gz
tar -xzf agent-guard_linux_amd64.tar.gz
chmod +x agent-guard
sudo mv agent-guard /usr/local/bin/
```

### 方式 3: Go 编译安装

```bash
go install github.com/imdlan/AIAgentGuard@latest
export PATH=$PATH:$(go env GOPATH)/bin:$PATH
```

### 方式 4: Docker 运行

```bash
docker run --rm -v $(pwd):/app:ro \
  imdlan/agent-guard:latest scan
```

---

## CLI 命令行工具使用

### 基础扫描

```bash
# 完整安全扫描（所有类别）
agent-guard scan

# JSON 格式输出
agent-guard scan --json

# 详细输出
agent-guard scan --verbose

# 使用自定义配置文件
agent-guard scan --config /path/to/policy.yaml
```

### 扫描特定类别

```bash
# 只扫描文件系统
agent-guard scan --category filesystem

# 只扫描依赖漏洞
agent-guard scan --category dependencies
agent-guard scan --category npmdeps
agent-guard scan --category pipdeps
agent-guard scan --category cargodeps

# 扫描多个类别
agent-guard scan --category filesystem --category shell --category network
```

### 沙箱执行

```bash
# 在隔离环境中运行命令
agent-guard run "curl https://api.example.com"

# 禁用网络访问
agent-guard run --disable-network "npm install"

# 限制目录访问
agent-guard run --allow-dirs /tmp,/data "node script.js"
```

### 生成报告

```bash
# 生成报告
agent-guard report

# JSON 格式报告
agent-guard report --json > report.json

# 保存到文件
agent-guard report --output security-audit-$(date +%Y%m%d).txt
```

### 初始化配置

```bash
# 生成默认配置文件
agent-guard init

# 强制覆盖
agent-guard init --force

# 指定配置文件路径
agent-guard init --path /etc/agent-guard/config.yaml
```bash
# 生成默认配置文件
agent-guard init

# 强制覆盖
agent-guard init --force

# 指定配置文件路径
agent-guard init --path /etc/agent-guard/config.yaml
```

### 安全修复向导（v1.4.0 新增）

**自动修复安全问题或获取修复指导**

```bash
# 预览修复而不执行（推荐第一步）
agent-guard fix --dry-run

# 自动修复所有问题
agent-guard fix --auto

# 修复特定类别
agent-guard fix --category filesystem --auto
```

**示例输出**:
```
🔧 安全修复向导
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

发现 3 个安全问题：

1. [文件系统] /Users/user/.ssh 可写入
   → 修复: chmod 700 /Users/user/.ssh
   → 优先级: HIGH

2. [Shell] 历史文件包含敏感命令
   → 修复: rm ~/.bash_history
   → 优先级: MEDIUM

3. [网络] 端口 22 对外开放
   → 需要手动检查
   → 优先级: LOW

执行修复？ [y/N]: y

✅ 已修复 2 个问题
⏭️  1 个问题需要手动检查
```

**安全特性**:
- 干运行模式在执行前预览更改
- 破坏性命令需要手动确认
- 修复命令记录到审计日志

### 风险趋势分析（v1.4.0 新增）

**对比历史扫描结果追踪安全态势**

```bash
# 显示最近 7 天的趋势
agent-guard trend

# 自定义时间范围
agent-guard trend --days 30

# 特定类别趋势
agent-guard trend --category filesystem

# JSON 格式输出
agent-guard trend --json > trend-data.json
```

**趋势分析功能**:
- 对比当前扫描与历史数据
- 识别改善/恶化趋势
- 显示类别级别变化
- 计算风险分数趋势
- 可视化趋势指示器（📈 改善、📉 恶化、➡️ 稳定）

---

## Web UI 仪表板

### 启动 Web UI

```bash
# 方式 1: Docker Compose（推荐）
cd webui
docker-compose up -d

# 访问
# 前端: http://localhost:3000
# 用户名: admin
# 密码: admin

# 方式 2: 手动启动
# 后端
cd webui/backend
go run main.go
# 默认端口: 8080

# 前端
cd webui/frontend
npm install
npm run dev
# 默认端口: 5173
```

### Web UI 功能

**实时监控面板**:
- 📊 扫描速率和持续时间
- 🎯 漏洞统计和趋势
- 📈 按语言分类的漏洞条形图
- 🔄 30秒自动刷新

**扫描结果展示**:
- 🔍 一键执行安全扫描
- 📋 完整权限分解（9 个类别）
- 🎨 颜色编码风险等级
- 📝 详细发现列表

**系统状态监控**:
- 📊 版本和运行状态
- ⏱️ 运行时间统计
- 🔧 扫描器状态

### API 端点

```bash
# 执行扫描
GET  /api/v1/scan

# 自定义扫描选项
POST /api/v1/scan
{
  "categories": ["filesystem", "shell"],
  "options": {
    "include_file_content": true,
    "timeout": 60
  }
}

# 获取监控指标
GET  /api/v1/metrics
GET  /api/v1/metrics/scan-rate
GET  /api/v1/metrics/vulnerabilities
GET  /api/v1/metrics/duration

# 系统状态
GET  /api/v1/status

# 安全告警
GET  /api/v1/alerts
```

---

## 监控与告警

### Prometheus 指标

**扫描指标**:
```promql
# 扫描总数
agent_guard_scans_total

# 扫描速率
rate(agent_guard_scans_total[5m])

# 扫描持续时间
agent_guard_scan_duration_seconds
histogram_quantile(0.95, agent_guard_scan_duration_seconds)
```

**漏洞指标**:
```promql
# 总漏洞数
sum(agent_guard_vulnerabilities_total)

# 按严重性
sum by (severity) (agent_guard_vulnerabilities_total)

# 按语言和严重性
agent_guard_vulnerabilities_total{severity="critical", language="go"}
```

### 启动监控服务

```bash
# AIAgentGuard 内置 Prometheus 服务器
agent-guard scan --metrics-addr :9090

# 访问指标
curl http://localhost:9090/metrics
```

### Grafana 仪表板导入

1. 访问 Grafana: http://localhost:3000
2. 登录: admin/admin
3. Dashboards → Import
4. 上传 `configs/grafana-dashboard.json`
5. 选择 Prometheus 数据源

### 完整监控栈部署

```bash
cd configs
docker-compose -f docker-compose.monitoring.yml up -d

# 服务端口
# AIAgentGuard Backend: 8080
# Web UI Frontend: 3000
# Prometheus: 9091
# Grafana: 3000
```

### 告警规则示例

```yaml
# 关键漏洞告警
- alert: CriticalVulnerabilitiesDetected
  expr: sum(agent_guard_vulnerabilities_total{severity="critical"}) > 0
  for: 5m
  annotations:
    summary: "Critical vulnerabilities detected"

# 扫描超时告警
- alert: ScanDurationTooHigh
  expr: histogram_quantile(0.95, agent_guard_scan_duration_seconds) > 300
  for: 15m
```

---

## 配置与定制

### 策略配置文件

**位置**:
- `.agent-guard.yaml` (当前目录)
- `~/.agent-guard.yaml` (用户目录)
- `/etc/agent-guard/config.yaml` (系统目录)

**配置示例**:

```yaml
version: 1

# 阻止的危险命令
blocked_commands:
  - "rm -rf /"
  - "dd if=/dev/zero"
  - "mkfs"
  - ":(){ :|:& };:"  # fork bomb

# 文件系统访问控制
filesystem:
  allow:
    - /tmp
    - /home/user/project
  deny:
    - /etc/passwd
    - /etc/shadow
    - ~/.ssh
    - /root

# Shell 命令控制
shell:
  allow:
    - cat
    - ls
    - grep
    "ps aux"
  deny:
    - "rm "
    - "dd "
    - ":(){"

# 网络访问控制
network:
  allow:
    - api.github.com
    - cdn.jsdelivr.net
  deny:
    - "*.malicious.com"
    - "10.0.0.0/8"

# 环境变量保护
secrets:
  block_env:
    - API_KEY
    - SECRET_TOKEN
    - DATABASE_URL
    - PRIVATE_KEY

# 沙箱执行配置
sandbox:
  disable_network: false
  readonly_root: false
```

### 环境变量配置

```bash
# 后端 API 配置
PORT=8080                    # API 服务端口
METRICS_ADDR=:9090          # Prometheus 指标端口

# 前端配置
VITE_API_URL=http://localhost:8080  # 后端 API 地址
```

---

## 部署指南

### Docker 部署

```bash
# 1. 构建镜像
docker build -t agent-guard:latest .

# 2. 运行容器
docker run --rm -v /path/to/scan:/app:ro \
  agent-guard:latest scan

# 3. 后台服务模式
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -p 8080:8080 \
  agent-guard:latest \
  scan --metrics-addr :9090
```

### Kubernetes 部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agent-guard-scanner
spec:
  replicas: 1
  selector:
    matchLabels:
      app: agent-guard
  template:
    metadata:
      labels:
        app: agent-guard
    spec:
      containers:
      - name: agent-guard
        image: imdlan/agent-guard:v1.2.0
        command: ["scan", "--metrics-addr", ":9090"]
        ports:
        - containerPort: 8080
          name: metrics
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: agent-guard-metrics
spec:
  selector:
    app: agent-guard
  ports:
    - port: 9090
      targetPort: 9090
```

### Docker Compose 完整栈

```bash
# 一键启动所有服务
cd webui
docker-compose -f docker-compose.yml up -d

# 包含服务
# - AIAgentGuard 后端 (8080)
# - React 前端 (3000)
# - Prometheus (9091)
# - Grafana (3000)
```

---

## 日常维护

### 定期安全扫描

```bash
# 每日扫描
0 2 * * * agent-guard scan >> /var/log/security-scan.log 2>&1

# 每周完整审计
0 3 * * 1 agent-guard report --json > /backups/security/weekly-$(date +\%Y%m%d).json

# 每月依赖检查
0 4 * * 1 agent-guard scan --category npmdeps,pipdeps,cargodeps
```

### 监控维护

```bash
# 检查 Prometheus 指标
curl http://localhost:9090/metrics | grep agent_guard

# 查看内存使用
curl http://localhost:9090/metrics | grep memory

# 查看运行时间
curl http://localhost:9090/metrics | grep uptime
```

### 日志管理

```bash
# 审计日志位置
~/.agent-guard/audit.log
/var/log/agent-guard/

# 日志轮转
logrotate ~/.agent-guard/audit.log {
  weekly
  rotate 52
  compress
  delaycompress
  missingok
  notifempty
}
```

### 更新维护

```bash
# 检查版本
agent-guard --help
grep "var version" main.go

# 更新到最新版本
brew upgrade agent-guard

# 或下载新版本
curl -LO https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_linux_amd64.tar.gz
```

### 数据库更新

```bash
# 更新 Go 漏洞数据库
go run golang.org/x/vuln/cmd/govulncheck@latest download

# 更新 npm 审计数据库
cd /path/to/npm-project
npm audit fix
```

---

## 故障排除

### 常见问题

#### 1. "command not found: agent-guard"

**解决方案**:
```bash
# 检查安装
which agent-guard

# 添加到 PATH (临时)
export PATH=$PATH:/usr/local/bin

# 永久添加
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. "permission denied" 错误

**解决方案**:
```bash
# 添加执行权限
chmod +x agent-guard

# 或使用绝对路径运行
./agent-guard scan
```

#### 3. npm/yarn/pip/cargo 扫描失败

**原因**: 对应工具未安装

**解决方案**:
```bash
# 安装 npm 审计工具
npm install -g audit-parser
npm audit fix

# 安装 Python 审计工具
pip install pip-audit
pip-audit

# 安装 Rust 审计工具
cargo install cargo-audit
```

#### 4. Prometheus 端点无数据

**检查**:
```bash
# 检查服务是否运行
curl http://localhost:9090/metrics

# 检查是否启用了指标
ps aux | grep "agent-guard.*metrics"
```

**解决方案**:
```bash
# 重新启动并启用指标
agent-guard scan --metrics-addr :9090
```

#### 5. Web UI 无法连接后端

**检查**:
```bash
# 检查后端状态
curl http://localhost:8080/api/v1/status

# 检查 CORS
curl -H "Origin: http://localhost:3000" \
  http://localhost:8080/api/v1/status -v
```

**解决方案**:
```bash
# 确保后端服务在 8080 端口
cd webui/backend
go run main.go

# 检查防火墙
sudo ufw allow 8080
```

#### 6. 扫描超时

**解决方案**:
```bash
# 增加超时时间
agent-guard scan --timeout 120

# 或跳过耗时操作
agent-guard scan --category filesystem --category shell
```

#### 7. Docker 扫描失败

**解决方案**:
```bash
# 挂载 Docker socket
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  agent-guard:latest scan

# 在容器中挂载目录
docker run --rm -v /path/to/project:/app:ro \
  agent-guard:latest scan --dir /app
```

### 调试模式

```bash
# 启用详细输出
agent-guard scan --verbose

# JSON 输出（便于解析）
agent-guard scan --json | jq .

# 检查配置文件
agent-guard init --dry-run
```

---

## 最佳实践

### 安全扫描最佳实践

1. **CI/CD 集成**
   - 每次代码推送前扫描
   - 阻止合并高风险代码
   - 自动化审计报告归档

2. **定期审计**
   - 每周完整扫描
   - 每月依赖更新检查
   - 每季度容器镜像扫描

3. **监控告警**
   - 配置关键指标告警
   - 集成到现有监控系统
   - 建立告警响应流程

4. **配置管理**
   - 使用版本控制策略文件
   - 环境特定配置（开发/测试/生产）
   - 定期审查安全策略

### 性能优化

```bash
# 只扫描需要的类别
agent-guard scan --category dependencies --category npmdeps

# 跳过耗时操作
agent-guard scan --category filesystem --category shell

# 并发扫描（多目录）
for dir in /path/to/projects/*; do
  agent-guard scan --dir "$dir" &
done
wait
```

### 安全加固

1. **最小权限原则**
   ```yaml
   # 只授予必要的权限
   filesystem:
     allow: ["/app", "/tmp"]
     deny: ["/etc", "/root"]
   ```

2. **网络隔离**
   ```bash
   # 扫描时禁用网络
   agent-guard run --disable-network "curl api.example.com"
   ```

3. **沙盒执行**
   ```bash
   # 在隔离环境中运行
   agent-guard run "npm install"
   ```

---

## 总结

AIAgentGuard v1.2.0 提供了：

### ✅ 已完成功能
1. **多语言依赖扫描** - Go、npm、pip、cargo
2. **实时监控面板** - Prometheus + Grafana
3. **Web UI 仪表板** - 可视化安全状态
4. **容器环境支持** - Docker、Kubernetes
5. **沙盒隔离执行** - 安全运行命令
6. **完整审计日志** - 安全事件追踪

### 🎯 适用场景
- CI/CD 流水线安全检查
- 开发环境实时监控
- 生产环境定期审计
- 容器化镜像扫描
- 企业级安全合规

### 📈 维护要点
- 定期更新依赖漏洞数据库
- 监控指标趋势分析
- 审查告警规则有效性
- 备份审计日志和报告
- 更新安全策略配置

### 🚀 下一步计划
- 实时监控集成 (WebSocket)
- 插件系统 (自定义扫描器)
- 高级告警 (机器学习异常检测)
- 性能优化 (指标批处理)

**快速开始只需 3 步**:
```bash
1. brew install agent-guard
2. agent-guard scan
3. 访问 http://localhost:3000 查看仪表板
```

**保护您的 AI Agents，从 AIAgentGuard 开始！** 🛡️
