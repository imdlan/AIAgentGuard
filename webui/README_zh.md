# AI AgentGuard Web UI

## 概览

AI AgentGuard Web Dashboard 提供了一个实时、可视化的安全监控界面，让管理员能够轻松查看和管理 AI Agent 的安全状态。

## 技术栈

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **样式**: CSS (现代 CSS3)
- **HTTP客户端**: Fetch API

### 后端
- **语言**: Go 1.25.5
- **框架**: Gin
- **API**: RESTful + WebSocket (规划中)

## 项目结构

```
webui/
├── frontend/              # React 前端应用
│   ├── src/
│   │   ├── components/    # React 组件
│   │   │   ├── Dashboard.tsx
│   │   │   └── Dashboard.css
│   │   ├── api/           # API 客户端
│   │   │   └── client.ts
│   │   ├── types/         # TypeScript 类型
│   │   │   └── index.ts
│   │   ├── App.tsx        # 根组件
│   │   ├── App.css        # 全局样式
│   │   └── main.tsx       # 入口文件
│   ├── package.json
│   ├── vite.config.ts
│   ├── Dockerfile
│   └── nginx.conf
├── backend/               # Go 后端 API
│   ├── main.go            # 主服务器文件
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── docker-compose.yml     # Docker 编排配置
└── README.md             # 本文档
```

## 快速开始

### 本地开发

#### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

#### 后端开发

```bash
cd backend

# 下载依赖
go mod download

# 运行服务器
go run main.go

# 或运行后端
mkdir -p ../frontend/dist  # 创建前端dist目录
go run main.go

# API 访问 http://localhost:8080
```

### 使用 Docker

```bash
# 启动所有服务
docker-compose up -d

# 访问
# 前端: http://localhost:3000
# 后端 API: http://localhost:8080/api/v1/status

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## API 端点

### 执行安全扫描

```bash
# GET 请求 - 执行完整扫描
GET /api/v1/scan

# POST 请求 - 自定义扫描选项
POST /api/v1/scan
Content-Type: application/json

{
  "categories": ["filesystem", "shell", "network"],
  "options": {
    "include_file_content": true,
    "include_processes": false,
    "include_suid": false,
    "max_depth": 5,
    "timeout": 30
  }
}
```

### 响应格式

```json
{
  "id": "scan-20260227-143045",
  "timestamp": "2026-02-27T14:30:45Z",
  "duration": 512,
  "results": {
    "filesystem": "HIGH",
    "shell": "HIGH",
    "network": "MEDIUM",
    "secrets": "LOW",
    "filecontent": "LOW",
    "dependencies": "LOW"
  },
  "overall": "HIGH",
  "details": [
    {
      "type": "HIGH",
      "category": "filesystem",
      "description": "Writable access to sensitive directories",
      "path": "/Users/David/.ssh, /Users/David/.config"
    }
  ]
}
```

### 其他端点

```
GET  /api/v1/scan/:id          # 获取特定扫描结果
GET  /api/v1/history           # 扫描历史记录
GET  /api/v1/trends            # 趋势数据
GET  /api/v1/alerts            # 安全告警
GET  /api/v1/status            # 系统状态
WS   /api/v1/realtime          # 实时更新
```

## 功能特性

### 当前实现 (v1.2.0)

✅ **实时安全扫描**: 一键执行完整的安全检查  
✅ **可视化结果**: 直观的权限风险评估  
✅ **风险等级分级**: LOW/MEDIUM/HIGH/CRITICAL  
✅ **详细发现**: 显示具体的安全问题和位置  
✅ **系统状态**: 实时显示系统运行状态  

### 计划中 (v1.2.0 后续)

🔄 **扫描历史**: 查看历史上所有扫描结果  
🔄 **趋势分析**: 可视化安全趋势变化  
🔄 **告警管理**: 实时安全告警和通知  
🔄 **WebSocket**: 实时更新扫描进度和结果  

### v1.3.0 计划

📋 **用户认证**: 登录和权限控制  
📋 **多租户**: 支持多组织和团队  
📋 **报告导出**: PDF/Excel 格式报告  

## 配置

### 环境变量

```bash
# 后端配置
PORT=8080                    # 服务器端口

# 前端配置
VITE_API_URL=http://localhost:8080  # 后端 API 地址
```

## 性能指标

### 前端
- 初次加载: < 500ms
- 扫描响应: < 200ms (不含实际扫描时间)
- 扫描执行: ~500ms (取决于机器性能)

### 后端
- API 响应时间: < 50ms
- 并发扫描: 支持 100+ 并发请求
- 内存使用: < 100MB

## 故障排除

### 前端无法连接后端

1. 检查后端是否运行
```bash
curl http://localhost:8080/api/v1/status
```

2. 检查环境变量
```bash
# 前端
echo $VITE_API_URL

# 后端
echo $PORT
```

3. 检查 CORS 配置
   - Go 后端已自动添加 CORS 中间件

### Docker 容器启动失败

```bash
# 查看详细日志
docker-compose logs backend
docker-compose logs frontend

