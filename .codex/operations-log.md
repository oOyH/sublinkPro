# 2026-04-28 分析记录

已完成仓库结构扫描，确认项目为 Go 后端加 React 前端的单仓库架构。

已梳理自动更新订阅链路，核心路径为 `services/scheduler/subscription_task.go` 到 `node/sub.go` 再到 `models/node.go`。

已确认任务中心前端只给 `speed_test` 显示停止按钮，`sub_update` 没有终止入口。

已确认订阅更新任务没有真正接入取消检查，后端 context 被创建但未沿执行链传递。

已确认 `TaskManager` 在任务终态落库失败时会继续移除内存任务，这会留下数据库 running 脏记录。

已确认 scheduler 对同机场订阅更新没有重入保护，存在重叠执行风险。

已使用本机 `D:\Program Files\Go\bin\go.exe` 做本地测试验证。

`go test ./...` 已执行完成，大部分包通过，`sublink/api` 包存在并发测试失败。

---

已完成修复落地：

`services/task_manager.go` 已补终态写库重试、终态补偿同步、孤儿 running 任务自动收敛，并避免终态写库失败后直接丢失内存任务。

`services/scheduler/subscription_task.go` 已补同机场订阅更新互斥保护，并把任务 context 接入订阅更新执行链。

`node/sub.go` 已补订阅下载请求的 context 取消、解析阶段取消检查、节点同步阶段取消检查。

`webs/src/components/TaskProgressPanel.jsx` 与 `webs/src/views/tasks/index.jsx` 已放开 `sub_update` 的停止入口。

已完成后端相关复测：

`go test ./services/... ./node/... ./models/...` 通过。

随后继续排查 `sublink/api` 并发测试失败。

确认根因不是本机环境，而是 `api/clients_test.go` 的测试夹具手工拼接 Windows 模板路径 JSON，导致 `sub.Config` 非法。

已改为用 `json.Marshal` 生成测试订阅配置。

复测结果：

`go test ./api/...` 通过。

`go test ./...` 通过。

2026-04-28 前端 UI 修复：已调整 webs/src/components/NodeNamePreprocessor.jsx 的规则行布局，改为可换行弹性布局，移除导致按钮被裁剪的固定最小宽度。
2026-04-28 验证：尝试执行前端 build，但 webs/node_modules 不存在，vite 命令不可用，当前无法在本地直接完成构建验证。

## 2026-04-28 安全审查
- 目标：审查新增订阅后是否存在可被他人盗取的风险。
- 结论：未发现匿名越权读取任意订阅的直接漏洞；主要风险来自自动创建默认永久分享链接、/c/ 仅凭 token 访问、允许自定义弱 token、缺少订阅访问限速。
- 证据文件：api/sub.go, models/subscription_share.go, api/clients.go, routers/clients.go, api/share.go, middlewares/ip.go, models/iplogs.go, webs/src/views/subscriptions/component/ShareManageDialog.jsx, services/telegram/handlers.go, middlewares/jwt.go, routers/subcription.go.

## 2026-04-28 订阅新增后的盗取风险审查补充

- 审查目标：分析新增订阅后是否存在被他人盗取订阅内容的风险。
- 结论：未发现匿名越权读取任意订阅的直接漏洞，但存在明显的订阅被盗用风险，核心原因是系统会在新增订阅后自动创建默认公开分享入口，且该入口默认启用、永不过期，只依赖 URL 中的 token 进行访问控制。
- 主要风险点：
  - 新增订阅时自动创建默认分享链接，默认配置为 `Name=默认分享链接`、`ExpireType=永不过期`、`IsLegacy=true`、`Enabled=true`。这意味着订阅创建后立即存在公开访问入口。
  - 公开访问入口为 `/c/?token=...`，该路由不走后台 JWT 鉴权，而是仅凭 token 访问订阅内容。任何拿到链接的人都可直接拉取订阅。
  - 前端分享弹窗会直接拼接并展示完整分享链接，可复制、可生成二维码，增加了截图、剪贴板、历史记录和外部转发导致 token 泄漏的风险。
  - Telegram 机器人会直接构造默认分享链接，若机器人消息被转发或环境不安全，token 也可能外泄。
  - 分享 token 放在 URL 查询参数中，存在被浏览器历史、反向代理日志、第三方监控、Referer 链路或应用日志间接带出的风险。
  - 当前允许管理员自定义 token，但后端仅检查唯一性，不强制长度和复杂度。如果被设置成短 token 或可猜 token，会进一步放大爆破和撞取风险。
  - 当前未看到 `/c/` 订阅访问链路有专门的失败限速、爆破惩罚或异常访问拦截机制。
  - `api/clients.go` 在无效 token 和过期 token 场景会写告警日志，并带出原始 token 字符串，属于额外泄漏面。
