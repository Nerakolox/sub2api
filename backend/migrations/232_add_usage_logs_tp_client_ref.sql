ALTER TABLE usage_logs
    ADD COLUMN IF NOT EXISTS tp_client_ref varchar(64);

CREATE INDEX IF NOT EXISTS idx_usage_logs_tp_client_ref
    ON usage_logs (tp_client_ref);
