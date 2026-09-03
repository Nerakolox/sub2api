// Package admin provides HTTP handlers for administrative operations.
package admin

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// CostReportHandler 成本报表 handler（Tokenports §1.6 §1.8 §3.6）
type CostReportHandler struct {
	db *sql.DB
}

// NewCostReportHandler creates a new cost report handler.
func NewCostReportHandler(db *sql.DB) *CostReportHandler {
	return &CostReportHandler{db: db}
}

// costByModelRow mirrors one row of v_cost_by_model_daily.
type costByModelRow struct {
	ReportDate    string  `json:"report_date"`
	ModelID       string  `json:"model_id"`
	APIKeyName    string  `json:"api_key_name"`
	PoolType      string  `json:"pool_type"`
	RequestCount  int64   `json:"request_count"`
	InputTokens   int64   `json:"input_tokens"`
	CacheTokens   int64   `json:"cache_tokens"`
	OutputTokens  int64   `json:"output_tokens"`
	TotalTokens   int64   `json:"total_tokens"`
	CostUSD       float64 `json:"cost_usd"`
}

// costByModelTotal holds the aggregated totals for the requested interval.
type costByModelTotal struct {
	RequestCount int64   `json:"request_count"`
	TotalTokens  int64   `json:"total_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// GetCostByModel handles GET /api/v1/admin/reports/cost/by-model
// 参数：start_date / end_date（YYYY-MM-DD），端点只返回 USD，汇率转换在前端展示层完成。
func (h *CostReportHandler) GetCostByModel(c *gin.Context) {
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	if startStr == "" || endStr == "" {
		response.BadRequest(c, "start_date and end_date are required (YYYY-MM-DD)")
		return
	}

	if _, err := time.Parse("2006-01-02", startStr); err != nil {
		response.BadRequest(c, "invalid start_date format, expected YYYY-MM-DD")
		return
	}
	if _, err := time.Parse("2006-01-02", endStr); err != nil {
		response.BadRequest(c, "invalid end_date format, expected YYYY-MM-DD")
		return
	}

	const query = `
SELECT
    report_date::text,
    COALESCE(model_id, ''),
    COALESCE(api_key_name, ''),
    pool_type,
    request_count,
    input_tokens,
    cache_tokens,
    output_tokens,
    total_tokens,
    cost_usd
FROM v_cost_by_model_daily
WHERE report_date BETWEEN $1::date AND $2::date
ORDER BY report_date, model_id, api_key_name
`

	rows, err := h.db.QueryContext(c.Request.Context(), query, startStr, endStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to query cost by model: "+err.Error())
		return
	}
	defer rows.Close()

	var data []costByModelRow
	var total costByModelTotal

	for rows.Next() {
		var r costByModelRow
		if err := rows.Scan(
			&r.ReportDate, &r.ModelID, &r.APIKeyName, &r.PoolType,
			&r.RequestCount, &r.InputTokens, &r.CacheTokens,
			&r.OutputTokens, &r.TotalTokens, &r.CostUSD,
		); err != nil {
			response.Error(c, http.StatusInternalServerError, "row scan error: "+err.Error())
			return
		}
		data = append(data, r)
		total.RequestCount += r.RequestCount
		total.TotalTokens += r.TotalTokens
		total.CostUSD += r.CostUSD
	}
	if err := rows.Err(); err != nil {
		response.Error(c, http.StatusInternalServerError, "rows iteration error: "+err.Error())
		return
	}

	if data == nil {
		data = []costByModelRow{}
	}

	response.Success(c, gin.H{
		"rows":  data,
		"total": total,
	})
}

// costByClientRow mirrors one row of v_cost_by_client_monthly.
type costByClientRow struct {
	ReportMonth    string  `json:"report_month"`
	NewAPIUserID   *int64  `json:"new_api_user_id"`
	NewAPITokenID  *int64  `json:"new_api_token_id"`
	IsUnattributed bool    `json:"is_unattributed"`
	UserLabel      string  `json:"user_label"`
	RequestCount   int64   `json:"request_count"`
	InputTokens    int64   `json:"input_tokens"`
	CacheTokens    int64   `json:"cache_tokens"`
	OutputTokens   int64   `json:"output_tokens"`
	TotalTokens    int64   `json:"total_tokens"`
	CostUSD        float64 `json:"cost_usd"`
	AttributionRate float64 `json:"attribution_rate"`
}

// GetCostByClient handles GET /api/v1/admin/reports/cost/by-client
// 参数：start_month / end_month（YYYY-MM），new_api_user_id 为 NULL 的行 user_label 为「未归因」。
func (h *CostReportHandler) GetCostByClient(c *gin.Context) {
	startStr := c.Query("start_month")
	endStr := c.Query("end_month")

	if startStr == "" || endStr == "" {
		response.BadRequest(c, "start_month and end_month are required (YYYY-MM)")
		return
	}

	// 把 YYYY-MM 转成月初日期，方便与视图的 report_month 比较
	startDate := startStr + "-01"
	endDate := endStr + "-01"

	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		response.BadRequest(c, "invalid start_month format, expected YYYY-MM")
		return
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		response.BadRequest(c, "invalid end_month format, expected YYYY-MM")
		return
	}

	const query = `
SELECT
    report_month::text,
    new_api_user_id,
    new_api_token_id,
    is_unattributed,
    request_count,
    input_tokens,
    cache_tokens,
    output_tokens,
    total_tokens,
    cost_usd,
    attribution_rate
FROM v_cost_by_client_monthly
WHERE report_month BETWEEN $1::date AND $2::date
ORDER BY report_month, new_api_user_id NULLS LAST
`

	rows, err := h.db.QueryContext(c.Request.Context(), query, startDate, endDate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to query cost by client: "+err.Error())
		return
	}
	defer rows.Close()

	var data []costByClientRow

	for rows.Next() {
		var r costByClientRow
		if err := rows.Scan(
			&r.ReportMonth, &r.NewAPIUserID, &r.NewAPITokenID, &r.IsUnattributed,
			&r.RequestCount, &r.InputTokens, &r.CacheTokens,
			&r.OutputTokens, &r.TotalTokens, &r.CostUSD, &r.AttributionRate,
		); err != nil {
			response.Error(c, http.StatusInternalServerError, "row scan error: "+err.Error())
			return
		}
		if r.IsUnattributed {
			r.UserLabel = "未归因"
		} else if r.NewAPIUserID != nil {
			// 前端会进一步格式化，这里只提供原始 ID 字符串方便展示
			r.UserLabel = ""
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		response.Error(c, http.StatusInternalServerError, "rows iteration error: "+err.Error())
		return
	}

	if data == nil {
		data = []costByClientRow{}
	}

	response.Success(c, gin.H{"rows": data})
}
