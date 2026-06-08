# Changelog

## [Unreleased]

### Changed
- **重构目录结构**：`tools/` 拆分为 `docker/`、`gitlab/`、`jenkins/`、`kube/` 平铺在根目录，更贴近 zhilo/ecommerce 风格
- **CLI 框架迁移**：`github.com/tiamxu/kit/cli` → `github.com/spf13/cobra`（解决本地 kit 模块 cli 包缺失问题）
- **client.go 拆分**：`client/client.go`（200行）拆为 `client.go` + `jenkins.go` + `gitlab.go`
- **DTO 解耦**：新增 `types/` 包，请求/响应结构体从 `client/` 抽离
- **API 客户端新增方法**：`ListJenkinsProjects` 用于批量创建时拉取真实项目列表
- **leister-api 新增接口**：`GET /api/jenkins/projects?group=xxx` 供 CLI 批量操作使用
- **leister-api 改造**：`JenkinsService` 注入 `GitlabService`，批量项目数据源从 GitLab API 获取

### Security
- **修复 Docker 凭证硬编码**：`RegistryPassword` 从硬编码改为读取 `DOCKER_REGISTRY_PASSWORD` 环境变量，缺失即报错

### Error Code Refactor
- **错误码迁到 leister-api**：按 zhilo 风格（5 位分段 + map[int]string 查表）实现 `leister-api/pkg/e/`
  - 段位：101xx 通用 / 201xx Jenkins / 301xx GitLab / 999xx 内部
  - 引入 `AppError` 类型 + `Unwrap()`，service 层返回业务错误，handler 层统一解析
  - 统一响应结构 `{code, msg, data}`，原始 err 不外泄到前端
  - handler 全部走 `e.OK` / `e.JSON` helper，禁止手写 `c.JSON(500, gin.H{"error": err.Error()})`
- **CLI 端同步改造**：`client.go` 解析 `apiEnvelope{Code, Msg, Data}`，业务失败透传 `[code] msg`

### Fixed
- **修复 `jks cts` 假批量 bug**：原代码硬编码 3 个 project1/2/3，现改为从 API 拉取真实项目列表
- **修复 `BaseURL` 默认值导致静默失败**：未设置 `LEISTER_API_URL` 时改为启动失败并清晰提示，不再默认连 localhost
- **修复 leister 编译失败**：本地 `../kit` HEAD 缺少 `cli` 包，通过切换到 `cobra` 彻底解决

### Removed
- **删除空文件** `tools/deploy/deploy.go` 及其目录
- **删除 `.gitignore` 中的 `cmd` 行**（不再使用 cmd 目录）
- **删除 `pkg/e/`**：CLI 端不需要错误码体系（错误信息由 API 透传），死代码清理
- **删除 `G:\go\src\go.work`**：改用 `go.mod` 的 `replace github.com/tiamxu/kit => ../kit` 替代 workspace

### Added
- **新增 `Makefile`**：build / run / test / clean 标准化命令
- **client 解析新响应结构**：CLI 端解析 `{code, msg, data}` 三段式，code != 200 视为业务错误并打印 `[code] msg`
- **client 修复 baseURL 尾斜杠**：`NewClient` 用 `strings.TrimRight` 去掉尾部 `/`，避免 `http://x/api//api/...` 拼接错误
