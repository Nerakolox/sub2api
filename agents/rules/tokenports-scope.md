# 本仓的职责边界

---

## 1. 在架构里的角色

Sub2API 是**后台成本层**。它对 New API 而言就是一个普通上游渠道，内部只有**一个用户（admin）**，所有客户流量都以 admin 的身份进入。

```
New API ──┬─ 渠道 1 ──> [本仓] 订阅分组 ──> oauth 账号池
          └─ 渠道 2 ──> [本仓] KEY 分组  ──> api_key 账号池
```

**本仓看不到客户身份。**这是架构的固有限制，不是缺陷。要做客户级归因，靠 `spec/03-request-ref.md` 定义的 header。

---

## 2. 本仓负责

| 职责 | 说明 |
|------|------|
| 上游账号池管理 | oauth 订阅账号 + api_key 账号 |
| 账号调度与粘性会话 | 健康度、冷却、优先级 |
| 供应商价卡 | 填在 Group 层 |
| 成本归集 | 每把 key 的 `quota_used` 累加 |
| 协议转换 | Claude / OpenAI / Gemini 等 |

---

## 3. 本仓不负责

| 不负责 | 归属 |
|--------|------|
| 客户身份、对客定价、客户账单 | New API |
| 预扣与并发闸门 | New API（本仓**不预扣**，允许透支） |
| 部署编排、数据库视图、迁移脚本 | 总管仓 |

---

## 4. 计费配置约定（配置不是代码，但改代码要知道）

统一使用**一个 admin 用户**承载全部链路。admin 在计费上没有特权——`CheckBillingEligibility` 与后扣逻辑都只看 `user_id`。

| 配置项 | 订阅池 | KEY 池 |
|--------|--------|--------|
| 账号类型 | `oauth` / `setup_token` | `api_key` |
| 分组 `rate_multiplier` | **`0`** | **`1.0`** |
| 分组价卡 | 留默认 | 供应商实际折扣价 |
| API Key | 一个 New API 渠道一把 | 同左 |
| API Key 额度 | 远高于实际用量 | 同左（见下方隐形开关） |

订阅池倍率为 `0` 的效果：`ActualCost = 0` 不扣余额，但 `usage_logs.total_cost` 保留影子金额，窗口余量监控照常。上游对 `RateMultiplier == 0` 有专门处理，这是被支持的用法。

---

## 5. 隐形开关（不配就静默失效）

**① API Key 额度决定成本是否入账**

```go
func (p *postUsageBillingParams) shouldDeductAPIKeyQuota() bool {
    return p.Cost.ActualCost > 0 && p.APIKey.Quota > 0 && p.APIKeyService != nil
}
```

key 上没设 `Quota`，`quota_used` 一分都不会累加。这是本项目**唯一的成本汇算依据**，不配等于成本账全空。

**② 账号配额决定 `AccountQuotaCost` 是否累加**

```go
return p.Cost.TotalCost > 0 && p.Account.IsAPIKeyOrBedrock() && p.Account.HasAnyQuotaLimit()
```

`IsAPIKeyOrBedrock()` 只认 `api_key` 和 `bedrock`——**订阅账号永远不会累加账号配额**。本项目用 key 口径汇算，这条不管。

---

## 6. 成本汇算口径（不要写出与之矛盾的逻辑）

- 用 `api_keys.quota_used` 累加值，**不要用 admin 余额的下降量**（余额会被充值打断且可能被手工污染）
- 采购成本取 `total_actual_cost`，不是 `total_cost`
- 订阅池的 `total_cost` 是影子金额，只用于分摊分析，**不进采购额**

---

## 7. admin 余额是可用性单点

本仓允许透支，预扣闸门又在 New API 侧，admin 余额会持续下降。跌破 0 时 `checkBalanceEligibility` 返回 `ErrInsufficientBalance`，**整条链路全部中断**。

改动计费相关代码时不要让这个边界更容易被触发。告警脚本在总管仓 `ops/`。

---

## 8. 计划中的改动

| 编号 | 内容 | 关联 SPEC |
|------|------|----------|
| 001 | `usage_logs.tp_client_ref` 字段 + 后扣时写入 | §3.2 |

写入点在后扣路径（`internal/service/gateway_usage_billing.go` · `usage_billing.go`）。缺失 header 时落 `NULL`，**不报错、不拒绝请求**。
