<template>
  <div class="disk-page">
    <div class="page-header">
      <h2>磁盘管理</h2>
      <p class="page-sub">查看磁盘分区及大文件扫描</p>
    </div>

    <div class="disk-row">
      <div v-for="(d, idx) in disks" :key="idx" class="disk-card">
        <div class="disk-top">
          <span class="disk-drive">{{ d.mountpoint }}</span>
          <el-tag size="small" type="info">{{ d.fstype }}</el-tag>
        </div>
        <div class="disk-gauge" :style="{ '--pct': d.used_percent + '%', '--color': diskColor(d.used_percent) }">
          <svg class="gauge-svg" viewBox="0 0 200 200">
            <circle class="gauge-track" cx="100" cy="100" r="88" />
            <circle class="gauge-fill" cx="100" cy="100" r="88" />
          </svg>
          <div class="gauge-center">
            <div class="gauge-val">{{ Math.round(d.used_percent) }}<span>%</span></div>
          </div>
        </div>
        <div class="disk-meta">
          <div class="meta-row"><span class="meta-key">已用</span><span class="meta-val">{{ formatBytes(d.used) }}</span></div>
          <div class="meta-row"><span class="meta-key">可用</span><span class="meta-val free">{{ formatBytes(d.free) }}</span></div>
          <div class="meta-row"><span class="meta-key">总量</span><span class="meta-val">{{ formatBytes(d.total) }}</span></div>
        </div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card">
        <div class="chart-header"><h3>分区空间对比</h3></div>
        <v-chart :option="barOption" style="height: 280px;" autoresize />
      </div>
      <div class="chart-card">
        <div class="chart-header"><h3>空间占比</h3></div>
        <v-chart :option="pieOption" style="height: 280px;" autoresize />
      </div>
    </div>

    <div class="scan-section">
      <div class="section-header">
        <h3>大文件扫描</h3>
        <div class="header-actions">
          <el-select v-model="scanDrive" size="small" style="width: 90px;">
            <el-option v-for="d in disks" :key="d.mountpoint" :label="d.mountpoint" :value="d.mountpoint" />
          </el-select>
          <el-select v-model="minSize" size="small" style="width: 110px;">
            <el-option :value="10" label="≥ 10 MB" />
            <el-option :value="50" label="≥ 50 MB" />
            <el-option :value="100" label="≥ 100 MB" />
            <el-option :value="500" label="≥ 500 MB" />
            <el-option :value="1024" label="≥ 1 GB" />
          </el-select>
          <div class="search-wrap">
            <span class="search-icon">🔍</span>
            <input v-model="fileKeyword" class="search-input" placeholder="搜索文件..." />
          </div>
          <button class="scan-btn" :class="{ loading: scanning }" :disabled="scanning" @click="handleScan">扫描</button>
        </div>
      </div>

      <table class="file-table">
        <thead>
          <tr>
            <th>#</th>
            <th>文件路径</th>
            <th>大小</th>
            <th>类型</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(f, idx) in filteredFiles" :key="f.path" class="file-row">
            <td class="td-index">{{ idx + 1 }}</td>
            <td class="td-path">{{ f.path }}</td>
            <td class="td-size">{{ formatBytes(f.size) }}</td>
            <td class="td-ext">{{ f.ext || '-' }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="scanResult" class="scan-tip">共找到 {{ scanResult.count }} 个大文件</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { use } from 'echarts/core'
import { BarChart, PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { api, type DiskInfo, type DiskScanResult } from '@/api/backend'

use([BarChart, PieChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const disks = ref<DiskInfo[]>([])
const scanning = ref(false)
const scanDrive = ref('C:\\')
const minSize = ref(50)
const fileKeyword = ref('')
const scanResult = ref<DiskScanResult | null>(null)

const diskColor = (pct: number) => {
  if (pct >= 85) return '#ef4444'
  return '#22c55e'
}

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const toGB = (b: number) => +(b / 1024 / 1024 / 1024).toFixed(1)

const filteredFiles = computed(() => {
  const files = scanResult.value?.files || []
  const kw = fileKeyword.value.toLowerCase()
  if (!kw) return files
  return files.filter(f => f.path.toLowerCase().includes(kw) || (f.ext && f.ext.toLowerCase().includes(kw)))
})

const barOption = computed(() => {
  const labels = disks.value.map(d => d.mountpoint.replace('\\', ''))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['已用', '可用'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: 48, right: 16, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: labels },
    yAxis: { type: 'value', name: 'GB' },
    series: [
      { name: '已用', type: 'bar', stack: 'total', data: disks.value.map(d => toGB(d.used)), itemStyle: { color: '#3b82f6', borderRadius: [3, 3, 0, 0] } },
      { name: '可用', type: 'bar', stack: 'total', data: disks.value.map(d => toGB(d.free)), itemStyle: { color: '#e2e8f0', borderRadius: [3, 3, 0, 0] } }
    ]
  }
})

const pieOption = computed(() => {
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} GB ({d}%)' },
    series: [{
      type: 'pie', radius: ['35%', '65%'], center: ['50%', '50%'],
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{d}%', fontSize: 11 },
      data: disks.value.map((d, i) => ({ value: toGB(d.used), name: d.mountpoint.replace('\\', ''), itemStyle: { color: ['#3b82f6', '#22c55e', '#f59e0b', '#ef4444', '#a855f7', '#14b8a6'][i % 6] } }))
    }]
  }
})

const handleScan = async () => {
  scanning.value = true
  try {
    scanResult.value = await api.scanLargeFiles(scanDrive.value, minSize.value)
  } catch { /* silent */ }
  finally { scanning.value = false }
}

onMounted(async () => {
  try {
    disks.value = await api.getDiskList()
    if (disks.value.length > 0) scanDrive.value = disks.value[0].mountpoint
  } catch { /* silent */ }
})
</script>

<style scoped>
.disk-page { padding: 0 4px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.disk-row { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 16px; margin-bottom: 20px; }

.disk-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; display: flex; flex-direction: column; align-items: center; }
.disk-top { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.disk-drive { font-size: 18px; font-weight: 700; color: #1e293b; }
.disk-gauge { position: relative; width: 120px; height: 120px; margin-bottom: 12px; }
.gauge-center { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; }
.gauge-val { font-size: 28px; font-weight: 700; color: #1e293b; line-height: 1; }
.gauge-val span { font-size: 14px; font-weight: 400; }
.gauge-svg { position: absolute; inset: 0; width: 100%; height: 100%; transform: rotate(-90deg); }
.gauge-track { fill: none; stroke: #ffffff; stroke-width: 10; }
.gauge-fill { fill: none; stroke: var(--color, #22c55e); stroke-width: 10; stroke-linecap: round; stroke-dasharray: 553; stroke-dashoffset: calc(553 - (553 * var(--pct, 0%) / 100)); transition: stroke-dashoffset 1s ease; }

.disk-meta { width: 100%; }
.meta-row { display: flex; justify-content: space-between; padding: 3px 0; font-size: 12px; }
.meta-key { color: #94a3b8; }
.meta-val { color: #1e293b; font-weight: 500; }
.meta-val.free { color: #22c55e; }

.charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
.chart-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px 20px; }
.chart-header { margin-bottom: 12px; }
.chart-header h3 { font-size: 15px; font-weight: 600; color: #1e293b; margin: 0; }

.scan-section { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px 20px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
.section-header h3 { font-size: 15px; font-weight: 600; color: #1e293b; margin: 0; }
.header-actions { display: flex; gap: 8px; align-items: center; }
.search-wrap { display: flex; align-items: center; gap: 6px; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 6px; padding: 0 10px; }
.search-icon { font-size: 12px; color: #94a3b8; }
.search-input { border: none; outline: none; background: transparent; font-size: 12px; color: #1e293b; width: 140px; padding: 6px 0; }

.scan-btn { padding: 7px 16px; border: none; border-radius: 6px; background: #0f172a; color: #fff; font-size: 12px; cursor: pointer; transition: all 0.2s ease; }
.scan-btn:hover:not(:disabled) { background: #1e293b; }
.scan-btn.loading { opacity: 0.7; cursor: not-allowed; }

.file-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.file-table th { background: #f8fafc; padding: 8px 10px; text-align: left; font-weight: 600; color: #64748b; font-size: 12px; border-bottom: 1px solid #e2e8f0; }
.file-table td { padding: 8px 10px; border-bottom: 1px solid #f8fafc; color: #1e293b; }
.file-row:last-child td { border-bottom: none; }
.file-row:hover td { background: #f8fafc; }
.td-index { color: #94a3b8; width: 40px; }
.td-path { max-width: 400px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.td-size { font-weight: 500; white-space: nowrap; }
.td-ext { color: #64748b; }

.scan-tip { margin-top: 8px; font-size: 12px; color: #94a3b8; }
</style>
