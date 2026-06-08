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

### Fixed
- **修复 `jks cts` 假批量 bug**：原代码硬编码 3 个 project1/2/3，现改为从 API 拉取真实项目列表
- **修复 `BaseURL` 默认值导致静默失败**：未设置 `LEISTER_API_URL` 时改为启动失败并清晰提示，不再默认连 localhost
- **修复 leister 编译失败**：本地 `../kit` HEAD 缺少 `cli` 包，通过切换到 `cobra` 彻底解决

### Removed
- **删除空文件** `tools/deploy/deploy.go` 及其目录
- **删除 `.gitignore` 中的 `cmd` 行**（不再使用 cmd 目录）
- **删除 `go.mod` 中的 `replace github.com/tiamxu/kit => ../kit`**（go.work 已接管）

### Added
- **新增 `pkg/e/` 统一错误码**：包含 Jenkins/GitLab/Docker/Kube 错误码常量及中文错误信息查询
- **新增 `Makefile`**：build / run / test / clean 标准化命令
- **新增 `G:\go\src\go.work`**：kit / leister / leister-api 三模块共享 workspace
