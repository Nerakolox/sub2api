/**
 * 成本报表 API（Tokenports §1.6 §1.8 §3.6）
 * 后端始终返回 USD，汇率转换在前端展示层完成。
 */

import { apiClient } from '../client'

// ==================== Types ====================

export interface CostByModelRow {
  report_date: string
  model_id: string
  api_key_name: string
  pool_type: string
  request_count: number
  input_tokens: number
  cache_tokens: number
  output_tokens: number
  total_tokens: number
  cost_usd: number
}

export interface CostByModelTotal {
  request_count: number
  total_tokens: number
  cost_usd: number
}

export interface CostByModelResponse {
  rows: CostByModelRow[]
  total: CostByModelTotal
}

export interface CostByClientRow {
  report_month: string
  new_api_user_id: number | null
  new_api_token_id: number | null
  is_unattributed: boolean
  user_label: string
  request_count: number
  input_tokens: number
  cache_tokens: number
  output_tokens: number
  total_tokens: number
  cost_usd: number
  attribution_rate: number
}

export interface CostByClientResponse {
  rows: CostByClientRow[]
}

// ==================== API calls ====================

const costReportAPI = {
  getCostByModel(params: { start_date: string; end_date: string }) {
    return apiClient.get<{ data: CostByModelResponse }>(
      '/admin/reports/cost/by-model',
      { params }
    )
  },

  getCostByClient(params: { start_month: string; end_month: string }) {
    return apiClient.get<{ data: CostByClientResponse }>(
      '/admin/reports/cost/by-client',
      { params }
    )
  }
}

export default costReportAPI