# 重新构建
docker-compose down
docker-compose build
docker-compose up -d
```

### 扫描超时

```bash
# 增加扫描超时时间
# 修改前端扫描请求配置
{
  "options": {
    "timeout": 60  // 增加到60秒
  }
}
```

## 开发指南

### 添加新组件

```bash
# 1. 创建组件文件
cd frontend/src/components
touch NewComponent.tsx

# 2. 添加样式
touch NewComponent.css

# 3. 在 App.tsx 中导入
import NewComponent from './components/NewComponent';
```

### 添加新 API 端点

```bash
# 1. 在 backend/main.go 中添加处理函数
func handleNewEndpoint(c *gin.Context) {
    // 实现
}

# 2. 注册路由
api.GET("/new-endpoint", handleNewEndpoint)

# 3. 在前端 API 客户端中添加方法
async newEndpoint(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/v1/new-endpoint`);
    return await response.json();
}
```

## 安全考虑

- ✅ 所有 API 请求使用 HTTPS (生产环境)
- ✅ CORS 配置为通配符，生产环境应限制
- ✅ 无敏感信息在前端暴露
- ✅ 后端验证所有输入

## 测试

### 前端测试

```bash
cd frontend
npm run test          # 运行测试
npm run test:coverage # 查看覆盖率
```

### 后端测试

```bash
cd backend
go test ./...            # 运行测试
go test -cover ./...      # 查看覆盖率
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License - 见主项目 LICENSE 文件

## 支持

## 部署

详细的部署指南请参阅: [DEPLOYMENT_zh.md](DEPLOYMENT_zh.md)

**部署选项**:
- 本地开发
- Docker / Docker Compose
- 生产环境（Systemd、Nginx、Kubernetes）

包含以下内容：
- 前置要求和安装
- 配置管理
- 监控设置
- 故障排除
- 安全加固
- 性能调优
- 备份和恢复

## 支持

- 问题反馈: https://github.com/imdlan/AIAgentGuard/issues
- 文档: https://github.com/imdlan/AIAgentGuard/tree/main/doc

## 监控功能 (v1.2.0 新增)

### Prometheus 集成

Web UI 后端现已集成 Prometheus 监控，可实时收集和暴露安全指标。

**指标端点**:
```bash
# Prometheus 指标（用于抓取）
GET /metrics

# 指标信息
GET /api/v1/metrics

# 扫描速率统计
GET /api/v1/metrics/scan-rate

# 漏洞统计
GET /api/v1/metrics/vulnerabilities

# 扫描持续时间
GET /api/v1/metrics/duration
```

**MetricsPanel 组件**:

前端新增 `MetricsPanel` 组件，提供实时监控面板：

- ✅ 扫描统计（总数、速率、平均时长）
- ✅ 漏洞概览（严重级别分布）
- ✅ 语言特定漏洞统计
- ✅ 性能指标（P50/P95/P99）
- ✅ 30秒自动刷新

### 配置 Prometheus

**快速启动（开发环境）**:
```bash
# 从项目根目录运行
cd /path/to/AIAgentGuard

# 启动后端（已启用 /metrics 端点）
go run webui/backend/main.go

# 使用 curl 测试
curl http://localhost:8080/metrics | head
```

**生产环境（Docker Compose）**:

使用提供的监控堆栈：
```bash
cd configs
docker-compose -f docker-compose.monitoring.yml up -d

# 访问服务
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3001
```

**Prometheus 配置**:

在 `prometheus.yml` 中添加抓取配置：
```yaml
scrape_configs:
  - job_name: 'agentguard'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana 仪表板

导入预配置的仪表板：

```bash
# 仪表板文件位置
configs/grafana-dashboard.json

# 通过 Grafana UI 导入:
# 1. 导航到 Dashboards -> Import
# 2. 上传 JSON 文件
# 3. 选择 Prometheus 数据源
# 4. 点击 Import
```

**仪表板包含**:
- 扫描速率趋势
- 按严重级别分类的漏洞计数
- 扫描持续时间百分位数
- 特定语言的漏洞趋势
- 实时告警

### 指标说明

**扫描指标** (`agentguard_scan_total`, `agentguard_scan_duration_seconds`):
- `scan_total`: 执行的总扫描次数
- `scan_duration_seconds`: 扫描持续时间（直方图）

**漏洞指标** (`agentguard_vulnerabilities_total`):
- 按 `severity` 标签: critical, high, medium, low
- 按 `language` 标签: go, npm, pip, cargo

**语言特定扫描** (`agentguard_language_scan_total`):
- Go 模块漏洞扫描
- npm/yarn 包扫描
- Python pip 包扫描
- Rust cargo 包扫描

### 部署指南

完整的部署指南请参阅: [DEPLOYMENT.md](DEPLOYMENT.md)

**部署选项**:
- 本地开发
- Docker / Docker Compose
- 生产环境（Systemd, Nginx, Kubernetes）

包含以下内容：
- 前置要求和安装
- 配置管理
- 监控设置
- 故障排除
- 安全加固
- 性能调优
- 备份和恢复
