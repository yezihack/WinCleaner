<template>
  <div class="memory">
    <h2 style="margin-bottom: 20px; color: #303133;">内存优化</h2>

    <el-row :gutter="16">
      <!-- 当前内存状态 -->
      <el-col :span="8">
        <el-card shadow="hover" class="mem-card">
          <template #header><span>当前内存状态</span></template>
          <el-progress
            type="dashboard"
            :percentage="Math.round(memPercent)"
            :width="160"
            :color="progressColor"
          />
          <p class="stat-text">{{ formatBytes(memUsed) }} / {{ formatBytes(memTotal) }}</p>
          <el-button
            type="primary" size="large"
            :loading="optimizing"
            style="margin-top: 12px;"
            @click="handleOptimize"
          >
            <el-icon><Cpu /></el-icon> 一键优化
          </el-button>
        </el-card>
      </el-col>

      <!-- 优化结果 / 说明 -->
      <el-col :span="8">
        <el-card v-if="optResult" shadow="hover" class="mem-card">
          <template #header><span>本次优化结果</span></template>
          <el-result icon="success" :title="`释放 ${optResult.freed_mb.toFixed(1)} MB`">
            <template #sub-title>
              <p>优化前: {{ optResult.before_percent.toFixed(1) }}%</p>
              <p>优化后: {{ optResult.after_percent.toFixed(1) }}%</p>
            </template>
          </el-result>
        </el-card>
        <el-card v-else shadow="hover" class="mem-card">
          <template #header><span>优化说明</span></template>
          <el-alert type="info" :closable="false" show-icon>
            <p>内存优化会收缩所有进程的工作集，将暂时不用的内存页释放回系统。</p>
            <p style="margin-top: 8px;">不会关闭任何程序，可临时释放物理内存。</p>
          </el-alert>
        </el-card>
      </el-col>

      <!-- 历史统计卡片 -->
      <el-col :span="8">
        <el-card shadow="hover" class="mem-card stat-overview">
          <template #header><span>优化统计</span></template>
          <div class="stat-grid">
            <div class="stat-item">
              <div class="stat-num">{{ optStats.total_count }}</div>
              <div class="stat-desc">累计优化次数</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ optStats.total_freed_mb.toFixed(0) }} MB</div>
              <div class="stat-desc">累计释放内存</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ optStats.last_opt_ago || '暂无' }}</div>
              <div class="stat-desc">距上次优化</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ optStats.last_opt_time || '-' }}</div>
              <div class="stat-desc">上次优化时间</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="16" style="margin-top: 16px;">
      <el-col :span="14">
        <el-card shadow="hover">
          <template #header><span>最近优化趋势</span></template>
          <v-chart :option="lineChartOption" style="height: 300px;" autoresize />
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="hover">
          <template #header><span>每日释放量（近 30 天）</span></template>
          <v-chart :option="dailyBarOption" style="height: 300px;" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px;">
      <el-col :span="14">
        <el-card shadow="hover">
          <template #header><span>月度优化统计</span></template>
          <v-chart :option="monthlyBarOption" style="height: 280px;" autoresize />
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="hover">
          <template #header><span>最近优化前后对比</span></template>
          <v-chart :option="compareChartOption" style="height: 280px;" autoresize />
        </el-card>
      </el-col>
    </el-row>
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

const progressColor = (pct: number) => {
  if (pct < 50) return '#67c23a'
  if (pct < 80) return '#e6a23c'
  return '#f56c6c'
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const lineChartOption = computed(() => {
  const records = optStats.recent_records || []
  const labels = records.map(r => r.date.slice(5) + ' ' + r.time.slice(0, 5))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['释放(MB)', '优化前%', '优化后%'], bottom: 0 },
    grid: { left: 50, right: 20, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 10, rotate: 30 } },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left' },
      { type: 'value', name: '%', position: 'right', max: 100 },
    ],
    series: [
      {
        name: '释放(MB)', type: 'bar', yAxisIndex: 0,
        data: records.map(r => +r.freed_mb.toFixed(1)),
        itemStyle: { color: '#67c23a', borderRadius: [4, 4, 0, 0] },
      },
      {
        name: '优化前%', type: 'line', yAxisIndex: 1, smooth: true,
        data: records.map(r => +r.before_percent.toFixed(1)),
        lineStyle: { color: '#f56c6c', width: 2 }, itemStyle: { color: '#f56c6c' },
      },
      {
        name: '优化后%', type: 'line', yAxisIndex: 1, smooth: true,
        data: records.map(r => +r.after_percent.toFixed(1)),
        lineStyle: { color: '#409eff', width: 2 }, itemStyle: { color: '#409eff' },
      },
    ]
  }
})

const dailyBarOption = computed(() => {
  const daily = optStats.daily_stats || []
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 50, right: 10, top: 16, bottom: 30 },
    xAxis: { type: 'category', data: daily.map(d => d.date.slice(5)), axisLabel: { fontSize: 10, rotate: 30 } },
    yAxis: { type: 'value', name: 'MB' },
    series: [{
      type: 'bar',
      data: daily.map(d => +d.freed_mb.toFixed(1)),
      itemStyle: {
        borderRadius: [4, 4, 0, 0],
        color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [{ offset: 0, color: '#67c23a' }, { offset: 1, color: '#b7eb8f' }]
        }
      },
    }]
  }
})

const monthlyBarOption = computed(() => {
  const monthly = optStats.monthly_stats || []
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['释放(MB)', '次数'], bottom: 0 },
    grid: { left: 50, right: 40, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: monthly.map(m => m.month) },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left' },
      { type: 'value', name: '次', position: 'right' },
    ],
    series: [
      {
        name: '释放(MB)', type: 'bar', yAxisIndex: 0,
        data: monthly.map(m => +m.freed_mb.toFixed(1)),
        itemStyle: { color: '#409eff', borderRadius: [4, 4, 0, 0] },
      },
      {
        name: '次数', type: 'line', yAxisIndex: 1, smooth: true,
        data: monthly.map(m => m.count),
        lineStyle: { color: '#e6a23c', width: 2 }, itemStyle: { color: '#e6a23c' },
      }
    ]
  }
})

const compareChartOption = computed(() => {
  const records = (optStats.recent_records || []).slice(-5)
  const labels = records.map(r => r.date.slice(5) + '\n' + r.time.slice(0, 5))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['优化前', '优化后'], bottom: 0 },
    grid: { left: 40, right: 10, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 10 } },
    yAxis: { type: 'value', name: '%', max: 100 },
    series: [
      {
        name: '优化前', type: 'bar',
        data: records.map(r => +r.before_percent.toFixed(1)),
        itemStyle: { color: '#f56c6c', borderRadius: [4, 4, 0, 0] },
      },
      {
        name: '优化后', type: 'bar',
        data: records.map(r => +r.after_percent.toFixed(1)),
        itemStyle: { color: '#67c23a', borderRadius: [4, 4, 0, 0] },
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
.mem-card { text-align: center; min-height: 300px; }
.stat-text { margin-top: 8px; color: #909399; font-size: 14px; }
.stat-overview { text-align: left; }
.stat-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; padding: 8px; }
.stat-item { text-align: center; }
.stat-num { font-size: 20px; font-weight: 700; color: #409eff; }
.stat-desc { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
