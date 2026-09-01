# 上游工程约定

> 违反这些会导致构建失败、CI 失败，或在 rebase 时制造不必要的冲突。

---

## 1. 工具链

| 项 | 版本 / 要求 |
|---|------------|
| Go | **1.27.0+**（比 new-api 高，装高的那个） |
| 前端包管理器 | **pnpm**，不是 npm 也不是 yarn |
| 前端框架 | Vue 3 |
| ORM | **ent**（不是 GORM） |

---

## 2. pnpm 相关的两个已知坑

**坑 1：`pnpm-lock.yaml` 必须同步提交**

改了 `frontend/package.json` 之后，上游 CI 的 `pnpm install --frozen-lockfile` 会因为 lock 不同步而失败。改完依赖记得：

```bash
cd frontend && pnpm install
git add pnpm-lock.yaml
```

**坑 2：npm 与 pnpm 的 `node_modules` 冲突**

之前用 npm 装过再切 pnpm 会报 `EPERM`。先删干净 `node_modules` 再 `pnpm install`。

> 附带：如果 `pnpm` 报「configured to use yarn because /Users/<你>/package.json has a packageManager field」，问题在**家目录**那个 `package.json`，不在本仓。把它挪走即可，不要改本仓的 `package.json`。

---

## 3. ent schema 改动流程

加字段不是改个 struct 就完了：

1. 改 ent schema
2. 重新 generate
3. **生成代码和 schema 改动放同一个 commit**

分成两个 commit 会导致 rebase 后出现不一致的中间态。

---

## 4. 目录结构

| 路径 | 内容 |
|------|------|
| `backend/internal/service/` | 计费、调度、定价解析。**高危区**，上游迭代最勤 |
| `backend/internal/domain/` | 常量与领域定义（账号类型等） |
| `backend/cmd/server/` | 入口 |
| `frontend/` | Vue 3 + pnpm |
| `deploy/` | 官方部署编排。**不动**，我们的编排在总管仓 |

---

## 5. 计费相关的关键文件

改计费前先读这几个，别凭记忆：

| 文件 | 内容 |
|------|------|
| `internal/service/billing_service.go` | 成本计算 |
| `internal/service/model_pricing_resolver.go` | 定价解析链 `Group → Channel → LiteLLM → Fallback` |
| `internal/service/gateway_usage_billing.go` · `usage_billing.go` | 后扣与五个账桶 |
| `internal/service/billing_cache_service.go` | 准入校验 |
| `internal/domain/constants.go` | 账号类型、角色常量 |

---

## 6. 运行模式

`RUN_MODE=simple` 会**跳过整个计费流程**，只保留调度与协议转换。

本项目用 `standard`。不要为了「简化调试」临时切 simple 然后忘记切回来——那会导致成本账出现无法解释的缺口。

---

## 7. 环境变量

本仓用**拆开的** `DATABASE_*` 系列（不是 DSN 字符串），Redis 用 `REDIS_HOST` / `REDIS_PORT` / `REDIS_DB`。

`AUTO_SETUP=true` 会自动生成配置并初始化，跳过安装向导。

`JWT_SECRET` 与 `TOTP_ENCRYPTION_KEY` 留空会每次启动随机生成——前者导致登录态失效，后者导致**已绑定 2FA 的账号永久失效**。开发环境也要固定。

---

## 8. 前端嵌入需要 `-tags embed`（已踩过）

`backend/internal/web/embed_on.go` 顶部是 `//go:build embed`，对应的 `embed_off.go` 是 `//go:build !embed`，后者的 `HasEmbeddedFrontend()` 固定返回 `false`。

路由注册时只有 `HasEmbeddedFrontend()` 为真才挂前端中间件（`internal/server/router.go`）。所以：

**不加 `-tags embed` 编译 → 前端路由整个不注册 → 访问 `/` 得到 Gin 默认 `404 page not found`，而 `/health` 返回 200。**

前端产物路径本身没问题：`frontend/vite.config.ts` 的 `outDir` 指向 `../backend/internal/web/dist`，与 `//go:embed all:dist` 期望的位置一致。**产物存在不等于被嵌入**——嵌入发生在编译期。

| 场景 | 做法 |
|------|------|
| 开发期 | 单独跑 `pnpm dev --port 3001`，Vite 自动把 API 代理到 8080。不需要嵌入 |
| 单端口验证 | `go build -tags embed -o ./tmp/sub2api ./cmd/server` |
| 构建镜像 | 无需处理，`Dockerfile` 已带该 tag |

项目的 air 配置**故意不带这个 tag**——带上会让每次热重载都重嵌近两百个文件。改动 air 配置前先想清楚这一点。

重新构建后**必须换新进程**才生效，当前运行的二进制不会动态读取新的 dist。

`RUN_MODE` 与此无关，它只影响计费逻辑。
