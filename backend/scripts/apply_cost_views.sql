-- Tokenports 成本报表视图部署说明
-- 将此文件下方的两个视图 apply 到 Sub2API 的 PostgreSQL 库：
--   psql $SUB2API_DATABASE_URL -f backend/scripts/apply_cost_views.sql
-- 视图依赖 usage_logs、api_keys、groups 三张表，均已存在。
-- 可重复执行（CREATE OR REPLACE）。

-- ======================================================
-- 1. 分模型 / 分渠道每日成本视图（SPEC §1.8）
-- ======================================================
CREATE OR REPLACE VIEW v_cost_by_model_daily AS
SELECT
    date_trunc('day', ul.created_at)::date                                          AS report_date,
    ul.model                                                                         AS model_id,
    ul.api_key_id,
    ak.name                                                                          AS api_key_name,
    CASE WHEN COALESCE(g.rate_multiplier, 1.0) = 0 THEN 'subscription' ELSE 'key' END AS pool_type,
    COUNT(*)                                                                         AS request_count,
    SUM(ul.input_tokens)                                                             AS input_tokens,
    SUM(ul.cache_creation_tokens + ul.cache_read_tokens)                             AS cache_tokens,
    SUM(ul.output_tokens)                                                            AS output_tokens,
    SUM(ul.input_tokens + ul.cache_creation_tokens + ul.cache_read_tokens + ul.output_tokens) AS total_tokens,
    SUM(ul.actual_cost)                                                              AS cost_usd
FROM usage_logs ul
JOIN api_keys ak ON ak.id = ul.api_key_id
LEFT JOIN groups g ON g.id = ak.group_id
GROUP BY
    date_trunc('day', ul.created_at),
    ul.model,
    ul.api_key_id,
    ak.name,
    COALESCE(g.rate_multiplier, 1.0);

-- ======================================================
-- 2. 客户维度每月成本视图（SPEC §3.6）
-- ======================================================
CREATE OR REPLACE VIEW v_cost_by_client_monthly AS
SELECT
    date_trunc('month', ul.created_at)::date                                        AS report_month,
    NULLIF(split_part(ul.tp_client_ref, ':', 1), '')::bigint                        AS new_api_user_id,
    NULLIF(split_part(ul.tp_client_ref, ':', 2), '')::bigint                        AS new_api_token_id,
    CASE WHEN ul.tp_client_ref IS NULL THEN true ELSE false END                     AS is_unattributed,
    COUNT(*)                                                                         AS request_count,
    SUM(ul.input_tokens)                                                             AS input_tokens,
    SUM(ul.cache_creation_tokens + ul.cache_read_tokens)                             AS cache_tokens,
    SUM(ul.output_tokens)                                                            AS output_tokens,
    SUM(ul.input_tokens + ul.cache_creation_tokens + ul.cache_read_tokens + ul.output_tokens) AS total_tokens,
    SUM(ul.actual_cost)                                                              AS cost_usd,
    COUNT(*) FILTER (WHERE ul.tp_client_ref IS NOT NULL) * 1.0 / COUNT(*)           AS attribution_rate
FROM usage_logs ul
GROUP BY
    date_trunc('month', ul.created_at),
    ul.tp_client_ref;
