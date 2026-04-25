<template>
  <div class="memory-page">
    <div class="page-header">
      <h2>内存优化</h2>
      <p class="page-sub">智能收缩进程工作集，释放物理内存</p>
    </div>

    <div class="main-gauge-wrap">
      <div class="gauge-container">
        <div class="gauge-ring" :style="{ '--pct': memPercent + '%', '--color': gaugeColor }">
          <div class="gauge-inner">
            <div class="gauge-value">{{ memPercent.toFixed(0) }}<span class="gauge-unit">%</span></div>
            <div class="gauge-label">内存占用</div>
          </div>
          <svg class="gauge-svg" viewBox="0 0 200 200">
            <circle class="gauge-track" cx="100" cy="100" r="88" />
            <circle class="gauge-fill" cx="100" cy="100" r="88" />
          </svg>
        </div>
        <div class="gauge-stats">
          <div class="gstat">
            <span class="gstat-val">{{ formatBytes(memUsed) }}</span>
            <span class="gstat-key">已用</span>
          </div>
          <div class="gstat-divider"></div>
          <div class="gstat">
            <span class="gstat-val">{{ formatBytes(memTotal - memUsed) }}</span>
            <span class="gstat-key">可用</span>
          </div>
          <div class="gstat-divider"></div>
          <div class="gstat">
            <span class="gstat-val">{{ formatBytes(memTotal) }}</span>
            <span class="gstat-key">总量</span>
          </div>
        </div>
      </div>

      <div class="action-panel">
        <button class="optimize-btn" :class="{ loading: optimizing }" :disabled="optimizing" @click="handleOptimize">
          <span class="btn-icon">⚡</span>
          <span class="btn-text">{{ optimizing ? '优化中...' : '一键优化' }}</span>
        </button>
        <p class="action-tip">安全收缩所有进程工作集，不关闭任何程序</p>
      </div>
    </div>

    <div v-if="optResult" class="result-banner" :class="optResult.freed_mb > 50 ? 'good' : 'ok'">
      <div class="result-icon">🎉</div>
      <div class="result-info">
        <div class="result-main">本次释放 <strong>{{ optResult.freed_mb.toFixed(1) }} MB</strong></div>
        <div class="result-sub">优化前 {{ optResult.before_percent.toFixed(1) }}% → 优化后 {{ optResult.after_percent.toFixed(1) }}%</div>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-value">{{ optStats.total_count }}</div>
          <div class="stat-label">累计优化次数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">💾</div>
        <div class="stat-content">
          <div class="stat-value">{{ optStats.total_freed_mb.toFixed(0) }} MB</div>
          <div class="stat-label">累计释放内存</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">⏱️</div>
        <div class="stat-content">
          <div class="stat-value">{{ optStats.last_opt_ago || '暂无' }}</div>
          <div class="stat-label">距上次优化</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🕐</div>
        <div class="stat-content">
          <div class="stat-value">{{ optStats.last_opt_time || '-' }}</div>
          <div class="stat-label">上次优化时间</div>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>优化趋势</h3>
        </div>
        <v-chart :option="trendOption" style="height: 260px;" autoresize />
      </div>
      <div class="chart-card">
        <div class="chart-header">
          <h3>每日释放量</h3>
          <span class="chart-sub">近 30 天</span>
        </div>
        <v-chart :option="dailyOption" style="height: 260px;" autoresize />
      </div>
      <div class="chart-card wide">
        <div class="chart-header">
          <h3>月度统计</h3>
        </div>
        <v-chart :option="monthlyOption" style="height: 240px;" autoresize />
      </div>
      <div class="chart-card wide">
        <div class="chart-header">
          <h3>最近对比</h3>
          <span class="chart-sub">最近 5 次</span>
        </div>
        <v-chart :option="compareOption" style="height: 240px;" autoresize />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { use } from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, GridComponent, LegendComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { api, type MemoryOptResult, type MemOptStats } from '@/api/backend'

