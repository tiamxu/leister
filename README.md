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
   - 批量创建 Jenkins 任务（从数据库读取）
   - 更新 Jenkins 任务配置

3. **GitLab 项目管理**
   - 获取 GitLab 项目信息
   - 将 GitLab 组项目同步到数据库

4. **Kubernetes 部署管理**
   - 获取部署信息
   - 重启部署
   - 创建部署

## 项目结构

```
leister/
├── main.go              # 主入口文件
├── config/              # 配置管理
│   ├── config.go        # 配置结构定义
│   └── config.yaml      # 配置文件示例
├── database/            # 数据库操作
│   ├── config.go        # 数据库配置
│   ├── init.sql         # 数据库初始化脚本
│   ├── item.go          # 项目数据模型
│   └── mysql.go         # MySQL 数据库操作
└── tools/               # 工具模块
    ├── build/           # 构建工具（兼容旧版）
    ├── deploy/          # 部署工具（兼容旧版）
    ├── docker/          # Docker 镜像管理
    ├── gitlab/          # GitLab 项目管理
    ├── jenkins/         # Jenkins 任务管理
    └── kube/            # Kubernetes 部署管理
```

## 安装

### 构建项目

```bash
go build -o bin/gigctl
cp bin/gigctl /usr/local/bin/gigctl
source ~/.bashrc
```

## 配置

Leister 使用 `config.yaml` 配置文件，配置文件示例：

```yaml
db:
  driver: mysql
  host: localhost
  port: 3306
  database: leister
  username: root
  password: password
  max_idle_conns: 10
  max_open_conns: 20
  conn_max_lifetime: 300

gitlab:
  url: https://gitlab.example.com
  token: your-gitlab-token

jenkins:
  url: https://jenkins.example.com
  username: admin
  password: admin-password
```

## 使用指南

### 查看命令帮助

```bash
gigctl help
```

### Docker 命令

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

#### 获取 GitLab 项目信息

```bash
gigctl git get -n project-name -g project-group
```

参数说明：
- `-n, --name`：项目名称（必需）
- `-g, --group`：项目组名称（必需）

#### 将 GitLab 组项目同步到数据库

```bash
gigctl git gen -g project-group
```

参数说明：
- `-g, --group`：项目组名称（必需）

### Jenkins 命令

#### 创建单个 Jenkins 任务

```bash
gigctl jks create -n job-name -g job-group
```

参数说明：
- `-n, --name`：任务名称（必需）
- `-g, --group`：任务组名称（必需）

#### 从数据库读取项目并创建任务

```bash
# 为所有项目创建任务
gigctl jks cts

# 为指定组的项目创建任务
gigctl jks cts -g project-group
```

参数说明：
- `-g, --group`：项目组名称（可选）

### Kubernetes 命令

#### 获取部署信息

```bash
gigctl kube get -n deployment-name
```

#### 重启部署

```bash
gigctl kube restart -n deployment-name
```

#### 创建部署

```bash
gigctl kube create -n deployment-name
```

## 技术栈

- **开发语言**：Go 1.19+
- **CLI 框架**：kit/cli（基于 cobra）
- **日志库**：kit/log
- **数据库**：MySQL
- **依赖库**：
  - github.com/bndr/gojenkins - Jenkins 客户端
  - github.com/xanzy/go-gitlab - GitLab 客户端
  - k8s.io/client-go - Kubernetes 客户端
  - github.com/koding/multiconfig - 配置管理

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

- Go 1.19 或更高版本
- Docker（用于构建和推送镜像）
- GitLab 访问权限（用于 GitLab 命令）
- Jenkins 访问权限（用于 Jenkins 命令）
- Kubernetes 配置（用于 Kubernetes 命令）

### 项目开发

1. 克隆项目
2. 配置 config.yaml
3. 构建项目：`go build -o bin/gigctl`
4. 运行测试：`go test ./...`

## 贡献指南

欢迎贡献代码、报告问题和提出建议！

## 许可证

本项目采用 MIT 许可证。
