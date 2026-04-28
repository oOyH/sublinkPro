# 2026-04-28 本地验证记录

使用 Go 路径 `D:\Program Files\Go\bin\go.exe`。

识别到版本为 `go1.26.2 windows/amd64`。

为避免默认 GOPATH 无权限，测试时将 `GOCACHE`、`GOMODCACHE`、`GOPATH` 重定向到仓库下 `.codex/` 目录。

执行命令：

`go test ./...`

结果：

大部分包通过。

`sublink/api` 失败。

失败用例：

`TestGetClientConcurrentClashRequestsKeepSubscriptionScoped`

`TestGetClientConcurrentSurgeRequestsKeepSubscriptionScoped`

失败现象：

响应体返回 `配置读取错误`，说明 `/c` 订阅导出链路在并发场景下存在回归。

---

## 2026-04-28 修复后复测

执行命令：

`go test ./services/... ./node/... ./models/...`

结果：

通过。

执行命令：

`go test ./...`

结果：

最初仍失败，失败点集中在 `sublink/api` 包的并发测试：

`TestGetClientConcurrentClashRequestsKeepSubscriptionScoped`

`TestGetClientConcurrentSurgeRequestsKeepSubscriptionScoped`

继续排查后确认不是本机环境问题，而是测试夹具代码问题。`api/clients_test.go` 手工拼接订阅 `Config` JSON，在 Windows 路径下把反斜杠直接写进 JSON 字符串，导致 `sub.Config` 不是合法 JSON，`GetClash` 和 `GetSurge` 在 `json.Unmarshal` 时返回 `配置读取错误`。

修复后改为使用 `json.Marshal` 生成测试订阅配置。

复测：

`go test ./api/...` 通过。

`go test ./...` 通过。
