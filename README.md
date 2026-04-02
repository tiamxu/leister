# Leister

Leister 是一个强大的 DevOps 工具集，用于简化 Docker 镜像构建、Jenkins 任务管理、GitLab 项目管理和 Kubernetes 部署操作。

## 功能特性

### 核心功能

1. **Docker 镜像管理**
   - 代码构建（支持 Go 和 Node.js）
   - Docker 镜像构建
   - 镜像推送至 Harbor/Registry

2. **Jenkins 任务管理**
   - 创建单个 Jenkins 任务
   - 批量创建 Jenkins 任务
   - 更新 Jenkins 任务配置

3. **GitLab 项目管理**
   - 获取 GitLab 项目信息
   - 生成 GitLab 项目数据

4. **Kubernetes 部署管理**
   - 获取部署信息
   - 重启部署
   - 创建部署

## 项目结构

```
leister/
├── main.go              # 主入口文件
├── config/              # 配置管理
│   └── config.go        # 配置结构定义（使用环境变量）
├── client/              # API 客户端
│   └── client.go        # HTTP API 客户端
├── tools/               # 工具模块
│   ├── docker/          # Docker 镜像管理
│   ├── gitlab/          # GitLab 项目管理
│   ├── jenkins/         # Jenkins 任务管理
│   └── kube/            # Kubernetes 部署管理
└── go.mod               # 依赖管理
```

## 安装

### 构建项目

```bash
go build -o gigctl
```

### 安装到系统路径

```bash
# Linux/Mac
cp gigctl /usr/local/bin/
chmod +x /usr/local/bin/gigctl

# Windows
copy gigctl.exe C:\Windows\System32\
```

## 配置

Leister 使用**环境变量**进行配置，无需配置文件。

### 环境变量

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `LEISTER_API_URL` | API 服务地址 | `http://localhost:8080` |
| `LEISTER_API_TIMEOUT` | API 超时时间（秒） | `30` |

### 配置示例

```bash
# 使用默认配置
./gigctl docker build

# 设置 API 地址
export LEISTER_API_URL=http://api.example.com:8080
./gigctl jks create -n myjob -g mygroup

# 临时指定环境变量
LEISTER_API_URL=http://api.example.com:8080 LEISTER_API_TIMEOUT=60 ./gigctl git get -n myproject -g mygroup
```

## 使用指南

### 查看命令帮助

```bash
gigctl help
```

### Docker 命令

Docker 命令直接执行本地 Docker 操作，不需要连接 API 服务。

#### 构建 Docker 镜像

```bash
# 基本构建
gigctl docker build

# 指定版本号和环境
gigctl docker build -t v1.0.0 -e prod

# 指定 Dockerfile
gigctl docker build -t v1.0.0 -f Dockerfile.prod
```

参数说明：
- `-t, --tag`：镜像标签（默认使用时间戳格式 202502031112）
- `-e, --env`：环境名称（默认 dev）
- `-l, --lang`：代码语言（默认 go）
- `-f, --dockerfile`：Dockerfile 路径（默认 Dockerfile-dev）

#### 推送 Docker 镜像

```bash
# 基本推送
gigctl docker push

# 指定版本号和环境
gigctl docker push -t v1.0.0 -e prod
```

### GitLab 命令

GitLab 命令通过 API 调用 leister-api 服务。

#### 获取 GitLab 项目信息

```bash
gigctl git get -n project-name -g project-group
```

参数说明：
- `-n, --name`：项目名称（必需）
- `-g, --group`：项目组名称（必需）

#### 生成 GitLab 项目数据

```bash
gigctl git gen -g project-group
```

参数说明：
- `-g, --group`：项目组名称（必需）

### Jenkins 命令

Jenkins 命令通过 API 调用 leister-api 服务。

#### 创建单个 Jenkins 任务

```bash
gigctl jks create -n job-name -g job-group
```

参数说明：
- `-n, --name`：任务名称（必需）
- `-g, --group`：任务组名称（必需）

#### 批量创建 Jenkins 任务

```bash
gigctl jks cts -g project-group
```

参数说明：
- `-g, --group`：项目组名称（可选）

#### 更新 Jenkins 任务

```bash
gigctl jks update -n job-name -g job-group
```

参数说明：
- `-n, --name`：任务名称（必需）
- `-g, --group`：任务组名称（必需）

### Kubernetes 命令

Kubernetes 命令直接执行本地 kubectl 操作，不需要连接 API 服务。

#### 获取部署信息

```bash
# 获取所有部署
gigctl kube get

# 获取指定部署
gigctl kube get -n deployment-name

# 指定命名空间
gigctl kube get -n deployment-name --namespace production
```

#### 重启部署

```bash
gigctl kube restart -n deployment-name
```

参数说明：
- `-n, --name`：部署名称（必需）
- `--namespace`：命名空间（默认 default）

#### 创建部署

```bash
gigctl kube create -n deployment-name --image nginx:latest --replicas 3
```

参数说明：
- `-n, --name`：部署名称（必需）
- `--image`：容器镜像（必需）
- `--replicas`：副本数（默认 1）
- `--namespace`：命名空间（默认 default）

## 技术栈

- **开发语言**：Go 1.25+
- **CLI 框架**：kit/cli（基于 cobra）
- **日志库**：kit/log
- **HTTP 客户端**：标准库 net/http
- **依赖库**：
  - github.com/tiamxu/kit - 工具库

## 架构说明

### CLI 与 API 分离

Leister 采用 CLI + API 的分离架构：

- **CLI 层**（leister）：负责命令行交互、参数解析、结果展示
- **API 层**（leister-api）：负责业务逻辑、状态管理、外部服务集成

### 命令分类

| 命令类型 | 执行方式 | 说明 |
|---------|---------|------|
| Docker | 本地执行 | 直接调用本地 docker 命令 |
| Kubernetes | 本地执行 | 直接调用本地 kubectl 命令 |
| Jenkins | API 调用 | 通过 HTTP 调用 leister-api 服务 |
| GitLab | API 调用 | 通过 HTTP 调用 leister-api 服务 |

## 特性说明

### 智能标签生成

当未指定镜像标签时，系统会自动生成时间戳格式的标签（如 202502031112）。

### 登录优化

Docker 工具会智能判断是否已登录到 Registry，避免重复登录，提高执行效率。

### 代码语言支持

- **Go**：自动执行 `go build` 命令
- **Node.js**：自动执行 `npm install` 和 `npm run build` 命令

## 开发

### 环境要求

- Go 1.25 或更高版本
- Docker（用于构建和推送镜像）
- kubectl（用于 Kubernetes 命令）
- leister-api 服务（用于 Jenkins 和 GitLab 命令）

### 项目开发

1. 克隆项目
2. 安装依赖：`go mod tidy`
3. 构建项目：`go build -o gigctl`
4. 运行测试：`go test ./...`

## 与 leister-api 配合使用

Leister CLI 需要与 leister-api 服务配合使用：

```bash
# 1. 启动 leister-api 服务
cd ../leister-api
./leister-api

# 2. 在另一个终端使用 leister CLI
cd ../leister
./gigctl docker build
./gigctl jks create -n myjob -g mygroup
```

## 贡献指南

欢迎贡献代码、报告问题和提出建议！

## 许可证

本项目采用 MIT 许可证。
