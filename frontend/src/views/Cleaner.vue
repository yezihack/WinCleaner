<template>
  <div class="cleaner">
    <h2 style="margin-bottom: 20px; color: #303133;">垃圾清理</h2>

    <el-card shadow="hover">
      <div class="action-bar">
        <el-button type="primary" :loading="scanning" @click="handleScan">
          <el-icon><Search /></el-icon> 扫描垃圾
        </el-button>
        <el-button
          type="danger"
          :loading="cleaning"
          :disabled="selectedCategories.length === 0"
          @click="handleClean"
        >
          <el-icon><Delete /></el-icon> 清理选中 ({{ formatBytes(selectedSize) }})
        </el-button>
      </div>

      <el-table
        v-if="results.length > 0"
        :data="results"
        style="margin-top: 16px;"
        row-key="category"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-files">
              <div class="expand-header">
                <span>{{ row.category }} - 共 {{ row.count }} 个文件</span>
                <el-input
                  v-model="filterText[row.category]"
                  placeholder="搜索文件名..."
                  clearable
                  size="small"
                  style="width: 240px;"
                  :prefix-icon="Search"
                />
              </div>
              <el-table :data="getFilteredItems(row)" size="small" max-height="360" stripe>
                <el-table-column label="文件路径" min-width="400">
                  <template #default="{ row: item }">
                    <span class="file-path">{{ item.path }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="大小" width="100" align="right">
                  <template #default="{ row: item }">
                    {{ formatBytes(item.size) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column type="selection" width="50" />
        <el-table-column label="分类">
          <template #default="{ row }">
            <span class="category-name">{{ row.category }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="count" label="文件数" width="100" />
        <el-table-column label="大小" width="120">
          <template #default="{ row }">{{ formatBytes(row.size) }}</template>
        </el-table-column>
      </el-table>

      <el-empty v-else-if="!scanning" description="点击扫描开始检测垃圾文件" />
    </el-card>

    <el-card v-if="cleanResult" shadow="hover" style="margin-top: 16px;">
      <template #header><span>清理结果</span></template>
      <el-result icon="success" :title="`释放空间: ${formatBytes(cleanResult.freed_size)}`">
        <template #sub-title>
          成功清理 {{ cleanResult.cleaned_count }} 个文件，
          {{ cleanResult.failed_count }} 个跳过（无权限）
        </template>
      </el-result>
    </el-card>

    <!-- 清理历史 -->
    <el-card shadow="hover" style="margin-top: 16px;">
      <template #header>
        <div class="history-header">
          <span>清理历史</span>
          <el-button text size="small" @click="loadHistory">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </template>

      <div v-if="history" class="history-section">
        <!-- 统计概览 -->
        <el-row :gutter="16" class="stat-cards">
          <el-col :span="8">
            <div class="stat-card">
              <div class="stat-value">{{ history.last_clean_ago || '从未清理' }}</div>
              <div class="stat-label">距上次清理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="stat-card">
              <div class="stat-value">{{ formatBytes(history.total_freed) }}</div>
              <div class="stat-label">累计释放空间</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="stat-card">
              <div class="stat-value">{{ history.total_count }}</div>
              <div class="stat-label">累计清理文件数</div>
            </div>
          </el-col>
        </el-row>

        <!-- 图表切换 -->
        <el-radio-group v-model="chartMode" style="margin: 16px 0;">
          <el-radio-button value="daily">近30天</el-radio-button>
          <el-radio-button value="monthly">按月统计</el-radio-button>
        </el-radio-group>

        <!-- ECharts 图表 -->
        <v-chart :option="chartOption" style="height: 300px;" autoresize />
      </div>

      <el-empty v-else description="暂无清理记录" :image-size="60" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
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

const selectedCategories = computed(() => selectedRows.value.map(r => r.category))
const selectedSize = computed(() => selectedRows.value.reduce((sum, r) => sum + r.size, 0))

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
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const p0 = params[0]
        const p1 = params[1]
        return `${p0.axisValue}<br/>
          ${p0.marker} 释放: ${p0.value} MB<br/>
          ${p1.marker} 文件数: ${p1.value}`
      }
    },
    legend: { data: ['释放空间 (MB)', '清理文件数'] },
    grid: { left: 50, right: 50, bottom: 40, top: 40 },
    xAxis: {
      type: 'category',
      data: xLabels,
      axisLabel: { rotate: isDaily ? 45 : 0, fontSize: 11 }
    },
    yAxis: [
      { type: 'value', name: 'MB', position: 'left' },
      { type: 'value', name: '文件数', position: 'right' }
    ],
    series: [
      {
        name: '释放空间 (MB)',
        type: 'bar',
        data: freedData,
        itemStyle: { color: '#409eff', borderRadius: [4, 4, 0, 0] },
        yAxisIndex: 0
      },
      {
        name: '清理文件数',
        type: 'line',
        data: countData,
        itemStyle: { color: '#67c23a' },
        smooth: true,
        yAxisIndex: 1
      }
    ]
  }
})

const handleScan = async () => {
  scanning.value = true
  cleanResult.value = null
  try {
    results.value = await api.scanJunk()
    ElMessage.success(`扫描完成，发现 ${results.value.length} 个分类`)
  } catch (e) {
    ElMessage.error('扫描失败')
  } finally {
    scanning.value = false
  }
}

const handleClean = async () => {
  try {
    await ElMessageBox.confirm(
      `确定清理选中的 ${selectedCategories.value.length} 个分类？`,
      '确认清理',
      { type: 'warning' }
    )
  } catch { return }

  cleaning.value = true
  try {
    cleanResult.value = await api.cleanJunk(selectedCategories.value)
    ElMessage.success('清理完成')
    results.value = await api.scanJunk()
    await loadHistory()
  } catch (e) {
    ElMessage.error('清理失败')
  } finally {
    cleaning.value = false
  }
}

const handleSelectionChange = (rows: ScanResult[]) => {
  selectedRows.value = rows
}

const loadHistory = async () => {
  try {
    const data = await api.getCleanHistory()
    if (data && (data.total_count > 0 || (data.records && data.records.length > 0))) {
      history.value = data
    } else {
      history.value = null
    }
  } catch {
    history.value = null
  }
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.action-bar {
  display: flex;
  gap: 12px;
}
.category-name {
  color: #409eff;
  cursor: pointer;
}
.expand-files {
  padding: 8px 16px;
}
.expand-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
  color: #606266;
}
.file-path {
  font-size: 12px;
  color: #606266;
  word-break: break-all;
}
.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.stat-cards {
  margin-bottom: 8px;
}
.stat-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}
.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}
.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