use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const memTotal = ref(0)
const memUsed = ref(0)
const memPercent = ref(0)
const optimizing = ref(false)
const optResult = ref<MemoryOptResult | null>(null)
const optStats = reactive<MemOptStats>({
  recent_records: null, daily_stats: null, monthly_stats: null,
  last_opt_time: '', last_opt_ago: '', total_freed_mb: 0, total_count: 0,
})

const gaugeColor = computed(() => {
  const p = memPercent.value
  if (p >= 85) return '#ef4444'
  return '#22c55e'
})

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const trendOption = computed(() => {
  const records = optStats.recent_records || []
  const labels = records.map(r => r.date.slice(5) + ' ' + r.time.slice(0, 5))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['释放(MB)', '优化前%', '优化后%'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 16, top: 12, bottom: 40 },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 10, rotate: 30 } },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left', axisLabel: { fontSize: 10 } },
      { type: 'value', name: '%', position: 'right', max: 100, axisLabel: { fontSize: 10 } },
    ],
    series: [
      {
        name: '释放(MB)', type: 'bar', yAxisIndex: 0,
        data: records.map(r => +r.freed_mb.toFixed(1)),
        itemStyle: { color: '#22c55e', borderRadius: [3, 3, 0, 0] },
      },
      {
        name: '优化前%', type: 'line', yAxisIndex: 1, smooth: true,
        data: records.map(r => +r.before_percent.toFixed(1)),
        lineStyle: { color: '#ef4444', width: 2 }, itemStyle: { color: '#ef4444' },
      },
      {
        name: '优化后%', type: 'line', yAxisIndex: 1, smooth: true,
        data: records.map(r => +r.after_percent.toFixed(1)),
        lineStyle: { color: '#3b82f6', width: 2 }, itemStyle: { color: '#3b82f6' },
      },
    ]
  }
})

const dailyOption = computed(() => {
  const daily = optStats.daily_stats || []
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 48, right: 10, top: 12, bottom: 28 },
    xAxis: { type: 'category', data: daily.map(d => d.date.slice(5)), axisLabel: { fontSize: 10, rotate: 30 } },
    yAxis: { type: 'value', name: 'MB', axisLabel: { fontSize: 10 } },
    series: [{
      type: 'bar',
      data: daily.map(d => +d.freed_mb.toFixed(1)),
      itemStyle: {
        borderRadius: [3, 3, 0, 0],
        color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [{ offset: 0, color: '#22c55e' }, { offset: 1, color: '#86efac' }]
        }
      },
    }]
  }
})

const monthlyOption = computed(() => {
  const monthly = optStats.monthly_stats || []
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['释放(MB)', '次数'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 40, top: 12, bottom: 40 },
    xAxis: { type: 'category', data: monthly.map(m => m.month) },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left', axisLabel: { fontSize: 10 } },
      { type: 'value', name: '次', position: 'right', axisLabel: { fontSize: 10 } },
    ],
    series: [
      {
        name: '释放(MB)', type: 'bar', yAxisIndex: 0,
        data: monthly.map(m => +m.freed_mb.toFixed(1)),
        itemStyle: { color: '#3b82f6', borderRadius: [3, 3, 0, 0] },
      },
      {
        name: '次数', type: 'line', yAxisIndex: 1, smooth: true,
        data: monthly.map(m => m.count),
        lineStyle: { color: '#f59e0b', width: 2 }, itemStyle: { color: '#f59e0b' },
      }
    ]
  }
})

const compareOption = computed(() => {
  const records = (optStats.recent_records || []).slice(-5)
  const labels = records.map(r => r.date.slice(5) + '\n' + r.time.slice(0, 5))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['优化前', '优化后'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 40, right: 10, top: 12, bottom: 40 },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 10 } },
    yAxis: { type: 'value', name: '%', max: 100, axisLabel: { fontSize: 10 } },
    series: [
      {
        name: '优化前', type: 'bar',
        data: records.map(r => +r.before_percent.toFixed(1)),
        itemStyle: { color: '#ef4444', borderRadius: [3, 3, 0, 0] },
      },
      {
        name: '优化后', type: 'bar',
        data: records.map(r => +r.after_percent.toFixed(1)),
        itemStyle: { color: '#22c55e', borderRadius: [3, 3, 0, 0] },
      }
    ]
  }
})

