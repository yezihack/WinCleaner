<template>
  <div class="cleaner-page">
    <div class="page-header">
      <h2>垃圾清理</h2>
      <p class="page-sub">扫描并清理系统垃圾文件，释放磁盘空间</p>
    </div>

    <div class="scan-action-row">
      <button class="scan-btn" :class="{ loading: scanning }" :disabled="scanning" @click="handleScan">
        <span class="btn-icon">🔍</span>
        <span>{{ scanning ? '扫描中...' : '扫描垃圾' }}</span>
      </button>
      <button
        class="clean-btn"
        :class="{ disabled: selectedCategories.length === 0 }"
        :disabled="selectedCategories.length === 0 || cleaning"
        @click="handleClean"
      >
        <span class="btn-icon">🗑️</span>
        <span>{{ cleaning ? '清理中...' : `清理选中 (${formatBytes(selectedSize)})` }}</span>
      </button>
    </div>

    <div v-if="results.length > 0" class="results-section">
      <div class="results-header">
        <span class="results-count">共 {{ results.length }} 个分类</span>
      </div>

      <div class="category-list">
        <div
          v-for="row in results"
          :key="row.category"
          class="category-item"
          :class="{ selected: isSelected(row.category) }"
          @click="toggleSelect(row)"
        >
          <div class="cat-check">
            <div class="check-box" :class="{ checked: isSelected(row.category) }">
              <span v-if="isSelected(row.category)">✓</span>
            </div>
          </div>
          <div class="cat-icon">{{ getCategoryIcon(row.category) }}</div>
          <div class="cat-info">
            <div class="cat-name">{{ row.category }}</div>
            <div class="cat-meta">{{ row.count }} 个文件</div>
          </div>
          <div class="cat-size">{{ formatBytes(row.size) }}</div>
          <div class="cat-expand" @click.stop="toggleExpand(row.category)">
            <span>{{ expandedCategories.has(row.category) ? '▲' : '▼' }}</span>
          </div>
        </div>

        <div v-if="expandedCategories.size > 0" class="expanded-panels">
          <div v-for="cat in expandedResults" :key="cat.category" class="expand-panel">
            <div class="expand-header">
              <span class="expand-title">{{ cat.category }} - {{ cat.count }} 个文件</span>
              <el-input
                v-model="filterText[cat.category]"
                placeholder="搜索文件..."
                clearable
                size="small"
                style="width: 200px;"
              />
            </div>
            <div class="file-list">
              <div
                v-for="item in getFilteredItems(cat)"
                :key="item.path"
                class="file-item"
              >
                <span class="file-path">{{ item.path }}</span>
                <span class="file-size">{{ formatBytes(item.size) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="!scanning" class="empty-state">
      <div class="empty-icon">📂</div>
      <div class="empty-text">点击上方扫描按钮开始检测垃圾文件</div>
    </div>

    <div v-if="cleanResult" class="result-banner" :class="cleanResult.freed_size > 1024 * 1024 * 100 ? 'great' : 'ok'">
      <div class="result-icon">✨</div>
      <div class="result-info">
        <div class="result-main">释放空间: <strong>{{ formatBytes(cleanResult.freed_size) }}</strong></div>
        <div class="result-sub">成功清理 {{ cleanResult.cleaned_count }} 个文件，{{ cleanResult.failed_count }} 个跳过</div>
      </div>
    </div>

    <div v-if="history" class="history-section">
      <div class="section-header">
        <h3>清理历史</h3>
        <button class="icon-btn" @click="loadHistory">🔄</button>
      </div>

      <div class="history-stats">
        <div class="hstat-card">
          <div class="hstat-val">{{ history.last_clean_ago || '从未清理' }}</div>
          <div class="hstat-key">距上次清理</div>
        </div>
        <div class="hstat-card">
          <div class="hstat-val">{{ formatBytes(history.total_freed) }}</div>
          <div class="hstat-key">累计释放</div>
        </div>
        <div class="hstat-card">
          <div class="hstat-val">{{ history.total_count }}</div>
          <div class="hstat-key">累计文件数</div>
        </div>
      </div>

      <div class="chart-toggle">
        <button
          v-for="m in ['daily', 'monthly']"
          :key="m"
          class="toggle-btn"
          :class="{ active: chartMode === m }"
          @click="chartMode = m as 'daily' | 'monthly'"
        >
          {{ m === 'daily' ? '近30天' : '按月' }}
        </button>
      </div>

      <div class="chart-wrap">
        <v-chart :option="chartOption" style="height: 280px;" autoresize />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { use } from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, GridComponent,
  LegendComponent, DataZoomComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { api, type ScanResult, type CleanResult, type CleanHistoryStats } from '@/api/backend'

use([
  BarChart, LineChart, TitleComponent, TooltipComponent,
  GridComponent, LegendComponent, DataZoomComponent, CanvasRenderer
])

const scanning = ref(false)
const cleaning = ref(false)
const results = ref<ScanResult[]>([])
const selectedRows = ref<ScanResult[]>([])
const cleanResult = ref<CleanResult | null>(null)
const filterText = reactive<Record<string, string>>({})
const history = ref<CleanHistoryStats | null>(null)
const chartMode = ref<'daily' | 'monthly'>('daily')
const expandedCategories = reactive(new Set<string>())

const selectedCategories = computed(() => selectedRows.value.map(r => r.category))
const selectedSize = computed(() => selectedRows.value.reduce((sum, r) => sum + r.size, 0))
const expandedResults = computed(() => results.value.filter(r => expandedCategories.has(r.category)))

const isSelected = (cat: string) => selectedRows.value.some(r => r.category === cat)

const toggleSelect = (row: ScanResult) => {
  const idx = selectedRows.value.findIndex(r => r.category === row.category)
  if (idx >= 0) {
    selectedRows.value.splice(idx, 1)
  } else {
    selectedRows.value.push(row)
  }
}

const toggleExpand = (cat: string) => {
  if (expandedCategories.has(cat)) {
    expandedCategories.delete(cat)
  } else {
    expandedCategories.add(cat)
  }
}

const getCategoryIcon = (cat: string): string => {
  const map: Record<string, string> = {
    '系统临时文件': '🗂️', 'Windows Update 缓存': '🔄',
    '缩略图缓存': '🖼️', '系统日志': '📋',
    '浏览器缓存': '🌐', '回收站': '🗑️', 'Windows 预读取': '⚡',
  }
  return map[cat] || '📁'
}

const getFilteredItems = (row: ScanResult) => {
  const keyword = (filterText[row.category] || '').toLowerCase()
  if (!keyword) return row.items
  return row.items.filter(item => item.path.toLowerCase().includes(keyword))
}

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const toMB = (bytes: number) => +(bytes / 1024 / 1024).toFixed(1)

const chartOption = computed(() => {
  if (!history.value) return {}
  const isDaily = chartMode.value === 'daily'
  const data = isDaily ? (history.value.daily_stats || []) : (history.value.monthly_stats || [])
  const xLabels = data.map((d: any) => isDaily ? d.date : d.month)
  const freedData = data.map((d: any) => toMB(d.freed_size))
  const countData = data.map((d: any) => d.count)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['释放(MB)', '文件数'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 48, bottom: 40, top: 24 },
    xAxis: { type: 'category', data: xLabels, axisLabel: { rotate: isDaily ? 45 : 0, fontSize: 11 } },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left', axisLabel: { fontSize: 10 } },
      { type: 'value', name: '文件', position: 'right', axisLabel: { fontSize: 10 } }
    ],
    series: [
      { name: '释放(MB)', type: 'bar', data: freedData, itemStyle: { color: '#22c55e', borderRadius: [3, 3, 0, 0] }, yAxisIndex: 0 },
      { name: '文件数', type: 'line', data: countData, itemStyle: { color: '#3b82f6' }, smooth: true, yAxisIndex: 1 }
    ]
  }
})

