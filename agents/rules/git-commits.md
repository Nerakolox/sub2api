# Git 提交规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)，格式：`type(scope): 描述`

> commit 纪律（一改一 commit、写意图不写行为、不格式化无关文件、ent generated 代码与 schema 同 commit）在 [fork-discipline.md §3](./fork-discipline.md) 里，本文件只定义格式与 scope。

## 类型

| 类型 | 用途 |
|------|------|
| `feat` | 新增功能或字段 |
| `fix` | 修复 bug |
| `docs` | 仅文档变更 |
| `refactor` | 重构，不改行为 |
| `test` | 测试相关 |
| `chore` | 构建、依赖、CI 等杂项 |

## Scope（本仓常用）

| scope | 对应层/功能 |
|-------|------------|
| `pool` | 订阅池与 KEY 池划分 |
| `schedule` | 账号调度与轮询逻辑 |
| `pricing` | 价卡与成本计算 |
| `account` | 账号管理与状态 |
| `ent` | ent schema 及生成代码（schema 与 generated 必须同 commit） |
| `service` | `internal/service/` 业务逻辑 |
| `web` | 前端页面与组件 |
| `middleware` | 中间件层 |
| `config` | 系统配置项 |
| `api` | HTTP 路由与 handler |

## 示例

```
feat(ent): 新增 tp_client_ref 字段并同步 generated 代码
fix(schedule): 修复账号耗尽后未触发备用池切换的问题
feat(pricing): 支持从总管仓 spec §2.3 读取成本倍率覆盖
refactor(service): 将余额扣减逻辑提取为独立 transaction helper
feat(web): 价卡列表新增成本列与隐形开关状态展示
```

## 补充约定

- 描述用中文，简短且完整，不超过 50 字
- ent schema 改动必须和对应 `generate` 产生的文件放同一个 commit，原因见 `fork-discipline.md §3`
- 完成后必须同步更新总管仓台账：`tokenports-relay/patches/sub2api.md`（漏了视为未完成）
