# Fork 纪律

> 本仓是 fork，长期要跟上游 rebase。所有纪律都服务于一个目标：**让下一次 rebase 尽可能不痛。**

---

## 1. 分支模型

| 分支 | 用途 | 纪律 |
|------|------|------|
| `main` | 上游镜像 | **永不提交自己的代码**。同步：`git fetch upstream && git reset --hard upstream/main` |
| `tokenports` | 二开分支 | 所有改动在这里 |

remote：

| 名称 | 地址 | 用途 |
|------|------|------|
| `origin` | 我们的 fork（SSH） | 推 |
| `upstream` | `https://github.com/Wei-Shaw/sub2api.git` | 只拉不推 |

---

## 2. 最小侵入（最重要的一条）

冲突面积和改动行数成正比。

**优先级**：加字段 > 加中间件 > 加配置开关 > 在现有函数里插一行 > 重写现有函数。

最后一种要在台账里说明为什么无法避免。

本仓的高危区是 `internal/service/`——计费、调度、定价解析都在这里，也是上游迭代最勤的地方。在这里改动要格外克制。

---

## 3. commit 纪律

- **一个改动一个 commit**，不 squash，不混杂重构
- message 写「**为什么改**」不是「改了什么」
- 不要顺手格式化无关文件
- **ent 的 generated 代码要和 schema 改动放同一个 commit**，否则 rebase 后会出现 schema 和生成代码不一致的中间态

---

## 4. 不碰上游文件的清单

| 文件 | 处置 |
|------|------|
| `openspec/` · `skills/` | 完全不动 |
| `DEV_GUIDE.md` · `README*.md` · `CLA.md` · `LICENSE` | 不动 |
| `.gitignore` | 不动。本地产物加进 `.git/info/exclude` |
| `deploy/` 下的官方 compose | 不动。我们的编排在总管仓 |

本地工具产物（`backend/.air.toml`、`backend/tmp/`）走 `.git/info/exclude`。

---

## 5. 命名前缀（硬约定）

**所有新增数据库列一律加 `tp_` 前缀**（如 `tp_client_ref`）。

上游以后很可能加同名字段。撞名的代价不是编译错误，是 migration 在生产环境静默改错列。

详见总管仓 `spec/04-naming.md`。

---

## 6. 台账义务

任何动了代码的任务，完成时更新：

```
/Users/nerakolo/sanxuninfo/Tokenports-Project/tokenports-relay/patches/sub2api.md
```

每条要写清：解决什么问题、改动形态、上游对应位置、rebase 时怎么判断还要不要保留。

**漏了视为任务未完成。**

---

## 7. 前端 API 路径约束

**不要在前端 API 调用里硬拼 `/api/v1/` 前缀。**

sub2api 后端路由基础路径已经是 `/api/v1`，若前端的 axios baseURL 或 API 客户端也拼了一次前缀，实际请求会变成 `/api/v1/api/v1/...` 导致 404。

规则：
- 前端 `api/` 目录下的路径字符串从 `/admin/...` 开始，不加 `/api/v1`
- baseURL 统一在一处配置（通常是 axios 实例），路径函数不重复拼前缀
- 新增接口写完后，用浏览器 DevTools Network 面板或 curl 确认实际请求 URL 是否正确，再提交

---

## 8. rebase 前的准备

1. 建备份分支（`tokenports-backup-<日期>`）
2. 读一遍台账，回顾每条改动的意图
3. `git fetch upstream && git rebase upstream/main tokenports`
4. 每处冲突先问：上游这段新代码是不是已经覆盖了我们的意图？
   - 是 → 丢弃我们的改动，台账标 `superseded`
   - 否 → 保留，必要时按上游新结构重写
5. **重新跑一次 ent generate**，确认生成代码与 schema 一致
6. 跑链路验证
7. 更新总管仓 README 的版本表