const handleScan = async () => {
  scanning.value = true
  cleanResult.value = null
  try {
    results.value = await api.scanJunk()
    ElMessage.success(`扫描完成，发现 ${results.value.length} 个分类`)
  } catch {
    ElMessage.error('扫描失败')
  } finally {
    scanning.value = false
  }
}

const handleClean = async () => {
  try {
    await ElMessageBox.confirm(`确定清理选中的 ${selectedCategories.value.length} 个分类？`, '确认清理', { type: 'warning' })
  } catch { return }
  cleaning.value = true
  try {
    cleanResult.value = await api.cleanJunk(selectedCategories.value)
    ElMessage.success('清理完成')
    results.value = await api.scanJunk()
    selectedRows.value = []
    await loadHistory()
  } catch {
    ElMessage.error('清理失败')
  } finally {
    cleaning.value = false
  }
}

const loadHistory = async () => {
  try {
    const data = await api.getCleanHistory()
    history.value = (data && (data.total_count > 0 || (data.records && data.records.length > 0))) ? data : null
  } catch {
    history.value = null
  }
}

onMounted(() => { loadHistory() })
</script>

<style scoped>
.cleaner-page { padding: 0 4px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.scan-action-row { display: flex; gap: 12px; margin-bottom: 20px; }

.scan-btn, .clean-btn {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 24px; border: none; border-radius: 10px;
  font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.2s ease;
}
.scan-btn { background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); color: #fff; }
.scan-btn:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 4px 16px rgba(15,23,42,0.3); }
.scan-btn.loading { opacity: 0.7; cursor: not-allowed; }