- 已确认的现有防线：
  - 后台订阅管理接口和分享管理接口本身走 `AuthToken`，不是匿名开放。
  - 默认随机 token 使用 `crypto/rand` 生成 16 字节随机值再 hex 编码，默认自动生成场景下 token 本身不属于弱随机。
  - 订阅链路支持 IP 黑白名单，可对公开访问做一定限制。
  - 访问日志表记录的是 IP、时间、地址、订阅 ID 和 ShareID，不直接持久化 token。
- 风险性质判定：这更像“默认公开分享设计风险”与“token 泄漏后可直接访问”的组合问题，而不是“可匿名遍历所有订阅”的通杀型漏洞。
- 证据文件：`api/sub.go`、`models/subscription_share.go`、`api/clients.go`、`routers/clients.go`、`api/share.go`、`middlewares/ip.go`、`models/iplogs.go`、`webs/src/views/subscriptions/component/ShareManageDialog.jsx`、`services/telegram/handlers.go`、`middlewares/jwt.go`、`routers/subcription.go`。

## 2026-04-28 原名预处理布局复查
- 用户反馈机场管理原名预处理区域仍存在布局错位，开始重新排查对应前端组件与样式。

- 处理结果：将 webs/src/components/NodeNamePreprocessor.jsx 的规则行从易受换行影响的 flex 布局改为稳定的响应式 grid 布局。桌面端固定为 拖拽柄 / 开关 / 匹配模式 / 查找框 / 箭头 / 替换框 / 删除按钮 七列，移动端改为首行控件加两行输入框，避免中间查找框被挤压或错位。
- 校验说明：按用户要求未继续执行前端构建，已完成静态代码复核，当前改动未引入额外状态或跨组件依赖。


## 2026-04-28 fork 仓库版本源排查
- 用户反馈 fork 后版本更新提示和更新日志仍显示上游仓库信息，开始排查版本源与更新日志实现。

- 修复内容：新增 webs/src/config/githubRepo.js 作为前端统一 GitHub 仓库源配置，默认从当前 fork 的 origin 仓库 oOyH/sublinkPro 拉取版本与发布信息，并支持用 VITE_GITHUB_REPO 覆盖。
- 替换范围：侧边栏版本更新卡改为读取共享仓库源；Dashboard 的更新日志与 Star 提醒卡也改为读取同一仓库源，避免继续显示上游 ZeroDeng01/sublinkPro 的版本和发布记录。
- 校验说明：未执行前端构建，已做静态代码复核，确认版本检查与发布日志请求地址已统一切换到共享配置。


## 2026-04-28 mihomo 检测失败排查
- 用户反馈同节点在 mihomo 中可正常测速，但本项目内检测失败，开始排查 mihomo 适配器构造与检测链路。

- 原因分析：截图里的红色“检测失败”不是节点主链路失败。项目内这两个节点的延迟测试其实成功了，失败的是延迟成功后的附加 IP 质量检测。当前实现一旦质量接口请求失败，就直接把 QualityStatus 写成 failed，前端又把它显示成“检测失败”，容易被误判成整个 mihomo 检测失败。
- 修复内容：为 IP 质量检测增加备用网址回退；当节点已经成功拿到出口 IP，但质量接口不可用时，将质量状态从 failed 降级为 partial，避免把可用节点误标成失败。
- 验证结果：`$env:GOTOOLCHAIN='local'; & 'D:\Program Files\Go\bin\go.exe' test ./services/mihomo ./services/scheduler` 通过。
