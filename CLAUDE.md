# CLAUDE.md — Tokenports fork of sub2api

> 上游仓库没有本文件，它由 Tokenports 项目添加。

本仓库不是独立项目，是 `Wei-Shaw/sub2api` 的 fork，在 Tokenports 架构里承担**账号池调度与供应商成本层**的角色。

---

## 强制阅读顺序

每个会话、每个新任务开工前，按顺序读完：

1. [`agents/rules/INDEX.md`](./agents/rules/INDEX.md) — 本 fork 的二开规则
2. `/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/spec/INDEX.md` — 跨系统契约**真理源**
3. `DEV_GUIDE.md`（上游）— 开发环境与已知坑

第 2 项**不在本仓内**。任何本仓代码与它冲突，停下来对齐它，不要边写边改。

---

## 三条红线

- **本仓只做 Sub2API 侧**。不要碰 `../new-api/`，也不要碰 `../tokenports-relay/`。需要另一侧配合时，产出任务描述交给总管。
- **契约先行**。跨系统的字段、header、口径约定必须先落 `spec/`，再写代码。
- **最小侵入**。优先钩子式改动（加字段、加中间件、加开关），避免重写上游现有函数。冲突面积和改动行数成正比。

---

## 上游资产不要动

| 路径 | 内容 |
|------|------|
| `openspec/` | 上游的变更规范流程 |
| `skills/sub2api-admin/` | 上游的 admin skill |
| `DEV_GUIDE.md` · `README*.md` | 上游文档 |
| `.gitignore` | 本地产物走 `.git/info/exclude` |

---

## 完成任务的最后一步

任何动了代码的任务，完成时必须更新总管仓的二开台账：

```
/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/patches/sub2api.md
```

**漏了这一步视为任务未完成。**台账记的是意图，不是 diff——rebase 解冲突时唯一能用的就是它。

---

## 沟通口径

所有沟通、文档、代码注释默认用中文。上游原有的英文注释不翻译。