.clean-btn { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: #fff; }
.clean-btn:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 4px 16px rgba(239,68,68,0.35); }
.clean-btn.disabled, .clean-btn:disabled { background: linear-gradient(135deg, #94a3b8 0%, #64748b 100%); cursor: not-allowed; }

.btn-icon { font-size: 16px; }

.results-section { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; margin-bottom: 20px; }
.results-header { padding: 12px 16px; border-bottom: 1px solid #f1f5f9; }
.results-count { font-size: 13px; color: #64748b; }

.category-list {}
.category-item {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 16px; border-bottom: 1px solid #f8fafc;
  cursor: pointer; transition: background 0.15s ease;
}
.category-item:hover { background: #f8fafc; }
.category-item:last-child { border-bottom: none; }
.category-item.selected { background: rgba(34,197,94,0.06); }

.cat-check { width: 20px; flex-shrink: 0; }
.check-box { width: 18px; height: 18px; border: 2px solid #cbd5e1; border-radius: 4px; display: flex; align-items: center; justify-content: center; transition: all 0.15s ease; }
.check-box.checked { background: #22c55e; border-color: #22c55e; color: #fff; font-size: 12px; }

.cat-icon { font-size: 20px; flex-shrink: 0; }
.cat-info { flex: 1; min-width: 0; }
.cat-name { font-size: 14px; font-weight: 500; color: #1e293b; }
.cat-meta { font-size: 12px; color: #94a3b8; margin-top: 2px; }
.cat-size { font-size: 14px; font-weight: 600; color: #1e293b; flex-shrink: 0; }
.cat-expand { width: 24px; text-align: center; color: #94a3b8; font-size: 10px; flex-shrink: 0; }

.expanded-panels { border-top: 1px solid #e2e8f0; }
.expand-panel { border-bottom: 1px solid #f1f5f9; }
.expand-panel:last-child { border-bottom: none; }
.expand-header { display: flex; justify-content: space-between; align-items: center; padding: 10px 16px; background: #f8fafc; }
.expand-title { font-size: 12px; color: #64748b; }
.file-list { max-height: 280px; overflow-y: auto; }
.file-item { display: flex; justify-content: space-between; padding: 6px 16px; border-bottom: 1px solid #f8fafc; }
.file-item:last-child { border-bottom: none; }
.file-path { font-size: 12px; color: #475569; word-break: break-all; flex: 1; }
.file-size { font-size: 12px; color: #94a3b8; flex-shrink: 0; margin-left: 12px; }

.empty-state { text-align: center; padding: 60px 0; }
.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-text { font-size: 14px; color: #94a3b8; }

.result-banner {
  display: flex; align-items: center; gap: 16px;
  padding: 16px 20px; border-radius: 12px; margin-bottom: 20px;
}
.result-banner.great { background: linear-gradient(135deg, rgba(34,197,94,0.12) 0%, rgba(34,197,94,0.06) 100%); border: 1px solid rgba(34,197,94,0.25); }
.result-banner.ok { background: linear-gradient(135deg, rgba(59,130,246,0.12) 0%, rgba(59,130,246,0.06) 100%); border: 1px solid rgba(59,130,246,0.25); }
.result-icon { font-size: 28px; }
.result-main { font-size: 15px; color: #1e293b; }
.result-main strong { font-weight: 700; color: #16a34a; }
.result-sub { font-size: 12px; color: #64748b; margin-top: 2px; }

.history-section { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.section-header h3 { font-size: 15px; font-weight: 600; color: #1e293b; margin: 0; }
.icon-btn { background: none; border: none; font-size: 16px; cursor: pointer; padding: 4px; }

.history-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 16px; }
.hstat-card { background: #f8fafc; border-radius: 8px; padding: 14px; text-align: center; }
.hstat-val { font-size: 18px; font-weight: 700; color: #1e293b; }
.hstat-key { font-size: 12px; color: #94a3b8; margin-top: 2px; }

.chart-toggle { display: flex; gap: 4px; margin-bottom: 12px; }
.toggle-btn { padding: 6px 14px; border: 1px solid #e2e8f0; background: #fff; border-radius: 6px; font-size: 12px; color: #64748b; cursor: pointer; transition: all 0.15s ease; }
.toggle-btn.active { background: #0f172a; color: #fff; border-color: #0f172a; }
.chart-wrap { }
</style>
