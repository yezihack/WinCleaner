<template>
  <div class="network">
    <h2 style="margin-bottom: 20px; color: #303133;">流量监控</h2>

    <!-- 实时总览卡片 -->
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card gradient-up">
          <div class="stat-icon">↑</div>
          <div class="stat-value">{{ formatSpeed(traffic.overview.up_speed) }}</div>
          <div class="stat-label">实时上传</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card gradient-down">
          <div class="stat-icon">↓</div>
          <div class="stat-value">{{ formatSpeed(traffic.overview.down_speed) }}</div>
          <div class="stat-label">实时下载</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card gradient-sent">
          <div class="stat-icon">📤</div>
          <div class="stat-value">{{ formatBytes(netStats.total_sent) }}</div>
          <div class="stat-label">历史总发送</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card gradient-recv">
          <div class="stat-icon">📥</div>
          <div class="stat-value">{{ formatBytes(netStats.total_recv) }}</div>
          <div class="stat-label">历史总接收</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="16" style="margin-top: 16px;">
      <!-- 近30天趋势（面积图） -->
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>流量趋势</span>
              <el-radio-group v-model="chartMode" size="small">
                <el-radio-button value="daily">近30天</el-radio-button>
                <el-radio-button value="monthly">按月</el-radio-button>
                <el-radio-button value="yearly">按年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <v-chart :option="trendChartOption" style="height: 320px;" autoresize />
        </el-card>
      </el-col>

      <!-- 上传/下载占比（饼图） -->
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>上传 / 下载占比</span></template>
          <v-chart :option="pieChartOption" style="height: 320px;" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 每日对比柱状图 -->
    <el-card shadow="hover" style="margin-top: 16px;">
      <template #header><span>每日流量对比（近 14 天）</span></template>
      <v-chart :option="barChartOption" style="height: 280px;" autoresize />
    </el-card>

    <!-- 应用流量监控 -->
    <el-card shadow="hover" style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <span>应用流量监控</span>
          <div class="header-actions">
            <el-input
              v-model="keyword"
              placeholder="搜索应用..."
              clearable
              size="small"
              style="width: 200px;"
              :prefix-icon="Search"
            />
            <el-button size="small" :loading="loading" @click="refresh">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        :data="filteredProcesses"
        :default-sort="{ prop: 'total', order: 'descending' }"
        max-height="360"
        stripe
      >
        <el-table-column prop="name" label="应用名称" min-width="160" sortable show-overflow-tooltip />
        <el-table-column prop="count" label="进程数" width="90" sortable />
        <el-table-column label="发送" width="140" sortable :sort-by="(row: any) => row.sent">
          <template #default="{ row }">
            <span class="traffic-up">↑ {{ formatBytes(row.sent) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="接收" width="140" sortable :sort-by="(row: any) => row.recv">
          <template #default="{ row }">
            <span class="traffic-down">↓ {{ formatBytes(row.recv) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="总计" prop="total" width="140" sortable>
          <template #default="{ row }">
            {{ formatBytes(row.sent + row.recv) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, GridComponent,
  LegendComponent, DataZoomComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { api, type NetTrafficResult, type NetTrafficStats } from '@/api/backend'

use([
  BarChart, LineChart, PieChart, TitleComponent, TooltipComponent,
  GridComponent, LegendComponent, DataZoomComponent, CanvasRenderer
])

const loading = ref(false)
const keyword = ref('')
const chartMode = ref<'daily' | 'monthly' | 'yearly'>('daily')

const traffic = reactive<NetTrafficResult>({
  overview: { total_sent: 0, total_recv: 0, up_speed: 0, down_speed: 0 },
  processes: [],
})

const netStats = reactive<NetTrafficStats>({
  daily_stats: [], monthly_stats: [], yearly_stats: [],
  total_sent: 0, total_recv: 0,
})

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

// 趋势面积图
const trendChartOption = computed(() => {
  let data: { label: string; sent: number; recv: number }[] = []

  if (chartMode.value === 'daily') {
    data = (netStats.daily_stats || []).map(d => ({ label: d.date.slice(5), sent: d.sent, recv: d.recv }))
  } else if (chartMode.value === 'monthly') {
    data = (netStats.monthly_stats || []).map(d => ({ label: d.month, sent: d.sent, recv: d.recv }))
  } else {
    data = (netStats.yearly_stats || []).map(d => ({ label: d.year, sent: d.sent, recv: d.recv }))
  }

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const s = params[0]; const r = params[1]
        return `${s.axisValue}<br/>${s.marker} 上传: ${formatBytes(s.value * 1024 * 1024)}<br/>${r.marker} 下载: ${formatBytes(r.value * 1024 * 1024)}`
      }
    },
    legend: { data: ['上传', '下载'], bottom: 0 },
    grid: { left: 50, right: 20, top: 20, bottom: 40 },
    xAxis: { type: 'category', data: data.map(d => d.label), axisLabel: { fontSize: 11 } },
    yAxis: { type: 'value', name: 'MB', axisLabel: { fontSize: 11 } },
    series: [
      {
        name: '上传', type: 'line', smooth: true,
        data: data.map(d => toMB(d.sent)),
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(103,194,58,0.4)' }, { offset: 1, color: 'rgba(103,194,58,0.05)' }] } },
        lineStyle: { color: '#67c23a', width: 2 },
        itemStyle: { color: '#67c23a' },
      },
      {
        name: '下载', type: 'line', smooth: true,
        data: data.map(d => toMB(d.recv)),
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(64,158,255,0.4)' }, { offset: 1, color: 'rgba(64,158,255,0.05)' }] } },
        lineStyle: { color: '#409eff', width: 2 },
        itemStyle: { color: '#409eff' },
      }
    ]
  }
})

// 饼图
const pieChartOption = computed(() => {
  const sent = netStats.total_sent || 0
  const recv = netStats.total_recv || 0
  return {
    tooltip: {
      trigger: 'item',
      formatter: (p: any) => `${p.name}: ${formatBytes(p.value * 1024 * 1024)} (${p.percent}%)`
    },
    legend: { bottom: 0, data: ['上传', '下载'] },
    series: [{
      type: 'pie', radius: ['40%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{d}%', fontSize: 12 },
      data: [
        { value: toMB(sent), name: '上传', itemStyle: { color: '#67c23a' } },
        { value: toMB(recv), name: '下载', itemStyle: { color: '#409eff' } },
      ]
    }]
  }
})

// 每日对比柱状图（近14天）
const barChartOption = computed(() => {
  const daily = (netStats.daily_stats || []).slice(-14)
  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const s = params[0]; const r = params[1]
        return `${s.axisValue}<br/>${s.marker} 上传: ${formatBytes(s.value * 1024 * 1024)}<br/>${r.marker} 下载: ${formatBytes(r.value * 1024 * 1024)}`
      }
    },
    legend: { data: ['上传', '下载'], bottom: 0 },
    grid: { left: 50, right: 20, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: daily.map(d => d.date.slice(5)), axisLabel: { fontSize: 11 } },
    yAxis: { type: 'value', name: 'MB', axisLabel: { fontSize: 11 } },
    series: [
      {
        name: '上传', type: 'bar', stack: 'total',
        data: daily.map(d => toMB(d.sent)),
        itemStyle: { color: '#67c23a', borderRadius: [4, 4, 0, 0] },
      },
      {
        name: '下载', type: 'bar', stack: 'total',
        data: daily.map(d => toMB(d.recv)),
        itemStyle: { color: '#409eff', borderRadius: [4, 4, 0, 0] },
      }
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

onMounted(() => {
  refresh()
  loadStats()
  timer = setInterval(refresh, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.stat-card {
  text-align: center;
  padding: 12px 0;
  border-radius: 12px;
  position: relative;
  overflow: hidden;
}
.stat-card .stat-icon {
  font-size: 24px;
  margin-bottom: 4px;
}
.stat-card .stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}
.stat-card .stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.gradient-up { background: linear-gradient(135deg, #f0fff0 0%, #e8fce8 100%); }
.gradient-down { background: linear-gradient(135deg, #f0f7ff 0%, #e0efff 100%); }
.gradient-sent { background: linear-gradient(135deg, #fff8f0 0%, #fff0e0 100%); }
.gradient-recv { background: linear-gradient(135deg, #f5f0ff 0%, #ece0ff 100%); }
.gradient-up .stat-value { color: #67c23a; }
.gradient-down .stat-value { color: #409eff; }
.gradient-sent .stat-value { color: #e6a23c; }
.gradient-recv .stat-value { color: #9b59b6; }

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.traffic-up { color: #67c23a; }
.traffic-down { color: #409eff; }
</style>
