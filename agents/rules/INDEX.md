# agents/rules — sub2api fork 二开规则

> 本目录是本 fork 的**本地规则**，由 Tokenports 项目添加，上游没有这个目录。
> 跨系统契约不在这里，在总管仓的 `spec/`。

---

## 规则清单

| 文件 | 内容 |
|------|------|
| [fork-discipline.md](./fork-discipline.md) | 分支模型、commit 纪律、最小侵入、台账义务 |
| [git-commits.md](./git-commits.md) | Conventional Commits 格式与本仓 scope 速查表 |
| [tokenports-scope.md](./tokenports-scope.md) | 本仓在架构里的职责边界、计费配置约定、隐形开关 |
| [upstream-constraints.md](./upstream-constraints.md) | 上游工具链与工程约定（ent / pnpm / Go 版本 / 已知坑） |

---

## 优先级

冲突时的裁决顺序（从高到低）：

1. **上游工程约定** — ent schema 流程、pnpm、CI 要求。违反会导致构建或 CI 失败。
2. **总管仓 `spec/`** — 跨系统契约。
3. **本目录规则** — 本 fork 的工程纪律。

真冲突了说明 `spec/` 写错了，去改 `spec/`，别在代码里绕。

---

## 真理源路径

| 内容 | 路径 |
|------|------|
| 跨系统契约 | `/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/spec/INDEX.md` |
| 二开台账 | `/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/patches/sub2api.md` |
| 跨仓任务 | `/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/joint-tasks/` |
