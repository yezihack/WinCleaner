<template>
  <div class="disk">
    <h2 style="margin-bottom: 20px; color: #303133;">磁盘管理</h2>

    <!-- 分区概览卡片 -->
    <el-row :gutter="16">
      <el-col :span="8" v-for="(d, idx) in disks" :key="idx">
        <el-card shadow="hover" class="disk-card">
          <div class="disk-header">
            <span class="disk-drive">{{ d.mountpoint }}</span>
            <el-tag size="small" type="info">{{ d.fstype }}</el-tag>
          </div>
          <el-progress
            type="dashboard"
            :percentage="Math.round(d.used_percent)"
            :width="120"
            :color="diskColor(d.used_percent)"
          />
          <div class="disk-detail">
            <span>已用 {{ formatBytes(d.used) }}</span>
            <span>共 {{ formatBytes(d.total) }}</span>
          </div>
          <div class="disk-free">可用 {{ formatBytes(d.free) }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 磁盘使用对比图 -->
    <el-row :gutter="16" style="margin-top: 16px;">
      <el-col :span="14">
        <el-card shadow="hover">
          <template #header><span>分区空间对比</span></template>
          <v-chart :option="barChartOption" style="height: 300px;" autoresize />
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="hover">
          <template #header><span>空间占比</span></template>
          <v-chart :option="pieChartOption" style="height: 300px;" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 大文件扫描 -->
    <el-card shadow="hover" style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <span>大文件扫描</span>
          <div class="header-actions">
            <el-select v-model="scanDrive" size="small" style="width: 100px;">
              <el-option
                v-for="d in disks" :key="d.mountpoint"
                :label="d.mountpoint" :value="d.mountpoint"
              />
            </el-select>
            <el-select v-model="minSize" size="small" style="width: 120px;">
              <el-option :value="10" label="≥ 10 MB" />
              <el-option :value="50" label="≥ 50 MB" />
              <el-option :value="100" label="≥ 100 MB" />
              <el-option :value="500" label="≥ 500 MB" />
              <el-option :value="1024" label="≥ 1 GB" />
            </el-select>
            <el-input
              v-model="fileKeyword"
              placeholder="搜索文件..."
              clearable size="small"
              style="width: 180px;"
              :prefix-icon="Search"
            />
            <el-button type="primary" size="small" :loading="scanning" @click="handleScan">
              扫描
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        :data="filteredFiles"
        :default-sort="{ prop: 'size', order: 'descending' }"
        max-height="400" stripe
      >
        <el-table-column type="index" width="50" />
        <el-table-column prop="path" label="文件路径" min-width="360" show-overflow-tooltip sortable />
        <el-table-column label="大小" width="120" sortable :sort-by="(row: any) => row.size">
          <template #default="{ row }">{{ formatBytes(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="ext" label="类型" width="100" sortable />
      </el-table>

      <div v-if="scanResult" style="margin-top: 8px; color: #909399; font-size: 12px;">
        共找到 {{ scanResult.count }} 个大文件
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { BarChart, PieChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, GridComponent, LegendComponent
} from 'echarts/components'
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
  if (pct < 60) return '#67c23a'
  if (pct < 85) return '#e6a23c'
  return '#f56c6c'
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
  return files.filter(f => f.path.toLowerCase().includes(kw) || f.ext.toLowerCase().includes(kw))
})

const barChartOption = computed(() => {
  const labels = disks.value.map(d => d.mountpoint.replace('\\', ''))
  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const u = params[0]; const f = params[1]
        return `${u.axisValue}<br/>${u.marker} 已用: ${u.value} GB<br/>${f.marker} 可用: ${f.value} GB`
      }
    },
    legend: { data: ['已用', '可用'], bottom: 0 },
    grid: { left: 50, right: 20, top: 16, bottom: 40 },
    xAxis: { type: 'category', data: labels },
    yAxis: { type: 'value', name: 'GB' },
    series: [
      {
        name: '已用', type: 'bar', stack: 'total',
        data: disks.value.map(d => toGB(d.used)),
        itemStyle: { color: '#409eff', borderRadius: [0, 0, 0, 0] },
      },
      {
        name: '可用', type: 'bar', stack: 'total',
        data: disks.value.map(d => toGB(d.free)),
        itemStyle: { color: '#e4e7ed', borderRadius: [4, 4, 0, 0] },
      }
    ]
  }
})

const pieChartOption = computed(() => {
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} GB ({d}%)' },
    series: [{
      type: 'pie', radius: ['35%', '65%'],
      center: ['50%', '50%'],
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{d}%', fontSize: 11 },
      data: disks.value.map((d, i) => ({
        value: toGB(d.used),
        name: d.mountpoint.replace('\\', ''),
        itemStyle: { color: ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#9b59b6', '#1abc9c'][i % 6] }
      }))
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
    if (disks.value.length > 0) {
      scanDrive.value = disks.value[0].mountpoint
    }
  } catch { /* silent */ }
})
</script>

<style scoped>
.disk-card { text-align: center; padding: 8px 0; }
.disk-header { display: flex; justify-content: center; align-items: center; gap: 8px; margin-bottom: 8px; }
.disk-drive { font-size: 18px; font-weight: 700; color: #303133; }
.disk-detail { display: flex; justify-content: space-between; padding: 0 16px; margin-top: 8px; font-size: 12px; color: #909399; }
.disk-free { margin-top: 4px; font-size: 13px; color: #67c23a; font-weight: 500; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; align-items: center; }
</style>
