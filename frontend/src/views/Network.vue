<template>
  <div class="network-page">
    <div class="page-header">
      <h2>流量监控</h2>
      <p class="page-sub">实时监控网络上传下载速度及流量统计</p>
    </div>

    <div class="stats-row">
      <div class="stat-card up">
        <div class="stat-icon">↑</div>
        <div class="stat-info">
          <div class="stat-val">{{ formatSpeed(traffic.overview.up_speed) }}</div>
          <div class="stat-key">实时上传</div>
        </div>
      </div>
      <div class="stat-card down">
        <div class="stat-icon">↓</div>
        <div class="stat-info">
          <div class="stat-val">{{ formatSpeed(traffic.overview.down_speed) }}</div>
          <div class="stat-key">实时下载</div>
        </div>
      </div>
      <div class="stat-card sent">
        <div class="stat-icon">📤</div>
        <div class="stat-info">
          <div class="stat-val">{{ formatBytes(netStats.total_sent) }}</div>
          <div class="stat-key">历史总发送</div>
        </div>
      </div>
      <div class="stat-card recv">
        <div class="stat-icon">📥</div>
        <div class="stat-info">
          <div class="stat-val">{{ formatBytes(netStats.total_recv) }}</div>
          <div class="stat-key">历史总接收</div>
        </div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card wide">
        <div class="chart-header">
          <h3>流量趋势</h3>
          <div class="toggle-group">
            <button v-for="m in ['daily', 'monthly', 'yearly']" :key="m" class="toggle-btn" :class="{ active: chartMode === m }" @click="chartMode = m as any">{{ m === 'daily' ? '近30天' : m === 'monthly' ? '按月' : '按年' }}</button>
          </div>
        </div>
        <v-chart :option="trendOption" style="height: 300px;" autoresize />
      </div>
      <div class="chart-card">
        <div class="chart-header">
          <h3>上传/下载占比</h3>
        </div>
        <v-chart :option="pieOption" style="height: 300px;" autoresize />
      </div>
    </div>

    <div class="chart-card" style="margin-bottom: 20px;">
      <div class="chart-header">
        <h3>每日流量对比（近 14 天）</h3>
      </div>
      <v-chart :option="barOption" style="height: 260px;" autoresize />
    </div>

    <div class="app-section">
      <div class="section-header">
        <h3>应用流量监控</h3>
        <div class="header-actions">
          <div class="search-wrap">
            <span class="search-icon">🔍</span>
            <input v-model="keyword" class="search-input" placeholder="搜索应用..." />
          </div>
          <button class="icon-btn" :class="{ loading }" @click="refresh">🔄</button>
        </div>
      </div>
      <table class="app-table">
        <thead>
          <tr>
            <th>应用名称</th>
            <th>进程数</th>
            <th>发送</th>
            <th>接收</th>
            <th>总计</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in filteredProcesses" :key="row.name" class="app-row">
            <td class="td-name">{{ row.name }}</td>
            <td>{{ row.count }}</td>
            <td class="traffic-up">↑ {{ formatBytes(row.sent) }}</td>
            <td class="traffic-down">↓ {{ formatBytes(row.recv) }}</td>
            <td class="td-total">{{ formatBytes(row.sent + row.recv) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { use } from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent, DataZoomComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { api, type NetTrafficResult, type NetTrafficStats } from '@/api/backend'

use([BarChart, LineChart, PieChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, DataZoomComponent, CanvasRenderer])

const loading = ref(false)
const keyword = ref('')
const chartMode = ref<'daily' | 'monthly' | 'yearly'>('daily')

const traffic = reactive<NetTrafficResult>({ overview: { total_sent: 0, total_recv: 0, up_speed: 0, down_speed: 0 }, processes: [] })
const netStats = reactive<NetTrafficStats>({ daily_stats: [], monthly_stats: [], yearly_stats: [], total_sent: 0, total_recv: 0 })
let timer: ReturnType<typeof setInterval> | null = null

const filteredProcesses = computed(() => {
  const kw = keyword.value.toLowerCase()
  const list = traffic.processes || []
  if (!kw) return list
  return list.filter(p => p.name.toLowerCase().includes(kw))
})

const toMB = (b: number) => +(b / 1024 / 1024).toFixed(2)

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatSpeed = (bps: number): string => {
  if (!bps || bps === 0) return '0 B/s'
  if (bps < 1024) return bps + ' B/s'
  if (bps < 1024 * 1024) return (bps / 1024).toFixed(1) + ' KB/s'
  return (bps / 1024 / 1024).toFixed(1) + ' MB/s'
}

const trendOption = computed(() => {
  let data: { label: string; sent: number; recv: number }[] = []
  if (chartMode.value === 'daily') data = (netStats.daily_stats || []).map(d => ({ label: d.date.slice(5), sent: d.sent, recv: d.recv }))
  else if (chartMode.value === 'monthly') data = (netStats.monthly_stats || []).map(d => ({ label: d.month, sent: d.sent, recv: d.recv }))
  else data = (netStats.yearly_stats || []).map(d => ({ label: d.year, sent: d.sent, recv: d.recv }))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['上传', '下载'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 16, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: data.map(d => d.label), axisLabel: { fontSize: 10 } },
    yAxis: { type: 'value', name: 'MB', axisLabel: { fontSize: 10 } },
    series: [
      { name: '上传', type: 'line', smooth: true, data: data.map(d => toMB(d.sent)), areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(34,197,94,0.3)' }, { offset: 1, color: 'rgba(34,197,94,0.02)' }] } }, lineStyle: { color: '#22c55e', width: 2 }, itemStyle: { color: '#22c55e' } },
      { name: '下载', type: 'line', smooth: true, data: data.map(d => toMB(d.recv)), areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(59,130,246,0.3)' }, { offset: 1, color: 'rgba(59,130,246,0.02)' }] } }, lineStyle: { color: '#3b82f6', width: 2 }, itemStyle: { color: '#3b82f6' } }
    ]
  }
})

const pieOption = computed(() => {
  const sent = netStats.total_sent || 0, recv = netStats.total_recv || 0
  return {
    tooltip: { trigger: 'item', formatter: (p: any) => `${p.name}: ${formatBytes(p.value * 1024 * 1024)} (${p.percent}%)` },
    legend: { bottom: 0, data: ['上传', '下载'] },
    series: [{ type: 'pie', radius: ['40%', '70%'], center: ['50%', '45%'], itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 2 }, label: { show: true, formatter: '{b}\n{d}%', fontSize: 12 }, data: [{ value: toMB(sent), name: '上传', itemStyle: { color: '#22c55e' } }, { value: toMB(recv), name: '下载', itemStyle: { color: '#3b82f6' } }] }]
  }
})

const barOption = computed(() => {
  const daily = (netStats.daily_stats || []).slice(-14)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['上传', '下载'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 16, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: daily.map(d => d.date.slice(5)), axisLabel: { fontSize: 10 } },
    yAxis: { type: 'value', name: 'MB', axisLabel: { fontSize: 10 } },
    series: [
      { name: '上传', type: 'bar', stack: 'total', data: daily.map(d => toMB(d.sent)), itemStyle: { color: '#22c55e', borderRadius: [3, 3, 0, 0] } },
      { name: '下载', type: 'bar', stack: 'total', data: daily.map(d => toMB(d.recv)), itemStyle: { color: '#3b82f6', borderRadius: [3, 3, 0, 0] } }
    ]
  }
})

const refresh = async () => {
  loading.value = true
  try {
    const data = await api.getNetTraffic()
    traffic.overview = data.overview || traffic.overview
    traffic.processes = data.processes || []
  } catch { /* silent */ }
  finally { loading.value = false }
}

const loadStats = async () => {
  try {
    const data = await api.getNetTrafficStats()
    Object.assign(netStats, data)
  } catch { /* silent */ }
}

onMounted(() => { refresh(); loadStats(); timer = setInterval(refresh, 5000) })
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.network-page { padding: 0 4px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 20px; }
.stat-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px; display: flex; align-items: center; gap: 12px; }
.stat-card.up { border-left: 3px solid #22c55e; }
.stat-card.down { border-left: 3px solid #3b82f6; }
.stat-card.sent { border-left: 3px solid #f59e0b; }
.stat-card.recv { border-left: 3px solid #a855f7; }
.stat-icon { font-size: 24px; }
.stat-val { font-size: 18px; font-weight: 700; color: #1e293b; }
.stat-key { font-size: 12px; color: #94a3b8; margin-top: 2px; }

.charts-row { display: grid; grid-template-columns: 2fr 1fr; gap: 16px; margin-bottom: 16px; }
.chart-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px 20px; }
.chart-card.wide {}
.chart-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.chart-header h3 { font-size: 15px; font-weight: 600; color: #1e293b; margin: 0; }
.toggle-group { display: flex; gap: 4px; }
.toggle-btn { padding: 5px 12px; border: 1px solid #e2e8f0; background: #fff; border-radius: 6px; font-size: 12px; color: #64748b; cursor: pointer; transition: all 0.15s ease; }
.toggle-btn.active { background: #0f172a; color: #fff; border-color: #0f172a; }

.app-section { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px 20px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
.section-header h3 { font-size: 15px; font-weight: 600; color: #1e293b; margin: 0; }
.header-actions { display: flex; gap: 8px; align-items: center; }
.search-wrap { display: flex; align-items: center; gap: 6px; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 6px; padding: 0 10px; }
.search-icon { font-size: 12px; color: #94a3b8; }
.search-input { border: none; outline: none; background: transparent; font-size: 12px; color: #1e293b; width: 140px; padding: 6px 0; }
.icon-btn { background: none; border: 1px solid #e2e8f0; border-radius: 6px; font-size: 14px; cursor: pointer; padding: 6px 8px; transition: all 0.15s ease; }
.icon-btn:hover { border-color: #3b82f6; }
.icon-btn.loading { opacity: 0.6; }

.app-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.app-table th { background: #f8fafc; padding: 8px 12px; text-align: left; font-weight: 600; color: #64748b; font-size: 12px; border-bottom: 1px solid #e2e8f0; }
.app-table td { padding: 9px 12px; border-bottom: 1px solid #f8fafc; color: #1e293b; }
.app-row:last-child td { border-bottom: none; }
.app-row:hover td { background: #f8fafc; }
.td-name { font-weight: 500; }
.td-total { font-weight: 600; color: #1e293b; }
.traffic-up { color: #22c55e; }
.traffic-down { color: #3b82f6; }
</style>
