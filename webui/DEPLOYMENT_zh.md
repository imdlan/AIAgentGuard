# AIAgentGuard v1.2.0 部署指南（中文）

> English version: [DEPLOYMENT.md](DEPLOYMENT.md)

本指南涵盖了在本地开发、Docker 和生产环境中部署具有监控功能的 AIAgentGuard Web UI。

## 目录

1. [前置要求](#前置要求)
2. [本地开发部署](#本地开发部署)
3. [Docker 部署](#docker-部署)
4. [生产环境部署](#生产环境部署)
5. [监控设置](#监控设置)
6. [配置](#配置)
7. [故障排除](#故障排除)
8. [安全考虑](#安全考虑)

---

## 前置要求

### 必需软件

- **Go**: 1.25.5 或更高版本
- **Node.js**: 18.x 或更高版本
- **npm**: 9.x 或更高版本
- **Git**: 用于克隆仓库

### 可选软件

- **Docker**: 20.x 或更高版本（用于容器化部署）
- **Docker Compose**: v2.x 或更高版本
- **Prometheus**: 用于指标收集和告警
- **Grafana**: 用于监控仪表板（可选）

### 系统要求

- **操作系统**: Linux、macOS 或带有 WSL2 的 Windows
- **内存**: 最小 2GB，推荐 4GB
- **磁盘**: 500MB 可用空间
- **网络**: 8080 端口（后端）、3000 端口（前端）、9090 端口（Prometheus）

---

## 本地开发部署

### 快速开始

```bash
# 克隆仓库
git clone https://github.com/imdlan/AIAgentGuard.git
cd AIAgentGuard

# 1. 启动后端（从父目录）
go run webui/backend/main.go

# 2. 在另一个终端中，启动前端（开发模式）
cd webui/frontend
npm install
npm run dev

# 3. 访问应用
# 前端: http://localhost:5173
# 后端 API: http://localhost:8080/api/v1/status
# Prometheus 指标: http://localhost:8080/metrics
```

### 后端设置

#### 从父目录运行（推荐）

```bash
cd /path/to/AIAgentGuard

# 安装依赖（仅首次需要）
go get github.com/gin-gonic/gin@v1.11.0
go mod tidy

# 启动后端服务器
go run webui/backend/main.go
```

后端将会：
- 在 8080 端口启动（可通过 `PORT` 环境变量配置）
- 在 `/api/v1/*` 提供 API 端点
- 在 `/metrics` 提供 Prometheus 指标
- 提供前端静态文件（如果 `./frontend/dist` 存在）

#### 后端环境变量

```bash
# 设置自定义端口
PORT=9090 go run webui/backend/main.go

# 启用 Go 模块下载
GOPROXY=direct go run webui/backend/main.go
```

### 前端设置

#### 开发模式

```bash
cd webui/frontend

# 安装依赖
npm install

# 启动开发服务器（支持热重载）
npm run dev

# 访问: http://localhost:5173
```

#### 生产构建

```bash
cd webui/frontend

# 构建生产版本
npm run build

# 输出: ./dist/
# - index.html
# - /assets/
#   - index-[hash].js
#   - index-[hash].css
```

### 测试设置

```bash
# 测试后端 API
curl http://localhost:8080/api/v1/status | jq .

# 测试指标端点
curl http://localhost:8080/metrics | head -20

# 执行扫描
curl http://localhost:8080/api/v1/scan | jq .

# 测试指标 API
curl http://localhost:8080/api/v1/metrics/scan-rate | jq .
```

---

## Docker 部署

### Docker Compose（Docker 环境推荐）

#### 快速开始

```bash
cd AIAgentGuard/webui

# 构建并启动所有服务
docker compose up -d

# 检查状态
docker compose ps

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

#### 服务说明

- **backend**: Go API 服务器（端口 8080）
- **frontend**: Nginx 提供 React 应用（端口 3000）

#### Docker Compose 配置

`docker-compose.yml` 已预配置：
- 正确的服务依赖关系
- 网络隔离
- 自动重启策略
- 环境变量配置

#### 构建选项

```bash
# 无缓存构建
docker compose build --no-cache

# 构建特定服务
docker compose build backend

# 重新构建并重启
docker compose up -d --build
```

### 手动 Docker 部署

#### 后端容器

```bash
cd webui/backend

# 构建镜像
docker build -t agentguard-backend:latest .

# 运行容器
docker run -d \
  --name agentguard-backend \
  -p 8080:8080 \
  -e PORT=8080 \
  --restart unless-stopped \
  agentguard-backend:latest
```

#### 前端容器

```bash
cd webui/frontend

# 构建镜像
docker build -t agentguard-frontend:latest .

# 运行容器
docker run -d \
  --name agentguard-frontend \
  -p 3000:80 \
  -e VITE_API_URL=http://backend:8080 \
  --restart unless-stopped \
  agentguard-frontend:latest
```

#### 网络设置

```bash
# 创建网络
docker network create agentguard-network

# 连接容器
docker network connect agentguard-network agentguard-backend
docker network connect agentguard-network agentguard-frontend
```

---

## 生产环境部署

### 部署架构

生产环境应使用以下架构部署：

```
                    ┌─────────────┐
                    │   Nginx     │
                    │  (反向代理)  │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
      ┌───────▼───────┐        ┌───────▼───────┐
      │   前端        │        │    后端        │
      │  (Nginx/React)│        │   (Go API)     │
      │   端口 3000   │        │   端口 8080    │
      └───────────────┘        └───────┬────────┘
                                       │
                              ┌────────▼────────┐
                              │  Prometheus     │
                              │   (指标)        │
                              │   端口 9090     │
                              └─────────────────┘
```

### Systemd 服务（Linux）

#### 后端服务

创建 `/etc/systemd/system/agentguard-backend.service`：

```ini
[Unit]
Description=AIAgentGuard 后端 API
After=network.target

[Service]
Type=simple
User=agentguard
Group=agentguard
WorkingDirectory=/opt/agentguard
Environment="PORT=8080"
ExecStart=/opt/agentguard/backend/server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启用并启动：
```bash
sudo systemctl daemon-reload
sudo systemctl enable agentguard-backend
sudo systemctl start agentguard-backend
sudo systemctl status agentguard-backend
```

#### 前端服务（Nginx）

创建 `/etc/systemd/system/agentguard-frontend.service`：

```ini
[Unit]
Description=AIAgentGuard 前端
After=network.target agentguard-backend.service

[Service]
Type=forking
PIDFile=/run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t
ExecStart=/usr/sbin/nginx
ExecReload=/bin/kill -s HUP $MAINPID
PrivateTmp=true
Restart=always

[Install]
WantedBy=multi-user.target
```

### Nginx 反向代理配置

创建 `/etc/nginx/sites-available/agentguard`：

```nginx
upstream backend {
    server localhost:8080;
}

upstream frontend {
    server localhost:3000;
}

server {
    listen 80;
    server_name agentguard.example.com;

    # 前端
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 后端 API
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Prometheus 指标（限制访问）
    location /metrics {
        proxy_pass http://backend;
        allow 127.0.0.1;
        allow 10.0.0.0/8;  # 内部网络
        deny all;
    }
}
```

启用站点：
```bash
sudo ln -s /etc/nginx/sites-available/agentguard /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Kubernetes 部署

#### 后端部署

创建 `backend-deployment.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agentguard-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: agentguard-backend
  template:
    metadata:
      labels:
        app: agentguard-backend
    spec:
      containers:
      - name: backend
        image: agentguard-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
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
  name: agentguard-backend
spec:
  selector:
    app: agentguard-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

#### 前端部署

创建 `frontend-deployment.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agentguard-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: agentguard-frontend
  template:
    metadata:
      labels:
        app: agentguard-frontend
    spec:
      containers:
      - name: frontend
        image: agentguard-frontend:latest
        ports:
        - containerPort: 80
        env:
        - name: VITE_API_URL
          value: "http://agentguard-backend:8080"
---
apiVersion: v1
kind: Service
metadata:
  name: agentguard-frontend
spec:
  selector:
    app: agentguard-frontend
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

部署：
```bash
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
```

---

## 监控设置

### Prometheus 配置

#### 抓取配置

添加到 `prometheus.yml`：

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'agentguard'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

#### 告警规则

创建 `alerts.yml`：

```yaml
groups:
  - name: agentguard_alerts
    interval: 30s
    rules:
      - alert: 高扫描率告警
        expr: scan_rate > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "检测到高扫描率"
          description: "扫描率为 {{ $value }} 次/秒"

      - alert: 漏洞检测告警
        expr: vulnerabilities_total{severity="critical"} > 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "检测到严重漏洞"
          description: "发现 {{ $value }} 个严重漏洞"
```

### Grafana 仪表板

导入预配置的仪表板：

```bash
# 仪表板位置
configs/grafana-dashboard.json

# 通过 Grafana UI 导入:
# 1. 进入 Dashboards -> Import
# 2. 上传 JSON 文件
# 3. 选择 Prometheus 数据源
# 4. 点击 Import
```

仪表板包含：
- 扫描速率趋势
- 按严重级别分类的漏洞计数
- 扫描持续时间百分位数
- 特定语言的漏洞趋势
- 实时告警

### 监控 API 端点

#### 指标信息

```bash
GET /api/v1/metrics
```

响应：
```json
{
  "message": "Prometheus 指标可在 /metrics 端点获取",
  "prometheus_url": "/metrics"
}
```

#### 扫描速率指标

```bash
GET /api/v1/metrics/scan-rate
```

响应：
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "scan_total": 125,
  "scan_rate": 2.5,
  "duration_avg": 0.85
}
```

#### 漏洞指标

```bash
GET /api/v1/metrics/vulnerabilities
```

响应：
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "vulnerabilities": {
    "critical": 0,
    "high": 2,
    "medium": 5,
    "low": 10
  },
  "by_language": {
    "go": 1,
    "npm": 3,
    "pip": 2,
    "cargo": 0
  }
}
```

#### 持续时间指标

```bash
GET /api/v1/metrics/duration
```

响应：
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "duration_p50": 0.65,
  "duration_p95": 1.2,
  "duration_p99": 1.8,
  "duration_avg": 0.85
}
```

---

## 配置

### 环境变量

#### 后端

```bash
# 服务器端口（默认: 8080）
PORT=8080

# 日志级别（默认: info）
LOG_LEVEL=debug

# CORS 源（默认: *）
CORS_ORIGINS=http://localhost:3000,https://agentguard.example.com
```

#### 前端

```bash
# API 基础 URL
VITE_API_URL=http://localhost:8080

# 指标刷新间隔（毫秒）
VITE_METRICS_REFRESH=30000
```

### 策略配置

创建 `.agent-guard.yaml`：

```yaml
# 阻止的命令
blocked_commands:
  - "rm -rf /"
  - "dd if=/dev/zero"
  - "mkfs"

# 允许的路径
allowed_paths:
  - /tmp
  - /home/user/project

# 拒绝的路径
denied_paths:
  - /etc/passwd
  - /etc/shadow
  - ~/.ssh

# 网络访问
network:
  allowed_domains:
    - api.github.com
    - cdn.jsdelivr.net
  denied_domains:
    - "*.malicious.com"
```

---

## 故障排除

### 后端问题

#### 端口已被占用

```bash
# 检查端口使用情况
lsof -ti:8080

# 终止进程
kill -9 $(lsof -ti:8080)

# 或使用不同端口
PORT=9090 go run webui/backend/main.go
```

#### 缺少依赖

```bash
# 清理并重新下载
go mod tidy
go mod download

# 验证依赖
go mod verify
```

#### 连接被拒绝

```bash
# 检查后端是否运行
ps aux | grep "go run"

# 检查日志
tail -f /tmp/backend-test.log

# 直接测试 API
curl -v http://localhost:8080/api/v1/status
```

### 前端问题

#### 构建失败

```bash
# 清理构建缓存
rm -rf node_modules dist
npm install
npm run build
```

#### API 连接错误

```bash
# 检查后端是否可访问
curl http://localhost:8080/api/v1/status

# 验证 VITE_API_URL
echo $VITE_API_URL

# 检查浏览器控制台的 CORS 错误
```

#### 代理问题

```bash
# 为本地主机绕过代理
export NO_PROXY="localhost,127.0.0.1"

# 或临时禁用代理
unset http_proxy https_proxy
```

### Docker 问题

#### 容器无法启动

```bash
# 查看日志
docker compose logs backend
docker compose logs frontend

# 无缓存重建
docker compose build --no-cache

# 检查容器状态
docker compose ps
```

#### 网络问题

```bash
# 检查网络
docker network inspect agentguard-network

# 测试连接性
docker compose exec backend ping frontend

# 重建网络
docker compose down
docker network prune
docker compose up -d
```

### 监控问题

#### Prometheus 未抓取

```bash
# 检查指标端点
curl http://localhost:8080/metrics | head

# 验证 Prometheus 配置
promtool check config prometheus.yml

# 检查 Prometheus 日志
docker compose logs prometheus
```

#### 无指标数据

```bash
# 验证扫描正在运行
curl http://localhost:8080/api/v1/scan | jq .

# 检查指标 API
curl http://localhost:8080/api/v1/metrics/scan-rate | jq .

# 查看 Prometheus 目标
# http://localhost:9090/targets
```

---

## 安全考虑

### 生产环境检查清单

- [ ] 为所有端点启用 HTTPS/TLS
- [ ] 限制 `/metrics` 端点仅限内部网络
- [ ] 设置强 CORS 策略
- [ ] 为 API 端点启用身份验证
- [ ] 定期安全更新
- [ ] 启用审计日志
- [ ] 配置速率限制
- [ ] 验证所有端点的输入
- [ ] 密钥管理（API 密钥、令牌）
- [ ] 定期漏洞扫描

### 加固指南

#### 后端安全

```go
// 添加身份验证中间件
router.Use(authMiddleware())

// 速率限制
router.Use(rateLimitMiddleware())

// 输入验证
router.Use(validationMiddleware())
```

#### Nginx 安全

```nginx
# 安全头
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;

# 隐藏版本
server_tokens off;

# 速率限制
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
```

### 防火墙规则

```bash
# 仅允许必要的端口
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 9090/tcp   # Prometheus（仅内部）
ufw enable
```

---

## 性能调优

### 后端优化

```bash
# 增加 Go 垃圾回收目标
export GOGC=100

# 启用 Go 性能分析
export ENABLE_PPROF=true

# 设置工作线程
export GOMAXPROCS=4
```

### 前端优化

```bash
# 使用优化构建
npm run build -- --mode production

# 启用压缩
gzip on;
gzip_types text/plain application/json application/javascript text/css;
```

### 数据库优化（未来）

- 连接池
- 查询优化
- 索引策略
- 缓存层

---

## 备份和恢复

### 配置备份

```bash
# 备份策略配置
cp .agent-guard.yaml .agent-guard.yaml.backup

# 备份 Prometheus 配置
cp prometheus.yml prometheus.yml.backup

# 备份 Grafana 仪表板
curl http://localhost:3000/api/dashboards/export > dashboards.json
```

### 数据备份

```bash
# 备份扫描历史
tar -czf scan-history-$(date +%Y%m%d).tar.gz ~/.agent-guard/

# 自动化备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d)
tar -czf /backup/agentguard-$DATE.tar.gz \
  ~/.agent-guard/ \
  .agent-guard.yaml \
  prometheus.yml
```

### 灾难恢复

```bash
# 从备份恢复
tar -xzf agentguard-20260228.tar.gz -C /

# 验证配置
go run webui/backend/main.go --config .agent-guard.yaml.backup
```

---

## 维护

### 定期任务

**每日**：
- 检查应用日志错误
- 监控指标仪表板
- 查看安全告警

**每周**：
- 查看扫描结果
- 更新漏洞数据库
- 检查磁盘空间使用

**每月**：
- 安全更新
- 性能审查
- 备份验证
- 依赖更新

### 日志轮转

创建 `/etc/logrotate.d/agentguard`：

```
/var/log/agentguard/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 agentguard agentguard
    sharedscripts
    postrotate
        systemctl reload agentguard-backend
    endscript
}
```

---

## 支持和资源

### 文档

- [主 README](../README.md)
- [监控指南](../doc/MONITORING.md)
- [使用指南](../USAGE.md)
- [路线图](../doc/ROADMAP.md)

### 获取帮助

- GitHub 问题: https://github.com/imdlan/AIAgentGuard/issues
- 文档: https://github.com/imdlan/AIAgentGuard/tree/main/doc
- 讨论: https://github.com/imdlan/AIAgentGuard/discussions

### 贡献

1. Fork 仓库
2. 创建功能分支
3. 进行更改
4. 提交 Pull Request

---

## 附录

### API 端点参考

| 端点 | 方法 | 描述 |
|----------|--------|-------------|
| `/api/v1/status` | GET | 系统状态 |
| `/api/v1/scan` | GET/POST | 执行扫描 |
| `/api/v1/scan/:id` | GET | 获取扫描结果 |
| `/api/v1/history` | GET | 扫描历史 |
| `/api/v1/trends` | GET | 趋势数据 |
| `/api/v1/alerts` | GET | 安全告警 |
| `/api/v1/metrics` | GET | 指标信息 |
| `/api/v1/metrics/scan-rate` | GET | 扫描统计 |
| `/api/v1/metrics/vulnerabilities` | GET | 漏洞计数 |
| `/api/v1/metrics/duration` | GET | 持续时间指标 |
| `/metrics` | GET | Prometheus 指标 |

### 默认端口

| 服务 | 端口 | 协议 |
|---------|------|----------|
| 后端 API | 8080 | HTTP |
| 前端 | 3000 | HTTP |
| Prometheus | 9090 | HTTP |
| Grafana | 3001 | HTTP |

### 文件位置

| 文件 | 位置 |
|------|----------|
| 后端代码 | `webui/backend/main.go` |
| 前端代码 | `webui/frontend/src/` |
| 策略配置 | `.agent-guard.yaml` |
| Prometheus 配置 | `configs/prometheus.yml` |
| Grafana 仪表板 | `configs/grafana-dashboard.json` |
| 日志 | `/var/log/agentguard/` 或 `~/.agent-guard/` |

---

**最后更新**: 2026-02-28
**版本**: v1.2.0
**维护者**: AIAgentGuard 团队