const refreshMemory = async () => {
  try {
    const info = await api.getSystemInfo()
    memTotal.value = info.mem_total
    memUsed.value = info.mem_used
    memPercent.value = info.mem_percent
  } catch { /* silent */ }
}

const loadStats = async () => {
  try {
    const data = await api.getMemOptStats()
    Object.assign(optStats, data)
  } catch { /* silent */ }
}

const handleOptimize = async () => {
  optimizing.value = true
  try {
    optResult.value = await api.optimizeMemory()
    ElMessage.success(`释放了 ${optResult.value.freed_mb.toFixed(1)} MB 内存`)
    await refreshMemory()
    await loadStats()
  } catch {
    ElMessage.error('内存优化失败')
  } finally {
    optimizing.value = false
  }
}

onMounted(() => {
  refreshMemory()
  loadStats()
})
</script>

<style scoped>
.memory-page {
  padding: 0 4px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  font-size: 22px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
}

.page-sub {
  font-size: 13px;
  color: #64748b;
  margin: 4px 0 0;
}

.main-gauge-wrap {
  display: flex;
  align-items: center;
  gap: 48px;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border-radius: 16px;
  padding: 32px 40px;
  margin-bottom: 20px;
}

.gauge-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.gauge-ring {
  position: relative;
  width: 180px;
  height: 180px;
}

.gauge-inner {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.gauge-value {
  font-size: 42px;
  font-weight: 700;
  color: #f8fafc;
  line-height: 1;
}

.gauge-unit {
  font-size: 20px;
  font-weight: 400;
}

.gauge-label {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
}

.gauge-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.gauge-track {
  fill: none;
  stroke: #ffffff;
  stroke-width: 12;
}

.gauge-fill {
  fill: none;
  stroke: var(--color, #22c55e);
  stroke-width: 12;
  stroke-linecap: round;
  stroke-dasharray: 553;
  stroke-dashoffset: calc(553 - (553 * var(--pct, 0%) / 100));
  transition: stroke-dashoffset 1s ease, stroke 0.5s ease;
}

.gauge-stats {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255,255,255,0.05);
  border-radius: 8px;
  padding: 10px 20px;
}

.gstat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.gstat-val {
  font-size: 14px;
  font-weight: 600;
  color: #f8fafc;
}

.gstat-key {
  font-size: 11px;
  color: #64748b;
}

.gstat-divider {
  width: 1px;
  height: 28px;
  background: #334155;
}

.action-panel {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
}

.optimize-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 32px;
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  border: none;
  border-radius: 12px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(34, 197, 94, 0.35);
}

.optimize-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 28px rgba(34, 197, 94, 0.5);
}

.optimize-btn:active:not(:disabled) {
  transform: translateY(0);
}

.optimize-btn.loading {
  opacity: 0.8;
  cursor: not-allowed;
}

.optimize-btn:disabled {
  background: linear-gradient(135deg, #64748b 0%, #475569 100%);
  box-shadow: none;
}

.btn-icon {
  font-size: 20px;
}

.action-tip {
  font-size: 12px;
  color: #94a3b8;
  margin: 0;
}

.result-banner {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.result-banner.good {
  background: linear-gradient(135deg, rgba(34,197,94,0.12) 0%, rgba(34,197,94,0.06) 100%);
  border: 1px solid rgba(34,197,94,0.25);
}

.result-banner.ok {
  background: linear-gradient(135deg, rgba(59,130,246,0.12) 0%, rgba(59,130,246,0.06) 100%);
  border: 1px solid rgba(59,130,246,0.25);
}

.result-icon {
  font-size: 28px;
}

.result-main {
  font-size: 15px;
  color: #1e293b;
}

.result-main strong {
  font-weight: 700;
  color: #16a34a;
}

.result-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px 18px;
  transition: box-shadow 0.2s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 16px rgba(0,0,0,0.06);
}

.stat-icon {
  font-size: 24px;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #1e293b;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.chart-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px 20px;
}

.chart-card.wide {
  grid-column: span 2;
}

.chart-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 12px;
}

.chart-header h3 {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.chart-sub {
  font-size: 12px;
  color: #94a3b8;
}
</style>
