<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Tab 切换 -->
      <div class="card">
        <div class="flex items-center border-b border-gray-200 px-2 dark:border-dark-700 sm:px-4">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            type="button"
            class="-mb-px inline-flex items-center gap-1.5 border-b-2 px-3 py-3 text-sm font-medium transition-colors sm:px-4"
            :class="activeTab === tab.key
              ? 'border-primary-500 text-primary-600 dark:text-primary-400'
              : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Tab 1：按模型 -->
        <div v-if="activeTab === 'by-model'" class="p-4 space-y-4">
          <div class="flex flex-wrap items-center gap-4">
            <DateRangePicker
              v-model:start-date="modelStartDate"
              v-model:end-date="modelEndDate"
              @change="loadModelData"
            />
            <button
              type="button"
              class="btn btn-secondary ml-auto"
              :disabled="modelLoading"
              @click="exportModelExcel"
            >
              导出 Excel
            </button>
          </div>

          <div v-if="modelLoading" class="py-12 text-center text-sm text-gray-500 dark:text-gray-400">
            加载中…
          </div>
          <div v-else-if="modelError" class="py-12 text-center text-sm text-red-500">
            {{ modelError }}
          </div>
          <div v-else class="overflow-x-auto">
            <table class="cost-table">
              <thead>
                <tr>
                  <th>日期</th>
                  <th>模型 ID</th>
                  <th>渠道 Key</th>
                  <th>池类型</th>
                  <th class="text-right">请求数</th>
                  <th class="text-right">输入 Tokens</th>
                  <th class="text-right">缓存 Tokens</th>
                  <th class="text-right">输出 Tokens</th>
                  <th class="text-right">总 Tokens</th>
                  <th class="text-right">成本(¥)</th>
                </tr>
              </thead>
              <tbody>
                <template v-for="(group, date) in modelGrouped" :key="date">
                  <tr
                    v-for="(row, i) in group.rows"
                    :key="`${date}-${i}`"
                  >
                    <td>{{ i === 0 ? row.report_date : '' }}</td>
                    <td>{{ row.model_id }}</td>
                    <td>{{ row.api_key_name }}</td>
                    <td>{{ poolTypeLabel(row.pool_type) }}</td>
                    <td class="text-right tabular-nums">{{ row.request_count.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ row.input_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ row.cache_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ row.output_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ row.total_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ toCNY(row.cost_usd) }}</td>
                  </tr>
                  <!-- 当日合计行 -->
                  <tr class="day-total">
                    <td>{{ date }} 合计</td>
                    <td colspan="3"></td>
                    <td class="text-right tabular-nums">{{ group.total.request_count.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ group.total.input_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ group.total.cache_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ group.total.output_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ group.total.total_tokens.toLocaleString() }}</td>
                    <td class="text-right tabular-nums">{{ toCNY(group.total.cost_usd) }}</td>
                  </tr>
                </template>
              </tbody>
              <!-- 区间合计 -->
              <tfoot v-if="modelData">
                <tr class="interval-total">
                  <td>区间合计</td>
                  <td colspan="3"></td>
                  <td class="text-right tabular-nums">{{ modelData.total.request_count.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">—</td>
                  <td class="text-right tabular-nums">—</td>
                  <td class="text-right tabular-nums">—</td>
                  <td class="text-right tabular-nums">{{ modelData.total.total_tokens.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ toCNY(modelData.total.cost_usd) }}</td>
                </tr>
              </tfoot>
            </table>
            <p v-if="!modelData?.rows?.length" class="py-8 text-center text-sm text-gray-500">该区间暂无数据</p>
          </div>
        </div>

        <!-- Tab 2：按客户 -->
        <div v-if="activeTab === 'by-client'" class="p-4 space-y-4">
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-2">
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300">开始月份</label>
              <input
                v-model="clientStartMonth"
                type="month"
                class="input input-sm"
                @change="loadClientData"
              />
            </div>
            <div class="flex items-center gap-2">
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300">结束月份</label>
              <input
                v-model="clientEndMonth"
                type="month"
                class="input input-sm"
                @change="loadClientData"
              />
            </div>
            <button
              type="button"
              class="btn btn-secondary ml-auto"
              :disabled="clientLoading"
              @click="exportClientExcel"
            >
              导出 Excel
            </button>
          </div>

          <div v-if="clientLoading" class="py-12 text-center text-sm text-gray-500 dark:text-gray-400">
            加载中…
          </div>
          <div v-else-if="clientError" class="py-12 text-center text-sm text-red-500">
            {{ clientError }}
          </div>
          <div v-else class="overflow-x-auto">
            <table class="cost-table">
              <thead>
                <tr>
                  <th>月份</th>
                  <th>New API 用户 ID</th>
                  <th class="text-right">请求数</th>
                  <th class="text-right">输入 Tokens</th>
                  <th class="text-right">缓存 Tokens</th>
                  <th class="text-right">输出 Tokens</th>
                  <th class="text-right">总 Tokens</th>
                  <th class="text-right">成本(¥)</th>
                  <th class="text-right">归因率</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(row, i) in clientData?.rows"
                  :key="i"
                  :class="{ 'unattributed': row.is_unattributed }"
                >
                  <td>{{ row.report_month }}</td>
                  <td :class="{ 'text-gray-400 dark:text-gray-500': row.is_unattributed }">
                    {{ row.is_unattributed ? '未归因' : row.new_api_user_id }}
                  </td>
                  <td class="text-right tabular-nums">{{ row.request_count.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ row.input_tokens.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ row.cache_tokens.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ row.output_tokens.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ row.total_tokens.toLocaleString() }}</td>
                  <td class="text-right tabular-nums">{{ toCNY(row.cost_usd) }}</td>
                  <td class="text-right tabular-nums">{{ (row.attribution_rate * 100).toFixed(1) }}%</td>
                </tr>
              </tbody>
            </table>
            <p v-if="!clientData?.rows?.length" class="py-8 text-center text-sm text-gray-500">该区间暂无数据</p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as XLSX from 'xlsx'
import { AppLayout } from '@/components/layout'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import { adminAPI } from '@/api/admin'
import type { CostByModelRow, CostByModelResponse, CostByClientResponse } from '@/api/admin/costReport'

// 展示层 USD → CNY 转换；后端接口始终返回 USD，汇率可后续做成 admin 配置
const USD_TO_CNY = 7.2

function toCNY(usd: number): string {
  return '¥' + (usd * USD_TO_CNY).toFixed(2)
}

function poolTypeLabel(t: string): string {
  return t === 'subscription' ? '订阅池' : 'KEY 池'
}

// ==================== Tab ====================
const tabs = [
  { key: 'by-model', label: '按模型' },
  { key: 'by-client', label: '按客户' }
] as const
const activeTab = ref<'by-model' | 'by-client'>('by-model')

// ==================== 按模型 ====================
const now = new Date()
const firstOfMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`
const today = now.toISOString().slice(0, 10)

const modelStartDate = ref(firstOfMonth)
const modelEndDate = ref(today)
const modelLoading = ref(false)
const modelError = ref('')
const modelData = ref<CostByModelResponse | null>(null)

interface DayGroup {
  rows: CostByModelRow[]
  total: {
    request_count: number
    input_tokens: number
    cache_tokens: number
    output_tokens: number
    total_tokens: number
    cost_usd: number
  }
}

const modelGrouped = computed((): Record<string, DayGroup> => {
  if (!modelData.value?.rows) return {}
  const groups: Record<string, DayGroup> = {}
  for (const row of modelData.value.rows) {
    if (!groups[row.report_date]) {
      groups[row.report_date] = {
        rows: [],
        total: { request_count: 0, input_tokens: 0, cache_tokens: 0, output_tokens: 0, total_tokens: 0, cost_usd: 0 }
      }
    }
    groups[row.report_date].rows.push(row)
    const t = groups[row.report_date].total
    t.request_count += row.request_count
    t.input_tokens += row.input_tokens
    t.cache_tokens += row.cache_tokens
    t.output_tokens += row.output_tokens
    t.total_tokens += row.total_tokens
    t.cost_usd += row.cost_usd
  }
  return groups
})

async function loadModelData() {
  modelLoading.value = true
  modelError.value = ''
  try {
    const res = await adminAPI.costReport.getCostByModel({
      start_date: modelStartDate.value,
      end_date: modelEndDate.value
    })
    modelData.value = (res.data as { data: CostByModelResponse }).data
  } catch (e: unknown) {
    modelError.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    modelLoading.value = false
  }
}

function exportModelExcel() {
  if (!modelData.value) return
  const COLS = ['日期', '模型 ID', '渠道 Key', '池类型', '请求数', '输入 Tokens', '缓存 Tokens', '输出 Tokens', '总 Tokens', '成本(¥)']
  const rows: (string | number)[][] = [COLS]
  for (const [date, group] of Object.entries(modelGrouped.value)) {
    for (const r of group.rows) {
      rows.push([r.report_date, r.model_id, r.api_key_name, poolTypeLabel(r.pool_type), r.request_count, r.input_tokens, r.cache_tokens, r.output_tokens, r.total_tokens, parseFloat((r.cost_usd * USD_TO_CNY).toFixed(2))])
    }
    rows.push([`${date} 合计`, '', '', '', group.total.request_count, group.total.input_tokens, group.total.cache_tokens, group.total.output_tokens, group.total.total_tokens, parseFloat((group.total.cost_usd * USD_TO_CNY).toFixed(2))])
  }
  rows.push(['区间合计', '', '', '', modelData.value.total.request_count, '', '', '', modelData.value.total.total_tokens, parseFloat((modelData.value.total.cost_usd * USD_TO_CNY).toFixed(2))])
  const ws = XLSX.utils.aoa_to_sheet(rows)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '每日明细')
  XLSX.writeFile(wb, `成本报告_按模型_${modelStartDate.value}_${modelEndDate.value}.xlsx`)
}

// ==================== 按客户 ====================
function monthStr(offsetMonths: number): string {
  const d = new Date(now.getFullYear(), now.getMonth() - offsetMonths, 1)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
}

const clientStartMonth = ref(monthStr(2))
const clientEndMonth = ref(monthStr(0))
const clientLoading = ref(false)
const clientError = ref('')
const clientData = ref<CostByClientResponse | null>(null)

async function loadClientData() {
  clientLoading.value = true
  clientError.value = ''
  try {
    const res = await adminAPI.costReport.getCostByClient({
      start_month: clientStartMonth.value,
      end_month: clientEndMonth.value
    })
    clientData.value = (res.data as { data: CostByClientResponse }).data
  } catch (e: unknown) {
    clientError.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    clientLoading.value = false
  }
}

function exportClientExcel() {
  if (!clientData.value) return
  const COLS = ['月份', 'New API 用户 ID', '请求数', '输入 Tokens', '缓存 Tokens', '输出 Tokens', '总 Tokens', '成本(¥)', '归因率']
  const rows: (string | number)[][] = [COLS]
  for (const r of clientData.value.rows) {
    rows.push([r.report_month, r.is_unattributed ? '未归因' : (r.new_api_user_id ?? ''), r.request_count, r.input_tokens, r.cache_tokens, r.output_tokens, r.total_tokens, parseFloat((r.cost_usd * USD_TO_CNY).toFixed(2)), (r.attribution_rate * 100).toFixed(1) + '%'])
  }
  const ws = XLSX.utils.aoa_to_sheet(rows)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '按客户明细')
  XLSX.writeFile(wb, `成本报告_按客户_${clientStartMonth.value}_${clientEndMonth.value}.xlsx`)
}

onMounted(() => {
  loadModelData()
  loadClientData()
})
</script>

<style scoped>
.cost-table {
  @apply w-full text-sm border-collapse;
}
.cost-table th {
  @apply px-3 py-2 text-left text-xs font-semibold uppercase tracking-wide
    text-gray-500 dark:text-gray-400
    border-b border-gray-200 dark:border-dark-600
    bg-gray-50 dark:bg-dark-800 whitespace-nowrap;
}
.cost-table td {
  @apply px-3 py-2 border-b border-gray-100 dark:border-dark-700
    text-gray-700 dark:text-gray-200 whitespace-nowrap;
}
.cost-table tr:hover td {
  @apply bg-gray-50 dark:bg-dark-800/50;
}
.day-total td {
  @apply bg-blue-50 dark:bg-blue-900/20 font-semibold text-gray-800 dark:text-gray-100;
}
.cost-table tfoot .interval-total td {
  @apply bg-primary-50 dark:bg-primary-900/20 font-bold text-primary-700 dark:text-primary-300;
}
.unattributed td {
  @apply text-gray-400 dark:text-gray-500;
}
</style>
